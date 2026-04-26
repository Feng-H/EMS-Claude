package model

import (
	"time"

	"gorm.io/gorm"
)

// Base model with common fields
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// =====================================================
// User & Org Models
// =====================================================

type UserRole string

const (
	RoleAdmin       UserRole = "admin"
	RoleSupervisor  UserRole = "supervisor"
	RoleEngineer    UserRole = "engineer"
	RoleMaintenance UserRole = "maintenance"
	RoleOperator    UserRole = "operator"
)

const (
	ApprovalStatusPending  = "pending"
	ApprovalStatusApproved = "approved"
	ApprovalStatusRejected = "rejected"
)

type User struct {
	BaseModel
	Username           string   `json:"username" gorm:"size:100;uniqueIndex;not null"`
	PasswordHash       string   `json:"-" gorm:"size:255;not null"`
	Name               string   `json:"name" gorm:"size:100;not null"`
	Role               UserRole `json:"role" gorm:"size:20;not null"`
	Phone              string   `json:"phone" gorm:"size:20"`
	IsActive           bool     `json:"is_active" gorm:"default:true"`
	ApprovalStatus     string   `json:"approval_status" gorm:"size:20;default:'pending'"`
	MustChangePassword bool     `json:"must_change_password" gorm:"default:true"`
	FirstLogin         bool     `json:"first_login" gorm:"default:true"`
	FactoryID          *uint    `json:"factory_id"`
	Factory            *Factory `json:"factory,omitempty" gorm:"foreignKey:FactoryID"`
}

type Base struct {
	BaseModel
	Code      string    `json:"code" gorm:"size:50;uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Factories []Factory `json:"factories,omitempty" gorm:"foreignKey:BaseID"`
}

type Factory struct {
	BaseModel
	BaseID    uint       `json:"base_id" gorm:"not null"`
	Base      *Base      `json:"base,omitempty" gorm:"foreignKey:BaseID"`
	Code      string     `json:"code" gorm:"size:50;uniqueIndex;not null"`
	Name      string     `json:"name" gorm:"size:100;not null"`
	Workshops []Workshop `json:"workshops,omitempty" gorm:"foreignKey:FactoryID"`
}

type Workshop struct {
	BaseModel
	FactoryID uint        `json:"factory_id" gorm:"not null"`
	Factory   *Factory    `json:"factory,omitempty" gorm:"foreignKey:FactoryID"`
	Code      string      `json:"code" gorm:"size:50;uniqueIndex;not null"`
	Name      string      `json:"name" gorm:"size:100;not null"`
	Equipments []Equipment `json:"equipments,omitempty" gorm:"foreignKey:WorkshopID"`
}

// =====================================================
// Equipment Models
// =====================================================

type EquipmentType struct {
	BaseModel
	Name                string              `json:"name" gorm:"size:100;not null"`
	Category            string              `json:"category" gorm:"size:100"`
	Description         string              `json:"description" gorm:"type:text"`
	InspectionTemplateID *uint              `json:"inspection_template_id"`
	InspectionTemplates []InspectionTemplate `json:"inspection_templates,omitempty" gorm:"foreignKey:EquipmentTypeID"`
}

type Equipment struct {
	BaseModel
	Code             string    `json:"code" gorm:"size:50;uniqueIndex;not null"`
	Name             string    `json:"name" gorm:"size:200;not null"`
	Model            string    `json:"model" gorm:"size:100"`
	Spec             string    `json:"spec" gorm:"size:255"`
	Status           string    `json:"status" gorm:"size:20;default:'running'"`
	QRCode           string    `json:"qr_code" gorm:"size:100;uniqueIndex"`
	
	PurchasePrice    float64    `json:"purchase_price" gorm:"type:decimal(12,2);default:0"`
	PurchaseDate     *time.Time `json:"purchase_date"`
	ServiceLifeYears int        `json:"service_life_years" gorm:"default:10"`
	ScrapValue       float64    `json:"scrap_value" gorm:"type:decimal(12,2);default:0"`
	HourlyLoss       float64    `json:"hourly_loss" gorm:"type:decimal(10,2);default:0"`
	
	TypeID           uint           `json:"type_id"`
	Type             *EquipmentType `json:"type,omitempty" gorm:"foreignKey:TypeID"`
	WorkshopID       uint           `json:"workshop_id"`
	Workshop         *Workshop      `json:"workshop,omitempty" gorm:"foreignKey:WorkshopID"`
	
	DedicatedMaintenanceID *uint `json:"dedicated_maintenance_id"`
	DedicatedMaintenance   *User `json:"dedicated_maintenance,omitempty" gorm:"foreignKey:DedicatedMaintenanceID"`
}

// =====================================================
// Inspection Models
// =====================================================

type InspectionTemplate struct {
	BaseModel
	Name            string           `json:"name" gorm:"size:200;not null"`
	EquipmentTypeID uint             `json:"equipment_type_id" gorm:"not null"`
	EquipmentType   *EquipmentType   `json:"equipment_type,omitempty" gorm:"foreignKey:EquipmentTypeID"`
	Items           []InspectionItem `json:"items,omitempty" gorm:"foreignKey:TemplateID"`
}

type InspectionItem struct {
	BaseModel
	TemplateID    uint   `json:"template_id" gorm:"not null"`
	Name          string `json:"name" gorm:"size:255;not null"`
	Method        string `json:"method" gorm:"type:text"`
	Criteria      string `json:"criteria" gorm:"type:text"`
	Standard      string `json:"standard" gorm:"type:text"`
	SequenceOrder int    `json:"sequence_order" gorm:"default:0"`
}

const (
	InspectionPending    = "pending"
	InspectionInProgress = "in_progress"
	InspectionCompleted  = "completed"
	InspectionMissed     = "missed"
)

type InspectionTask struct {
	BaseModel
	EquipmentID   uint                `json:"equipment_id" gorm:"not null"`
	Equipment     *Equipment          `json:"equipment,omitempty" gorm:"foreignKey:EquipmentID"`
	TemplateID    uint                `json:"template_id" gorm:"not null"`
	Template      *InspectionTemplate `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
	AssignedTo    uint                `json:"assigned_to" gorm:"not null"`
	Assignee      *User               `json:"assignee,omitempty" gorm:"foreignKey:AssignedTo"`
	ScheduledDate time.Time           `json:"scheduled_date" gorm:"not null"`
	Status        string              `json:"status" gorm:"size:20;default:'pending'"`
	StartedAt     *time.Time          `json:"started_at"`
	CompletedAt   *time.Time          `json:"completed_at"`
	Latitude      *float64            `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude     *float64            `json:"longitude" gorm:"type:decimal(11,8)"`
	Records       []InspectionRecord `json:"records,omitempty" gorm:"foreignKey:TaskID"`
}

type InspectionRecord struct {
	BaseModel
	TaskID   uint            `json:"task_id" gorm:"not null"`
	ItemID   uint            `json:"item_id" gorm:"not null"`
	Item     *InspectionItem `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Result   string          `json:"result" gorm:"size:20;not null"`
	Remark   string          `json:"remark" gorm:"type:text"`
	PhotoURL string          `json:"photo_url" gorm:"size:500"`
	Images   []string        `json:"images" gorm:"type:text[]"`
}

// =====================================================
// Repair Models
// =====================================================

type RepairStatus string

const (
	RepairPending    RepairStatus = "pending"
	RepairAssigned   RepairStatus = "assigned"
	RepairInProgress RepairStatus = "in_progress"
	RepairTesting    RepairStatus = "testing"
	RepairConfirmed  RepairStatus = "confirmed"
	RepairAudited    RepairStatus = "audited"
	RepairClosed     RepairStatus = "closed"
)

type RepairOrder struct {
	BaseModel
	EquipmentID      uint       `json:"equipment_id" gorm:"not null"`
	Equipment        *Equipment `json:"equipment,omitempty" gorm:"foreignKey:EquipmentID"`
	FaultDescription string     `json:"fault_description" gorm:"type:text;not null"`
	FaultCode        string     `json:"fault_code" gorm:"size:50"`
	Priority         int        `json:"priority" gorm:"default:3"`
	Status           RepairStatus `json:"status" gorm:"size:20;default:'pending'"`
	ReporterID       uint       `json:"reporter_id" gorm:"not null"`
	Reporter         *User      `json:"reporter,omitempty" gorm:"foreignKey:ReporterID"`
	AssignedTo       *uint      `json:"assigned_to"`
	Assignee         *User      `json:"assignee,omitempty" gorm:"foreignKey:AssignedTo"`
	StartedAt        *time.Time `json:"started_at"`
	CompletedAt      *time.Time `json:"completed_at"`
	ConfirmedAt      *time.Time `json:"confirmed_at"`
	AuditedAt        *time.Time `json:"audited_at"`
	ClosedAt         *time.Time `json:"closed_at"`
	Solution         string     `json:"solution" gorm:"type:text"`
	Photos           []string   `json:"photos" gorm:"type:text[]"`
	Logs             []RepairLog `json:"logs,omitempty" gorm:"foreignKey:OrderID"`
	CostDetail       *RepairCostDetail `json:"cost_detail,omitempty" gorm:"foreignKey:OrderID"`
}

type RepairCostDetail struct {
	ID             uint      `json:"id" gorm:"primarykey"`
	OrderID        uint      `json:"order_id" gorm:"uniqueIndex;not null"`
	SparePartCost  float64   `json:"spare_part_cost" gorm:"type:decimal(12,2);default:0"`
	LaborCost      float64   `json:"labor_cost" gorm:"type:decimal(12,2);default:0"`
	OtherCost      float64   `json:"other_cost" gorm:"type:decimal(12,2);default:0"`
	Currency       string    `json:"currency" gorm:"size:10;default:'CNY'"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RepairLog struct {
	BaseModel
	OrderID uint         `json:"order_id" gorm:"not null"`
	Order   *RepairOrder `json:"-" gorm:"foreignKey:OrderID"`
	UserID  uint         `json:"user_id" gorm:"not null"`
	User    *User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Action  string       `json:"action" gorm:"size:100;not null"`
	Content string       `json:"content" gorm:"type:text"`
}

// =====================================================
// Maintenance Models
// =====================================================

type MaintenancePlan struct {
	BaseModel
	Name            string                `json:"name" gorm:"size:200;not null"`
	EquipmentTypeID uint                  `json:"equipment_type_id" gorm:"not null"`
	EquipmentType   *EquipmentType        `json:"equipment_type,omitempty" gorm:"foreignKey:EquipmentTypeID"`
	CycleDays       int                   `json:"cycle_days" gorm:"not null"`
	FlexibleDays    int                   `json:"flexible_days" gorm:"default:3"`
	WorkHours       float64               `json:"work_hours" gorm:"type:decimal(5,2);default:0"`
	Level           int                   `json:"level" gorm:"default:1"`
	Description     string                `json:"description" gorm:"type:text"`
	Items           []MaintenancePlanItem `json:"items,omitempty" gorm:"foreignKey:PlanID"`
}

type MaintenancePlanItem struct {
	BaseModel
	PlanID        uint   `json:"plan_id" gorm:"not null"`
	Name          string `json:"name" gorm:"size:255;not null"`
	Method        string `json:"method" gorm:"type:text"`
	Criteria      string `json:"criteria" gorm:"type:text"`
	Content       string `json:"content" gorm:"type:text"`
	Standard      string `json:"standard" gorm:"type:text"`
	SequenceOrder int    `json:"sequence_order" gorm:"default:0"`
}

const (
	MaintenancePending   = "pending"
	MaintenanceInProgress = "in_progress"
	MaintenanceCompleted = "completed"
	MaintenanceOverdue   = "overdue"
)

type MaintenanceTask struct {
	BaseModel
	EquipmentID   uint             `json:"equipment_id" gorm:"not null"`
	Equipment     *Equipment       `json:"equipment,omitempty" gorm:"foreignKey:EquipmentID"`
	PlanID        uint             `json:"plan_id" gorm:"not null"`
	Plan          *MaintenancePlan `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
	ScheduledDate string           `json:"scheduled_date" gorm:"size:10;not null"`
	DueDate       string           `json:"due_date" gorm:"size:10"`
	Status        string           `json:"status" gorm:"size:20;default:'pending'"`
	AssignedTo    uint             `json:"assigned_to" gorm:"not null"`
	Assignee      *User            `json:"assignee,omitempty" gorm:"foreignKey:AssignedTo"`
	StartedAt     *time.Time       `json:"started_at"`
	CompletedAt   *time.Time       `json:"completed_at"`
	ActualHours   float64          `json:"actual_hours" gorm:"type:decimal(5,2)"`
	Latitude      *float64         `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude     *float64         `json:"longitude" gorm:"type:decimal(11,8)"`
	Remark        string           `json:"remark" gorm:"type:text"`
	Records       []MaintenanceRecord `json:"records,omitempty" gorm:"foreignKey:TaskID"`
}

type MaintenanceRecord struct {
	BaseModel
	TaskID   uint                 `json:"task_id" gorm:"not null"`
	ItemID   uint                 `json:"item_id" gorm:"not null"`
	Item     *MaintenancePlanItem `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Result   string               `json:"result" gorm:"size:20;not null"`
	Remark   string               `json:"remark" gorm:"type:text"`
	PhotoURL string               `json:"photo_url" gorm:"size:500"`
}

// =====================================================
// Spare Part Models
// =====================================================

type SparePart struct {
	BaseModel
	Code          string `json:"code" gorm:"size:50;uniqueIndex;not null"`
	Name          string `json:"name" gorm:"size:200;not null"`
	Specification string `json:"specification" gorm:"size:255"`
	Unit          string `json:"unit" gorm:"size:20"`
	Category      string `json:"category" gorm:"size:50"`
	SafetyStock   int    `json:"safety_stock" gorm:"default:0"`
	FactoryID     *uint  `json:"factory_id"`
}

type SparePartInventory struct {
	BaseModel
	SparePartID uint       `json:"spare_part_id" gorm:"uniqueIndex:idx_spare_factory"`
	SparePart   *SparePart `json:"spare_part,omitempty" gorm:"foreignKey:SparePartID"`
	FactoryID   uint       `json:"factory_id" gorm:"uniqueIndex:idx_spare_factory"`
	Factory     *Factory   `json:"factory,omitempty" gorm:"foreignKey:FactoryID"`
	Quantity    int        `json:"quantity" gorm:"default:0"`
}

type SparePartConsumption struct {
	BaseModel
	SparePartID uint       `json:"spare_part_id" gorm:"not null"`
	SparePart   *SparePart `json:"spare_part,omitempty" gorm:"foreignKey:SparePartID"`
	OrderID     *uint      `json:"order_id"`
	Order       *RepairOrder `json:"-" gorm:"foreignKey:OrderID"`
	TaskID      *uint      `json:"task_id"`
	Quantity    int        `json:"quantity" gorm:"not null"`
	UserID      uint       `json:"user_id" gorm:"not null"`
	User        *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// =====================================================
// Knowledge & Document Models
// =====================================================

type KnowledgeArticle struct {
	BaseModel
	Title           string         `json:"title" gorm:"size:500;not null"`
	EquipmentTypeID *uint          `json:"equipment_type_id"`
	EquipmentType   *EquipmentType `json:"equipment_type,omitempty" gorm:"foreignKey:EquipmentTypeID"`
	FaultPhenomenon string         `json:"fault_phenomenon" gorm:"type:text"`
	CauseAnalysis   string         `json:"cause_analysis" gorm:"type:text"`
	Solution        string         `json:"solution" gorm:"type:text"`
	SourceType      string         `json:"source_type" gorm:"size:50"`
	SourceID        *uint          `json:"source_id"`
	Tags            []string       `json:"tags" gorm:"type:text[]"`
	CreatedBy       uint           `json:"created_by" gorm:"not null"`
	Creator         User           `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}

type ManualDocument struct {
	BaseModel
	Title           string         `json:"title" gorm:"size:200;not null"`
	EquipmentTypeID *uint          `json:"equipment_type_id"`
	EquipmentType   *EquipmentType `json:"equipment_type,omitempty" gorm:"foreignKey:EquipmentTypeID"`
	FilePath        string         `json:"file_path" gorm:"size:500"`
	Chunks          []ManualChunk  `json:"chunks,omitempty" gorm:"foreignKey:DocumentID"`
}

type ManualChunk struct {
	BaseModel
	DocumentID   uint            `json:"document_id" gorm:"not null"`
	SectionTitle string          `json:"section_title" gorm:"size:200"`
	Content      string          `json:"content" gorm:"type:text"`
	PageNumber   int             `json:"page_number"`
}

type EquipmentRuntimeSnapshot struct {
	BaseModel
	EquipmentID  uint       `json:"equipment_id" gorm:"not null;index"`
	Status       string     `json:"status" gorm:"size:20"`
	LoadFactor   float64    `json:"load_factor" gorm:"type:decimal(5,2)"`
	SnapshotDate string     `json:"snapshot_date" gorm:"size:10;index"`
}

// =====================================================
// Agent Brain Models (Phase 2 & 3)
// =====================================================

type AgentSkill struct {
	BaseModel
	Name                string   `json:"name" gorm:"size:200;not null"`
	Description         string   `json:"description" gorm:"type:text"`
	ApplicableTo        string   `json:"applicable_to" gorm:"type:jsonb"`
	ApplicableScenarios string   `json:"applicable_scenarios" gorm:"type:jsonb"`
	Steps               string   `json:"steps" gorm:"type:jsonb;not null"`
	Version             int      `json:"version" gorm:"default:1"`
	Status              string   `json:"status" gorm:"size:20;default:'active';index"`
	UsageCount          int      `json:"usage_count" gorm:"default:0"`
	SuccessRate         float64  `json:"success_rate" gorm:"type:decimal(5,4);default:0"`
	CreatedBy           string   `json:"created_by" gorm:"size:100"`
}

type AgentKnowledge struct {
	ID               string    `json:"id" gorm:"primarykey;size:100"`
	Title            string    `json:"title" gorm:"size:500;not null"`
	Type             string    `json:"type" gorm:"size:50"`
	Summary          string    `json:"summary" gorm:"type:text"`
	Details          string    `json:"details" gorm:"type:text"`
	RelatedSkillID   string    `json:"related_skill_id" gorm:"size:100"`
	Confidence       float64   `json:"confidence" gorm:"type:decimal(5,4);default:0"`
	Status           string    `json:"status" gorm:"size:20;default:'draft';index"`
	CreatedBy        string    `json:"created_by" gorm:"size:100"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	VerifiedBy       *uint     `json:"verified_by"`
	Verifier         *User     `json:"verifier,omitempty" gorm:"foreignKey:VerifiedBy"`
	VerifiedAt       *time.Time `json:"verified_at"`
	ReferencedCount  int       `json:"referenced_count" gorm:"default:0"`
	LastReferenced   *time.Time `json:"last_referenced"`
}

type AgentExperience struct {
	BaseModel
	UserID     uint     `json:"user_id" gorm:"not null;index"`
	User       *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Category   string   `json:"category" gorm:"size:50"`
	Content    string   `json:"content" gorm:"type:text"`
	Weight     float64  `json:"weight" gorm:"type:decimal(5,4);default:1.0"`
	DecayRate  float64  `json:"decay_rate" gorm:"type:decimal(5,4);default:0.01"`
	Status     string   `json:"status" gorm:"size:20;default:'active'"`
}

type AgentConversation struct {
	BaseModel
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title     string         `json:"title" gorm:"size:200"`
	Status    string         `json:"status" gorm:"size:20;default:'active'"`
	Messages  []AgentMessage `json:"messages,omitempty" gorm:"foreignKey:ConversationID"`
}

type AgentMessage struct {
	ID             uint      `json:"id" gorm:"primarykey"`
	ConversationID uint      `json:"conversation_id" gorm:"not null;index"`
	Conversation   *AgentConversation `json:"-" gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE;"`
	Role           string    `json:"role" gorm:"size:20;not null"`
	Content        string    `json:"content" gorm:"type:text;not null"`
	ImageURL       string    `json:"image_url" gorm:"size:500"`
	ToolCalls      string    `json:"tool_calls" gorm:"type:jsonb"`
	SkillID        string    `json:"skill_id" gorm:"size:100"`
	KnowledgeIDs   string    `json:"knowledge_ids" gorm:"type:jsonb;default:'[]'"`
	CreatedAt      time.Time `json:"created_at"`
}

type AgentPushSubscription struct {
	BaseModel
	UserID   uint   `json:"user_id" gorm:"not null;index"`
	PushType string `json:"push_type" gorm:"size:50;not null"`
	Enabled  bool   `json:"enabled" gorm:"default:true"`
	Scope    string `json:"scope" gorm:"type:jsonb"`
}

type AgentUsage struct {
	BaseModel
	SessionID      uint      `json:"session_id" gorm:"not null;index"`
	UserID         uint      `json:"user_id" gorm:"not null;index"`
	Scenario       string    `json:"scenario" gorm:"size:50"`
	Model          string    `json:"model" gorm:"size:100"`
	PromptTokens   int       `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens    int       `json:"total_tokens"`
	ResponseTimeMs int64     `json:"response_time_ms"`
}

type AgentSession struct {
	BaseModel
	UserID        uint            `json:"user_id" gorm:"not null;index"`
	Scenario      string          `json:"scenario" gorm:"size:100"`
	FactoryID     *uint           `json:"factory_id"`
	Language      string          `json:"language" gorm:"size:10"`
	InputSnapshot string          `json:"input_snapshot" gorm:"type:text"`
	Status        string          `json:"status" gorm:"size:20"`
	TraceID       string          `json:"trace_id" gorm:"size:100;uniqueIndex"`
	Artifacts     []AgentArtifact `json:"artifacts,omitempty" gorm:"foreignKey:SessionID"`
}

type AgentArtifact struct {
	BaseModel
	SessionID     uint                 `json:"session_id" gorm:"not null;index"`
	ArtifactType  string               `json:"artifact_type" gorm:"size:50"`
	Title         string               `json:"title" gorm:"size:200"`
	Summary       string               `json:"summary" gorm:"type:text"`
	ResultJSON    string               `json:"result_json" gorm:"type:text"`
	RiskLevel     string               `json:"risk_level" gorm:"size:20"`
	EvidenceLinks []AgentEvidenceLink `json:"evidence_links,omitempty" gorm:"foreignKey:ArtifactID"`
}

type AgentEvidenceLink struct {
	BaseModel
	ArtifactID   uint    `json:"artifact_id" gorm:"not null;index"`
	EvidenceType string  `json:"evidence_type" gorm:"size:50"`
	SourceTable  string  `json:"source_table" gorm:"size:100"`
	SourceID     uint    `json:"source_id"`
	Excerpt      string  `json:"excerpt" gorm:"type:text"`
	Score        float64 `json:"score" gorm:"type:decimal(5,4)"`
}
