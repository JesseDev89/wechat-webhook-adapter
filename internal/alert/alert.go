// Package alert 提供告警信息的格式化功能
package alert

import (
	"bytes"
	"fmt"
	"time"

	"jesse.com/wechat-webhook-adapter/internal/types"
)

// Formatter 告警格式化器
type Formatter struct{}

// NewFormatter 创建格式化器
func NewFormatter() *Formatter {
	return &Formatter{}
}

// Format 格式化告警信息为 Markdown
func (f *Formatter) Format(payload *types.AlertmanagerPayload) string {
	var header string
	if payload.Status == "firing" {
		header = "🚨 **Prometheus 告警触发**\n"
	} else {
		header = "✅ **Prometheus 告警恢复**\n"
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	var sb bytes.Buffer
	sb.WriteString(header)
	sb.WriteString(fmt.Sprintf("> 时间: `%s`\n", now))

	for _, alert := range payload.Alerts {
		f.formatSingleAlert(&sb, &alert)
	}

	if payload.ExternalURL != "" {
		sb.WriteString(fmt.Sprintf("\n[查看详情](%s)\n", payload.ExternalURL))
	}

	return sb.String()
}

// formatSingleAlert 格式化单条告警
func (f *Formatter) formatSingleAlert(sb *bytes.Buffer, alert *types.Alert) {
	severity := alert.Labels["severity"]
	alertname := alert.Labels["alertname"]
	if alertname == "" {
		alertname = "Unknown"
	}
	instance := alert.Labels["instance"]
	if instance == "" {
		instance = "unknown"
	}
	summary := alert.Annotations["summary"]
	if summary == "" {
		summary = "无描述信息"
	}

	icon := f.severityIcon(severity)

	sb.WriteString("\n---\n")
	sb.WriteString(fmt.Sprintf("%s **%s** (%s)\n", icon, alertname, alert.Status))
	sb.WriteString(fmt.Sprintf("> 严重程度: `%s`\n", severity))
	sb.WriteString(fmt.Sprintf("> 实例: `%s`\n", instance))
	sb.WriteString(fmt.Sprintf("> 描述: %s\n", summary))
}

// severityIcon 根据严重程度返回对应的 emoji
func (f *Formatter) severityIcon(severity string) string {
	switch severity {
	case "critical":
		return "🔴"
	case "warning":
		return "🟡"
	case "info":
		return "🔵"
	default:
		return "⚪"
	}
}
