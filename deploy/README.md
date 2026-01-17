# EMS 生产环境部署指南

## 概述

本指南介绍如何将 EMS 设备管理系统部署到生产环境。

## 前置要求

- Docker 20.10+
- Docker Compose 2.0+
- 至少 2GB 可用内存
- 至少 10GB 可用磁盘空间

## 快速部署

### 1. 一键部署

```bash
# 给部署脚本添加执行权限
chmod +x deploy/deploy.sh

# 执行部署
./deploy/deploy.sh
```

### 2. 手动部署

```bash
# 构建并启动所有服务
docker-compose -f deploy/docker-compose.yml up -d

# 查看服务状态
docker-compose -f deploy/docker-compose.yml ps

# 查看日志
docker-compose -f deploy/docker-compose.yml logs -f
```

## 服务说明

部署后包含以下服务：

| 服务 | 端口 | 说明 |
|------|------|------|
| 前端 | 80 | Vue3 构建的生产环境前端 |
| 后端 | 8080 | Go API 服务 |
| PostgreSQL | 5432 | 数据库 |
| Redis | 6379 | 缓存服务 |

## 访问地址

- **前端**: http://localhost
- **后端 API**: http://localhost:8080
- **API 文档**: http://localhost:8080/swagger/index.html

## 默认账号

- **管理员**: admin / password123

⚠️ **重要**: 首次登录后请立即修改默认密码！

## 数据库迁移

如果是从旧版本升级，需要运行数据库迁移：

```bash
# 方式1: 使用 Docker 执行迁移
docker-compose -f deploy/docker-compose.yml exec postgres psql -U ems -d ems_db -f /docker-entrypoint-initdb.d/01-schema.sql

# 方式2: 直接执行迁移文件
psql -h localhost -U ems -d ems_db -f db/migrations/001_add_user_approval.sql
```

## 常用命令

### 查看服务状态

```bash
docker-compose -f deploy/docker-compose.yml ps
```

### 查看日志

```bash
# 查看所有服务日志
docker-compose -f deploy/docker-compose.yml logs -f

# 查看特定服务日志
docker-compose -f deploy/docker-compose.yml logs -f frontend
docker-compose -f deploy/docker-compose.yml logs -f backend
docker-compose -f deploy/docker-compose.yml logs -f postgres
```

### 重启服务

```bash
# 重启所有服务
docker-compose -f deploy/docker-compose.yml restart

# 重启特定服务
docker-compose -f deploy/docker-compose.yml restart backend
docker-compose -f deploy/docker-compose.yml restart frontend
```

### 停止服务

```bash
# 停止所有服务
docker-compose -f deploy/docker-compose.yml down

# 停止并删除数据卷
docker-compose -f deploy/docker-compose.yml down -v
```

### 更新服务

```bash
# 重新构建并启动
docker-compose -f deploy/docker-compose.yml up -d --build

# 仅更新前端
docker-compose -f deploy/docker-compose.yml up -d --build frontend

# 仅更新后端
docker-compose -f deploy/docker-compose.yml up -d --build backend
```

## 数据备份

### 备份数据库

```bash
# 创建备份目录
mkdir -p backups

# 备份数据库
docker-compose -f deploy/docker-compose.yml exec postgres pg_dump -U ems ems_db > backups/ems_backup_$(date +%Y%m%d_%H%M%S).sql
```

### 恢复数据库

```bash
# 从备份恢复
docker-compose -f deploy/docker-compose.yml exec -T postgres psql -U ems ems_db < backups/ems_backup_20250117_100000.sql
```

## 监控和维护

### 健康检查

```bash
# 检查所有服务健康状态
docker-compose -f deploy/docker-compose.yml ps
```

### 清理未使用的资源

```bash
# 清理未使用的镜像
docker image prune -a

# 清理未使用的容器
docker container prune

# 清理未使用的卷
docker volume prune
```

### 查看资源使用

```bash
# 查看 Docker 资源使用情况
docker stats
```

## 故障排查

### 前端无法访问

1. 检查前端服务状态：
   ```bash
   docker-compose -f deploy/docker-compose.yml logs frontend
   ```

2. 检查端口占用：
   ```bash
   lsof -i :80
   ```

### 后端无法访问

1. 检查后端服务状态：
   ```bash
   docker-compose -f deploy/docker-compose.yml logs backend
   ```

2. 检查数据库连接：
   ```bash
   docker-compose -f deploy/docker-compose.yml logs postgres
   ```

### 数据库连接失败

1. 检查 PostgreSQL 是否启动：
   ```bash
   docker-compose -f deploy/docker-compose.yml ps postgres
   ```

2. 测试数据库连接：
   ```bash
   docker-compose -f deploy/docker-compose.yml exec postgres psql -U ems -d ems_db -c "SELECT 1"
   ```

## 安全建议

1. **修改默认密码**: 首次部署后立即修改所有默认密码
2. **配置防火墙**: 限制对数据库端口的访问
3. **使用 HTTPS**: 生产环境建议配置 SSL 证书
4. **定期备份**: 设置自动备份任务
5. **更新镜像**: 定期更新 Docker 镜像以获取安全补丁

## 性能优化

1. **调整资源限制**: 在 docker-compose.yml 中添加资源限制
2. **启用缓存**: 确保 Redis 正常运行
3. **数据库优化**: 根据负载调整 PostgreSQL 配置
4. **负载均衡**: 使用 nginx 配置负载均衡

## 支持

如有问题，请查看：
- 项目文档: /CLAUDE.md
- API 文档: http://localhost:8080/swagger/index.html
- 日志文件: `docker-compose logs`
