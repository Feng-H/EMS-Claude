#!/bin/bash

# EMS 快速同步到生产环境
# 将开发环境的修改同步到生产环境

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}  同步开发环境 → 生产环境${NC}"
echo -e "${BLUE}=========================================${NC}"
echo ""

# 进入项目根目录
cd "$(dirname "$0")"

# 检查是否有修改
if [ -n "$(git status --porcelain 2>/dev/null)" ]; then
    echo -e "${YELLOW}⚠️  检测到未提交的修改${NC}"
    echo "建议先提交代码，或者按 Ctrl+C 取消"
    echo ""
    read -p "继续同步？(y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${RED}已取消${NC}"
        exit 1
    fi
fi

# 步骤 1: 停止开发环境
echo -e "${YELLOW}[1/4] 停止开发环境...${NC}"
if [ -f stop-dev.sh ]; then
    ./stop-dev.sh > /dev/null 2>&1 || true
fi
pkill -f "vite" 2>/dev/null || true
echo -e "${GREEN}✓ 开发环境已停止${NC}"

# 步骤 2: 构建前端
echo ""
echo -e "${YELLOW}[2/4] 构建前端生产环境...${NC}"
cd frontend
/opt/homebrew/Cellar/node@22/22.16.0/bin/npm run build
cd ..
echo -e "${GREEN}✓ 前端构建完成${NC}"

# 步骤 3: 检查 Docker
echo ""
echo -e "${YELLOW}[3/4] 检查 Docker 服务...${NC}"
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}错误: Docker 未运行${NC}"
    echo "请先启动 Docker Desktop"
    exit 1
fi
echo -e "${GREEN}✓ Docker 运行正常${NC}"

# 步骤 4: 部署生产环境
echo ""
echo -e "${YELLOW}[4/4] 部署到生产环境...${NC}"

# 停止现有容器
docker-compose -f deploy/docker-compose.yml down 2>/dev/null || true

# 构建并启动
docker-compose -f deploy/docker-compose.yml up -d --build

# 等待服务启动
echo ""
echo -e "${YELLOW}等待服务启动...${NC}"
sleep 10

# 检查服务状态
echo ""
echo -e "${YELLOW}检查服务状态...${NC}"
docker-compose -f deploy/docker-compose.yml ps

echo ""
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}  ✓ 同步完成！${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""
echo -e "生产环境访问地址："
echo -e "  🌐 前端: ${BLUE}http://localhost${NC}"
echo -e "  🔧 后端: ${BLUE}http://localhost:8080${NC}"
echo ""
echo -e "查看日志："
echo -e "  docker-compose -f deploy/docker-compose.yml logs -f"
echo ""
echo -e "${YELLOW}提示：如需返回开发环境，运行：${NC}"
echo -e "  ${GREEN}./start-dev.sh${NC}"
echo ""
