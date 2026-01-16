package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ems/backend/internal/dto"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
	"github.com/ems/backend/internal/service"
	"github.com/gin-gonic/gin"
)

var (
	inspectionTemplateService *service.InspectionTemplateService
	inspectionItemService     *service.InspectionItemService
	inspectionTaskService     *service.InspectionTaskService
)

func InitInspection() {
	inspectionTemplateService = service.NewInspectionTemplateService()
	inspectionItemService = service.NewInspectionItemService()
	inspectionTaskService = service.NewInspectionTaskService()
}

// =====================================================
// Inspection Template APIs
// =====================================================

// ListInspectionTemplates returns all inspection templates
// @Summary Get inspection templates
// @Tags inspection
// @Router /inspection/templates [get]
func ListInspectionTemplates(c *gin.Context) {
	templates, err := inspectionTemplateService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.InspectionTemplateResponse, len(templates))
	for i, t := range templates {
		response[i] = dto.InspectionTemplateResponse{
			ID:              t.ID,
			Name:            t.Name,
			EquipmentTypeID: t.EquipmentTypeID,
			ItemCount:       len(t.Items),
			CreatedAt:       t.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetInspectionTemplate returns a template with items
// @Summary Get inspection template detail
// @Tags inspection
// @Router /inspection/templates/:id [get]
func GetInspectionTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	template, err := inspectionTemplateService.GetByID(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items, _ := inspectionItemService.GetByTemplateID(template.ID)
	itemResponses := make([]dto.InspectionItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = dto.InspectionItemResponse{
			ID:            item.ID,
			TemplateID:    item.TemplateID,
			Name:          item.Name,
			Method:        item.Method,
			Criteria:      item.Criteria,
			SequenceOrder: item.SequenceOrder,
		}
	}

	c.JSON(http.StatusOK, dto.InspectionTemplateDetailResponse{
		ID:              template.ID,
		Name:            template.Name,
		EquipmentTypeID: template.EquipmentTypeID,
		Items:           itemResponses,
		CreatedAt:       template.CreatedAt,
	})
}

// CreateInspectionTemplate creates a new template
// @Summary Create inspection template
// @Tags inspection
// @Router /inspection/templates [post]
func CreateInspectionTemplate(c *gin.Context) {
	var req dto.InspectionTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	template, err := inspectionTemplateService.Create(req.Name, req.EquipmentTypeID)
	if err != nil {
		handleInspectionServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.InspectionTemplateResponse{
		ID:              template.ID,
		Name:            template.Name,
		EquipmentTypeID: template.EquipmentTypeID,
		CreatedAt:       template.CreatedAt,
	})
}

// =====================================================
// Inspection Item APIs
// =====================================================

// CreateInspectionItem creates a new inspection item
// @Summary Create inspection item
// @Tags inspection
// @Router /inspection/items [post]
func CreateInspectionItem(c *gin.Context) {
	var req dto.InspectionItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	item, err := inspectionItemService.Create(
		req.TemplateID,
		req.Name,
		req.Method,
		req.Criteria,
		req.SequenceOrder,
	)
	if err != nil {
		handleInspectionServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.InspectionItemResponse{
		ID:            item.ID,
		TemplateID:    item.TemplateID,
		Name:          item.Name,
		Method:        item.Method,
		Criteria:      item.Criteria,
		SequenceOrder: item.SequenceOrder,
	})
}

// =====================================================
// Inspection Task APIs
// =====================================================

// ListInspectionTasks returns inspection tasks
// @Summary Get inspection tasks
// @Tags inspection
// @Router /inspection/tasks [get]
func ListInspectionTasks(c *gin.Context) {
	var query dto.InspectionTaskQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default pagination
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	filter := &service.InspectionTaskFilter{
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	if query.AssignedTo > 0 {
		filter.AssignedTo = query.AssignedTo
	}
	if query.Status != "" {
		filter.Status = query.Status
	}
	if query.DateFrom != "" {
		if t, err := time.Parse("2006-01-02", query.DateFrom); err == nil {
			filter.DateFrom = t
		}
	}
	if query.DateTo != "" {
		if t, err := time.Parse("2006-01-02", query.DateTo); err == nil {
			filter.DateTo = t
		}
	}

	result, err := inspectionTaskService.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.InspectionTaskResponse, len(result.Items))
	for i, task := range result.Items {
		response[i] = inspectionTaskToResponse(&task)
	}

	c.JSON(http.StatusOK, dto.InspectionTaskListResponse{
		Total: result.Total,
		Items: response,
	})
}

// GetInspectionTask returns a single task
// @Summary Get inspection task
// @Tags inspection
// @Router /inspection/tasks/:id [get]
func GetInspectionTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	task, err := inspectionTaskService.GetByID(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inspectionTaskToResponse(task))
}

// GetMyTasks returns current user's tasks for today
// @Summary Get my inspection tasks
// @Tags inspection
// @Router /inspection/my-tasks [get]
func GetMyTasks(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current date in local timezone (Asia/Shanghai)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	today := time.Now().In(loc).Truncate(24 * time.Hour)
	tasks, err := inspectionTaskService.GetMyTasks(userID, today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.InspectionTaskResponse, len(tasks))
	for i, task := range tasks {
		response[i] = inspectionTaskToResponse(&task)
	}

	c.JSON(http.StatusOK, response)
}

// GetMyTaskStatistics returns current user's task statistics
// @Summary Get my task statistics
// @Tags inspection
// @Router /inspection/my-stats [get]
func GetMyTaskStatistics(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	stats, err := inspectionTaskService.GetMyStatistics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MyTasksResponse{
		PendingCount:    stats.PendingCount,
		InProgressCount: stats.InProgressCount,
		TodayTasks:      stats.TodayTasks,
	})
}

// StartInspection starts an inspection (with anti-fraud QR code verification)
// @Summary Start inspection
// @Tags inspection
// @Router /inspection/start [post]
func StartInspection(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.StartInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startReq := &service.StartInspectionRequest{
		EquipmentID: req.EquipmentID,
		QRCode:      req.QRCode,
		Latitude:    &req.Latitude,
		Longitude:   &req.Longitude,
		Timestamp:   req.Timestamp,
	}

	task, items, err := inspectionTaskService.StartInspection(userID, startReq)
	if err != nil {
		if err == service.ErrInvalidQRCode {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid QR code or equipment mismatch"})
			return
		}
		if err == service.ErrInvalidTimestamp {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp - request expired"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If task is already completed, return existing records
	if task.Status == "completed" {
		recordRepo := repository.NewInspectionRecordRepository()
		records, _ := recordRepo.GetByTaskID(task.ID)
		recordResponses := make([]dto.InspectionRecordResponse, len(records))
		for i, r := range records {
			recordResponses[i] = dto.InspectionRecordResponse{
				ID:       r.ID,
				TaskID:   r.TaskID,
				Result:   r.Result,
				Remark:   r.Remark,
				PhotoURL: r.PhotoURL,
			}
		}

		c.JSON(http.StatusOK, dto.StartInspectionResponse{
			TaskID:      task.ID,
			EquipmentID: task.EquipmentID,
			StartedAt:   *task.StartedAt,
		})
		return
	}

	itemResponses := make([]dto.InspectionItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = dto.InspectionItemResponse{
			ID:            item.ID,
			TemplateID:    item.TemplateID,
			Name:          item.Name,
			Method:        item.Method,
			Criteria:      item.Criteria,
			SequenceOrder: item.SequenceOrder,
		}
	}

	c.JSON(http.StatusOK, dto.StartInspectionResponse{
		TaskID:      task.ID,
		EquipmentID: task.EquipmentID,
		Items:       itemResponses,
		StartedAt:   time.Now(),
	})
}

// CompleteInspection completes an inspection task
// @Summary Complete inspection
// @Tags inspection
// @Router /inspection/complete [post]
func CompleteInspection(c *gin.Context) {
	var req dto.CompleteInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	records := make([]service.CompleteInspectionRecord, len(req.Records))
	for i, r := range req.Records {
		records[i] = service.CompleteInspectionRecord{
			ItemID:   r.ItemID,
			Result:   r.Result,
			Remark:   r.Remark,
			PhotoURL: r.PhotoURL,
		}
	}

	result, err := inspectionTaskService.CompleteInspection(
		req.TaskID,
		records,
		req.Latitude,
		req.Longitude,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: If there are NG items, trigger repair workflow
	if result.NGCount > 0 {
		// Create repair order for NG items
		// This will be implemented in the repair module
	}

	c.JSON(http.StatusOK, dto.CompleteInspectionResponse{
		TaskID:      result.TaskID,
		CompletedAt: result.CompletedAt,
		TotalCount:  result.TotalCount,
		OKCount:     result.OKCount,
		NGCount:     result.NGCount,
		NGItems:     result.NGItemIDs,
	})
}

// GetInspectionStatistics returns inspection statistics
// @Summary Get inspection statistics
// @Tags inspection
// @Router /inspection/statistics [get]
func GetInspectionStatistics(c *gin.Context) {
	stats, err := inspectionTaskService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.InspectionStatistics{
		TotalTasks:      stats.TotalTasks,
		PendingTasks:    stats.PendingTasks,
		InProgressTasks: stats.InProgressTasks,
		CompletedTasks:  stats.CompletedTasks,
		OverdueTasks:    stats.OverdueTasks,
		TodayCompleted:  stats.TodayCompleted,
		CompletionRate:  stats.CompletionRate,
	})
}

// =====================================================
// Helper Functions
// =====================================================

func inspectionTaskToResponse(task *model.InspectionTask) dto.InspectionTaskResponse {
	r := dto.InspectionTaskResponse{
		ID:            task.ID,
		EquipmentID:   task.EquipmentID,
		TemplateID:    task.TemplateID,
		AssignedTo:    task.AssignedTo,
		ScheduledDate: task.ScheduledDate,
		Status:        string(task.Status),
		StartedAt:     task.StartedAt,
		CompletedAt:   task.CompletedAt,
		Latitude:      task.Latitude,
		Longitude:     task.Longitude,
	}

	if task.Equipment.ID > 0 {
		r.EquipmentCode = task.Equipment.Code
		r.EquipmentName = task.Equipment.Name
	}
	if task.Template.ID > 0 {
		r.TemplateName = task.Template.Name
	}
	if task.Assignee != nil {
		r.AssigneeName = task.Assignee.Name
	}

	// Count items
	r.ItemCount = len(task.Template.Items)
	r.CompletedCount = len(task.Records)

	return r
}

func handleInspectionServiceError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	if err == service.ErrNotFound {
		status = http.StatusNotFound
	} else if err == service.ErrDuplicateCode {
		status = http.StatusConflict
	}
	c.JSON(status, gin.H{"error": err.Error()})
}
