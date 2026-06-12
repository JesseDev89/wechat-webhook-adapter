// Package handler 提供 HTTP 请求处理功能
package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"jesse.com/wechat-webhook-adapter/internal/alert"
	"jesse.com/wechat-webhook-adapter/internal/config"
	"jesse.com/wechat-webhook-adapter/internal/types"
	"jesse.com/wechat-webhook-adapter/internal/wechat"
)

// Handler HTTP 请求处理器
type Handler struct {
	config    *config.Config
	wechat    *wechat.Client
	formatter *alert.Formatter

	mu               sync.Mutex
	buffer           []types.Alert
	timer            *time.Timer
	debounceDuration time.Duration
}

// New 创建处理器
func New(cfg *config.Config) *Handler {
	return &Handler{
		config:           cfg,
		wechat:           wechat.NewClient(cfg.WebhookKey),
		formatter:        alert.NewFormatter(),
		debounceDuration: 5 * time.Second,
	}
}

// RegisterRoutes 注册 HTTP 路由
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/webhook", h.webhook)
	mux.HandleFunc("/health", h.health)
}

// health 健康检查
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		log.Printf("写入响应失败: %v\n", err)
	}
}

// webhook 接收 Alertmanager 告警
func (h *Handler) webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/webhook" {
		http.NotFound(w, r)
		return
	}

	var payload types.AlertmanagerPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Json解析失败: %v", err)
		http.Error(w, `{"status":"invalid json"}`, http.StatusBadRequest)
		return
	}

	// 打印原始请求内容
	requestBody, _ := json.MarshalIndent(payload, "", "  ")
	log.Printf("接收到的请求内容: %s", string(requestBody))

	log.Printf("接收到告警信息: status=%s, alerts=%d", payload.Status, len(payload.Alerts))

	if err := h.CanSend(); err != nil {
		log.Printf("发送失败: %v", err)
		http.Error(w, `{"status":"failed"}`, http.StatusInternalServerError)
		return
	}

	h.bufferAlerts(payload.Alerts)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(`{"status":"success"}`)); err != nil {
		log.Printf("写入响应失败: %v", err)
	}
}

// bufferAlerts 将告警缓冲并设置 debounce 定时器
func (h *Handler) bufferAlerts(alerts []types.Alert) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.buffer = append(h.buffer, alerts...)

	if h.timer != nil {
		h.timer.Reset(h.debounceDuration)
	} else {
		h.timer = time.AfterFunc(h.debounceDuration, h.flush)
	}
}

// flush 将缓冲的告警合并发送
func (h *Handler) flush() {
	h.mu.Lock()
	alerts := h.buffer
	h.buffer = nil
	h.timer = nil
	h.mu.Unlock()

	if len(alerts) == 0 {
		return
	}

	// 推断整体状态：有 firing 则为 firing
	status := "resolved"
	for _, a := range alerts {
		if a.Status == "firing" {
			status = "firing"
			break
		}
	}

	message := h.formatter.Format(status, alerts, "")

	if err := h.wechat.Send(message); err != nil {
		log.Printf("发送失败: %v", err)
	}
}

// Shutdown 优雅关闭，确保缓冲的告警被发送
func (h *Handler) Shutdown() {
	h.mu.Lock()
	if h.timer != nil {
		h.timer.Stop()
		h.timer = nil
	}
	h.mu.Unlock()

	h.flush()
}

// CanSend 检查是否可以发送消息（用于预检查）
func (h *Handler) CanSend() error {
	if h.config.WebhookKey == "" {
		return fmt.Errorf("WECHAT_WEBHOOK_KEY 未设置")
	}
	return nil
}
