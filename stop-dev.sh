#!/bin/bash

# EMS 开发环境停止脚本

echo "========================================="
echo "  停止 EMS 开发环境"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 从 PID 文件读取进程 ID
if [ -f /tmp/ems-backend.pid ]; then
    BACKEND_PID=$(cat /tmp/ems-backend.pid)
    if ps -p $BACKEND_PID > /dev/null 2>&1; then
        echo "停止后端服务 (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
        echo -e "${GREEN}✓ 后端服务已停止${NC}"
    else
        echo -e "${YELLOW}后端服务未运行${NC}"
    fi
    rm /tmp/ems-backend.pid
fi

if [ -f /tmp/ems-frontend.pid ]; then
    FRONTEND_PID=$(cat /tmp/ems-frontend.pid)
    if ps -p $FRONTEND_PID > /dev/null 2>&1; then
        echo "停止前端服务 (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID
        echo -e "${GREEN}✓ 前端服务已停止${NC}"
    else
        echo -e "${YELLOW}前端服务未运行${NC}"
    fi
    rm /tmp/ems-frontend.pid
fi

# 清理可能残留的进程
pkill -f "go run main.go" 2>/dev/null && echo -e "${GREEN}✓ 已清理残留的后端进程${NC}"
pkill -f "vite" 2>/dev/null && echo -e "${GREEN}✓ 已清理残留的前端进程${NC}"

echo ""
echo -e "${GREEN}所有服务已停止${NC}"
