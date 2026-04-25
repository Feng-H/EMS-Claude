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
- **UI 风格**: 现代化工业风，支持暗色模式，针对智能运维优化的交互设计
- **智能大脑**: 基于 LLM 的设备管理 Agent，具备对话、审计、建议及自我学习能力
- **部署**: Docker + Docker Compose

## 🤖 智能运维助手 (Agent Phase 2)

本项目集成了一个具有“专家心智”的智能助手，不仅仅是简单的对话，它深度集成了 EMS 系统的各项数据：

### 核心能力：
1. **多轮专家对话**: 支持上下文关联的设备状态深挖，自动匹配工业排查逻辑。
2. **三层记忆体系**:
   - **技能库 (Skill Store)**: 沉淀专家级的排查步骤（SOP），AI 按步骤执行工具调用。
   - **知识库 (Knowledge Wiki)**: 对话后自动提炼故障根因与预防措施，实现“越用越聪明”。
   - **经验库 (Experience Store)**: 记录并学习用户的管理偏好，实现个性化建议。
3. **主动感知审计**: 自动识别“级联失效”、“保养不力”等管理层容易忽视的隐蔽风险。
4. **工业数据模拟**: 内置 180 天具有逻辑因果关系的模拟数据，完美还原“救火式维修”与“预防性维护”的 ROI 对比。

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

```bash
./start-dev.sh
```

- 前端开发服务: http://localhost:5173
- 后端 API: http://localhost:8080
- API 文档: http://localhost:8080/swagger/index.html

开发环境会自动启动 PostgreSQL 和 Redis，并与 Docker 部署共用同一套数据卷。

### Docker 部署

```bash
# 启动所有服务（适合 VPS 上由宿主机 Nginx 反代）
docker compose up -d --build

# 查看日志
docker compose logs -f

# 停止服务
docker compose down
```

- 前端容器默认监听: `127.0.0.1:3000`
- 后端容器默认监听: `127.0.0.1:9000`
- 如需调整端口，可覆盖 `EMS_FRONTEND_PORT`、`EMS_BACKEND_PORT`

## 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | password123 | 管理员 |

## 开发文档

详细开发文档请参阅 [CLAUDE.md](./CLAUDE.md)

## License

MIT
