package config

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// 保存原始值
	origKey := os.Getenv("WECHAT_WEBHOOK_KEY")
	origPort := os.Getenv("SERVER_PORT")
	defer func() {
		os.Setenv("WECHAT_WEBHOOK_KEY", origKey)
		os.Setenv("SERVER_PORT", origPort)
	}()

	os.Setenv("WECHAT_WEBHOOK_KEY", "test-key")
	os.Setenv("SERVER_PORT", "9090")

	cfg := New()
	if cfg.WebhookKey != "test-key" {
		t.Errorf("期望 WebhookKey='test-key', 实际='%s'", cfg.WebhookKey)
	}
	if cfg.Port != "9090" {
		t.Errorf("期望 Port='9090', 实际='%s'", cfg.Port)
	}
}

func TestNew_Defaults(t *testing.T) {
	os.Unsetenv("WECHAT_WEBHOOK_KEY")
	os.Unsetenv("SERVER_PORT")

	cfg := New()
	if cfg.WebhookKey != "" {
		t.Errorf("期望 WebhookKey='', 实际='%s'", cfg.WebhookKey)
	}
	if cfg.Port != "80" {
		t.Errorf("期望 Port='80', 实际='%s'", cfg.Port)
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	if v := getEnv("TEST_KEY", "default"); v != "test_value" {
		t.Errorf("期望 'test_value', 实际='%s'", v)
	}

	if v := getEnv("NON_EXISTENT_KEY", "default"); v != "default" {
		t.Errorf("期望 'default', 实际='%s'", v)
	}
}
