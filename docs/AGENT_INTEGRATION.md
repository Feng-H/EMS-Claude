# EMS Agent 平台接入与配置指南 (v1.0)

EMS (Equipment Management System) 不仅仅是一个管理后台，更是一个标准的 **Agent 能力平台**。本指南详细说明如何通过 API 和标准化协议将外部智能体（如 OpenClaw, Hermes, AutoGPT 等）接入 EMS 系统。

---

## 1. 身份认证 (API Key)

外部 Agent 访问 EMS 必须通过 API Key 进行身份验证。

### 配置步骤：
1. **生成密钥**：登录 EMS Web 端，进入 `智能助手` -> `Agent 集成` -> `API 密钥`。
2. **创建密钥**：点击 "创建新密钥"，输入应用名称（如 "My-External-Agent"）并选择有效期。
3. **保存密钥**：系统会生成一个以 `ems_` 开头的密钥（例如 `ems_7f8a...`）。**该密钥仅显示一次**，请务必妥善保存。

### 使用方法：
在所有请求的 HTTP Header 中添加 `X-API-KEY` 字段：

```http
GET /api/v1/agent/tools HTTP/1.1
Host: ems.yourdomain.com
X-API-KEY: ems_您的密钥内容
```

---

## 2. 工具集发现 (Tool Discovery)

EMS 遵循类 **MCP (Model Context Protocol)** 架构，通过标准化 Schema 暴露其领域能力。

### 获取工具列表：
**Endpoint**: `GET /api/v1/agent/tools`

该接口返回系统当前支持的所有工具及其参数定义（JSON Schema）。

**响应示例**：
```json
{
  "tools": [
    {
      "name": "search_equipment",
      "description": "通过名称或编码搜索设备",
      "input_schema": {
        "type": "object",
        "properties": {
          "keyword": { "type": "string", "description": "搜索关键词" }
        },
        "required": ["keyword"]
      }
    },
    {
      "name": "get_equipment_health",
      "description": "获取设备实时健康度预测 (RUL)",
      "input_schema": {
        "type": "object",
        "properties": {
          "equipment_id": { "type": "integer" }
        },
        "required": ["equipment_id"]
      }
    }
  ]
}
```

---

## 3. 工具调用 (Tool Call)

Agent 根据发现的 Schema 构造参数，调用具体功能。

### 调用接口：
**Endpoint**: `POST /api/v1/agent/tools/call`

**请求示例 (报修申请)**：
```json
{
  "name": "report_repair",
  "arguments": {
    "equipment_id": 105,
    "fault_description": "2号轴承运行噪音异常，温度升高至 85℃",
    "priority": 1
  }
}
```

**成功响应**：
```json
{
  "content": "Repair order #20260507001 created successfully",
  "is_error": false
}
```

---

## 4. 主动事件推送 (Proactive Push)

您可以配置 EMS 在特定事件发生时主动通知 Agent。

### 支持的推送类型：
- `ng_inspection`: **点检异常**。当现场人员提交 NG 结果时触发。
- `repair_request`: **新报修工单**。当有新的报修申请需要分析时触发。
- `low_stock`: **备件预警**。当备件库存低于安全水位时触发。

### 配置方法：
在 `Agent 集成` -> `主动推送配置` 中：
1. 开启对应开关。
2. 配置 **推送范围 (Scope)**：使用 JSON 格式过滤。例如 `{"factory_id": 1}` 表示仅接收 1 号工厂的推送。
3. **Webhook 接收**：系统会将事件推送到您配置的 Webhook 地址或关联的飞书机器人。

---

## 5. 进阶：RAG 知识检索

Agent 可以利用 EMS 优化的混合检索能力来获取专业维修建议。

**Endpoint**: `GET /api/v1/agent/knowledges`
**参数**: `query` (问题关键词), `equipment_type_id` (可选，过滤设备类型)

**检索逻辑**：
- **权重 1.0**: 专家审核过的 `知识库文章`。
- **权重 0.8**: 原始 `设备技术手册` 切片。
- **关联分析**: 系统会自动关联该设备的历史维修记录，生成全证据链。

---

## 6. 安全建议

- **权限隔离**：API Key 继承其所属用户的角色权限。建议为外部 Agent 创建专门的 "智能体账号" 并赋予 `engineer` 或 `viewer` 角色。
- **频率限制**：API 默认设有速率限制，请避免高频轮询，优先使用推送机制。
- **HTTPS**：生产环境务必全程强制 HTTPS，防止 API Key 被截获。

---

## 开发者支持

如有集成问题，请查阅 [API 交互时序图](./docs/archive/ems_agent_design.md) 或联系系统管理员。
