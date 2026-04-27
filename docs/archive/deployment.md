# EMS 设备管理系统 - 部署文档

## 目录

1. [系统要求](#系统要求)
2. [开发环境搭建](#开发环境搭建)
3. [Docker 部署](#docker-部署)
4. [生产环境部署](#生产环境部署)
5. [数据库优化](#数据库优化)
6. [监控与维护](#监控与维护)
7. [故障排查](#故障排查)

---

## 系统要求

### 最低配置

| 组件 | 最低配置 | 推荐配置 |
|------|---------|---------|
| CPU | 2核 | 4核+ |
| 内存 | 4GB | 8GB+ |
| 磁盘 | 20GB | 50GB+ SSD |
| 操作系统 | Linux/Windows/macOS | Linux (Ubuntu 20.04+) |

### 软件依赖

- **后端**: Go 1.23+, PostgreSQL 15+, Redis 7+
- **前端**: Node.js 20+, npm/yarn/pnpm
- **容器**: Docker 20+, Docker Compose 2+

---

## 开发环境搭建

### 1. 克隆项目

```bash
git clone https://github.com/your-org/EMS-Claude.git
cd EMS-Claude
```

### 2. 后端开发环境

```bash
cd backend

# 安装依赖
go mod download

# 复制配置文件
cp config/config.yaml.example config/config.yaml

# 编辑配置
vim config/config.yaml
```

配置文件说明：

```yaml
server:
  port: 8080        # API 服务端口
  mode: debug       # debug/release/test

database:
  host: localhost
  port: 5432
  user: ems
  password: ems123
  dbname: ems_db

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-secret-key-change-in-production  # 生产环境必须修改
  expire_hours: 24
```

运行后端：

```bash
# 运行迁移
go run main.go

# 或使用 air 进行热重载
air
```

### 3. 前端开发环境

```bash
cd frontend

# 安装依赖
npm install

# 开发模式运行
npm run dev

# 构建生产版本
npm run build
```

### 4. 数据库初始化

```bash
# 创建数据库
createdb -U postgres ems_db

# 运行初始脚本
psql -U postgres -d ems_db -f db/schema.sql

# (可选) 导入种子数据
psql -U postgres -d ems_db -f db/seeds/data.sql

# 创建索引优化
psql -U postgres -d ems_db -f backend/scripts/migrations/add_indexes.sql
```

---

## Docker 部署

### 快速启动

```bash
# 在项目根目录
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

### 服务说明

| 服务 | 端口 | 说明 |
|------|------|------|
| frontend | 80 | Web 前端 |
| backend | 8080 | API 服务 |
| postgres | 5432 | PostgreSQL 数据库 |
| redis | 6379 | Redis 缓存 |
| nginx | 8081 | 反向代理 (可选) |

### 构建镜像

```bash
# 构建后端镜像
docker build -t ems-backend:latest ./backend

# 构建前端镜像
docker build -t ems-frontend:latest ./frontend

# 使用自定义镜像启动
docker-compose up -d
```

### 生产环境配置

创建 `docker-compose.prod.yml`：

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: always

  backend:
    image: ems-backend:latest
    environment:
      - GIN_MODE=release
    restart: always
    depends_on:
      - postgres
      - redis

  frontend:
    image: ems-frontend:latest
    restart: always
    depends_on:
      - backend
```

使用环境变量文件 `.env`：

```env
POSTGRES_DB=ems_db
POSTGRES_USER=ems
POSTGRES_PASSWORD=strong-password-here
```

启动：

```bash
docker-compose -f docker-compose.prod.yml --env-file .env up -d
```

---

## 生产环境部署

### 1. 使用 Nginx 反向代理

安装 Nginx：

```bash
sudo apt update
sudo apt install nginx
```

配置站点 `/etc/nginx/sites-available/ems`：

```nginx
upstream backend {
    server localhost:8080;
}

upstream frontend {
    server localhost:3000;
}

server {
    listen 80;
    server_name ems.example.com;

    # 前端静态文件
    location / {
        proxy_pass http://frontend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # API 代理
    location /api/ {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 文件上传大小限制
    client_max_body_size 10M;
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/ems /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 2. HTTPS 配置 (Let's Encrypt)

```bash
# 安装 certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d ems.example.com

# 自动续期
sudo certbot renew --dry-run
```

### 3. 使用 Systemd 管理服务

后端服务 `/etc/systemd/system/ems-backend.service`：

```ini
[Unit]
Description=EMS Backend API
After=network.target postgresql.service

[Service]
Type=simple
User=ems
WorkingDirectory=/opt/ems/backend
ExecStart=/opt/ems/backend/ems-api
Restart=always
RestartSec=5
Environment="GIN_MODE=release"

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable ems-backend
sudo systemctl start ems-backend
```

---

## 数据库优化

### 1. 连接池配置

```yaml
# config/config.yaml
database:
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
```

### 2. 索引优化

运行索引迁移脚本：

```bash
psql -U ems -d ems_db -f backend/scripts/migrations/add_indexes.sql
```

### 3. 启用查询缓存

```yaml
# config/config.yaml
redis:
  enabled: true
  host: localhost
  port: 6379
```

### 4. 定期维护

```bash
# 定期清理死元组
psql -U ems -d ems_db -c "VACUUM ANALYZE;"

# 重建索引
REINDEX DATABASE CONCURRENTLY ems_db;
```

### 5. 备份策略

```bash
# 每日备份脚本
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/postgres"
mkdir -p $BACKUP_DIR

pg_dump -U ems -d ems_db -Fc -f $BACKUP_DIR/ems_$DATE.dump

# 保留最近 7 天的备份
find $BACKUP_DIR -name "ems_*.dump" -mtime +7 -delete
```

---

## 监控与维护

### 1. 健康检查端点

- `GET /health` - 服务健康状态
- `GET /api/v1/auth/me` - 认证状态

### 2. 日志管理

```bash
# 查看后端日志
tail -f backend/logs/ems.log

# 查看 Docker 日志
docker-compose logs -f backend

# 查看 Nginx 日志
tail -f /var/log/nginx/access.log
```

### 3. 性能监控

使用 Prometheus + Grafana 监控：

```yaml
# docker-compose.monitoring.yml
services:
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    ports:
      - "3001:3000"
```

---

## 故障排查

### 常见问题

**1. 数据库连接失败**

```bash
# 检查 PostgreSQL 状态
sudo systemctl status postgresql

# 检查连接
psql -U ems -d ems_db -h localhost

# 检查防火墙
sudo ufw status
```

**2. Redis 连接失败**

```bash
# 检查 Redis 状态
redis-cli ping

# 检查配置
redis-cli config get requirepass
```

**3. 前端构建失败**

```bash
# 清理缓存
rm -rf node_modules package-lock.json
npm install

# 检查 Node 版本
node -v  # 应该是 v20+
```

**4. Docker 容器启动失败**

```bash
# 查看容器日志
docker-compose logs backend

# 重建容器
docker-compose up -d --force-recreate

# 清理并重启
docker-compose down -v
docker-compose up -d
```

---

## API 文档

启动服务后访问：

- Swagger UI: `http://localhost:8080/swagger/index.html`

---

## 安全建议

1. **修改默认密码**: 更改数据库、Redis、JWT secret
2. **启用 HTTPS**: 生产环境必须使用 SSL
3. **限制访问**: 使用防火墙限制数据库访问
4. **定期更新**: 保持依赖包最新
5. **备份**: 设置自动备份策略
6. **监控**: 配置日志和性能监控
