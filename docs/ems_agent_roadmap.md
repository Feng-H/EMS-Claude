# EMS Agent Roadmap

## 1. Purpose

This roadmap connects the EMS Agent PRD, phase-one design, and phase-two extension documents into one delivery sequence.

It supports:

- [ems_agent_prd.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_prd.md)
- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)
- [ems_agent_schema_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_schema_design.md)
- [ems_agent_api_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_api_design.md)
- [ems_agent_phase1_task_breakdown.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase1_task_breakdown.md)
- [ems_agent_phase2_schema_extension.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase2_schema_extension.md)
- [ems_agent_phase2_api_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase2_api_design.md)

Its purpose is to answer:

- what we build first
- what depends on what
- what belongs to phase one versus later phases
- how the product evolves from a management assistant into a learning agent

## 2. Evolution Path

The EMS Agent should evolve in four steps:

1. `Phase 1`
   Recommendation and audit management assistant
2. `Phase 2.1`
   Multi-turn conversation and knowledge accumulation
3. `Phase 2.2`
   Skill system and reusable analysis methods
4. `Phase 2.3`
   Experience calibration and proactive push
5. `Phase 2.4`
   Advanced retrieval and semantic matching

This means the product does not start as a fully autonomous learning agent.

It starts as a controlled, evidence-based assistant and grows into one.

## 3. Phase Summary

| Phase | Core Capability | Product Shape | Main User Value |
|------|------------------|---------------|-----------------|
| `Phase 1` | 建议 + 审计 + 分析 | 管理助手 | 快速得到有依据的结论和建议 |
| `Phase 2.1` | 对话 + 知识沉淀 | 会话助手 | 从一次分析中沉淀可复用知识 |
| `Phase 2.2` | 技能系统 | 学习型助手 | 从成功路径中形成分析方法 |
| `Phase 2.3` | 经验系统 + 主动推送 | 主动协作助手 | 更懂工程师偏好，能主动提醒 |
| `Phase 2.4` | 高级检索 | 高精度记忆助手 | 更准确地召回知识、技能和手册 |

## 4. Phase 1

## 4.1 Goal

Build a reliable, bounded management assistant for engineers and supervisors.

## 4.2 What is included

- maintenance optimization recommendation
- repair and maintenance audit
- structured analysis support
- evidence-based outputs
- knowledge-base and manual retrieval
- Chinese-first output
- strict permission scope
- audit persistence for sessions, artifacts, and evidence

## 4.3 What is not included

- multi-turn general chat
- skill store
- knowledge wiki generated from conversations
- experience store
- proactive push subscriptions
- semantic vector retrieval as a hard dependency

## 4.4 Main supporting docs

- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)
- [ems_agent_schema_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_schema_design.md)
- [ems_agent_api_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_api_design.md)
- [ems_agent_phase1_task_breakdown.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase1_task_breakdown.md)

## 5. Phase 2.1

## 5.1 Goal

Turn the management assistant into a conversation-capable assistant that can accumulate draft knowledge from useful analyses.

## 5.2 What is included

- conversation creation and history
- chat interface
- structured intent recognition from conversation
- analysis through conversation
- automatic draft knowledge generation after valuable analyses
- knowledge review and confirmation workflow
- event-triggered device context summary

## 5.3 Dependencies

Depends on phase one:

- tool layer
- analyzer outputs
- permission enforcement
- evidence persistence

## 5.4 Main supporting docs

- [ems_agent_prd.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_prd.md)
- [ems_agent_phase2_schema_extension.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase2_schema_extension.md)
- [ems_agent_phase2_api_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase2_api_design.md)

## 6. Phase 2.2

## 6.1 Goal

Extract reusable analysis methods and turn them into skills.

## 6.2 What is included

- skill store
- skill draft generation
- skill confirmation and editing
- skill matching
- skill-driven tool execution

## 6.3 Dependencies

Depends on phase 2.1:

- conversation history
- knowledge accumulation
- post-analysis reflection flow

Depends on phase one:

- stable tools
- analyzers
- scope-safe execution path

## 7. Phase 2.3

## 7.1 Goal

Close the human-in-the-loop loop with experience calibration and proactive push.

## 7.2 What is included

- experience store
- preference and correction capture
- experience decay
- daily or event-driven push
- push subscriptions
- monthly and quarterly reports
- knowledge and skill review notifications

## 7.3 Dependencies

Depends on phase 2.1 and 2.2:

- conversation
- knowledge
- skills
- feedback data

Requires additional scheduling and push infrastructure.

## 8. Phase 2.4

## 8.1 Goal

Improve retrieval accuracy for manuals, knowledge, and skills through semantic matching.

## 8.2 What is included

- `pgvector` integration or equivalent vector retrieval layer
- embeddings for manual chunks
- embeddings for agent knowledge
- embeddings for skills where useful
- hybrid retrieval combining scope filters, keyword search, and semantic ranking

## 8.3 Dependencies

Depends on:

- enough data volume in manuals and knowledge
- stable retrieval baselines from earlier phases

This phase should not block phase one or early phase two delivery.

## 9. Dependency Map

```text
Phase 1
  ├─ tools
  ├─ analyzers
  ├─ policy and scope
  ├─ session and artifact persistence
  └─ structured recommendation and audit output
       ↓
Phase 2.1
  ├─ conversation
  ├─ knowledge draft generation
  └─ knowledge review
       ↓
Phase 2.2
  ├─ reflection
  ├─ skill extraction
  └─ skill execution
       ↓
Phase 2.3
  ├─ experience calibration
  ├─ scheduler
  └─ proactive push
       ↓
Phase 2.4
  └─ semantic retrieval upgrade
```

## 10. Recommended Delivery Order

Recommended implementation order:

1. finish phase one fully
2. add phase 2.1 conversation and knowledge draft flow
3. add phase 2.2 skills
4. add phase 2.3 experience and push
5. upgrade retrieval in phase 2.4 when scale justifies it

This order keeps the product grounded in real data and user feedback before adding learning and autonomy.

## 11. Risks by Phase

### Phase 1 risks

- recommendation quality depends on data completeness
- permission leakage risk if tool scoping is weak
- weak retrieval quality if knowledge and manual search are poor

### Phase 2.1 risks

- conversations may become generic if not tied tightly to tools and scope
- draft knowledge volume may grow faster than review capacity

### Phase 2.2 risks

- poorly abstracted skills may become brittle
- too many low-quality skills can reduce matching quality

### Phase 2.3 risks

- push fatigue if the signal threshold is too low
- stale preferences can distort outputs if decay is weak

### Phase 2.4 risks

- vector retrieval can add complexity before there is enough value
- semantic recall without scope control can hurt precision

## 12. Exit Criteria by Phase

### Phase 1 exit criteria

- engineers can use recommendation and audit flows in production-like conditions
- outputs are evidence-linked and scope-safe
- core artifact persistence is complete

### Phase 2.1 exit criteria

- multi-turn conversation works reliably
- useful analyses can generate draft knowledge
- engineers can confirm or reject knowledge

### Phase 2.2 exit criteria

- at least a small set of active skills is reusable across sessions
- skill matching improves analysis efficiency or consistency

### Phase 2.3 exit criteria

- user preferences measurably improve output usefulness
- push reports deliver clear value without causing fatigue

### Phase 2.4 exit criteria

- semantic retrieval outperforms keyword-only retrieval for real knowledge recall cases

## 13. Recommended Communication Rule

When discussing implementation with engineering or stakeholders:

- use `Phase 1` to mean bounded management assistant only
- use `Phase 2.x` to mean learning and memory capabilities

This avoids mixing today's deliverables with the long-term product vision.

## 14. Reference

This roadmap supports:

- [ems_agent_prd.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_prd.md)
- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)
- [ems_agent_phase1_task_breakdown.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase1_task_breakdown.md)
- [ems_agent_phase2_schema_extension.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase2_schema_extension.md)
- [ems_agent_phase2_api_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase2_api_design.md)
