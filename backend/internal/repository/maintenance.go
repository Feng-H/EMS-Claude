package repository

import (
	"time"

	"github.com/ems/backend/internal/model"
	"gorm.io/gorm"
)

// MaintenancePlan Repository
type MaintenancePlanRepository struct {
	db *gorm.DB
}

func NewMaintenancePlanRepository() *MaintenancePlanRepository {
	return &MaintenancePlanRepository{db: DB}
}

func (r *MaintenancePlanRepository) Create(plan *model.MaintenancePlan) error {
	return r.db.Create(plan).Error
}

func (r *MaintenancePlanRepository) GetByID(id uint) (*model.MaintenancePlan, error) {
	var plan model.MaintenancePlan
	err := r.db.Preload("EquipmentType").Preload("Items").First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *MaintenancePlanRepository) List() ([]model.MaintenancePlan, error) {
	var plans []model.MaintenancePlan
	err := r.db.Preload("EquipmentType").Order("level ASC, created_at DESC").Find(&plans).Error
	return plans, err
}

func (r *MaintenancePlanRepository) GetByEquipmentTypeID(typeID uint) ([]model.MaintenancePlan, error) {
	var plans []model.MaintenancePlan
	err := r.db.Where("equipment_type_id = ?", typeID).
		Preload("EquipmentType").Preload("Items").
		Find(&plans).Error
	return plans, err
}

func (r *MaintenancePlanRepository) Update(plan *model.MaintenancePlan) error {
	return r.db.Save(plan).Error
}

func (r *MaintenancePlanRepository) Delete(id uint) error {
	return r.db.Delete(&model.MaintenancePlan{}, id).Error
}

// MaintenancePlanItem Repository
type MaintenancePlanItemRepository struct {
	db *gorm.DB
}

func NewMaintenancePlanItemRepository() *MaintenancePlanItemRepository {
	return &MaintenancePlanItemRepository{db: DB}
}

func (r *MaintenancePlanItemRepository) Create(item *model.MaintenancePlanItem) error {
	return r.db.Create(item).Error
}

func (r *MaintenancePlanItemRepository) GetByPlanID(planID uint) ([]model.MaintenancePlanItem, error) {
	var items []model.MaintenancePlanItem
	err := r.db.Where("plan_id = ?", planID).
		Order("sequence_order").
		Find(&items).Error
	return items, err
}

func (r *MaintenancePlanItemRepository) Update(item *model.MaintenancePlanItem) error {
	return r.db.Save(item).Error
}

func (r *MaintenancePlanItemRepository) Delete(id uint) error {
	return r.db.Delete(&model.MaintenancePlanItem{}, id).Error
}

// MaintenanceTask Repository
type MaintenanceTaskRepository struct {
	db *gorm.DB
}

func NewMaintenanceTaskRepository() *MaintenanceTaskRepository {
	return &MaintenanceTaskRepository{db: DB}
}

func (r *MaintenanceTaskRepository) Create(task *model.MaintenanceTask) error {
	return r.db.Create(task).Error
}

func (r *MaintenanceTaskRepository) CreateBatch(tasks []model.MaintenanceTask) error {
	return r.db.Create(&tasks).Error
}

func (r *MaintenanceTaskRepository) GetByID(id uint) (*model.MaintenanceTask, error) {
	var task model.MaintenanceTask
	err := r.db.Preload("Plan").Preload("Plan.Items").
		Preload("Equipment").Preload("Equipment.Type").
		Preload("Assignee").
		Preload("Records").Preload("Records.Item").
		First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

type MaintenanceTaskFilter struct {
	Status     string
	AssignedTo uint
	DateFrom   time.Time
	DateTo     time.Time
	Page       int
	PageSize   int
}

func (r *MaintenanceTaskRepository) List(filter MaintenanceTaskFilter) ([]model.MaintenanceTask, int64, error) {
	var tasks []model.MaintenanceTask
	var total int64

	query := r.db.Model(&model.MaintenanceTask{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.AssignedTo > 0 {
		query = query.Where("assigned_to = ?", filter.AssignedTo)
	}
	if !filter.DateFrom.IsZero() {
		query = query.Where("scheduled_date >= ?", filter.DateFrom)
	}
	if !filter.DateTo.IsZero() {
		query = query.Where("scheduled_date <= ?", filter.DateTo)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (filter.Page - 1) * filter.PageSize
	err := query.Preload("Plan").Preload("Equipment").
		Preload("Assignee").
		Order("scheduled_date DESC, created_at DESC").
		Offset(offset).Limit(filter.PageSize).Find(&tasks).Error

	return tasks, total, err
}

func (r *MaintenanceTaskRepository) Update(task *model.MaintenanceTask) error {
	return r.db.Save(task).Error
}

func (r *MaintenanceTaskRepository) Delete(id uint) error {
	return r.db.Delete(&model.MaintenanceTask{}, id).Error
}

// Get tasks by equipment and date range
func (r *MaintenanceTaskRepository) GetByEquipmentAndDateRange(equipmentID uint, start, end time.Time) ([]model.MaintenanceTask, error) {
	var tasks []model.MaintenanceTask
	startStr := start.Format("2006-01-02")
	endStr := end.Format("2006-01-02")
	err := r.db.Where("equipment_id = ? AND scheduled_date >= ? AND scheduled_date <= ?",
		equipmentID, startStr, endStr).
		Preload("Plan").
		Order("scheduled_date ASC").
		Find(&tasks).Error
	return tasks, err
}

// Get pending tasks for a user
func (r *MaintenanceTaskRepository) GetPendingTasksForUser(userID uint, date time.Time) ([]model.MaintenanceTask, error) {
	var tasks []model.MaintenanceTask
	dateStr := date.Format("2006-01-02")
	endDateStr := date.AddDate(0, 0, 7).Format("2006-01-02")
	err := r.db.Where("assigned_to = ? AND scheduled_date >= ? AND scheduled_date <= ? AND status IN ?",
		userID, dateStr, endDateStr, []string{"pending", "in_progress"}).
		Preload("Equipment").Preload("Equipment.Type").
		Preload("Plan").
		Order("scheduled_date ASC").
		Find(&tasks).Error
	return tasks, err
}

// Get overdue tasks
func (r *MaintenanceTaskRepository) GetOverdueTasks() ([]model.MaintenanceTask, error) {
	var tasks []model.MaintenanceTask
	today := time.Now().Format("2006-01-02")
	err := r.db.Where("due_date < ? AND status IN ?", today, []string{"pending", "in_progress"}).
		Preload("Equipment").Preload("Assignee").
		Find(&tasks).Error
	return tasks, err
}

// Update overdue status
func (r *MaintenanceTaskRepository) UpdateOverdueStatus(ids []uint) error {
	return r.db.Model(&model.MaintenanceTask{}).
		Where("id IN ?", ids).
		Update("status", "overdue").Error
}

// MaintenanceRecord Repository
type MaintenanceRecordRepository struct {
	db *gorm.DB
}

func NewMaintenanceRecordRepository() *MaintenanceRecordRepository {
	return &MaintenanceRecordRepository{db: DB}
}

func (r *MaintenanceRecordRepository) Create(record *model.MaintenanceRecord) error {
	return r.db.Create(record).Error
}

func (r *MaintenanceRecordRepository) CreateBatch(records []model.MaintenanceRecord) error {
	return r.db.Create(&records).Error
}

func (r *MaintenanceRecordRepository) GetByTaskID(taskID uint) ([]model.MaintenanceRecord, error) {
	var records []model.MaintenanceRecord
	err := r.db.Where("task_id = ?", taskID).Preload("Item").Find(&records).Error
	return records, err
}

// Get statistics
func (r *MaintenanceTaskRepository) GetStatistics() (map[string]int64, error) {
	stats := make(map[string]int64)

	var totalPlans, totalTasks, pending, inProgress, completed, overdue, todayCompleted int64
	r.db.Model(&model.MaintenancePlan{}).Count(&totalPlans)
	r.db.Model(&model.MaintenanceTask{}).Count(&totalTasks)
	r.db.Model(&model.MaintenanceTask{}).Where("status = ?", "pending").Count(&pending)
	r.db.Model(&model.MaintenanceTask{}).Where("status = ?", "in_progress").Count(&inProgress)
	r.db.Model(&model.MaintenanceTask{}).Where("status = ?", "completed").Count(&completed)
	r.db.Model(&model.MaintenanceTask{}).Where("status = ?", "overdue").Count(&overdue)

	// Today completed
	today := time.Now().Format("2006-01-02")
	r.db.Model(&model.MaintenanceTask{}).
		Where("status = ? AND DATE(completed_at) = ?", "completed", today).
		Count(&todayCompleted)

	stats["total_plans"] = totalPlans
	stats["total_tasks"] = totalTasks
	stats["pending"] = pending
	stats["in_progress"] = inProgress
	stats["completed"] = completed
	stats["overdue"] = overdue
	stats["today_completed"] = todayCompleted

	return stats, nil
}

// Get user statistics
func (r *MaintenanceTaskRepository) GetUserStatistics(userID uint) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Pending for today and this week
	today := time.Now()
	todayStart := today.Format("2006-01-02")
	weekEnd := today.AddDate(7, 0, 0).Format("2006-01-02")

	var pending, inProgress, todayCompleted, weekCompleted int64
	r.db.Model(&model.MaintenanceTask{}).
		Where("assigned_to = ? AND scheduled_date >= ? AND scheduled_date < ? AND status = ?",
			userID, todayStart, weekEnd, "pending").
		Count(&pending)

	r.db.Model(&model.MaintenanceTask{}).
		Where("assigned_to = ? AND status = ?", userID, "in_progress").
		Count(&inProgress)

	// Completed today
	r.db.Model(&model.MaintenanceTask{}).
		Where("assigned_to = ? AND status = ? AND DATE(completed_at) = ?", userID, "completed", todayStart).
		Count(&todayCompleted)

	// Completed this week
	r.db.Model(&model.MaintenanceTask{}).
		Where("assigned_to = ? AND status = ? AND completed_at >= ? AND completed_at < ?",
			userID, "completed", todayStart, weekEnd).
		Count(&weekCompleted)

	stats["pending"] = pending
	stats["in_progress"] = inProgress
	stats["today_completed"] = todayCompleted
	stats["week_completed"] = weekCompleted

	return stats, nil
}
