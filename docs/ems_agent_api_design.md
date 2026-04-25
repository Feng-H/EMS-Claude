# EMS Agent API Design

## 1. Purpose

This document defines the phase-one API design for the EMS management assistant agent.

It is a supporting document for the main design:

- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)

This API design focuses on:

- recommendation and audit request patterns
- response structure and evidence traceability
- localization metadata
- permission-aware request handling

## 2. API Principles

### 2.1 Read-only in phase one

Phase-one APIs should not directly mutate core business records.

They should return:

- recommendations
- anomaly findings
- analysis results
- draft-style artifacts

### 2.2 Structured request first

The UI should prefer structured filters and scenario-specific parameters instead of sending only free-form prompts.

### 2.3 Structured response plus readable summary

Each response should contain:

- a readable summary
- structured metrics
- evidence links
- risk indicators
- traceable metadata

### 2.4 Scope must come from authentication context

The backend should not trust client-supplied organization scope blindly.

The authenticated user context should determine:

- allowed factory scope
- allowed workshop scope
- accessible equipment and work order range

## 3. Common Response Envelope

Suggested phase-one response envelope:

```json
{
  "success": true,
  "trace_id": "agt_20260425_xxxx",
  "language": "zh-CN",
  "scenario": "maintenance_recommendation",
  "scope_summary": {
    "factory_id": 2,
    "workshop_id": 8
  },
  "summary": "建议将二级保养周期从30天调整为21天，并补充润滑回路检查项。",
  "risk_level": "medium",
  "artifact_id": 1001,
  "evidence_count": 4,
  "data": {}
}
```

Recommended common fields:

- `success`
- `trace_id`
- `language`
- `scenario`
- `scope_summary`
- `summary`
- `risk_level`
- `artifact_id`
- `evidence_count`
- `data`

## 4. Common Supporting Types

### 4.1 Evidence object

```json
{
  "evidence_type": "knowledge",
  "source_table": "knowledge_articles",
  "source_id": 21,
  "title": "数控主轴振动故障处理最佳实践",
  "excerpt": "建议优先检查润滑回路和轴承温升。",
  "score": 0.93
}
```

### 4.2 Recommendation item

```json
{
  "type": "cycle_adjustment",
  "target": "maintenance_plan",
  "target_id": 18,
  "title": "缩短二级保养周期",
  "description": "建议从30天调整为21天。",
  "reason": "近90天重复故障增加，且同类设备平均周期为21天。",
  "impact": "预计降低重复故障率和平均停机时长。"
}
```

### 4.3 Anomaly item

```json
{
  "anomaly_type": "repeat_failure_after_maintenance",
  "severity": "high",
  "target_type": "repair_order",
  "target_id": 301,
  "title": "保养后短期重复故障",
  "description": "同一设备在保养后3天内再次出现相同故障。",
  "suggested_action": "复核保养项是否覆盖润滑回路检查，并核对维修方案是否参考知识库最佳实践。"
}
```

## 5. Endpoints

## 5.1 `POST /api/v1/agent/maintenance/recommend`

Purpose:

- generate maintenance optimization recommendations for equipment, equipment type, or similar equipment group

Suggested request:

```json
{
  "factory_id": 2,
  "workshop_id": 8,
  "equipment_type_id": 5,
  "equipment_ids": [101, 102, 103],
  "time_range": {
    "start_date": "2026-01-01",
    "end_date": "2026-03-31"
  },
  "question": "请评估这批设备的保养项目和周期是否需要调整。",
  "language": "zh-CN"
}
```

Request notes:

- `factory_id` and `workshop_id` are advisory filters and still require backend permission validation
- `equipment_ids` should be optional
- if neither `equipment_ids` nor `equipment_type_id` is provided, the API should reject the request

Suggested response:

```json
{
  "success": true,
  "trace_id": "agt_x1",
  "language": "zh-CN",
  "scenario": "maintenance_recommendation",
  "scope_summary": {
    "factory_id": 2,
    "workshop_id": 8
  },
  "summary": "建议缩短二级保养周期，并补充润滑回路检查和轴承温升记录。",
  "risk_level": "medium",
  "artifact_id": 1001,
  "evidence_count": 5,
  "data": {
    "current_plan": {
      "plan_id": 18,
      "plan_name": "CNC 二级保养",
      "cycle_days": 30,
      "item_count": 8,
      "completion_rate": 81.2
    },
    "recommendations": [],
    "expected_benefits": [
      "降低重复故障率",
      "降低平均停机时长"
    ],
    "evidence": []
  }
}
```

`data.recommendations` should be an array of recommendation items.

## 5.2 `POST /api/v1/agent/audit/repair`

Purpose:

- audit repair reasonableness, repeat failures, and cost anomalies

Suggested request:

```json
{
  "factory_id": 2,
  "workshop_id": 8,
  "equipment_type_id": 5,
  "time_range": {
    "start_date": "2026-03-01",
    "end_date": "2026-03-31"
  },
  "anomaly_types": [
    "repeat_failure_after_maintenance",
    "cost_deviation",
    "best_practice_mismatch"
  ],
  "language": "zh-CN"
}
```

Suggested response:

```json
{
  "success": true,
  "trace_id": "agt_x2",
  "language": "zh-CN",
  "scenario": "repair_audit",
  "scope_summary": {
    "factory_id": 2,
    "workshop_id": 8
  },
  "summary": "发现3条高风险维修异常，其中2条涉及保养后短期重复故障，1条涉及费用偏高。",
  "risk_level": "high",
  "artifact_id": 1002,
  "evidence_count": 8,
  "data": {
    "anomalies": [],
    "stats": {
      "checked_orders": 56,
      "high_risk_count": 3,
      "medium_risk_count": 5
    },
    "evidence": []
  }
}
```

`data.anomalies` should be an array of anomaly items.

## 5.3 `POST /api/v1/agent/audit/maintenance`

Purpose:

- audit maintenance plan density, execution quality, and mismatch with failure patterns

Suggested request:

```json
{
  "factory_id": 2,
  "equipment_type_id": 5,
  "time_range": {
    "start_date": "2026-01-01",
    "end_date": "2026-03-31"
  },
  "focus": [
    "over_maintenance",
    "under_maintenance",
    "task_execution_quality"
  ],
  "language": "zh-CN"
}
```

Suggested response fields in `data`:

- `audit_summary`
- `anomalies`
- `plan_comparisons`
- `evidence`

## 5.4 `POST /api/v1/agent/analyze`

Purpose:

- answer structured management questions with controlled tools

Suggested request:

```json
{
  "factory_id": 2,
  "workshop_id": 8,
  "question": "为什么这个月维修费用明显上升？",
  "time_range": {
    "start_date": "2026-03-01",
    "end_date": "2026-03-31"
  },
  "language": "zh-CN"
}
```

Suggested response fields in `data`:

- `key_findings`
- `metric_comparisons`
- `top_entities`
- `evidence`
- `recommended_actions`

## 5.5 `GET /api/v1/agent/sessions/:id`

Purpose:

- retrieve the persisted session metadata for audit and UI history

Suggested response fields in `data`:

- `session`
- `artifacts`

## 5.6 `GET /api/v1/agent/artifacts/:id`

Purpose:

- retrieve one recommendation, audit report, or analysis report

Suggested response fields in `data`:

- `artifact`
- `evidence`
- `related_session`

## 6. Error Handling

Suggested error envelope:

```json
{
  "success": false,
  "trace_id": "agt_x3",
  "error": {
    "code": "FORBIDDEN_SCOPE",
    "message": "当前用户无权访问该工厂范围的数据。"
  }
}
```

Suggested error codes:

- `INVALID_ARGUMENT`
- `FORBIDDEN_SCOPE`
- `NOT_FOUND`
- `UNSUPPORTED_STORAGE_MODE`
- `ANALYSIS_NOT_AVAILABLE`
- `INTERNAL_ERROR`

## 7. Localization Rules

The phase-one agent should return Chinese by default.

API handling recommendations:

- if `language` is omitted, default to `zh-CN`
- if the user profile later stores preference, the backend may override request defaults
- all summaries, recommendations, and anomaly descriptions should be generated in the resolved language

Recommended response metadata:

- always include `language`

## 8. Knowledge Base Integration

The API layer should reflect that the agent can use both manuals and knowledge articles.

Recommendations:

- expose `knowledge` as an evidence type in all relevant responses
- in audit responses, include best-practice mismatch findings where applicable
- allow future drill-down from evidence to the original knowledge article

## 9. Memory Mode Strategy

Two acceptable phase-one approaches:

### 9.1 Database-only support

Return:

- `UNSUPPORTED_STORAGE_MODE`

when the agent is called under memory mode.

### 9.2 Mock-compatible support

If memory mode support is implemented:

- return clearly simulated results
- use seeded mock anomaly scenarios
- still preserve the same response shape as database mode

## 10. Service-side Enforcement

The service layer should enforce:

- authenticated user context resolution
- scope injection into all tools
- validation of requested equipment and work order identifiers
- persistence of `agent_sessions` and `agent_artifacts`
- evidence count and trace ID generation

The LLM should not be allowed to invent identifiers or widen scope.

## 11. Versioning Recommendation

Suggested versioning path:

- keep endpoints under `/api/v1/agent`
- if response schema changes substantially later, add a result schema version field inside `data`

Suggested optional field:

```json
{
  "data_schema_version": "1.0"
}
```

## 12. Reference

This document supports:

- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)
