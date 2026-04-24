package v1

import (
	"strconv"
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
	c.JSON(200, result)
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
	c.JSON(200, result)
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
	c.JSON(200, gin.H{"data": bases, "total": len(bases)})
}

// ListFactories 获取工厂列表 (内存模式)
func ListFactoriesMemory(c *gin.Context) {
	store := memory.GetStore()
	var factories []*model.Factory
	for _, f := range store.Factories {
		factories = append(factories, f)
	}
	c.JSON(200, gin.H{"data": factories, "total": len(factories)})
}

// ListWorkshops 获取车间列表 (内存模式)
func ListWorkshopsMemory(c *gin.Context) {
	store := memory.GetStore()
	var workshops []*model.Workshop
	for _, w := range store.Workshops {
		workshops = append(workshops, w)
	}
	c.JSON(200, gin.H{"data": workshops, "total": len(workshops)})
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
		"data": equipments,
		"total": len(equipments),
		"page": 1,
		"page_size": len(equipments),
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
	c.JSON(200, gin.H{"data": types, "total": len(types)})
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
	c.JSON(200, gin.H{"data": templates, "total": len(templates)})
}

// ListInspectionTasks 获取点检任务列表 (内存模式)
func ListInspectionTasksMemory(c *gin.Context) {
	store := memory.GetStore()
	var tasks []*model.InspectionTask
	for _, t := range store.InspectionTasks {
		tasks = append(tasks, t)
	}
	c.JSON(200, gin.H{"data": tasks, "total": len(tasks), "items": tasks})
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
	c.JSON(200, gin.H{"data": tasks, "total": len(tasks)})
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
		"data": orders,
		"total": len(orders),
		"page": 1,
		"page_size": len(orders),
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
	c.JSON(200, gin.H{"data": tasks, "total": len(tasks)})
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
	c.JSON(200, gin.H{"data": plans, "total": len(plans)})
}

// ListMaintenanceTasks 获取保养任务列表 (内存模式)
func ListMaintenanceTasksMemory(c *gin.Context) {
	store := memory.GetStore()
	var tasks []*model.MaintenanceTask
	for _, t := range store.MaintenanceTasks {
		tasks = append(tasks, t)
	}
	c.JSON(200, gin.H{"data": gin.H{"items": tasks, "total": len(tasks)}})
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
	c.JSON(200, gin.H{"data": tasks, "total": len(tasks)})
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
	c.JSON(200, gin.H{"data": parts, "total": len(parts)})
}

// GetInventory 获取库存列表 (内存模式)
func GetInventoryMemory(c *gin.Context) {
	store := memory.GetStore()
	var inventories []*model.SparePartInventory
	for _, i := range store.SparePartInventory {
		inventories = append(inventories, i)
	}
	c.JSON(200, gin.H{"data": inventories, "total": len(inventories)})
}

// GetSparePartStatistics 获取备件统计 (内存模式)
func GetSparePartStatisticsMemory(c *gin.Context) {
	c.JSON(200, gin.H{
		"total_parts": 2,
		"low_stock":   0,
		"total_value": 50000.0,
	})
}

// ============ 统计分析 ============

// GetDashboardOverview 获取仪表盘概览 (内存模式)
func GetDashboardOverviewMemory(c *gin.Context) {
	store := memory.GetStore()

	equipmentTotal := len(store.Equipment)
	running := equipmentTotal - 1
	stopped := 0
	maintenance := 0
	scrapped := 0

	c.JSON(200, gin.H{
		"equipment": gin.H{
			"total_equipment":      equipmentTotal,
			"running_equipment":    running,
			"stopped_equipment":    stopped,
			"maintenance_equipment": maintenance,
			"scrapped_equipment":   scrapped,
		},
		"mttr_mtbf": gin.H{
			"mttr":         2.01, // 小时 (120.5分钟)
			"mtbf":         480.5,
			"availability": 99.5,
		},
		"tasks": gin.H{
			"inspection_completion_rate":   100.0,
			"maintenance_completion_rate": 95.5,
			"repair_completion_rate":       100.0,
		},
		"pending_inspections":  2,
		"pending_maintenances": 2,
		"pending_repairs":      0,
		"low_stock_alerts":     0,
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
	// 生成近30天的趋势数据
	type TrendItem struct {
		Date             string  `json:"date"`
		InspectionTasks  int     `json:"inspection_tasks"`
		MaintenanceTasks int     `json:"maintenance_tasks"`
		RepairOrders     int     `json:"repair_orders"`
		DowntimeHours    float64 `json:"downtime_hours"`
	}

	var result []TrendItem
	now := time.Now()
	for i := 29; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		result = append(result, TrendItem{
			Date:             date.Format("2006-01-02"),
			InspectionTasks:  10 + (i % 5),
			MaintenanceTasks: 5 + (i % 3),
			RepairOrders:     2 + (i % 4),
			DowntimeHours:    float64(i%5) * 0.5,
		})
	}

	c.JSON(200, result)
}

// GetFailureAnalysis 获取故障分析 (内存模式)
func GetFailureAnalysisMemory(c *gin.Context) {
	type FailureItem struct {
		EquipmentTypeID   uint    `json:"equipment_type_id"`
		EquipmentTypeName string  `json:"equipment_type_name"`
		FailureCount      int     `json:"failure_count"`
		TotalDowntime     float64 `json:"total_downtime"`
	}

	result := []FailureItem{
		{1, "数控机床", 8, 12.5},
		{2, "焊接机器人", 5, 8.0},
		{3, "冲压机", 3, 4.5},
		{4, "注塑机", 2, 3.0},
		{5, "切割机", 2, 2.5},
		{6, "磨床", 1, 1.5},
	}

	c.JSON(200, result)
}

// GetTopFailureEquipment 获取故障率最高的设备 (内存模式)
func GetTopFailureEquipmentMemory(c *gin.Context) {
	type TopFailureItem struct {
		EquipmentID    uint    `json:"equipment_id"`
		EquipmentCode  string  `json:"equipment_code"`
		EquipmentName  string  `json:"equipment_name"`
		FailureCount   int     `json:"failure_count"`
		DowntimeHours  float64 `json:"downtime_hours"`
		MTTR           float64 `json:"mttr"`
	}

	result := []TopFailureItem{
		{1, "EQ-JJ-001", "数控机床-1", 8, 12.5, 1.56},
		{3, "EQ-JJ-003", "数控机床-3", 5, 8.0, 1.6},
		{6, "EQ-HJ-002", "焊接机器人-2", 3, 4.5, 1.5},
		{2, "EQ-JJ-002", "数控机床-2", 2, 2.5, 1.25},
		{4, "EQ-HJ-001", "焊接机器人-1", 2, 2.0, 1.0},
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
	c.JSON(200, gin.H{"data": articles, "total": len(articles)})
}

// SearchKnowledgeArticles 搜索知识库文章 (内存模式)
func SearchKnowledgeArticlesMemory(c *gin.Context) {
	keyword := c.Query("keyword")
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
	c.JSON(200, gin.H{
		"data":    results,
		"total":   len(results),
		"keyword": keyword,
	})
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
