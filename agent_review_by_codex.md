# EMS Agent 实现探索与改进建议

日期：2026-05-08  
作者：Codex

## 1. 结论摘要

项目已经实现了一个相对完整的 Agent 雏形：后端有独立 `backend/internal/agent` 模块，包含对话、技能、知识草稿、专项分析、预测性维护、工具发现与工具调用；前端提供了智能助手工作台和 Agent 集成页；认证层也已支持 `X-API-KEY` 形式的外部 Agent 接入。

但当前实现更接近“演示级 Agent 能力底座”，还没有达到可稳定开放给外部 Agent 的生产级水平。主要差距集中在四类：权限隔离不一致、多个核心分析仍是 mock/占位、外部工具协议较薄、主动推送只有配置没有完整投递链路。尤其是部分详情接口和工具调用绕过了用户工厂权限，需要优先修复。

## 2. 探索范围

重点阅读了以下文件：

- `backend/internal/agent/service/agent.go`
- `backend/internal/agent/service/tool_service.go`
- `backend/internal/agent/controller/agent.go`
- `backend/internal/agent/analyzer/*.go`
- `backend/internal/agent/tool/*.go`
- `backend/internal/agent/policy/policy.go`
- `backend/internal/agent/repository/*.go`
- `backend/internal/model/model.go`
- `backend/internal/middleware/auth.go`
- `backend/api/v1/auth.go`
- `backend/internal/service/lark.go`
- `frontend/src/views/agent/ManagementAssistantView.vue`
- `frontend/src/views/agent/AgentIntegrationView.vue`
- `frontend/src/api/agent.ts`
- `docs/AGENT_INTEGRATION.md`
- `db/schema.sql`

## 3. 当前实现了什么

### 3.1 Agent 服务骨架

`AgentService` 会根据存储模式选择 DB 或 memory repository，并组装 retrieval、maintenance、repair、prompt、LLM client、各类 analyzer。LLM 使用 OpenAI-compatible chat completions，未配置 API key 时进入预览/规则回退模式。

核心入口包括：

- 保养推荐：`RecommendMaintenance`
- 维修审计：`AuditRepair`
- 保养审计：`AuditMaintenance`
- 通用分析：`Analyze`
- 多轮对话：`Chat`
- 技能管理与执行：`CreateSkill` / `UpdateSkill` / `ExecuteSkill`
- 预测性维护：`GetEquipmentPrediction`
- 外部工具发现/调用：`ListTools` / `CallTool`
- 主动推送订阅：`Subscribe` / `ListSubscriptions`

### 3.2 对话与自主学习

`Chat` 支持创建或复用会话、保存用户与助手消息、读取最近 10 条历史、调用 LLM 生成回复，并异步触发 `ReflectAndLearn`。

学习链路目前包括：

- 从对话中提取知识草稿，写入 `agent_knowledges`
- 从对话中提取技能草稿，写入 `agent_skills`
- 预留用户经验提取接口，但 `asyncCollectExperience` 为空实现

前端 `ManagementAssistantView.vue` 已提供专家对话、专项审计、知识审核、飞书绑定、右侧预测洞察面板。

### 3.3 技能系统

技能以 `AgentSkill` 存储，`steps` 是 JSON 字符串。执行时按 step tool 名称分发到固定工具：

- `get_equipment_profile`
- `get_failure_stats`
- `get_cost_analysis`
- `get_maintenance_compliance`
- `get_failure_distribution`
- `search_manual_knowledge`
- `predict_remaining_life`
- `detect_symptoms`
- `get_tco_analysis`
- `get_retirement_recommendation`

如果配置了 LLM，会把收集到的证据链再交给 LLM 总结。

### 3.4 专项分析与预测

已经有三个 analyzer：

- `MaintenanceAnalyzer`：读取保养计划、搜索手册/知识库，按周期大于 30 天给出缩短周期建议。
- `RepairAuditAnalyzer`：返回重复故障审计结果，但目前是固定 mock 案例。
- `PredictiveAnalyzer`：提供 RUL、亚健康征兆、TCO、退役建议。

预测性维护目前混合了真实查询和演示常量，例如固定超负荷系数、固定标准 MTBF、固定已运行小时数、固定保养无效性征兆。

### 3.5 外部 Agent 接入

外部接入层已经具备基本能力：

- API Key 创建、列表、删除。
- `AuthMiddleware` 优先识别 `X-API-KEY`，再回退 JWT。
- `GET /api/v1/agent/tools` 返回工具 JSON Schema。
- `POST /api/v1/agent/tools/call` 执行工具。

当前暴露工具包括：

- `search_equipment`
- `get_equipment_health`
- `get_spare_part_inventory`
- `report_repair`

前端 `AgentIntegrationView.vue` 已提供 API Key 管理、工具发现展示、主动推送配置 UI。

### 3.6 飞书集成

飞书 webhook 收到消息后会定位 bot owner 和发送者。如果发送者已绑定 EMS 用户，则转成 `ChatRequest` 调用 Agent；如果未绑定，则发送绑定链接。前端也集成了飞书机器人配置与绑定引导。

## 4. 能力成熟度判断

| 能力 | 当前成熟度 | 判断 |
| --- | --- | --- |
| 多轮对话 | MVP | 能持久化历史并调用 LLM，但缺少真实业务上下文检索和权限校验 |
| 保养推荐 | MVP | 有规则、证据与 LLM 总结，但规则较粗 |
| 维修审计 | 演示级 | 返回固定重复故障案例，未真正扫描维修单 |
| 保养审计 | 占位 | 返回“开发中” |
| 通用分析 | 占位 | 返回“开发中” |
| 技能管理 | MVP | CRUD 与执行链路存在，但匹配和工具编排较弱 |
| 自主学习 | 雏形 | 可提取知识/技能草稿，缺少去重、审核流约束和经验学习 |
| RAG 检索 | 初级 | LIKE/包含匹配 + 简单权重，未做向量/全文索引 |
| 预测性维护 | 演示级 | 有 RUL/TCO 接口，但大量硬编码假设 |
| 外部工具发现 | MVP | 有 JSON Schema，但不是完整 MCP/任务协议 |
| API Key | MVP | 可用，但明文存储、无 scope、无限流 |
| 主动推送 | 配置雏形 | 有订阅表和 UI，没有完整事件触发与投递 |
| 飞书 Agent | MVP | 可以桥接消息到 Chat，但安全与上下文能力仍需增强 |

## 5. 优先级最高的问题

### P0：权限隔离存在绕过风险

1. `GetSession`、`GetArtifact`、`GetConversation` 只按 ID 查询，没有校验资源是否属于当前用户。Controller 层也没有传入 userID 做归属检查。普通用户如果猜到 ID，可能读取他人的会话、产物、证据。

2. `GetEquipmentPrediction` 在 service 中直接构造 `model.User{Role: "admin"}`，注释还写着“默认按管理员权限”。这会绕过工厂隔离。更严重的是外部工具 `get_equipment_health` 也调用这个方法，所以 API Key 用户也可能查询跨工厂设备健康数据。

3. `report_repair` 创建维修单前没有校验 equipment 是否在当前用户授权工厂内。`get_spare_part_inventory` 如果传入 `factory_id`，也没有校验该 factory 是否属于当前用户可访问范围。

4. `ListKnowledges`、`ListSkills`、`CreateSkill`、`UpdateSkill` 等管理能力只要求登录，没有明确角色授权。知识审核和技能发布属于高影响操作，至少应限制 admin/manager 或专门权限。

建议：

- 所有 by-id 查询接口改为 `GetX(user, id)`，repository 查询中带 `user_id` 或 factory join。
- `GetEquipmentPrediction` 接收真实 user，并在 analyzer/tool 层延续权限检查。
- 外部工具统一走 Tool Registry，每个工具声明 `scope`、`read/write`、`required_role`、`factory_bound`。
- 对写操作工具如 `report_repair` 增加 equipment ownership 校验。

### P0：API Key 安全模型不足

API Key 当前以明文存储在 `user_api_keys.key`，认证时直接按 key 查询。它可继承用户全部角色权限，没有细粒度 scope，没有 rate limit，没有 IP / origin 限制，也没有独立审计日志。

建议：

- 只存 key hash，明文仅创建时显示一次。
- 增加 `scopes`、`allowed_tools`、`allowed_factory_ids`、`rate_limit_per_minute`、`last_used_ip`。
- 外部 Agent 默认只读，写工具需要显式授权。
- 为 `CallTool` 记录结构化审计日志，包括 tool name、arguments 摘要、user_id、api_key_id、trace_id、结果状态。

### P0：用户可覆盖 system prompt

`Chat`、`RecommendMaintenance`、`AuditRepair` 等请求 DTO 都允许传入 `system_prompt`，service 中会直接替换或拼接系统提示词。这对普通用户和外部 Agent 都是高风险入口，可能绕过安全策略、改变输出边界、诱导泄露数据。

建议：

- 生产接口不接受任意 `system_prompt`。
- 如需调试，限制为 admin-only debug 参数，并在非生产环境开启。
- 用服务端固定 policy prompt + 用户问题 + 检索上下文组合，不允许客户端覆盖系统角色。

## 6. 功能正确性与产品落差

### P1：多个接口仍是占位或 mock

`AuditMaintenance` 和 `Analyze` 明确返回“开发中”占位响应。`RepairAuditAnalyzer` 写明是 MVP 假设案例，并固定返回 `repair_order:101`、`EQ-JJ-001` 等内容。这会让前端看起来完成了审计，实际不具备可信业务价值。

建议：

- 将占位接口在前端标记为“实验/开发中”，或后端返回 `501 Not Implemented`。
- 维修审计先落地三类真实规则：短期重复故障、异常维修成本、维修后 72 小时再次报修。
- 每个审计结论必须引用真实 `repair_orders`、`knowledge_articles`、`maintenance_tasks` 记录。

### P1：预测性维护算法演示味较重

RUL 预测使用固定 `loadFactor=1.15`、`avgMTBFHours=300`、`currentUsedHours=240`。症状识别固定插入“二级保养后 72 小时液压系统报警”。TCO 的使用年限对 `PRESS-05` 做特殊判断。

建议：

- 引入 `equipment_runtime_snapshots`、设备采购日期、维修工单时间序列、保养任务时间序列作为真实特征。
- 将预测结果标注为规则模型，并输出置信度来源。
- 把硬编码 demo 逻辑迁移到 seed/demo 数据或 feature flag，避免生产误报。

### P1：Chat 说“引用证据”，但常规对话没有注入证据

标准对话 fallback 只把历史消息交给 LLM，system prompt 要求“结论必须引用系统中的财务与技术证据”，但没有检索设备、维修、备件、知识库上下文。只有匹配到技能时才会收集证据。

建议：

- 在 Chat 前增加 intent + entity extraction：设备、工厂、时间范围、指标类型。
- 对常规问题也执行 retrieval plan，把业务数据作为 `context` 注入。
- 回复结构中返回 evidence ids，让前端可展开溯源。

### P1：工具发现还不是稳定 Agent 协议

当前工具 schema 只有 `name`、`description`、`input_schema`。缺少 output schema、错误码、幂等性、权限、版本、读写属性、分页、trace、dry-run、schema 版本。`CallToolResponse.Content` 是任意类型，外部 Agent 很难稳定解析。

建议：

- 引入 tool manifest：`name/version/input_schema/output_schema/scopes/side_effect/idempotency/errors/examples`。
- 写工具支持 `dry_run`，并返回可确认的 preview。
- 统一错误 envelope：`success=false, error.code, error.message, trace_id`。
- 对外部 Agent 提供 OpenAPI 或 MCP adapter 文档。

### P1：主动推送只有订阅配置，没有完整闭环

代码中 `Subscribe` 和 `ListSubscriptions` 可保存配置，`NotifyEvent` 只会在 RUL 小于 7 天时创建 artifact，没有查询订阅、没有 webhook 地址、没有飞书/邮件投递、没有失败重试。代码搜索也没有发现业务事件调用 `NotifyEvent`。

建议：

- 建立事件总线或 domain hooks：点检 NG、维修创建、低库存、预测高风险。
- `agent_push_subscriptions` 增加 destination、webhook_url、secret、delivery_status 等字段。
- 增加 delivery worker、签名、重试、死信队列和投递日志。

### P1：文档与接口行为不一致

`docs/AGENT_INTEGRATION.md` 写 `GET /api/v1/agent/knowledges` 支持 `query` 和 `equipment_type_id`，但 controller 只读取 `status`。文档还写默认有频率限制、主动推送到 webhook/飞书，但代码未见完整实现。

建议：

- 将文档标记为“当前已实现”和“规划中”两部分。
- 对已发布接口加契约测试，保证文档示例和响应结构一致。

## 7. 工程质量与可维护性问题

### P1：Service 内部直接 new 依赖，不利于测试

`NewAgentService` 直接读取全局 config、database、repository、LLM client。`tool_service.go` 和 subscription 逻辑也直接调用 `database.GetDB()`。这会让单元测试和权限回归测试很难写。

建议：

- 引入 `AgentServiceOptions` 或 constructor 注入：repo、policy、tool registry、llm client、clock、logger。
- 外部工具调用和主动推送也通过 repository/service 接口，不直接拿全局 DB。

### P1：错误处理偏松

多处 `_ = repo.CreateMessage`、`_ = json.Unmarshal`、`_ = repo.CreateKnowledge` 会吞掉关键错误。异步学习失败仅日志，前端无法知道提取失败原因。

建议：

- 写入会话/消息失败应返回错误或至少影响 trace 状态。
- JSON 解析失败记录原始模型响应片段、错误原因和 convID。
- 产物、证据链、usage 写入建议使用事务。

### P2：数据模型与 schema 文件不一致

GORM model 中有 `AgentPushSubscription` 和 `UserAPIKey`，`main.go` 会 AutoMigrate，但 `db/schema.sql` 没有完整创建这些表；同时 schema 中 `agent_skills` 是 JSONB，而 GORM model 使用 text 字符串。这会增加本地初始化、迁移和生产部署的不确定性。

建议：

- 统一使用显式 migration，减少 AutoMigrate 与 SQL schema 分叉。
- JSON 字段统一建模为 JSONB 或 datatypes.JSON，并补充索引。

### P2：测试覆盖不足

仓库里没有 agent 相关 `_test.go`。目前只有 middleware/config/inspection 少量测试。考虑到 Agent 涉及权限、写操作和外部 API，测试缺口较大。

建议优先补：

- API Key 认证、过期、禁用、last_used_at。
- `GetConversation/GetSession/GetArtifact` 越权访问。
- `CallTool` 四个工具的权限边界。
- Chat 无 LLM、LLM 失败、技能匹配、知识提取失败。
- Repair audit / maintenance recommendation 的规则输出。

## 8. 建议改进路线

### 第一阶段：先把安全边界补牢

1. 修复所有 by-id agent 资源的用户归属校验。
2. `GetEquipmentPrediction` 和 `get_equipment_health` 改为真实用户上下文。
3. `report_repair`、`get_spare_part_inventory` 增加工厂权限校验。
4. 禁止普通请求传入 `system_prompt`。
5. API Key 改 hash 存储，并增加最小 scope 字段。

### 第二阶段：把演示能力改成可信业务能力

1. `AuditRepair` 接入真实维修单扫描。
2. `AuditMaintenance` 和 `Analyze` 要么实现，要么显式返回未实现。
3. Chat fallback 接入业务数据检索和证据链。
4. RUL/TCO 去掉硬编码，明确规则模型输入与置信度。

### 第三阶段：外部 Agent 平台化

1. 建立统一 Tool Registry 和 tool manifest。
2. 工具响应结构固定化，增加 output schema 和错误码。
3. 写工具支持 dry-run / confirmation。
4. 主动推送补齐事件触发、投递、签名、重试、日志。

### 第四阶段：工程化与可观测性

1. 依赖注入改造 AgentService。
2. 关键链路加 trace_id 贯穿日志、usage、artifact、tool call。
3. LLM token 用量、耗时、失败率、工具成功率入库。
4. 补齐 agent 单元测试与接口契约测试。

## 9. 可以优先拆出的任务清单

1. `P0-agent-authz`: 修复会话、产物、预测、工具调用权限绕过。
2. `P0-api-key-hardening`: API Key hash 存储、scope、审计日志。
3. `P1-tool-registry`: 把工具定义和执行从 switch 改成注册表。
4. `P1-repair-audit-real-data`: 维修审计从 mock 改为真实规则扫描。
5. `P1-chat-evidence-context`: Chat 增加实体识别、数据检索、证据返回。
6. `P1-proactive-push-worker`: 主动推送补齐触发和投递 worker。
7. `P2-agent-tests`: 增加权限、工具、对话、学习链路测试。

## 10. 总体评价

这个项目的 Agent 方向是清晰的：它已经把 EMS 从普通 CRUD 后台往“工业领域能力平台”推进了一步，尤其是会话、技能、知识、工具发现、飞书桥接这些模块都已成形。现在最需要的不是继续堆更多 Agent 概念，而是把现有能力做实：权限可信、证据可信、接口可信、文档可信。

如果先完成 P0 和 P1 的前半部分，这个 Agent 模块就能从“好看的演示”提升到“内部可用的智能运维助手”。之后再完善工具协议和主动推送，才适合对 OpenClaw、Hermes Agent 或其他外部智能体稳定开放。
