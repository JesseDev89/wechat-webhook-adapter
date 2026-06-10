#!/bin/bash
set -euo pipefail

# ============================================================================
# 构建并推送企业微信 Webhook 服务镜像
# ============================================================================
# 前置条件:
#   1. Docker 已安装并可运行
#   2. 可访问镜像仓库 (如没有，可以先 save 到本地再导入 K8s)
#
# 用法:
#   # 方式一: 推送到远程镜像仓库
#   REGISTRY=registry.example.com/jesse ./build.sh
#
#   # 方式二: 仅构建本地镜像
#   ./build.sh
#
#   # 方式三: 保存为 tar 包 (用于离线环境)
#   ./build.sh save
# ============================================================================

# ==================== 可配置变量 ====================
REGISTRY="${REGISTRY:-registry.cn-guangzhou.aliyuncs.com}"                          # 镜像仓库地址，如 registry.example.com
IMAGE_NAME="jesse-dnmp/wechat-webhook-adapter"
IMAGE_TAG="${IMAGE_TAG:-v1.0.0}"
DOCKERFILE_DIR="$(cd "$(dirname "$0")" && pwd)"

# ==================== 颜色输出 ====================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'
#!/bin/bash
set -euo pipefail

# ============================================================================
# 构建并推送企业微信 Webhook 服务镜像
# ============================================================================
# 前置条件:
#   1. Docker 已安装并可运行
#   2. 可访问镜像仓库 (如没有，可以先 save 到本地再导入 K8s)
#
# 用法:
#   # 方式一: 推送到远程镜像仓库
#   REGISTRY=registry.example.com/jesse ./build.sh
#
#   # 方式二: 仅构建本地镜像
#   ./build.sh
#
#   # 方式三: 保存为 tar 包 (用于离线环境)
#   ./build.sh save
# ============================================================================

# ==================== 可配置变量 ====================
REGISTRY="${REGISTRY:-registry.cn-guangzhou.aliyuncs.com}"                          # 镜像仓库地址，如 registry.example.com
IMAGE_NAME="jesse-dnmp/wechat-webhook-adapter"
IMAGE_TAG="${IMAGE_TAG:-v1.0.0}"
DOCKERFILE_DIR="$(cd "$(dirname "$0")" && pwd)"

# ==================== 颜色输出 ====================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

# ==================== 构建镜像 ====================
build_image() {
    log_info "开始构建镜像..."
    log_info "目录: ${DOCKERFILE_DIR}"
    log_info "镜像: ${IMAGE_NAME}:${IMAGE_TAG}"

    cd "${DOCKERFILE_DIR}"
    docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" .

    log_info "镜像构建完成"
}

# ==================== 推送镜像 ====================
push_image() {
    if [ -z "${REGISTRY}" ]; then
        log_warn "REGISTRY 未设置，跳过推送"
        return
    fi

    local full_image="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
    log_info "推送镜像到: ${full_image}"

    docker tag "${IMAGE_NAME}:${IMAGE_TAG}" "${full_image}"
    docker push "${full_image}"

    log_info "镜像推送完成"
}

# ==================== 保存为 tar 包 ====================
save_image() {
    local tar_file="${IMAGE_NAME}-${IMAGE_TAG}.tar"
    log_info "保存镜像为 tar 包: ${tar_file}"

    docker save "${IMAGE_NAME}:${IMAGE_TAG}" -o "${tar_file}"

    log_info "镜像保存完成: ${tar_file}"
    log_info "导入命令: docker load -i ${tar_file}"
}

# ==================== 显示帮助 ====================
show_help() {
    echo "用法: $0 [save]"
    echo ""
    echo "环境变量:"
    echo "  REGISTRY    镜像仓库地址 (可选)"
    echo "  IMAGE_TAG   镜像标签 (默认: v1.0.0)"
    echo ""
    echo "示例:"
    echo "  # 仅构建本地镜像"
    echo "  $0"
    echo ""
    echo "  # 构建并推送到远程仓库"
    echo "  REGISTRY=registry.example.com/jesse $0"
    echo ""
    echo "  # 保存为 tar 包 (离线环境)"
    echo "  $0 save"
}

# ==================== 主流程 ====================
main() {
    case "${1:-}" in
        help|--help|-h)
            show_help
            exit 0
            ;;
        save)
            build_image
            save_image
            ;;
        *)
            build_image
            push_image
            ;;
    esac
}

main "$@"

log_info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }

# ==================== 构建镜像 ====================
build_image() {
    log_info "开始构建镜像..."
    log_info "目录: ${DOCKERFILE_DIR}"
    log_info "镜像: ${IMAGE_NAME}:${IMAGE_TAG}"

    cd "${DOCKERFILE_DIR}"
    docker build -t "${IMAGE_NAME}:${IMAGE_TAG}" .

    log_info "镜像构建完成"
}

# ==================== 推送镜像 ====================
push_image() {
    if [ -z "${REGISTRY}" ]; then
        log_warn "REGISTRY 未设置，跳过推送"
        return
    fi

    local full_image="${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
    log_info "推送镜像到: ${full_image}"

    docker tag "${IMAGE_NAME}:${IMAGE_TAG}" "${full_image}"
    docker push "${full_image}"

    log_info "镜像推送完成"
}

# ==================== 保存为 tar 包 ====================
save_image() {
    local tar_file="${IMAGE_NAME}-${IMAGE_TAG}.tar"
    log_info "保存镜像为 tar 包: ${tar_file}"

    docker save "${IMAGE_NAME}:${IMAGE_TAG}" -o "${tar_file}"

    log_info "镜像保存完成: ${tar_file}"
    log_info "导入命令: docker load -i ${tar_file}"
}

# ==================== 显示帮助 ====================
show_help() {
    echo "用法: $0 [save]"
    echo ""
    echo "环境变量:"
    echo "  REGISTRY    镜像仓库地址 (可选)"
    echo "  IMAGE_TAG   镜像标签 (默认: v1.0.0)"
    echo ""
    echo "示例:"
    echo "  # 仅构建本地镜像"
    echo "  $0"
    echo ""
    echo "  # 构建并推送到远程仓库"
    echo "  REGISTRY=registry.example.com/jesse $0"
    echo ""
    echo "  # 保存为 tar 包 (离线环境)"
    echo "  $0 save"
}

# ==================== 主流程 ====================
main() {
    case "${1:-}" in
        help|--help|-h)
            show_help
            exit 0
            ;;
        save)
            build_image
            save_image
            ;;
        *)
            build_image
            push_image
            ;;
    esac
}

main "$@"
