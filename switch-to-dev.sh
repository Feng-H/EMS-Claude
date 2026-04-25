#!/bin/bash

# 从生产环境切换回开发环境

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

resolve_compose_cmd() {
    if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
        COMPOSE_CMD=(docker compose)
        return
    fi

    if command -v docker-compose >/dev/null 2>&1; then
        COMPOSE_CMD=(docker-compose)
        return
    fi

    echo -e "${RED}错误: 未找到 Docker Compose${NC}"
    exit 1
}

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}  切换到开发环境${NC}"
echo -e "${BLUE}=========================================${NC}"
echo ""

# 停止生产环境
echo -e "${YELLOW}[1/2] 停止生产环境...${NC}"
cd "$(dirname "$0")"
resolve_compose_cmd
"${COMPOSE_CMD[@]}" stop backend frontend >/dev/null 2>&1 || true
echo -e "${GREEN}✓ 生产环境已停止${NC}"

# 启动开发环境
echo ""
echo -e "${YELLOW}[2/2] 启动开发环境...${NC}"
./start-dev.sh

echo ""
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}  ✓ 已切换到开发环境${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""
echo -e "开发环境访问地址："
echo -e "  🌐 前端: ${BLUE}http://localhost:5173${NC}"
echo -e "  🔧 后端: ${BLUE}http://localhost:8080${NC}"
echo ""
echo -e "${YELLOW}提示：修改代码后，运行以下命令同步到生产环境：${NC}"
echo -e "  ${GREEN}./sync-to-prod.sh${NC}"
echo ""
