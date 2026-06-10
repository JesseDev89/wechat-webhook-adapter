// Package handler 提供 HTTP 请求处理功能
package handler

import (
	"encoding/json"
	"log"
	"net/http"

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
}

// New 创建处理器
func New(cfg *config.Config) *Handler {
	return &Handler{
		config:    cfg,
		wechat:    wechat.NewClient(cfg.WebhookKey),
		formatter: alert.NewFormatter(),
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

	log.Printf("接收到告警信息: message %+v, status=%s, alerts=%d", payload, payload.Status, len(payload.Alerts))

	message := h.formatter.Format(&payload)

	if err := h.wechat.Send(message); err != nil {
		log.Printf("发送失败: %v", err)
		http.Error(w, `{"status":"failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(`{"status":"success"}`)); err != nil {
		log.Printf("写入响应失败: %v", err)
	}
}
