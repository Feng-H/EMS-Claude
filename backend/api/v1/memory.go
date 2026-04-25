package v1

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/jwt"
	"github.com/ems/backend/pkg/memory"
)

var (
	store *memory.Store
)

// InitMemory 初始化内存模式
func InitMemory() {
	store = memory.GetStore()
	store.InitMockData()
}

// ============ 认证相关 ============

// Login 登录 (内存模式)
func LoginMemory(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()
	for _, user := range store.Users {
		if user.Username == req.Username && user.IsActive {
			// 密码验证在内存存储初始化时已设置，使用bcrypt
			if user.PasswordHash != "" {
				token, err := jwt.GenerateToken(user.ID, string(user.Role))
				if err != nil {
					c.JSON(500, gin.H{"error": "生成token失败"})
					return
				}
				c.JSON(200, gin.H{
					"token": token,
					"expire_at": time.Now().Add(24 * time.Hour).Unix(),
					"user_info": gin.H{
						"id":       user.ID,
						"username": user.Username,
						"name":     user.Name,
						"role":     string(user.Role),
						"factory_id": user.FactoryID,
					},
					"must_change_password": user.MustChangePassword,
				})
				return
			}
		}
	}
	c.JSON(401, gin.H{"error": "用户名或密码错误"})
}

// GetCurrentUser 获取当前用户 (内存模式)
func GetCurrentUserMemory(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(uint); ok {
				store := memory.GetStore()
				if user := store.FindUser(uid); user != nil {
					c.JSON(200, gin.H{
						"id":                  user.ID,
						"username":            user.Username,
						"name":                user.Name,
						"role":                string(user.Role),
						"factory_id":          user.FactoryID,
						"phone":               user.Phone,
						"is_active":           user.IsActive,
						"approval_status":     string(user.ApprovalStatus),
						"must_change_password": user.MustChangePassword,
					})
					return
				}
			}
		}
		c.JSON(404, gin.H{"error": "用户不存在"})
		return
	}

	userID, _ := strconv.ParseUint(userIDStr, 10, 32)
	store := memory.GetStore()
	if user := store.FindUser(uint(userID)); user != nil {
		c.JSON(200, gin.H{
			"id":                  user.ID,
			"username":            user.Username,
			"name":                user.Name,
			"role":                string(user.Role),
			"factory_id":          user.FactoryID,
			"phone":               user.Phone,
			"is_active":           user.IsActive,
			"approval_status":     string(user.ApprovalStatus),
			"must_change_password": user.MustChangePassword,
		})
		return
	}
	c.JSON(404, gin.H{"error": "用户不存在"})
}

// ============ 用户管理 (内存模式) ============

// GetUsers 获取用户列表 (内存模式)
func GetUsersMemory(c *gin.Context) {
	store := memory.GetStore()
	users := store.GetUsers()
	var result []gin.H
	for _, u := range users {
		result = append(result, gin.H{
			"id":                  u.ID,
			"username":            u.Username,
			"name":                u.Name,
			"role":                string(u.Role),
			"phone":               u.Phone,
			"is_active":           u.IsActive,
			"approval_status":     string(u.ApprovalStatus),
			"must_change_password": u.MustChangePassword,
			"factory_id":          u.FactoryID,
			"created_at":          u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	c.JSON(200, gin.H{"items": result, "total": len(result)})
}

// CreateUser 创建用户 (内存模式)
func CreateUserMemory(c *gin.Context) {
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

	// 检查用户名是否已存在
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
		PasswordHash:       req.Password, // 简化处理，实际应该hash
		Name:               req.Name,
		Role:               model.UserRole(req.Role),
		Phone:              req.Phone,
		IsActive:           true,
		ApprovalStatus:     model.ApprovalStatusApproved,
		MustChangePassword: true,
		FirstLogin:         true,
	}

	store.AddUser(newUser.ID, newUser)

	c.JSON(201, gin.H{"id": newUser.ID, "message": "用户创建成功"})
}

// UpdateUser 更新用户 (内存模式)
func UpdateUserMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Name     string `json:"name"`
		Role     string `json:"role"`
		Phone    string `json:"phone"`
		IsActive *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()

	updated := store.UpdateUser(uint(id), func(user *model.User) {
		if req.Name != "" {
			user.Name = req.Name
		}
		if req.Role != "" {
			user.Role = model.UserRole(req.Role)
		}
		if req.Phone != "" {
			user.Phone = req.Phone
		}
		if req.IsActive != nil {
			user.IsActive = *req.IsActive
		}
		user.UpdatedAt = time.Now()
	})

	if updated {
		c.JSON(200, gin.H{"message": "用户更新成功"})
	} else {
		c.JSON(404, gin.H{"error": "用户不存在"})
	}
}

// GetPendingApplications 获取待审核申请 (内存模式)
func GetPendingApplicationsMemory(c *gin.Context) {
	store := memory.GetStore()
	users := store.GetUsers()
	var result []gin.H
	for _, u := range users {
		if u.ApprovalStatus == model.ApprovalStatusPending {
			result = append(result, gin.H{
				"id":              u.ID,
				"username":        u.Username,
				"name":            u.Name,
				"role":            string(u.Role),
				"phone":           u.Phone,
				"is_active":       u.IsActive,
				"approval_status": string(u.ApprovalStatus),
				"created_at":      u.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}
	c.JSON(200, gin.H{"items": result, "total": len(result)})
}

// ApproveApplication 审核用户申请 (内存模式)
func ApproveApplicationMemory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Approve bool   `json:"approve"`
		Reason   string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	store := memory.GetStore()

	updated := store.UpdateUser(uint(id), func(user *model.User) {
		if req.Approve {
			user.ApprovalStatus = model.ApprovalStatusApproved
		} else {
			user.ApprovalStatus = model.ApprovalStatusRejected
			user.RejectionReason = req.Reason
		}
	})

	if updated {
		c.JSON(200, gin.H{"message": "审核完成"})
	} else {
		c.JSON(404, gin.H{"error": "用户不存在"})
	}
}

// ============ 组织架构 ============

// ListBases 获取基地列表 (内存模式)
func ListBasesMemory(c *gin.Context) {
	store := memory.GetStore()
	var bases []*model.Base
	for _, b := range store.Bases {
		bases = append(bases, b)
	}
	c.JSON(200, gin.H{"items": bases, "total": len(bases)})
}

// ListFactories 获取工厂列表 (内存模式)
func ListFactoriesMemory(c *gin.Context) {
	store := memory.GetStore()
	var factories []*model.Factory
	for _, f := range store.Factories {
		factories = append(factories, f)
	}
	c.JSON(200, gin.H{"items": factories, "total": len(factories)})
}

// ListWorkshops 获取车间列表 (内存模式)
func ListWorkshopsMemory(c *gin.Context) {
	store := memory.GetStore()
	var workshops []*model.Workshop
	for _, w := range store.Workshops {
		workshops = append(workshops, w)
	}
	c.JSON(200, gin.H{"items": workshops, "total": len(workshops)})
}

// ============ 设备相关 ============

// ListEquipment 获取设备列表 (内存模式)
func ListEquipmentMemory(c *gin.Context) {
	store := memory.GetStore()
	var equipments []*model.Equipment
	for _, e := range store.Equipment {
		equipments = append(equipments, e)
	}
	c.JSON(200, gin.H{
		"items": equipments,
		"total": len(equipments),
	})
}

// GetEquipmentStatistics 获取设备统计 (内存模式)
func GetEquipmentStatisticsMemory(c *gin.Context) {
	store := memory.GetStore()
	total := len(store.Equipment)

	running := 0
	stopped := 0
	maintenance := 0
	scrapped := 0
	for _, e := range store.Equipment {
		switch e.Status {
		case "running":
			running++
		case "stopped":
			stopped++
		case "maintenance":
			maintenance++
		case "scrapped":
			scrapped++
		default:
			running++
		}
	}

	c.JSON(200, gin.H{
		"total":       total,
		"running":     running,
		"stopped":     stopped,
		"maintenance": maintenance,
		"scrapped":    scrapped,
	})
}

// ListEquipmentTypes 获取设备类型列表 (内存模式)
func ListEquipmentTypesMemory(c *gin.Context) {
	store := memory.GetStore()
	var types []*model.EquipmentType
	for _, t := range store.EquipmentTypes {
		types = append(types, t)
	}
	c.JSON(200, gin.H{"items": types, "total": len(types)})
}

// GetEquipmentByQRCode 通过二维码获取设备 (内存模式)
func GetEquipmentByQRCodeMemory(c *gin.Context) {
	code := c.Param("code")
	store := memory.GetStore()
	for _, e := range store.Equipment {
		if e.QRCode == code || e.Code == code {
			c.JSON(200, e)
			return
		}
	}
	c.JSON(404, gin.H{"error": "设备不存在"})
}

// ============ 点检相关 ============

// ListInspectionTemplates 获取点检模板列表 (内存模式)
func ListInspectionTemplatesMemory(c *gin.Context) {
	store := memory.GetStore()
	var templates []*model.InspectionTemplate
	for _, t := range store.InspectionTemplates {
		templates = append(templates, t)
	}
	c.JSON(200, gin.H{"items": templates, "total": len(templates)})
}

// ListInspectionTasks 获取点检任务列表 (内存模式)
func ListInspectionTasksMemory(c *gin.Context) {
	store := memory.GetStore()
	var tasks []*model.InspectionTask
	for _, t := range store.InspectionTasks {
		tasks = append(tasks, t)
	}
	c.JSON(200, gin.H{"items": tasks, "total": len(tasks)})
}

// GetMyTasks 获取我的点检任务 (内存模式)
func GetMyTasksMemory(c *gin.Context) {
	store := memory.GetStore()
	var tasks []*model.InspectionTask
	for _, t := range store.InspectionTasks {
		if t.Status == model.InspectionPending || t.Status == model.InspectionInProgress {
			tasks = append(tasks, t)
		}
	}
	c.JSON(200, gin.H{"items": tasks, "total": len(tasks)})
}

// GetMyTaskStatistics 获取我的点检统计 (内存模式)
func GetMyTaskStatisticsMemory(c *gin.Context) {
	c.JSON(200, gin.H{
		"pending_count":      2,
		"in_progress_count":  0,
		"today_tasks":        5,
	})
}

// GetInspectionStatistics 获取点检统计 (内存模式)
func GetInspectionStatisticsMemory(c *gin.Context) {
	store := memory.GetStore()
	totalTasks := len(store.InspectionTasks)
	pending := 0
	inProgress := 0
	completed := 0
	overdue := 0

	for _, t := range store.InspectionTasks {
		switch t.Status {
		case model.InspectionPending:
			pending++
		case model.InspectionInProgress:
			inProgress++
		case model.InspectionCompleted:
			completed++
		}
	}

	completionRate := 0.0
	if totalTasks > 0 {
		completionRate = float64(completed) / float64(totalTasks) * 100
	}

	c.JSON(200, gin.H{
		"total_tasks":        totalTasks,
		"pending_tasks":      pending,
		"in_progress_tasks":  inProgress,
		"completed_tasks":    completed,
		"overdue_tasks":      overdue,
		"today_completed":    completed,
		"completion_rate":    completionRate,
	})
}

// ============ 维修相关 ============

// ListRepairOrders 获取维修工单列表 (内存模式)
func ListRepairOrdersMemory(c *gin.Context) {
	store := memory.GetStore()
	var orders []*model.RepairOrder
	for _, o := range store.RepairOrders {
		orders = append(orders, o)
	}
	c.JSON(200, gin.H{
		"items": orders,
		"total": len(orders),
	})
}

// GetMyRepairTasks 获取我的维修任务 (内存模式)
func GetMyRepairTasksMemory(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)
	store := memory.GetStore()
	var tasks []*model.RepairOrder
	for _, o := range store.RepairOrders {
		if o.AssignedTo != nil && *o.AssignedTo == uint(userID) {
			if o.Status != model.RepairAudited && o.Status != model.RepairClosed {
				tasks = append(tasks, o)
			}
		}
	}
	c.JSON(200, gin.H{"items": tasks, "total": len(tasks)})
}

// GetRepairStatistics 获取维修统计 (内存模式)
func GetRepairStatisticsMemory(c *gin.Context) {
	store := memory.GetStore()
	total := len(store.RepairOrders)
	pending := 0
	completed := 0
	avgMTTR := 120.0 // 分钟

	for _, o := range store.RepairOrders {
		if o.Status == model.RepairPending || o.Status == model.RepairAssigned {
			pending++
		}
		if o.Status == model.RepairAudited || o.Status == model.RepairClosed {
			completed++
		}
	}

	c.JSON(200, gin.H{
		"total":    total,
		"pending":  pending,
		"completed": completed,
		"avg_mttr": avgMTTR,
	})
}

// CreateRepairOrder 创建维修工单 (内存模式)
func CreateRepairOrderMemory(c *gin.Context) {
	var req struct {
		EquipmentID uint   `json:"equipment_id" binding:"required"`
		FaultDescription string `json:"fault_description" binding:"required"`
		Priority    int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求参数"})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 32)
	store := memory.GetStore()

	now := time.Now()
	order := &model.RepairOrder{
		BaseModel: model.BaseModel{
			ID:        store.NextID(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		EquipmentID:       req.EquipmentID,
		FaultDescription:  req.FaultDescription,
		ReporterID:        uint(userID),
		Status:            model.RepairPending,
		Priority:          req.Priority,
	}

	store.AddRepairOrder(order)

	c.JSON(201, order)
}

// ============ 保养相关 ============

// ListMaintenancePlans 获取保养计划列表 (内存模式)
func ListMaintenancePlansMemory(c *gin.Context) {
	store := memory.GetStore()
	var plans []*model.MaintenancePlan
	for _, p := range store.MaintenancePlans {
		plans = append(plans, p)
	}
	c.JSON(200, gin.H{"items": plans, "total": len(plans)})
}

// ListMaintenanceTasks 获取保养任务列表 (内存模式)
func ListMaintenanceTasksMemory(c *gin.Context) {
	store := memory.GetStore()
	var tasks []*model.MaintenanceTask
	for _, t := range store.MaintenanceTasks {
		tasks = append(tasks, t)
	}
	c.JSON(200, gin.H{"items": tasks, "total": len(tasks)})
}

// GetMyMaintenanceTasks 获取我的保养任务 (内存模式)
func GetMyMaintenanceTasksMemory(c *gin.Context) {
	store := memory.GetStore()
	var tasks []*model.MaintenanceTask
	for _, t := range store.MaintenanceTasks {
		if t.Status != model.MaintenanceCompleted {
			tasks = append(tasks, t)
		}
	}
	c.JSON(200, gin.H{"items": tasks, "total": len(tasks)})
}

// GetMaintenanceStatistics 获取保养统计 (内存模式)
func GetMaintenanceStatisticsMemory(c *gin.Context) {
	store := memory.GetStore()
	totalTasks := len(store.MaintenanceTasks)
	pending := 0
	inProgress := 0
	completed := 0
	overdue := 0

	for _, t := range store.MaintenanceTasks {
		switch t.Status {
		case "pending":
			pending++
		case "in_progress":
			inProgress++
		case "completed":
			completed++
		case "overdue":
			overdue++
		}
	}

	completionRate := 0.0
	if totalTasks > 0 {
		completionRate = float64(completed) / float64(totalTasks) * 100
	}

	c.JSON(200, gin.H{
		"total_plans":       1,
		"total_tasks":       totalTasks,
		"pending_tasks":     pending,
		"in_progress_tasks": inProgress,
		"completed_tasks":   completed,
		"overdue_tasks":     overdue,
		"today_completed":   completed,
		"completion_rate":   completionRate,
	})
}

// ============ 备件相关 ============

// ListSpareParts 获取备件列表 (内存模式)
func ListSparePartsMemory(c *gin.Context) {
	store := memory.GetStore()
	var parts []*model.SparePart
	for _, p := range store.SpareParts {
		parts = append(parts, p)
	}
	c.JSON(200, gin.H{"items": parts, "total": len(parts)})
}

// GetInventory 获取库存列表 (内存模式)
func GetInventoryMemory(c *gin.Context) {
	store := memory.GetStore()
	var inventories []*model.SparePartInventory
	for _, i := range store.SparePartInventory {
		inventories = append(inventories, i)
	}
	c.JSON(200, gin.H{"items": inventories, "total": len(inventories)})
}

// GetSparePartStatistics 获取备件统计 (内存模式)
func GetSparePartStatisticsMemory(c *gin.Context) {
	c.JSON(200, gin.H{
		"total_parts":        2,
		"low_stock_count":    0,
		"total_stock_value":  50000.0,
		"monthly_consumption": 10,
	})
}

// ============ 统计分析 ============

// GetDashboardOverview 获取仪表盘概览 (内存模式)
func GetDashboardOverviewMemory(c *gin.Context) {
	store := memory.GetStore()

	// 统计设备状态
	totalEquipment := len(store.Equipment)
	running := 0
	stopped := 0
	maintenance := 0
	scrapped := 0

	for _, eq := range store.Equipment {
		switch eq.Status {
		case "running":
			running++
		case "stopped":
			stopped++
		case "maintenance":
			maintenance++
		case "scrapped":
			scrapped++
		}
	}

	// 计算MTTR/MTBF
	totalDowntime := 0.0
	repairCount := 0

	// 从 RepairLog.Content 中提取停机时长
	for _, log := range store.RepairLogs {
		// Content格式: "停机时长: X.XX小时"
		if strings.Contains(log.Content, "停机时长:") {
			parts := strings.Split(log.Content, ":")
			if len(parts) > 1 {
				hourStr := strings.TrimSuffix(strings.TrimSpace(parts[1]), "小时")
				if hours, err := strconv.ParseFloat(hourStr, 64); err == nil {
					totalDowntime += hours
					repairCount++
				}
			}
		}
	}

	// MTTR = 平均修复时间（小时）
	mttr := 0.0
	if repairCount > 0 {
		mttr = totalDowntime / float64(repairCount)
	}

	// MTBF = 平均故障间隔时间（小时）
	// 计算：运行总时长 / 故障次数
	now := time.Now()
	totalOperatingHours := 0.0
	for _, eq := range store.Equipment {
		if !eq.PurchaseDate.IsZero() {
			hours := now.Sub(*eq.PurchaseDate).Hours()
			totalOperatingHours += hours
		}
	}
	mtbf := 0.0
	if repairCount > 0 {
		// MTBF = (总运行时间 - 总停机时间) / 故障次数
		mtbf = (totalOperatingHours - totalDowntime) / float64(repairCount)
	}

	// 可用率
	availability := 0.0
	if totalOperatingHours > 0 {
		uptime := totalOperatingHours - totalDowntime
		availability = (uptime / totalOperatingHours) * 100
	}

	// 任务完成率统计
	inspectionTotal := len(store.InspectionTasks)
	inspectionCompleted := 0
	for _, task := range store.InspectionTasks {
		if task.Status == model.InspectionCompleted {
			inspectionCompleted++
		}
	}
	inspectionRate := 0.0
	if inspectionTotal > 0 {
		inspectionRate = (float64(inspectionCompleted) / float64(inspectionTotal)) * 100
	}

	maintenanceTotal := len(store.MaintenanceTasks)
	maintenanceCompleted := 0
	for _, task := range store.MaintenanceTasks {
		if task.Status == model.MaintenanceCompleted {
			maintenanceCompleted++
		}
	}
	maintenanceRate := 0.0
	if maintenanceTotal > 0 {
		maintenanceRate = (float64(maintenanceCompleted) / float64(maintenanceTotal)) * 100
	}

	repairTotal := len(store.RepairOrders)
	repairCompleted := 0
	for _, order := range store.RepairOrders {
		if order.Status == model.RepairAudited || order.Status == model.RepairClosed {
			repairCompleted++
		}
	}
	repairRate := 0.0
	if repairTotal > 0 {
		repairRate = (float64(repairCompleted) / float64(repairTotal)) * 100
	}

	// 待处理任务统计
	pendingInspections := 0
	for _, task := range store.InspectionTasks {
		if task.Status == model.InspectionPending {
			pendingInspections++
		}
	}

	pendingMaintenances := 0
	for _, task := range store.MaintenanceTasks {
		if task.Status == model.MaintenancePending {
			pendingMaintenances++
		}
	}

	pendingRepairs := 0
	for _, order := range store.RepairOrders {
		if order.Status == model.RepairPending || order.Status == model.RepairAssigned {
			pendingRepairs++
		}
	}

	// 低库存预警
	lowStockAlerts := 0
	for _, inv := range store.SparePartInventory {
		if part, ok := store.SpareParts[inv.SparePartID]; ok {
			if inv.Quantity < part.SafetyStock {
				lowStockAlerts++
			}
		}
	}

	c.JSON(200, gin.H{
		"equipment": gin.H{
			"total_equipment":       totalEquipment,
			"running_equipment":     running,
			"stopped_equipment":     stopped,
			"maintenance_equipment": maintenance,
			"scrapped_equipment":    scrapped,
		},
		"mttr_mtbf": gin.H{
			"mttr":         mttr,
			"mtbf":         mtbf,
			"availability": availability,
		},
		"tasks": gin.H{
			"inspection_completion_rate":   inspectionRate,
			"maintenance_completion_rate":  maintenanceRate,
			"repair_completion_rate":       repairRate,
		},
		"pending_inspections":  pendingInspections,
		"pending_maintenances": pendingMaintenances,
		"pending_repairs":      pendingRepairs,
		"low_stock_alerts":     lowStockAlerts,
	})
}

// GetMTTRMTBF 获取MTTR/MTBF (内存模式)
func GetMTTRMTBFMemory(c *gin.Context) {
	c.JSON(200, gin.H{
		"mttr":         120.5, // 分钟
		"mtbf":         480.5, // 小时
		"availability": 99.5,
	})
}

// GetTrendData 获取趋势数据 (内存模式)
func GetTrendDataMemory(c *gin.Context) {
	store := memory.GetStore()

	type TrendItem struct {
		Date             string  `json:"date"`
		InspectionTasks  int     `json:"inspection_tasks"`
		MaintenanceTasks int     `json:"maintenance_tasks"`
		RepairOrders     int     `json:"repair_orders"`
		DowntimeHours    float64 `json:"downtime_hours"`
	}

	// 获取日期范围参数
	days := 30 // 默认30天
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 365 {
			days = parsed
		}
	}

	var result []TrendItem
	now := time.Now()

	// 按日期统计数据
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		dayStart := date.Truncate(24 * time.Hour)
		dayEnd := dayStart.Add(24 * time.Hour)

		// 统计当天的点检任务
		inspectionCount := 0
		for _, task := range store.InspectionTasks {
			if (task.ScheduledDate.Equal(dayStart) || task.ScheduledDate.After(dayStart)) && task.ScheduledDate.Before(dayEnd) {
				inspectionCount++
			}
		}

		// 统计当天的保养任务（ScheduledDate是字符串格式"2006-01-02"）
		maintenanceCount := 0
		for _, task := range store.MaintenanceTasks {
			if task.ScheduledDate == dateStr {
				maintenanceCount++
			}
		}

		// 统计当天的维修工单
		repairCount := 0
		downtimeHours := 0.0
		for _, order := range store.RepairOrders {
			if (order.CreatedAt.Equal(dayStart) || order.CreatedAt.After(dayStart)) && order.CreatedAt.Before(dayEnd) {
				repairCount++
			}
		}

		// 统计当天的停机时长（从RepairLog.Content中提取）
		for _, log := range store.RepairLogs {
			if (log.CreatedAt.Equal(dayStart) || log.CreatedAt.After(dayStart)) && log.CreatedAt.Before(dayEnd) {
				// Content格式: "停机时长: X.XX小时"
				if strings.Contains(log.Content, "停机时长:") {
					parts := strings.Split(log.Content, ":")
					if len(parts) > 1 {
						hourStr := strings.TrimSuffix(strings.TrimSpace(parts[1]), "小时")
						if hours, err := strconv.ParseFloat(hourStr, 64); err == nil {
							downtimeHours += hours
						}
					}
				}
			}
		}

		result = append(result, TrendItem{
			Date:             dateStr,
			InspectionTasks:  inspectionCount,
			MaintenanceTasks: maintenanceCount,
			RepairOrders:     repairCount,
			DowntimeHours:    downtimeHours,
		})
	}

	c.JSON(200, result)
}

// GetFailureAnalysis 获取故障分析 (内存模式)
func GetFailureAnalysisMemory(c *gin.Context) {
	store := memory.GetStore()

	type FailureItem struct {
		EquipmentTypeID   uint    `json:"equipment_type_id"`
		EquipmentTypeName string  `json:"equipment_type_name"`
		FailureCount      int     `json:"failure_count"`
		TotalDowntime     float64 `json:"total_downtime"`
	}

	// 按设备类型统计故障
	typeStats := make(map[uint]*FailureItem)

	for _, order := range store.RepairOrders {
		if eq, ok := store.Equipment[order.EquipmentID]; ok {
			if _, exists := typeStats[eq.TypeID]; !exists {
				typeStats[eq.TypeID] = &FailureItem{
					EquipmentTypeID:   eq.TypeID,
					EquipmentTypeName: "", // 后面更新
					FailureCount:      0,
					TotalDowntime:     0,
				}
			}

			typeStats[eq.TypeID].FailureCount++

			// 累加停机时长（从RepairLog.Content中提取）
			for _, log := range store.RepairLogs {
				if log.OrderID == order.ID {
					// Content格式: "停机时长: X.XX小时"
					if strings.Contains(log.Content, "停机时长:") {
						parts := strings.Split(log.Content, ":")
						if len(parts) > 1 {
							hourStr := strings.TrimSuffix(strings.TrimSpace(parts[1]), "小时")
							if hours, err := strconv.ParseFloat(hourStr, 64); err == nil {
								typeStats[eq.TypeID].TotalDowntime += hours
							}
						}
					}
				}
			}
		}
	}

	// 填充设备类型名称
	for _, item := range typeStats {
		if eqType, ok := store.EquipmentTypes[item.EquipmentTypeID]; ok {
			item.EquipmentTypeName = eqType.Name
		}
	}

	// 转换为数组并按故障次数排序
	var result []*FailureItem
	for _, item := range typeStats {
		result = append(result, item)
	}

	// 排序：故障次数降序
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].FailureCount > result[i].FailureCount {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	// 限制返回数量
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if len(result) > limit {
		result = result[:limit]
	}

	c.JSON(200, result)
}

// GetTopFailureEquipment 获取故障率最高的设备 (内存模式)
func GetTopFailureEquipmentMemory(c *gin.Context) {
	store := memory.GetStore()

	type TopFailureItem struct {
		EquipmentID   uint    `json:"equipment_id"`
		EquipmentCode string  `json:"equipment_code"`
		EquipmentName string  `json:"equipment_name"`
		FailureCount  int     `json:"failure_count"`
		DowntimeHours float64 `json:"downtime_hours"`
		MTTR          float64 `json:"mttr"`
	}

	// 按设备统计故障
	equipmentStats := make(map[uint]*TopFailureItem)

	for _, order := range store.RepairOrders {
		if eq, ok := store.Equipment[order.EquipmentID]; ok {
			if _, exists := equipmentStats[eq.ID]; !exists {
				equipmentStats[eq.ID] = &TopFailureItem{
					EquipmentID:   eq.ID,
					EquipmentCode: eq.Code,
					EquipmentName: eq.Name,
					FailureCount:  0,
					DowntimeHours: 0,
					MTTR:          0,
				}
			}

			equipmentStats[eq.ID].FailureCount++

			// 累加停机时长并计算MTTR（从RepairLog.Content中提取）
			totalDowntime := 0.0
			repairCount := 0
			for _, log := range store.RepairLogs {
				if log.OrderID == order.ID {
					// Content格式: "停机时长: X.XX小时"
					if strings.Contains(log.Content, "停机时长:") {
						parts := strings.Split(log.Content, ":")
						if len(parts) > 1 {
							hourStr := strings.TrimSuffix(strings.TrimSpace(parts[1]), "小时")
							if hours, err := strconv.ParseFloat(hourStr, 64); err == nil {
								totalDowntime += hours
								repairCount++
							}
						}
					}
				}
			}

			equipmentStats[eq.ID].DowntimeHours += totalDowntime
			if repairCount > 0 {
				equipmentStats[eq.ID].MTTR = equipmentStats[eq.ID].DowntimeHours / float64(equipmentStats[eq.ID].FailureCount)
			}
		}
	}

	// 转换为数组并按故障次数排序
	var result []*TopFailureItem
	for _, item := range equipmentStats {
		result = append(result, item)
	}

	// 排序：故障次数降序
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].FailureCount > result[i].FailureCount {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	// 限制返回数量
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if len(result) > limit {
		result = result[:limit]
	}

	c.JSON(200, result)
}

// ============ 知识库 ============

// ListKnowledgeArticles 获取知识库文章列表 (内存模式)
func ListKnowledgeArticlesMemory(c *gin.Context) {
	store := memory.GetStore()
	var articles []*model.KnowledgeArticle
	for _, a := range store.KnowledgeArticles {
		articles = append(articles, a)
	}
	c.JSON(200, gin.H{"items": articles, "total": len(articles)})
}

// SearchKnowledgeArticles 搜索知识库文章 (内存模式)
func SearchKnowledgeArticlesMemory(c *gin.Context) {
	_ = c.Query("keyword") // TODO: 实现关键词搜索
	store := memory.GetStore()
	var results []gin.H
	for _, a := range store.KnowledgeArticles {
		results = append(results, gin.H{
			"id":         a.ID,
			"title":      a.Title,
			"category":   "故障处理",
			"summary":    a.CauseAnalysis,
			"created_at": a.CreatedAt,
		})
	}
	c.JSON(200, results)
}

// ============ 健康检查 ============

// HealthCheck 健康检查 (内存模式)
func HealthCheckMemory(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "ems-api",
		"mode":    "memory",
	})
}
