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

// User roles
type UserRole string

const (
	RoleAdmin       UserRole = "admin"
	RoleSupervisor  UserRole = "supervisor"
	RoleEngineer    UserRole = "engineer"
	RoleMaintenance UserRole = "maintenance"
	RoleOperator    UserRole = "operator"
)

// Organization structure
type Base struct {
	BaseModel
	Code string `json:"code" gorm:"uniqueIndex;not null"`
	Name string `json:"name" gorm:"not null"`
}

type Factory struct {
	BaseModel
	BaseID uint `json:"base_id" gorm:"not null"`
	Base   Base `json:"-" gorm:"foreignKey:BaseID"`
	Code   string `json:"code" gorm:"not null"`
	Name   string `json:"name" gorm:"not null"`
}

type Workshop struct {
	BaseModel
	FactoryID uint    `json:"factory_id" gorm:"not null"`
	Factory   Factory `json:"-" gorm:"foreignKey:FactoryID"`
	Code      string  `json:"code" gorm:"not null"`
	Name      string  `json:"name" gorm:"not null"`
}

// User
type User struct {
	BaseModel
	Username     string   `json:"username" gorm:"uniqueIndex;size:50;not null"`
	PasswordHash string   `json:"-" gorm:"size:255;not null"`
	Name         string   `json:"name" gorm:"size:100;not null"`
	Role         UserRole `json:"role" gorm:"type:varchar(20);not null;default:'operator'"`
	FactoryID    *uint    `json:"factory_id"`
	Factory      *Factory `json:"factory,omitempty" gorm:"foreignKey:FactoryID"`
	Phone        string   `json:"phone" gorm:"size:20"`
	IsActive     bool     `json:"is_active" gorm:"default:true"`
}

// Equipment
type Equipment struct {
	BaseModel
	Code                    string     `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Name                    string     `json:"name" gorm:"size:200;not null"`
	TypeID                  uint       `json:"type_id" gorm:"not null"`
	Type                    EquipmentType `json:"type,omitempty" gorm:"foreignKey:TypeID"`
	WorkshopID              uint       `json:"workshop_id" gorm:"not null"`
	Workshop                Workshop   `json:"workshop,omitempty" gorm:"foreignKey:WorkshopID"`
	QRCode                  string     `json:"qr_code" gorm:"uniqueIndex;size:255;not null"`
	Spec                    string     `json:"spec" gorm:"type:text"`
	PurchaseDate            *time.Time `json:"purchase_date"`
	Status                  string     `json:"status" gorm:"size:20;default:'running'"`
	DedicatedMaintenanceID  *uint      `json:"dedicated_maintenance_id"`
	DedicatedMaintenance    *User      `json:"dedicated_maintenance,omitempty" gorm:"foreignKey:DedicatedMaintenanceID"`
}

type EquipmentType struct {
	BaseModel
	Name                 string  `json:"name" gorm:"size:100;not null"`
	Category             string  `json:"category" gorm:"size:50"`
	InspectionTemplateID *uint   `json:"inspection_template_id"`
}

// Inspection
type InspectionTemplate struct {
	BaseModel
	Name            string             `json:"name" gorm:"size:200;not null"`
	EquipmentTypeID uint               `json:"equipment_type_id" gorm:"not null"`
	EquipmentType   EquipmentType      `json:"equipment_type,omitempty" gorm:"foreignKey:EquipmentTypeID"`
	Items           []InspectionItem   `json:"items,omitempty" gorm:"foreignKey:TemplateID"`
}

type InspectionItem struct {
	BaseModel
	TemplateID    uint   `json:"template_id" gorm:"not null"`
	Name          string `json:"name" gorm:"size:200;not null"`
	Method        string `json:"method" gorm:"size:500"`
	Criteria      string `json:"criteria" gorm:"size:500"`
	SequenceOrder int    `json:"sequence_order" gorm:"default:0"`
}

type InspectionTaskStatus string

const (
	InspectionPending    InspectionTaskStatus = "pending"
	InspectionInProgress InspectionTaskStatus = "in_progress"
	InspectionCompleted  InspectionTaskStatus = "completed"
	InspectionOverdue    InspectionTaskStatus = "overdue"
)

type InspectionTask struct {
	BaseModel
	EquipmentID  uint                 `json:"equipment_id" gorm:"not null"`
	Equipment    Equipment            `json:"equipment,omitempty" gorm:"foreignKey:EquipmentID"`
	TemplateID   uint                 `json:"template_id" gorm:"not null"`
	Template     InspectionTemplate   `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
	AssignedTo   uint                 `json:"assigned_to" gorm:"not null"`
	Assignee     *User                `json:"assignee,omitempty" gorm:"foreignKey:AssignedTo"`
	ScheduledDate time.Time          `json:"scheduled_date" gorm:"not null"`
	Status       InspectionTaskStatus `json:"status" gorm:"size:20;default:'pending'"`
	StartedAt    *time.Time           `json:"started_at"`
	CompletedAt  *time.Time           `json:"completed_at"`
	Latitude     *float64             `json:"latitude"`
	Longitude    *float64             `json:"longitude"`
	Records      []InspectionRecord   `json:"records,omitempty" gorm:"foreignKey:TaskID"`
}

type InspectionRecord struct {
	BaseModel
	TaskID  uint   `json:"task_id" gorm:"not null"`
	Task    InspectionTask `json:"-" gorm:"foreignKey:TaskID"`
	ItemID  uint   `json:"item_id" gorm:"not null"`
	Result  string `json:"result" gorm:"size:10;not null"` // OK or NG
	Remark  string `json:"remark" gorm:"type:text"`
	PhotoURL string `json:"photo_url" gorm:"size:500"`
}

// Repair
type RepairStatus string

const (
	RepairPending   RepairStatus = "pending"
	RepairAssigned  RepairStatus = "assigned"
	RepairInProgress RepairStatus = "in_progress"
	RepairTesting   RepairStatus = "testing"
	RepairConfirmed RepairStatus = "confirmed"
	RepairAudited   RepairStatus = "audited"
	RepairClosed    RepairStatus = "closed"
)

type RepairOrder struct {
	BaseModel
	EquipmentID       uint        `json:"equipment_id" gorm:"not null"`
	Equipment         Equipment   `json:"equipment,omitempty" gorm:"foreignKey:EquipmentID"`
	FaultDescription string      `json:"fault_description" gorm:"type:text;not null"`
	ReporterID        uint        `json:"reporter_id" gorm:"not null"`
	Reporter          User        `json:"reporter,omitempty" gorm:"foreignKey:ReporterID"`
	AssignedTo        *uint       `json:"assigned_to"`
	Assignee          *User       `json:"assignee,omitempty" gorm:"foreignKey:AssignedTo"`
	Status            RepairStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Priority          int         `json:"priority" gorm:"default:3"` // 1=high, 2=medium, 3=low
	Photos            []string    `json:"photos" gorm:"type:text[]"`
	FaultCode         string      `json:"fault_code" gorm:"size:50"`
	Solution          string      `json:"solution" gorm:"type:text"`
	StartedAt         *time.Time  `json:"started_at"`
	CompletedAt       *time.Time  `json:"completed_at"`
	ConfirmedAt       *time.Time  `json:"confirmed_at"`
	AuditedAt         *time.Time  `json:"audited_at"`
	Logs              []RepairLog `json:"logs,omitempty" gorm:"foreignKey:OrderID"`
}

type RepairLog struct {
	BaseModel
	OrderID uint   `json:"order_id" gorm:"not null"`
	Order   RepairOrder `json:"-" gorm:"foreignKey:OrderID"`
	UserID  uint   `json:"user_id" gorm:"not null"`
	User    User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Action  string `json:"action" gorm:"size:50"`
	Content string `json:"content" gorm:"type:text"`
}

// Maintenance
type MaintenancePlan struct {
	BaseModel
	Name            string                `json:"name" gorm:"size:200;not null"`
	EquipmentTypeID uint                  `json:"equipment_type_id" gorm:"not null"`
	EquipmentType   *EquipmentType        `json:"equipment_type,omitempty" gorm:"foreignKey:EquipmentTypeID"`
	Level           int                   `json:"level" gorm:"not null"` // 1, 2, 3
	CycleDays       int                   `json:"cycle_days" gorm:"not null"`
	FlexibleDays    int                   `json:"flexible_days" gorm:"default:3"`
	WorkHours       float64               `json:"work_hours"`
	Items           []MaintenancePlanItem `json:"items,omitempty" gorm:"foreignKey:PlanID"`
}

type MaintenanceTaskStatus string

const (
	MaintenancePending   MaintenanceTaskStatus = "pending"
	MaintenanceInProgress MaintenanceTaskStatus = "in_progress"
	MaintenanceCompleted MaintenanceTaskStatus = "completed"
	MaintenanceOverdue    MaintenanceTaskStatus = "overdue"
)

type MaintenanceTask struct {
	BaseModel
	PlanID        uint                 `json:"plan_id" gorm:"not null"`
	Plan          MaintenancePlan      `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
	EquipmentID   uint                 `json:"equipment_id" gorm:"not null"`
	Equipment     Equipment            `json:"equipment,omitempty" gorm:"foreignKey:EquipmentID"`
	AssignedTo    uint                 `json:"assigned_to" gorm:"not null"`
	Assignee      *User                `json:"assignee,omitempty" gorm:"foreignKey:AssignedTo"`
	ScheduledDate string               `json:"scheduled_date" gorm:"size:10;not null"` // YYYY-MM-DD
	DueDate       string               `json:"due_date" gorm:"size:10;not null"`        // YYYY-MM-DD
	Status        MaintenanceTaskStatus `json:"status" gorm:"size:20;default:'pending'"`
	StartedAt     *time.Time           `json:"started_at"`
	CompletedAt   *time.Time           `json:"completed_at"`
	ActualHours   float64              `json:"actual_hours"`
	Latitude      *float64             `json:"latitude"`
	Longitude     *float64             `json:"longitude"`
	Remark        string               `json:"remark" gorm:"type:text"`
	Records       []MaintenanceRecord  `json:"records,omitempty" gorm:"foreignKey:TaskID"`
}

type MaintenancePlanItem struct {
	BaseModel
	PlanID        uint            `json:"plan_id" gorm:"not null"`
	Plan          MaintenancePlan `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
	Name          string          `json:"name" gorm:"size:200;not null"`
	Method        string          `json:"method" gorm:"type:text"`
	Criteria      string          `json:"criteria" gorm:"type:text"`
	SequenceOrder int             `json:"sequence_order" gorm:"default:0"`
}

type MaintenanceRecord struct {
	BaseModel
	TaskID     uint                 `json:"task_id" gorm:"not null"`
	Task       *MaintenanceTask     `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	ItemID     uint                 `json:"item_id" gorm:"not null"`
	Item       *MaintenancePlanItem `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Result     string               `json:"result" gorm:"size:10;not null"` // OK/NG
	Remark     string               `json:"remark" gorm:"type:text"`
	PhotoURL   string               `json:"photo_url" gorm:"size:500"`
}

// Spare Parts
type SparePart struct {
	BaseModel
	Code         string  `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Name         string  `json:"name" gorm:"size:200;not null"`
	Specification string `json:"specification" gorm:"type:text"`
	Unit         string  `json:"unit" gorm:"size:20"`
	FactoryID    *uint   `json:"factory_id"`
	SafetyStock  int     `json:"safety_stock" gorm:"default:0"`
}

type SparePartInventory struct {
	BaseModel
	SparePartID uint       `json:"spare_part_id" gorm:"not null"`
	SparePart   SparePart  `json:"spare_part,omitempty" gorm:"foreignKey:SparePartID"`
	FactoryID   uint       `json:"factory_id" gorm:"not null"`
	Factory     Factory    `json:"factory,omitempty" gorm:"foreignKey:FactoryID"`
	Quantity    int        `json:"quantity" gorm:"default:0"`
}

type SparePartConsumption struct {
	BaseModel
	SparePartID uint         `json:"spare_part_id" gorm:"not null"`
	SparePart   SparePart    `json:"spare_part,omitempty" gorm:"foreignKey:SparePartID"`
	OrderID     *uint        `json:"order_id"` // repair order
	Order       *RepairOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	TaskID      *uint        `json:"task_id"` // maintenance task
	Task        *MaintenanceTask `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	Quantity    int           `json:"quantity" gorm:"not null"`
	UserID      uint          `json:"user_id" gorm:"not null"`
	User        User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Knowledge Base
type KnowledgeArticle struct {
	BaseModel
	Title             string   `json:"title" gorm:"size:200;not null"`
	EquipmentTypeID   *uint    `json:"equipment_type_id"`
	EquipmentType     *EquipmentType `json:"equipment_type,omitempty" gorm:"foreignKey:EquipmentTypeID"`
	FaultPhenomenon   string   `json:"fault_phenomenon" gorm:"type:text"`
	CauseAnalysis     string   `json:"cause_analysis" gorm:"type:text"`
	Solution          string   `json:"solution" gorm:"type:text"`
	SourceType        string   `json:"source_type" gorm:"size:20"` // repair, manual, other
	SourceID          *uint    `json:"source_id"`
	Tags              []string `json:"tags" gorm:"type:text[]"`
	CreatedBy         uint     `json:"created_by" gorm:"not null"`
	Creator           User     `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}
