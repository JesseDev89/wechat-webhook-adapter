package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"jesse.com/wechat-webhook-adapter/internal/config"
	"jesse.com/wechat-webhook-adapter/internal/types"
)

func TestHandler_Health(t *testing.T) {
	cfg := &config.Config{}
	h := New(cfg)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("期望状态码 %d, 实际 %d", http.StatusOK, rec.Code)
	}

	body, _ := io.ReadAll(rec.Body)
	if !strings.Contains(string(body), "ok") {
		t.Errorf("期望响应包含 'ok', 实际: %s", string(body))
	}
}

func TestHandler_Webhook_InvalidMethod(t *testing.T) {
	cfg := &config.Config{}
	h := New(cfg)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/webhook", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("期望状态码 %d, 实际 %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestHandler_Webhook_InvalidJSON(t *testing.T) {
	cfg := &config.Config{}
	h := New(cfg)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("期望状态码 %d, 实际 %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandler_Webhook_WithoutWebhookKey(t *testing.T) {
	cfg := &config.Config{WebhookKey: ""}
	h := New(cfg)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	payload := types.AlertmanagerPayload{
		Status: "firing",
		Alerts: []types.Alert{
			{
				Status: "firing",
				Labels: map[string]string{
					"alertname": "TestAlert",
					"severity":  "critical",
					"instance":  "localhost:9090",
				},
				Annotations: map[string]string{
					"summary": "测试告警",
				},
			},
		},
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("期望状态码 %d, 实际 %d", http.StatusInternalServerError, rec.Code)
	}
}
