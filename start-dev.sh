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

# 进入项目根目录
cd "$(dirname "$0")"

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo -e "${RED}错误: Node.js 未安装${NC}"
    exit 1
fi

# 检查 Go
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}警告: Go 未安装，后端将无法启动${NC}"
fi

echo -e "${BLUE}1. 启动后端服务...${NC}"
if command -v go &> /dev/null; then
    cd backend
    go run main.go > /tmp/ems-backend.log 2>&1 &
    BACKEND_PID=$!
    echo "后端 PID: $BACKEND_PID"
    cd ..
fi

echo -e "${BLUE}2. 启动前端服务...${NC}"
cd frontend
/opt/homebrew/Cellar/node@22/22.16.0/bin/npm run dev > /tmp/ems-frontend.log 2>&1 &
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
