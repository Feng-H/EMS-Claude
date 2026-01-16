package repository

import (
	"time"

	"github.com/ems/backend/internal/model"
	"gorm.io/gorm"
)

// AnalyticsRepository provides analytics queries
type AnalyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository() *AnalyticsRepository {
	return &AnalyticsRepository{db: DB}
}

// GetEquipmentStats returns equipment status statistics
func (r *AnalyticsRepository) GetEquipmentStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	var total, running, stopped, maintenance, scrapped int64
	r.db.Model(&model.Equipment{}).Count(&total)
	r.db.Model(&model.Equipment{}).Where("status = ?", "running").Count(&running)
	r.db.Model(&model.Equipment{}).Where("status = ?", "stopped").Count(&stopped)
	r.db.Model(&model.Equipment{}).Where("status = ?", "maintenance").Count(&maintenance)
	r.db.Model(&model.Equipment{}).Where("status = ?", "scrapped").Count(&scrapped)

	stats["total"] = total
	stats["running"] = running
	stats["stopped"] = stopped
	stats["maintenance"] = maintenance
	stats["scrapped"] = scrapped

	return stats, nil
}

// GetMTTRMTBF calculates MTTR and MTBF
func (r *AnalyticsRepository) GetMTTRMTBF(factoryID *uint) (map[string]float64, error) {
	stats := make(map[string]float64)

	// MTTR: Mean Time To Repair = Total repair time / Number of completed repairs
	type RepairTimeResult struct {
		TotalHours float64
		Count      int64
	}

	var repairTime RepairTimeResult
	r.db.Raw(`
		SELECT COALESCE(SUM(EXTRACT(EPOCH FROM (completed_at - created_at))/3600), 0) as total_hours,
		       COUNT(*) as count
		FROM repair_orders
		WHERE status = 'closed' AND completed_at IS NOT NULL
	`).Scan(&repairTime)

	if repairTime.Count > 0 {
		stats["mttr"] = repairTime.TotalHours / float64(repairTime.Count)
	} else {
		stats["mttr"] = 0
	}

	// MTBF: Mean Time Between Failures = Total running time / Number of failures
	// Simplified: Days between first and last repair / (number of repairs - 1)
	// For now, use inverse of failure rate per day
	var mtbfData struct {
		TotalDays int64
		Failures  int64
	}

	r.db.Raw(`
		SELECT COALESCE(EXTRACT(DAY FROM (NOW() - MIN(created_at)))::integer, 30) as total_days,
		       COUNT(*) as failures
		FROM repair_orders
		WHERE status = 'closed' AND created_at > NOW() - INTERVAL '1 year'
	`).Scan(&mtbfData)

	if mtbfData.Failures > 1 && mtbfData.TotalDays > 0 {
		stats["mtbf"] = float64(mtbfData.TotalDays) / float64(mtbfData.Failures) * 24 // hours
	} else {
		stats["mtbf"] = 720 // Default: 30 days = 720 hours
	}

	// Availability = MTBF / (MTTR + MTBF)
	mttr := stats["mttr"]
	mtbf := stats["mtbf"]
	if mttr+mtbf > 0 {
		stats["availability"] = mtbf / (mttr + mtbf) * 100
	} else {
		stats["availability"] = 0
	}

	return stats, nil
}

// GetCompletionRates returns completion rates for different task types
func (r *AnalyticsRepository) GetCompletionRates() (map[string]int64, error) {
	stats := make(map[string]int64)

	// Inspection completion rate
	var inspectionTotal, inspectionCompleted int64
	r.db.Model(&model.InspectionTask{}).Count(&inspectionTotal)
	r.db.Model(&model.InspectionTask{}).Where("status = ?", "completed").Count(&inspectionCompleted)
	stats["inspection_total"] = inspectionTotal
	stats["inspection_completed"] = inspectionCompleted

	// Maintenance completion rate
	var maintenanceTotal, maintenanceCompleted int64
	r.db.Model(&model.MaintenanceTask{}).Count(&maintenanceTotal)
	r.db.Model(&model.MaintenanceTask{}).Where("status = ?", "completed").Count(&maintenanceCompleted)
	stats["maintenance_total"] = maintenanceTotal
	stats["maintenance_completed"] = maintenanceCompleted

	// Repair completion rate
	var repairTotal, repairCompleted int64
	r.db.Model(&model.RepairOrder{}).Count(&repairTotal)
	r.db.Model(&model.RepairOrder{}).Where("status = ?", "closed").Count(&repairCompleted)
	stats["repair_total"] = repairTotal
	stats["repair_completed"] = repairCompleted

	return stats, nil
}

// GetPendingTasks returns count of pending tasks
func (r *AnalyticsRepository) GetPendingTasks() (map[string]int64, error) {
	stats := make(map[string]int64)

	today := time.Now().Format("2006-01-02")
	weekEnd := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	var pendingInspections, pendingMaintenances, pendingRepairs int64
	r.db.Model(&model.InspectionTask{}).
		Where("scheduled_date >= ? AND scheduled_date < ? AND status IN ?",
			today, weekEnd, []string{"pending", "in_progress"}).
		Count(&pendingInspections)

	r.db.Model(&model.MaintenanceTask{}).
		Where("scheduled_date >= ? AND scheduled_date < ? AND status IN ?",
			today, weekEnd, []string{"pending", "in_progress"}).
		Count(&pendingMaintenances)

	r.db.Model(&model.RepairOrder{}).Where("status IN ?", []string{"pending", "assigned"}).Count(&pendingRepairs)

	stats["pending_inspections"] = pendingInspections
	stats["pending_maintenances"] = pendingMaintenances
	stats["pending_repairs"] = pendingRepairs

	return stats, nil
}

// GetTrendData returns daily trend data
func (r *AnalyticsRepository) GetTrendData(startDate, endDate string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			date::text as date,
			COALESCE(inspections, 0) as inspection_tasks,
			COALESCE(maintenances, 0) as maintenance_tasks,
			COALESCE(repairs, 0) as repair_orders
		FROM (
			SELECT generate_series(
				?::date,
				?::date,
				'1 day'::interval
			) as date
		) d
		LEFT JOIN (
			SELECT DATE(completed_at) as date, COUNT(*) as inspections
			FROM inspection_tasks
			WHERE completed_at >= ?::date AND completed_at < ?::date + interval '1 day'
			GROUP BY DATE(completed_at)
		) i ON d.date = i.date
		LEFT JOIN (
			SELECT DATE(completed_at) as date, COUNT(*) as maintenances
			FROM maintenance_tasks
			WHERE completed_at >= ?::date AND completed_at < ?::date + interval '1 day'
			GROUP BY DATE(completed_at)
		) m ON d.date = m.date
		LEFT JOIN (
			SELECT DATE(completed_at) as date, COUNT(*) as repairs
			FROM repair_orders
			WHERE completed_at >= ?::date AND completed_at < ?::date + interval '1 day'
			GROUP BY DATE(completed_at)
		) r ON d.date = r.date
		ORDER BY date
	`

	err := r.db.Raw(query, startDate, endDate, startDate, endDate, startDate, endDate, startDate, endDate).
		Scan(&results).Error

	return results, err
}

// GetFailureAnalysis returns failure analysis by equipment type
func (r *AnalyticsRepository) GetFailureAnalysis(limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			et.id as equipment_type_id,
			et.name as equipment_type_name,
			COUNT(ro.id) as failure_count,
			COALESCE(SUM(
				EXTRACT(EPOCH FROM (COALESCE(ro.completed_at, NOW()) - ro.created_at))/3600
			), 0) as total_downtime
		FROM equipment_types et
		LEFT JOIN equipment e ON e.type_id = et.id
		LEFT JOIN repair_orders ro ON ro.equipment_id = e.id
			AND ro.status = 'closed'
			AND ro.created_at > NOW() - INTERVAL '6 months'
		GROUP BY et.id, et.name
		HAVING COUNT(ro.id) > 0
		ORDER BY failure_count DESC
		LIMIT ?
	`

	err := r.db.Raw(query, limit).Scan(&results).Error
	return results, err
}

// GetTopFailureEquipment returns equipment with most failures
func (r *AnalyticsRepository) GetTopFailureEquipment(limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			e.id as equipment_id,
			e.code as equipment_code,
			e.name as equipment_name,
			COUNT(ro.id) as failure_count,
			COALESCE(SUM(
				EXTRACT(EPOCH FROM (COALESCE(ro.completed_at, NOW()) - ro.created_at))/3600
			), 0) as downtime_hours,
			COALESCE(AVG(
				EXTRACT(EPOCH FROM (COALESCE(ro.completed_at, NOW()) - ro.created_at))/3600
			), 0) as mttr
		FROM equipment e
		INNER JOIN repair_orders ro ON ro.equipment_id = e.id
			AND ro.status = 'closed'
			AND ro.created_at > NOW() - INTERVAL '6 months'
		GROUP BY e.id, e.code, e.name
		ORDER BY failure_count DESC
		LIMIT ?
	`

	err := r.db.Raw(query, limit).Scan(&results).Error
	return results, err
}

// GetLowStockCount returns count of low stock items
func (r *AnalyticsRepository) GetLowStockCount() (int64, error) {
	var count int64

	r.db.Raw(`
		SELECT COUNT(*)
		FROM spare_part_inventories spi
		INNER JOIN spare_parts sp ON spi.spare_part_id = sp.id
		WHERE spi.quantity < COALESCE(sp.safety_stock, 0)
	`).Scan(&count)

	return count, nil
}
