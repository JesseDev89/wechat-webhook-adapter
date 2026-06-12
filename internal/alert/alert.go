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
	container   string
	endpoint    string
	pod         string
	service     string
	severity    string
	alertname   string
	condition   string
	effect      string
	key         string
	status      string
	alertLabels map[string]string
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
		container:   alert.Labels["container"],
		endpoint:    alert.Labels["endpoint"],
		pod:         alert.Labels["pod"],
		service:     alert.Labels["service"],
		severity:    alert.Labels["severity"],
		alertname:   alert.Labels["alertname"],
		condition:   alert.Labels["condition"],
		effect:      alert.Labels["effect"],
		key:         alert.Labels["key"],
		status:      alert.Labels["status"],
		alertLabels: alert.Labels,
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
	f.addAlertDetails(sb, item)
	sb.WriteString(fmt.Sprintf("> 描述: %s\n", displayDesc(item)))
}

// formatMultiAlert 格式化多条聚合告警
func (f *Formatter) formatMultiAlert(sb *bytes.Buffer, group alertGroup, icon, alertname string) {
	sb.WriteString("\n---\n")
	sb.WriteString(fmt.Sprintf("%s **%s** (%s) [%d 条]\n", icon, alertname, group.status, len(group.items)))
	sb.WriteString(fmt.Sprintf("> 严重程度: `%s`\n", group.severity))
	sb.WriteString("> 告警详情:\n")
	for i, item := range group.items {
		sb.WriteString(fmt.Sprintf("> \n> --- 第 %d 条 ---\n", i+1))
		f.addAlertDetails(sb, item)
	}
	sb.WriteString(fmt.Sprintf("> 描述: %s\n", displayDesc(group.items[0])))
}

// addAlertDetails 添加告警详情信息
func (f *Formatter) addAlertDetails(sb *bytes.Buffer, item alertItem) {
	if item.node != "" {
		sb.WriteString(fmt.Sprintf("> 节点: `%s`\n", item.node))
	}
	if item.instance != "" {
		sb.WriteString(fmt.Sprintf("> 实例: `%s`\n", item.instance))
	}
	if item.job != "" {
		sb.WriteString(fmt.Sprintf("> 任务: `%s`\n", item.job))
	}
	if item.container != "" {
		sb.WriteString(fmt.Sprintf("> 容器: `%s`\n", item.container))
	}
	if item.pod != "" {
		sb.WriteString(fmt.Sprintf("> Pod: `%s`\n", item.pod))
	}
	if item.service != "" {
		sb.WriteString(fmt.Sprintf("> 服务: `%s`\n", item.service))
	}
	if item.condition != "" {
		sb.WriteString(fmt.Sprintf("> 条件: `%s`\n", item.condition))
	}
	if item.effect != "" {
		sb.WriteString(fmt.Sprintf("> 影响: `%s`\n", item.effect))
	}
	if item.namespace != "" {
		sb.WriteString(fmt.Sprintf("> 命名空间: `%s`\n", item.namespace))
	}
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
