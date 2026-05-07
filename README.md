# EMS 设备管理系统 v1.0

> **NEW: Agent 平台能力升级!** 企业级设备管理平台，现已支持标准化 Agent 接口 (MCP-Ready)、多维度管理驾驶舱、全生命周期资产管理。
>
> 🌐 **在线演示**: [https://ems.317316.xyz](https://ems.317316.xyz)

## 📄 许可证与授权

**Copyright © 2025 Feng-H**

本项目采用 **GNU General Public License v3.0 (GPL-3.0)** 许可证开源。详见 [LICENSE](./LICENSE) 文件。

---

## 🤖 核心能力：Agent 运维平台

本项目已从传统的管理系统演进为 **Agent 能力底座**，为外部智能体提供结构化的工业领域能力。

- **标准化工具集 (Tool Discovery)**: 提供 MCP 兼容的工具 Schema，外部 Agent 可通过 API 自动发现并调用 `查询设备`、`报修申请`、`预测健康度` 等功能。
- **工业级 API Key 管理**: 支持第三方系统/Agent 的安全接入，提供一世一签的密钥分发与生命周期管理。
- **主动事件推送 (Proactive Push)**: 配置化推送规则，当发生点检 NG、备件低库存或关键故障时，系统主动向 Agent 触发 Webhook 或飞书通知。
- **混合 RAG 智能诊断**: 优化的混合检索算法（专家知识库 + 技术手册分片），为 LLM 提供高置信度的维修建议。

---

## 📊 决策支持：管理驾驶舱

全新升级的统计分析系统，为管理层提供实时、精准的数据洞察：

- **核心指标分析**: 实时监控 MTTR (平均维修时间)、MTBF (平均无故障时间) 及资产可用性评分。
- **全工厂维度过滤**: 支持按集团、工厂、车间多级下钻分析，实时刷新关键业务指标。
- **设备可靠性排名**: 自动生成故障 TOP10、停机损失排名及可用性黑榜，辅助决策资产更新与备件储备。
- **任务趋势看板**: 近 30 天点检、维修、保养任务的动态趋势对比。

---

## 🚀 快速部署

### 1. 克隆项目

```bash
git clone https://github.com/Feng-H/EMS-Claude.git
cd EMS-Claude
```

### 2. 配置环境变量

复制模板并生成密钥：

```bash
cp .env.example .env

# 生成 JWT 密钥（必填）
echo "EMS_JWT_SECRET=$(openssl rand -hex 32)" >> .env

# 生成数据库密码
echo "EMS_DATABASE_PASSWORD=$(openssl rand -hex 16)" >> .env
```

编辑 `.env`，按需配置 LLM：

```bash
# LLM 智能助手配置 (默认使用 SiliconFlow/DeepSeek)
EMS_LLM_PROVIDER=openai
EMS_LLM_BASE_URL=https://api.siliconflow.cn/v1
EMS_LLM_API_KEY=sk-xxxx...
EMS_LLM_MODEL=deepseek-ai/DeepSeek-V3

# 域名配置 (可选，用于生成飞书回调基地址)
EMS_DOMAIN=ems.yourdomain.com
```

### 3. 启动服务

```bash
docker compose up -d --build
```

首次启动会自动：
- 创建数据库表结构
- 填充演示数据（设备、维修工单、保养计划、知识库等）

启动后访问：`http://你的IP:3000`

默认账号：`admin` / `admin123`（首次登录需修改密码）

### 4. 常用命令

```bash
# 查看日志
docker compose logs -f backend

# 重启服务
docker compose restart

# 清空数据重新部署（会重新填充演示数据）
docker compose down -v
docker compose up -d --build

# 仅更新代码重新构建
git pull origin main
docker compose up -d --build
```

---

## 📱 飞书机器人集成

通过飞书机器人可以直接在飞书中与 AI 智能助手对话，查询设备状态、维修记录等。

### 配置步骤

#### 第一步：创建飞书应用
1. 登录 [飞书开放平台](https://open.feishu.cn/)
2. 点击「创建企业自建应用」，选择「机器人」
3. 进入「凭证与基础信息」，记录下 **App ID** 和 **App Secret**
4. 进入「事件与回调」，在页面中找到并记录 **Verification Token** (校验令牌)

#### 第二步：在 EMS 系统中预配置
> **关键步骤**：必须先在 EMS 保存配置，飞书平台的回调地址才能验证通过。
1. 登录 EMS 系统，进入「个人中心」或「系统设置」->「飞书集成」
2. 填入第一步获取的 **App ID**、**App Secret** 和 **Verification Token**
3. 如果计划使用加密，也请填入 **Encrypt Key**
4. 点击 **「保存」**。此时系统已准备好处理飞书的验证请求。
5. 复制页面显示的 **Webhook URL** (例如：`https://你的域名/api/v1/lark/webhook/1`)

#### 第三步：配置飞书事件订阅
1. 返回飞书开放平台 ->「事件与回调」->「事件配置」
2. 点击「编辑」请求地址，填入上一步复制的 **Webhook URL**
3. 点击 **「保存」**。飞书会立即发送验证挑战，由于 EMS 已提前保存配置，验证将顺利通过。
   > **注意**：如果先在飞书保存而未在 EMS 配置，飞书会提示「挑战码验证失败」。

#### 第四步：添加订阅事件
1. 在同一页面点击「添加事件」
2. 搜索并添加 **`im.message.receive_v1`**（接收消息）事件
3. 如有版本管理，创建新版本并发布应用。

#### 第五步：完成绑定
1. 在手机飞书中搜索并打开你的机器人对话框。
2. 发送任意消息，机器人会回复一个绑定链接。
3. 点击链接并登录 EMS 账号，即可完成飞书账号与 EMS 账号的互关联。

> **新功能**：
> - **多机器人模式**：每个用户都可以配置属于自己的飞书机器人，实现私密的 AI 运维助手交互，互不干扰。
> - **海量演示数据**：系统默认填充 100+ 台模拟设备及 15 天以上的历史点检、维修、保养记录，支持深度 AI 分析与故障预测演示。

#### 第六步：开始使用

绑定完成后，直接在飞书中给机器人发消息即可获得 AI 回复。机器人会先回复 "👍 收到，正在分析中..."，处理完成后发送完整分析结果。

![飞书机器人对话效果](飞书机器人agent.jpeg)

### 配置架构

```
用户手机飞书 -> 飞书开放平台 -> Webhook -> Nginx(前端) -> 后端 API
                                                  -> Agent AI -> LLM -> 飞书API -> 回复用户
```

---

## 技术栈

- **后端**: Go 1.23 + Gin + GORM + PostgreSQL + Redis
- **前端**: Vue 3 + TypeScript + Vite + Element Plus (PC) + Vant 4 (H5)
- **Agent 协议**: 标准 JSON Schema 工具描述，MCP (Model Context Protocol) 架构兼容。
- **AI 能力**: 深度集成 DeepSeek-V3/R1，支持预测性维护 (RUL 预测) 与故障链分析。
- **H5 移动端**: 
  - **全屏体验**: 独立于 PC 端的全屏移动端布局，专为现场操作优化。
  - **闭环作业**: 扫码点检、现场报修、移动执行、实时通知。
- **部署**: Docker + Nginx (生产级配置)

## 项目结构

```
EMS-Claude/
├── backend/          # Go 后端服务
├── frontend/         # Vue3 前端应用 (PC + H5 移动端)
├── db/               # 数据库结构参考与种子数据 (文档用)
├── docs/             # 开发文档
└── .env.example      # 环境变量模板
```

## 开发指引

详细开发规范请参阅 [CLAUDE.md](./CLAUDE.md)
