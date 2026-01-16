package dto

import "time"

// Inspection Template DTOs

type InspectionTemplateRequest struct {
	Name            string `json:"name" binding:"required"`
	EquipmentTypeID uint   `json:"equipment_type_id" binding:"required"`
}

type InspectionTemplateResponse struct {
	ID              uint             `json:"id"`
	Name            string           `json:"name"`
	EquipmentTypeID uint             `json:"equipment_type_id"`
	EquipmentTypeName string         `json:"equipment_type_name,omitempty"`
	ItemCount       int              `json:"item_count"`
	CreatedAt       time.Time        `json:"created_at"`
}

type InspectionTemplateDetailResponse struct {
	ID              uint                      `json:"id"`
	Name            string                    `json:"name"`
	EquipmentTypeID uint                      `json:"equipment_type_id"`
	EquipmentTypeName string                 `json:"equipment_type_name,omitempty"`
	Items           []InspectionItemResponse  `json:"items"`
	CreatedAt       time.Time                 `json:"created_at"`
}

// Inspection Item DTOs

type InspectionItemRequest struct {
	TemplateID uint   `json:"template_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Method     string `json:"method"`
	Criteria   string `json:"criteria"`
	SequenceOrder int  `json:"sequence_order"`
}

type InspectionItemResponse struct {
	ID            uint   `json:"id"`
	TemplateID    uint   `json:"template_id"`
	TemplateName  string `json:"template_name,omitempty"`
	Name          string `json:"name"`
	Method        string `json:"method"`
	Criteria      string `json:"criteria"`
	SequenceOrder int    `json:"sequence_order"`
}

// Inspection Task DTOs

type InspectionTaskQuery struct {
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"page_size" binding:"min=1,max=100"`
	AssignedTo uint   `form:"assigned_to"`
	Status     string `form:"status"`
	DateFrom   string `form:"date_from"`
	DateTo     string `form:"date_to"`
}

type InspectionTaskResponse struct {
	ID            uint                 `json:"id"`
	EquipmentID   uint                 `json:"equipment_id"`
	EquipmentCode string              `json:"equipment_code,omitempty"`
	EquipmentName string              `json:"equipment_name,omitempty"`
	TemplateID    uint                 `json:"template_id"`
	TemplateName  string              `json:"template_name,omitempty"`
	AssignedTo    uint                 `json:"assigned_to"`
	AssigneeName  string              `json:"assignee_name,omitempty"`
	ScheduledDate time.Time          `json:"scheduled_date"`
	Status        string              `json:"status"`
	StartedAt     *time.Time          `json:"started_at"`
	CompletedAt   *time.Time          `json:"completed_at"`
	Latitude      *float64            `json:"latitude"`
	Longitude     *float64            `json:"longitude"`
	ItemCount     int                 `json:"item_count"`
	CompletedCount int                `json:"completed_count"`
	Records       []InspectionRecordResponse `json:"records,omitempty"`
}

type InspectionTaskListResponse struct {
	Total int64                     `json:"total"`
	Items []InspectionTaskResponse  `json:"items"`
}

// Inspection Record DTOs

type InspectionRecordRequest struct {
	TaskID  uint   `json:"task_id" binding:"required"`
	ItemID  uint   `json:"item_id" binding:"required"`
	Result  string `json:"result" binding:"required,oneof=OK NG"`
	Remark  string `json:"remark"`
	PhotoURL string `json:"photo_url"`
}

type InspectionRecordResponse struct {
	ID          uint      `json:"id"`
	TaskID      uint      `json:"task_id"`
	ItemID      uint      `json:"item_id"`
	ItemName    string    `json:"item_name,omitempty"`
	Result      string    `json:"result"`
	Remark      string    `json:"remark"`
	PhotoURL    string    `json:"photo_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// Start Inspection Request (with anti-fraud)

type StartInspectionRequest struct {
	EquipmentID uint   `json:"equipment_id" binding:"required"`
	QRCode      string `json:"qr_code" binding:"required"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timestamp   int64   `json:"timestamp" binding:"required"`
}

type StartInspectionResponse struct {
	TaskID      uint                      `json:"task_id"`
	EquipmentID uint                      `json:"equipment_id"`
	Equipment   EquipmentResponse         `json:"equipment"`
	Items       []InspectionItemResponse   `json:"items"`
	StartedAt   time.Time                 `json:"started_at"`
}

// Complete Inspection Request

type CompleteInspectionRequest struct {
	TaskID  uint                           `json:"task_id" binding:"required"`
	Records []InspectionRecordRequest       `json:"records" binding:"required"`
	Latitude *float64                      `json:"latitude"`
	Longitude *float64                     `json:"longitude"`
}

type CompleteInspectionResponse struct {
	TaskID         uint      `json:"task_id"`
	CompletedAt    time.Time `json:"completed_at"`
	TotalCount     int       `json:"total_count"`
	OKCount        int       `json:"ok_count"`
	NGCount        int       `json:"ng_count"`
	NGItems        []uint    `json:"ng_items"` // Record IDs with NG result
}

// Statistics

type InspectionStatistics struct {
	TotalTasks      int64   `json:"total_tasks"`
	PendingTasks    int64   `json:"pending_tasks"`
	InProgressTasks int64   `json:"in_progress_tasks"`
	CompletedTasks  int64   `json:"completed_tasks"`
	OverdueTasks    int64   `json:"overdue_tasks"`
	TodayCompleted  int64   `json:"today_completed"`
	CompletionRate  float64 `json:"completion_rate"` // Percentage
}

type MyTasksResponse struct {
	PendingCount    int `json:"pending_count"`
	InProgressCount int `json:"in_progress_count"`
	TodayTasks      int `json:"today_tasks"`
}

// Task Generation Request (for auto-generation)

type GenerateTasksRequest struct {
	EquipmentIDs []uint `json:"equipment_ids"`
	Date         string `json:"date" binding:"required"` // YYYY-MM-DD
}

type GenerateTasksResponse struct {
	CreatedCount int      `json:"created_count"`
	TaskIDs      []uint   `json:"task_ids"`
	Errors       []string `json:"errors"`
}
