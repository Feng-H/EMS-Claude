package repository

import (
	"time"

	"github.com/ems/backend/internal/model"
	"gorm.io/gorm"
)

// InspectionTemplate Repository
type InspectionTemplateRepository struct {
	db *gorm.DB
}

func NewInspectionTemplateRepository() *InspectionTemplateRepository {
	return &InspectionTemplateRepository{db: DB}
}

func (r *InspectionTemplateRepository) Create(template *model.InspectionTemplate) error {
	return r.db.Create(template).Error
}

func (r *InspectionTemplateRepository) GetByID(id uint) (*model.InspectionTemplate, error) {
	var template model.InspectionTemplate
	err := r.db.Preload("EquipmentType").Preload("Items").First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *InspectionTemplateRepository) List() ([]model.InspectionTemplate, error) {
	var templates []model.InspectionTemplate
	err := r.db.Preload("EquipmentType").Order("created_at DESC").Find(&templates).Error
	return templates, err
}

func (r *InspectionTemplateRepository) GetByEquipmentTypeID(typeID uint) ([]model.InspectionTemplate, error) {
	var templates []model.InspectionTemplate
	err := r.db.Where("equipment_type_id = ?", typeID).
		Preload("EquipmentType").Find(&templates).Error
	return templates, err
}

func (r *InspectionTemplateRepository) Update(template *model.InspectionTemplate) error {
	return r.db.Save(template).Error
}

func (r *InspectionTemplateRepository) Delete(id uint) error {
	return r.db.Delete(&model.InspectionTemplate{}, id).Error
}

// InspectionItem Repository
type InspectionItemRepository struct {
	db *gorm.DB
}

func NewInspectionItemRepository() *InspectionItemRepository {
	return &InspectionItemRepository{db: DB}
}

func (r *InspectionItemRepository) Create(item *model.InspectionItem) error {
	return r.db.Create(item).Error
}

func (r *InspectionItemRepository) GetByID(id uint) (*model.InspectionItem, error) {
	var item model.InspectionItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *InspectionItemRepository) GetByTemplateID(templateID uint) ([]model.InspectionItem, error) {
	var items []model.InspectionItem
	err := r.db.Where("template_id = ?", templateID).
		Order("sequence_order").Find(&items).Error
	return items, err
}

func (r *InspectionItemRepository) Update(item *model.InspectionItem) error {
	return r.db.Save(item).Error
}

func (r *InspectionItemRepository) Delete(id uint) error {
	return r.db.Delete(&model.InspectionItem{}, id).Error
}

// InspectionTask Repository
type InspectionTaskRepository struct {
	db *gorm.DB
}

func NewInspectionTaskRepository() *InspectionTaskRepository {
	return &InspectionTaskRepository{db: DB}
}

func (r *InspectionTaskRepository) Create(task *model.InspectionTask) error {
	return r.db.Create(task).Error
}

func (r *InspectionTaskRepository) GetByID(id uint) (*model.InspectionTask, error) {
	var task model.InspectionTask
	err := r.db.Preload("Equipment").Preload("Equipment.Type").Preload("Equipment.Workshop").
		Preload("Template").Preload("Template.Items").
		Preload("Assignee").
		Preload("Records").Preload("Records.Item").
		First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

type InspectionTaskFilter struct {
	AssignedTo uint
	Status     string
	DateFrom   time.Time
	DateTo     time.Time
	Page       int
	PageSize   int
}

func (r *InspectionTaskRepository) List(filter InspectionTaskFilter) ([]model.InspectionTask, int64, error) {
	var tasks []model.InspectionTask
	var total int64

	query := r.db.Model(&model.InspectionTask{})

	if filter.AssignedTo > 0 {
		query = query.Where("assigned_to = ?", filter.AssignedTo)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
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
	err := query.Preload("Equipment").Preload("Equipment.Type").
		Preload("Template").
		Preload("Assignee").
		Order("scheduled_date DESC, created_at DESC").
		Offset(offset).Limit(filter.PageSize).Find(&tasks).Error

	return tasks, total, err
}

func (r *InspectionTaskRepository) GetByEquipmentAndDate(equipmentID uint, date time.Time) (*model.InspectionTask, error) {
	var task model.InspectionTask
	err := r.db.Where("equipment_id = ? AND scheduled_date = ?", equipmentID, date).
		Preload("Equipment").Preload("Template").
		First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *InspectionTaskRepository) Update(task *model.InspectionTask) error {
	return r.db.Save(task).Error
}

func (r *InspectionTaskRepository) Delete(id uint) error {
	return r.db.Delete(&model.InspectionTask{}, id).Error
}

// Get pending tasks for a user on a specific date
func (r *InspectionTaskRepository) GetPendingTasksForUser(userID uint, date time.Time) ([]model.InspectionTask, error) {
	var tasks []model.InspectionTask
	// Use date_trunc on date column to compare dates only, avoiding timezone issues
	err := r.db.Where("assigned_to = ? AND date(scheduled_date) = date(?) AND status IN ?",
		userID, date, []string{"pending", "in_progress"}).
		Preload("Equipment").Preload("Equipment.Type").
		Preload("Template").Preload("Template.Items").
		Order("id ASC").
		Find(&tasks).Error
	return tasks, err
}

// Get overdue tasks
func (r *InspectionTaskRepository) GetOverdueTasks() ([]model.InspectionTask, error) {
	var tasks []model.InspectionTask
	today := time.Now().Format("2006-01-02")
	err := r.db.Where("scheduled_date < ? AND status IN ?", today, []string{"pending", "in_progress"}).
		Preload("Equipment").Preload("Assignee").
		Find(&tasks).Error
	return tasks, err
}

// InspectionRecord Repository
type InspectionRecordRepository struct {
	db *gorm.DB
}

func NewInspectionRecordRepository() *InspectionRecordRepository {
	return &InspectionRecordRepository{db: DB}
}

func (r *InspectionRecordRepository) Create(record *model.InspectionRecord) error {
	return r.db.Create(record).Error
}

func (r *InspectionRecordRepository) CreateBatch(records []model.InspectionRecord) error {
	return r.db.Create(&records).Error
}

func (r *InspectionRecordRepository) GetByID(id uint) (*model.InspectionRecord, error) {
	var record model.InspectionRecord
	err := r.db.Preload("Task").Preload("Item").First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *InspectionRecordRepository) GetByTaskID(taskID uint) ([]model.InspectionRecord, error) {
	var records []model.InspectionRecord
	err := r.db.Where("task_id = ?", taskID).Preload("Item").Find(&records).Error
	return records, err
}

func (r *InspectionRecordRepository) Update(record *model.InspectionRecord) error {
	return r.db.Save(record).Error
}

func (r *InspectionRecordRepository) Delete(id uint) error {
	return r.db.Delete(&model.InspectionRecord{}, id).Error
}

// Get NG records from a task
func (r *InspectionRecordRepository) GetNGRecords(taskID uint) ([]model.InspectionRecord, error) {
	var records []model.InspectionRecord
	err := r.db.Where("task_id = ? AND result = ?", taskID, "NG").Find(&records).Error
	return records, err
}

// Check if all items in a task are completed
func (r *InspectionRecordRepository) GetTaskProgress(taskID uint) (total int, completed int, err error) {
	var totalCount int64
	var completedCount int64

	// Get item count from task
	var task model.InspectionTask
	if err := r.db.First(&task, taskID).Error; err != nil {
		return 0, 0, err
	}

	// Count template items
	var itemCount int64
	r.db.Model(&model.InspectionItem{}).Where("template_id = ?", task.TemplateID).Count(&itemCount)
	totalCount = itemCount

	// Count completed records
	r.db.Model(&model.InspectionRecord{}).Where("task_id = ?", taskID).Count(&completedCount)

	return int(totalCount), int(completedCount), nil
}

// Statistics
func (r *InspectionTaskRepository) GetStatistics() (map[string]int64, error) {
	stats := make(map[string]int64)

	var total, pending, inProgress, completed, overdue, todayCompleted int64
	r.db.Model(&model.InspectionTask{}).Count(&total)
	r.db.Model(&model.InspectionTask{}).Where("status = ?", "pending").Count(&pending)
	r.db.Model(&model.InspectionTask{}).Where("status = ?", "in_progress").Count(&inProgress)
	r.db.Model(&model.InspectionTask{}).Where("status = ?", "completed").Count(&completed)
	r.db.Model(&model.InspectionTask{}).Where("status = ?", "overdue").Count(&overdue)

	// Today completed
	today := time.Now().Format("2006-01-02")
	r.db.Model(&model.InspectionTask{}).
		Where("status = ? AND DATE(completed_at) = ?", "completed", today).
		Count(&todayCompleted)

	stats["total"] = total
	stats["pending"] = pending
	stats["in_progress"] = inProgress
	stats["completed"] = completed
	stats["overdue"] = overdue
	stats["today_completed"] = todayCompleted

	return stats, nil
}

func (r *InspectionTaskRepository) GetUserStatistics(userID uint) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Pending count for today
	today := time.Now()
	todayStr := today.Format("2006-01-02")

	var todayPending, inProgress, todayCompleted int64
	r.db.Model(&model.InspectionTask{}).
		Where("assigned_to = ? AND scheduled_date = ? AND status = ?", userID, todayStr, "pending").
		Count(&todayPending)

	// In progress count
	r.db.Model(&model.InspectionTask{}).
		Where("assigned_to = ? AND status = ?", userID, "in_progress").
		Count(&inProgress)

	// Completed today
	r.db.Model(&model.InspectionTask{}).
		Where("assigned_to = ? AND status = ? AND DATE(completed_at) = ?", userID, "completed", todayStr).
		Count(&todayCompleted)

	stats["today_pending"] = todayPending
	stats["in_progress"] = inProgress
	stats["today_completed"] = todayCompleted

	return stats, nil
}
