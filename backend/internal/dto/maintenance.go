package dto

import "time"

// =====================================================
// Maintenance Plan DTOs
// =====================================================

// CreateMaintenancePlanRequest creates a new maintenance plan
type CreateMaintenancePlanRequest struct {
	Name            string `json:"name" binding:"required"`
	EquipmentTypeID uint   `json:"equipment_type_id" binding:"required"`
	Level           int    `json:"level" binding:"required,min=1,max=3"`           // 1=一级,2=二级,3=精度
	CycleDays       int    `json:"cycle_days" binding:"required,min=1"`          // 周期天数
	FlexibleDays    int    `json:"flexible_days"`                                // 弹性窗口
	WorkHours       float64 `json:"work_hours"`                                  // 工时定额
}

// MaintenancePlanResponse represents a maintenance plan in API responses
type MaintenancePlanResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	EquipmentTypeID uint   `json:"equipment_type_id"`
	EquipmentTypeName string `json:"equipment_type_name,omitempty"`
	Level           int    `json:"level"`
	LevelName       string `json:"level_name,omitempty"`
	CycleDays       int    `json:"cycle_days"`
	FlexibleDays    int    `json:"flexible_days"`
	WorkHours       float64 `json:"work_hours"`
	ItemCount       int    `json:"item_count"`
	CreatedAt       time.Time `json:"created_at"`
}

// =====================================================
// Maintenance Plan Item DTOs
// =====================================================

// CreateMaintenanceItemRequest creates a new maintenance item
type CreateMaintenanceItemRequest struct {
	PlanID        uint   `json:"plan_id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Method        string `json:"method"`
	Criteria      string `json:"criteria"`
	SequenceOrder int    `json:"sequence_order" binding:"required,min=1"`
}

// MaintenanceItemResponse represents a maintenance item in API responses
type MaintenanceItemResponse struct {
	ID            uint   `json:"id"`
	PlanID        uint   `json:"plan_id"`
	Name          string `json:"name"`
	Method        string `json:"method,omitempty"`
	Criteria      string `json:"criteria,omitempty"`
	SequenceOrder int    `json:"sequence_order"`
}

// =====================================================
// Maintenance Task DTOs
// =====================================================

// MaintenanceTaskStatus represents the status of a maintenance task
type MaintenanceTaskStatus string

const (
	MaintenanceStatusPending    MaintenanceTaskStatus = "pending"     // 待执行
	MaintenanceStatusInProgress MaintenanceTaskStatus = "in_progress" // 进行中
	MaintenanceStatusCompleted  MaintenanceTaskStatus = "completed"   // 已完成
	MaintenanceStatusOverdue    MaintenanceTaskStatus = "overdue"     // 已逾期
)

// GenerateMaintenanceTasksRequest generates maintenance tasks
type GenerateMaintenanceTasksRequest struct {
	PlanID       uint   `json:"plan_id" binding:"required"`
	EquipmentIDs []uint `json:"equipment_ids" binding:"required"`
	Date         string `json:"date" binding:"required"` // YYYY-MM-DD
}

// MaintenanceTaskResponse represents a maintenance task in API responses
type MaintenanceTaskResponse struct {
	ID             uint                     `json:"id"`
	PlanID         uint                     `json:"plan_id"`
	PlanName       string                   `json:"plan_name,omitempty"`
	EquipmentID    uint                     `json:"equipment_id"`
	EquipmentCode  string                   `json:"equipment_code,omitempty"`
	EquipmentName  string                   `json:"equipment_name,omitempty"`
	AssignedTo     uint                     `json:"assigned_to,omitempty"`
	AssigneeName   string                   `json:"assignee_name,omitempty"`
	ScheduledDate  string                   `json:"scheduled_date"`
	DueDate        string                   `json:"due_date"`
	Status         MaintenanceTaskStatus    `json:"status"`
	StartedAt      *time.Time                `json:"started_at,omitempty"`
	CompletedAt    *time.Time                `json:"completed_at,omitempty"`
	ActualHours    float64                  `json:"actual_hours"`
	Remark         string                   `json:"remark,omitempty"`
	ItemCount      int                      `json:"item_count"`
	CompletedCount int                      `json:"completed_count"`
	CreatedAt      time.Time                `json:"created_at"`
}

// MaintenanceTaskListResponse represents a paginated list of maintenance tasks
type MaintenanceTaskListResponse struct {
	Total int64                       `json:"total"`
	Items []MaintenanceTaskResponse `json:"items"`
}

// =====================================================
// Maintenance Execution DTOs
// =====================================================

// StartMaintenanceRequest starts a maintenance task
type StartMaintenanceRequest struct {
	TaskID    uint     `json:"task_id" binding:"required"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

// MaintenanceItemRecord represents a completed maintenance item
type MaintenanceItemRecord struct {
	ItemID    uint   `json:"item_id" binding:"required"`
	Result    string `json:"result" binding:"required"` // OK/NG
	Remark    string `json:"remark"`
	PhotoURL  string `json:"photo_url"`
}

// CompleteMaintenanceRequest completes a maintenance task
type CompleteMaintenanceRequest struct {
	TaskID      uint                     `json:"task_id" binding:"required"`
	Records     []MaintenanceItemRecord   `json:"records" binding:"required"`
	Latitude    *float64                 `json:"latitude"`
	Longitude   *float64                 `json:"longitude"`
	ActualHours float64                  `json:"actual_hours"`
	Remark       string                   `json:"remark"`
}

// CompleteMaintenanceResponse represents the result of completing maintenance
type CompleteMaintenanceResponse struct {
	TaskID       uint   `json:"task_id"`
	CompletedAt  string `json:"completed_at"`
	TotalCount   int    `json:"total_count"`
	OKCount      int    `json:"ok_count"`
	NGCount      int    `json:"ng_count"`
	NGItemIDs   []uint `json:"ng_item_ids"`
}

// =====================================================
// Maintenance Statistics DTOs
// =====================================================

// MaintenanceStatistics represents maintenance statistics
type MaintenanceStatistics struct {
	TotalPlans     int64 `json:"total_plans"`
	TotalTasks     int64 `json:"total_tasks"`
	PendingTasks   int64 `json:"pending_tasks"`
	InProgressTasks int64 `json:"in_progress_tasks"`
	CompletedTasks int64 `json:"completed_tasks"`
	OverdueTasks   int64 `json:"overdue_tasks"`
	TodayCompleted int64 `json:"today_completed"`
	CompletionRate float64 `json:"completion_rate"`
}

// MyMaintenanceStatistics represents current user's maintenance statistics
type MyMaintenanceStatistics struct {
	PendingCount   int `json:"pending_count"`
	InProgressCount int `json:"in_progress_count"`
	TodayTasks     int `json:"today_tasks"`
	ThisWeekTasks  int `json:"this_week_tasks"`
}

// MaintenanceTaskQuery represents query parameters for maintenance tasks
type MaintenanceTaskQuery struct {
	Status     string `form:"status"`
	AssignedTo uint   `form:"assigned_to"`
	DateFrom   string `form:"date_from"`
	DateTo     string `form:"date_to"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
}
