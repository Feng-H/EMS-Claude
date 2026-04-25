#!/bin/bash

# EMS 生产环境部署脚本

set -e

echo "========================================="
echo "  EMS 生产环境部署"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

    echo -e "${RED}错误: Docker Compose 未安装${NC}"
    exit 1
}

if ! command -v docker >/dev/null 2>&1; then
    echo -e "${RED}错误: Docker 未安装${NC}"
    exit 1
fi

# 进入部署目录
cd "$(dirname "$0")/.."
resolve_compose_cmd

echo -e "${YELLOW}1. 停止现有服务...${NC}"
"${COMPOSE_CMD[@]}" stop backend frontend >/dev/null 2>&1 || true

echo -e "${YELLOW}2. 清理旧镜像...${NC}"
"${COMPOSE_CMD[@]}" build --no-cache backend frontend

echo -e "${YELLOW}3. 启动服务...${NC}"
"${COMPOSE_CMD[@]}" up -d postgres redis backend frontend

echo -e "${YELLOW}4. 等待服务启动...${NC}"
sleep 10

echo -e "${YELLOW}5. 检查服务状态...${NC}"
"${COMPOSE_CMD[@]}" ps

echo ""
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}  部署完成！${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""
echo "访问地址："
echo "  前端: http://127.0.0.1:3000"
echo "  后端: http://127.0.0.1:9000"
echo ""
echo "查看日志："
echo "  ${COMPOSE_CMD[*]} logs -f"
echo ""
