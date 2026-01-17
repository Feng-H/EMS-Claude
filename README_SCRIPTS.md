# EMS 设备管理系统

## 🎯 快速开始

### 开发环境（推荐）

```bash
# 一键启动
./start-dev.sh

# 访问地址
http://localhost:5173
```

### 生产环境

```bash
# 同步并部署
./sync-to-prod.sh

# 访问地址
http://localhost
```

## 📖 文档导航

| 文档 | 说明 |
|------|------|
| [WORKFLOW.md](./WORKFLOW.md) | **开发工作流程**（必读） |
| [CLAUDE.md](./CLAUDE.md) | 项目总览和技术架构 |
| [deploy/README.md](./deploy/README.md) | 生产环境部署指南 |

## 🛠️ 快捷命令

| 命令 | 说明 |
|------|------|
| `./start-dev.sh` | 启动开发环境 |
| `./stop-dev.sh` | 停止开发环境 |
| `./sync-to-prod.sh` | 同步代码到生产环境 |
| `./switch-to-dev.sh` | 切换回开发环境 |

## 🌟 新功能

### 1. 人员管理

- ✅ 用户账号管理（新增、编辑、启用/禁用）
- ✅ 账号申请审核流程
- ✅ 角色权限管理

### 2. 账号申请

- ✅ 登录页面支持账号申请
- ✅ 管理员审核机制
- ✅ 申请状态追踪

### 3. 密码管理

- ✅ 首次登录强制修改密码
- ✅ 默认密码：`password123`
- ✅ 密码修改页面

## 🔑 默认账号

- **用户名**: `admin`
- **密码**: `password123`

⚠️ **重要**: 首次登录后请立即修改密码！

## 📂 项目结构

```
EMS-Claude/
├── backend/           # Go 后端服务
├── frontend/          # Vue3 前端应用
├── deploy/            # 生产环境配置
│   ├── nginx.conf
│   ├── Dockerfile.*
│   └── docker-compose.yml
├── db/                # 数据库脚本
│   ├── schema.sql
│   └── migrations/
├── start-dev.sh       # 启动开发环境
├── stop-dev.sh        # 停止开发环境
├── sync-to-prod.sh    # 同步到生产环境
└── switch-to-dev.sh   # 切换到开发环境
```

## 🔧 技术栈

### 前端
- Vue 3 + TypeScript
- Vite
- Element Plus
- Pinia (状态管理)
- Vue Router

### 后端
- Go 1.23
- Gin (Web 框架)
- GORM (ORM)
- PostgreSQL (数据库)
- Redis (缓存)

### 部署
- Docker
- Docker Compose
- Nginx

## 📞 支持

如有问题，请查看：
1. [开发工作流程](./WORKFLOW.md)
2. [项目总览](./CLAUDE.md)
3. [部署指南](./deploy/README.md)

## 📄 许可证

GPL-3.0
