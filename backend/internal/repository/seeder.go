package repository

import (
	"log"
	"time"

	"github.com/ems/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase 全模块演示数据补全逻辑 (动态哈希加固版)
func SeedDatabase(db *gorm.DB) error {
	// 1. 检查核心管理员是否存在
	var adminUser model.User
	err := db.Where("username = ?", "admin").First(&adminUser).Error
	
	if err == nil {
		log.Println("Admin user already exists, skipping seeding.")
		return nil
	}

	log.Println("🚀 Database seems empty or missing core users. Starting seeding...")
	now := time.Now()
	
	// 动态生成哈希，确保 100% 匹配
	cost := bcrypt.DefaultCost
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte("admin123"), cost)
	if err != nil {
		return err
	}
	hp := string(hashedBytes)

	// 2. Organization
	base := model.Base{Code: "BASE-HQ", Name: "集团总部基地"}
	db.Create(&base)
	fac := model.Factory{BaseID: base.ID, Code: "FAC-SZ", Name: "苏州智能工厂"}
	db.Create(&fac)
	ws1 := model.Workshop{FactoryID: fac.ID, Code: "WS-MCH", Name: "精密机加车间"}
	ws2 := model.Workshop{FactoryID: fac.ID, Code: "WS-ASM", Name: "全自动装配车间"}
	db.Create([]*model.Workshop{&ws1, &ws2})

	// 3. Users
	admin := model.User{Username: "admin", PasswordHash: hp, Name: "系统管理员", Role: model.RoleAdmin, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
	liSi := model.User{Username: "maint_li", PasswordHash: hp, Name: "预防型-李四", Role: model.RoleMaintenance, FactoryID: &fac.ID, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
	zs := model.User{Username: "maint_zhang", PasswordHash: hp, Name: "救火型-张三", Role: model.RoleMaintenance, FactoryID: &fac.ID, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
	op := model.User{Username: "operator", PasswordHash: hp, Name: "操作员小王", Role: model.RoleOperator, FactoryID: &fac.ID, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
	db.Create([]*model.User{&admin, &liSi, &zs, &op})

	// 4. Equipment Types & Templates
	cncType := model.EquipmentType{Name: "高精度数控机床", Category: "加工设备"}
	pressType := model.EquipmentType{Name: "全自动冲床", Category: "成型设备"}
	robotType := model.EquipmentType{Name: "工业机器人", Category: "智能装备"}
	db.Create([]*model.EquipmentType{&cncType, &pressType, &robotType})

	insTemp := model.InspectionTemplate{Name: "CNC日常点检标准", EquipmentTypeID: cncType.ID}
	db.Create(&insTemp)
	db.Create([]*model.InspectionItem{
		{TemplateID: insTemp.ID, Name: "液压压力检查", Method: "观察压力表", Criteria: "4.0-5.0 MPa", SequenceOrder: 1},
		{TemplateID: insTemp.ID, Name: "主轴温度", Method: "手感/红外", Criteria: "无烫感", SequenceOrder: 2},
	})

	// 5. Equipment
	e1 := model.Equipment{
		Code: "CNC-001", Name: "李四维护-A区机床", TypeID: cncType.ID, WorkshopID: ws1.ID, QRCode: "QR_A",
		PurchasePrice: 280000.0, PurchaseDate: timePtr(now.AddDate(-3, 0, 0)), ServiceLifeYears: 8, ScrapValue: 28000.0, HourlyLoss: 150.0, Status: "running", DedicatedMaintenanceID: &liSi.ID,
	}
	e2 := model.Equipment{
		Code: "CNC-002", Name: "张三维护-B区机床", TypeID: cncType.ID, WorkshopID: ws1.ID, QRCode: "QR_B",
		PurchasePrice: 280000.0, PurchaseDate: timePtr(now.AddDate(-3, 0, 0)), ServiceLifeYears: 8, ScrapValue: 28000.0, HourlyLoss: 150.0, Status: "stopped", DedicatedMaintenanceID: &zs.ID,
	}
	eP := model.Equipment{
		Code: "PRESS-05", Name: "12年老旧冲床", TypeID: pressType.ID, WorkshopID: ws1.ID, QRCode: "QR_OLD",
		PurchasePrice: 150000.0, PurchaseDate: timePtr(now.AddDate(-12, 0, 0)), ServiceLifeYears: 10, ScrapValue: 5000.0, HourlyLoss: 80.0, Status: "maintenance",
	}
	eR := model.Equipment{
		Code: "ROBOT-01", Name: "ABB码垛机器人", TypeID: robotType.ID, WorkshopID: ws2.ID, QRCode: "QR_R1",
		PurchasePrice: 450000.0, PurchaseDate: timePtr(now.AddDate(-1, 0, 0)), ServiceLifeYears: 5, ScrapValue: 45000.0, HourlyLoss: 300.0, Status: "running",
	}
	db.Create([]*model.Equipment{&e1, &e2, &eP, &eR})

	// 6. Spare Parts
	pump := model.SparePart{Code: "PUMP-01", Name: "高压柱塞泵", FactoryID: &fac.ID, Specification: "Rexroth A10V", Unit: "台", SafetyStock: 2}
	filter := model.SparePart{Code: "FLT-CHEAP", Name: "普通滤芯(降本件)", FactoryID: &fac.ID, Unit: "个", SafetyStock: 50}
	oil := model.SparePart{Code: "OIL-HM46", Name: "液压油", FactoryID: &fac.ID, Unit: "桶", SafetyStock: 100}
	db.Create([]*model.SparePart{&pump, &filter, &oil})
	db.Create([]*model.SparePartInventory{
		{SparePartID: pump.ID, FactoryID: fac.ID, Quantity: 3},
		{SparePartID: filter.ID, FactoryID: fac.ID, Quantity: 120},
		{SparePartID: oil.ID, FactoryID: fac.ID, Quantity: 250},
	})

	// 7. 180 Days of Maintenance
	mPlan := model.MaintenancePlan{Name: "CNC二级保养", EquipmentTypeID: cncType.ID, CycleDays: 30, Level: 2, WorkHours: 4.0}
	db.Create(&mPlan)
	for i := 1; i <= 6; i++ {
		date := now.AddDate(0, 0, -i*30)
		db.Create(&model.MaintenanceTask{PlanID: mPlan.ID, EquipmentID: e1.ID, AssignedTo: liSi.ID, Status: "completed", CompletedAt: &date, ActualHours: 4.2, ScheduledDate: date.Format("2006-01-02")})
		db.Create(&model.MaintenanceTask{PlanID: mPlan.ID, EquipmentID: e2.ID, AssignedTo: zs.ID, Status: "completed", CompletedAt: &date, ActualHours: 0.2, ScheduledDate: date.Format("2006-01-02")})
	}

	// 8. Repairs (The 5.5W Incident)
	failDate := now.AddDate(0, 0, -10)
	order := model.RepairOrder{
		EquipmentID: e2.ID, Status: model.RepairClosed, Priority: 1, AssignedTo: &zs.ID, ReporterID: op.ID,
		FaultDescription: "液压泵损坏", Solution: "更换泵及滤芯", StartedAt: &failDate, CompletedAt: timePtr(failDate.Add(4 * time.Hour)), ClosedAt: timePtr(failDate.Add(24 * time.Hour)),
	}
	db.Create(&order)
	db.Create(&model.RepairCostDetail{OrderID: order.ID, SparePartCost: 52000.0, LaborCost: 3000.0})

	// 9. Agent Brain
	db.Create(&model.AgentSkill{
		Name: "级联失效审计", Description: "分析连锁损失风险", Status: "active",
		Steps: `[{"step": 1, "tool": "get_cost_analysis"}, {"step": 2, "tool": "get_maintenance_compliance"}]`,
	})
	db.Create(&model.AgentExperience{
		UserID: admin.ID, Category: "reporting", Content: "该用户偏好简洁的财务摘要，重点关注 TCO 异常。", Weight: 1.0, Status: "active",
	})
	db.Create(&model.KnowledgeArticle{
		Title: "Rexroth泵早期损坏特征", EquipmentTypeID: &cncType.ID, FaultPhenomenon: "压力异常波动", CauseAnalysis: "滤芯堵塞导致空化", Solution: "更换原厂滤芯", CreatedBy: admin.ID, SourceType: "expert_summary",
	})

	log.Println("✅ Comprehensive Seeding Complete. Login with admin/admin123 now!")
	return nil
}
