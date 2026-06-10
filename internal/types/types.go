// Package types 定义项目中的所有数据结构
package types

// WeChatMessage 企业微信机器人消息结构体
type WeChatMessage struct {
	MsgType  string          `json:"msgtype"`
	Markdown MarkdownContent `json:"markdown"`
}

// MarkdownContent Markdown 内容
type MarkdownContent struct {
	Content string `json:"content"`
}

// AlertmanagerPayload Alertmanager 推送的告警负载
type AlertmanagerPayload struct {
	Status      string  `json:"status"`
	Alerts      []Alert `json:"alerts"`
	ExternalURL string  `json:"external_url"`
}

// Alert 单条告警信息
type Alert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}
