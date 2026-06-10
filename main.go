package main

import (
	"log"
	"net/http"

	"jesse.com/wechat-webhook-adapter/internal/config"
	"jesse.com/wechat-webhook-adapter/internal/handler"
)

func main() {
	cfg := config.New()

	if cfg.WebhookKey == "" {
		log.Println("警告: WECHAT_WEBHOOK_KEY 未设置，消息将无法发送")
	}

	mux := http.NewServeMux()
	h := handler.New(cfg)
	h.RegisterRoutes(mux)

	addr := ":" + cfg.Port

	log.Printf("Webhook 服务已启动: http://0.0.0.0%s/webhook", addr)
	log.Printf("企业微信 Key: %s", map[bool]string{true: "已设置", false: "未设置"}[cfg.WebhookKey != ""])

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
