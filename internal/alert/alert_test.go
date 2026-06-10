package alert

import (
	"strings"
	"testing"

	"jesse.com/wechat-webhook-adapter/internal/types"
)

func TestFormatter_Format(t *testing.T) {
	formatter := NewFormatter()
	payload := &types.AlertmanagerPayload{
		Status: "firing",
		Alerts: []types.Alert{
			{
				Status: "firing",
				Labels: map[string]string{
					"alertname": "HighCPU",
					"severity":  "critical",
					"instance":  "server-01",
				},
				Annotations: map[string]string{
					"summary": "CPU 使用率过高",
				},
			},
		},
		ExternalURL: "http://prometheus.example.com",
	}

	result := formatter.Format(payload)

	checks := []string{
		"告警触发",
		"HighCPU",
		"critical",
		"server-01",
		"CPU 使用率过高",
		"prometheus.example.com",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("格式化结果应包含 '%s', 实际结果:\n%s", check, result)
		}
	}
}

func TestFormatter_Format_Resolved(t *testing.T) {
	formatter := NewFormatter()
	payload := &types.AlertmanagerPayload{
		Status: "resolved",
		Alerts: []types.Alert{
			{
				Status: "resolved",
				Labels: map[string]string{
					"alertname": "TestAlert",
					"severity":  "warning",
				},
				Annotations: map[string]string{},
			},
		},
	}

	result := formatter.Format(payload)

	if !strings.Contains(result, "告警恢复") {
		t.Error("期望包含 '告警恢复' 字样")
	}
}

func TestFormatter_Format_EmptyFields(t *testing.T) {
	formatter := NewFormatter()
	payload := &types.AlertmanagerPayload{
		Status: "firing",
		Alerts: []types.Alert{
			{
				Status: "firing",
				Labels: map[string]string{
					"severity": "info",
				},
				Annotations: map[string]string{},
			},
		},
	}

	result := formatter.Format(payload)

	if !strings.Contains(result, "Unknown") {
		t.Error("期望 alertname 为空时显示 'Unknown'")
	}
	if !strings.Contains(result, "unknown") {
		t.Error("期望 instance 为空时显示 'unknown'")
	}
	if !strings.Contains(result, "无描述信息") {
		t.Error("期望 summary 为空时显示 '无描述信息'")
	}
}

func TestFormatter_Format_MultipleAlerts(t *testing.T) {
	formatter := NewFormatter()
	payload := &types.AlertmanagerPayload{
		Status: "firing",
		Alerts: []types.Alert{
			{
				Status: "firing",
				Labels: map[string]string{
					"alertname": "Alert1",
					"severity":  "critical",
					"instance":  "host1",
				},
				Annotations: map[string]string{"summary": "问题1"},
			},
			{
				Status: "firing",
				Labels: map[string]string{
					"alertname": "Alert2",
					"severity":  "warning",
					"instance":  "host2",
				},
				Annotations: map[string]string{"summary": "问题2"},
			},
		},
	}

	result := formatter.Format(payload)

	if !strings.Contains(result, "---") {
		t.Error("多条告警应包含分隔符 '---'")
	}

	if !strings.Contains(result, "Alert1") || !strings.Contains(result, "Alert2") {
		t.Error("应包含所有告警的名称")
	}
}
