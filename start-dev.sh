#!/bin/bash

# EMS 开发环境快速启动脚本

set -e

echo "========================================="
echo "  EMS 开发环境启动"
echo "========================================="

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

    echo -e "${RED}错误: 未找到 Docker Compose，请先安装 Docker Desktop 或 docker-compose${NC}"
    exit 1
}

wait_for_postgres() {
    local attempts=30

    while [ $attempts -gt 0 ]; do
        if "${COMPOSE_CMD[@]}" exec -T postgres pg_isready -U ems -d ems_db >/dev/null 2>&1; then
            return 0
        fi
        sleep 2
        attempts=$((attempts - 1))
    done

    echo -e "${RED}错误: PostgreSQL 启动超时${NC}"
    exit 1
}

# 进入项目根目录
cd "$(dirname "$0")"

# 检查 Node.js / npm
if ! command -v node >/dev/null 2>&1; then
    echo -e "${RED}错误: Node.js 未安装${NC}"
    exit 1
fi

if ! command -v npm >/dev/null 2>&1; then
    echo -e "${RED}错误: npm 未安装${NC}"
    exit 1
fi

# 检查 Go
if ! command -v go >/dev/null 2>&1; then
    echo -e "${RED}错误: Go 未安装，无法启动开发后端${NC}"
    exit 1
fi

resolve_compose_cmd

echo -e "${BLUE}1. 启动数据库与缓存服务...${NC}"
"${COMPOSE_CMD[@]}" up -d postgres redis
echo "等待 PostgreSQL 就绪..."
wait_for_postgres
echo -e "${GREEN}✓ PostgreSQL / Redis 已就绪${NC}"

echo -e "${BLUE}2. 停止容器化前后端，释放本地开发端口...${NC}"
"${COMPOSE_CMD[@]}" stop backend frontend >/dev/null 2>&1 || true
echo -e "${GREEN}✓ 容器化前后端已停止${NC}"

echo -e "${BLUE}3. 启动后端服务...${NC}"
cd backend
go run main.go > /tmp/ems-backend.log 2>&1 &
BACKEND_PID=$!
echo "后端 PID: $BACKEND_PID"
cd ..

echo -e "${BLUE}4. 启动前端服务...${NC}"
cd frontend
npm run dev > /tmp/ems-frontend.log 2>&1 &
FRONTEND_PID=$!
echo "前端 PID: $FRONTEND_PID"
cd ..

echo ""
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}  开发环境已启动！${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""
echo "访问地址："
echo -e "  前端: ${BLUE}http://localhost:5173${NC}"
echo -e "  后端: ${BLUE}http://localhost:8080${NC}"
echo ""
echo "停止服务："
echo "  kill $BACKEND_PID $FRONTEND_PID"
echo ""
echo "查看日志："
echo "  tail -f /tmp/ems-backend.log"
echo "  tail -f /tmp/ems-frontend.log"
echo ""

# 保存 PID 到文件
echo "$BACKEND_PID" > /tmp/ems-backend.pid
echo "$FRONTEND_PID" > /tmp/ems-frontend.pid
