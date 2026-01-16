# EMS 设备管理系统

企业级设备管理平台，支持点检、维修、保养、备件管理、统计分析等功能。

## 技术栈

- **后端**: Go 1.23 + Gin + GORM + PostgreSQL + Redis
- **前端**: Vue 3 + TypeScript + Vite + Element Plus
- **移动端**: Vue 3 + Vant 4 (H5)
- **部署**: Docker + Docker Compose

## 项目结构

```
EMS-Claude/
├── backend/          # Go后端服务
├── frontend/         # Vue3前端应用
├── db/              # 数据库脚本
├── docker/          # Docker配置
└── docs/            # 项目文档
```

## 快速开始

### 前置要求

- Go 1.23+
- Node.js 20+
- PostgreSQL 15+
- Redis 7+
- Docker (可选)

### 本地开发

#### 1. 启动数据库服务

```bash
cd docker
docker-compose up postgres redis
```

#### 2. 启动后端服务

```bash
cd backend
go mod download
go run main.go
```

后端服务将在 http://localhost:8080 启动

API文档: http://localhost:8080/swagger/index.html

#### 3. 启动前端服务

```bash
cd frontend
npm install
npm run dev
```

前端服务将在 http://localhost:5173 启动

### Docker 部署

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

## 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | password123 | 管理员 |

## 开发文档

详细开发文档请参阅 [CLAUDE.md](./CLAUDE.md)

## License

MIT
