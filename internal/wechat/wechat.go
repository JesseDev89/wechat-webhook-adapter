// Package wechat 提供企业微信机器人消息发送功能
package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"jesse.com/wechat-webhook-adapter/internal/types"
)

// WebhookURL 企业微信机器人 Webhook 地址
const WebhookURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"

// Client 企业微信机器人客户端
type Client struct {
	key string
}

// NewClient 创建微信客户端
func NewClient(key string) *Client {
	return &Client{key: key}
}

// Send 发送消息到企业微信
func (c *Client) Send(content string) error {
	if c.key == "" {
		return fmt.Errorf("WECHAT_WEBHOOK_KEY 未设置")
	}

	url := fmt.Sprintf("%s?key=%s", WebhookURL, c.key)
	msg := types.WeChatMessage{
		MsgType: "markdown",
		Markdown: types.MarkdownContent{
			Content: content,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json encode failed: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("http post failed: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response failed: %w", err)
	}

	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		return fmt.Errorf("wechat api error: %v", result)
	}

	log.Println("消息发送成功")
	return nil
}
