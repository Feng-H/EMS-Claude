package service

import (
	"github.com/ems/backend/internal/repository"
)

// AnalyticsService
type AnalyticsService struct {
	analyticsRepo *repository.AnalyticsRepository
}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: repository.NewAnalyticsRepository(),
	}
}

// GetDashboardOverview returns the dashboard overview data
func (s *AnalyticsService) GetDashboardOverview() (*DashboardOverview, error) {
	// Get equipment stats
	equipmentStats, err := s.analyticsRepo.GetEquipmentStats()
	if err != nil {
		return nil, err
	}

	// Get MTTR/MTBF
	mttrMtbf, err := s.analyticsRepo.GetMTTRMTBF(nil)
	if err != nil {
		return nil, err
	}

	// Get completion rates
	completionRates, err := s.analyticsRepo.GetCompletionRates()
	if err != nil {
		return nil, err
	}

	// Get pending tasks
	pendingTasks, err := s.analyticsRepo.GetPendingTasks()
	if err != nil {
		return nil, err
	}

	// Get low stock count
	lowStockCount, err := s.analyticsRepo.GetLowStockCount()
	if err != nil {
		return nil, err
	}

	// Calculate completion rates
	var inspectionRate, maintenanceRate, repairRate float64
	if completionRates["inspection_total"] > 0 {
		inspectionRate = float64(completionRates["inspection_completed"]) / float64(completionRates["inspection_total"]) * 100
	}
	if completionRates["maintenance_total"] > 0 {
		maintenanceRate = float64(completionRates["maintenance_completed"]) / float64(completionRates["maintenance_total"]) * 100
	}
	if completionRates["repair_total"] > 0 {
		repairRate = float64(completionRates["repair_completed"]) / float64(completionRates["repair_total"]) * 100
	}

	return &DashboardOverview{
		Equipment: EquipmentAnalytics{
			TotalEquipment:     equipmentStats["total"],
			RunningEquipment:   equipmentStats["running"],
			StoppedEquipment:   equipmentStats["stopped"],
			MaintenanceEquipment: equipmentStats["maintenance"],
			ScrappedEquipment:  equipmentStats["scrapped"],
		},
		MTTR_MTBF: MTTRMTBF{
			MTTR:        mttrMtbf["mttr"],
			MTBF:        mttrMtbf["mtbf"],
			Availability: mttrMtbf["availability"],
		},
		Tasks: CompletionRate{
			InspectionCompletionRate:  inspectionRate,
			MaintenanceCompletionRate: maintenanceRate,
			RepairCompletionRate:     repairRate,
		},
		PendingInspections:  pendingTasks["pending_inspections"],
		PendingMaintenances: pendingTasks["pending_maintenances"],
		PendingRepairs:      pendingTasks["pending_repairs"],
		LowStockAlerts:      lowStockCount,
	}, nil
}

// GetMTTRMTBF returns MTTR and MTBF data
func (s *AnalyticsService) GetMTTRMTBF(factoryID *uint) (*MTTRMTBF, error) {
	stats, err := s.analyticsRepo.GetMTTRMTBF(factoryID)
	if err != nil {
		return nil, err
	}

	return &MTTRMTBF{
		MTTR:        stats["mttr"],
		MTBF:        stats["mtbf"],
		Availability: stats["availability"],
	}, nil
}

// GetTrendData returns trend data for the specified date range
func (s *AnalyticsService) GetTrendData(startDate, endDate string) ([]TrendData, error) {
	results, err := s.analyticsRepo.GetTrendData(startDate, endDate)
	if err != nil {
		return nil, err
	}

	trends := make([]TrendData, len(results))
	for i, r := range results {
		trends[i] = TrendData{
			Date:            r["date"].(string),
			InspectionTasks:  int(r["inspection_tasks"].(int64)),
			MaintenanceTasks: int(r["maintenance_tasks"].(int64)),
			RepairOrders:     int(r["repair_orders"].(int64)),
			DowntimeHours:    0, // Can be calculated from repair orders
		}
	}

	return trends, nil
}

// GetFailureAnalysis returns failure analysis by equipment type
func (s *AnalyticsService) GetFailureAnalysis(limit int) ([]FailureAnalysis, error) {
	results, err := s.analyticsRepo.GetFailureAnalysis(limit)
	if err != nil {
		return nil, err
	}

	analysis := make([]FailureAnalysis, len(results))
	for i, r := range results {
		analysis[i] = FailureAnalysis{
			EquipmentTypeID:   uint(r["equipment_type_id"].(int64)),
			EquipmentTypeName: r["equipment_type_name"].(string),
			FailureCount:      r["failure_count"].(int64),
			TotalDowntime:     r["total_downtime"].(int64),
		}
	}

	return analysis, nil
}

// GetTopFailureEquipment returns equipment with most failures
func (s *AnalyticsService) GetTopFailureEquipment(limit int) ([]TopFailureEquipment, error) {
	results, err := s.analyticsRepo.GetTopFailureEquipment(limit)
	if err != nil {
		return nil, err
	}

	equipment := make([]TopFailureEquipment, len(results))
	for i, r := range results {
		equipment[i] = TopFailureEquipment{
			EquipmentID:   uint(r["equipment_id"].(int64)),
			EquipmentCode: r["equipment_code"].(string),
			EquipmentName: r["equipment_name"].(string),
			FailureCount:  r["failure_count"].(int64),
			DowntimeHours: r["downtime_hours"].(int64),
			MTTR:          r["mttr"].(float64),
		}
	}

	return equipment, nil
}

// Types
type EquipmentAnalytics struct {
	TotalEquipment     int64 `json:"total_equipment"`
	RunningEquipment   int64 `json:"running_equipment"`
	StoppedEquipment   int64 `json:"stopped_equipment"`
	MaintenanceEquipment int64 `json:"maintenance_equipment"`
	ScrappedEquipment  int64 `json:"scrapped_equipment"`
}

type MTTRMTBF struct {
	MTTR         float64 `json:"mttr"`
	MTBF         float64 `json:"mtbf"`
	Availability float64 `json:"availability"`
}

type CompletionRate struct {
	InspectionCompletionRate  float64 `json:"inspection_completion_rate"`
	MaintenanceCompletionRate float64 `json:"maintenance_completion_rate"`
	RepairCompletionRate      float64 `json:"repair_completion_rate"`
}

type DashboardOverview struct {
	Equipment          EquipmentAnalytics `json:"equipment"`
	MTTR_MTBF          MTTRMTBF           `json:"mttr_mtbf"`
	Tasks              CompletionRate     `json:"tasks"`
	PendingInspections int64              `json:"pending_inspections"`
	PendingMaintenances int64             `json:"pending_maintenances"`
	PendingRepairs     int64              `json:"pending_repairs"`
	LowStockAlerts     int64              `json:"low_stock_alerts"`
}

type TrendData struct {
	Date            string `json:"date"`
	InspectionTasks int    `json:"inspection_tasks"`
	MaintenanceTasks int    `json:"maintenance_tasks"`
	RepairOrders     int    `json:"repair_orders"`
	DowntimeHours    int    `json:"downtime_hours"`
}

type FailureAnalysis struct {
	EquipmentTypeID   uint   `json:"equipment_type_id"`
	EquipmentTypeName string `json:"equipment_type_name"`
	FailureCount      int64  `json:"failure_count"`
	TotalDowntime     int64  `json:"total_downtime"`
}

type TopFailureEquipment struct {
	EquipmentID   uint    `json:"equipment_id"`
	EquipmentCode string  `json:"equipment_code"`
	EquipmentName string  `json:"equipment_name"`
	FailureCount  int64   `json:"failure_count"`
	DowntimeHours int64   `json:"downtime_hours"`
	MTTR          float64 `json:"mttr"`
}
