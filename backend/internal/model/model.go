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

// User model
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

// Organization models
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

// Equipment models
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
	Name             string    `json:"name" gorm:"size:100;not null"`
	Model            string    `json:"model" gorm:"size:100"`
	Spec             string    `json:"spec" gorm:"size:255"`
	Status           string    `json:"status" gorm:"size:20;default:'running'"`
	QRCode           string    `json:"qr_code" gorm:"size:100;uniqueIndex"`
	
	// 资产财务与寿命字段
	PurchasePrice    float64    `json:"purchase_price" gorm:"type:decimal(12,2);default:0"`
	PurchaseDate     *time.Time `json:"purchase_date"`
	ServiceLifeYears int        `json:"service_life_years" gorm:"default:10"`
	ScrapValue       float64    `json:"scrap_value" gorm:"type:decimal(12,2);default:0"`
	HourlyLoss       float64    `json:"hourly_loss" gorm:"type:decimal(10,2);default:0"`
	
	// Relationships
	TypeID           uint           `json:"type_id"`
	Type             *EquipmentType `json:"type,omitempty" gorm:"foreignKey:TypeID"`
	WorkshopID       uint           `json:"workshop_id"`
	Workshop         *Workshop      `json:"workshop,omitempty" gorm:"foreignKey:WorkshopID"`
	
	// Dedicated maintenance person
	DedicatedMaintenanceID *uint `json:"dedicated_maintenance_id"`
	DedicatedMaintenance   *User `json:"dedicated_maintenance,omitempty" gorm:"foreignKey:DedicatedMaintenanceID"`
}

// Inspection models
type InspectionTemplate struct {
	BaseModel
	Name            string           `json:"name" gorm:"size:100;not null"`
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
	Result   string          `json:"result" gorm:"size:20;not null"` // OK, NG
	Remark   string          `json:"remark" gorm:"type:text"`
	PhotoURL string          `json:"photo_url" gorm:"size:500"`
	Images   []string        `json:"images" gorm:"type:text[]"`
}

// Repair models
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

// Maintenance models
type MaintenancePlan struct {
	BaseModel
	Name            string                `json:"name" gorm:"size:100;not null"`
	EquipmentTypeID uint                  `json:"equipment_type_id" gorm:"not null"`
	EquipmentType   *EquipmentType        `json:"equipment_type,omitempty" gorm:"foreignKey:EquipmentTypeID"`
	CycleDays       int                   `json:"cycle_days" gorm:"not null"`
	FlexibleDays    int                   `json:"flexible_days" gorm:"default:0"`
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
	Result   string               `json:"result" gorm:"size:20;not null"` // 合格, 不合格
	Remark   string               `json:"remark" gorm:"type:text"`
	PhotoURL string               `json:"photo_url" gorm:"size:500"`
}

// SparePart models
type SparePart struct {
	BaseModel
	Code          string `json:"code" gorm:"size:50;uniqueIndex;not null"`
	Name          string `json:"name" gorm:"size:100;not null"`
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

// Knowledge models
type KnowledgeArticle struct {
	BaseModel
	Title           string         `json:"title" gorm:"size:500;not null"`
	EquipmentTypeID *uint          `json:"equipment_type_id"`
	EquipmentType   *EquipmentType `json:"equipment_type,omitempty" gorm:"foreignKey:EquipmentTypeID"`
	FaultPhenomenon string         `json:"fault_phenomenon" gorm:"type:text"`
	CauseAnalysis   string         `json:"cause_analysis" gorm:"type:text"`
	Solution        string         `json:"solution" gorm:"type:text"`
	SourceType      string         `json:"source_type" gorm:"size:50"` // repair, manual, expert_summary
	SourceID        *uint          `json:"source_id"`
	Tags            []string       `json:"tags" gorm:"type:text[]"`
	CreatedBy       uint           `json:"created_by" gorm:"not null"`
	Creator         User           `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}
