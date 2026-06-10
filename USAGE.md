 # 使用文档

## 目录

- [前置条件](#前置条件)
- [获取企业微信机器人 Key](#获取企业微信机器人-key)
- [本地开发](#本地开发)
- [Docker 部署](#docker-部署)
- [Kubernetes 部署](#kubernetes-部署)
- [Alertmanager 配置](#alertmanager-配置)
- [测试验证](#测试验证)
- [常见问题](#常见问题)

---

## 前置条件

- Go 1.25+（本地开发）
- Docker（容器化部署）
- kubectl + Kubernetes 集群（K8s 部署）
- 企业微信群组（用于创建群机器人）

---

## 获取企业微信机器人 Key

1. 打开企业微信，进入目标群聊
2. 点击右上角【群设置】（三个点图标）
3. 选择【添加群机器人】
4. 点击【新建机器人】，设置名称
5. 复制 Webhook 地址，提取其中的 `key` 参数

```
https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
                                              └─── 复制这串 key ────┘
```

---

## 本地开发

### 1. 克隆项目

```bash
git clone <repository-url>
cd wechat-webhook-adapter
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 运行服务

```bash
export WECHAT_WEBHOOK_KEY="your-robot-key"
export SERVER_PORT="80"

go run main.go
```

服务启动后输出：

```
Webhook 服务已启动: http://0.0.0.0:80/webhook
企业微信 Key: 已设置
```

### 4. 验证服务

```bash
curl http://localhost:80/health
# 预期输出: {"status":"ok"}
```

---

## Docker 部署

### 方式一：直接构建运行

```bash
# 构建镜像
docker build -t wechat-webhook-adapter:v1.0.0 .

# 运行容器
docker run -d \
  --name wechat-webhook \
  -p 80:80 \
  -e WECHAT_WEBHOOK_KEY="your-robot-key" \
  -e SERVER_PORT="80" \
  --restart unless-stopped \
  wechat-webhook-adapter:v1.0.0
```

### 方式二：使用 build.sh 脚本

```bash
# 仅构建本地镜像
./build.sh

# 构建并推送到远程镜像仓库
REGISTRY=registry.example.com/jesse ./build.sh

# 保存为 tar 包（离线环境）
./build.sh save
```

### 查看容器日志

```bash
docker logs -f wechat-webhook
```

---

## Kubernetes 部署

### 1. 修改配置

编辑 `wechat-webhook-k8s.yaml`，修改以下两项：

```yaml
# 修改为你的镜像地址
image: your-registry.com/jesse/wechat-webhook-adapter:v1.0.0

# 修改为你的企业微信机器人 Key
- name: WECHAT_WEBHOOK_KEY
  value: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
```

### 2. 部署到集群

```bash
kubectl apply -f wechat-webhook-k8s.yaml
```

### 3. 检查部署状态

```bash
# 查看 Pod 状态
kubectl get pods -n monitoring -l app=wechat-webhook

# 查看服务状态
kubectl get svc -n monitoring wechat-webhook

# 查看日志
kubectl logs -n monitoring -l app=wechat-webhook --tail=50
```

### 4. 验证服务可用性

```bash
# 端口转发到本地
kubectl port-forward -n monitoring svc/wechat-webhook 80:80

# 测试健康检查
curl http://localhost:80/health
```

---

## Alertmanager 配置

项目已提供完整的 Alertmanager 配置示例，参考 `alertmanager-config.yaml`。

### 核心配置说明

#### 路由配置

```yaml
route:
  receiver: 'wechat-default'
  group_by: ['alertname', 'severity', 'namespace']
  group_wait: 30s        # 组等待时间
  group_interval: 5m     # 组发送间隔
  repeat_interval: 4h      # 重复告警间隔
```

#### Webhook 接收器

```yaml
receivers:
  - name: 'wechat-default'
    webhook_configs:
      - url: 'http://wechat-webhook:80/webhook'
        send_resolved: true  # 告警恢复时也发送通知
```

### 应用配置

将配置更新到 Alertmanager Secret：

```bash
kubectl create secret generic alertmanager-prometheus-kube-prometheus-alertmanager \
  --from-file=alertmanager.yaml=alertmanager-config.yaml \
  -n monitoring --dry-run=client -o yaml | kubectl apply -f -
```

等待 Alertmanager Pod 自动重启后生效。

---

## 测试验证

### 模拟 Alertmanager 推送

```bash
curl -X POST http://localhost:80/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "status": "firing",
    "alerts": [
      {
        "status": "firing",
        "labels": {
          "alertname": "HighMemoryUsage",
          "severity": "critical",
          "instance": "web-server-01"
        },
        "annotations": {
          "summary": "内存使用率超过 90%"
        }
      }
    ],
    "external_url": "http://alertmanager.example.com"
  }'
```

预期响应：

```json
{"status":"success"}
```

此时企业微信群应收到告警消息。

### 测试告警恢复

```bash
curl -X POST http://localhost:80/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "status": "resolved",
    "alerts": [
      {
        "status": "resolved",
        "labels": {
          "alertname": "HighMemoryUsage",
          "severity": "critical",
          "instance": "web-server-01"
        },
        "annotations": {
          "summary": "内存使用率已恢复正常"
        }
      }
    ],
    "external_url": "http://alertmanager.example.com"
  }'
```

---

## 常见问题

### Q1: 服务启动后收不到告警消息？

检查清单：

- [ ] `WECHAT_WEBHOOK_KEY` 是否正确设置
- [ ] Alertmanager webhook URL 是否正确指向服务地址
- [ ] 网络是否可达（服务与 Alertmanager 之间）
- [ ] 查看服务日志是否有报错信息

### Q2: 告警消息格式不正确？

确保 Alertmanager 推送的数据结构符合标准格式。服务会解析以下字段：

- `status`: 告警状态（firing/resolved）
- `alerts[].labels.alertname`: 告警名称
- `alerts[].labels.severity`: 严重程度
- `alerts[].labels.instance`: 实例信息
- `alerts[].annotations.summary`: 描述信息

### Q3: 如何在 K8s 中调试？

```bash
# 进入 Pod 内部
kubectl exec -it -n monitoring deployment/wechat-webhook -- sh

# 查看服务日志
kubectl logs -n monitoring -l app=wechat-webhook -f

# 测试 webhook 连通性
kubectl run -it --rm debug --image=busybox --restart=Never -- sh
wget -qO- http://wechat-webhook:80/health
```

### Q4: 如何更新机器人 Key？

修改环境变量后重启服务即可：

```bash
# Docker
docker restart wechat-webhook

# Kubernetes
kubectl rollout restart deployment/wechat-webhook -n monitoring
```
