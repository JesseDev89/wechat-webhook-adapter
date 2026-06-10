# ==================== 构建阶段 ====================
FROM golang:1.25-alpine AS builder

WORKDIR /www/

# 先复制依赖文件，利用 Docker 缓存层
COPY go.mod go.sum* ./
RUN go mod download

# 再复制源代码
COPY main.go .
COPY internal/ ./internal/

# 编译（静态链接，去除调试信息，减小体积）
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o wechat-webhook-adapter .

# ==================== 运行阶段 ====================
FROM alpine:3.20

# 创建非 root 用户组和用户（GID/UID 均为 1000）
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

WORKDIR /www

# 复制编译后的二进制文件
COPY --from=builder /www/wechat-webhook-adapter .

# 切换到非 root 用户
USER appuser

# 暴露端口（使用非特权端口，非 root 用户也能运行）
EXPOSE 80

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80/health || exit 1

ENTRYPOINT ["./wechat-webhook-adapter"]
