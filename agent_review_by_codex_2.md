# EMS Agent 第二轮复审报告

日期：2026-05-08  
作者：Codex  
复审基准：`fe61781 docs: add agent implementation review` -> 当前 `HEAD` (`3d57e18`)

## 1. 结论摘要

Gemini 已经按上一版建议做了一轮明显增强：新增 Tool Registry、API Key 哈希存储、API Key scopes/rate_limit 字段、会话/产物/对话详情归属校验、知识/技能管理角色限制、维修审计真实工单扫描雏形、保养审计和通用分析接口替代了原来的占位响应，并补了少量单元测试。

但当前代码有一个硬阻断：后端无法编译，`go test ./...` 失败。因此这轮更新尚不能合入或发布。除此之外，若只看设计意图，P0 安全问题已经部分收敛；但 scopes/rate_limit 没有真正生效、维修审计仍可能跨工厂读数据、主动推送仍未闭环、部分“实现”仍是演示逻辑或固定数据。

## 2. 验证结果

执行命令：

```bash
cd backend
go test ./...
```

结果：失败，构建未通过。

主要错误：

```text
internal/agent/prompt/prompt.go:75:1: syntax error: unexpected ..., expected }
internal/agent/prompt/prompt.go:76:2: syntax error: unexpected keyword return after top level declaration
internal/agent/tool/registry.go:4:2: "context" imported and not used
```

这说明新增的 agent 测试目前没有真正形成有效保护，因为代码在进入测试执行前已经构建失败。

## 3. 已改善的地方

### 3.1 Agent 资源详情加了归属校验

上一版指出 `GetSession`、`GetArtifact`、`GetConversation` 可被猜 ID 越权读取。当前 service 已新增 userID/role 参数，并对非 admin 用户做归属检查：

- `GetSession(id, userID, role)`
- `GetArtifact(id, userID, role)`
- `GetConversation(id, userID, role)`

这是有效改进。下一步应补对应接口测试，尤其是“用户 A 不能读取用户 B 会话/产物”。

### 3.2 设备预测接口改为真实用户上下文

`GetEquipmentPrediction` 已从原来的 admin 绕过改成接收 `user model.User`，controller 也会加载当前用户后调用。这修复了 Web 端和 `get_equipment_health` 工具最明显的跨工厂绕过风险。

仍需注意：方法内部吞掉了 `PredictRUL`、`CalculateTCO`、`DetectSymptoms` 的错误，权限失败时可能返回 `{rul:nil, tco:nil, symptoms:nil}` 而不是明确 403/错误。

### 3.3 report_repair 增加了设备工厂校验

`handleReportRepair` 创建维修单前会加载设备，并对非 admin 用户检查设备所属 workshop.factory_id。这比上一版安全很多。

仍需注意：如果 `equipment_id` 参数缺失或类型不对，当前默认 `equipID=0`，最后返回 “equipment not found”，错误不够结构化；priority 也缺少范围校验。

### 3.4 API Key 改为哈希存储

创建 API Key 时，现在只把 SHA-256 hash 写入 `user_api_keys.key`，认证时 hash 请求头后查询。明文 key 只在创建时返回一次。这比上一版明文存储更好。

新增字段：

- `scopes`
- `rate_limit`

前端也新增了 scopes 多选和 rate limit 输入。

### 3.5 知识与技能管理加了角色限制

知识审核、知识列表、技能列表、技能创建、技能详情、技能更新都限制为 `admin` 或 `manager`。这修复了上一版“普通登录用户可以管理 Agent 知识/技能”的问题。

### 3.6 Chat 开始注入业务上下文

标准对话 fallback 会尝试：

- 从消息中提取设备 ID
- 注入设备 profile 和 health
- 搜索知识库/手册并追加前两条参考

这比上一版“要求引用证据但不注入证据”更接近真实 Agent 行为。

### 3.7 维修审计从固定 mock 前进到真实工单扫描

`RepairAuditAnalyzer` 现在会读取维修单，按设备分组识别 72 小时内相似故障，并检测高额维修成本。这比原来的固定 `repair_order:101 / EQ-JJ-001` 好很多。

### 3.8 主动推送补了基础投递结构

新增了 `AgentPushLog`，`AgentPushSubscription` 增加 `WebhookURL` 和 `Secret` 字段，`deliverPush` 会向 webhook POST payload，并记录 success/failed。

这是从“只有订阅配置”到“有最小投递雏形”的进步。

## 4. 阻断问题

### P0：当前后端无法编译

文件：`backend/internal/agent/prompt/prompt.go`

`BuildKnowledgeExtractionPrompt` 函数体里多了一行裸 `...`，Go 语法不允许，导致整个后端构建失败。

文件：`backend/internal/agent/tool/registry.go`

导入了 `context` 但未使用，也会导致构建失败。

建议立即修复：

1. 删除 `prompt.go` 第 75 行的 `...`。
2. 删除 `registry.go` 未使用的 `context` import。
3. 重新跑 `go test ./...`。

## 5. 高优先级残留问题

### P0：API Key scopes 和 rate_limit 只是字段，未真正生效

虽然 `UserAPIKey` 增加了 `Scopes` 和 `RateLimit`，前端也能配置，但认证中间件只设置了 `user_id` 和 `role`，没有把 api_key_id/scopes/rate_limit 放进 context。

`ToolRegistry.Call` 接收 `userScopes []string`，但 `CallTool` 直接传 `nil`，registry 内部也明确写着“production 再检查”。因此一个只有 `read:equipment` 的 API Key 仍可能调用 `report_repair`。

影响：

- scope 是 UI/数据字段，不是安全边界。
- rate_limit 完全没有执行。
- 写工具依然只依赖用户角色/工厂，不依赖 API Key 授权范围。

建议：

- AuthMiddleware 设置 `api_key_id`、`api_key_scopes`、`api_key_rate_limit`。
- `CallTool` 从 gin context 或 service 上下文传入 scopes。
- Tool Registry 对 read/write 工具统一检查 required scopes。
- rate limit 至少按 api_key_id + minute 在 Redis 或内存中计数。

### P0：维修审计仍可能跨工厂读取工单

文件：`backend/internal/agent/tool/repair.go`

`GetOrdersByFilter` 注释说要按用户工厂过滤，但实际没有过滤，DB 模式直接 `return t.orderRepo.List(filter)`。`RepairAuditAnalyzer` 依赖这个方法读取工单，因此非 admin 用户可能审计到其他工厂维修单。

这是上一版“工厂级隔离”在新维修审计实现中的回归风险。

建议：

- `RepairOrderFilter` 增加 `FactoryID` 或在 repository 层 join equipment/workshops 过滤。
- memory 模式也要按 equipment -> workshop -> factory 过滤。
- 增加非 admin 用户跨工厂审计测试。

### P0：主动推送仍使用 admin 权限做预测

文件：`backend/internal/agent/service/agent.go`

`NotifyEvent` 中仍用 `model.User{Role: "admin"}` 调用 `PredictRUL`。虽然主动事件可能由系统触发，但后续订阅匹配和投递必须按订阅用户 scope 隔离；当前代码查出所有 `predictive_maintenance` 订阅后直接投递，没有验证 target 是否属于订阅用户/订阅 scope。

建议：

- 对每个 subscription 加载订阅用户，按该用户上下文重新校验 target。
- 实现 `scope` 匹配，如 factory_id/workshop_id/equipment_ids。
- 如果需要系统级分析，分析结果也不能直接推给无权限订阅者。

### P1：system_prompt 对 admin 仍可完全覆盖

当前限制了非 admin 用户覆盖 system prompt，这是改进。但 admin 请求仍可把 system prompt 完全替换成任意内容。生产环境里 admin 也可能通过被盗 token、误操作或外部调试入口触发 prompt 越权。

建议：

- 生产接口完全移除 `system_prompt`，或仅在 `EMS_AGENT_DEBUG_PROMPT=true` 时允许。
- 即便允许 admin 调试，也应只追加 debug instruction，不应替换安全 policy prompt。

### P1：通用分析和保养审计仍偏演示

保养审计 `MaintenanceAnalyzer.Audit` 仍是固定数字：

- 总检查 45
- 延期任务 3
- 合规率 0.93
- “当前工厂有 3 个保养任务超过预计开始时间 48 小时未启动”

通用分析 `Analyze` 虽然接了 context，但 `AnalyzeData.KeyFindings` 仍固定为“设备处于次健康状态”“维修成本有上升趋势”，Evidence 为空。

建议：

- 保养审计读取真实 `maintenance_tasks`，按计划、工厂、时间范围统计延期/漏检/超期。
- `AnalyzeData` 应从 contextMap 或 analyzer 结果生成，至少返回真实 evidence。

### P1：预测性维护仍保留核心硬编码

`PredictRUL` 仍使用固定：

- `loadFactor := 1.15`
- `avgMTBFHours := 300.0`
- `currentUsedHours := 240.0`

`CalculateTCO` 仍对 `PRESS-05` 特判 12 年使用年限。虽然症状识别的一部分改为真实任务/工单关联，但最后仍对 `equipmentID == 1` 注入 demo finding。

建议：

- RUL 输入改为 runtime snapshots、维修频率、保养执行、设备年龄等真实特征。
- demo fallback 仅在 memory/demo mode 生效，不能在 DB 生产数据中按 ID 注入。
- TCO 使用采购日期/折旧年限计算真实使用年限。

### P1：主动推送前后端协议不完整

后端 `Subscribe` 支持 `webhook_url`，但前端 `agentApi.subscribe` 类型不包含 `webhook_url`，页面也没有输入 webhook URL。结果用户无法通过当前 UI 配置真正的 webhook 投递地址。

另外 UI 中 push types 是：

- `ng_inspection`
- `repair_request`
- `low_stock`

但 `NotifyEvent` 固定查询的是 `predictive_maintenance`。这意味着前端配置的三类订阅不会被当前 `NotifyEvent` 命中。

建议：

- UI 增加 webhook_url/secret 输入。
- 后端和前端统一 push_type 枚举。
- 业务事件真正调用 `NotifyEvent(eventType, targetID, context)`，不要只留方法。

### P1：Tool Registry 形态有了，但协议仍偏薄

Tool Registry 现在集中注册工具，这是正确方向。但 `ToolDefinition` 仍只有 name/description/input_schema，没有 output_schema、错误码、side effect、幂等性、版本、required scopes 等可发现元数据。

建议：

- ToolDefinition 增加 `output_schema`、`required_scopes`、`side_effect`、`idempotent`、`examples`。
- `ListTools` 按 API Key scopes 过滤可见工具。
- `CallToolResponse` 固定 envelope，避免 `content` 任意结构让外部 Agent 难解析。

## 6. 中优先级问题

### P2：新增测试覆盖太浅

新增了：

- `tool/registry_test.go`
- `analyzer/maintenance_test.go`

但 registry 测试没有覆盖 scope 拒绝，因为实现也没有拒绝。maintenance audit 测试只验证固定 mock 返回，不能保护真实业务逻辑。

建议增加：

- API Key hash 认证成功/失败/过期。
- API Key scope 不足时拒绝 `report_repair`。
- 非 admin 维修审计不能读取其他工厂工单。
- `GetEquipmentPrediction` 权限失败时返回明确错误。
- `NotifyEvent` 只投递给 scope 命中的订阅者。

### P2：错误处理仍大量吞掉

示例：

- Chat 保存消息失败被忽略。
- `GetEquipmentPrediction` 忽略三个 analyzer 错误。
- `AuditMaintenance` 创建 session/artifact 失败被忽略。
- `deliverPush` 创建 push log 失败被忽略。

建议：

- 关键持久化用事务。
- 权限错误不要吞掉，否则容易把“无权限”伪装成空结果。

### P2：schema 与迁移仍不完整

GORM model 新增了 `AgentPushLog`、`WebhookURL`、`Secret`、API Key scopes/rate_limit，但 `db/schema.sql` 仍没有完整同步。项目依赖 AutoMigrate 会增加部署不确定性。

建议：

- 补 SQL migration。
- 明确哪些字段是历史兼容新增字段，哪些需要 backfill。

## 7. 上一版建议完成度

| 上一版问题 | 当前状态 | 评价 |
| --- | --- | --- |
| by-id 会话/产物/对话越权 | 已部分修复 | 需要补测试 |
| 预测接口 admin 绕过 | 已部分修复 | 内部错误仍被吞 |
| report_repair 无设备权限校验 | 已修复雏形 | 参数校验仍弱 |
| 知识/技能管理无角色限制 | 已修复雏形 | 角色硬编码，后续可抽权限 |
| API Key 明文存储 | 已修复 | 使用 SHA-256 hash，建议后续加 salt/pepper 或 HMAC |
| API Key scopes/rate limit | 字段已加，未生效 | 仍是 P0 |
| system_prompt 任意覆盖 | 非 admin 已限制 | admin/debug 仍需治理 |
| RepairAudit mock | 部分改为真实扫描 | 工厂隔离缺失 |
| MaintenanceAudit 占位 | 改为固定规则响应 | 仍是演示逻辑 |
| Analyze 占位 | 改为 context + LLM | data/evidence 仍固定/为空 |
| Chat 无证据上下文 | 已改善 | entity extraction 仍粗糙 |
| 主动推送只有配置 | 有投递雏形 | 类型不一致、无 scope 校验、UI 缺 webhook |
| Agent 测试缺失 | 已加少量测试 | 当前构建失败，测试保护不足 |

## 8. 建议下一步修复顺序

### 立即修复

1. 修复 `prompt.go` 裸 `...` 和 `registry.go` 未使用 import。
2. 跑通 `go test ./...`。
3. 增加一个 CI 必跑测试，避免不可编译代码再次进入 main。

### 第一批安全修复

1. API Key scopes 真正接入 AuthMiddleware -> Controller/Service -> ToolRegistry。
2. rate_limit 真正执行。
3. `GetOrdersByFilter` 加工厂过滤，修复维修审计跨工厂读取。
4. `NotifyEvent` 对订阅者 scope 做投递前校验。

### 第二批功能补实

1. 保养审计改为真实 maintenance_tasks 查询。
2. 通用分析返回真实 evidence。
3. RUL/TCO 移除核心硬编码。
4. 前端主动推送补 webhook_url/secret，统一 push_type。

### 第三批工程化

1. ToolDefinition 增加 output_schema/scopes/side_effect/version。
2. agent service 依赖注入，降低全局 DB/config 耦合。
3. 补 migration，减少 AutoMigrate 和 schema.sql 分叉。

## 9. 总体评价

Gemini 这轮方向是对的：它没有只改文档，而是开始把上一版 review 里的安全边界、工具注册、审计逻辑、API Key 安全、主动推送都往代码里落。但实现质量还不稳，最明显的问题是当前后端无法构建；其次是若干“字段/UI/注释已经存在，但安全逻辑没有真正接上”。

如果先把构建修掉，再把 API Key scope/rate-limit 和维修审计工厂隔离补实，这个 Agent 模块会比上一版明显更接近内部可用。现在还不建议把它描述为生产可用的外部 Agent 平台。
