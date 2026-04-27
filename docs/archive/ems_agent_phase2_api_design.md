# EMS Agent Phase 2 API Design

## 1. Purpose

This document extends the phase-one API design for the learning and proactive EMS Agent described in the PRD.

It supports:

- [ems_agent_prd.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_prd.md)
- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)
- [ems_agent_api_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_api_design.md)

This document covers:

- multi-turn chat
- skill management
- knowledge management
- experience management
- push subscription management

## 2. API Principles

Phase-two APIs should preserve all phase-one constraints:

- permission-scoped access
- Chinese-first output
- traceability
- evidence-aware responses

Additional phase-two principles:

- support human confirmation and editing loops
- separate draft content from confirmed content
- keep learning actions auditable

## 3. Chat APIs

## 3.1 `POST /api/v1/agent/chat`

Purpose:

- submit a user message into a conversation
- run intent recognition, skill matching, tool orchestration, and response generation

Suggested request:

```json
{
  "conversation_id": 12,
  "message": "三车间最近报修特别多，帮我看看什么情况",
  "context": {
    "factory_id": 3,
    "workshop_id": 8
  },
  "language": "zh-CN"
}
```

Suggested response:

```json
{
  "success": true,
  "trace_id": "agt_chat_001",
  "language": "zh-CN",
  "conversation_id": 12,
  "message_id": 88,
  "summary": "三车间近30天报修量较上月上升62%，主要集中在B线CNC主轴相关问题。",
  "data": {
    "reply": "我已经分析了三车间近30天的数据...",
    "used_skill": {
      "id": "repeat_failure_investigation",
      "name": "重复故障排查"
    },
    "referenced_knowledge": [],
    "evidence": []
  }
}
```

## 3.2 `POST /api/v1/agent/conversations`

Purpose:

- create a new conversation

Suggested request:

```json
{
  "title": "三车间设备异常分析",
  "language": "zh-CN"
}
```

## 3.3 `GET /api/v1/agent/conversations`

Purpose:

- list current user's conversations

## 3.4 `GET /api/v1/agent/conversations/:id`

Purpose:

- retrieve conversation metadata and messages

Suggested response fields:

- `conversation`
- `messages`

## 4. Skill APIs

## 4.1 `GET /api/v1/agent/skills`

Purpose:

- list skills visible to the current user

Suggested filters:

- `status`
- `domain`
- `keyword`

## 4.2 `GET /api/v1/agent/skills/:id`

Purpose:

- retrieve one skill definition and usage summary

## 4.3 `POST /api/v1/agent/skills`

Purpose:

- create a skill draft manually

Suggested request:

```json
{
  "name": "备件批次质量排查",
  "description": "当故障率异常上升时，排查是否与备件批次相关。",
  "applicable_to": ["fault_rate_anomaly"],
  "applicable_scenarios": ["同类型设备故障率异常上升"],
  "steps": [],
  "scope": {
    "domain": "spare_part"
  }
}
```

## 4.4 `PUT /api/v1/agent/skills/:id`

Purpose:

- edit a skill draft or active skill

## 4.5 `POST /api/v1/agent/skills/:id/activate`

Purpose:

- confirm a draft skill and move it to active state

## 4.6 `POST /api/v1/agent/skills/:id/deprecate`

Purpose:

- mark a skill as deprecated

## 5. Knowledge APIs

## 5.1 `GET /api/v1/agent/knowledge`

Purpose:

- search and list knowledge items

Suggested filters:

- `status`
- `type`
- `factory_id`
- `workshop_id`
- `equipment_type`
- `keyword`

## 5.2 `GET /api/v1/agent/knowledge/:id`

Purpose:

- retrieve one knowledge item with details and references

## 5.3 `POST /api/v1/agent/knowledge/:id/confirm`

Purpose:

- confirm a draft knowledge item

Suggested request:

```json
{
  "comment": "结论成立，可以作为后续同类问题参考。"
}
```

## 5.4 `POST /api/v1/agent/knowledge/:id/reject`

Purpose:

- reject a draft knowledge item

## 5.5 `PUT /api/v1/agent/knowledge/:id`

Purpose:

- edit a knowledge item before or during confirmation workflow

## 6. Experience APIs

## 6.1 `GET /api/v1/agent/experiences`

Purpose:

- list experience items for the current user or authorized scope

Suggested filters:

- `type`
- `status`
- `domain`

## 6.2 `POST /api/v1/agent/experiences`

Purpose:

- create or record an experience item explicitly

Suggested request:

```json
{
  "type": "preference",
  "category": "analysis_depth",
  "content": {
    "preference": "备件相关分析默认展示同比和环比"
  }
}
```

## 6.3 `PUT /api/v1/agent/experiences/:id`

Purpose:

- update or deactivate an experience item

## 7. Push APIs

## 7.1 `GET /api/v1/agent/push/subscriptions`

Purpose:

- list current user's push subscriptions

## 7.2 `POST /api/v1/agent/push/subscriptions`

Purpose:

- create or upsert a push subscription

Suggested request:

```json
{
  "push_type": "daily_summary",
  "scope": {
    "factory_id": 3,
    "workshop_ids": [8, 9]
  },
  "frequency": "daily",
  "quiet_hours_start": "21:00:00",
  "quiet_hours_end": "08:00:00",
  "enabled": true
}
```

## 7.3 `PUT /api/v1/agent/push/subscriptions/:id`

Purpose:

- update subscription scope, frequency, or quiet hours

## 7.4 `DELETE /api/v1/agent/push/subscriptions/:id`

Purpose:

- disable or remove a push subscription

## 8. Feedback APIs

## 8.1 `POST /api/v1/agent/feedback`

Purpose:

- collect lightweight feedback from agent outputs

Suggested request:

```json
{
  "artifact_id": 1001,
  "feedback_type": "useful",
  "reason": null,
  "comment": null
}
```

Supported `feedback_type` values:

- `useful`
- `not_useful`
- `supplement`

## 9. Common Error Codes

Suggested additions beyond phase one:

- `CONVERSATION_NOT_FOUND`
- `SKILL_NOT_FOUND`
- `KNOWLEDGE_NOT_FOUND`
- `EXPERIENCE_NOT_FOUND`
- `PUSH_SUBSCRIPTION_NOT_FOUND`
- `INVALID_SKILL_STATE`
- `INVALID_KNOWLEDGE_STATE`

## 10. Permission Notes

Recommended access boundary:

- engineers and above can view and confirm skills and knowledge within authorized scope
- ordinary users should not manage shared skills by default
- experiences are user-scoped unless explicitly elevated for admin review
- push subscriptions are user-owned

## 11. API Evolution Note

Recommended split:

- keep phase-one analysis endpoints in [ems_agent_api_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_api_design.md)
- keep conversation and memory-management endpoints in this phase-two extension document

## 12. Reference

This document supports:

- [ems_agent_prd.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_prd.md)
- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)
- [ems_agent_api_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_api_design.md)
