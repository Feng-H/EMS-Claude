# EMS Agent Phase 1 Task Breakdown

## 1. Purpose

This document breaks the phase-one EMS agent design into implementable work items.

It supports the following design documents:

- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)
- [ems_agent_schema_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_schema_design.md)
- [ems_agent_api_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_api_design.md)

The goal is to turn the current recommendation and audit assistant design into a staged delivery plan that engineering can execute.

## 2. Phase-One Delivery Goal

Phase one should deliver a usable management assistant for engineers and supervisors with two MVP capabilities:

- maintenance optimization recommendation
- repair cost and repair reasonableness audit

The first version should be:

- read-only against core business records
- evidence-based
- Chinese-first
- permission-scoped
- auditable

## 3. Workstream Overview

Recommended workstreams:

- backend agent foundation
- data model and migration
- retrieval and knowledge integration
- recommendation and audit analyzers
- frontend management assistant UI
- authentication and permission boundary enforcement
- testing and demo data
- rollout and observability

## 4. Milestone Plan

## 4.1 Milestone A: Agent foundation

Goal:

- make the backend able to receive agent requests and return scaffolded responses

Deliverables:

- `backend/internal/agent/` module structure
- base controller, service, repository, policy, prompt, and tool packages
- common response envelope
- trace ID generation
- session and artifact persistence flow

Suggested tasks:

- create `backend/internal/agent/controller`
- create `backend/internal/agent/service`
- create `backend/internal/agent/repository`
- create `backend/internal/agent/policy`
- create `backend/internal/agent/prompt`
- create `backend/internal/agent/tool`
- define common request context object carrying `user_id`, `factory_id`, role, and language
- define common response DTOs for summary, evidence, recommendation, and anomaly items
- add route registration under `/api/v1/agent`

Definition of done:

- a protected agent endpoint can be called successfully
- the endpoint writes `agent_sessions`
- the endpoint returns a structured envelope with `trace_id`

## 4.2 Milestone B: Schema and persistence

Goal:

- add the new storage needed by the agent

Deliverables:

- migrations for phase-one agent tables
- repository methods for sessions, artifacts, and evidence
- optional extensions to existing business tables

Suggested tasks:

- add migration for `equipment_manual_documents`
- add migration for `equipment_manual_chunks`
- add migration for `repair_cost_details`
- add migration for `equipment_runtime_snapshots`
- add migration for `agent_sessions`
- add migration for `agent_artifacts`
- add migration for `agent_evidence_links`
- add repository methods for create session
- add repository methods for create artifact
- add repository methods for batch insert evidence links
- evaluate whether `knowledge_articles` needs extra indexing or fields in phase one

Definition of done:

- migrations run successfully in database mode
- new tables are queryable
- agent request flow can persist session, artifact, and evidence placeholders

## 4.3 Milestone C: Retrieval and knowledge integration

Goal:

- make manuals and knowledge articles usable as evidence sources

Deliverables:

- manual ingestion path
- retrieval service for manuals and knowledge articles
- preference logic that prioritizes `knowledge_articles`

Suggested tasks:

- define manual ingestion process from uploaded file to chunk records
- add search support for `equipment_manual_chunks`
- add search support for `knowledge_articles`
- implement `search_manual_knowledge` tool
- rank knowledge-base best practices above generic manual chunks when both match
- return evidence objects with source table, source ID, excerpt, and score
- define a fallback path using PostgreSQL full-text search
- leave semantic retrieval upgrade path ready for later `pgvector`

Definition of done:

- the agent can return manual and knowledge-based evidence
- audit and recommendation flows can cite `knowledge_articles`
- search results are scoped to relevant equipment type or equipment context

## 4.4 Milestone D: Analyzer implementation

Goal:

- implement the deterministic business logic behind recommendations and audits

Deliverables:

- maintenance recommendation analyzer
- repair audit analyzer
- anomaly and recommendation scoring outputs

Suggested tasks:

- implement maintenance cycle comparison against manual guidance
- implement maintenance item gap detection against failure patterns
- implement peer-equipment benchmark comparison
- implement repeat-failure-after-maintenance detection
- implement repair cost deviation analysis
- implement best-practice mismatch detection using `knowledge_articles`
- define risk level thresholds such as `low`, `medium`, `high`
- define analyzer output schema before LLM summarization

Definition of done:

- analyzers can run without LLM output
- analyzers return structured recommendation and anomaly objects
- outputs include enough evidence references for later summarization

## 4.5 Milestone E: LLM summarization and localization

Goal:

- turn structured analyzer output into readable Chinese-first management output

Deliverables:

- prompt templates
- localization rules
- summary rendering logic

Suggested tasks:

- define default language as `zh-CN`
- create prompt templates for maintenance recommendation
- create prompt templates for repair audit
- create prompt templates for analysis explanation
- ensure terminology matches domestic factory management usage
- prevent the model from inventing evidence not present in analyzer output
- inject evidence snippets and structured metrics into summarization input
- add language metadata to response and artifact persistence

Definition of done:

- recommendation and audit summaries are readable in Chinese
- response language is explicit
- summaries do not introduce unsupported conclusions outside analyzer evidence

## 4.6 Milestone F: Permission boundary enforcement

Goal:

- ensure the agent never accesses unauthorized data

Deliverables:

- centralized scope injection
- policy checks on all tools
- forbidden-scope error handling

Suggested tasks:

- derive organization scope from authenticated user context
- inject `user_id`, `factory_id`, and role into every agent service request
- validate requested equipment IDs against accessible scope
- validate requested work order IDs against accessible scope
- block cross-factory queries centrally
- add `scope_summary` to response metadata
- persist resolved scope in `agent_sessions.input_snapshot`
- add tests for forbidden access cases

Definition of done:

- no agent tool can query outside resolved scope
- invalid scope requests return `FORBIDDEN_SCOPE`
- logs and artifacts reflect the resolved scope

## 4.7 Milestone G: Frontend management assistant

Goal:

- provide a usable PC-side entry point for engineers and supervisors

Deliverables:

- management assistant page
- scenario tabs
- result rendering panels

Suggested tasks:

- add `frontend/src/views/agent/ManagementAssistantView.vue`
- add route entry in the PC router
- add role-based visibility for engineer and supervisor
- create filters for factory, workshop, equipment type, equipment, and time range
- add scenario tabs for maintenance optimization, audit review, and analysis assistant
- add result sections for conclusion, evidence, risk, and detailed items
- add loading, empty, and error states
- add trace ID display for support and audit use

Definition of done:

- engineers and supervisors can submit a structured request
- the page renders recommendation and audit outputs correctly
- the UI feels like a management tool rather than a raw chat box

## 4.8 Milestone H: Testing and demo readiness

Goal:

- ensure the feature can be tested in both engineering and demo contexts

Deliverables:

- backend tests
- frontend interaction verification
- optional memory-mode demo data

Suggested tasks:

- add unit tests for analyzers
- add repository tests for agent persistence
- add API tests for permission failures and valid responses
- add mock data scenarios for repeated failure after maintenance
- add mock data scenarios for abnormal spare part cost
- add mock data scenarios for best-practice mismatch
- decide whether memory mode returns mock results or explicit unsupported-mode errors
- add documentation for demo scenarios

Definition of done:

- core analyzer logic has automated coverage
- permission checks are tested
- product demos can show at least one convincing recommendation case and one convincing audit case

## 4.9 Milestone I: Observability and rollout

Goal:

- make the phase-one agent safe to observe and iterate after release

Deliverables:

- operational metrics
- error visibility
- usage tracking

Suggested tasks:

- log trace IDs across controller and service layers
- log scenario type, runtime, and result status
- track counts for recommendation calls and audit calls
- track analyzer failure rate
- track unsupported storage mode rate if relevant
- add basic admin-facing reporting for agent session volume if needed

Definition of done:

- agent usage and failure patterns can be observed
- support teams can trace one result through `trace_id`

## 5. Team-by-Team Breakdown

## 5.1 Backend tasks

- create agent module structure
- add agent routes and DTOs
- implement service orchestration
- implement repositories for sessions, artifacts, and evidence
- implement retrieval tools
- implement analyzers
- enforce permission scope
- integrate LLM summarization

## 5.2 Data and migration tasks

- create migrations for new tables
- prepare indexing strategy
- define search-vector update strategy
- extend existing tables only where phase-one value is clear

## 5.3 Frontend tasks

- add management assistant page
- add structured request forms
- add result rendering components
- add role-based visibility
- add error and loading handling

## 5.4 QA and testing tasks

- validate permission boundaries
- validate Chinese output readability
- validate evidence rendering
- validate trace ID and artifact retrieval
- validate memory-mode behavior

## 6. Suggested Sprint Split

One reasonable split for implementation:

### Sprint 1

- milestone A
- milestone B
- milestone F

Outcome:

- backend skeleton, persistence, and permission boundary are in place

### Sprint 2

- milestone C
- milestone D
- partial milestone E

Outcome:

- recommendation and audit logic work with structured outputs and evidence

### Sprint 3

- milestone G
- remaining milestone E
- milestone H
- milestone I

Outcome:

- the feature is usable from the frontend and is testable and observable

## 7. Risks and Watchpoints

- retrieval quality may be weak if manual and knowledge content are not well structured
- cost audit quality will be limited if repair cost data is incomplete
- recommendations may be too generic if runtime data is not yet available
- memory mode can create extra maintenance burden if treated as full parity too early
- if permission injection is not centralized, agent tools may accidentally leak cross-factory data

## 8. Recommended MVP Cut Line

If scope must be reduced, keep only:

- maintenance recommendation for equipment type
- repair audit with repeated failure and cost deviation checks
- manual and knowledge evidence retrieval
- PC-side structured UI

Postpone if needed:

- open-ended analysis assistant
- memory-mode mock support
- semantic vector retrieval
- optional extensions to existing business tables beyond what is required

## 9. Acceptance Checklist

- agent endpoints exist and are protected
- sessions, artifacts, and evidence are persisted
- maintenance recommendation works for at least one realistic scenario
- repair audit works for at least one realistic scenario
- outputs are Chinese by default
- evidence includes knowledge or manual sources
- cross-factory access is blocked
- frontend engineers and supervisors can use the page
- trace IDs are visible in API responses and logs

## 10. Reference

This task breakdown supports:

- [ems_agent_design.md](/Users/apple/claudecode/EMS-Claude/ems_agent_design.md)
- [ems_agent_schema_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_schema_design.md)
- [ems_agent_api_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_api_design.md)
