# EMS Agent 系统：架构、能力与集成指南

> EMS 的 Agent 系统不是一个简单的聊天机器人，而是一个**具备领域知识、可自我进化的工业设备运维智能体**。它同时服务于两类使用者：系统内部的运维人员（通过对话界面）和外部的 AI Agent（通过标准化 API）。

---

## 目录

1. [整体架构](#1-整体架构)
2. [内部 Agent：智能运维助手](#2-内部-agent智能运维助手)
3. [外部 Agent 集成：Tool Protocol](#3-外部-agent-集成tool-protocol)
4. [技能系统：可编排的工具链](#4-技能系统可编排的工具链)
5. [预测性分析引擎](#5-预测性分析引擎)
6. [自我进化机制](#6-自我进化机制)
7. [权限与数据隔离](#7-权限与数据隔离)
8. [API 参考](#8-api-参考)

---

## 1. 整体架构

### 1.1 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                     前端 (Vue 3 + Element Plus)              │
│  ManagementAssistantView    AgentIntegrationView            │
│  (对话/审计/知识审核)       (API Key/工具发现/推送配置)       │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTP REST API
┌──────────────────────────▼──────────────────────────────────┐
│                   Controller 层 (Gin)                        │
│         认证校验 → 参数绑定 → 调用 Service → 统一响应         │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                    Service 层 (核心编排)                      │
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │ Policy   │  │ Analyzer │  │  Tool    │  │  Prompt  │   │
│  │ 权限策略  │  │ 规则分析  │  │ 数据获取  │  │ 提示词   │   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
│       │             │             │              │          │
│       ▼             ▼             ▼              ▼          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              LLM Client (OpenAI 兼容)                 │   │
│  │         SiliconFlow / DeepSeek-V3 / 其他              │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │            Repository 层 (双模式存储)                  │   │
│  │     DBAgentRepository (GORM) / MemoryAgentRepository  │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 目录结构

```
backend/internal/agent/
├── controller/agent.go       # 20 个 HTTP 端点
├── service/
│   ├── agent.go              # 核心业务：Chat、Skill 执行、反思学习
│   └── tool_service.go       # 外部 Agent 的工具发现与调用
├── repository/
│   ├── agent.go              # GORM 实现 (PostgreSQL)
│   └── memory_agent.go       # 内存实现 (开发/演示)
├── analyzer/
│   ├── maintenance.go        # 保养分析器
│   ├── repair.go             # 维修审计分析器
│   └── predictive.go         # 预测性分析器 (RUL/TCO/症状/退役)
├── tool/
│   ├── retrieval.go          # 检索工具 (设备档案 + 混合 RAG)
│   ├── maintenance.go        # 保养工具 (合规率/计划查询)
│   └── repair.go             # 维修工具 (故障统计/成本分析)
├── policy/policy.go          # 工厂级数据隔离
├── prompt/prompt.go          # 6 个 LLM 提示词模板
└── dto/agent.go              # 全部请求/响应结构体
```

### 1.3 两条业务路径

Agent 系统有两条核心业务路径，它们共享底层的 Tool、Analyzer、LLM 组件，但面向不同的使用场景：

| 路径 | 触发方式 | 用途 | 特点 |
|------|----------|------|------|
| **结构化分析** | 前端表单 / 外部 API | 保养建议、维修审计、通用分析 | 规则引擎主导，LLM 只做摘要生成 |
| **对话式交互** | Chat 端点 | 多轮问答、设备诊断、决策支持 | 意图识别 → 技能匹配 → 工具编排 → LLM 生成 |

---

## 2. 内部 Agent：智能运维助手

内部 Agent 是面向系统用户的 AI 助手，通过 `ManagementAssistantView.vue` 前端页面提供三种交互模式。

### 2.1 专家对话 (Chat)

用户可以与"智能资产专家"进行多轮对话，聚焦于设备运维领域的决策支持。

**完整调用链：**

```
用户输入 "分析 CNC-001 的健康状态"
    │
    ▼
POST /api/v1/agent/chat
    { "message": "分析 CNC-001 的健康状态" }
    │
    ▼
Controller: 认证 → 加载 User 模型
    │
    ▼
Service.Chat():
    │
    ├── 1. 创建/复用 AgentConversation
    ├── 2. 持久化用户消息到 AgentMessage
    ├── 3. 加载用户个性化经验 (AgentExperience)
    │
    ├── 4. 意图识别与技能匹配
    │   └── MatchSkills("分析 CNC-001 的健康状态")
    │       └── 匹配到技能 "设备深度诊断" (status=active)
    │
    ├── 5. 执行技能 ExecuteSkill()
    │   ├── 解析 Steps JSON:
    │   │   [{step:1, tool:"get_equipment_profile"},
    │   │    {step:2, tool:"get_failure_stats"},
    │   │    {step:3, tool:"predict_remaining_life"},
    │   │    {step:4, tool:"detect_symptoms"}]
    │   │
    │   ├── 从消息中提取设备 ID: "CNC-001" → equipment_id=1
    │   │
    │   ├── Step 1: RetrievalTool.GetEquipmentProfile(1)
    │   │   └── 返回: {name:"CNC加工中心", status:"running", ...}
    │   │
    │   ├── Step 2: RepairTool.GetFailureStats(1)
    │   │   └── 返回: {repair_count:5, mttr:4.2, total_downtime:21}
    │   │
    │   ├── Step 3: PredictiveAnalyzer.PredictRUL(1)
    │   │   └── 返回: {health_score:72, estimated_rul_days:15, ...}
    │   │
    │   ├── Step 4: PredictiveAnalyzer.DetectSymptoms(1)
    │   │   └── 返回: [{type:"micro_stop", severity:"medium", ...}]
    │   │
    │   └── 收集所有 EvidenceItem
    │
    ├── 6. LLM 生成综合分析摘要
    │   └── 注入 Prompt: "用户意图 + 技能描述 + 证据链"
    │   └── 调用 LLM ChatCompletion
    │
    ├── 7. 持久化助手回复
    ├── 8. 异步触发 ReflectAndLearn()  ← 自我进化
    └── 9. 返回 { conversation_id, reply, trace_id }
```

**如果未匹配到技能**，系统退回到标准 LLM 对话模式：
- 注入系统提示词（"顶级工业资产战略专家"角色）
- 加载用户经验上下文
- 取最近 10 条历史消息作为上下文
- 调用 LLM 生成回复

### 2.2 专项审计 (Audit)

面向工程师和管理者的结构化分析能力。

**维修合理性审计 (`POST /agent/audit/repair`)：**

```
请求: { factory_id: 1, equipment_type_id: 2, time_range: {...} }
    │
    ▼
Policy.ValidateScope() → 确认用户有权访问该工厂数据
    │
    ▼
RepairAuditAnalyzer.Analyze()
    ├── 识别重复故障 (repeat_failure)
    ├── 检测维修成本异常
    └── 搜索知识库中的相关标准 (RetrievalTool.SearchManualKnowledge)
    │
    ▼
LLM 生成审计报告摘要
    │
    ▼
持久化: AgentSession → AgentArtifact → AgentEvidenceLink[]

响应: {
    "summary": "发现 3 项维修异常...",
    "risk_level": "high",
    "data": {
        "anomalies": [...],
        "stats": {...},
        "evidence": [...]      ← 每条证据可追溯到具体数据源
    }
}
```

### 2.3 知识审核 (Knowledge)

Agent 从对话中自动提炼的知识草稿，需要人工审核后才能入库。

```
GET /api/v1/agent/knowledges?status=draft
    → 返回待审核的知识列表

PUT /api/v1/agent/knowledge/:id/status
    { "status": "approved" }   ← 确认入库
    { "status": "rejected" }   ← 驳回
```

---

## 3. 外部 Agent 集成：Tool Protocol

EMS 通过标准化的 Tool Protocol 向外部 AI Agent（如 LangChain、AutoGPT、Claude 等）暴露领域能力。这套协议遵循 MCP（Model Context Protocol）的设计理念。

### 3.1 认证方式

外部 Agent 通过 `X-API-KEY` Header 进行认证：

```http
GET /api/v1/agent/tools HTTP/1.1
X-API-KEY: ems_7f8a9b2c...
```

API Key 的特性：
- 以 `ems_` 前缀标识
- 继承创建者用户的角色权限
- 支持设置过期时间（30天/90天/1年/永不过期）
- 创建后仅显示一次

### 3.2 工具发现 (Tool Discovery)

```http
GET /api/v1/agent/tools
```

返回系统当前支持的所有工具及其 JSON Schema 定义：

```json
{
  "tools": [
    {
      "name": "search_equipment",
      "description": "Search for equipment by name, code or model",
      "input_schema": {
        "type": "object",
        "properties": {
          "keyword": {
            "type": "string",
            "description": "Search keyword (name, code, or model)"
          }
        },
        "required": ["keyword"]
      }
    },
    {
      "name": "get_equipment_health",
      "description": "Get real-time health analysis and remaining useful life (RUL) prediction",
      "input_schema": {
        "type": "object",
        "properties": {
          "equipment_id": { "type": "integer" }
        },
        "required": ["equipment_id"]
      }
    },
    {
      "name": "get_spare_part_inventory",
      "description": "Check stock levels of spare parts across different factories",
      "input_schema": {
        "type": "object",
        "properties": {
          "spare_part_id": { "type": "integer" },
          "factory_id": { "type": "integer", "description": "Optional filter" }
        },
        "required": ["spare_part_id"]
      }
    },
    {
      "name": "report_repair",
      "description": "Submit a new repair request for a faulty equipment",
      "input_schema": {
        "type": "object",
        "properties": {
          "equipment_id": { "type": "integer" },
          "fault_description": { "type": "string" },
          "priority": { "type": "integer", "description": "1=High, 2=Medium, 3=Low" }
        },
        "required": ["equipment_id", "fault_description"]
      }
    }
  ]
}
```

### 3.3 工具调用 (Tool Call)

```http
POST /api/v1/agent/tools/call
Content-Type: application/json

{
  "name": "get_equipment_health",
  "arguments": { "equipment_id": 1 }
}
```

响应：
```json
{
  "content": {
    "rul": {
      "equipment_id": 1,
      "health_score": 72.5,
      "estimated_rul_days": 15,
      "reliability": 0.85,
      "risk_factors": ["超负荷运行 (115%)", "近期发生重复故障"],
      "recommendation": "预警：健康分较低，建议在下周内安排油质监测和同心度校准。"
    },
    "tco": {
      "accumulated_repair_cost": 45000,
      "downtime_loss": 120000,
      "depreciated_value": 280000,
      "total_cost_of_ownership": 185000
    },
    "symptoms": [
      {
        "type": "micro_stop",
        "title": "频发微停预警",
        "severity": "medium",
        "description": "设备在近期出现 4 次短时停机..."
      }
    ]
  },
  "is_error": false
}
```

**重要**：`report_repair` 是写操作，会真正创建维修工单。外部 Agent 应谨慎调用。

### 3.4 与 LLM Agent 框架的集成示例

以下是一个外部 LLM Agent 如何利用 EMS Tool Protocol 的典型流程：

```
用户对 ChatGPT 说: "帮我查一下工厂里 CNC 设备的健康状态"
    │
    ▼
外部 LLM Agent (如 ChatGPT + Function Calling):
    │
    ├── 1. GET /api/v1/agent/tools          → 发现可用工具
    ├── 2. POST /api/v1/agent/tools/call    → 调用 search_equipment
    │       { "name": "search_equipment", "arguments": { "keyword": "CNC" } }
    │       → 返回设备列表
    │
    ├── 3. POST /api/v1/agent/tools/call    → 调用 get_equipment_health
    │       { "name": "get_equipment_health", "arguments": { "equipment_id": 1 } }
    │       → 返回健康分析
    │
    └── 4. LLM 综合分析结果，生成自然语言回复
```

### 3.5 主动事件推送

外部 Agent 可以订阅系统事件，当特定条件满足时接收推送通知。

**支持的推送类型：**

| 推送类型 | 触发条件 | 用途 |
|----------|----------|------|
| `ng_inspection` | 现场人员提交点检 NG 结果 | 及时发现设备异常 |
| `repair_request` | 有新的报修工单创建 | 触发外部 Agent 分析 |
| `low_stock` | 备件库存低于安全水位 | 自动触发采购流程 |

**订阅配置：**

```http
POST /api/v1/agent/subscribe
{
  "push_type": "ng_inspection",
  "enabled": true,
  "scope": { "factory_id": 1 }
}
```

**推送触发逻辑**（`NotifyEvent` 方法）：
- 当事件发生时，系统自动执行 `PredictRUL` 预测
- 如果设备 RUL < 7 天，创建 `proactive_push` 类型的 Artifact（"设备停机风险预警"）
- 通知所有匹配 scope 的订阅者

---

## 4. 技能系统：可编排的工具链

技能系统是 Agent 的核心能力载体。一个技能定义了一系列有序的工具调用步骤，用于完成特定的分析任务。

### 4.1 技能定义

```go
// AgentSkill 模型
type AgentSkill struct {
    Name                string   // 技能名称
    Description         string   // 技能描述
    ApplicableTo        []string // 适用设备类型
    ApplicableScenarios []string // 适用场景描述
    Steps               []Step   // 执行步骤 (JSON 数组)
    Version             int      // 版本号
    Status              string   // draft / active / disabled
    UsageCount          int      // 使用次数
    SuccessRate         float64  // 成功率
}
```

**Steps 结构示例：**

```json
[
  {"step": 1, "action": "获取设备档案", "tool": "get_equipment_profile"},
  {"step": 2, "action": "分析故障统计", "tool": "get_failure_stats"},
  {"step": 3, "action": "预测剩余寿命", "tool": "predict_remaining_life"},
  {"step": 4, "action": "检测亚健康征兆", "tool": "detect_symptoms"},
  {"step": 5, "action": "计算总持有成本", "tool": "get_tco_analysis"},
  {"step": 6, "action": "评估退役建议", "tool": "get_retirement_recommendation"}
]
```

### 4.2 内置工具清单

技能 Steps 中可调用的工具：

| 工具名 | 实现层 | 功能 |
|--------|--------|------|
| `get_equipment_profile` | RetrievalTool | 获取设备档案（编号、状态、财务字段） |
| `get_failure_stats` | RepairTool | 故障统计（维修次数、MTTR、总停机时长） |
| `get_cost_analysis` | RepairTool | 成本分析（备件成本 + 人工成本） |
| `get_maintenance_compliance` | MaintenanceTool | 保养合规率（已完成/总任务数） |
| `get_failure_distribution` | RepairAuditAnalyzer | 故障分布分析 |
| `search_manual_knowledge` | RetrievalTool | 混合 RAG 检索（知识库 + 手册） |
| `predict_remaining_life` | PredictiveAnalyzer | RUL 预测 |
| `detect_symptoms` | PredictiveAnalyzer | 亚健康征兆识别 |
| `get_tco_analysis` | PredictiveAnalyzer | 全生命周期总成本计算 |
| `get_retirement_recommendation` | PredictiveAnalyzer | 退役评估 |

### 4.3 技能匹配与执行

当用户发送消息时，系统通过以下流程匹配并执行技能：

```
用户消息: "CNC-001 最近老是出问题，帮我全面分析一下"
    │
    ▼
MatchSkills(): 在数据库中搜索 status='active' 的技能
    ├── SQL: WHERE name ILIKE '%分析%' OR applicable_scenarios ILIKE '%分析%'
    └── 按 success_rate DESC 排序，取第 1 个
    │
    ▼
命中技能: "设备深度诊断" (success_rate: 0.92)
    │
    ▼
extractEquipmentID(): 从消息中提取设备 ID
    ├── 正则匹配设备编码模式 (如 CNC-001)
    └── 或模糊匹配设备名称
    │
    ▼
ExecuteSkill(): 按 Steps 顺序执行
    ├── Step 1: get_equipment_profile → 收集证据
    ├── Step 2: get_failure_stats → 收集证据
    ├── Step 3: predict_remaining_life → 收集证据
    ├── Step 4: detect_symptoms → 收集证据
    │
    ▼
LLM 综合所有证据生成分析摘要
    │
    ▼
记录使用统计: usage_count++, 更新 success_rate
```

### 4.4 技能管理 API

| 操作 | 方法 | 端点 |
|------|------|------|
| 创建技能 | POST | `/api/v1/agent/skills` |
| 技能列表 | GET | `/api/v1/agent/skills?status=active` |
| 技能详情 | GET | `/api/v1/agent/skills/:id` |
| 更新技能 | PUT | `/api/v1/agent/skills/:id` |

---

## 5. 预测性分析引擎

预测性分析是 Agent 系统的核心差异化能力，由 `PredictiveAnalyzer` 实现，包含 4 个分析维度。

### 5.1 RUL 预测（剩余健康寿命）

**算法：**

```
输入:
  - 设备故障统计 (MTBF)
  - 负荷系数 (load_factor, 默认 1.15)
  - 近期故障频率

计算:
  avgMTBFHours = 基准 MTBF (300h)
  if 近期故障 > 3 次:
      avgMTBFHours *= 0.7     // 动态下调

  rulHours = (avgMTBFHours - currentUsedHours) / loadFactor
  healthScore = (rulHours / avgMTBFHours) * 100

输出:
  - health_score: 0-100 健康分
  - estimated_rul_days: 预计剩余天数
  - reliability: 置信度 (0-1.0)
  - risk_factors: 风险因子列表
  - recommendation: 维护建议
```

### 5.2 TCO 计算（全生命周期总成本）

```
TCO = 累计维修费 + 停机损失 + 折旧

其中:
  累计维修费 = Σ(备件成本 + 人工成本)  ← 从 RepairCostDetail 表获取
  停机损失 = 累计停机小时 × 每小时产值损失  ← 从 Equipment.HourlyLoss 获取
  折旧 = 购置价 - 当前净值  ← 直线折旧法
```

### 5.3 亚健康征兆检测

系统通过分析历史维修记录，识别以下亚健康模式：

| 征兆类型 | 检测规则 | 严重程度 |
|----------|----------|----------|
| 频发微停 (micro_stop) | 30 分钟以内的短时停机 ≥ 3 次 | medium |
| 保养无效 (pm_ineffective) | 执行保养后 72 小时内再次故障 | high |
| MTTR 漂移 (mttr_drift) | 近期维修时间显著增加 | medium |

### 5.4 退役评估

```
maintenanceRatio = 累计维修费 / 购置原值

if maintenanceRatio > 0.6:
    decision = "retire"     → 强烈建议退役
elif maintenanceRatio > 0.4:
    decision = "evaluate"   → 列入观察名单
else:
    decision = "continue"   → 继续使用
```

### 5.5 前端展示

右侧面板以可视化方式展示预测结果：
- **健康评分**：仪表盘组件，直观显示 0-100 分
- **RUL 天数**：大数字展示，配合颜色编码（红/黄/绿）
- **TCO 构成**：堆叠柱状图（维修费 + 停机损失 + 折旧）
- **风险征兆**：标签列表，严重程度用颜色区分

---

## 6. 自我进化机制

Agent 的自我进化是系统最核心的设计亮点。通过 `ReflectAndLearn` 方法，Agent 在每次对话后**异步**进行反思学习。

### 6.1 进化闭环

```
对话完成
    │
    ▼ (异步 goroutine)
ReflectAndLearn()
    │
    ├── asyncExtractKnowledge()     ← 知识提取
    │   ├── LLM 分析对话记录
    │   ├── 提取: title, type, summary, details, confidence
    │   └── 创建 AgentKnowledge (status=draft)
    │
    ├── asyncExtractSkill()         ← 技能提炼
    │   ├── LLM 分析对话中的排查逻辑
    │   ├── 提炼: name, description, applicable_scenarios, steps
    │   └── 创建 AgentSkill (status=draft)
    │
    └── asyncCollectExperience()    ← 经验收集 (待实现)
        └── 从对话中收集用户偏好
    │
    ▼
人工审核 (ManagementAssistantView → 知识审核)
    ├── 管理员查看草稿: 标题、类型、置信度
    ├── 确认入库 (approved) 或 驳回 (rejected)
    │
    ▼
已审核的知识/技能被激活
    ├── 知识: 在 RAG 检索中被优先召回 (权重 1.0)
    └── 技能: 在 Chat 中被自动匹配执行
    │
    ▼
Agent 能力增强 → 下次对话表现更好
```

### 6.2 知识提取细节

**提示词模板** (`BuildKnowledgeExtractionPrompt`)：
- 要求 LLM 从对话中提取有价值的工业知识
- 输出 JSON 格式: `{title, type, summary, details, confidence}`
- 知识类型: 故障根因、预防措施、最佳实践等

**提取结果示例：**

```json
{
  "title": "CNC 主轴轴承异响的早期识别方法",
  "type": "fault_diagnosis",
  "summary": "当 CNC 主轴在高速运转时出现周期性异响，且频率与转速成正比，通常是轴承内圈磨损的早期征兆...",
  "details": { "symptoms": [...], "root_cause": "...", "prevention": "..." },
  "confidence": 0.85
}
```

### 6.3 技能提炼细节

**提示词模板** (`BuildSkillExtractionPrompt`)：
- 要求 LLM 从对话中提炼通用的排查套路
- 输出 JSON 格式: `{name, description, applicable_scenarios, steps}`
- Steps 中的 tool 字段必须是系统内置工具名

### 6.4 经验衰减机制

`AgentExperience` 模型支持时间衰减，确保过时的经验逐渐降权：

```go
// 衰减公式
weight = weight * (1 - decay_rate)

// 示例
初始权重: 1.0, 衰减率: 0.1
第 1 次衰减: 1.0 * 0.9 = 0.9
第 2 次衰减: 0.9 * 0.9 = 0.81
第 3 次衰减: 0.81 * 0.9 = 0.729
...
```

衰减后的经验在 Chat 时被注入到上下文中，影响 LLM 的回复风格和关注点。

### 6.5 使用统计追踪

每次 Agent 调用都会记录 `AgentUsage`：

```go
type AgentUsage struct {
    SessionID       uint    // 关联会话
    UserID          uint    // 调用用户
    Scenario        string  // 场景类型
    Model           string  // 使用的 LLM 模型
    PromptTokens    int     // 输入 token 数
    CompletionTokens int    // 输出 token 数
    ResponseTimeMs  int64   // 响应时间 (毫秒)
}
```

这些数据可用于：
- 成本分析（token 消耗）
- 性能监控（响应时间）
- 使用趋势分析

---

## 7. 权限与数据隔离

### 7.1 工厂级作用域

`PolicyService` 实现了严格的工厂级数据隔离：

```
Admin 用户
    └── 可以访问所有工厂的数据

其他角色用户
    └── 只能访问自己所属工厂 (factory_id) 的数据
        └── 请求其他工厂数据 → FORBIDDEN_SCOPE 错误

无工厂绑定的非 Admin 用户
    └── 完全拒绝访问
```

### 7.2 双重防线

1. **Policy 层**：在 Service 入口处校验 `ValidateScope()`
2. **Tool 层**：每个 Tool 内部的 `checkPermission()` 方法

```go
// Tool 层的权限检查示例
func (t *RepairTool) checkPermission(equipmentID uint, user model.User) error {
    if user.Role == "admin" || user.FactoryID == nil {
        return nil
    }
    // 检查设备是否属于用户所在工厂
    e, _ := repo.GetByID(equipmentID)
    if e.Workshop.FactoryID != *user.FactoryID {
        return fmt.Errorf("access denied: equipment belongs to another factory")
    }
    return nil
}
```

### 7.3 证据链追溯

每个分析产出物（`AgentArtifact`）都通过 `AgentEvidenceLink` 关联到具体的数据源：

```
AgentArtifact (分析报告)
    │
    ├── AgentEvidenceLink: type="knowledge", source_table="knowledge_articles", source_id=42
    ├── AgentEvidenceLink: type="manual", source_table="equipment_manual_chunks", source_id=15
    └── AgentEvidenceLink: type="repair_record", source_table="repair_orders", source_id=201
```

这确保了"每个结论都有据可查"，用户可以追溯到原始数据。

---

## 8. API 参考

### 8.1 对话 API

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/agent/chat` | 发送对话消息 |
| GET | `/agent/conversations` | 会话列表 |
| GET | `/agent/conversations/:id` | 会话详情（含消息） |

**Chat 请求体：**

```json
{
  "conversation_id": 1,           // 可选，不传则创建新会话
  "message": "分析 CNC-001 的健康状态",
  "context": { "page": "equipment_detail" },  // 可选，补充上下文
  "system_prompt": "..."          // 可选，自定义系统提示词
}
```

**Chat 响应：**

```json
{
  "conversation_id": 1,
  "reply": "根据分析，CNC-001 当前健康评分为 72 分...",
  "trace_id": "tr_abc123",
  "artifact_id": 42,
  "suggested_actions": ["查看维修历史", "创建保养计划"]
}
```

### 8.2 分析 API

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/agent/maintenance/recommend` | 保养优化建议 |
| POST | `/agent/audit/repair` | 维修合理性审计 |
| POST | `/agent/audit/maintenance` | 保养计划审计 |
| POST | `/agent/analyze` | 通用分析 |
| GET | `/agent/equipment/:id/prediction` | 设备预测（RUL+TCO+症状） |

### 8.3 知识与技能 API

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/agent/knowledges` | 知识列表 |
| PUT | `/agent/knowledge/:id/status` | 审核知识 |
| GET | `/agent/skills` | 技能列表 |
| POST | `/agent/skills` | 创建技能 |
| GET | `/agent/skills/:id` | 技能详情 |
| PUT | `/agent/skills/:id` | 更新技能 |

### 8.4 外部 Agent API

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/agent/tools` | 工具发现 |
| POST | `/agent/tools/call` | 工具调用 |
| POST | `/agent/subscribe` | 推送订阅 |
| GET | `/agent/subscriptions` | 订阅列表 |

### 8.5 会话历史 API

| 方法 | 端点 | 说明 |
|------|------|------|
| GET | `/agent/sessions` | 会话列表 |
| GET | `/agent/sessions/:id` | 会话详情 |
| GET | `/agent/artifacts/:id` | 产出物详情 |

### 8.6 统一响应格式

**成功响应 (`AgentResponseEnvelope`)：**

```json
{
  "success": true,
  "trace_id": "tr_abc123",
  "language": "zh",
  "scenario": "maintenance_recommendation",
  "scope_summary": { "factory_id": 1 },
  "summary": "建议缩短保养周期...",
  "risk_level": "medium",
  "artifact_id": 42,
  "evidence_count": 3,
  "data": { ... }
}
```

**错误响应 (`AgentErrorEnvelope`)：**

```json
{
  "success": false,
  "trace_id": "tr_abc123",
  "error": {
    "code": "FORBIDDEN_SCOPE",
    "Message": "access denied: equipment belongs to another factory"
  }
}
```

---

## 附录：技术栈依赖

| 组件 | 技术 | 说明 |
|------|------|------|
| LLM 客户端 | OpenAI 兼容 API | 支持 SiliconFlow/DeepSeek-V3，超时 60s |
| 数据库 | PostgreSQL + GORM | Agent 专属表 10+ 张 |
| 缓存 | Redis | 可选，用于会话缓存 |
| 前端 | Vue 3 + Element Plus | 管理助手 + 集成管理两个页面 |
| HTTP 框架 | Gin | 统一中间件链 |
