# wechat-webhook-adapter

Prometheus Alertmanager 企业微信机器人 Webhook 适配器，将 Alertmanager 告警以 Markdown 格式推送到企业微信群。

---

## 功能特性

- **接收告警**：接收 Alertmanager 推送的 webhook 告警
- **消息格式化**：将告警信息自动格式化为易读的 Markdown 消息
- **同类告警聚合**：相同告警名称 + 状态的告警自动合并展示，避免刷屏
- **智能告警对象**：按 node > daemonset > instance > job 优先级自动识别告警对象
- **分级图标**：根据告警严重程度（critical/warning/info）显示不同颜色图标
- **容器化部署**：提供 Dockerfile 和 Kubernetes 部署配置
- **健康检查**：内置 `/health` 健康检查端点

---

## 技术栈

- **语言**：Go 1.25
- **构建**：多阶段 Docker 构建，静态编译，无依赖
- **部署**：Docker / Kubernetes

---

## 项目结构

```
.
├── main.go                          # 程序入口
├── internal/
│   ├── alert/alert.go               # 告警格式化（Markdown 转换）
│   ├── config/config.go             # 配置管理（环境变量）
│   ├── handler/handler.go           # HTTP 路由处理器
│   ├── types/types.go               # 数据结构定义
│   └── wechat/wechat.go             # 企业微信机器人消息发送
├── Dockerfile                       # 多阶段构建镜像
├── build.sh                         # 镜像构建脚本
├── wechat-webhook-k8s.yaml          # Kubernetes 部署配置
├── alertmanager-config.yaml         # Alertmanager 配置示例
├── go.mod                           # Go 模块定义
└── go.sum                           # 依赖校验
```

---

## 架构设计

```
┌─────────────────┐     ┌─────────────────────────┐     ┌─────────────────┐
│  Alertmanager   │────▶│  wechat-webhook-adapter │────▶│  企业微信机器人  │
│  (Webhook 推送) │     │  接收 /webhook          │     │  (Markdown 消息)│
└─────────────────┘     └─────────────────────────┘     └─────────────────┘
                                  │
                                  ▼
                         ┌────────────────────┐
                         │  格式化告警消息   │
                         │  - 告警名称       │
                         │  - 严重程度       │
                         │  - 告警对象       │
                         │  - 命名空间       │
                         │  - 描述信息       │
                         │  - 同类告警聚合   │
                         └────────────────────┘
```

---

## 快速开始

### 1. 获取企业微信机器人 Key

1. 打开企业微信群聊
2. 进入【群设置】→【添加群机器人】
3. 复制 Webhook URL 中的 `key` 参数

### 2. 本地运行

```bash
# 设置环境变量
export WECHAT_WEBHOOK_KEY="your-robot-key"
export SERVER_PORT="80"

# 运行
go run main.go
```

### 3. Docker 运行

```bash
# 构建镜像
docker build -t wechat-webhook-adapter .

# 运行容器
docker run -d \
  -p 80:80 \
  -e WECHAT_WEBHOOK_KEY="your-robot-key" \
  registry.cn-guangzhou.aliyuncs.com/jessedev/wechat-webhook-adapter:v1.0.0
```

### 4. Kubernetes 部署

```bash
# 修改 wechat-webhook-k8s.yaml 中的 WECHAT_WEBHOOK_KEY
kubectl apply -f wechat-webhook-k8s.yaml
```

---

## 配置说明

### 环境变量

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `WECHAT_WEBHOOK_KEY` | 企业微信机器人 Key | *(必填)* |
| `SERVER_PORT` | HTTP 服务端口 | `80` |

### HTTP 端点

| 端点 | 方法 | 说明 |
|------|------|------|
| `/webhook` | POST | 接收 Alertmanager 告警 |
| `/health` | GET | 健康检查 |

---

## 消息格式

企业微信机器人收到的消息格式示例：

### 单条告警

```
🚨 **Prometheus 告警触发**
> 时间: `2026-06-12 12:33:11`

---
🟡 **KubeNodeNotReady** (firing)
> 严重程度: `warning`
> 告警对象: `k8s-node-02`
> 命名空间: `monitoring`
> 描述: k8s-node-02 has been unready for more than 15 minutes on cluster .

[查看详情](http://alertmanager.example.com)
```

### 同类告警聚合

```
✅ **Prometheus 告警恢复**
> 时间: `2026-06-12 12:12:56`

---
🟡 **KubeDaemonSetRolloutStuck** (resolved) [2 条]
> 严重程度: `warning`
> 告警对象:
> - calico-node (kube-system)
> - kube-proxy (kube-system)
> 描述: DaemonSet kube-system/calico-node has not finished or progressed for at least 15m.
```

---

## 相关文档

- [USAGE.md](USAGE.md) — 详细使用文档
- [alertmanager-config.yaml](alertmanager-config.yaml) — Alertmanager 配置示例
- [wechat-webhook-k8s.yaml](wechat-webhook-k8s.yaml) — Kubernetes 部署配置

---

## License

MIT
