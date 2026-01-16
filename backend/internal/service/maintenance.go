package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
)

// MaintenancePlanService
type MaintenancePlanService struct {
	planRepo      *repository.MaintenancePlanRepository
	itemRepo      *repository.MaintenancePlanItemRepository
	taskRepo      *repository.MaintenanceTaskRepository
	equipRepo     *repository.EquipmentRepository
	equipTypeRepo *repository.EquipmentTypeRepository
	userRepo      *UserRepository
}

func NewMaintenancePlanService() *MaintenancePlanService {
	return &MaintenancePlanService{
		planRepo:      repository.NewMaintenancePlanRepository(),
		itemRepo:      repository.NewMaintenancePlanItemRepository(),
		taskRepo:      repository.NewMaintenanceTaskRepository(),
		equipRepo:     repository.NewEquipmentRepo(),
		equipTypeRepo: repository.NewEquipmentTypeRepo(),
		userRepo:      NewUserRepository(),
	}
}

func (s *MaintenancePlanService) CreatePlan(name string, equipmentTypeID uint, level int, cycleDays int, flexibleDays int, workHours float64) (*model.MaintenancePlan, error) {
	// Verify equipment type exists
	types, _ := s.equipTypeRepo.List()
	hasType := false
	for _, t := range types {
		if t.ID == equipmentTypeID {
			hasType = true
			break
		}
	}
	if !hasType {
		return nil, ErrNotFound
	}

	plan := &model.MaintenancePlan{
		Name:            name,
		EquipmentTypeID: equipmentTypeID,
		Level:           level,
		CycleDays:       cycleDays,
		FlexibleDays:    flexibleDays,
		WorkHours:       workHours,
	}

	if err := s.planRepo.Create(plan); err != nil {
		return nil, err
	}

	return s.planRepo.GetByID(plan.ID)
}

func (s *MaintenancePlanService) GetByID(id uint) (*model.MaintenancePlan, error) {
	return s.planRepo.GetByID(id)
}

func (s *MaintenancePlanService) List() ([]model.MaintenancePlan, error) {
	return s.planRepo.List()
}

func (s *MaintenancePlanService) CreateItem(planID uint, name string, method string, criteria string, sequenceOrder int) (*model.MaintenancePlanItem, error) {
	// Verify plan exists
	if _, err := s.planRepo.GetByID(planID); err != nil {
		return nil, ErrNotFound
	}

	item := &model.MaintenancePlanItem{
		PlanID:        planID,
		Name:          name,
		Method:        method,
		Criteria:      criteria,
		SequenceOrder: sequenceOrder,
	}

	if err := s.itemRepo.Create(item); err != nil {
		return nil, err
	}

	// Return by getting through plan
	plan, _ := s.planRepo.GetByID(planID)
	for _, i := range plan.Items {
		if i.ID == item.ID {
			return &i, nil
		}
	}

	return item, nil
}

func (s *MaintenancePlanService) DeletePlan(id uint) error {
	// Check if has tasks
	tasks, _, _ := s.taskRepo.List(repository.MaintenanceTaskFilter{
		PageSize: 1,
	})
	if len(tasks) > 0 {
		return errors.New("cannot delete plan with existing tasks")
	}

	return s.planRepo.Delete(id)
}

// MaintenanceTaskService
type MaintenanceTaskService struct {
	taskRepo        *repository.MaintenanceTaskRepository
	planRepo        *repository.MaintenancePlanRepository
	recordRepo      *repository.MaintenanceRecordRepository
	equipRepo       *repository.EquipmentRepository
	userRepo        *UserRepository
	planItemRepo    *repository.MaintenancePlanItemRepository
}

func NewMaintenanceTaskService() *MaintenanceTaskService {
	return &MaintenanceTaskService{
		taskRepo:     repository.NewMaintenanceTaskRepository(),
		planRepo:     repository.NewMaintenancePlanRepository(),
		recordRepo:   repository.NewMaintenanceRecordRepository(),
		equipRepo:    repository.NewEquipmentRepo(),
		userRepo:     NewUserRepository(),
		planItemRepo: repository.NewMaintenancePlanItemRepository(),
	}
}

// MaintenanceGenerateTasksRequest generates maintenance tasks for multiple equipment
type MaintenanceGenerateTasksRequest struct {
	PlanID       uint
	EquipmentIDs []uint
	BaseDate     time.Time
}

func (s *MaintenanceTaskService) GenerateTasks(userID uint, req *MaintenanceGenerateTasksRequest) (*MaintenanceGenerateTasksResult, error) {
	// Verify plan exists
	plan, err := s.planRepo.GetByID(req.PlanID)
	if err != nil {
		return nil, ErrNotFound
	}

	var createdTasks []uint
	var errors []string

	for _, equipID := range req.EquipmentIDs {
		// Get equipment
		equipment, err := s.equipRepo.GetByID(equipID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Equipment %d: not found", equipID))
			continue
		}

		// Check if equipment type matches plan
		if equipment.TypeID != plan.EquipmentTypeID {
			errors = append(errors, fmt.Sprintf("Equipment %s: type mismatch", equipment.Code))
			continue
		}

		// Check if task already exists for this date range
		existingTasks, _ := s.taskRepo.GetByEquipmentAndDateRange(
			equipID,
			req.BaseDate,
			req.BaseDate.AddDate(0, 0, int(plan.FlexibleDays)+1),
		)

		if len(existingTasks) > 0 {
			continue // Skip if task already exists in flexible window
		}

		// Create task
		task := &model.MaintenanceTask{
			PlanID:        plan.ID,
			EquipmentID:   equipID,
			AssignedTo:    userID, // Will be reassigned later
			ScheduledDate: req.BaseDate.Format("2006-01-02"),
			DueDate:       req.BaseDate.AddDate(0, 0, int(plan.FlexibleDays)).Format("2006-01-02"),
			Status:        "pending",
		}

		if err := s.taskRepo.Create(task); err != nil {
			errors = append(errors, fmt.Sprintf("Equipment %s: %v", equipment.Code, err))
			continue
		}

		createdTasks = append(createdTasks, task.ID)
	}

	return &MaintenanceGenerateTasksResult{
		CreatedCount: len(createdTasks),
		TaskIDs:      createdTasks,
		Errors:       errors,
	}, nil
}

func (s *MaintenanceTaskService) GetByID(id uint) (*model.MaintenanceTask, error) {
	return s.taskRepo.GetByID(id)
}

func (s *MaintenanceTaskService) List(filter *MaintenanceTaskFilter) (*MaintenanceTaskListResult, error) {
	repoFilter := repository.MaintenanceTaskFilter{
		Status:     filter.Status,
		AssignedTo: filter.AssignedTo,
		DateFrom:   filter.DateFrom,
		DateTo:     filter.DateTo,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}

	tasks, total, err := s.taskRepo.List(repoFilter)
	if err != nil {
		return nil, err
	}

	return &MaintenanceTaskListResult{
		Items: tasks,
		Total: total,
	}, nil
}

func (s *MaintenanceTaskService) GetMyTasks(userID uint, date time.Time) ([]model.MaintenanceTask, error) {
	return s.taskRepo.GetPendingTasksForUser(userID, date)
}

// StartMaintenance starts a maintenance task
type StartMaintenanceRequest struct {
	TaskID    uint
	UserID    uint
	Latitude  *float64
	Longitude *float64
}

func (s *MaintenanceTaskService) StartMaintenance(req *StartMaintenanceRequest) error {
	task, err := s.taskRepo.GetByID(req.TaskID)
	if err != nil {
		return err
	}

	if task.Status != "pending" {
		return errors.New("task is not in pending status")
	}

	now := time.Now()
	task.Status = "in_progress"
	task.StartedAt = &now
	task.Latitude = req.Latitude
	task.Longitude = req.Longitude

	return s.taskRepo.Update(task)
}

// MaintenanceItemRecord represents a completed maintenance item
type MaintenanceItemRecord struct {
	ItemID   uint
	Result   string
	Remark   string
	PhotoURL string
}

// CompleteMaintenanceRequest completes a maintenance task
type CompleteMaintenanceRequest struct {
	TaskID      uint
	UserID      uint
	Records     []MaintenanceItemRecord
	Latitude    *float64
	Longitude   *float64
	ActualHours float64
	Remark      string
}

func (s *MaintenanceTaskService) CompleteMaintenance(req *CompleteMaintenanceRequest) (*CompleteMaintenanceResult, error) {
	task, err := s.taskRepo.GetByID(req.TaskID)
	if err != nil {
		return nil, err
	}

	if task.Status != "in_progress" && task.Status != "pending" {
		return nil, errors.New("task is not in a completable state")
	}

	// Create records
	var recordModels []model.MaintenanceRecord
	var ngItemIDs []uint

	for _, r := range req.Records {
		record := model.MaintenanceRecord{
			TaskID:   req.TaskID,
			ItemID:   r.ItemID,
			Result:   r.Result,
			Remark:   r.Remark,
			PhotoURL: r.PhotoURL,
		}
		recordModels = append(recordModels, record)

		if r.Result == "NG" {
			ngItemIDs = append(ngItemIDs, r.ItemID)
		}
	}

	// Save records
	if err := s.recordRepo.CreateBatch(recordModels); err != nil {
		return nil, err
	}

	// Update task
	now := time.Now()
	task.Status = "completed"
	task.CompletedAt = &now
	task.Latitude = req.Latitude
	task.Longitude = req.Longitude
	if req.ActualHours > 0 {
		task.ActualHours = req.ActualHours
	}
	task.Remark = req.Remark

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	// Get total items for count
	items, _ := s.planItemRepo.GetByPlanID(task.PlanID)

	return &CompleteMaintenanceResult{
		TaskID:      req.TaskID,
		CompletedAt: now,
		TotalCount:  len(items),
		OKCount:     len(req.Records) - len(ngItemIDs),
		NGCount:     len(ngItemIDs),
		NGItemIDs:   ngItemIDs,
	}, nil
}

func (s *MaintenanceTaskService) GetStatistics() (*MaintenanceStatistics, error) {
	stats, err := s.taskRepo.GetStatistics()
	if err != nil {
		return nil, err
	}

	total := stats["total_tasks"]
	completed := stats["completed"]

	var completionRate float64
	if total > 0 {
		completionRate = float64(completed) / float64(total) * 100
	}

	return &MaintenanceStatistics{
		TotalPlans:      stats["total_plans"],
		TotalTasks:      stats["total_tasks"],
		PendingTasks:    stats["pending"],
		InProgressTasks: stats["in_progress"],
		CompletedTasks:  stats["completed"],
		OverdueTasks:    stats["overdue"],
		TodayCompleted:  stats["today_completed"],
		CompletionRate:  completionRate,
	}, nil
}

func (s *MaintenanceTaskService) GetUserStatistics(userID uint) (*MyMaintenanceStatistics, error) {
	stats, err := s.taskRepo.GetUserStatistics(userID)
	if err != nil {
		return nil, err
	}

	return &MyMaintenanceStatistics{
		PendingCount:    int(stats["pending"]),
		InProgressCount: int(stats["in_progress"]),
		TodayTasks:      int(stats["today_completed"]),
		ThisWeekTasks:   int(stats["week_completed"]),
	}, nil
}

func (s *MaintenanceTaskService) UpdateOverdueTasks() (int, error) {
	tasks, err := s.taskRepo.GetOverdueTasks()
	if err != nil {
		return 0, err
	}

	var ids []uint
	for _, task := range tasks {
		ids = append(ids, task.ID)
	}

	if len(ids) > 0 {
		err := s.taskRepo.UpdateOverdueStatus(ids)
		if err != nil {
			return 0, err
		}
	}

	return len(ids), nil
}

// Types
type MaintenanceTaskFilter struct {
	Status     string
	AssignedTo uint
	DateFrom   time.Time
	DateTo     time.Time
	Page       int
	PageSize   int
}

type MaintenanceTaskListResult struct {
	Items []model.MaintenanceTask
	Total int64
}

type MaintenanceStatistics struct {
	TotalPlans      int64
	TotalTasks      int64
	PendingTasks    int64
	InProgressTasks int64
	CompletedTasks  int64
	OverdueTasks    int64
	TodayCompleted  int64
	CompletionRate  float64
}

type MyMaintenanceStatistics struct {
	PendingCount    int
	InProgressCount int
	TodayTasks      int
	ThisWeekTasks   int
}

type CompleteMaintenanceResult struct {
	TaskID      uint
	CompletedAt time.Time
	TotalCount  int
	OKCount     int
	NGCount     int
	NGItemIDs   []uint
}

type MaintenanceGenerateTasksResult struct {
	CreatedCount int
	TaskIDs      []uint
	Errors       []string
}
