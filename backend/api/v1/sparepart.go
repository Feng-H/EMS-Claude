package v1

import (
	"net/http"
	"strconv"

	"github.com/ems/backend/internal/dto"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/internal/repository"
	"github.com/ems/backend/internal/service"
	"github.com/gin-gonic/gin"
)

var (
	sparePartService *service.SparePartService
)

func InitSparePart() {
	sparePartService = service.NewSparePartService()
}

// =====================================================
// Spare Part APIs
// =====================================================

// ListSpareParts returns all spare parts
// @Summary Get spare parts
// @Tags spareparts
// @Router /spareparts [get]
func ListSpareParts(c *gin.Context) {
	var query dto.SparePartQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	filter := repository.SparePartFilter{
		Code:     query.Code,
		Name:     query.Name,
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	result, err := sparePartService.ListParts(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.SparePartResponse, len(result.Items))
	for i, p := range result.Items {
		response[i] = dto.SparePartResponse{
			ID:            p.ID,
			Code:          p.Code,
			Name:          p.Name,
			Specification: p.Specification,
			Unit:          p.Unit,
			SafetyStock:   p.SafetyStock,
			CreatedAt:     p.CreatedAt,
		}
		if p.FactoryID != nil {
			response[i].FactoryID = p.FactoryID
		}
	}

	c.JSON(http.StatusOK, dto.SparePartListResponse{
		Total: result.Total,
		Items: response,
	})
}

// CreateSparePart creates a new spare part
// @Summary Create spare part
// @Tags spareparts
// @Router /spareparts [post]
func CreateSparePart(c *gin.Context) {
	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req dto.CreateSparePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	part, err := sparePartService.CreatePart(
		req.Code,
		req.Name,
		req.Specification,
		req.Unit,
		req.FactoryID,
		req.SafetyStock,
	)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.SparePartResponse{
		ID:            part.ID,
		Code:          part.Code,
		Name:          part.Name,
		Specification: part.Specification,
		Unit:          part.Unit,
		SafetyStock:   part.SafetyStock,
		CreatedAt:     part.CreatedAt,
	})
}

// UpdateSparePart updates a spare part
// @Summary Update spare part
// @Tags spareparts
// @Router /spareparts/:id [put]
func UpdateSparePart(c *gin.Context) {
	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.UpdateSparePartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sparePartService.UpdatePart(uint(id), req.Code, req.Name, req.Specification, req.Unit, req.FactoryID, req.SafetyStock); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

// DeleteSparePart deletes a spare part
// @Summary Delete spare part
// @Tags spareparts
// @Router /spareparts/:id [delete]
func DeleteSparePart(c *gin.Context) {
	if role, _ := middleware.GetUserRole(c); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := sparePartService.DeletePart(uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// =====================================================
// Inventory APIs
// =====================================================

// GetInventory returns inventory list
// @Summary Get inventory
// @Tags spareparts
// @Router /spareparts/inventory [get]
func GetInventory(c *gin.Context) {
	var query dto.InventoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	filter := repository.InventoryFilter{
		SparePartID: query.SparePartID,
		FactoryID:   query.FactoryID,
		Page:        query.Page,
		PageSize:    query.PageSize,
	}

	result, err := sparePartService.GetInventory(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.InventoryResponse, len(result.Items))
	for i, inv := range result.Items {
		isLowStock := false
		if inv.SparePart.ID > 0 && inv.Quantity < inv.SparePart.SafetyStock {
			isLowStock = true
		}

		response[i] = dto.InventoryResponse{
			ID:            inv.ID,
			SparePartID:   inv.SparePartID,
			SparePartCode: inv.SparePart.Code,
			SparePartName: inv.SparePart.Name,
			FactoryID:     inv.FactoryID,
			FactoryName:   inv.Factory.Name,
			Quantity:      inv.Quantity,
			IsLowStock:    isLowStock,
			UpdatedAt:     inv.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.InventoryListResponse{
		Total: result.Total,
		Items: response,
	})
}

// StockIn adds stock to inventory
// @Summary Stock in
// @Tags spareparts
// @Router /spareparts/stock-in [post]
func StockIn(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req dto.StockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sparePartService.StockIn(req.SparePartID, req.FactoryID, req.Quantity, req.Remark, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock added successfully"})
}

// StockOut removes stock from inventory
// @Summary Stock out
// @Tags spareparts
// @Router /spareparts/stock-out [post]
func StockOut(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.StockOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sparePartService.StockOut(req.SparePartID, req.FactoryID, req.Quantity, req.OrderID, req.TaskID, req.Remark, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock removed successfully"})
}

// GetLowStockAlerts returns low stock alerts
// @Summary Get low stock alerts
// @Tags spareparts
// @Router /spareparts/alerts [get]
func GetLowStockAlerts(c *gin.Context) {
	alerts, err := sparePartService.GetLowStockAlerts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.LowStockAlert, len(alerts))
	for i, a := range alerts {
		response[i] = dto.LowStockAlert{
			SparePartID:   a.SparePartID,
			SparePartCode: a.SparePartCode,
			SparePartName: a.SparePartName,
			FactoryID:     a.FactoryID,
			FactoryName:   a.FactoryName,
			CurrentStock:  a.CurrentStock,
			SafetyStock:   a.SafetyStock,
			Shortage:      a.Shortage,
		}
	}

	c.JSON(http.StatusOK, response)
}

// =====================================================
// Consumption APIs
// =====================================================

// GetConsumptionRecords returns consumption records
// @Summary Get consumption records
// @Tags spareparts
// @Router /spareparts/consumptions [get]
func GetConsumptionRecords(c *gin.Context) {
	var query dto.ConsumptionQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	filter := repository.ConsumptionFilter{
		SparePartID: query.SparePartID,
		OrderID:     query.OrderID,
		TaskID:      query.TaskID,
		DateFrom:    query.DateFrom,
		DateTo:      query.DateTo,
		Page:        query.Page,
		PageSize:    query.PageSize,
	}

	result, err := sparePartService.GetConsumption(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.ConsumptionResponse, len(result.Items))
	for i, cons := range result.Items {
		response[i] = dto.ConsumptionResponse{
			ID:            cons.ID,
			SparePartID:   cons.SparePartID,
			SparePartCode: cons.SparePart.Code,
			SparePartName: cons.SparePart.Name,
			OrderID:       cons.OrderID,
			TaskID:        cons.TaskID,
			Quantity:      cons.Quantity,
			UserID:        cons.UserID,
			UserName:      cons.User.Name,
			CreatedAt:     cons.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.ConsumptionListResponse{
		Total: result.Total,
		Items: response,
	})
}

// CreateConsumption creates a consumption record
// @Summary Create consumption
// @Tags spareparts
// @Router /spareparts/consumptions [post]
func CreateConsumption(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.CreateConsumptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sparePartService.RecordConsumption(req.SparePartID, req.Quantity, req.OrderID, req.TaskID, userID, req.Remark); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Consumption recorded successfully"})
}

// =====================================================
// Statistics APIs
// =====================================================

// GetSparePartStatistics returns spare part statistics
// @Summary Get spare part statistics
// @Tags spareparts
// @Router /spareparts/statistics [get]
func GetSparePartStatistics(c *gin.Context) {
	stats, err := sparePartService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SparePartStatistics{
		TotalParts:        stats.TotalParts,
		LowStockCount:     stats.LowStockCount,
		TotalStockValue:   stats.TotalStockValue,
		MonthlyConsumption: stats.MonthlyConsumption,
	})
}
