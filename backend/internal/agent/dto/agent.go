package dto

import "time"

// =====================================================
// Common Agent DTOs
// =====================================================

// AgentResponseEnvelope represents the common response structure for all agent APIs
type AgentResponseEnvelope struct {
	Success       bool                   `json:"success"`
	TraceID       string                 `json:"trace_id"`
	Language      string                 `json:"language"`
	Scenario      string                 `json:"scenario"`
	ScopeSummary  map[string]interface{} `json:"scope_summary"`
	Summary       string                 `json:"summary"`
	RiskLevel     string                 `json:"risk_level"`
	ArtifactID    uint                   `json:"artifact_id,omitempty"`
	EvidenceCount int                    `json:"evidence_count"`
	Data          interface{}            `json:"data"`
}

// AgentError represents a structured error in agent responses
type AgentError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// AgentErrorEnvelope represents the common error structure
type AgentErrorEnvelope struct {
	Success bool       `json:"success"`
	TraceID string     `json:"trace_id"`
	Error   AgentErrDetail `json:"error"`
}

type AgentErrDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// EvidenceItem represents a single piece of evidence
type EvidenceItem struct {
	EvidenceType string  `json:"evidence_type"`
	SourceTable  string  `json:"source_table"`
	SourceID     uint    `json:"source_id"`
	Title        string  `json:"title"`
	Excerpt      string  `json:"excerpt"`
	Score        float64 `json:"score"`
}

// RecommendationItem represents a single recommendation
type RecommendationItem struct {
	Type        string `json:"type"`
	Target      string `json:"target"`
	TargetID    uint   `json:"target_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Reason      string `json:"reason"`
	Impact      string `json:"impact"`
}

// AnomalyItem represents a single anomaly finding
type AnomalyItem struct {
	AnomalyType     string `json:"anomaly_type"`
	Severity        string `json:"severity"`
	TargetType      string `json:"target_type"`
	TargetID        uint   `json:"target_id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	SuggestedAction string `json:"suggested_action"`
}

// =====================================================
// Maintenance Recommendation DTOs
// =====================================================

type TimeRange struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type MaintenanceRecommendRequest struct {
	FactoryID       uint      `json:"factory_id"`
	WorkshopID      uint      `json:"workshop_id"`
	EquipmentTypeID uint      `json:"equipment_type_id"`
	EquipmentIDs    []uint    `json:"equipment_ids"`
	TimeRange       TimeRange `json:"time_range"`
	Question        string    `json:"question"`
	Language        string    `json:"language"`
}

type MaintenanceRecommendData struct {
	CurrentPlan      interface{}          `json:"current_plan"`
	Recommendations []RecommendationItem `json:"recommendations"`
	ExpectedBenefits []string             `json:"expected_benefits"`
	Evidence        []EvidenceItem       `json:"evidence"`
}

// =====================================================
// Repair Audit DTOs
// =====================================================

type RepairAuditRequest struct {
	FactoryID       uint      `json:"factory_id"`
	WorkshopID      uint      `json:"workshop_id"`
	EquipmentTypeID uint      `json:"equipment_type_id"`
	TimeRange       TimeRange `json:"time_range"`
	AnomalyTypes    []string  `json:"anomaly_types"`
	Language        string    `json:"language"`
}

type RepairAuditData struct {
	Anomalies []AnomalyItem  `json:"anomalies"`
	Stats     interface{}    `json:"stats"`
	Evidence  []EvidenceItem `json:"evidence"`
}

// =====================================================
// Maintenance Audit DTOs
// =====================================================

type MaintenanceAuditRequest struct {
	FactoryID       uint      `json:"factory_id"`
	EquipmentTypeID uint      `json:"equipment_type_id"`
	TimeRange       TimeRange `json:"time_range"`
	Focus           []string  `json:"focus"`
	Language        string    `json:"language"`
}

type MaintenanceAuditData struct {
	AuditSummary     string         `json:"audit_summary"`
	Anomalies        []AnomalyItem  `json:"anomalies"`
	PlanComparisons  interface{}    `json:"plan_comparisons"`
	Evidence         []EvidenceItem `json:"evidence"`
}

// =====================================================
// Analysis Assistant DTOs
// =====================================================

type AnalyzeRequest struct {
	FactoryID  uint      `json:"factory_id"`
	WorkshopID uint      `json:"workshop_id"`
	Question   string    `json:"question"`
	TimeRange  TimeRange `json:"time_range"`
	Language   string    `json:"language"`
}

type AnalyzeData struct {
	KeyFindings       []string           `json:"key_findings"`
	MetricComparisons interface{}        `json:"metric_comparisons"`
	TopEntities       interface{}        `json:"top_entities"`
	Evidence          []EvidenceItem     `json:"evidence"`
	RecommendedActions []string           `json:"recommended_actions"`
}

// =====================================================
// Session & Artifact DTOs
// =====================================================

type AgentSessionResponse struct {
	ID            uint           `json:"id"`
	UserID        uint           `json:"user_id"`
	Scenario      string         `json:"scenario"`
	FactoryID     uint           `json:"factory_id"`
	WorkshopID    uint           `json:"workshop_id"`
	Language      string         `json:"language"`
	QueryText     string         `json:"query_text"`
	Status        string         `json:"status"`
	TraceID       string         `json:"trace_id"`
	CreatedAt     time.Time      `json:"created_at"`
	Artifacts     []uint         `json:"artifacts,omitempty"`
}

type AgentArtifactResponse struct {
	ID             uint           `json:"id"`
	SessionID      uint           `json:"session_id"`
	ArtifactType   string         `json:"artifact_type"`
	Title          string         `json:"title"`
	Summary        string         `json:"summary"`
	ResultJSON     interface{}    `json:"result_json"`
	RiskLevel      string         `json:"risk_level"`
	CreatedAt      time.Time      `json:"created_at"`
	Evidence       []EvidenceItem `json:"evidence,omitempty"`
	RelatedSession interface{}    `json:"related_session,omitempty"`
}
