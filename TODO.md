# EMS Agent Phase 1 Progress

## 📋 TODO List

### Sprint 1: Backend Foundation, Schema & Security
- [x] **Milestone A: Agent foundation**
    - [x] Create `backend/internal/agent/` module structure (controller, service, repository, policy, prompt, tool).
    - [x] Define common request context and response DTOs.
    - [x] Implement trace ID generation.
    - [x] Register routes under `/api/v1/agent`.
    - [x] Implement base session and artifact persistence flow.
- [x] **Milestone B: Schema and persistence**
    - [x] Add migrations for `equipment_manual_documents` and `equipment_manual_chunks`.
    - [x] Add migration for `repair_cost_details`.
    - [x] Add migration for `equipment_runtime_snapshots`.
    - [x] Add migration for `agent_sessions`, `agent_artifacts`, and `agent_evidence_links`.
    - [x] Implement repository methods for the new tables (supporting both DB and Memory modes).
- [x] **Milestone F: Permission boundary enforcement**
    - [x] Derive organization scope from auth context.
    - [x] Inject `user_id`, `factory_id`, and role into service requests.
    - [x] Validate equipment and work order IDs against scope.
    - [x] Block cross-factory queries centrally (via `PolicyService`).

### Sprint 2: Logic & Knowledge
- [x] **Milestone C: Retrieval and knowledge integration**
    - [x] Implement `RetrievalTool` supporting knowledge articles and manual chunks.
    - [x] Add search support for `equipment_manual_chunks` and `knowledge_articles`.
    - [x] Implement `search_manual_knowledge` tool (integrated in `RetrievalTool`).
- [x] **Milestone D: Analyzer implementation**
    - [x] Implement maintenance recommendation analyzer (basic rule-based).
    - [x] Implement repair audit analyzer (repeat-failure detection).
    - [x] Implement supporting tools (`MaintenanceTool`, `RepairTool`).
- [x] **Milestone E: LLM summarization & Localization**
    - [x] Create `llm` package with OpenAI-compatible client.
    - [x] Create `PromptTool` with Chinese templates for recommendations and audits.
    - [x] Integrated LLM summarization into `AgentService`.

### Sprint 3: UI & Rollout
- [x] **Milestone G: Frontend management assistant**
    - [x] Add `ManagementAssistantView.vue` and route.
    - [x] Implement filters and scenario tabs (Maintenance optimization, Repair audit).
    - [x] Render results (conclusion, evidence, risk, items) with professional industrial styling.
- [x] **Milestone H: Testing and demo readiness**
    - [x] Implement backend session and artifact persistence logic.
    - [x] Add mock data for representative agent scenarios (repeat failures, cost deviation).
    - [x] Support historical analysis loading in frontend.
- [x] **Milestone I: Observability and rollout**
    - [x] Implement `AgentUsage` model to track tokens and response time.
    - [x] Integrated usage logging into all agent scenarios.

---

## 📈 Progress Record

### 2026-04-25
- 🚀 Initialized `TODO.md` based on `ems_agent_design.md` and supporting documents.
- 🔍 Reviewed project architecture and confirmed the plan for Sprint 1.
- ✅ **Milestone A completed**: Scaffolded backend agent module structure, DTOs, trace utility, controller, and service. Registered routes in `main.go`.
- ✅ **Milestone B completed**: Defined agent models (GORM) and implemented repositories for both database and memory modes. Updated `AutoMigrate` in `main.go` and `Store` in `pkg/memory`.
- ✅ **Milestone F completed**: Implemented `PolicyService` for organization scope derivation and validation. Injected authorized context into agent service requests.
- ✅ **Milestone C completed**: Implemented `RetrievalTool` and `MaintenanceTool` for data gathering from knowledge base and maintenance plans.
- ✅ **Milestone D completed**: Implemented `MaintenanceAnalyzer` and `RepairAuditAnalyzer` with deterministic logic. Integrated analyzers and supporting tools into `AgentService`.
- ✅ **Milestone E completed**: Implemented `llm` package and `PromptTool`. Integrated LLM summarization into the agent workflow for Chinese-first professional outputs.
- ✅ **Milestone G completed**: Built the PC Management Assistant frontend with dual-pane layout, supporting Maintenance Recommendation and Repair Audit scenarios.
- ✅ **Milestone H completed**: Implemented full persistence flow for Agent sessions and artifacts. Added "story-driven" mock data for realistic demos of repeat-failure and cost-deviation auditing.
- ✅ **Milestone I completed**: Implemented `AgentUsage` tracking to monitor LLM performance and cost. Phase 1 is now feature-complete and ready for production-like verification.
