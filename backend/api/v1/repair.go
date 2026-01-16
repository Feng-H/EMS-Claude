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
	repairOrderService *service.RepairOrderService
)

func InitRepair() {
	repairOrderService = service.NewRepairOrderService()
}

// =====================================================
// Repair Order APIs
// =====================================================

// CreateRepairOrder creates a new repair order
// @Summary Create repair order
// @Tags repair
// @Router /repair/orders [post]
func CreateRepairOrder(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.CreateRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createReq := &service.CreateOrderRequest{
		EquipmentID:      req.EquipmentID,
		FaultDescription: req.FaultDescription,
		FaultCode:        req.FaultCode,
		Photos:           req.Photos,
		Priority:         req.Priority,
		ReporterID:       userID,
	}

	order, err := repairOrderService.CreateOrder(createReq)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, orderToResponse(order))
}

// ListRepairOrders returns repair orders
// @Summary Get repair orders
// @Tags repair
// @Router /repair/orders [get]
func ListRepairOrders(c *gin.Context) {
	var query dto.RepairTaskQuery
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

	filter := &service.RepairOrderFilter{
		Status:     query.Status,
		Priority:   query.Priority,
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

	result, err := repairOrderService.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.RepairOrderResponse, len(result.Items))
	for i, order := range result.Items {
		response[i] = orderToResponse(&order)
	}

	c.JSON(http.StatusOK, dto.RepairOrderListResponse{
		Total: result.Total,
		Items: response,
	})
}

// GetRepairOrder returns a single repair order
// @Summary Get repair order
// @Tags repair
// @Router /repair/orders/:id [get]
func GetRepairOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	order, err := repairOrderService.GetByID(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.RepairOrderDetailResponse{
		RepairOrderResponse: orderToResponse(order),
	})
}

// AssignRepairOrder assigns a repair order to a technician
// @Summary Assign repair order
// @Tags repair
// @Router /repair/orders/:id/assign [post]
func AssignRepairOrder(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.AssignRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := repairOrderService.AssignOrder(uint(id), req.AssignTo, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderToResponse(order))
}

// StartRepair starts a repair task
// @Summary Start repair
// @Tags repair
// @Router /repair/orders/:id/start [post]
func StartRepair(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.StartRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repairOrderService.StartRepair(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repair started"})
}

// UpdateRepair updates repair progress
// @Summary Update repair progress
// @Tags repair
// @Router /repair/orders/:id/update [post]
func UpdateRepair(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.UpdateRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateReq := &service.UpdateRepairRequest{
		Solution:    req.Solution,
		SpareParts:  req.SpareParts,
		ActualHours: req.ActualHours,
		Photos:      req.Photos,
		NextStatus:  req.NextStatus,
	}

	if err := repairOrderService.UpdateRepair(uint(id), userID, updateReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repair updated"})
}

// ConfirmRepair confirms or rejects a completed repair
// @Summary Confirm repair
// @Tags repair
// @Router /repair/orders/:id/confirm [post]
func ConfirmRepair(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.ConfirmRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repairOrderService.ConfirmRepair(uint(id), userID, req.Accepted, req.Comment, req.Photos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repair confirmed"})
}

// AuditRepair audits a confirmed repair
// @Summary Audit repair
// @Tags repair
// @Router /repair/orders/:id/audit [post]
func AuditRepair(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check permissions
	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" && role != "supervisor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.AuditRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repairOrderService.AuditRepair(uint(id), userID, req.Approved, req.Comment, &req.ActualHours); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repair audited"})
}

// GetMyRepairTasks returns current user's active repair tasks
// @Summary Get my repair tasks
// @Tags repair
// @Router /repair/my-tasks [get]
func GetMyRepairTasks(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	tasks, err := repairOrderService.GetMyTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.RepairOrderResponse, len(tasks))
	for i, task := range tasks {
		response[i] = orderToResponse(&task)
	}

	c.JSON(http.StatusOK, response)
}

// GetMyRepairStatistics returns current user's repair statistics
// @Summary Get my repair statistics
// @Tags repair
// @Router /repair/my-stats [get]
func GetMyRepairStatistics(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	stats, err := repairOrderService.GetUserStatistics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MyRepairStatistics{
		PendingCount:    stats.PendingCount,
		InProgressCount: stats.InProgressCount,
		CompletedCount:  stats.CompletedCount,
		TodayCompleted:  stats.TodayCompleted,
	})
}

// GetRepairStatistics returns repair statistics
// @Summary Get repair statistics
// @Tags repair
// @Router /repair/statistics [get]
func GetRepairStatistics(c *gin.Context) {
	stats, err := repairOrderService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.RepairStatistics{
		TotalOrders:      stats.TotalOrders,
		PendingOrders:    stats.PendingOrders,
		InProgressOrders: stats.InProgressOrders,
		CompletedOrders:  stats.CompletedOrders,
		TodayCompleted:   stats.TodayCompleted,
		TodayCreated:     stats.TodayCreated,
		AvgRepairTime:    stats.AvgRepairTime,
		AvgResponseTime:  stats.AvgResponseTime,
	})
}

// =====================================================
// Helper Functions
// =====================================================

func orderToResponse(order *model.RepairOrder) dto.RepairOrderResponse {
	r := dto.RepairOrderResponse{
		ID:               order.ID,
		EquipmentID:      order.EquipmentID,
		FaultDescription: order.FaultDescription,
		FaultCode:        order.FaultCode,
		ReporterID:       order.ReporterID,
		Status:           dto.RepairStatus(order.Status),
		Priority:         order.Priority,
		Photos:           order.Photos,
		Solution:         order.Solution,
		CreatedAt:        order.CreatedAt,
		StartedAt:        order.StartedAt,
		CompletedAt:      order.CompletedAt,
		ConfirmedAt:      order.ConfirmedAt,
		AuditedAt:        order.AuditedAt,
	}

	if order.Equipment.ID > 0 {
		r.EquipmentCode = order.Equipment.Code
		r.EquipmentName = order.Equipment.Name
	}
	if order.Reporter.ID > 0 {
		r.ReporterName = order.Reporter.Name
	}
	if order.Assignee != nil && order.Assignee.ID > 0 {
		r.AssignedTo = order.Assignee.ID
		r.AssigneeName = order.Assignee.Name
	}

	return r
}
