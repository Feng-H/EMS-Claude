package repository

import (
	"time"

	"github.com/ems/backend/internal/model"
	"gorm.io/gorm"
)

// RepairOrder Repository
type RepairOrderRepository struct {
	db *gorm.DB
}

func NewRepairOrderRepository() *RepairOrderRepository {
	return &RepairOrderRepository{db: DB}
}

func (r *RepairOrderRepository) Create(order *model.RepairOrder) error {
	return r.db.Create(order).Error
}

func (r *RepairOrderRepository) GetByID(id uint) (*model.RepairOrder, error) {
	var order model.RepairOrder
	err := r.db.Preload("Equipment").Preload("Equipment.Type").
		Preload("Reporter").Preload("Assignee").
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

type RepairOrderFilter struct {
	Status     string
	Priority   int
	AssignedTo uint
	DateFrom   time.Time
	DateTo     time.Time
	Page       int
	PageSize   int
}

func (r *RepairOrderRepository) List(filter RepairOrderFilter) ([]model.RepairOrder, int64, error) {
	var orders []model.RepairOrder
	var total int64

	query := r.db.Model(&model.RepairOrder{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority > 0 {
		query = query.Where("priority = ?", filter.Priority)
	}
	if filter.AssignedTo > 0 {
		query = query.Where("assigned_to = ?", filter.AssignedTo)
	}
	if !filter.DateFrom.IsZero() {
		query = query.Where("created_at >= ?", filter.DateFrom)
	}
	if !filter.DateTo.IsZero() {
		query = query.Where("created_at <= ?", filter.DateTo)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (filter.Page - 1) * filter.PageSize
	err := query.Preload("Equipment").Preload("Equipment.Type").
		Preload("Reporter").Preload("Assignee").
		Order("priority ASC, created_at DESC").
		Offset(offset).Limit(filter.PageSize).Find(&orders).Error

	return orders, total, err
}

func (r *RepairOrderRepository) Update(order *model.RepairOrder) error {
	return r.db.Save(order).Error
}

func (r *RepairOrderRepository) Delete(id uint) error {
	return r.db.Delete(&model.RepairOrder{}, id).Error
}

// Get by equipment ID
func (r *RepairOrderRepository) GetByEquipmentID(equipmentID uint, limit int) ([]model.RepairOrder, error) {
	var orders []model.RepairOrder
	err := r.db.Where("equipment_id = ?", equipmentID).
		Order("created_at DESC").
		Limit(limit).
		Find(&orders).Error
	return orders, err
}

// Get pending orders for assignment
func (r *RepairOrderRepository) GetPendingOrders() ([]model.RepairOrder, error) {
	var orders []model.RepairOrder
	err := r.db.Where("status = ?", "pending").
		Preload("Equipment").Preload("Equipment.Type").
		Preload("Reporter").
		Order("priority ASC, created_at ASC").
		Find(&orders).Error
	return orders, err
}

// Get orders by assignee
func (r *RepairOrderRepository) GetByAssignee(assigneeID uint) ([]model.RepairOrder, error) {
	var orders []model.RepairOrder
	err := r.db.Where("assigned_to = ? AND status IN ?", assigneeID, []string{"assigned", "in_progress", "testing"}).
		Preload("Equipment").
		Order("priority ASC, created_at ASC").
		Find(&orders).Error
	return orders, err
}

// Get statistics
func (r *RepairOrderRepository) GetStatistics() (map[string]int64, error) {
	stats := make(map[string]int64)

	var total, pending, inProgress, completed, todayCompleted, todayCreated int64
	r.db.Model(&model.RepairOrder{}).Count(&total)
	r.db.Model(&model.RepairOrder{}).Where("status = ?", "pending").Count(&pending)
	r.db.Model(&model.RepairOrder{}).Where("status IN ?", []string{"assigned", "in_progress", "testing"}).Count(&inProgress)
	r.db.Model(&model.RepairOrder{}).Where("status = ?", "closed").Count(&completed)

	// Today completed
	today := time.Now().Format("2006-01-02")
	r.db.Model(&model.RepairOrder{}).
		Where("status = ? AND DATE(closed_at) = ?", "closed", today).
		Count(&todayCompleted)

	// Today created
	r.db.Model(&model.RepairOrder{}).
		Where("DATE(created_at) = ?", today).
		Count(&todayCreated)

	stats["total"] = total
	stats["pending"] = pending
	stats["in_progress"] = inProgress
	stats["completed"] = completed
	stats["today_completed"] = todayCompleted
	stats["today_created"] = todayCreated

	return stats, nil
}

// Get user statistics
func (r *RepairOrderRepository) GetUserStatistics(userID uint) (map[string]int64, error) {
	stats := make(map[string]int64)

	var pending, inProgress, completed, todayCompleted int64
	// Assigned pending
	r.db.Model(&model.RepairOrder{}).
		Where("assigned_to = ? AND status = ?", userID, "assigned").
		Count(&pending)

	// In progress
	r.db.Model(&model.RepairOrder{}).
		Where("assigned_to = ? AND status IN ?", userID, []string{"in_progress", "testing"}).
		Count(&inProgress)

	// Completed
	r.db.Model(&model.RepairOrder{}).
		Where("assigned_to = ? AND status = ?", userID, "closed").
		Count(&completed)

	// Today completed
	today := time.Now().Format("2006-01-02")
	r.db.Model(&model.RepairOrder{}).
		Where("assigned_to = ? AND status = ? AND DATE(closed_at) = ?", userID, "closed", today).
		Count(&todayCompleted)

	stats["pending"] = pending
	stats["in_progress"] = inProgress
	stats["completed"] = completed
	stats["today_completed"] = todayCompleted

	return stats, nil
}

// Calculate average repair time (in hours)
func (r *RepairOrderRepository) GetAvgRepairTime() (float64, error) {
	var avgHours float64
	err := r.db.Model(&model.RepairOrder{}).
		Select("AVG(EXTRACT(EPOCH FROM (closed_at - created_at)) / 3600)").
		Where("status = ? AND closed_at IS NOT NULL", "closed").
		Scan(&avgHours).Error
	return avgHours, err
}

// Calculate average response time (in minutes) from created to started
func (r *RepairOrderRepository) GetAvgResponseTime() (float64, error) {
	var avgMinutes float64
	err := r.db.Model(&model.RepairOrder{}).
		Select("AVG(EXTRACT(EPOCH FROM (started_at - created_at)) / 60)").
		Where("started_at IS NOT NULL").
		Scan(&avgMinutes).Error
	return avgMinutes, err
}

// RepairLog Repository
type RepairLogRepository struct {
	db *gorm.DB
}

func NewRepairLogRepository() *RepairLogRepository {
	return &RepairLogRepository{db: DB}
}

func (r *RepairLogRepository) Create(log *model.RepairLog) error {
	return r.db.Create(log).Error
}

func (r *RepairLogRepository) GetByOrderID(orderID uint) ([]model.RepairLog, error) {
	var logs []model.RepairLog
	err := r.db.Where("order_id = ?", orderID).
		Preload("User").
		Order("created_at ASC").
		Find(&logs).Error
	return logs, err
}
