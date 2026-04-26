# EMS 设备管理系统

> 企业级设备管理平台，支持点检、维修、保养、备件管理、统计分析等功能。

## 📄 许可证与授权

**Copyright © 2025 Feng-H**

本项目采用 **GNU General Public License v3.0 (GPL-3.0)** 许可证开源。详见 [LICENSE](./LICENSE) 文件。

---

## 🤖 智能运维助手 (Agent Phase 3 完全体)

本项目集成了一个具有“专家心智”的智能助手，具备 **L4 级主动洞察** 能力：

### 核心能力：
1. **多轮资产战略对话**: 基于 TCO (总持有成本) 和 RUL (剩余健康寿命) 进行深度财务与技术对标。
2. **三层记忆体系**:
   - **技能库**: 沉淀级联失效审计、退役 ROI 评价等专家 SOP。
   - **知识库**: 自动从对话中提炼“劣质滤芯导致泵报废”等结构化经验。
   - **经验库**: 自动学习用户偏好，实现个性化管理建模。
3. **180 天逻辑仿真数据**: 
   - 完美模拟“李四 (预防专家)”与“张三 (救火队长)”的行为差异。
   - 提供 PRESS-05 等超期服役设备的财务退役预警案例。
4. **双模一致性部署**: 无论是内存模式还是 Docker PostgreSQL 模式，均能获得完全一致的演示体验。

---

## 🚀 部署与配置 (Deployment)

### 1. 配置环境变量 (.env)
在项目根目录创建或编辑 `.env` 文件，填入您的 LLM 密钥。本项目已通过 Docker 的 `env_file` 机制实现密钥的安全穿透。

```bash
## LLM config (以 SiliconFlow 为例)
EMS_LLM_PROVIDER=openai
EMS_LLM_BASE_URL=https://api.siliconflow.cn/v1
EMS_LLM_API_KEY=sk-xxxx... # 在此填入您的密钥
EMS_LLM_MODEL=deepseek-ai/DeepSeek-V3
```

### 2. 生产环境部署 (Docker)
```bash
# 启动所有服务
docker compose up -d --build

# 重启后端以应用新的 .env 配置
docker compose restart backend
```

### 🛡️ 安全性保障 (Security)
- **密钥不泄露**: API Key 仅在**后端容器内部**内存中加载。前端浏览器无法访问，且 Key 不会被打包进 Docker 镜像。
- **环境隔离**: 敏感配置通过 `.env` 管理并由 Docker 注入，符合 12-Factor App 安全规范。

---

## 技术栈

- **后端**: Go 1.23 + Gin + GORM + PostgreSQL + Redis
- **前端**: Vue 3 + TypeScript + Vite + Element Plus
- **移动端**: Vue 3 + Vant 4 (H5)
- **UI 风格**: 现代化工业风，针对智能运维优化的交互设计

## 默认账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | admin123 | 管理员 |

## 项目结构
```
EMS-Claude/
├── backend/          # Go后端服务
├── frontend/         # Vue3前端应用
├── db/              # 数据库结构与 180天 仿真种子数据
├── docker/          # Docker部署配置
└── docs/            # Phase 1-3 深度开发文档
```

## 开发指引
详细开发规范请参阅 [CLAUDE.md](./CLAUDE.md)
