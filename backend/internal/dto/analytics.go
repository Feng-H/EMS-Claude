package dto

// =====================================================
// Analytics DTOs
// =====================================================

// EquipmentAnalytics represents equipment statistics
type EquipmentAnalytics struct {
	TotalEquipment    int64   `json:"total_equipment"`
	RunningEquipment  int64   `json:"running_equipment"`
	StoppedEquipment  int64   `json:"stopped_equipment"`
	MaintenanceEquipment int64 `json:"maintenance_equipment"`
	ScrappedEquipment int64   `json:"scrapped_equipment"`
}

// MTTRMTBF represents Mean Time To Repair and Mean Time Between Failures
type MTTRMTBF struct {
	MTTR   float64 `json:"mttr"`   // Mean Time To Repair (hours)
	MTBF   float64 `json:"mtbf"`   // Mean Time Between Failures (hours)
	Availability float64 `json:"availability"` // MTBF / (MTTR + MTBF) * 100
}

// OEEData represents Overall Equipment Effectiveness data
type OEEData struct {
	EquipmentID   uint    `json:"equipment_id"`
	EquipmentCode string  `json:"equipment_code"`
	EquipmentName string  `json:"equipment_name"`
	Availability  float64 `json:"availability"`   // 可用率
	Performance   float64 `json:"performance"`     // 性能率
	Quality       float64 `json:"quality"`        // 质量率
	OEE           float64 `json:"oee"`            // OEE = Availability * Performance * Quality
}

// CompletionRate represents task completion rates
type CompletionRate struct {
	InspectionCompletionRate float64 `json:"inspection_completion_rate"`
	MaintenanceCompletionRate float64 `json:"maintenance_completion_rate"`
	RepairCompletionRate     float64 `json:"repair_completion_rate"`
}

// TrendData represents trend data over time
type TrendData struct {
	Date     string `json:"date"`
	InspectionTasks int `json:"inspection_tasks"`
	MaintenanceTasks int `json:"maintenance_tasks"`
	RepairOrders     int `json:"repair_orders"`
	DowntimeHours    int `json:"downtime_hours"`
}

// FailureAnalysis represents failure analysis data
type FailureAnalysis struct {
	EquipmentTypeID   uint   `json:"equipment_type_id"`
	EquipmentTypeName string `json:"equipment_type_name"`
	FailureCount      int64  `json:"failure_count"`
	TotalDowntime     int64  `json:"total_downtime"` // hours
}

// TopFailureEquipment represents equipment with most failures
type TopFailureEquipment struct {
	EquipmentID    uint   `json:"equipment_id"`
	EquipmentCode  string `json:"equipment_code"`
	EquipmentName  string `json:"equipment_name"`
	FailureCount   int64  `json:"failure_count"`
	DowntimeHours  int64  `json:"downtime_hours"`
	MTTR           float64 `json:"mttr"`
}

// DashboardOverview represents the dashboard overview
type DashboardOverview struct {
	Equipment        EquipmentAnalytics `json:"equipment"`
	Tasks            CompletionRate     `json:"tasks"`
	MTTR_MTBF        MTTRMTBF           `json:"mttr_mtbf"`
	PendingInspections int64            `json:"pending_inspections"`
	PendingMaintenances int64           `json:"pending_maintenances"`
	PendingRepairs   int64              `json:"pending_repairs"`
	LowStockAlerts   int64              `json:"low_stock_alerts"`
}

// AnalyticsQuery represents query parameters for analytics
type AnalyticsQuery struct {
	StartDate string `form:"start_date"` // YYYY-MM-DD
	EndDate   string `form:"end_date"`   // YYYY-MM-DD
	FactoryID *uint  `form:"factory_id"`
	EquipmentTypeID *uint `form:"equipment_type_id"`
}
