// Package config 负责读取和管理应用配置
package config

import "os"

// Config 应用配置
type Config struct {
	WebhookKey string
	Port       string
}

// New 从环境变量读取配置
func New() *Config {
	return &Config{
		WebhookKey: getEnv("WECHAT_WEBHOOK_KEY", ""),
		Port:       getEnv("SERVER_PORT", "80"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
