# EMS Agent Phase 2 Schema Extension

## 1. Purpose

This document extends the phase-one schema design for the broader EMS Agent PRD direction.

It supports:

- [ems_agent_prd.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_prd.md)
- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)
- [ems_agent_schema_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_schema_design.md)

This document focuses on phase-two persistence for:

- skill store
- knowledge wiki
- experience store
- multi-turn conversations
- push subscriptions

## 2. Positioning

Phase-one schema covers:

- manual retrieval
- cost audit support
- agent session persistence
- evidence linkage

Phase-two extends this into a learning agent with persistent memory and proactive behavior.

## 3. New Tables

## 3.1 `agent_skills`

Purpose:

- store reusable analysis methods
- support draft, active, and deprecated lifecycle
- track usage and success metrics

Suggested SQL:

```sql
CREATE TABLE agent_skills (
    id VARCHAR(100) PRIMARY KEY,
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

CREATE INDEX idx_agent_skills_status
    ON agent_skills(status);

CREATE INDEX idx_agent_skills_applicable
    ON agent_skills USING GIN(applicable_to);
```

Recommended status values:

- `draft`
- `active`
- `deprecated`

Notes:

- `id` should be stable and readable to support reuse and version evolution
- `steps` should store tool-oriented execution templates rather than free-form prose only

## 3.2 `agent_knowledge`

Purpose:

- store knowledge items extracted from successful analyses
- support verification and aging

Suggested SQL:

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

CREATE INDEX idx_agent_knowledge_status
    ON agent_knowledge(status);

CREATE INDEX idx_agent_knowledge_type
    ON agent_knowledge(type);

CREATE INDEX idx_agent_knowledge_scope
    ON agent_knowledge USING GIN(scope);

CREATE INDEX idx_agent_knowledge_search
    ON agent_knowledge USING GIN(search_vector);
```

Recommended status values:

- `draft`
- `confirmed`
- `rejected`
- `archived`

Recommended type values:

- `root_cause_analysis`
- `pattern`
- `equipment_profile`
- `maintenance_insight`
- `spare_part_insight`
- `benchmark_comparison`

Notes:

- `scope` allows factory, workshop, equipment type, equipment IDs, and time period tagging
- `search_vector` supports full-text search in the early phase

## 3.3 `agent_experiences`

Purpose:

- store user-specific and behavior-specific calibration signals

Suggested SQL:

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

CREATE INDEX idx_agent_experiences_user
    ON agent_experiences(user_id, type);

CREATE INDEX idx_agent_experiences_status
    ON agent_experiences(status);
```

Recommended type values:

- `preference`
- `correction`
- `boundary`
- `quality_feedback`

Notes:

- `content` should support domain-level and user-level scoping
- the scheduler can periodically recalculate effective weight based on decay

## 3.4 `agent_conversations`

Purpose:

- persist multi-turn user sessions separately from phase-one analysis sessions

Suggested SQL:

```sql
CREATE TABLE agent_conversations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(500),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_agent_conversations_user
    ON agent_conversations(user_id, updated_at DESC);
```

Recommended status values:

- `active`
- `closed`
- `archived`

## 3.5 `agent_messages`

Purpose:

- persist message history, tool usage, and referenced memory items

Suggested SQL:

```sql
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

CREATE INDEX idx_agent_messages_conversation
    ON agent_messages(conversation_id, created_at);
```

Recommended role values:

- `user`
- `agent`
- `system`

Notes:

- `tool_calls` should store an audit-friendly structured summary, not raw provider response blobs
- `knowledge_ids` can later be normalized if query volume grows

## 3.6 `agent_push_subscriptions`

Purpose:

- persist proactive push configuration

Suggested SQL:

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

Recommended push types:

- `daily_summary`
- `equipment_alert`
- `periodic_report`
- `knowledge_review_notice`

## 4. Relationship with Phase-One Tables

Recommended relationship model:

- `agent_sessions` remains request-level audit persistence
- `agent_artifacts` remains output-level persistence
- `agent_evidence_links` remains evidence-level persistence
- `agent_conversations` and `agent_messages` persist multi-turn chat history
- `agent_knowledge` becomes a reusable memory source for later chat and analysis
- `agent_skills` becomes a reusable execution source for matching and orchestration
- `agent_experiences` calibrate ranking, verbosity, and prioritization

## 5. Suggested Repository Split

Recommended repository ownership:

- `repository/session_repository.go`
- `repository/artifact_repository.go`
- `repository/evidence_repository.go`
- `repository/skill_repository.go`
- `repository/knowledge_repository.go`
- `repository/experience_repository.go`
- `repository/conversation_repository.go`
- `repository/push_repository.go`

## 6. Search and Retrieval Notes

Recommended retrieval order in phase two:

1. scoped structured business data
2. confirmed `agent_knowledge`
3. `knowledge_articles`
4. manual chunks
5. draft knowledge only when explicitly requested by engineers

If `pgvector` is introduced later, phase-two memory items should also support embedding-based recall.

## 7. Migration Advice

Recommended migration order:

1. `agent_skills`
2. `agent_knowledge`
3. `agent_experiences`
4. `agent_conversations`
5. `agent_messages`
6. `agent_push_subscriptions`

## 8. Open Questions

- should `agent_knowledge.related_skill_id` reference `agent_skills(id)` as a formal foreign key
- should `agent_messages.knowledge_ids` stay JSONB or move to a join table later
- should push delivery history get its own table in phase two or phase three

## 9. Reference

This document supports:

- [ems_agent_prd.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_prd.md)
- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)
- [ems_agent_schema_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_schema_design.md)
