package repository

import (
	"fmt"
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
func (r *AnalyticsRepository) GetEquipmentStats(factoryID *uint) (map[string]int64, error) {
	stats := make(map[string]int64)

	query := r.db.Model(&model.Equipment{})
	if factoryID != nil {
		query = query.Joins("INNER JOIN workshops w ON equipment.workshop_id = w.id").Where("w.factory_id = ?", *factoryID)
	}

	var total, running, stopped, maintenance, scrapped int64
	query.Count(&total)
	query.Where("status = ?", "running").Count(&running)
	query.Where("status = ?", "stopped").Count(&stopped)
	query.Where("status = ?", "maintenance").Count(&maintenance)
	query.Where("status = ?", "scrapped").Count(&scrapped)

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
	query := `
		SELECT COALESCE(SUM(GREATEST(EXTRACT(EPOCH FROM (completed_at - created_at))/3600, 0)), 0) as total_hours,
		       COUNT(*) as count
		FROM repair_orders ro
		INNER JOIN equipment e ON ro.equipment_id = e.id
		INNER JOIN workshops w ON e.workshop_id = w.id
		WHERE ro.status = 'closed' AND ro.completed_at IS NOT NULL
	`
	args := []interface{}{}
	if factoryID != nil {
		query += " AND w.factory_id = ?"
		args = append(args, *factoryID)
	}

	r.db.Raw(query, args...).Scan(&repairTime)

	if repairTime.Count > 0 {
		stats["mttr"] = repairTime.TotalHours / float64(repairTime.Count)
	} else {
		stats["mttr"] = 0
	}

	// MTBF: Mean Time Between Failures
	var mtbfData struct {
		TotalDays int64
		Failures  int64
	}

	mtbfQuery := `
		SELECT COALESCE(EXTRACT(DAY FROM (NOW() - MIN(ro.created_at)))::integer, 30) as total_days,
		       COUNT(*) as failures
		FROM repair_orders ro
		INNER JOIN equipment e ON ro.equipment_id = e.id
		INNER JOIN workshops w ON e.workshop_id = w.id
		WHERE ro.status = 'closed' AND ro.created_at > NOW() - INTERVAL '1 year'
	`
	if factoryID != nil {
		mtbfQuery += " AND w.factory_id = ?"
	}

	r.db.Raw(mtbfQuery, args...).Scan(&mtbfData)

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

// GetMTBFRanking returns equipment ranking by MTBF
func (r *AnalyticsRepository) GetMTBFRanking(limit int, factoryID *uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			e.id as equipment_id,
			e.code as equipment_code,
			e.name as equipment_name,
			COALESCE(EXTRACT(DAY FROM (NOW() - MIN(ro.created_at)))::integer, 30) * 24 / GREATEST(COUNT(ro.id), 1) as mtbf
		FROM equipment e
		LEFT JOIN repair_orders ro ON ro.equipment_id = e.id AND ro.status = 'closed'
		INNER JOIN workshops w ON e.workshop_id = w.id
		WHERE 1=1
	`
	args := []interface{}{}
	if factoryID != nil {
		query += " AND w.factory_id = ?"
		args = append(args, *factoryID)
	}
	query += ` GROUP BY e.id, e.code, e.name ORDER BY mtbf DESC LIMIT ?`
	args = append(args, limit)

	err := r.db.Raw(query, args...).Scan(&results).Error
	return results, err
}

// GetDowntimeRanking returns equipment ranking by total downtime
func (r *AnalyticsRepository) GetDowntimeRanking(limit int, factoryID *uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			e.id as equipment_id,
			e.code as equipment_code,
			e.name as equipment_name,
			COALESCE(SUM(GREATEST(EXTRACT(EPOCH FROM (ro.completed_at - ro.created_at))/3600, 0)), 0) as total_downtime
		FROM equipment e
		INNER JOIN repair_orders ro ON ro.equipment_id = e.id AND ro.status = 'closed'
		INNER JOIN workshops w ON e.workshop_id = w.id
		WHERE 1=1
	`
	args := []interface{}{}
	if factoryID != nil {
		query += " AND w.factory_id = ?"
		args = append(args, *factoryID)
	}
	query += ` GROUP BY e.id, e.code, e.name ORDER BY total_downtime DESC LIMIT ?`
	args = append(args, limit)

	err := r.db.Raw(query, args...).Scan(&results).Error
	return results, err
}

// GetCompletionRates returns completion rates for different task types
func (r *AnalyticsRepository) GetCompletionRates(factoryID *uint) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Inspection completion rate
	insQuery := r.db.Model(&model.InspectionTask{})
	if factoryID != nil {
		insQuery = insQuery.Joins("INNER JOIN equipment e ON inspection_tasks.equipment_id = e.id").
			Joins("INNER JOIN workshops w ON e.workshop_id = w.id").Where("w.factory_id = ?", *factoryID)
	}
	var inspectionTotal, inspectionCompleted int64
	insQuery.Count(&inspectionTotal)
	insQuery.Where("inspection_tasks.status = ?", "completed").Count(&inspectionCompleted)
	stats["inspection_total"] = inspectionTotal
	stats["inspection_completed"] = inspectionCompleted

	// Maintenance completion rate
	maintQuery := r.db.Model(&model.MaintenanceTask{})
	if factoryID != nil {
		maintQuery = maintQuery.Joins("INNER JOIN equipment e ON maintenance_tasks.equipment_id = e.id").
			Joins("INNER JOIN workshops w ON e.workshop_id = w.id").Where("w.factory_id = ?", *factoryID)
	}
	var maintenanceTotal, maintenanceCompleted int64
	maintQuery.Count(&maintenanceTotal)
	maintQuery.Where("maintenance_tasks.status = ?", "completed").Count(&maintenanceCompleted)
	stats["maintenance_total"] = maintenanceTotal
	stats["maintenance_completed"] = maintenanceCompleted

	// Repair completion rate
	repQuery := r.db.Model(&model.RepairOrder{})
	if factoryID != nil {
		repQuery = repQuery.Joins("INNER JOIN equipment e ON repair_orders.equipment_id = e.id").
			Joins("INNER JOIN workshops w ON e.workshop_id = w.id").Where("w.factory_id = ?", *factoryID)
	}
	var repairTotal, repairCompleted int64
	repQuery.Count(&repairTotal)
	repQuery.Where("repair_orders.status = ?", "closed").Count(&repairCompleted)
	stats["repair_total"] = repairTotal
	stats["repair_completed"] = repairCompleted

	return stats, nil
}

// GetPendingTasks returns count of pending tasks
func (r *AnalyticsRepository) GetPendingTasks(factoryID *uint) (map[string]int64, error) {
	stats := make(map[string]int64)

	today := time.Now().Format("2006-01-02")
	weekEnd := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	insQuery := r.db.Model(&model.InspectionTask{})
	maintQuery := r.db.Model(&model.MaintenanceTask{})
	repQuery := r.db.Model(&model.RepairOrder{})

	if factoryID != nil {
		insQuery = insQuery.Joins("INNER JOIN equipment e ON inspection_tasks.equipment_id = e.id").
			Joins("INNER JOIN workshops w ON e.workshop_id = w.id").Where("w.factory_id = ?", *factoryID)
		maintQuery = maintQuery.Joins("INNER JOIN equipment e ON maintenance_tasks.equipment_id = e.id").
			Joins("INNER JOIN workshops w ON e.workshop_id = w.id").Where("w.factory_id = ?", *factoryID)
		repQuery = repQuery.Joins("INNER JOIN equipment e ON repair_orders.equipment_id = e.id").
			Joins("INNER JOIN workshops w ON e.workshop_id = w.id").Where("w.factory_id = ?", *factoryID)
	}

	var pendingInspections, pendingMaintenances, pendingRepairs int64
	insQuery.Where("inspection_tasks.scheduled_date >= ? AND inspection_tasks.scheduled_date < ? AND inspection_tasks.status IN ?",
			today, weekEnd, []string{"pending", "in_progress"}).
		Count(&pendingInspections)

	maintQuery.Where("maintenance_tasks.scheduled_date >= ? AND maintenance_tasks.scheduled_date < ? AND maintenance_tasks.status IN ?",
			today, weekEnd, []string{"pending", "in_progress"}).
		Count(&pendingMaintenances)

	repQuery.Where("repair_orders.status IN ?", []string{"pending", "assigned"}).Count(&pendingRepairs)

	stats["pending_inspections"] = pendingInspections
	stats["pending_maintenances"] = pendingMaintenances
	stats["pending_repairs"] = pendingRepairs

	return stats, nil
}

// GetTrendData returns daily trend data
func (r *AnalyticsRepository) GetTrendData(startDate, endDate string, factoryID *uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			d.date::text as date,
			COALESCE(i.inspections, 0) as inspection_tasks,
			COALESCE(m.maintenances, 0) as maintenance_tasks,
			COALESCE(rep.repairs, 0) as repair_orders
		FROM (
			SELECT generate_series(
				?::date,
				?::date,
				'1 day'::interval
			) as date
		) d
		LEFT JOIN (
			SELECT DATE(it.completed_at) as date, COUNT(*) as inspections
			FROM inspection_tasks it
			INNER JOIN equipment e ON it.equipment_id = e.id
			INNER JOIN workshops w ON e.workshop_id = w.id
			WHERE it.completed_at >= ?::date AND it.completed_at < ?::date + interval '1 day'
			%s
			GROUP BY DATE(it.completed_at)
		) i ON d.date = i.date
		LEFT JOIN (
			SELECT DATE(mt.completed_at) as date, COUNT(*) as maintenances
			FROM maintenance_tasks mt
			INNER JOIN equipment e ON mt.equipment_id = e.id
			INNER JOIN workshops w ON e.workshop_id = w.id
			WHERE mt.completed_at >= ?::date AND mt.completed_at < ?::date + interval '1 day'
			%s
			GROUP BY DATE(mt.completed_at)
		) m ON d.date = m.date
		LEFT JOIN (
			SELECT DATE(ro.completed_at) as date, COUNT(*) as repairs
			FROM repair_orders ro
			INNER JOIN equipment e ON ro.equipment_id = e.id
			INNER JOIN workshops w ON e.workshop_id = w.id
			WHERE ro.completed_at >= ?::date AND ro.completed_at < ?::date + interval '1 day'
			%s
			GROUP BY DATE(ro.completed_at)
		) rep ON d.date = rep.date
		ORDER BY d.date
	`

	factoryFilter := ""
	args := []interface{}{startDate, endDate, startDate, endDate, startDate, endDate, startDate, endDate}
	if factoryID != nil {
		factoryFilter = " AND w.factory_id = ?"
		// We need to insert the factoryID for each subquery
		// Since we have 3 subqueries with 2 date placeholders each, we need to be careful with the order.
		// Subquery i: DATE(it.completed_at) >= ? AND < ? AND w.factory_id = ?
		// Subquery m: DATE(mt.completed_at) >= ? AND < ? AND w.factory_id = ?
		// Subquery rep: DATE(ro.completed_at) >= ? AND < ? AND w.factory_id = ?
		
		// Wait, the query string construction with %s is better.
		query = fmt.Sprintf(query, factoryFilter, factoryFilter, factoryFilter)
		
		// New args list:
		args = []interface{}{
			startDate, endDate, // generate_series
			startDate, endDate, *factoryID, // i
			startDate, endDate, *factoryID, // m
			startDate, endDate, *factoryID, // rep
		}
	} else {
		query = fmt.Sprintf(query, "", "", "")
	}

	err := r.db.Raw(query, args...).Scan(&results).Error

	return results, err
}

// GetEquipmentPerformanceRanking returns equipment ranking by performance (availability)
func (r *AnalyticsRepository) GetEquipmentPerformanceRanking(limit int, factoryID *uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT
			e.id as equipment_id,
			e.code as equipment_code,
			e.name as equipment_name,
			COALESCE(AVG(CASE WHEN ro.id IS NULL THEN 100 ELSE 
				(720 - GREATEST(EXTRACT(EPOCH FROM (ro.completed_at - ro.created_at))/3600, 0)) / 720 * 100 
			END), 100) as performance_score
		FROM equipment e
		LEFT JOIN repair_orders ro ON ro.equipment_id = e.id 
			AND ro.status = 'closed'
			AND ro.created_at > NOW() - INTERVAL '30 days'
		INNER JOIN workshops w ON e.workshop_id = w.id
		WHERE 1=1
	`
	args := []interface{}{}
	if factoryID != nil {
		query += " AND w.factory_id = ?"
		args = append(args, *factoryID)
	}
	query += ` GROUP BY e.id, e.code, e.name ORDER BY performance_score DESC LIMIT ?`
	args = append(args, limit)

	err := r.db.Raw(query, args...).Scan(&results).Error
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
				GREATEST(EXTRACT(EPOCH FROM (COALESCE(ro.completed_at, NOW()) - ro.created_at))/3600, 0)
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
				GREATEST(EXTRACT(EPOCH FROM (COALESCE(ro.completed_at, NOW()) - ro.created_at))/3600, 0)
			), 0) as downtime_hours,
			COALESCE(AVG(
				GREATEST(EXTRACT(EPOCH FROM (COALESCE(ro.completed_at, NOW()) - ro.created_at))/3600, 0)
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
