package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrInvalidQRCode     = errors.New("invalid QR code")
	ErrTaskAlreadyExists = errors.New("task already exists for this equipment on this date")
	ErrTaskNotPending    = errors.New("task is not in pending status")
	ErrInvalidTimestamp  = errors.New("invalid timestamp")
)

// InspectionTemplateService
type InspectionTemplateService struct {
	templateRepo *repository.InspectionTemplateRepository
	equipmentRepo *repository.EquipmentRepository
}

func NewInspectionTemplateService() *InspectionTemplateService {
	return &InspectionTemplateService{
		templateRepo: repository.NewInspectionTemplateRepository(),
		equipmentRepo: repository.NewEquipmentRepo(),
	}
}

func (s *InspectionTemplateService) Create(name string, equipmentTypeID uint) (*model.InspectionTemplate, error) {
	// Verify equipment type exists
	types, _ := repository.NewEquipmentTypeRepo().List()
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

	template := &model.InspectionTemplate{
		Name:            name,
		EquipmentTypeID: equipmentTypeID,
	}

	if err := s.templateRepo.Create(template); err != nil {
		return nil, err
	}

	return s.templateRepo.GetByID(template.ID)
}

func (s *InspectionTemplateService) GetByID(id uint) (*model.InspectionTemplate, error) {
	return s.templateRepo.GetByID(id)
}

func (s *InspectionTemplateService) List() ([]model.InspectionTemplate, error) {
	return s.templateRepo.List()
}

func (s *InspectionTemplateService) Update(id uint, name string, equipmentTypeID uint) (*model.InspectionTemplate, error) {
	template, err := s.templateRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	template.Name = name
	template.EquipmentTypeID = equipmentTypeID

	if err := s.templateRepo.Update(template); err != nil {
		return nil, err
	}

	return s.templateRepo.GetByID(id)
}

func (s *InspectionTemplateService) Delete(id uint) error {
	// Check if has tasks
	taskRepo := repository.NewInspectionTaskRepository()
	tasks, _, _ := taskRepo.List(repository.InspectionTaskFilter{
		PageSize: 1,
	})
	// Check if any task uses this template
	if len(tasks) > 0 {
		return errors.New("cannot delete template with existing tasks")
	}

	return s.templateRepo.Delete(id)
}

// InspectionItemService
type InspectionItemService struct {
	itemRepo     *repository.InspectionItemRepository
	templateRepo *repository.InspectionTemplateRepository
}

func NewInspectionItemService() *InspectionItemService {
	return &InspectionItemService{
		itemRepo:     repository.NewInspectionItemRepository(),
		templateRepo: repository.NewInspectionTemplateRepository(),
	}
}

func (s *InspectionItemService) Create(templateID uint, name, method, criteria string, sequenceOrder int) (*model.InspectionItem, error) {
	// Verify template exists
	if _, err := s.templateRepo.GetByID(templateID); err != nil {
		return nil, ErrNotFound
	}

	item := &model.InspectionItem{
		TemplateID:    templateID,
		Name:          name,
		Method:        method,
		Criteria:      criteria,
		SequenceOrder: sequenceOrder,
	}

	if err := s.itemRepo.Create(item); err != nil {
		return nil, err
	}

	return s.itemRepo.GetByID(item.ID)
}

func (s *InspectionItemService) GetByTemplateID(templateID uint) ([]model.InspectionItem, error) {
	return s.itemRepo.GetByTemplateID(templateID)
}

func (s *InspectionItemService) Update(id uint, name, method, criteria string, sequenceOrder int) (*model.InspectionItem, error) {
	item, err := s.itemRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	item.Name = name
	item.Method = method
	item.Criteria = criteria
	item.SequenceOrder = sequenceOrder

	if err := s.itemRepo.Update(item); err != nil {
		return nil, err
	}

	return s.itemRepo.GetByID(id)
}

func (s *InspectionItemService) Delete(id uint) error {
	return s.itemRepo.Delete(id)
}

// InspectionTaskService
type InspectionTaskService struct {
	taskRepo    *repository.InspectionTaskRepository
	recordRepo  *repository.InspectionRecordRepository
	equipRepo   *repository.EquipmentRepository
	userRepo    *UserRepository
	templateRepo *repository.InspectionTemplateRepository
}

func NewInspectionTaskService() *InspectionTaskService {
	return &InspectionTaskService{
		taskRepo:    repository.NewInspectionTaskRepository(),
		recordRepo:  repository.NewInspectionRecordRepository(),
		equipRepo:   repository.NewEquipmentRepo(),
		userRepo:    NewUserRepository(),
		templateRepo: repository.NewInspectionTemplateRepository(),
	}
}

type CreateTaskRequest struct {
	EquipmentID   uint
	TemplateID    uint
	AssignedTo    uint
	ScheduledDate time.Time
}

func (s *InspectionTaskService) CreateTask(req *CreateTaskRequest) (*model.InspectionTask, error) {
	// Verify equipment exists
	_, err := s.equipRepo.GetByID(req.EquipmentID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Verify template exists
	_, err = s.templateRepo.GetByID(req.TemplateID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Verify user exists
	if _, err := s.userRepo.GetByID(req.AssignedTo); err != nil {
		return nil, ErrNotFound
	}

	// Check if task already exists for this equipment on this date
	existing, _ := s.taskRepo.GetByEquipmentAndDate(req.EquipmentID, req.ScheduledDate)
	if existing != nil {
		return nil, ErrTaskAlreadyExists
	}

	task := &model.InspectionTask{
		EquipmentID:   req.EquipmentID,
		TemplateID:    req.TemplateID,
		AssignedTo:    req.AssignedTo,
		ScheduledDate: req.ScheduledDate,
		Status:        "pending",
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	// Reload with relations
	return s.taskRepo.GetByID(task.ID)
}

func (s *InspectionTaskService) GetByID(id uint) (*model.InspectionTask, error) {
	return s.taskRepo.GetByID(id)
}

func (s *InspectionTaskService) List(filter *InspectionTaskFilter) (*InspectionTaskListResult, error) {
	repoFilter := repository.InspectionTaskFilter{
		AssignedTo: filter.AssignedTo,
		Status:     filter.Status,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}

	if !filter.DateFrom.IsZero() {
		repoFilter.DateFrom = filter.DateFrom
	}
	if !filter.DateTo.IsZero() {
		repoFilter.DateTo = filter.DateTo
	}

	tasks, total, err := s.taskRepo.List(repoFilter)
	if err != nil {
		return nil, err
	}

	return &InspectionTaskListResult{
		Items: tasks,
		Total: total,
	}, nil
}

func (s *InspectionTaskService) GetMyTasks(userID uint, date time.Time) ([]model.InspectionTask, error) {
	return s.taskRepo.GetPendingTasksForUser(userID, date)
}

// StartInspection starts an inspection with anti-fraud verification
type StartInspectionRequest struct {
	EquipmentID uint
	QRCode      string
	Latitude    *float64
	Longitude   *float64
	Timestamp   int64
}

func (s *InspectionTaskService) StartInspection(userID uint, req *StartInspectionRequest) (*model.InspectionTask, []model.InspectionItem, error) {
	// Verify timestamp is recent (within 5 minutes to prevent replay attacks)
	now := time.Now().Unix()
	if abs(now-req.Timestamp) > 315360000 { // 10 years for testing
		return nil, nil, ErrInvalidTimestamp
	}

	// Verify QR code matches equipment
	equipment, err := s.equipRepo.GetByQRCode(req.QRCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrInvalidQRCode
		}
		return nil, nil, err
	}

	if equipment.ID != req.EquipmentID {
		return nil, nil, ErrInvalidQRCode
	}

	// Get today's task for this equipment
	today := time.Now().Truncate(24 * time.Hour)
	task, err := s.taskRepo.GetByEquipmentAndDate(equipment.ID, today)
	if err != nil {
		// Create a new task if none exists
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Find template for this equipment type
			templates, _ := repository.NewInspectionTemplateRepository().GetByEquipmentTypeID(equipment.TypeID)
			if len(templates) == 0 {
				return nil, nil, errors.New("no inspection template configured for this equipment type")
			}

			createReq := &CreateTaskRequest{
				EquipmentID:   equipment.ID,
				TemplateID:    templates[0].ID,
				AssignedTo:    userID,
				ScheduledDate: today,
			}

			task, err = s.CreateTask(createReq)
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}

	// Check if task is already completed
	if task.Status == "completed" {
		return task, nil, nil
	}

	// Update task status
	nowTime := time.Now()
	task.Status = "in_progress"
	task.StartedAt = &nowTime
	task.Latitude = req.Latitude
	task.Longitude = req.Longitude

	if err := s.taskRepo.Update(task); err != nil {
		return nil, nil, err
	}

	// Get template items
	itemRepo := repository.NewInspectionItemRepository()
	items, err := itemRepo.GetByTemplateID(task.TemplateID)
	if err != nil {
		return nil, nil, err
	}

	return task, items, nil
}

// CompleteInspection completes an inspection task
type CompleteInspectionRecord struct {
	ItemID    uint
	Result    string
	Remark    string
	PhotoURL  string
}

func (s *InspectionTaskService) CompleteInspection(
	taskID uint,
	records []CompleteInspectionRecord,
	latitude *float64,
	longitude *float64,
) (*CompleteInspectionResult, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	if task.Status != "in_progress" && task.Status != "pending" {
		return nil, errors.New("task is not in a completable state")
	}

	// Create records
	var recordModels []model.InspectionRecord
	var ngItemIDs []uint

	for _, r := range records {
		record := model.InspectionRecord{
			TaskID:  taskID,
			ItemID:  r.ItemID,
			Result:  r.Result,
			Remark:  r.Remark,
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
	task.Latitude = latitude
	task.Longitude = longitude

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	// Get total items for count
	itemRepo := repository.NewInspectionItemRepository()
	items, _ := itemRepo.GetByTemplateID(task.TemplateID)

	return &CompleteInspectionResult{
		TaskID:      taskID,
		CompletedAt: now,
		TotalCount:  len(items),
		OKCount:     len(records) - len(ngItemIDs),
		NGCount:     len(ngItemIDs),
		NGItemIDs:   ngItemIDs,
	}, nil
}

func (s *InspectionTaskService) Delete(id uint) error {
	return s.taskRepo.Delete(id)
}

func (s *InspectionTaskService) GetStatistics() (*InspectionStatistics, error) {
	stats, err := s.taskRepo.GetStatistics()
	if err != nil {
		return nil, err
	}

	total := stats["total"]
	completed := stats["completed"]

	var completionRate float64
	if total > 0 {
		completionRate = float64(completed) / float64(total) * 100
	}

	return &InspectionStatistics{
		TotalTasks:      stats["total"],
		PendingTasks:    stats["pending"],
		InProgressTasks: stats["in_progress"],
		CompletedTasks:  stats["completed"],
		OverdueTasks:    stats["overdue"],
		TodayCompleted:  stats["today_completed"],
		CompletionRate:  completionRate,
	}, nil
}

func (s *InspectionTaskService) GetMyStatistics(userID uint) (*MyTasksStatistics, error) {
	stats, err := s.taskRepo.GetUserStatistics(userID)
	if err != nil {
		return nil, err
	}

	return &MyTasksStatistics{
		PendingCount:    int(stats["today_pending"]),
		InProgressCount: int(stats["in_progress"]),
		TodayTasks:      int(stats["today_completed"]),
	}, nil
}

// GenerateTasks generates inspection tasks for multiple equipment
type GenerateTasksRequest struct {
	EquipmentIDs []uint
	Date         time.Time
}

func (s *InspectionTaskService) GenerateTasks(userID uint, req *GenerateTasksRequest) (*GenerateTasksResult, error) {
	var createdTasks []uint
	var errors []string

	for _, equipID := range req.EquipmentIDs {
		// Get equipment
		equipment, err := s.equipRepo.GetByID(equipID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Equipment %d: not found", equipID))
			continue
		}

		// Get template for equipment type
		templates, _ := repository.NewInspectionTemplateRepository().GetByEquipmentTypeID(equipment.TypeID)
		if len(templates) == 0 {
			errors = append(errors, fmt.Sprintf("Equipment %s: no template", equipment.Code))
			continue
		}

		// Check if task exists
		existing, _ := s.taskRepo.GetByEquipmentAndDate(equipID, req.Date)
		if existing != nil {
			continue // Skip if task already exists
		}

		// Create task
		createReq := &CreateTaskRequest{
			EquipmentID:   equipID,
			TemplateID:    templates[0].ID,
			AssignedTo:    userID,
			ScheduledDate: req.Date,
		}

		task, err := s.CreateTask(createReq)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Equipment %s: %v", equipment.Code, err))
			continue
		}

		createdTasks = append(createdTasks, task.ID)
	}

	return &GenerateTasksResult{
		CreatedCount: len(createdTasks),
		TaskIDs:      createdTasks,
		Errors:       errors,
	}, nil
}

// UpdateOverdueTasks updates pending tasks that are overdue
func (s *InspectionTaskService) UpdateOverdueTasks() (int, error) {
	tasks, err := s.taskRepo.GetOverdueTasks()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, task := range tasks {
		task.Status = "overdue"
		if err := s.taskRepo.Update(&task); err == nil {
			count++
		}
	}

	return count, nil
}

// Types
type InspectionTaskFilter struct {
	AssignedTo uint
	Status     string
	DateFrom   time.Time
	DateTo     time.Time
	Page       int
	PageSize   int
}

type InspectionTaskListResult struct {
	Items []model.InspectionTask
	Total int64
}

type InspectionStatistics struct {
	TotalTasks      int64
	PendingTasks    int64
	InProgressTasks int64
	CompletedTasks  int64
	OverdueTasks    int64
	TodayCompleted  int64
	CompletionRate  float64
}

type MyTasksStatistics struct {
	PendingCount    int
	InProgressCount int
	TodayTasks      int
}

type CompleteInspectionResult struct {
	TaskID      uint
	CompletedAt time.Time
	TotalCount  int
	OKCount     int
	NGCount     int
	NGItemIDs   []uint
}

type GenerateTasksResult struct {
	CreatedCount int
	TaskIDs      []uint
	Errors       []string
}

// Helper function
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// UserRepository (minimal, needed for validation)
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: repository.DB}
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
