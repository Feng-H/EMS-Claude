# EMS Agent Design

Related supporting documents:

- [EMS Agent PRD](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_prd.md)
- [EMS Agent Roadmap](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_roadmap.md)
- [Agent Schema Design](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_schema_design.md)
- [Agent Phase 2 Schema Extension](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase2_schema_extension.md)
- [Agent API Design](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_api_design.md)
- [Agent Phase 2 API Extension](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase2_api_design.md)
- [Agent Phase 1 Task Breakdown](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase1_task_breakdown.md)

## 1. Background

This agent is designed for the EMS project as a management assistant for equipment engineers and supervisors.

The initial focus is not on frontline execution, but on:

- maintenance recommendation
- maintenance and repair auditing
- operational analysis support

Its role is to help engineers and supervisors optimize maintenance content, review whether maintenance and repair actions and costs are reasonable, and support data-driven equipment governance.

## 2. Target Users

Primary users in phase one:

- equipment engineers
- equipment supervisors

The agent is positioned as a management assistant rather than an operator assistant.

## 3. Product Positioning

The phase-one agent is a device operations management assistant.

Its value is not generic chat, but turning scattered equipment data, maintenance history, repair records, spare part consumption, manuals, and knowledge into actionable and auditable recommendations.

## 4. Scope

### 4.1 What the agent should do

- recommend maintenance items and maintenance cycles based on manuals and historical equipment behavior
- audit whether maintenance, repair, and related costs appear reasonable
- provide engineering and management analysis support for equipment operations
- generate recommendation drafts, audit drafts, and analysis reports
- provide evidence-based conclusions rather than unsupported opinions

### 4.2 What the agent should not do

- directly modify official maintenance plans
- directly close repair orders
- directly confirm that a record is compliant or non-compliant without human review
- directly deduct spare part inventory or approve costs
- bypass permission boundaries to access cross-factory data
- replace on-site anti-fraud mechanisms such as QR-driven execution

The principle is: the agent gives recommendations and audit clues, while business staff confirm and execute.

## 5. Phase-One Core Scenarios

### 5.1 Maintenance Optimization Assistant

Input:

- equipment type
- equipment ID
- a group of similar equipment
- optional factory, workshop, and time range filters

Output:

- current maintenance plan summary
- suggested maintenance item additions, removals, or adjustments
- suggested cycle adjustments
- explanation of supporting evidence
- risk reminders such as over-maintenance, under-maintenance, or poor executability

Key evidence sources:

- manufacturer manual guidance
- historical failures and repair records
- performance of similar equipment
- maintenance execution quality and completion rate

### 5.2 Maintenance and Repair Audit Assistant

Input:

- time range
- factory or workshop
- equipment type
- target equipment or work order range

Output:

- suspicious maintenance or repair records
- anomaly types
- risk level
- evidence chain
- recommended review actions

Typical anomaly categories:

- repeat failures after maintenance
- excessive spare part cost
- abnormal labor hours
- maintenance items not matching actual failure modes
- high maintenance input with weak outcome
- unusual deviation from peer equipment

### 5.3 Analysis Assistant

This is for engineers and supervisors to ask structured management questions such as:

- why did repair cost increase this month in a workshop
- which equipment should be prioritized for governance
- which maintenance actions are most effective in reducing failures
- which equipment show high maintenance input but poor benefit

Outputs should include:

- natural language conclusions
- structured evidence
- compared populations
- key metrics
- abnormal findings

## 6. Interaction Model

Phase one should not rely on a pure chat box.

Recommended interaction mode:

- structured filters first
- guided templates second
- natural language enhancement third

Recommended entry tabs:

- Maintenance Optimization
- Audit Review
- Analysis Assistant

For each result, the UI should consistently show:

- conclusion
- evidence
- risk level
- recommended action

## 7. Implementation Principles

### 7.1 Recommendation, not direct execution

The agent should not directly change core business records in phase one.

It should produce:

- recommendations
- review suggestions
- draft reports
- draft plan changes

Human confirmation is required before any business changes take effect.

### 7.2 Evidence-based outputs

Every recommendation should include clear evidence from one or more of the following:

- manuals
- historical maintenance data
- historical repair data
- spare part and cost data
- benchmark comparison with similar equipment
- explicit audit or recommendation rules

### 7.3 Structured analysis first, LLM explanation second

The agent should not rely on the model alone to decide whether something is reasonable.

Recommended architecture:

- rules and metrics identify issues and recommendation candidates
- structured analysis computes comparisons, anomaly scores, and trends
- the LLM explains findings and writes readable recommendations

The LLM is the explanation layer, not the sole decision maker.

### 7.4 Localization-first output

The current EMS frontend and business documentation are primarily Chinese-oriented.

Phase one should treat Chinese as the default output language for:

- recommendations
- audit findings
- analysis reports
- evidence summaries

The prompt layer should explicitly enforce:

- Chinese output by default
- optional language switching based on user preference
- terminology aligned with domestic factory and equipment management usage

This helps engineers and supervisors consume outputs directly without secondary rewriting.

## 8. Technical Architecture

The backend should host the agent service rather than letting the frontend talk directly to an LLM.

Recommended module layout:

```text
backend/
  internal/
    agent/
      controller/
      service/
      analyzer/
      tool/
      policy/
      prompt/
      repository/
```

### 8.1 Module responsibilities

`controller`

- provide phase-one APIs

`service`

- orchestrate request flow
- run permission checks
- gather data
- trigger analyzers
- call LLM summarization
- persist audit logs

`analyzer`

- perform business logic independent of the LLM
- maintenance recommendation scoring
- anomaly detection
- benchmark comparison
- cost deviation analysis
- repeat failure detection

`tool`

- expose EMS data capabilities as controlled tools

`policy`

- apply role-based access control
- limit data by factory and organization scope
- support masking and export rules

`prompt`

- manage prompt templates
- prompts should organize expression rather than replace analysis logic
- enforce localization rules and response style for Chinese business scenarios

`repository`

- persist agent sessions
- persist outputs and evidence links
- support retrieval of manual and knowledge data

## 9. Tool Design

Phase one should expose a small set of well-bounded tools returning structured data.

Suggested tools:

- `get_equipment_context`
- `get_maintenance_profile`
- `get_repair_history`
- `get_sparepart_cost_summary`
- `get_equipment_benchmark`
- `search_manual_knowledge`
- `run_maintenance_recommendation`
- `run_audit_checks`

### 9.1 Tool responsibilities

`get_equipment_context`

- equipment profile
- equipment type
- organization location
- current status
- dedicated maintenance owner

`get_maintenance_profile`

- maintenance plans
- maintenance items
- maintenance cycles
- completion rate
- overdue rate

`get_repair_history`

- failures over a time range
- repair frequency
- repeat failures
- downtime indicators

`get_sparepart_cost_summary`

- spare part consumption
- cost distribution
- abnormal fluctuations

`get_equipment_benchmark`

- comparison with peer equipment of the same type

`search_manual_knowledge`

- retrieve relevant manual or knowledge-base evidence snippets
- prioritize project knowledge articles before falling back to manual content

`run_maintenance_recommendation`

- output structured recommendation results

`run_audit_checks`

- output structured anomaly findings and risk levels

### 9.2 Tool security boundary

All tool calls must be constrained by the current user context.

The service layer should inject:

- `user_id`
- `factory_id`
- allowed organization scope
- role information

Tool implementations should never rely on the LLM to pass correct scope parameters.

The agent should not know or infer equipment IDs outside its permission boundary.

This means:

- list queries must always be filtered by authorized organization scope
- single-record queries must verify accessibility before returning data
- cross-factory access must be blocked centrally rather than relying on UI controls

## 10. Data Foundation and Gaps

The current EMS data domains already include:

- equipment
- repair
- maintenance
- spare parts
- analytics
- knowledge

This is a strong starting point for a management assistant.

However, to make recommendation and audit outputs reliable, several data capabilities should be added.

### 10.1 Manual and standard data

The agent needs structured access to equipment manuals and manufacturer recommendations, not just uploaded files.

It should be able to retrieve:

- recommended maintenance cycles
- recommended maintenance items
- restrictions and cautions
- common faults and suggested treatments

### 10.2 Runtime and operating data

For better cycle recommendations, the system should gradually collect:

- runtime hours
- downtime hours
- load rate
- output or workload indicators when available

Without this, recommendations will be more experience-based than condition-based.

### 10.3 Cost detail data

Current spare part consumption data is helpful, but cost audit will be incomplete without:

- labor cost
- spare part cost
- outsourcing cost
- other repair-related costs

### 10.4 Agent audit persistence

The project should persist:

- who asked the agent
- what data was used
- what recommendation or audit result was returned
- whether the result was later confirmed or acted upon

### 10.5 Deep linkage with the existing knowledge base

The current project already includes a knowledge domain and related article model.

Phase one should not treat manuals as the only source of maintenance and repair guidance.

The agent should deeply integrate with the existing knowledge base by:

- searching `knowledge_articles` before or alongside manual retrieval
- comparing current maintenance or repair actions against existing best-practice articles
- identifying whether repeated faults already have recommended solutions in the knowledge base
- suggesting promotion of excellent repair cases into the knowledge base

In maintenance and repair audit scenarios, the agent should be able to explain:

- whether the repair approach aligns with known best practice
- whether an existing knowledge article was ignored
- whether a new case is worth promoting into the knowledge base

### 10.6 Memory mode compatibility

The current backend supports `storage.mode: memory`.

The agent module should explicitly define its behavior in memory mode.

Recommended phase-one options:

- option A: declare the full agent available only in database mode
- option B: provide a lightweight mock mode for demos and development

If mock support is provided, `InitMockData` should include representative agent scenarios, such as:

- the same equipment failing twice shortly after maintenance
- spare part consumption abnormally high compared with peer equipment
- a repair record inconsistent with a knowledge-base best practice

This makes the agent easier to demonstrate and test before full production data is available.

## 11. Suggested New Tables

Detailed schema design reference:

- [docs/ems_agent_schema_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_schema_design.md)

### 11.1 `equipment_manual_documents`

Store document metadata.

Suggested fields:

- `id`
- `equipment_type_id`
- `equipment_id` nullable
- `name`
- `file_url`
- `version`
- `source`
- `uploaded_by`
- `created_at`

### 11.2 `equipment_manual_chunks`

Store chunked manual content and retrieval references.

Suggested fields:

- `id`
- `document_id`
- `section_title`
- `content`
- `embedding_ref` or vector field
- `created_at`

For phase-one retrieval implementation:

- preferred long-term option: PostgreSQL with `pgvector`
- short-term transition option: PostgreSQL full-text search plus keyword matching

Since the current project uses PostgreSQL 15, vector retrieval should be planned as an incremental enhancement rather than an immediate hard dependency.

### 11.3 `repair_cost_details`

Store structured repair cost details.

Suggested fields:

- `id`
- `order_id`
- `labor_cost`
- `spare_part_cost`
- `outsource_cost`
- `other_cost`
- `currency`
- `created_at`

### 11.4 `equipment_runtime_snapshots`

Store equipment runtime snapshots.

Suggested fields:

- `id`
- `equipment_id`
- `date`
- `runtime_hours`
- `downtime_hours`
- `load_rate`
- `output_qty` nullable

### 11.5 `agent_sessions`

Store request-level agent session logs.

Suggested fields:

- `id`
- `user_id`
- `scenario`
- `factory_id`
- `query_text`
- `status`
- `created_at`

### 11.6 `agent_artifacts`

Store agent-generated artifacts.

Suggested fields:

- `id`
- `session_id`
- `artifact_type`
- `input_snapshot`
- `result_json`
- `summary`
- `confirmed_by` nullable
- `created_at`

Typical artifact types:

- `recommendation`
- `audit_report`
- `analysis_report`

### 11.7 `agent_evidence_links`

Link recommendations and audit outputs to their evidence.

Suggested fields:

- `id`
- `artifact_id`
- `evidence_type`
- `source_id`
- `excerpt`
- `score`

Suggested evidence types:

- `manual`
- `maintenance`
- `repair`
- `cost`
- `benchmark`
- `rule`

## 12. Phase-One Analysis Rules

Phase one should begin with deterministic rules and structured comparisons rather than advanced autonomous reasoning.

### 12.1 Maintenance recommendation rules

- compare current cycle with manual-recommended cycle
- when failure interval shortens, suggest shorter cycles or additional maintenance items
- when maintenance execution is stable and failures are low, suggest evaluating a looser cycle
- when frequent fault patterns are not covered by current maintenance items, suggest adding targeted tasks
- when planned workload is not executable in practice, suggest optimizing content instead of only adding work

### 12.2 Audit rules

- same-type equipment repair cost deviates significantly from the peer group
- failure repeats shortly after maintenance
- expensive spare parts are replaced too frequently
- the same fault recurs while the solution pattern never changes
- labor hours are abnormally high or low
- maintenance execution rate is low while the plan is overly dense
- low-value equipment has excessive maintenance input
- repair handling does not align with an existing knowledge-base best practice
- similar fault cases exist in the knowledge base but were not referenced

### 12.3 Analysis rules

- top high-failure equipment
- top high-cost equipment
- MTBF and MTTR comparison across same-type equipment
- correlation between maintenance input and failure reduction
- priority governance candidate list

## 13. Frontend Recommendation

Phase one should be PC-first.

A suggested new page:

- `frontend/src/views/agent/ManagementAssistantView.vue`

Recommended page structure:

- left panel for filters
- center panel for question input or analysis templates
- right panel for results

Recommended result sections:

- conclusion
- key evidence
- recommended actions
- detailed table

This should feel like a management tool rather than a generic chat interface.

## 14. Suggested APIs

Detailed API design reference:

- [docs/ems_agent_api_design.md](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_api_design.md)

Suggested phase-one API directions:

- `POST /api/v1/agent/maintenance/recommend`
- `POST /api/v1/agent/audit/repair`
- `POST /api/v1/agent/audit/maintenance`
- `POST /api/v1/agent/analyze`
- `GET /api/v1/agent/sessions/:id`
- `GET /api/v1/agent/artifacts/:id`

The response shape should include:

- result summary
- evidence list
- structured metrics
- anomaly list or recommendation list
- trace identifiers for auditability

Suggested response metadata:

- `language`
- `trace_id`
- `evidence_count`
- `scope_summary`

## 15. Delivery Strategy

### 15.1 Phase-one delivery sequence

1. Build read-only analysis flow
2. Add manual and knowledge retrieval support
3. Deliver two fixed templates first
4. Add open-ended but tool-constrained analysis later

Recommended sequencing detail:

1. Start with database mode support
2. Add optional memory-mode mock scenarios for demos
3. Use full-text retrieval first if vector capability is not yet ready
4. Upgrade to `pgvector` when manual and knowledge retrieval scale requires semantic ranking

### 15.2 Minimum viable scope

To control scope, the best phase-one MVP is:

- maintenance optimization recommendation
- repair cost anomaly audit

These two scenarios are aligned with current business goals and match the existing EMS data domains most directly.

## 16. Summary

The EMS phase-one agent should be a recommendation and audit assistant for engineers and supervisors.

It should:

- optimize maintenance content and cycles
- audit maintenance, repair, and cost reasonableness
- support engineering and management analysis

It should not directly execute sensitive business changes.

The recommended implementation path is:

- structured data retrieval
- rule-based and metric-based analysis
- evidence-linked LLM explanation
- auditable outputs with human confirmation

This keeps the agent useful, controllable, and aligned with EMS's strong workflow and accountability requirements.

For concrete implementation details, use these supporting documents together with this main design:

- [EMS Agent PRD](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_prd.md)
- [EMS Agent Roadmap](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_roadmap.md)
- [Agent Schema Design](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_schema_design.md)
- [Agent Phase 2 Schema Extension](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase2_schema_extension.md)
- [Agent API Design](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_api_design.md)
- [Agent Phase 2 API Extension](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase2_api_design.md)
- [Agent Phase 1 Task Breakdown](/Users/apple/claudecode/EMS-Claude/docs/ems_agent_phase1_task_breakdown.md)
