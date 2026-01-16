package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ems/backend/internal/dto"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/service"
	"github.com/gin-gonic/gin"
)

var (
	maintenancePlanService *service.MaintenancePlanService
	maintenanceTaskService *service.MaintenanceTaskService
)

func InitMaintenance() {
	maintenancePlanService = service.NewMaintenancePlanService()
	maintenanceTaskService = service.NewMaintenanceTaskService()
}

// =====================================================
// Maintenance Plan APIs
// =====================================================

// ListMaintenancePlans returns all maintenance plans
// @Summary Get maintenance plans
// @Tags maintenance
// @Router /maintenance/plans [get]
func ListMaintenancePlans(c *gin.Context) {
	plans, err := maintenancePlanService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.MaintenancePlanResponse, len(plans))
	for i, p := range plans {
		response[i] = dto.MaintenancePlanResponse{
			ID:              p.ID,
			Name:            p.Name,
			EquipmentTypeID: p.EquipmentTypeID,
			Level:           p.Level,
			LevelName:       getLevelName(p.Level),
			CycleDays:       p.CycleDays,
			FlexibleDays:    p.FlexibleDays,
			WorkHours:       p.WorkHours,
			ItemCount:       len(p.Items),
			CreatedAt:       p.CreatedAt,
		}
		if p.EquipmentType.ID > 0 {
			response[i].EquipmentTypeName = p.EquipmentType.Name
		}
	}

	c.JSON(http.StatusOK, response)
}

// CreateMaintenancePlan creates a new maintenance plan
// @Summary Create maintenance plan
// @Tags maintenance
// @Router /maintenance/plans [post]
func CreateMaintenancePlan(c *gin.Context) {
	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req dto.CreateMaintenancePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := maintenancePlanService.CreatePlan(
		req.Name,
		req.EquipmentTypeID,
		req.Level,
		req.CycleDays,
		req.FlexibleDays,
		req.WorkHours,
	)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MaintenancePlanResponse{
		ID:              plan.ID,
		Name:            plan.Name,
		EquipmentTypeID: plan.EquipmentTypeID,
		Level:           plan.Level,
		LevelName:       getLevelName(plan.Level),
		CycleDays:       plan.CycleDays,
		FlexibleDays:    plan.FlexibleDays,
		WorkHours:       plan.WorkHours,
		ItemCount:       0,
		CreatedAt:       plan.CreatedAt,
	})
}

// CreateMaintenanceItem creates a new maintenance item
// @Summary Create maintenance item
// @Tags maintenance
// @Router /maintenance/items [post]
func CreateMaintenanceItem(c *gin.Context) {
	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req dto.CreateMaintenanceItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := maintenancePlanService.CreateItem(
		req.PlanID,
		req.Name,
		req.Method,
		req.Criteria,
		req.SequenceOrder,
	)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MaintenanceItemResponse{
		ID:            item.ID,
		PlanID:        item.PlanID,
		Name:          item.Name,
		Method:        item.Method,
		Criteria:      item.Criteria,
		SequenceOrder: item.SequenceOrder,
	})
}

// =====================================================
// Maintenance Task APIs
// =====================================================

// GenerateMaintenanceTasks generates maintenance tasks
// @Summary Generate maintenance tasks
// @Tags maintenance
// @Router /maintenance/tasks/generate [post]
func GenerateMaintenanceTasks(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" && role != "supervisor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req dto.GenerateMaintenanceTasksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	baseDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	genReq := &service.MaintenanceGenerateTasksRequest{
		PlanID:       req.PlanID,
		EquipmentIDs: req.EquipmentIDs,
		BaseDate:     baseDate,
	}

	result, err := maintenanceTaskService.GenerateTasks(userID, genReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"created_count": result.CreatedCount,
		"task_ids":      result.TaskIDs,
		"errors":        result.Errors,
	})
}

// ListMaintenanceTasks returns maintenance tasks
// @Summary Get maintenance tasks
// @Tags maintenance
// @Router /maintenance/tasks [get]
func ListMaintenanceTasks(c *gin.Context) {
	var query dto.MaintenanceTaskQuery
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

	filter := &service.MaintenanceTaskFilter{
		Status:     query.Status,
		AssignedTo: query.AssignedTo,
		Page:       query.Page,
		PageSize:   query.PageSize,
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

	result, err := maintenanceTaskService.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.MaintenanceTaskResponse, len(result.Items))
	for i, task := range result.Items {
		response[i] = taskToResponse(&task)
	}

	c.JSON(http.StatusOK, dto.MaintenanceTaskListResponse{
		Total: result.Total,
		Items: response,
	})
}

// GetMaintenanceTask returns a single task
// @Summary Get maintenance task
// @Tags maintenance
// @Router /maintenance/tasks/:id [get]
func GetMaintenanceTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	task, err := maintenanceTaskService.GetByID(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, taskToResponse(task))
}

// GetMyMaintenanceTasks returns current user's tasks
// @Summary Get my maintenance tasks
// @Tags maintenance
// @Router /maintenance/my-tasks [get]
func GetMyMaintenanceTasks(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	tasks, err := maintenanceTaskService.GetMyTasks(userID, today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.MaintenanceTaskResponse, len(tasks))
	for i, task := range tasks {
		response[i] = taskToResponse(&task)
	}

	c.JSON(http.StatusOK, response)
}

// StartMaintenance starts a maintenance task
// @Summary Start maintenance
// @Tags maintenance
// @Router /maintenance/start [post]
func StartMaintenance(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.StartMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startReq := &service.StartMaintenanceRequest{
		TaskID:    req.TaskID,
		UserID:    userID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	if err := maintenanceTaskService.StartMaintenance(startReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Maintenance started"})
}

// CompleteMaintenance completes a maintenance task
// @Summary Complete maintenance
// @Tags maintenance
// @Router /maintenance/complete [post]
func CompleteMaintenance(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.CompleteMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	records := make([]service.MaintenanceItemRecord, len(req.Records))
	for i, r := range req.Records {
		records[i] = service.MaintenanceItemRecord{
			ItemID:   r.ItemID,
			Result:   r.Result,
			Remark:   r.Remark,
			PhotoURL: r.PhotoURL,
		}
	}

	completeReq := &service.CompleteMaintenanceRequest{
		TaskID:      req.TaskID,
		UserID:      userID,
		Records:     records,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		ActualHours: req.ActualHours,
		Remark:      req.Remark,
	}

	result, err := maintenanceTaskService.CompleteMaintenance(completeReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// TODO: If there are NG items, trigger inspection or repair

	c.JSON(http.StatusOK, dto.CompleteMaintenanceResponse{
		TaskID:      result.TaskID,
		CompletedAt:  result.CompletedAt.Format("2006-01-02T15:04:05Z"),
		TotalCount:  result.TotalCount,
		OKCount:     result.OKCount,
		NGCount:     result.NGCount,
		NGItemIDs:  result.NGItemIDs,
	})
}

// GetMaintenanceStatistics returns maintenance statistics
// @Summary Get maintenance statistics
// @Tags maintenance
// @Router /maintenance/statistics [get]
func GetMaintenanceStatistics(c *gin.Context) {
	stats, err := maintenanceTaskService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MaintenanceStatistics{
		TotalPlans:      stats.TotalPlans,
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

func getLevelName(level int) string {
	names := map[int]string{
		1: "一级保养",
		2: "二级保养",
		3: "精度保养",
	}
	if name, ok := names[level]; ok {
		return name
	}
	return "未知"
}

func taskToResponse(task *model.MaintenanceTask) dto.MaintenanceTaskResponse {
	r := dto.MaintenanceTaskResponse{
		ID:            task.ID,
		PlanID:        task.PlanID,
		EquipmentID:   task.EquipmentID,
		ScheduledDate: task.ScheduledDate,
		DueDate:       task.DueDate,
		Status:        dto.MaintenanceTaskStatus(task.Status),
		ActualHours:   task.ActualHours,
		Remark:        task.Remark,
		StartedAt:     task.StartedAt,
		CompletedAt:   task.CompletedAt,
		CreatedAt:     task.CreatedAt,
	}

	if task.Plan.ID > 0 {
		r.PlanName = task.Plan.Name
	}
	if task.Equipment.ID > 0 {
		r.EquipmentCode = task.Equipment.Code
		r.EquipmentName = task.Equipment.Name
	}
	if task.Assignee != nil && task.Assignee.ID > 0 {
		r.AssignedTo = task.Assignee.ID
		r.AssigneeName = task.Assignee.Name
	}

	// Count items
	if task.Plan.ID > 0 {
		r.ItemCount = len(task.Plan.Items)
	}
	r.CompletedCount = len(task.Records)

	return r
}
