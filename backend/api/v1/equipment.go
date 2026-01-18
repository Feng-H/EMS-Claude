package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ems/backend/internal/dto"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
	"github.com/ems/backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	baseService       *service.BaseService
	factoryService    *service.FactoryService
	workshopService   *service.WorkshopService
	equipmentTypeService *service.EquipmentTypeService
	equipmentService  *service.EquipmentService
)

func InitEquipment(database *gorm.DB) {
	repository.Init(database)

	baseService = service.NewBaseService()
	factoryService = service.NewFactoryService()
	workshopService = service.NewWorkshopService()
	equipmentTypeService = service.NewEquipmentTypeService()
	equipmentService = service.NewEquipmentService()
}

// =====================================================
// Organization APIs
// =====================================================

// ListBases returns all bases
// @Summary Get all bases
// @Description Get list of all bases
// @Tags organization
// @Produce json
// @Security Bearer
// @Success 200 {array} dto.BaseResponse
// @Router /organization/bases [get]
func ListBases(c *gin.Context) {
	bases, err := baseService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.BaseResponse, len(bases))
	for i, b := range bases {
		response[i] = dto.BaseResponse{
			ID:        b.ID,
			Code:      b.Code,
			Name:      b.Name,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// CreateBase creates a new base
// @Summary Create base
// @Description Create a new base
// @Tags organization
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.BaseRequest true "Base data"
// @Success 201 {object} dto.BaseResponse
// @Router /organization/bases [post]
func CreateBase(c *gin.Context) {
	var req dto.BaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check admin permission
	if role, _ := middleware.GetUserRole(c); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can create bases"})
		return
	}

	base, err := baseService.Create(req.Code, req.Name)
	if err != nil {
		if err == service.ErrDuplicateCode {
			c.JSON(http.StatusConflict, gin.H{"error": "Code already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.BaseResponse{
		ID:        base.ID,
		Code:      base.Code,
		Name:      base.Name,
		CreatedAt: base.CreatedAt,
		UpdatedAt: base.UpdatedAt,
	})
}

// UpdateBase updates a base
// @Summary Update base
// @Description Update an existing base
// @Tags organization
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Base ID"
// @Param request body dto.BaseRequest true "Base data"
// @Success 200 {object} dto.BaseResponse
// @Router /organization/bases/{id} [put]
func UpdateBase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.BaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can update bases"})
		return
	}

	base, err := baseService.Update(uint(id), req.Code, req.Name)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Base not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.BaseResponse{
		ID:        base.ID,
		Code:      base.Code,
		Name:      base.Name,
		CreatedAt: base.CreatedAt,
		UpdatedAt: base.UpdatedAt,
	})
}

// DeleteBase deletes a base
// @Summary Delete base
// @Description Delete a base
// @Tags organization
// @Produce json
// @Security Bearer
// @Param id path int true "Base ID"
// @Success 204
// @Router /organization/bases/{id} [delete]
func DeleteBase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can delete bases"})
		return
	}

	if err := baseService.Delete(uint(id)); err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Base not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListFactories returns all factories
// @Summary Get all factories
// @Description Get list of all factories
// @Tags organization
// @Produce json
// @Security Bearer
// @Param base_id query int false "Filter by base ID"
// @Success 200 {array} dto.FactoryResponse
// @Router /organization/factories [get]
func ListFactories(c *gin.Context) {
	var factories []model.Factory
	var err error

	baseIDStr := c.Query("base_id")
	if baseIDStr != "" {
		baseID, _ := strconv.ParseUint(baseIDStr, 10, 32)
		factories, err = factoryService.ListByBaseID(uint(baseID))
	} else {
		factories, err = factoryService.List()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.FactoryResponse, len(factories))
	for i, f := range factories {
		response[i] = dto.FactoryResponse{
			ID:        f.ID,
			BaseID:    f.BaseID,
			BaseName:  f.Base.Name,
			Code:      f.Code,
			Name:      f.Name,
			CreatedAt: f.CreatedAt,
			UpdatedAt: f.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// CreateFactory creates a new factory
// @Summary Create factory
// @Description Create a new factory
// @Tags organization
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.FactoryRequest true "Factory data"
// @Success 201 {object} dto.FactoryResponse
// @Router /organization/factories [post]
func CreateFactory(c *gin.Context) {
	var req dto.FactoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can create factories"})
		return
	}

	factory, err := factoryService.Create(req.BaseID, req.Code, req.Name)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Base not found"})
			return
		}
		if err == service.ErrDuplicateCode {
			c.JSON(http.StatusConflict, gin.H{"error": "Code already exists in this base"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.FactoryResponse{
		ID:        factory.ID,
		BaseID:    factory.BaseID,
		Code:      factory.Code,
		Name:      factory.Name,
		CreatedAt: factory.CreatedAt,
		UpdatedAt: factory.UpdatedAt,
	})
}

// UpdateFactory updates a factory
// @Summary Update factory
// @Tags organization
// @Router /organization/factories/{id} [put]
func UpdateFactory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.FactoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can update factories"})
		return
	}

	factory, err := factoryService.Update(uint(id), req.BaseID, req.Code, req.Name)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.FactoryResponse{
		ID:        factory.ID,
		BaseID:    factory.BaseID,
		Code:      factory.Code,
		Name:      factory.Name,
		CreatedAt: factory.CreatedAt,
		UpdatedAt: factory.UpdatedAt,
	})
}

// DeleteFactory deletes a factory
// @Summary Delete factory
// @Router /organization/factories/{id} [delete]
func DeleteFactory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can delete factories"})
		return
	}

	if err := factoryService.Delete(uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ListWorkshops returns all workshops
// @Summary Get all workshops
// @Tags organization
// @Router /organization/workshops [get]
func ListWorkshops(c *gin.Context) {
	var workshops []model.Workshop
	var err error

	factoryIDStr := c.Query("factory_id")
	if factoryIDStr != "" {
		factoryID, _ := strconv.ParseUint(factoryIDStr, 10, 32)
		workshops, err = workshopService.ListByFactoryID(uint(factoryID))
	} else {
		workshops, err = workshopService.List()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.WorkshopResponse, len(workshops))
	for i, w := range workshops {
		response[i] = dto.WorkshopResponse{
			ID:          w.ID,
			FactoryID:   w.FactoryID,
			FactoryName: w.Factory.Name,
			Code:        w.Code,
			Name:        w.Name,
			CreatedAt:   w.CreatedAt,
			UpdatedAt:   w.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// CreateWorkshop creates a new workshop
// @Summary Create workshop
// @Tags organization
// @Router /organization/workshops [post]
func CreateWorkshop(c *gin.Context) {
	var req dto.WorkshopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can create workshops"})
		return
	}

	workshop, err := workshopService.Create(req.FactoryID, req.Code, req.Name)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.WorkshopResponse{
		ID:        workshop.ID,
		FactoryID: workshop.FactoryID,
		Code:      workshop.Code,
		Name:      workshop.Name,
		CreatedAt: workshop.CreatedAt,
		UpdatedAt: workshop.UpdatedAt,
	})
}

// UpdateWorkshop updates a workshop
// @Summary Update workshop
// @Router /organization/workshops/{id} [put]
func UpdateWorkshop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.WorkshopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can update workshops"})
		return
	}

	workshop, err := workshopService.Update(uint(id), req.FactoryID, req.Code, req.Name)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.WorkshopResponse{
		ID:        workshop.ID,
		FactoryID: workshop.FactoryID,
		Code:      workshop.Code,
		Name:      workshop.Name,
		CreatedAt: workshop.CreatedAt,
		UpdatedAt: workshop.UpdatedAt,
	})
}

// DeleteWorkshop deletes a workshop
// @Summary Delete workshop
// @Router /organization/workshops/{id} [delete]
func DeleteWorkshop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can delete workshops"})
		return
	}

	if err := workshopService.Delete(uint(id)); err != nil {
		handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// =====================================================
// Equipment Type APIs
// =====================================================

// ListEquipmentTypes returns all equipment types
// @Summary Get all equipment types
// @Tags equipment
// @Router /equipment/types [get]
func ListEquipmentTypes(c *gin.Context) {
	types, err := equipmentTypeService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.EquipmentTypeResponse, len(types))
	for i, t := range types {
		response[i] = dto.EquipmentTypeResponse{
			ID:         t.ID,
			Name:       t.Name,
			Category:   t.Category,
			TemplateID: t.InspectionTemplateID,
			CreatedAt:  t.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// CreateEquipmentType creates a new equipment type
// @Summary Create equipment type
// @Tags equipment
// @Router /equipment/types [post]
func CreateEquipmentType(c *gin.Context) {
	var req dto.EquipmentTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	equipmentType, err := equipmentTypeService.Create(req.Name, req.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.EquipmentTypeResponse{
		ID:        equipmentType.ID,
		Name:      equipmentType.Name,
		Category:  equipmentType.Category,
		CreatedAt: equipmentType.CreatedAt,
	})
}

// =====================================================
// Equipment APIs
// =====================================================

// ListEquipment returns paginated equipment list
// @Summary Get equipment list
// @Tags equipment
// @Router /equipment [get]
func ListEquipment(c *gin.Context) {
	var query dto.EquipmentQuery
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

	filter := &service.EquipmentFilter{
		Code:       query.Code,
		Name:       query.Name,
		TypeID:     query.TypeID,
		FactoryID:  query.FactoryID,
		WorkshopID: query.WorkshopID,
		Status:     query.Status,
		Page:       query.Page,
		PageSize:   query.PageSize,
	}

	result, err := equipmentService.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.EquipmentResponse, len(result.Items))
	for i, e := range result.Items {
		response[i] = equipmentToResponse(&e)
	}

	c.JSON(http.StatusOK, dto.EquipmentListResponse{
		Total: result.Total,
		Items: response,
	})
}

// GetEquipment returns a single equipment by ID
// @Summary Get equipment by ID
// @Tags equipment
// @Router /equipment/{id} [get]
func GetEquipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	equipment, err := equipmentService.GetByID(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, equipmentToResponse(equipment))
}

// GetEquipmentByQRCode returns equipment by QR code
// @Summary Get equipment by QR code
// @Tags equipment
// @Router /equipment/qr/{code} [get]
func GetEquipmentByQRCode(c *gin.Context) {
	qrCode := c.Param("code")

	equipment, err := equipmentService.GetByQRCode(qrCode)
	if err != nil {
		if err == service.ErrNotFound || errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, equipmentToResponse(equipment))
}

// CreateEquipment creates new equipment
// @Summary Create equipment
// @Tags equipment
// @Router /equipment [post]
func CreateEquipment(c *gin.Context) {
	var req dto.EquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	createReq := &service.CreateEquipmentRequest{
		Code:                  req.Code,
		Name:                  req.Name,
		TypeID:                req.TypeID,
		WorkshopID:            req.WorkshopID,
		Spec:                  req.Spec,
		PurchaseDate:          req.PurchaseDate,
		Status:                req.Status,
		DedicatedMaintenanceID: req.DedicatedMaintenanceID,
	}

	equipment, err := equipmentService.Create(createReq)
	if err != nil {
		if err == service.ErrDuplicateCode {
			c.JSON(http.StatusConflict, gin.H{"error": "Equipment code already exists"})
			return
		}
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, equipmentToResponse(equipment))
}

// UpdateEquipment updates equipment
// @Summary Update equipment
// @Tags equipment
// @Router /equipment/{id} [put]
func UpdateEquipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.EquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	updateReq := &service.UpdateEquipmentRequest{
		Code:                  req.Code,
		Name:                  req.Name,
		TypeID:                req.TypeID,
		WorkshopID:            req.WorkshopID,
		Spec:                  req.Spec,
		PurchaseDate:          req.PurchaseDate,
		Status:                req.Status,
		DedicatedMaintenanceID: req.DedicatedMaintenanceID,
	}

	equipment, err := equipmentService.Update(uint(id), updateReq)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, equipmentToResponse(equipment))
}

// DeleteEquipment deletes equipment
// @Summary Delete equipment
// @Tags equipment
// @Router /equipment/{id} [delete]
func DeleteEquipment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can delete equipment"})
		return
	}

	if err := equipmentService.Delete(uint(id)); err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetEquipmentStatistics returns equipment statistics
// @Summary Get equipment statistics
// @Tags equipment
// @Router /equipment/statistics [get]
func GetEquipmentStatistics(c *gin.Context) {
	stats, err := equipmentService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.EquipmentStatistics{
		Total:      stats.Total,
		Running:    stats.Running,
		Stopped:    stats.Stopped,
		Maintenance: stats.Maintenance,
		Scrapped:   stats.Scrapped,
	})
}

// =====================================================
// Helper Functions
// =====================================================

func equipmentToResponse(e *model.Equipment) dto.EquipmentResponse {
	r := dto.EquipmentResponse{
		ID:                      e.ID,
		Code:                    e.Code,
		Name:                    e.Name,
		TypeID:                  e.TypeID,
		WorkshopID:              e.WorkshopID,
		QRCode:                  e.QRCode,
		Spec:                    e.Spec,
		PurchaseDate:            e.PurchaseDate,
		Status:                  e.Status,
		DedicatedMaintenanceID:  e.DedicatedMaintenanceID,
		CreatedAt:               e.CreatedAt,
		UpdatedAt:               e.UpdatedAt,
	}

	if e.Type.ID > 0 {
		r.TypeName = e.Type.Name
	}
	if e.Workshop.ID > 0 {
		r.WorkshopName = e.Workshop.Name
		r.FactoryID = e.Workshop.FactoryID
		if e.Workshop.Factory.ID > 0 {
			r.FactoryName = e.Workshop.Factory.Name
		}
	}
	if e.DedicatedMaintenance != nil {
		r.DedicatedMaintenanceName = e.DedicatedMaintenance.Name
	}

	return r
}

func handleServiceError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	if err == service.ErrNotFound {
		status = http.StatusNotFound
	} else if err == service.ErrDuplicateCode {
		status = http.StatusConflict
	}
	c.JSON(status, gin.H{"error": err.Error()})
}
