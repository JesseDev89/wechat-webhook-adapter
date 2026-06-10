package wechat

import (
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-key")
	if client == nil {
		t.Fatal("NewClient 不应返回 nil")
	}
	if client.key != "test-key" {
		t.Errorf("期望 key='test-key', 实际='%s'", client.key)
	}
}

func TestClient_Send_WithoutKey(t *testing.T) {
	client := NewClient("")
	err := client.Send("测试消息")
	if err == nil {
		t.Error("期望返回错误，实际 nil")
	}
	if !strings.Contains(err.Error(), "WECHAT_WEBHOOK_KEY") {
		t.Errorf("错误信息应包含 'WECHAT_WEBHOOK_KEY', 实际: %v", err)
	}
}
