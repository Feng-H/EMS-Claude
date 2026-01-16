package dto

import "time"

// =====================================================
// Repair Order DTOs
// =====================================================

// RepairStatus represents the status of a repair order
type RepairStatus string

const (
	RepairStatusPending    RepairStatus = "pending"     // 待派单
	RepairStatusAssigned   RepairStatus = "assigned"    // 已派单
	RepairStatusInProgress RepairStatus = "in_progress" // 维修中
	RepairStatusTesting    RepairStatus = "testing"     // 待测试
	RepairStatusConfirmed  RepairStatus = "confirmed"   // 待审核
	RepairStatusAudited    RepairStatus = "audited"     // 已审核
	RepairStatusClosed      RepairStatus = "closed"      // 已关闭
)

// CreateRepairRequest creates a new repair order
type CreateRepairRequest struct {
	EquipmentID      uint     `json:"equipment_id" binding:"required"`
	FaultDescription string   `json:"fault_description" binding:"required"`
	FaultCode        string   `json:"fault_code"`
	Photos           []string `json:"photos"`
	Priority         int      `json:"priority"` // 1=高,2=中,3=低
}

// RepairOrderResponse represents a repair order in API responses
type RepairOrderResponse struct {
	ID               uint         `json:"id"`
	EquipmentID      uint         `json:"equipment_id"`
	EquipmentCode    string       `json:"equipment_code,omitempty"`
	EquipmentName    string       `json:"equipment_name,omitempty"`
	FaultDescription string       `json:"fault_description"`
	FaultCode        string       `json:"fault_code,omitempty"`
	ReporterID       uint         `json:"reporter_id"`
	ReporterName     string       `json:"reporter_name,omitempty"`
	AssignedTo       uint         `json:"assigned_to,omitempty"`
	AssigneeName     string       `json:"assignee_name,omitempty"`
	Status           RepairStatus `json:"status"`
	Priority         int          `json:"priority"`
	Photos           []string     `json:"photos,omitempty"`
	Solution         string       `json:"solution,omitempty"`
	SpareParts       string       `json:"spare_parts,omitempty"`
	ActualHours      float64      `json:"actual_hours,omitempty"`
	CreatedAt        time.Time    `json:"created_at"`
	StartedAt        *time.Time   `json:"started_at,omitempty"`
	CompletedAt      *time.Time   `json:"completed_at,omitempty"`
	ConfirmedAt      *time.Time   `json:"confirmed_at,omitempty"`
	AuditedAt        *time.Time   `json:"audited_at,omitempty"`
}

// RepairOrderListResponse represents a paginated list of repair orders
type RepairOrderListResponse struct {
	Total int64                    `json:"total"`
	Items []RepairOrderResponse    `json:"items"`
}

// RepairOrderDetailResponse represents a repair order with logs
type RepairOrderDetailResponse struct {
	RepairOrderResponse
	Logs []RepairLogResponse `json:"logs,omitempty"`
}

// =====================================================
// Repair Assignment DTOs
// =====================================================

// AssignRepairRequest assigns a repair order to a technician
type AssignRepairRequest struct {
	AssignTo uint `json:"assign_to" binding:"required"`
}

// =====================================================
// Repair Execution DTOs
// =====================================================

// StartRepairRequest starts a repair task
type StartRepairRequest struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

// UpdateRepairRequest updates repair progress
type UpdateRepairRequest struct {
	Solution    string   `json:"solution"`
	SpareParts  string   `json:"spare_parts"`
	ActualHours float64  `json:"actual_hours"`
	Photos      []string `json:"photos"`
	NextStatus  string   `json:"next_status"` // testing, confirmed
}

// ConfirmRepairRequest confirms a completed repair
type ConfirmRepairRequest struct {
	Accepted  bool   `json:"accepted"`
	Comment   string `json:"comment"`
	Photos    []string `json:"photos"`
}

// AuditRepairRequest audits a confirmed repair
type AuditRepairRequest struct {
	Approved   bool   `json:"approved"`
	Comment    string `json:"comment"`
	ActualHours float64 `json:"actual_hours,omitempty"`
}

// =====================================================
// Repair Log DTOs
// =====================================================

// CreateRepairLogRequest creates a new repair log
type CreateRepairLogRequest struct {
	Action  string `json:"action" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// RepairLogResponse represents a repair log in API responses
type RepairLogResponse struct {
	ID        uint      `json:"id"`
	OrderID   uint      `json:"order_id"`
	UserID    uint      `json:"user_id"`
	UserName  string    `json:"user_name,omitempty"`
	Action    string    `json:"action"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// =====================================================
// Repair Statistics DTOs
// =====================================================

// RepairStatistics represents repair statistics
type RepairStatistics struct {
	TotalOrders      int64   `json:"total_orders"`
	PendingOrders    int64   `json:"pending_orders"`
	InProgressOrders int64   `json:"in_progress_orders"`
	CompletedOrders  int64   `json:"completed_orders"`
	TodayCompleted   int64   `json:"today_completed"`
	TodayCreated     int64   `json:"today_created"`
	AvgRepairTime    float64 `json:"avg_repair_time"` // hours
	AvgResponseTime  float64 `json:"avg_response_time"` // minutes
}

// MyRepairStatistics represents current user's repair statistics
type MyRepairStatistics struct {
	PendingCount    int `json:"pending_count"`
	InProgressCount int `json:"in_progress_count"`
	CompletedCount  int `json:"completed_count"`
	TodayCompleted  int `json:"today_completed"`
}

// RepairTaskQuery represents query parameters for repair tasks
type RepairTaskQuery struct {
	Status     string `form:"status"`
	Priority   int    `form:"priority"`
	AssignedTo uint   `form:"assigned_to"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	DateFrom   string `form:"date_from"`
	DateTo     string `form:"date_to"`
}
