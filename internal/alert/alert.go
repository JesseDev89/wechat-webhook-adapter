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

// alertItem 提取的告警关键信息
type alertItem struct {
	node        string
	daemonset   string
	instance    string
	job         string
	namespace   string
	summary     string
	description string
}

// alertGroup 聚合后的告警组
type alertGroup struct {
	alertname string
	status    string
	severity  string
	items     []alertItem
}

// Format 格式化告警信息为 Markdown
func (f *Formatter) Format(status string, alerts []types.Alert, externalURL string) string {
	var header string
	if status == "firing" {
		header = "🚨 **Prometheus 告警触发**\n"
	} else {
		header = "✅ **Prometheus 告警恢复**\n"
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	var sb bytes.Buffer
	sb.WriteString(header)
	sb.WriteString(fmt.Sprintf("> 时间: `%s`\n", now))

	groups := f.groupAlerts(alerts)
	for _, group := range groups {
		f.formatAlertGroup(&sb, group)
	}

	if externalURL != "" {
		sb.WriteString(fmt.Sprintf("\n[查看详情](%s)\n", externalURL))
	}

	return sb.String()
}

// groupAlerts 按 alertname + status 聚合告警
func (f *Formatter) groupAlerts(alerts []types.Alert) []alertGroup {
	groupMap := make(map[string]*alertGroup)
	var order []string

	for i := range alerts {
		alert := &alerts[i]
		key := alert.Labels["alertname"] + "|" + alert.Status
		if g, ok := groupMap[key]; ok {
			g.items = append(g.items, f.extractAlertItem(alert))
		} else {
			groupMap[key] = &alertGroup{
				alertname: alert.Labels["alertname"],
				status:    alert.Status,
				severity:  alert.Labels["severity"],
				items:     []alertItem{f.extractAlertItem(alert)},
			}
			order = append(order, key)
		}
	}

	result := make([]alertGroup, 0, len(order))
	for _, key := range order {
		result = append(result, *groupMap[key])
	}
	return result
}

func (f *Formatter) extractAlertItem(alert *types.Alert) alertItem {
	return alertItem{
		node:        alert.Labels["node"],
		daemonset:   alert.Labels["daemonset"],
		instance:    alert.Labels["instance"],
		job:         alert.Labels["job"],
		namespace:   alert.Labels["namespace"],
		summary:     alert.Annotations["summary"],
		description: alert.Annotations["description"],
	}
}

// displayName 获取告警对象的展示名称
// 优先级: node > daemonset > instance > job > unknown
func displayName(item alertItem) string {
	switch {
	case item.node != "":
		return item.node
	case item.daemonset != "":
		return item.daemonset
	case item.instance != "":
		return item.instance
	case item.job != "":
		return item.job
	default:
		return "unknown"
	}
}

// displayDesc 获取告警描述，优先使用 description
func displayDesc(item alertItem) string {
	if item.description != "" {
		return item.description
	}
	if item.summary != "" {
		return item.summary
	}
	return "无描述信息"
}

// formatAlertGroup 格式化聚合后的告警组
func (f *Formatter) formatAlertGroup(sb *bytes.Buffer, group alertGroup) {
	alertname := group.alertname
	if alertname == "" {
		alertname = "Unknown"
	}

	icon := f.severityIcon(group.severity)

	if len(group.items) == 1 {
		f.formatSingleAlert(sb, group, icon, alertname)
	} else {
		f.formatMultiAlert(sb, group, icon, alertname)
	}
}

// formatSingleAlert 格式化单条告警
func (f *Formatter) formatSingleAlert(sb *bytes.Buffer, group alertGroup, icon, alertname string) {
	item := group.items[0]

	sb.WriteString("\n---\n")
	sb.WriteString(fmt.Sprintf("%s **%s** (%s)\n", icon, alertname, group.status))
	sb.WriteString(fmt.Sprintf("> 严重程度: `%s`\n", group.severity))
	sb.WriteString(fmt.Sprintf("> 告警对象: `%s`\n", displayName(item)))
	if item.namespace != "" {
		sb.WriteString(fmt.Sprintf("> 命名空间: `%s`\n", item.namespace))
	}
	sb.WriteString(fmt.Sprintf("> 描述: %s\n", displayDesc(item)))
}

// formatMultiAlert 格式化多条聚合告警
func (f *Formatter) formatMultiAlert(sb *bytes.Buffer, group alertGroup, icon, alertname string) {
	sb.WriteString("\n---\n")
	sb.WriteString(fmt.Sprintf("%s **%s** (%s) [%d 条]\n", icon, alertname, group.status, len(group.items)))
	sb.WriteString(fmt.Sprintf("> 严重程度: `%s`\n", group.severity))
	sb.WriteString("> 告警对象:\n")
	for _, item := range group.items {
		name := displayName(item)
		if item.namespace != "" {
			name = fmt.Sprintf("%s (%s)", name, item.namespace)
		}
		sb.WriteString(fmt.Sprintf("> - %s\n", name))
	}
	sb.WriteString(fmt.Sprintf("> 描述: %s\n", displayDesc(group.items[0])))
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
