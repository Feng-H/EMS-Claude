# EMS 设备管理系统

> 企业级设备管理平台，支持点检、维修、保养、备件管理、统计分析等功能，集成飞书智能对话。
>
> 🌐 **在线演示**: [https://ems.317316.xyz](https://ems.317316.xyz)

## 📄 许可证与授权

**Copyright © 2025 Feng-H**

本项目采用 **GNU General Public License v3.0 (GPL-3.0)** 许可证开源。详见 [LICENSE](./LICENSE) 文件。

---

## 🤖 智能运维助手

本项目集成了 AI 智能助手，支持多轮对话、设备分析、预测性维护等功能，可通过 Web 界面或飞书机器人直接交互。

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
3. 记录下 **App ID** 和 **App Secret**

#### 第二步：配置事件订阅

1. 进入应用 -> 「事件与回调」->「事件配置」
2. **请求地址**填写：`https://你的域名/api/v1/lark/webhook/你的UID`
   > **如何获取 UID？**：登录 EMS 系统，进入「个人设置」->「飞书集成」，页面会显示你专属的 Webhook URL。
3. 记录页面显示的 **Verification Token**
4. 如果启用了 **Encrypt Key**，也一并记录

#### 第三步：配置 EMS 系统

1. 登录 EMS 系统，进入「个人中心」或「系统设置」->「飞书集成」
2. 填入飞书应用的 **App ID**、**App Secret**、**Verification Token** 和 **Encrypt Key**
3. 点击保存。系统会自动激活该机器人的 Webhook 处理逻辑。

#### 第四步：订阅消息事件

1. 在「事件与回调」->「事件配置」中点击「添加事件」
2. 搜索并添加 **`im.message.receive_v1`**（接收消息）事件
3. 如有版本管理，创建新版本并发布

#### 第五步：绑定 EMS 账号与配置机器人

1. **配置个人机器人**：登录 EMS 系统，在「个人中心」或「系统设置」中填写你的飞书应用凭证（App ID, App Secret 等）。系统会为你生成唯一的 Webhook 回调地址。
2. **在飞书开放平台配置 Webhook**：将生成的 Webhook 地址填入飞书应用的「事件与回调」配置中。
3. **完成绑定**：在手机飞书中打开机器人对话，发送任意消息，根据回复的链接完成账号绑定。

> **新功能**：系统现在支持**多机器人模式**。每个用户都可以配置属于自己的飞书机器人，实现私密的 AI 运维助手交互，互不干扰。

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
- **AI**: OpenAI 兼容接口（支持 DeepSeek、SiliconFlow 等）
- **部署**: Docker + Nginx

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
