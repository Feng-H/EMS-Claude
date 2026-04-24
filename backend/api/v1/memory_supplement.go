package v1

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/memory"
)

// ============ Auth (补充) ============

// LogoutMemory 登出 (内存模式)
func LogoutMemory(c *gin.Context) {
	// 内存模式无需清理服务端token
	c.JSON(200, gin.H{"message": "登出成功"})
}

// RefreshTokenMemory 刷新token (内存模式)
func RefreshTokenMemory(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}
	// 简化处理：直接返回成功
	c.JSON(200, gin.H{"message": "token刷新成功"})
}

// ChangePasswordMemory 修改密码 (内存模式)
func ChangePasswordMemory(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	updated := store.UpdateUser(uint(userID), func(user *model.User) {
		if req.OldPassword == "admin123" { // 简化验证
			user.PasswordHash = req.NewPassword
		}
	})

	if updated {
		c.JSON(200, gin.H{"message": "密码修改成功"})
	} else {
		c.JSON(404, gin.H{"error": "用户不存在"})
	}
}

// ApplyAccountMemory 申请账号 (内存模式)
func ApplyAccountMemory(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Role     string `json:"role" binding:"required"`
		Phone    string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	if store.FindUserByUsername(req.Username) != nil {
		c.JSON(409, gin.H{"error": "用户名已存在"})
		return
	}

	now := time.Now()
	newUser := &model.User{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Username:           req.Username,
		PasswordHash:       req.Password,
		Name:               req.Name,
		Role:               model.UserRole(req.Role),
		Phone:              req.Phone,
		IsActive:           true,
		ApprovalStatus:     model.ApprovalStatusPending,
		MustChangePassword: true,
		FirstLogin:         true,
	}

	store.AddUser(newUser.ID, newUser)
	c.JSON(201, gin.H{"message": "申请提交成功", "id": newUser.ID})
}

// ============ Organization CRUD ============

// CreateBaseMemory 创建基地
func CreateBaseMemory(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	base := &model.Base{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Code: req.Code,
		Name: req.Name,
	}

	store.AddBase(base.ID, base)
	c.JSON(201, base)
}

// UpdateBaseMemory 更新基地
func UpdateBaseMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	updated := store.UpdateBase(uint(id), func(base *model.Base) {
		if req.Code != "" {
			base.Code = req.Code
		}
		if req.Name != "" {
			base.Name = req.Name
		}
		base.UpdatedAt = time.Now()
	})

	if updated {
		c.JSON(200, gin.H{"message": "更新成功"})
	} else {
		c.JSON(404, gin.H{"error": "基地不存在"})
	}
}

// DeleteBaseMemory 删除基地
func DeleteBaseMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if store.DeleteBase(uint(id)) {
		c.JSON(200, gin.H{"message": "删除成功"})
	} else {
		c.JSON(404, gin.H{"error": "基地不存在"})
	}
}

// CreateFactoryMemory 创建工厂
func CreateFactoryMemory(c *gin.Context) {
	var req struct {
		BaseID uint   `json:"base_id" binding:"required"`
		Code   string `json:"code" binding:"required"`
		Name   string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	factory := &model.Factory{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		BaseID: req.BaseID,
		Code:   req.Code,
		Name:   req.Name,
	}

	store.AddFactory(factory.ID, factory)
	c.JSON(201, factory)
}

// UpdateFactoryMemory 更新工厂
func UpdateFactoryMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		BaseID uint   `json:"base_id"`
		Code   string `json:"code"`
		Name   string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	updated := store.UpdateFactory(uint(id), func(factory *model.Factory) {
		if req.BaseID != 0 {
			factory.BaseID = req.BaseID
		}
		if req.Code != "" {
			factory.Code = req.Code
		}
		if req.Name != "" {
			factory.Name = req.Name
		}
		factory.UpdatedAt = time.Now()
	})

	if updated {
		c.JSON(200, gin.H{"message": "更新成功"})
	} else {
		c.JSON(404, gin.H{"error": "工厂不存在"})
	}
}

// DeleteFactoryMemory 删除工厂
func DeleteFactoryMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if store.DeleteFactory(uint(id)) {
		c.JSON(200, gin.H{"message": "删除成功"})
	} else {
		c.JSON(404, gin.H{"error": "工厂不存在"})
	}
}

// CreateWorkshopMemory 创建车间
func CreateWorkshopMemory(c *gin.Context) {
	var req struct {
		FactoryID uint   `json:"factory_id" binding:"required"`
		Code      string `json:"code" binding:"required"`
		Name      string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	workshop := &model.Workshop{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		FactoryID: req.FactoryID,
		Code:      req.Code,
		Name:      req.Name,
	}

	store.AddWorkshop(workshop.ID, workshop)
	c.JSON(201, workshop)
}

// UpdateWorkshopMemory 更新车间
func UpdateWorkshopMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		FactoryID uint   `json:"factory_id"`
		Code      string `json:"code"`
		Name      string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	updated := store.UpdateWorkshop(uint(id), func(workshop *model.Workshop) {
		if req.FactoryID != 0 {
			workshop.FactoryID = req.FactoryID
		}
		if req.Code != "" {
			workshop.Code = req.Code
		}
		if req.Name != "" {
			workshop.Name = req.Name
		}
		workshop.UpdatedAt = time.Now()
	})

	if updated {
		c.JSON(200, gin.H{"message": "更新成功"})
	} else {
		c.JSON(404, gin.H{"error": "车间不存在"})
	}
}

// DeleteWorkshopMemory 删除车间
func DeleteWorkshopMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if store.DeleteWorkshop(uint(id)) {
		c.JSON(200, gin.H{"message": "删除成功"})
	} else {
		c.JSON(404, gin.H{"error": "车间不存在"})
	}
}

// ============ Equipment CRUD ============

// GetEquipmentMemory 获取单个设备
func GetEquipmentMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if equipment := store.FindEquipment(uint(id)); equipment != nil {
		c.JSON(200, equipment)
		return
	}
	c.JSON(404, gin.H{"error": "设备不存在"})
}

// CreateEquipmentMemory 创建设备
func CreateEquipmentMemory(c *gin.Context) {
	var req struct {
		Code     string  `json:"code" binding:"required"`
		Name     string  `json:"name" binding:"required"`
		TypeID   uint    `json:"type_id" binding:"required"`
		WorkshopID uint `json:"workshop_id" binding:"required"`
		Spec     string `json:"spec"`
		Status   string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	equipment := &model.Equipment{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Code:       req.Code,
		Name:       req.Name,
		TypeID:     req.TypeID,
		WorkshopID: req.WorkshopID,
		Spec:       req.Spec,
		QRCode:     "QR-" + req.Code,
		Status:     req.Status,
	}

	store.AddEquipment(equipment.ID, equipment)
	c.JSON(201, equipment)
}

// UpdateEquipmentMemory 更新设备
func UpdateEquipmentMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Code       string `json:"code"`
		Name       string `json:"name"`
		TypeID     uint   `json:"type_id"`
		WorkshopID uint   `json:"workshop_id"`
		Spec       string `json:"spec"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	updated := store.UpdateEquipment(uint(id), func(equipment *model.Equipment) {
		if req.Code != "" {
			equipment.Code = req.Code
		}
		if req.Name != "" {
			equipment.Name = req.Name
		}
		if req.TypeID != 0 {
			equipment.TypeID = req.TypeID
		}
		if req.WorkshopID != 0 {
			equipment.WorkshopID = req.WorkshopID
		}
		if req.Spec != "" {
			equipment.Spec = req.Spec
		}
		if req.Status != "" {
			equipment.Status = req.Status
		}
		equipment.UpdatedAt = time.Now()
	})

	if updated {
		c.JSON(200, gin.H{"message": "更新成功"})
	} else {
		c.JSON(404, gin.H{"error": "设备不存在"})
	}
}

// DeleteEquipmentMemory 删除设备
func DeleteEquipmentMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if store.DeleteEquipment(uint(id)) {
		c.JSON(200, gin.H{"message": "删除成功"})
	} else {
		c.JSON(404, gin.H{"error": "设备不存在"})
	}
}

// CreateEquipmentTypeMemory 创建设备类型
func CreateEquipmentTypeMemory(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Category string `json:"category"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	equipmentType := &model.EquipmentType{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:     req.Name,
		Category: req.Category,
	}

	store.AddEquipmentType(equipmentType.ID, equipmentType)
	c.JSON(201, equipmentType)
}

// UpdateEquipmentTypeMemory 更新设备类型
func UpdateEquipmentTypeMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	updated := store.UpdateEquipmentType(uint(id), func(et *model.EquipmentType) {
		if req.Name != "" {
			et.Name = req.Name
		}
		if req.Category != "" {
			et.Category = req.Category
		}
		et.UpdatedAt = time.Now()
	})

	if updated {
		c.JSON(200, gin.H{"message": "更新成功"})
	} else {
		c.JSON(404, gin.H{"error": "设备类型不存在"})
	}
}

// DeleteEquipmentTypeMemory 删除设备类型
func DeleteEquipmentTypeMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if store.DeleteEquipmentType(uint(id)) {
		c.JSON(200, gin.H{"message": "删除成功"})
	} else {
		c.JSON(404, gin.H{"error": "设备类型不存在"})
	}
}

// ============ Inspection 补充 ============

// GetInspectionTemplateMemory 获取点检模板详情
func GetInspectionTemplateMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if template := store.FindInspectionTemplate(uint(id)); template != nil {
		items := store.GetInspectionItemsByTemplate(template.ID)
		c.JSON(200, gin.H{
			"id":                 template.ID,
			"name":               template.Name,
			"equipment_type_id":  template.EquipmentTypeID,
			"item_count":         len(items),
			"items":              items,
			"created_at":         template.CreatedAt.Format("2006-01-02 15:04:05"),
		})
		return
	}
	c.JSON(404, gin.H{"error": "模板不存在"})
}

// CreateInspectionTemplateMemory 创建点检模板
func CreateInspectionTemplateMemory(c *gin.Context) {
	var req struct {
		Name            string `json:"name" binding:"required"`
		EquipmentTypeID uint   `json:"equipment_type_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	template := &model.InspectionTemplate{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:            req.Name,
		EquipmentTypeID: req.EquipmentTypeID,
	}

	store.AddInspectionTemplate(template.ID, template)
	c.JSON(201, template)
}

// CreateInspectionItemMemory 创建点检项目
func CreateInspectionItemMemory(c *gin.Context) {
	var req struct {
		TemplateID uint   `json:"template_id" binding:"required"`
		Name       string `json:"name" binding:"required"`
		Method     string `json:"method"`
		Criteria   string `json:"criteria"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	item := &model.InspectionItem{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		TemplateID: req.TemplateID,
		Name:       req.Name,
		Method:     req.Method,
		Criteria:   req.Criteria,
	}

	store.AddInspectionItem(item.ID, item)
	c.JSON(201, item)
}

// GetInspectionTaskMemory 获取点检任务详情
func GetInspectionTaskMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if task := store.FindInspectionTask(uint(id)); task != nil {
		c.JSON(200, task)
		return
	}
	c.JSON(404, gin.H{"error": "任务不存在"})
}

// StartInspectionMemory 开始点检
func StartInspectionMemory(c *gin.Context) {
	var req struct {
		EquipmentID uint   `json:"equipment_id" binding:"required"`
		QRCode      string `json:"qr_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()

	// 查找模板
	templateID := uint(0)
	for _, t := range store.InspectionTemplates {
		templateID = t.ID
		break
	}

	task := &model.InspectionTask{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		EquipmentID:    req.EquipmentID,
		TemplateID:     templateID,
		ScheduledDate:  now,
		Status:         model.InspectionInProgress,
	}
	task.StartedAt = &now

	store.AddInspectionTask(task.ID, task)

	// 获取检查项目
	items := store.GetInspectionItemsByTemplate(templateID)

	c.JSON(200, gin.H{
		"task_id":    task.ID,
		"equipment":  store.FindEquipment(req.EquipmentID),
		"items":      items,
		"started_at": now.Format("2006-01-02 15:04:05"),
	})
}

// CompleteInspectionMemory 完成点检
func CompleteInspectionMemory(c *gin.Context) {
	var req struct {
		TaskID  uint `json:"task_id" binding:"required"`
		Records []struct {
			ItemID  uint   `json:"item_id" binding:"required"`
			Result string `json:"result" binding:"required"`
			Remark string `json:"remark"`
		}
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()

	// 更新任务状态
	store.UpdateInspectionTask(req.TaskID, func(task *model.InspectionTask) {
		task.Status = model.InspectionCompleted
		task.CompletedAt = &now
		task.UpdatedAt = now
	})

	// 创建记录
	ngCount := 0
	for _, record := range req.Records {
		if record.Result == "NG" {
			ngCount++
		}
		rec := &model.InspectionRecord{
			BaseModel: model.BaseModel{
				ID:        store.NextID(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			TaskID:  req.TaskID,
			ItemID:  record.ItemID,
			Result:  record.Result,
			Remark:  record.Remark,
		}
		store.AddInspectionRecord(rec.ID, rec)
	}

	c.JSON(200, gin.H{
		"task_id":      req.TaskID,
		"completed_at": now.Format("2006-01-02 15:04:05"),
		"total_count":  len(req.Records),
		"ok_count":     len(req.Records) - ngCount,
		"ng_count":     ngCount,
	})
}

// ============ Repair 补充 ============

// GetRepairOrderMemory 获取维修工单详情
func GetRepairOrderMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if order := store.FindRepairOrder(uint(id)); order != nil {
		c.JSON(200, order)
		return
	}
	c.JSON(404, gin.H{"error": "工单不存在"})
}

// AssignRepairOrderMemory 派发维修工单
func AssignRepairOrderMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		AssignTo uint `json:"assign_to" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	updated := store.UpdateRepairOrder(uint(id), func(order *model.RepairOrder) {
		order.AssignedTo = &req.AssignTo
		order.Status = model.RepairAssigned
		order.UpdatedAt = time.Now()
	})

	if updated {
		c.JSON(200, gin.H{"message": "派单成功"})
	} else {
		c.JSON(404, gin.H{"error": "工单不存在"})
	}
}

// StartRepairMemory 开始维修
func StartRepairMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	now := time.Now()
	updated := store.UpdateRepairOrder(uint(id), func(order *model.RepairOrder) {
		order.Status = model.RepairInProgress
		order.StartedAt = &now
		order.UpdatedAt = now
	})

	if updated {
		c.JSON(200, gin.H{"message": "开始维修"})
	} else {
		c.JSON(404, gin.H{"error": "工单不存在"})
	}
}

// UpdateRepairMemory 更新维修进度
func UpdateRepairMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Solution    string `json:"solution"`
		SpareParts  string `json:"spare_parts"`
		ActualHours float64 `json:"actual_hours"`
		NextStatus  string `json:"next_status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	updated := store.UpdateRepairOrder(uint(id), func(order *model.RepairOrder) {
		if req.Solution != "" {
			order.Solution = req.Solution
		}
		// SpareParts 和 ActualHours 字段不存在于模型中，暂时忽略
		// 实际项目中可以存储在 Solution 或 Logs 中
		if req.NextStatus == "testing" {
			order.Status = model.RepairTesting
		}
		order.UpdatedAt = time.Now()
	})

	if updated {
		c.JSON(200, gin.H{"message": "更新成功"})
	} else {
		c.JSON(404, gin.H{"error": "工单不存在"})
	}
}

// ConfirmRepairMemory 确认维修
func ConfirmRepairMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Accepted bool   `json:"accepted" binding:"required"`
		Comment  string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	updated := store.UpdateRepairOrder(uint(id), func(order *model.RepairOrder) {
		if req.Accepted {
			order.Status = model.RepairConfirmed
		}
		order.ConfirmedAt = &now
		order.UpdatedAt = now
	})

	if updated {
		c.JSON(200, gin.H{"message": "确认成功"})
	} else {
		c.JSON(404, gin.H{"error": "工单不存在"})
	}
}

// AuditRepairMemory 审核维修
func AuditRepairMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Approved    bool    `json:"approved" binding:"required"`
		Comment     string  `json:"comment"`
		ActualHours float64 `json:"actual_hours"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	updated := store.UpdateRepairOrder(uint(id), func(order *model.RepairOrder) {
		if req.Approved {
			order.Status = model.RepairAudited
		}
		// ActualHours 字段不存在于模型中，暂时忽略
		order.AuditedAt = &now
		order.UpdatedAt = now
	})

	if updated {
		c.JSON(200, gin.H{"message": "审核成功"})
	} else {
		c.JSON(404, gin.H{"error": "工单不存在"})
	}
}

// GetMyRepairStatisticsMemory 获取我的维修统计
func GetMyRepairStatisticsMemory(c *gin.Context) {
	store := memory.GetStore()

	pendingCount := 0
	inProgressCount := 0
	completedCount := 0

	for _, order := range store.RepairOrders {
		switch order.Status {
		case model.RepairPending:
			pendingCount++
		case model.RepairInProgress, model.RepairTesting:
			inProgressCount++
		case model.RepairAudited, model.RepairClosed:
			completedCount++
		}
	}

	c.JSON(200, gin.H{
		"pending_count":     pendingCount,
		"in_progress_count": inProgressCount,
		"completed_count":   completedCount,
		"today_completed":   0,
	})
}

// ============ Maintenance 补充 ============

// CreateMaintenancePlanMemory 创建保养计划
func CreateMaintenancePlanMemory(c *gin.Context) {
	var req struct {
		Name            string `json:"name" binding:"required"`
		EquipmentTypeID uint   `json:"equipment_type_id" binding:"required"`
		Level           int    `json:"level" binding:"required"`
		CycleDays       int    `json:"cycle_days" binding:"required"`
		FlexibleDays    int    `json:"flexible_days"`
		WorkHours       float64 `json:"work_hours"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	plan := &model.MaintenancePlan{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:            req.Name,
		EquipmentTypeID: req.EquipmentTypeID,
		Level:           req.Level,
		CycleDays:       req.CycleDays,
		FlexibleDays:    req.FlexibleDays,
		WorkHours:       req.WorkHours,
	}

	store.AddMaintenancePlan(plan.ID, plan)
	c.JSON(201, plan)
}

// CreateMaintenanceItemMemory 创建保养项目
func CreateMaintenanceItemMemory(c *gin.Context) {
	var req struct {
		PlanID uint   `json:"plan_id" binding:"required"`
		Name   string `json:"name" binding:"required"`
		Method string `json:"method"`
		Criteria string `json:"criteria"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	item := &model.MaintenancePlanItem{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		PlanID: req.PlanID,
		Name:   req.Name,
		Method: req.Method,
		Criteria: req.Criteria,
	}

	store.AddMaintenanceItem(item.ID, item)
	c.JSON(201, item)
}

// GenerateMaintenanceTasksMemory 生成保养任务
func GenerateMaintenanceTasksMemory(c *gin.Context) {
	var req struct {
		PlanID       uint   `json:"plan_id" binding:"required"`
		EquipmentIDs []uint `json:"equipment_ids" binding:"required"`
		Date         string `json:"date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	var taskIDs []uint

	for _, equipmentID := range req.EquipmentIDs {
		task := &model.MaintenanceTask{
			BaseModel: model.BaseModel{
				ID:        store.NextID(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			PlanID:        req.PlanID,
			EquipmentID:   equipmentID,
			ScheduledDate: req.Date,
			DueDate:       req.Date,
			Status:        model.MaintenancePending,
		}
		store.AddMaintenanceTask(task.ID, task)
		taskIDs = append(taskIDs, task.ID)
	}

	c.JSON(200, gin.H{
		"created_count": len(taskIDs),
		"task_ids":      taskIDs,
		"errors":        []string{},
	})
}

// GetMaintenanceTaskMemory 获取保养任务详情
func GetMaintenanceTaskMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if task := store.FindMaintenanceTask(uint(id)); task != nil {
		c.JSON(200, task)
		return
	}
	c.JSON(404, gin.H{"error": "任务不存在"})
}

// StartMaintenanceMemory 开始保养
func StartMaintenanceMemory(c *gin.Context) {
	var req struct {
		TaskID uint `json:"task_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	updated := store.UpdateMaintenanceTask(req.TaskID, func(task *model.MaintenanceTask) {
		task.Status = model.MaintenanceInProgress
		task.StartedAt = &now
		task.UpdatedAt = now
	})

	if updated {
		c.JSON(200, gin.H{"message": "开始保养"})
	} else {
		c.JSON(404, gin.H{"error": "任务不存在"})
	}
}

// CompleteMaintenanceMemory 完成保养
func CompleteMaintenanceMemory(c *gin.Context) {
	var req struct {
		TaskID     uint `json:"task_id" binding:"required"`
		Records    []struct {
			ItemID  uint   `json:"item_id"`
			Result  string `json:"result"`
			Remark  string `json:"remark"`
		}
		ActualHours float64 `json:"actual_hours"`
		Remark      string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()

	// 更新任务状态
	store.UpdateMaintenanceTask(req.TaskID, func(task *model.MaintenanceTask) {
		task.Status = model.MaintenanceCompleted
		task.CompletedAt = &now
		task.ActualHours = req.ActualHours
		task.Remark = req.Remark
		task.UpdatedAt = now
	})

	// 创建记录
	ngCount := 0
	var ngItemIDs []uint
	for _, record := range req.Records {
		if record.Result == "NG" {
			ngCount++
			ngItemIDs = append(ngItemIDs, record.ItemID)
		}
		rec := &model.MaintenanceRecord{
			BaseModel: model.BaseModel{
				ID:        store.NextID(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			TaskID: req.TaskID,
			ItemID: record.ItemID,
			Result: record.Result,
			Remark: record.Remark,
		}
		store.AddMaintenanceRecord(rec.ID, rec)
	}

	c.JSON(200, gin.H{
		"task_id":       req.TaskID,
		"completed_at":  now.Format("2006-01-02 15:04:05"),
		"total_count":   len(req.Records),
		"ok_count":      len(req.Records) - ngCount,
		"ng_count":      ngCount,
		"ng_item_ids":   ngItemIDs,
	})
}

// ============ Spare Parts 补充 ============

// CreateSparePartMemory 创建备件
func CreateSparePartMemory(c *gin.Context) {
	var req struct {
		Code          string `json:"code" binding:"required"`
		Name          string `json:"name" binding:"required"`
		Specification string `json:"specification"`
		Unit          string `json:"unit"`
		SafetyStock   int    `json:"safety_stock"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()
	part := &model.SparePart{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Code:          req.Code,
		Name:          req.Name,
		Specification: req.Specification,
		Unit:          req.Unit,
		SafetyStock:   req.SafetyStock,
	}

	store.AddSparePart(part.ID, part)
	c.JSON(201, part)
}

// UpdateSparePartMemory 更新备件
func UpdateSparePartMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Code          string `json:"code"`
		Name          string `json:"name"`
		Specification string `json:"specification"`
		Unit          string `json:"unit"`
		SafetyStock   int    `json:"safety_stock"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	updated := store.UpdateSparePart(uint(id), func(part *model.SparePart) {
		if req.Code != "" {
			part.Code = req.Code
		}
		if req.Name != "" {
			part.Name = req.Name
		}
		if req.Specification != "" {
			part.Specification = req.Specification
		}
		if req.Unit != "" {
			part.Unit = req.Unit
		}
		part.SafetyStock = req.SafetyStock
		part.UpdatedAt = time.Now()
	})

	if updated {
		c.JSON(200, gin.H{"message": "更新成功"})
	} else {
		c.JSON(404, gin.H{"error": "备件不存在"})
	}
}

// DeleteSparePartMemory 删除备件
func DeleteSparePartMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if store.DeleteSparePart(uint(id)) {
		c.JSON(200, gin.H{"message": "删除成功"})
	} else {
		c.JSON(404, gin.H{"error": "备件不存在"})
	}
}

// StockInMemory 备件入库
func StockInMemory(c *gin.Context) {
	var req struct {
		SparePartID uint   `json:"spare_part_id" binding:"required"`
		FactoryID   uint   `json:"factory_id" binding:"required"`
		Quantity    int    `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()

	// 查找现有库存
	var inventory *model.SparePartInventory
	for _, inv := range store.SparePartInventory {
		if inv.SparePartID == req.SparePartID && inv.FactoryID == req.FactoryID {
			inventory = inv
			break
		}
	}

	if inventory != nil {
		inventory.Quantity += req.Quantity
		inventory.UpdatedAt = now
	} else {
		inventory = &model.SparePartInventory{
			BaseModel: model.BaseModel{
				ID:        store.NextID(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			SparePartID: req.SparePartID,
			FactoryID:   req.FactoryID,
			Quantity:    req.Quantity,
		}
		store.AddSparePartInventory(inventory.ID, inventory)
	}

	c.JSON(200, gin.H{"message": "入库成功"})
}

// StockOutMemory 备件出库
func StockOutMemory(c *gin.Context) {
	var req struct {
		SparePartID uint `json:"spare_part_id" binding:"required"`
		FactoryID   uint `json:"factory_id" binding:"required"`
		Quantity    int  `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()

	// 查找库存
	for _, inv := range store.SparePartInventory {
		if inv.SparePartID == req.SparePartID && inv.FactoryID == req.FactoryID {
			if inv.Quantity < req.Quantity {
				c.JSON(400, gin.H{"error": "库存不足"})
				return
			}
			inv.Quantity -= req.Quantity
			inv.UpdatedAt = time.Now()
			c.JSON(200, gin.H{"message": "出库成功"})
			return
		}
	}

	c.JSON(404, gin.H{"error": "库存不存在"})
}

// GetLowStockAlertsMemory 获取低库存预警
func GetLowStockAlertsMemory(c *gin.Context) {
	store := memory.GetStore()
	var alerts []gin.H

	for _, inv := range store.SparePartInventory {
		part := store.FindSparePart(inv.SparePartID)
		if part != nil && inv.Quantity < part.SafetyStock {
			alerts = append(alerts, gin.H{
				"spare_part_id":   part.ID,
				"spare_part_code": part.Code,
				"spare_part_name": part.Name,
				"factory_id":      inv.FactoryID,
				"current_stock":   inv.Quantity,
				"safety_stock":    part.SafetyStock,
				"shortage":        part.SafetyStock - inv.Quantity,
			})
		}
	}

	c.JSON(200, alerts)
}

// GetConsumptionsMemory 获取消耗记录
func GetConsumptionsMemory(c *gin.Context) {
	store := memory.GetStore()
	var consumptions []gin.H

	for _, cons := range store.SparePartConsumption {
		part := store.FindSparePart(cons.SparePartID)
		if part != nil {
			consumptions = append(consumptions, gin.H{
				"id":              cons.ID,
				"spare_part_id":   cons.SparePartID,
				"spare_part_code": part.Code,
				"spare_part_name": part.Name,
				"quantity":        cons.Quantity,
				"user_id":         cons.UserID,
				"created_at":      cons.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	c.JSON(200, gin.H{
		"data":  consumptions,
		"total": len(consumptions),
	})
}

// CreateConsumptionMemory 创建消耗记录
func CreateConsumptionMemory(c *gin.Context) {
	var req struct {
		SparePartID uint `json:"spare_part_id" binding:"required"`
		Quantity    int  `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	now := time.Now()

	consumption := &model.SparePartConsumption{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		SparePartID: req.SparePartID,
		Quantity:    req.Quantity,
	}

	store.AddSparePartConsumption(consumption.ID, consumption)
	c.JSON(201, consumption)
}

// ============ Knowledge 补充 ============

// GetKnowledgeArticleMemory 获取知识库文章详情
func GetKnowledgeArticleMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if article := store.FindKnowledgeArticle(uint(id)); article != nil {
		c.JSON(200, article)
		return
	}
	c.JSON(404, gin.H{"error": "文章不存在"})
}

// CreateKnowledgeArticleMemory 创建知识库文章
func CreateKnowledgeArticleMemory(c *gin.Context) {
	var req struct {
		Title            string   `json:"title" binding:"required"`
		EquipmentTypeID  *uint    `json:"equipment_type_id"`
		FaultPhenomenon  string   `json:"fault_phenomenon"`
		CauseAnalysis    string   `json:"cause_analysis"`
		Solution         string   `json:"solution" binding:"required"`
		SourceType       string   `json:"source_type"`
		Tags             []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	store := memory.GetStore()
	now := time.Now()

	article := &model.KnowledgeArticle{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Title:            req.Title,
		EquipmentTypeID:  req.EquipmentTypeID,
		FaultPhenomenon:  req.FaultPhenomenon,
		CauseAnalysis:    req.CauseAnalysis,
		Solution:         req.Solution,
		SourceType:       req.SourceType,
		Tags:             req.Tags,
		CreatedBy:        uint(userID),
	}

	store.AddKnowledgeArticle(article.ID, article)
	c.JSON(201, article)
}

// UpdateKnowledgeArticleMemory 更新知识库文章
func UpdateKnowledgeArticleMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Title           string   `json:"title"`
		EquipmentTypeID *uint    `json:"equipment_type_id"`
		FaultPhenomenon string   `json:"fault_phenomenon"`
		CauseAnalysis   string   `json:"cause_analysis"`
		Solution        string   `json:"solution"`
		Tags            []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	updated := store.UpdateKnowledgeArticle(uint(id), func(article *model.KnowledgeArticle) {
		if req.Title != "" {
			article.Title = req.Title
		}
		if req.EquipmentTypeID != nil {
			article.EquipmentTypeID = req.EquipmentTypeID
		}
		if req.FaultPhenomenon != "" {
			article.FaultPhenomenon = req.FaultPhenomenon
		}
		if req.CauseAnalysis != "" {
			article.CauseAnalysis = req.CauseAnalysis
		}
		if req.Solution != "" {
			article.Solution = req.Solution
		}
		if req.Tags != nil {
			article.Tags = req.Tags
		}
		article.UpdatedAt = time.Now()
	})

	if updated {
		c.JSON(200, gin.H{"message": "更新成功"})
	} else {
		c.JSON(404, gin.H{"error": "文章不存在"})
	}
}

// DeleteKnowledgeArticleMemory 删除知识库文章
func DeleteKnowledgeArticleMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	store := memory.GetStore()
	if store.DeleteKnowledgeArticle(uint(id)) {
		c.JSON(200, gin.H{"message": "删除成功"})
	} else {
		c.JSON(404, gin.H{"error": "文章不存在"})
	}
}

// ConvertFromRepairMemory 从维修工单转换
func ConvertFromRepairMemory(c *gin.Context) {
	var req struct {
		OrderID         uint   `json:"order_id" binding:"required"`
		Title           string `json:"title" binding:"required"`
		FaultPhenomenon string `json:"fault_phenomenon" binding:"required"`
		CauseAnalysis   string `json:"cause_analysis" binding:"required"`
		Tags            []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	order := store.FindRepairOrder(req.OrderID)
	if order == nil {
		c.JSON(404, gin.H{"error": "工单不存在"})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)

	equipment := store.FindEquipment(order.EquipmentID)
	var equipmentTypeID *uint
	if equipment != nil {
		equipmentTypeID = &equipment.TypeID
	}

	now := time.Now()
	article := &model.KnowledgeArticle{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Title:            req.Title,
		EquipmentTypeID:  equipmentTypeID,
		FaultPhenomenon:  req.FaultPhenomenon,
		CauseAnalysis:    req.CauseAnalysis,
		Solution:         order.Solution,
		SourceType:       "repair",
		SourceID:         &req.OrderID,
		Tags:             req.Tags,
		CreatedBy:        uint(userID),
	}

	store.AddKnowledgeArticle(article.ID, article)
	c.JSON(201, article)
}
