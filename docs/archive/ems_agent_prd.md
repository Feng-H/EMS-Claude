---
title: EMS Agent 产品需求文档（PRD）
tags: [format-article, technique-llm, system-eam, tool, concept]
---

# EMS Agent 产品需求文档（PRD）

> 版本：v1.0
> 日期：2026-04-25
> 状态：草案
> 关联文档：
> - [Phase 1 设计](../ems_agent_design.md)
> - [Roadmap](./ems_agent_roadmap.md)
> - [Schema 设计](./ems_agent_schema_design.md)
> - [Phase 2 Schema 扩展](./ems_agent_phase2_schema_extension.md)
> - [API 设计](./ems_agent_api_design.md)
> - [Phase 2 API 扩展](./ems_agent_phase2_api_design.md)
> - [Phase 1 任务拆解](./ems_agent_phase1_task_breakdown.md)

---

## 1. 产品定位

### 1.1 一句话定义

EMS Agent 是一个嵌入设备管理系统的、具有学习能力的数据分析智能体——它通过系统数据和工程师对话两个渠道获取信息，持续积累分析方法和领域知识，越用越聪明。

### 1.2 与通用 AI 助手的区别

| 维度 | 通用 AI 助手 | EMS Agent |
|------|-------------|-----------|
| **领域** | 开放领域 | 设备管理专精 |
| **数据来源** | 以用户输入为主 | 嵌入业务系统，系统数据是主食 |
| **记忆** | 会话级或有限记忆 | 三层持久化存储（技能+知识+经验） |
| **成长** | 依赖模型升级 | 自主学习，越用越强 |
| **主动性** | 被动响应 | 主动发现、主动推送 |
| **协作** | 单向输出 | 人机闭环，工程师校准行为 |

### 1.3 核心特质

1. **领域专精**：只做设备管理，拥有专业基座（维修策略、故障分析方法、指标体系），不越界
2. **数据原生**：嵌在业务系统里，设备台账、点检、维修、保养、备件、运行数据都是它的数据源
3. **会学习**：技能库+知识库+经验库三层沉淀，不是每次从零开始
4. **人机协作**：工程师校准，不是单机跑；人在环里保证质量
5. **能成长**：用得越多越聪明，部署日不是能力上限

## 2. 目标用户

### 2.1 主要用户：设备工程师

- **使用频率**：日常使用，每日/每周
- **核心诉求**：快速定位设备问题、获得可执行的维修方案、掌握设备状态全貌
- **使用场景**：
  - 接到报修单时，快速了解设备"前世今生"
  - 定期分析车间/工厂设备健康状态
  - 通过对话深挖特定问题
  - 校准 Agent 的分析结论和技能

### 2.2 辅助用户：管理者（工厂厂长/设备主管）

- **使用频率**：周期性使用，月度/季度
- **核心诉求**：了解设备管理整体状况、异常趋势、资源分配合理性
- **使用场景**：
  - 查看 Agent 推送的周期性分析报告
  - 按需查询特定设备或车间的管理状况
  - 审阅 Agent 的分析结论和改进建议

### 2.3 用户画像与 Hermes Agent 的对比

| 维度 | Hermes Agent 的用户（黄峰） | EMS Agent 的用户（设备工程师） |
|------|--------------------------|------------------------------|
| **专业性** | 跨领域，写作+技术+运维 | 设备管理垂直领域 |
| **交互深度** | 深度对话，探索式 | 任务导向，解决问题 |
| **时间投入** | 有时间深度思考 | 被问题追着跑，时间碎片化 |
| **期望输出** | 文章、方案、洞察 | 判断依据、行动建议 |
| **反馈意愿** | 主动纠正、讨论 | 被动确认，精力有限 |

## 3. 能力模型

### 3.1 能力层级

```text
┌─────────────────────────────────────────┐
│  L4: 主动洞察 — 发现人没想到要看的东西    │
├─────────────────────────────────────────┤
│  L3: 方法积累 — 沉淀分析思路为可复用技能  │
├─────────────────────────────────────────┤
│  L2: 跨域关联 — 组合多数据域做综合判断    │
├─────────────────────────────────────────┤
│  L1: 数据查询 — 从系统数据中提取信息      │
└─────────────────────────────────────────┘
```

- **L1（数据查询）**：Phase 1 已具备。通过 Tool 层查询设备、维修、保养、备件等数据。
- **L2（跨域关联）**：Phase 2 核心。跨数据域组合分析，例如：维修记录×保养执行×点检NG→发现保养遗漏与故障的关联。
- **L3（方法积累）**：Phase 2 高阶。从成功的分析过程中提炼可复用的分析路径，存入技能库。
- **L4（主动洞察）**：Phase 2 目标。基于历史经验和当前数据，主动发现异常模式并推送。

### 3.2 专业基座

Agent 内置设备管理领域的专业知识框架，这些是它的"常识"：

**设备生命周期知识**
- 浴缸曲线（早期故障期→偶然故障期→耗损故障期）
- 设备老化规律与维护策略的关系
- 不同设备类型的典型故障模式

**维修策略体系**
- 事后维修（Breakdown Maintenance）
- 预防性维修（Preventive Maintenance）
- 预测性维修（Predictive Maintenance）
- 以可靠性为中心的维修（RCM）
- TPM（全员生产维护）理念

**故障分析方法**
- 鱼骨图（因果分析）
- 5Why 分析法
- FTA（故障树分析）
- PM 分析（现象-机理分析）

**核心指标体系**
- MTBF（平均故障间隔时间）
- MTTR（平均修复时间）
- OEE（设备综合效率）
- 可用率、利用率、性能率
- 保养完成率、点检完成率

**备件管理逻辑**
- ABC 分类法
- 安全库存模型
- 经济订货量
- 备件与故障的关联分析

专业基座通过 System Prompt 注入，确保 Agent 的分析视角和术语符合设备工程领域的标准。

### 3.3 能力边界

#### 应该做的

- 设备故障分析与诊断
- 维修策略评估与优化建议
- 备件消耗分析与异常预警
- 保养计划合理性检视
- 设备健康状态综合评估
- 故障模式识别与根因分析
- 跨设备/跨时间的趋势对比
- 行业基准数据解读

#### 不应该碰的

- 生产排程优化（MES 职责）
- 质量缺陷分析（QMS 职责）
- 采购比价与供应商管理（ERP 职责）
- 人员排班与绩效（HR 职责）
- 任何与设备管理无关的通用问答

当用户提出超出边界的问题时，Agent 应识别并说明："这个问题超出设备管理的范围，建议使用 xxx 系统/工具处理。"

## 4. 三层存储体系

这是 Agent 实现"越用越聪明"的核心机制。借鉴 Hermes Agent 的 memory + skills + wiki 三层架构，针对设备管理场景设计。

### 4.1 总览

| 层 | 名称 | 存什么 | 怎么产生 | 怎么用 | 生命周期 |
|---|---|---|---|---|---|
| **技能库** | Skill Store | 分析方法、排查流程、诊断框架 | 对话提炼 + 主动归纳 + 工程师创建 | 分析时匹配执行步骤 | 长期有效，持续迭代 |
| **知识库** | Knowledge Wiki | 分析结论、故障案例、设备画像、规律 | 分析结果沉淀 + 系统数据提炼 | 分析时注入上下文 | 会过期，需更新 |
| **经验库** | Experience Store | 工程师修正、偏好、校准反馈、行为模式 | 工程师交互反馈 + 满意度统计 | 调整 Agent 行为权重 | 需衰减，需校准 |

三层的流动关系：

```text
对话/分析 → 产出结论 → 沉淀知识库
         → 提炼方法 → 存入技能库
         → 收集反馈 → 更新经验库

分析执行时：技能库决定怎么分析
            知识库补充历史上下文
            经验库校准行为偏好
```

### 4.2 技能库（Skill Store）

#### 4.2.1 什么是技能

技能是一套**可执行的分析方法**。它定义了面对特定类型的问题时，应该按照什么步骤、调用什么工具、做什么判断。

与知识的区别：
- **知识**：关于某个具体设备/车间/时间段的分析结论（"3号CNC主轴轴承批次BT20250803有问题"）
- **技能**：分析方法的通用框架（"如何排查备件批次质量问题的系统化方法"）

技能可以跨设备、跨时间复用。

#### 4.2.2 技能数据结构

```json
{
  "id": "spare_part_batch_investigation",
  "name": "备件批次质量排查",
  "description": "当某类设备故障率异常上升时，排查是否与特定批次备件相关的系统化方法",

  "applicable_to": [
    "fault_rate_anomaly",
    "repeated_failure",
    "spare_part_consumption_anomaly"
  ],

  "applicable_scenarios": [
    "同类型设备故障率异常上升",
    "特定备件消耗量突然放大",
    "维修后短期内重复故障"
  ],

  "scope": {
    "domain": "spare_part",
    "requires_tools": [
      "get_failure_distribution",
      "get_spare_part_consumption_detail",
      "get_equipment_runtime_trend"
    ]
  },

  "steps": [
    {
      "step": 1,
      "action": "按设备类型和时间段筛选故障记录，统计故障率趋势",
      "tool": "get_failure_distribution",
      "params": { "group_by": "spare_part_batch", "time_range": "last_90_days" },
      "expected": "确认故障率上升的起始时间和幅度",
      "decision_point": null
    },
    {
      "step": 2,
      "action": "调取故障设备使用的备件批次信息",
      "tool": "get_spare_part_consumption_detail",
      "params": { "merge_with": "repair_records", "time_range": "last_90_days" },
      "expected": "找出共同使用的备件批次号",
      "decision_point": null
    },
    {
      "step": 3,
      "action": "对比不同批次备件的故障率差异",
      "tool": "compare_failure_rate_by_batch",
      "params": { "threshold": 0.3 },
      "expected": "批次间差异是否超过30%",
      "decision_point": {
        "condition": "批次间故障率差异 > 30%",
        "if_true": "初步确认批次相关性，进入第4步",
        "if_false": "批次因素排除，考虑其他原因（运行环境、操作变更等）"
      }
    },
    {
      "step": 4,
      "action": "排除运行时长变化的干扰因素",
      "tool": "get_equipment_runtime_trend",
      "params": { "compare_with": "peer_equipment" },
      "expected": "确认故障率上升不是单纯因为使用强度增加",
      "decision_point": null
    }
  ],

  "meta": {
    "version": 3,
    "created_at": "2025-12-15T10:00:00Z",
    "created_by": "agent:auto",
    "usage_count": 12,
    "success_rate": 0.83,
    "last_used": "2026-03-20T14:30:00Z",
    "status": "active"
  },

  "changelog": [
    {
      "version": 1,
      "source": "agent:auto",
      "note": "从对话记录提炼：工程师排查heater故障的路径",
      "timestamp": "2025-12-15T10:00:00Z"
    },
    {
      "version": 2,
      "source": "engineer:张三",
      "note": "补充第4步：排除运行时长干扰，避免误判",
      "timestamp": "2026-01-20T16:00:00Z"
    },
    {
      "version": 3,
      "source": "agent:self_improve",
      "note": "标准化工具调用参数，统一时间范围默认值",
      "timestamp": "2026-02-28T09:00:00Z"
    }
  ]
}
```

#### 4.2.3 技能的产生路径

**路径一：对话中自动提炼（主要路径）**

工程师通过对话引导 Agent 完成一次分析。对话结束后，Agent 进行反思：

1. 回顾分析路径：初始分析方向 vs 最终被工程师引导到的方向
2. 识别关键转折点：工程师做了什么修正，引入了什么新视角
3. 提炼通用步骤：将具体操作抽象为可复用的步骤序列
4. 生成技能草稿（status: `draft`）

工程师在下次使用时确认，draft → `active`。

**路径二：从成功案例中自动归纳**

Agent 完成一次分析，工程师确认"结论有用"（通过反馈机制）。Agent 回溯分析过程：
1. 提取使用的工具调用序列
2. 识别哪些步骤对得出结论是关键的
3. 结构化为技能模板
4. 生成技能草稿

**路径三：工程师直接创建（辅助路径）**

高级工程师有成熟的诊断方法，可以通过界面直接创建技能。适合将隐性经验显性化。

#### 4.2.4 技能的匹配与执行

当 Agent 收到分析任务时：

1. **意图识别**：从用户输入中提取分析目标
2. **技能检索**：根据意图 + 设备上下文，检索匹配的技能（基于 `applicable_to` 和 `applicable_scenarios`）
3. **技能排序**：按 `success_rate` × `recency_score` 排序
4. **执行**：按技能定义的步骤序列执行工具调用
5. **灵活调整**：如果执行过程中遇到 unexpected 数据，可以偏离技能步骤，但记录偏离原因

如果匹配不到已有技能，Agent 退回到通用分析模式（L2 能力），并标记这次分析为潜在的技能来源。

#### 4.2.5 技能的迭代与淘汰

- **自动优化**：每次使用后，根据实际效果更新 `success_rate` 和 `usage_count`
- **人工修正**：工程师可以在技能详情页编辑步骤
- **自动合并**：检测到功能相似的技能时，提示工程师合并
- **淘汰机制**：`success_rate` 持续低于 0.4 且 90 天未使用的技能，自动标记为 `deprecated`

### 4.3 知识库（Knowledge Wiki）

#### 4.3.1 什么是知识

知识是分析过程中产出的**有价值的结论和发现**。每条知识是一个独立的知识条目，可被检索和引用。

与技能的区别：
- 技能是"怎么做"（方法论），知识是"发现了什么"（结论）
- 技能跨场景复用，知识针对具体对象
- 技能长期有效，知识会过时

#### 4.3.2 知识数据结构

```json
{
  "id": "k_20251215_001",
  "title": "三车间B线CNC冬季heater故障与备件批次相关",
  "type": "root_cause_analysis",
  "scope": {
    "factory_id": 3,
    "workshop_id": 8,
    "equipment_type": "CNC",
    "equipment_ids": ["EQ-003-021", "EQ-003-023", "EQ-003-025"],
    "time_period": { "start": "2025-11-01", "end": "2026-02-28" }
  },

  "summary": "三车间B线3台CNC在2025年11月至2026年2月间 heater 故障率上升47%，根因为批次BT20250803的加热管质量缺陷。已更换批次后故障率恢复正常。",

  "details": {
    "evidence": [
      "维修记录：heater类设备11-2月故障单数环比+47%",
      "备件记录：故障heater均使用批次BT20250803的加热管",
      "对比数据：使用批次BT20250805的同类设备无异常"
    ],
    "root_cause": "批次BT20250803加热管存在热稳定性缺陷",
    "resolution": "更换为批次BT20250805，3月后故障率回落至正常水平",
    "prevention": "建议对批次BT20250803进行隔离检查，通知供应商"
  },

  "related_skill_id": "spare_part_batch_investigation",

  "meta": {
    "created_at": "2025-12-15T10:00:00Z",
    "created_by": "agent:analysis_session_1024",
    "verified_by": "engineer:张三",
    "verified_at": "2026-03-15T14:00:00Z",
    "status": "confirmed",
    "confidence": 0.92,
    "referenced_count": 5,
    "last_referenced": "2026-04-10T09:00:00Z",
    "expire_at": null
  }
}
```

#### 4.3.3 知识类型

| type | 说明 | 示例 |
|------|------|------|
| `root_cause_analysis` | 根因分析结论 | "heater故障根因是备件批次缺陷" |
| `pattern` | 周期性规律 | "每年11-2月 heater 故障率上升15%" |
| `equipment_profile` | 设备画像 | "3号CNC近6个月主轴相关故障占75%" |
| `maintenance_insight` | 保养有效性洞察 | "二级保养执行率低于80%的设备故障率是正常水平的2.3倍" |
| `spare_part_insight` | 备件消耗洞察 | "某型号密封圈消耗量同比上升200%，排查发现安装工艺变更" |
| `benchmark_comparison` | 对标分析 | "A车间CNC设备MTBF比B车间低40%，差距主要来自预防性保养完成率" |

#### 4.3.4 知识的生命周期

```text
生成（draft） → 工程师确认（confirmed） → 引用和更新 → 过期/归档
                                      → 工程师否认（rejected）
```

- **自动生成**：分析结束后，如果结论满足一定置信度阈值，自动生成 draft 知识条目
- **工程师确认**：工程师查看后确认或修正
- **引用**：后续分析中自动检索相关知识注入上下文
- **自动过期**：超过 `expire_at` 或 180 天未被引用的知识，降低检索权重
- **主动归档**：工程师可以手动归档过时知识

### 4.4 经验库（Experience Store）

#### 4.4.1 什么是经验

经验是 Agent 在交互过程中收集到的**行为校准信息**。它不包含具体的分析内容，而是记录"什么有用、什么没用、用户偏好什么"。

#### 4.4.2 经验数据结构

```json
{
  "id": "exp_20260120_001",
  "type": "preference",
  "category": "analysis_depth",

  "content": {
    "observation": "工程师张三在查看备件分析时，总是要求看同比数据",
    "preference": "备件相关分析默认包含同比和环比两个维度",
    "scope": {
      "user_id": 5,
      "domain": "spare_part"
    }
  },

  "meta": {
    "source": "engineer_feedback",
    "created_at": "2026-01-20T16:00:00Z",
    "weight": 0.8,
    "decay_rate": 0.05,
    "last_applied": "2026-04-15T10:00:00Z",
    "status": "active"
  }
}
```

#### 4.4.3 经验类型

| type | 说明 | 示例 |
|------|------|------|
| `preference` | 工程师偏好 | "默认包含同比数据"、"输出要简洁" |
| `correction` | 行为修正 | "不要用故障次数排序，用停机时长排序" |
| `boundary` | 边界确认 | "这个方向不需要，关注备件即可" |
| `quality_feedback` | 输出质量反馈 | "这次分析有用"、"结论不准确" |

#### 4.4.4 经验的衰减机制

- **偏好类**：每 90 天衰减一次，如果期间没有同类反馈，权重降低 20%
- **修正类**：每 60 天衰减一次，衰减更快，因为情况可能已变化
- **边界类**：不自动衰减，需要工程师手动更新
- **质量反馈**：每次使用后刷新，长期不使用自然降权

衰减公式：`effective_weight = weight × (1 - decay_rate) ^ (days_since_last_applied / decay_cycle)`

## 5. 交互模式

### 5.1 两种触发方式

#### 5.1.1 主动推送（Agent 发起）

Agent 定期或基于事件自动分析数据，将结果推送给工程师。

**定时巡检**
- 频率：每日凌晨自动运行
- 范围：根据工程师关注的车间/设备范围
- 输出：发现异常或值得注意的变化，推送摘要
- 不打扰原则：只推送有实质发现的内容，没有异常则不推送

**事件触发**
- 触发条件：报修单创建、保养任务逾期、备件低于安全库存等
- 输出：与事件相关的设备上下文和历史分析
- 场景：工程师打开报修单时，自动展示设备诊断摘要

**周期报告**
- 频率：月度/季度
- 对象：管理者（工厂厂长）
- 内容：设备健康状况总览、异常趋势、改进建议

#### 5.1.2 对话分析（用户发起）

工程师通过自然语言对话，向 Agent 提出分析需求。

**典型对话场景**：

```text
工程师：三车间最近报修特别多，帮我看看什么情况
Agent：正在分析三车间近30天报修数据...

[Agent 自动匹配到相关技能，执行分析，注入相关知识]

分析结果：
1. 报修总量：较上月上升62%（47单→76单）
2. 主要集中在B线CNC设备，占总量的68%
3. B线6台CNC中有4台在近30天内报修，主轴相关故障占75%
4. 对比发现：这4台设备上一次保养（3月15日）跳过了"主轴轴承检查"项
5. 点检记录：3月20日、4月2日均报"主轴温度偏高"NG

建议：优先检查4台设备的主轴轴承状态，确认备件仓NSK 6205库存。

相关经验：去年11月 heater 类设备也出现类似模式（保养遗漏→点检NG→集中故障）
```

### 5.2 交互设计原则

1. **结论先行**：每条输出以结论开头，数据作为支撑，不是反过来
2. **可操作性**：给出的建议要具体到"检查什么、联系谁、带什么备件"，不是泛泛而谈
3. **证据可追溯**：每个结论都关联到具体的数据来源
4. **简洁克制**：工程师时间有限，不做无关延伸
5. **渐进展开**：先给摘要，工程师需要深挖时再展开细节

### 5.3 反馈机制

工程师对 Agent 输出的反馈是人机闭环的关键接口。

**轻量反馈（不打断工作流）**
- 每条输出附带快捷操作：👍 有用 / 👎 没用 / 💡 补充
- 点击"有用"：知识条目 `confirmed`，技能 `success_rate` +1
- 点击"没用"：触发简短原因选择（方向不对 / 数据不准确 / 结论太浅 / 其他）
- 点击"补充"：输入修正或补充内容

**深度反馈（在技能/知识详情页）**
- 工程师可以编辑技能步骤
- 工程师可以修正知识结论
- 工程师可以查看和调整自己的经验偏好

## 6. Agent 与 Analyzer 的关系

### 6.1 分工

```text
┌─────────────────────────────────────────────────┐
│  Analyzer（现有）                               │
│  - 固定分析维度，预定义规则                     │
│  - 定时跑，输出格式固定                         │
│  - 例如：维修审核异常检测、保养完成率统计       │
│  - 是基础计算能力                               │
├─────────────────────────────────────────────────┤
│  Agent（新增）                                  │
│  - 灵活的分析，根据场景组合 Analyzer 结果       │
│  - 从数据和对话两个来源学习                     │
│  - 能把一次分析过程沉淀为可复用的技能           │
│  - 是上层智能                                   │
└─────────────────────────────────────────────────┘
```

### 6.2 Agent 如何使用 Analyzer

Agent 可以将 Analyzer 的输出作为工具调用的一部分：

- `run_audit_checks` → 调用维修审核 Analyzer，获取异常检测结果
- `run_maintenance_recommendation` → 调用保养建议 Analyzer，获取建议结果
- `get_analytics_summary` → 调用统计分析模块，获取 MTBF/MTTR/OEE 等指标

Analyzer 负责"算"，Agent 负责"判断和表达"。

## 7. 数据来源

### 7.1 系统内部数据（主要数据源）

Agent 与 EMS 系统共享数据库，可以直接访问以下数据域：

| 数据域 | 表/模块 | Agent 能分析出什么 |
|--------|---------|-------------------|
| 设备台账 | equipment, equipment_types | 设备年龄分布、同型号集中度、归属关系 |
| 点检记录 | inspection_tasks, inspection_records | NG趋势、早期征兆、执行质量 |
| 维修记录 | repair_orders, repair_logs | 故障频率、维修成本、人员能力、重复故障 |
| 保养记录 | maintenance_plans, maintenance_tasks | 执行率、遗漏项、与故障的关联 |
| 备件消耗 | spare_parts, spare_part_consumptions | 消耗趋势、批次质量、库存风险 |
| 运行数据 | equipment_runtime_snapshots | 利用率、OEE、运行强度与故障关联 |
| 知识库 | knowledge_articles | 最佳实践匹配、经验案例检索 |
| 手册文档 | equipment_manual_documents, equipment_manual_chunks | 标准作业程序、厂家建议 |

### 7.2 对话数据（辅助数据源）

工程师与 Agent 的对话也是重要的数据输入：
- 分析意图和关注方向
- 工程师的领域知识和经验
- 对 Agent 输出的修正和补充

### 7.3 数据安全与权限

- Agent 的所有数据访问严格遵循 EMS 现有的权限体系
- 工程师只能看到自己权限范围内的数据
- 跨工厂/车间数据访问由 Policy 层统一拦截
- Agent Session 记录完整的数据访问链路，支持审计

## 8. 技术架构

### 8.1 整体架构

```text
┌─────────────────────────────────────────────────────┐
│                    用户界面层                       │
│  ┌───────────┐  ┌───────────┐  ┌──────────────────┐ │
│  │ 对话界面   │  │ 推送通知   │  │ 报告/分析页       │ │
│  └─────┬─────┘  └─────┬─────┘  └────────┬─────────┘ │
└────────┼──────────────┼──────────────────┼───────────┘
         │              │                  │
         ▼              ▼                  ▼
┌─────────────────────────────────────────────────────┐
│                   API 层 (Gin)                      │
│  /api/v1/agent/chat       对话接口                  │
│  /api/v1/agent/push       推送管理                  │
│  /api/v1/agent/skills     技能管理                  │
│  /api/v1/agent/knowledge  知识管理                  │
│  /api/v1/agent/experience 经验管理                  │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────┴──────────────────────────────┐
│                  Agent 核心引擎                     │
│                                                      │
│  ┌────────────┐  ┌────────────┐  ┌───────────────┐  │
│  │ 对话管理器   │  │ 技能调度器  │  │ 推送引擎      │  │
│  │ Session Mgr │  │ Skill Sched│  │ Push Engine  │  │
│  └─────┬──────┘  └─────┬──────┘  └───────┬───────┘  │
│        │               │                  │          │
│  ┌─────┴───────────────┴──────────────────┴───────┐  │
│  │              LLM 推理层                         │  │
│  │     (外部API: DeepSeek / OpenAI / 兼容接口)     │  │
│  │  - 意图识别  - 技能匹配  - 结果组织  - 技能提炼 │  │
│  └─────────────────────┬──────────────────────────┘  │
│                        │                              │
│  ┌─────────────────────┴──────────────────────────┐  │
│  │              工具层 (Tools)                    │  │
│  │  ┌──────────┐ ┌──────────┐ ┌────────────────┐  │  │
│  │  │ Analyzer │ │ 数据查询  │ │ 知识检索       │  │  │
│  │  │ 调用     │ │ Tools    │ │ (全文+语义)    │  │  │
│  │  └──────────┘ └──────────┘ └────────────────┘  │  │
│  └─────────────────────────────────────────────────┘  │
│                                                      │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐             │
│  │ Policy   │ │ 学习引擎  │ │ 经验引擎  │             │
│  │ 权限控制  │ │ Learning │ │ Experience│            │
│  └──────────┘ └──────────┘ └──────────┘             │
└──────────────────────────────────────────────────────┘
                       │
┌──────────────────────┴──────────────────────────────┐
│                    存储层                           │
│  ┌──────────┐ ┌──────────┐ ┌──────────────────────┐ │
│  │ PostgreSQL│ │  Redis   │ │ 三层存储              │ │
│  │ 业务数据   │ │ 会话缓存  │ │ ┌────────────────┐  │ │
│  │ Agent数据  │ │ 推送队列  │ │ │ 技能库          │  │ │
│  │           │ │          │ │ │ Skill Store     │  │ │
│  │           │ │          │ │ ├────────────────┤  │ │
│  │           │ │          │ │ │ 知识库          │  │ │
│  │           │ │          │ │ │ Knowledge Wiki  │  │ │
│  │           │ │          │ │ ├────────────────┤  │ │
│  │           │ │          │ │ │ 经验库          │  │ │
│  │           │ │          │ │ │ Experience Store│  │ │
│  │           │ │          │ │ └────────────────┘  │ │
│  └──────────┘ └──────────┘ └──────────────────────┘ │
└──────────────────────────────────────────────────────┘
```

### 8.2 模块职责

在现有 `internal/agent/` 结构基础上扩展：

| 模块 | 现有 | 新增职责 |
|------|------|---------|
| `controller/` | API handler | 新增：对话、技能管理、知识管理、推送相关 API |
| `service/` | 编排分析流程 | 新增：对话管理、技能调度、学习触发 |
| `analyzer/` | 规则分析 | 保持：作为工具被 Agent 调用 |
| `tool/` | 数据查询工具 | 扩展：增加技能/知识/经验的读写工具 |
| `policy/` | 权限控制 | 扩展：技能和知识的访问权限 |
| `prompt/` | Prompt 模板 | 扩展：技能提炼、知识生成、对话引导的 Prompt |
| `repository/` | 持久化 | 扩展：技能库、知识库、经验库的 CRUD |
| **`scheduler/`** | **新增** | 定时巡检、推送调度、经验衰减计算 |
| **`learning/`** | **新增** | 对话反思、技能提炼、知识抽取 |
| **`push/`** | **新增** | 推送内容生成、推送规则管理 |

### 8.3 LLM 调用策略

- **推理模型**：使用外部 API（DeepSeek / OpenAI 兼容接口），不依赖本地模型
- **API 可配置**：通过 `config.yaml` 配置 provider、model、endpoint
- **调用场景**：
  - 意图识别（对话输入 → 结构化意图）
  - 技能匹配（根据意图检索和排序技能）
  - 结果组织（结构化数据 → 自然语言输出）
  - 技能提炼（对话记录 → 结构化技能步骤）
  - 知识抽取（分析过程 → 结构化知识条目）
  - 推送摘要生成
- **降级策略**：LLM 不可用时，退回到纯 Analyzer 模式（规则引擎 + 固定模板输出），不影响基本分析能力

## 9. Tool 层设计

### 9.1 数据查询工具

基于现有 Tool 层扩展，所有工具返回**清洗过的结构化数据**，不是原始 SQL 结果。

```go
// 设备维度
get_equipment_profile(equipment_id)               // 设备基本信息+当前状态
get_equipment_repair_history(equipment_id, range) // 维修记录（含备件消耗）
get_equipment_maintenance_records(equipment_id, range) // 保养执行记录
get_equipment_inspection_records(equipment_id, range)  // 点检记录（含NG项）
get_equipment_runtime(equipment_id, range)        // 运行时长/停机时长

// 车间/工厂维度
get_failure_distribution(scope, range, group_by)  // 故障分布
get_maintenance_compliance(scope, range)          // 保养完成情况
get_inspection_ng_summary(scope, range)           // 点检NG汇总
get_spare_part_consumption(scope, range)          // 备件消耗统计

// 对比维度
get_similar_failure_cases(equipment_type, symptom) // 同类型同症状历史案例
get_equipment_benchmark(equipment_id)              // 同型号设备对标

// Analyzer 输出
run_audit_checks(scope, range)                     // 调用维修审核 Analyzer
run_maintenance_recommendation(equipment_type_id)  // 调用保养建议 Analyzer
```

### 9.2 知识与技能工具

```go
// 知识检索
search_knowledge(query, scope)                     // 从知识库检索相关知识
get_knowledge_detail(knowledge_id)                 // 获取知识详情

// 技能工具
match_skills(intent, context)                      // 匹配适用于当前任务的技能
execute_skill(skill_id, context)                   // 执行技能步骤
get_skill_detail(skill_id)                         // 获取技能详情

// 经验工具
get_user_preferences(user_id, domain)              // 获取用户偏好
apply_experience(experience_id, context)           // 应用经验校准
```

### 9.3 Tool 安全边界

- 所有工具调用受 Policy 层约束，注入用户权限上下文
- 工具不接收 LLM 生成的任意参数，参数经过白名单校验
- 跨工厂/车间访问在工具层统一拦截
- 每次工具调用记录到 Agent Session 日志

## 10. 学习引擎

### 10.1 对话后反思流程

每次对话结束后，Learning 模块执行反思流程：

```text
对话结束
  │
  ├─ 1. 判断是否值得提炼
  │     - 对话是否产生了有价值的结论？
  │     - 工程师是否给出了引导性修正？
  │     - 分析路径是否有可复用性？
  │     → 如果否，结束
  │
  ├─ 2. 提炼知识（如果结论有价值）
  │     - 提取结论和证据
  │     - 生成知识条目（status: draft）
  │     - 关联相关设备/车间/时间段
  │
  ├─ 3. 提炼技能（如果分析路径有复用性）
  │     - 对比初始路径 vs 最终路径
  │     - 识别关键修正点
  │     - 抽象为通用步骤
  │     - 生成技能草稿（status: draft）
  │
  └─ 4. 收集经验（如果工程师有反馈）
        - 记录偏好、修正、边界
        - 更新经验库
```

### 10.2 技能自动优化

Learning 模块定期（每月）执行技能优化：

- 检查 `success_rate` 低于 0.5 的技能，分析失败原因
- 检查功能相似的技能，建议合并
- 基于近期使用数据，建议步骤调整
- 生成优化建议，由工程师确认后执行

### 10.3 知识自动维护

- 每月检查未确认的 draft 知识，提示工程师审核
- 超过 90 天未确认的 draft 知识自动归档
- confirmed 知识引用时验证是否仍然适用（对比最新数据）

## 11. 推送引擎

### 11.1 推送原则

- **有才推**：没有实质发现不推送，避免信息疲劳
- **精准推**：推送给相关的人，不群发
- **可配置**：工程师可以配置关注范围和推送频率
- **不打扰**：避免工作时间外推送，支持免打扰时段

### 11.2 推送类型

| 类型 | 频率 | 触发条件 | 目标用户 |
|------|------|---------|---------|
| 日报摘要 | 每日 | 有异常发现时推送 | 设备工程师 |
| 设备预警 | 实时 | 报修单创建时关联设备上下文 | 相关工程师 |
| 周期报告 | 月度/季度 | 自动生成 | 管理者 |
| 知识沉淀通知 | 按需 | 新知识/技能需要确认时 | 相关工程师 |

### 11.3 推送内容格式

```text
📋 设备健康日报 — 三车间 — 2026年4月25日

⚠️ 需要关注（2项）：
1. B线4号CNC（EQ-003-024）：近7天点检2次报"液压温度偏高"NG，
   该设备上次保养已逾期5天。建议优先安排保养。
2. 备件预警：NSK 6205轴承当前库存3件，安全库存5件，
   按近30天消耗速率预计10天后缺货。

📊 数据快照：
- 报修单：昨日3单（日均水平）
- 保养逾期：5项（较上周-2）
- 点检NG率：3.2%（较上周+0.4%）

🔗 详细分析请点击查看
```

## 12. 前端设计建议

### 12.1 Agent 入口

**主入口：Agent 对话页**

位置：导航栏新增"AI 助手"入口

布局：
- 左侧：会话列表 + 快捷分析模板
- 中间：对话区域（消息式交互）
- 右侧：关联信息面板（设备信息、知识卡片、技能引用）

**辅助入口：嵌入式面板**

- 报修单详情页：右侧嵌入"设备诊断摘要"面板
- 设备详情页：新增"AI 洞察"标签页
- 首页仪表盘：新增"Agent 发现"卡片

### 12.2 技能管理页

仅对工程师及以上角色可见：

- 技能列表：名称、适用场景、成功率、使用次数、状态
- 技能详情：完整步骤、执行历史、版本记录
- 技能编辑：修改步骤、调整参数、更新适用条件

### 12.3 知识管理页

- 知识搜索：全文检索 + 设备/车间/时间筛选
- 知识详情：完整内容、证据链、引用关系
- 知识审核：draft 知识的确认/修正/拒绝操作

## 13. 数据库 Schema 设计（新增表）

在 Phase 1 已有表的基础上新增：

### 13.1 `agent_skills`

```sql
CREATE TABLE agent_skills (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    applicable_to JSONB DEFAULT '[]',
    applicable_scenarios TEXT[] DEFAULT '{}',
    steps JSONB NOT NULL,
    scope JSONB DEFAULT '{}',
    version INTEGER DEFAULT 1,
    status VARCHAR(20) DEFAULT 'draft',
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    usage_count INTEGER DEFAULT 0,
    success_rate NUMERIC(5,4) DEFAULT 0,
    last_used TIMESTAMP,
    changelog JSONB DEFAULT '[]'
);

CREATE INDEX idx_agent_skills_status ON agent_skills(status);
CREATE INDEX idx_agent_skills_applicable ON agent_skills USING GIN(applicable_to);
```

### 13.2 `agent_knowledge`

```sql
CREATE TABLE agent_knowledge (
    id VARCHAR(50) PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    type VARCHAR(50) NOT NULL,
    scope JSONB DEFAULT '{}',
    summary TEXT NOT NULL,
    details JSONB,
    related_skill_id VARCHAR(100),
    confidence NUMERIC(5,4) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'draft',
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    verified_at TIMESTAMP,
    referenced_count INTEGER DEFAULT 0,
    last_referenced TIMESTAMP,
    expire_at TIMESTAMP,
    search_vector tsvector
);

CREATE INDEX idx_agent_knowledge_status ON agent_knowledge(status);
CREATE INDEX idx_agent_knowledge_type ON agent_knowledge(type);
CREATE INDEX idx_agent_knowledge_scope ON agent_knowledge USING GIN(scope);
CREATE INDEX idx_agent_knowledge_search ON agent_knowledge USING GIN(search_vector);
```

### 13.3 `agent_experiences`

```sql
CREATE TABLE agent_experiences (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    category VARCHAR(100),
    content JSONB NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    weight NUMERIC(5,4) DEFAULT 0.8,
    decay_rate NUMERIC(5,4) DEFAULT 0.05,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_applied TIMESTAMP,
    expire_at TIMESTAMP
);

CREATE INDEX idx_agent_experiences_user ON agent_experiences(user_id, type);
CREATE INDEX idx_agent_experiences_status ON agent_experiences(status);
```

### 13.4 `agent_conversations`

```sql
CREATE TABLE agent_conversations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE agent_messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL REFERENCES agent_conversations(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    tool_calls JSONB,
    skill_id VARCHAR(100),
    knowledge_ids JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_agent_messages_conversation ON agent_messages(conversation_id, created_at);
```

### 13.5 `agent_push_subscriptions`

```sql
CREATE TABLE agent_push_subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    push_type VARCHAR(50) NOT NULL,
    scope JSONB DEFAULT '{}',
    frequency VARCHAR(20) DEFAULT 'daily',
    quiet_hours_start TIME,
    quiet_hours_end TIME,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, push_type)
);
```

## 14. 分阶段交付

### Phase 2.1：对话 + 知识沉淀

**目标**：实现基础的对话分析能力，分析结论自动沉淀为知识条目。

- 对话接口（多轮对话，上下文管理）
- 知识库 CRUD + 检索
- 分析后自动生成 draft 知识
- 工程师确认/修正知识
- 推送：事件触发型（报修单关联设备摘要）

### Phase 2.2：技能系统

**目标**：实现技能的自动提炼和复用。

- 技能库 CRUD + 匹配执行
- 对话反思引擎（自动提炼技能草稿）
- 技能详情页和编辑功能
- 技能自动匹配和排序

### Phase 2.3：经验系统 + 主动推送

**目标**：完善人机闭环，实现主动发现和推送。

- 经验库 + 衰减机制
- 定时巡检引擎
- 推送订阅管理
- 月度/季度报告生成
- 技能自动优化

### Phase 2.4：高级检索

**目标**：引入语义检索，提升知识/技能匹配精度。

- pgvector 集成
- 知识和手册的向量索引
- 语义+结构化混合检索

## 15. 成功指标

| 指标 | 目标 | 衡量方式 |
|------|------|---------|
| **知识采纳率** | confirmed / (confirmed + rejected) > 70% | 知识库统计 |
| **技能成功率** | 平均 success_rate > 0.7 | 技能库统计 |
| **推送有效性** | 推送点击率 > 40% | 推送统计 |
| **工程师参与度** | 月活用户 > 60% | 用户统计 |
| **分析采纳率** | 工程师标记"有用"的比例 > 65% | 反馈统计 |
| **冷启动时间** | 部署3个月内技能库 ≥ 10条有效技能 | 技能库统计 |
