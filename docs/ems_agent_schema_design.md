# EMS Agent Schema Design

## 1. Purpose

This document defines the phase-one database design for the EMS management assistant agent.

It is a supporting document for the main design:

- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)

This schema design focuses on:

- manual and knowledge retrieval support
- repair and maintenance cost audit support
- agent session and artifact persistence
- evidence traceability
- compatibility with existing EMS data domains

## 2. Design Principles

### 2.1 Read-heavy and audit-friendly

The phase-one agent is recommendation- and audit-oriented, so the schema should prioritize:

- read performance
- traceability
- evidence linkage
- easy rollback and evolution

### 2.2 Minimal disruption to existing business tables

The current EMS project already includes:

- equipment
- maintenance
- repair
- spare parts
- analytics
- knowledge

The new schema should extend these domains instead of forcing major rewrites of existing tables.

### 2.3 Recommendation and audit data should be persisted separately

Agent-generated results should not overwrite official business data directly.

Instead, they should be stored as:

- sessions
- artifacts
- evidence links
- optional confirmation metadata

## 3. Existing Tables Reused by the Agent

The phase-one agent should reuse the following existing tables:

- `equipment`
- `equipment_types`
- `maintenance_plans`
- `maintenance_tasks`
- `repair_orders`
- `repair_logs`
- `spare_parts`
- `spare_part_consumptions`
- `knowledge_articles`
- `users`
- `bases`
- `factories`
- `workshops`

These are part of the existing EMS system and provide the operational context for recommendation and audit logic.

## 4. New Tables

## 4.1 `equipment_manual_documents`

Purpose:

- store manual metadata
- support versioning and document lifecycle
- link manuals to equipment type or specific equipment

Suggested SQL:

```sql
CREATE TABLE equipment_manual_documents (
    id SERIAL PRIMARY KEY,
    equipment_type_id INTEGER REFERENCES equipment_types(id) ON DELETE CASCADE,
    equipment_id INTEGER REFERENCES equipment(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    version VARCHAR(50),
    source VARCHAR(100),
    uploaded_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    parsing_status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (equipment_type_id IS NOT NULL OR equipment_id IS NOT NULL)
);

CREATE INDEX idx_equipment_manual_documents_type
    ON equipment_manual_documents(equipment_type_id);

CREATE INDEX idx_equipment_manual_documents_equipment
    ON equipment_manual_documents(equipment_id);
```

Notes:

- `equipment_type_id` is preferred for shared manuals
- `equipment_id` supports device-specific documents
- `parsing_status` can track whether text extraction and chunking are complete

## 4.2 `equipment_manual_chunks`

Purpose:

- store chunked manual content
- support full-text search now
- support semantic retrieval later

Suggested SQL:

```sql
CREATE TABLE equipment_manual_chunks (
    id SERIAL PRIMARY KEY,
    document_id INTEGER NOT NULL REFERENCES equipment_manual_documents(id) ON DELETE CASCADE,
    section_title VARCHAR(255),
    content TEXT NOT NULL,
    chunk_index INTEGER NOT NULL DEFAULT 0,
    embedding_ref VARCHAR(255),
    search_vector tsvector,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_equipment_manual_chunks_document
    ON equipment_manual_chunks(document_id);

CREATE INDEX idx_equipment_manual_chunks_search
    ON equipment_manual_chunks USING GIN(search_vector);
```

Notes:

- `search_vector` supports PostgreSQL full-text search as the phase-one baseline
- `embedding_ref` allows future semantic retrieval without forcing an immediate vector dependency
- if `pgvector` is introduced later, a vector column can be added in a follow-up migration

## 4.3 `repair_cost_details`

Purpose:

- store structured repair cost details for audit analysis

Suggested SQL:

```sql
CREATE TABLE repair_cost_details (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES repair_orders(id) ON DELETE CASCADE,
    labor_cost NUMERIC(12,2) DEFAULT 0,
    spare_part_cost NUMERIC(12,2) DEFAULT 0,
    outsource_cost NUMERIC(12,2) DEFAULT 0,
    other_cost NUMERIC(12,2) DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'CNY',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(order_id)
);
```

Notes:

- phase one can keep this as a single summary record per repair order
- if the business later needs line-item granularity, a child table can be added

## 4.4 `equipment_runtime_snapshots`

Purpose:

- capture runtime and downtime indicators for recommendation logic

Suggested SQL:

```sql
CREATE TABLE equipment_runtime_snapshots (
    id SERIAL PRIMARY KEY,
    equipment_id INTEGER NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    snapshot_date DATE NOT NULL,
    runtime_hours NUMERIC(10,2) DEFAULT 0,
    downtime_hours NUMERIC(10,2) DEFAULT 0,
    load_rate NUMERIC(5,2),
    output_qty NUMERIC(12,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(equipment_id, snapshot_date)
);

CREATE INDEX idx_equipment_runtime_snapshots_equipment_date
    ON equipment_runtime_snapshots(equipment_id, snapshot_date DESC);
```

Notes:

- this table is intentionally simple for phase one
- data can come from manual input, interface sync, or future IoT integration

## 4.5 `agent_sessions`

Purpose:

- track who asked the agent, under what scenario, and under which scope

Suggested SQL:

```sql
CREATE TABLE agent_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    scenario VARCHAR(50) NOT NULL,
    factory_id INTEGER REFERENCES factories(id) ON DELETE SET NULL,
    workshop_id INTEGER REFERENCES workshops(id) ON DELETE SET NULL,
    language VARCHAR(20) DEFAULT 'zh-CN',
    query_text TEXT,
    input_snapshot JSONB,
    status VARCHAR(20) DEFAULT 'completed',
    trace_id VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_agent_sessions_user
    ON agent_sessions(user_id, created_at DESC);

CREATE INDEX idx_agent_sessions_factory
    ON agent_sessions(factory_id, created_at DESC);

CREATE INDEX idx_agent_sessions_scenario
    ON agent_sessions(scenario, created_at DESC);
```

Notes:

- `trace_id` is used to trace one request across service, analyzer, and response logs
- `input_snapshot` keeps the resolved filters used for the analysis

## 4.6 `agent_artifacts`

Purpose:

- persist agent outputs independently from business tables

Suggested SQL:

```sql
CREATE TABLE agent_artifacts (
    id SERIAL PRIMARY KEY,
    session_id INTEGER NOT NULL REFERENCES agent_sessions(id) ON DELETE CASCADE,
    artifact_type VARCHAR(30) NOT NULL,
    title VARCHAR(255),
    summary TEXT,
    input_snapshot JSONB,
    result_json JSONB NOT NULL,
    risk_level VARCHAR(20),
    confirmed_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    confirmed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_agent_artifacts_session
    ON agent_artifacts(session_id);

CREATE INDEX idx_agent_artifacts_type
    ON agent_artifacts(artifact_type, created_at DESC);
```

Recommended `artifact_type` values:

- `recommendation`
- `audit_report`
- `analysis_report`

Notes:

- `result_json` should hold structured outputs rather than raw text blobs
- `summary` is optimized for UI display and export previews

## 4.7 `agent_evidence_links`

Purpose:

- link each result to concrete supporting evidence

Suggested SQL:

```sql
CREATE TABLE agent_evidence_links (
    id SERIAL PRIMARY KEY,
    artifact_id INTEGER NOT NULL REFERENCES agent_artifacts(id) ON DELETE CASCADE,
    evidence_type VARCHAR(30) NOT NULL,
    source_table VARCHAR(100),
    source_id INTEGER,
    excerpt TEXT,
    score NUMERIC(5,4),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_agent_evidence_links_artifact
    ON agent_evidence_links(artifact_id);

CREATE INDEX idx_agent_evidence_links_source
    ON agent_evidence_links(source_table, source_id);
```

Recommended `evidence_type` values:

- `manual`
- `knowledge`
- `maintenance`
- `repair`
- `cost`
- `benchmark`
- `rule`

Notes:

- `source_table + source_id` keeps the design flexible across multiple domains
- `knowledge` is split from `manual` because knowledge articles should be a first-class evidence source

## 5. Recommended Changes to Existing Tables

These are optional but helpful improvements if the project is ready for small extensions.

### 5.1 `knowledge_articles`

Recommended additions:

- `equipment_type_id` nullable
- `fault_code` nullable
- `is_best_practice` boolean default false
- `status` such as `draft`, `published`, `archived`

Why:

- makes knowledge retrieval more relevant
- supports best-practice comparison in audit scenarios

### 5.2 `repair_orders`

Recommended additions:

- `fault_category` nullable
- `downtime_hours` numeric
- `knowledge_article_id` nullable

Why:

- supports structured audit and failure pattern comparison
- allows marking that a repair explicitly referenced a knowledge article

### 5.3 `maintenance_plans`

Recommended additions:

- `source_type` such as `manual`, `knowledge`, `custom`
- `version`
- `effective_from`

Why:

- helps explain why a plan exists
- supports later comparison between original source recommendation and current configured plan

## 6. Retrieval Strategy

## 6.1 Phase-one baseline

Use PostgreSQL full-text search plus keyword matching.

Recommended approach:

- maintain `search_vector` on manual chunks
- maintain a similar text retrieval strategy for `knowledge_articles`
- rank matches using keyword relevance and domain heuristics

## 6.2 Future enhancement

If semantic retrieval becomes necessary at scale, introduce `pgvector`.

Recommended path:

1. keep current schema compatible via `embedding_ref`
2. later add an actual vector column
3. backfill embeddings for manual chunks and knowledge articles
4. blend vector ranking with structured business filters

## 7. Permission and Scope Model

The schema itself is not sufficient for security, but it should support scope-aware queries.

The service layer should always inject:

- `user_id`
- `factory_id`
- role context
- allowed workshop or organization scope

Query guidelines:

- do not query by arbitrary equipment ID from model output
- resolve candidate records only within authorized organization scope
- record scope snapshot in `agent_sessions.input_snapshot`

## 8. Memory Mode Recommendation

Two acceptable phase-one strategies:

### 8.1 Database-mode only

The simplest option is to enable full agent capability only when `storage.mode=database`.

Pros:

- fewer special cases
- easier audit consistency

Cons:

- demos are less convenient before real data is ready

### 8.2 Mock-compatible mode

Provide limited support in memory mode for demos.

If this is chosen, mock data should include:

- repeated failure after maintenance
- abnormal spare part cost pattern
- a repair case that conflicts with best practice

## 9. Migration Plan

Recommended migration order:

1. add manual document and chunk tables
2. add repair cost details
3. add runtime snapshots
4. add agent session, artifact, and evidence tables
5. optionally extend knowledge and repair tables

## 10. Open Questions

- should `knowledge_articles` be extended directly or joined through a mapping table
- should `repair_cost_details` remain summary-level or move to line-item detail early
- should `agent_artifacts.result_json` be versioned with a schema version field
- should memory mode support be delivered in phase one or postponed

## 11. Reference

This document supports:

- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)
