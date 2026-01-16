# EMS 设备管理系统

> 企业级设备管理平台，支持点检、维修、保养、备件管理、统计分析等功能。

## 📄 许可证与授权

**Copyright © 2025 Feng-H**

本项目采用 **GNU General Public License v3.0 (GPL-3.0)** 许可证开源。

### ⚠️ 重要声明

**要求开源**：
- 如果你使用、修改或分发本项目的代码，你的项目也必须以 **GPL-3.0** 协议开源
- 你必须开源所有基于本项目的修改和衍生作品
- 你需要提供源代码，并保留原始版权声明

**商业使用需授权**：
- 如果你不希望开源你的修改，你需要获得项目作者的**书面授权**
- 联系方式：通过 GitHub Issues 联系 [Feng-H](https://github.com/Feng-H)

### 为什么选择 GPL-3.0？

- ✅ **保护开源**：确保基于本项目的所有改进也必须开源
- ✅ **防止私有化**：禁止将本代码集成到闭源商业产品中
- ✅ **社区贡献**：鼓励社区共同改进，所有改进回馈社区
- ✅ **作者权益**：保护原作者的劳动成果

详见 [LICENSE](./LICENSE) 文件。

---

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
