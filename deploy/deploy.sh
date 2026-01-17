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

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: Docker 未安装${NC}"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}错误: Docker Compose 未安装${NC}"
    exit 1
fi

# 进入部署目录
cd "$(dirname "$0")/.."

echo -e "${YELLOW}1. 停止现有服务...${NC}"
docker-compose -f deploy/docker-compose.yml down

echo -e "${YELLOW}2. 清理旧镜像...${NC}"
docker-compose -f deploy/docker-compose.yml build --no-cache

echo -e "${YELLOW}3. 启动服务...${NC}"
docker-compose -f deploy/docker-compose.yml up -d

echo -e "${YELLOW}4. 等待服务启动...${NC}"
sleep 10

echo -e "${YELLOW}5. 检查服务状态...${NC}"
docker-compose -f deploy/docker-compose.yml ps

echo ""
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}  部署完成！${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""
echo "访问地址："
echo "  前端: http://localhost"
echo "  后端: http://localhost:8080"
echo ""
echo "查看日志："
echo "  docker-compose -f deploy/docker-compose.yml logs -f"
echo ""
