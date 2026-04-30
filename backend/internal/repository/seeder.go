package repository

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ems/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func mustCreateBatch(db *gorm.DB, name string, value interface{}, batchSize int) error {
	if err := db.CreateInBatches(value, batchSize).Error; err != nil {
		log.Printf("❌ Seeder: failed to create batch %s: %v", name, err)
		return fmt.Errorf("create batch %s: %w", name, err)
	}
	return nil
}

func SeedDatabase(db *gorm.DB) error {
	now := time.Now()
	rnd := rand.New(rand.NewSource(now.UnixNano()))

	hp, err := bcryptHash("admin123")
	if err != nil {
		return fmt.Errorf("bcrypt hash: %w", err)
	}

	// Organization
	base, err := findOrCreateBase(db, "BASE-HQ", "集团总部基地")
	if err != nil { return err }
	fac, err := findOrCreateFactory(db, "FAC-SZ", "苏州智能工厂", base.ID)
	if err != nil { return err }
	ws1, err := findOrCreateWorkshop(db, "WS-MCH", "精密机加车间", fac.ID)
	if err != nil { return err }
	ws2, err := findOrCreateWorkshop(db, "WS-ASM", "全自动装配车间", fac.ID)
	if err != nil { return err }
	workshops := []uint{ws1.ID, ws2.ID}

	// Users
	admin, err := findOrCreateUser(db, "admin", hp, "系统管理员", model.RoleAdmin, nil, true)
	if err != nil { return err }
	db.Model(&admin).Updates(map[string]interface{}{"password_hash": hp, "is_active": true, "approval_status": model.ApprovalStatusApproved})

	var operators []model.User
	for i := 1; i <= 10; i++ {
		u, err := findOrCreateUser(db, fmt.Sprintf("operator_%d", i), hp, fmt.Sprintf("操作工%d号", i), model.RoleOperator, &fac.ID, true)
		if err != nil { return err }
		operators = append(operators, *u)
	}

	var maints []model.User
	for i := 1; i <= 5; i++ {
		u, err := findOrCreateUser(db, fmt.Sprintf("maint_%d", i), hp, fmt.Sprintf("维修工%d号", i), model.RoleMaintenance, &fac.ID, true)
		if err != nil { return err }
		maints = append(maints, *u)
	}

	for i := 1; i <= 2; i++ {
		_, err := findOrCreateUser(db, fmt.Sprintf("eng_%d", i), hp, fmt.Sprintf("工程师%d号", i), model.RoleEngineer, &fac.ID, true)
		if err != nil { return err }
	}

	var equipCount int64
	db.Model(&model.Equipment{}).Count(&equipCount)
	if equipCount >= 100 {
		log.Println("ℹ️  Business data already exists (>=100 equipments), skipping full seeding.")
		return nil
	}

	log.Println("🚀 Starting production-scale business data seeding...")

	// 5 Equipment Types
	typeNames := []string{"数控机床", "全自动冲床", "工业机器人", "注塑机", "自动化流水线"}
	var equipTypes []model.EquipmentType
	for _, name := range typeNames {
		et := model.EquipmentType{Name: name, Category: "生产设备"}
		db.Create(&et)
		equipTypes = append(equipTypes, et)

		// Templates
		inspTemp := model.InspectionTemplate{Name: name + "日常点检", EquipmentTypeID: et.ID}
		db.Create(&inspTemp)
		items := []model.InspectionItem{
			{TemplateID: inspTemp.ID, Name: "外观检查", Method: "目视", Criteria: "无破损", SequenceOrder: 1},
			{TemplateID: inspTemp.ID, Name: "运行状态", Method: "耳听/目视", Criteria: "平稳无异响", SequenceOrder: 2},
			{TemplateID: inspTemp.ID, Name: "关键参数", Method: "仪表读取", Criteria: "在标准范围内", SequenceOrder: 3},
		}
		db.Create(&items)

		maintPlan := model.MaintenancePlan{Name: name + "月度保养", EquipmentTypeID: et.ID, CycleDays: 30, Level: 1, WorkHours: 2.0}
		db.Create(&maintPlan)
		maintItems := []model.MaintenancePlanItem{
			{PlanID: maintPlan.ID, Name: "全面清洁", Method: "人工", Criteria: "无油污、灰尘", SequenceOrder: 1},
			{PlanID: maintPlan.ID, Name: "加注润滑油", Method: "注油", Criteria: "油位正常", SequenceOrder: 2},
			{PlanID: maintPlan.ID, Name: "紧固螺栓", Method: "工具", Criteria: "无松动", SequenceOrder: 3},
		}
		db.Create(&maintItems)
	}

	// 20+ Spare Parts
	var spareParts []model.SparePart
	partNames := []string{
		"深沟球轴承 6205", "圆柱滚子轴承 NU205", "推力球轴承 51105", 
		"高压柱塞泵 A10V", "齿轮泵 CBT-F4", "叶片泵 CBQ-G5", 
		"交流接触器 CJX2-18", "热继电器 JR28", "微型断路器 DZ47", 
		"光电开关 PR18", "接近开关 E2E", "压力传感器 PT100", 
		"普通滤芯", "高精度滤芯", "回油滤芯", 
		"抗磨液压油 HM46", "齿轮油 L-CKC 150", "润滑脂 2#", 
		"O型密封圈 NBR50", "骨架油封 TC 35x50x8", "密封垫片", 
		"步进电机 400W", "伺服电机 750W", "变频器 1.5kW",
	}
	for i, name := range partNames {
		sp := model.SparePart{
			Code: fmt.Sprintf("SP-%03d", i+1), Name: name, 
			FactoryID: &fac.ID, Unit: "件", SafetyStock: rnd.Intn(20) + 5,
		}
		spareParts = append(spareParts, sp)
	}
	mustCreateBatch(db, "SpareParts", spareParts, 100)
	db.Find(&spareParts)

	var inventories []model.SparePartInventory
	for _, sp := range spareParts {
		qty := rnd.Intn(100) + 50
		inventories = append(inventories, model.SparePartInventory{SparePartID: sp.ID, FactoryID: fac.ID, Quantity: qty})
		
		// Transaction with BaseModel properly initialized
		tx := model.SparePartTransaction{
			SparePartID: sp.ID, FactoryID: fac.ID, Type: "in", Quantity: qty, OperatorID: admin.ID, Remark: "初始入库",
		}
		tx.CreatedAt = now.AddDate(0, -3, 0)
		db.Create(&tx)
	}
	mustCreateBatch(db, "SparePartInventories", inventories, 100)

	// 100 Equipments
	var equipments []model.Equipment
	for i := 1; i <= 100; i++ {
		eType := equipTypes[rnd.Intn(len(equipTypes))]
		maint := maints[rnd.Intn(len(maints))]
		ws := workshops[rnd.Intn(len(workshops))]
		status := "running"
		if rnd.Float32() < 0.05 { status = "stopped" } else if rnd.Float32() < 0.1 { status = "maintenance" }

		eq := model.Equipment{
			Code: fmt.Sprintf("EQ-%04d", i), Name: fmt.Sprintf("%s-%03d", eType.Name, i),
			TypeID: eType.ID, WorkshopID: ws, QRCode: fmt.Sprintf("EQ-%04d", i),
			PurchasePrice: float64(rnd.Intn(500000) + 50000), PurchaseDate: timePtr(now.AddDate(-rnd.Intn(5), -rnd.Intn(12), 0)),
			ServiceLifeYears: 10, ScrapValue: 5000.0, HourlyLoss: float64(rnd.Intn(400) + 50),
			Status: status, DedicatedMaintenanceID: &maint.ID,
		}
		equipments = append(equipments, eq)
	}
	mustCreateBatch(db, "Equipments", equipments, 100)
	db.Find(&equipments)

	// Load templates for referencing
	var inspTemplates []model.InspectionTemplate
	db.Preload("Items").Find(&inspTemplates)
	tempMap := make(map[uint]*model.InspectionTemplate)
	for i := range inspTemplates { tempMap[inspTemplates[i].EquipmentTypeID] = &inspTemplates[i] }

	var maintPlans []model.MaintenancePlan
	db.Preload("Items").Find(&maintPlans)
	planMap := make(map[uint]*model.MaintenancePlan)
	for i := range maintPlans { planMap[maintPlans[i].EquipmentTypeID] = &maintPlans[i] }

	log.Println("⏳ Generating historical data (15 days)...")

	historyDays := 15
	for d := historyDays; d >= 0; d-- {
		date := now.AddDate(0, 0, -d)
		dateStr := date.Format("2006-01-02")
		
		for _, eq := range equipments {
			op := operators[rnd.Intn(len(operators))]
			maint := maints[rnd.Intn(len(maints))]

			// 1. Inspection
			temp := tempMap[eq.TypeID]
			iTask := model.InspectionTask{
				EquipmentID: eq.ID, TemplateID: temp.ID, AssignedTo: op.ID,
				ScheduledDate: date, Status: "completed", CompletedAt: timePtr(date.Add(time.Minute * time.Duration(rnd.Intn(30)+5))),
			}
			db.Create(&iTask)

			hasNG := false
			for _, item := range temp.Items {
				res := "OK"
				if rnd.Float32() < 0.03 { res = "NG"; hasNG = true }
				db.Create(&model.InspectionRecord{TaskID: iTask.ID, ItemID: item.ID, Result: res})
			}

			// 2. Maintenance
			if d == 15 || d == 0 { // at start and end of history
				mPlan := planMap[eq.TypeID]
				mTask := model.MaintenanceTask{
					PlanID: mPlan.ID, EquipmentID: eq.ID, AssignedTo: maint.ID,
					Status: "completed", CompletedAt: timePtr(date.Add(time.Hour * 2)), ActualHours: float64(rnd.Intn(3)+1), ScheduledDate: dateStr,
				}
				db.Create(&mTask)
				for _, item := range mPlan.Items {
					db.Create(&model.MaintenanceRecord{TaskID: mTask.ID, ItemID: item.ID, Result: "OK"})
				}
			}

			// 3. Repair if NG
			if hasNG {
				ro := model.RepairOrder{
					EquipmentID: eq.ID, Status: model.RepairClosed, Priority: rnd.Intn(3)+1,
					AssignedTo: &maint.ID, ReporterID: op.ID, FaultDescription: "日常点检发现异常，需要紧急修复",
					StartedAt: timePtr(date.Add(time.Hour)), Solution: "修复完成并已恢复生产",
					CompletedAt: timePtr(date.Add(time.Hour * 4)), ClosedAt: timePtr(date.Add(time.Hour * 6)),
				}
				db.Create(&ro)
				db.Create(&model.RepairCostDetail{OrderID: ro.ID, SparePartCost: float64(rnd.Intn(2000)), LaborCost: float64(rnd.Intn(1000))})
				
				// Random spare parts
				for i := 0; i < rnd.Intn(2)+1; i++ {
					sp := spareParts[rnd.Intn(len(spareParts))]
					qty := rnd.Intn(2)+1
					db.Create(&model.SparePartConsumption{SparePartID: sp.ID, OrderID: &ro.ID, Quantity: qty, UserID: maint.ID})
				}
			}
		}
	}

	// Final Static Data
	manual := model.ManualDocument{Title: "通用生产设备维护手册", EquipmentTypeID: &equipTypes[0].ID}
	db.Create(&manual)
	db.Create(&model.ManualChunk{DocumentID: manual.ID, SectionTitle: "通用维护", Content: "定期检查和润滑是保持设备长寿命的关键...", PageNumber: 1})

	articles := []model.KnowledgeArticle{
		{Title: "高频故障：传感器失效分析", EquipmentTypeID: &equipTypes[0].ID, FaultPhenomenon: "信号丢失或波动", CauseAnalysis: "环境油污干扰", Solution: "清理并加装防护罩", CreatedBy: admin.ID, SourceType: "expert_summary"},
	}
	db.Create(&articles)

	db.Create(&model.AgentSkill{Name: "生产效率分析", Description: "分析停机损失与产出比", Status: "active", Steps: `[{"step": 1, "tool": "get_cost_analysis"}]`})
	
	log.Println("✅ Production-scale seeding complete. Login: admin / admin123")
	return nil
}

func bcryptHash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { return "", err }
	return string(hashedBytes), nil
}

func findOrCreateBase(db *gorm.DB, code, name string) (*model.Base, error) {
	var base model.Base
	if err := db.Where("code = ?", code).First(&base).Error; err != nil {
		base = model.Base{Code: code, Name: name}
		if err := db.Create(&base).Error; err != nil { return nil, err }
	}
	return &base, nil
}

func findOrCreateFactory(db *gorm.DB, code, name string, baseID uint) (*model.Factory, error) {
	var fac model.Factory
	if err := db.Where("code = ?", code).First(&fac).Error; err != nil {
		fac = model.Factory{BaseID: baseID, Code: code, Name: name}
		if err := db.Create(&fac).Error; err != nil { return nil, err }
	}
	return &fac, nil
}

func findOrCreateWorkshop(db *gorm.DB, code, name string, factoryID uint) (*model.Workshop, error) {
	var ws model.Workshop
	if err := db.Where("code = ?", code).First(&ws).Error; err != nil {
		ws = model.Workshop{FactoryID: factoryID, Code: code, Name: name}
		if err := db.Create(&ws).Error; err != nil { return nil, err }
	}
	return &ws, nil
}

func findOrCreateUser(db *gorm.DB, username, passwordHash, name string, role model.UserRole, factoryID *uint, isActive bool) (*model.User, error) {
	var user model.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		user = model.User{
			Username: username, PasswordHash: passwordHash, Name: name,
			Role: role, FactoryID: factoryID, IsActive: isActive,
			ApprovalStatus: model.ApprovalStatusApproved,
		}
		if err := db.Create(&user).Error; err != nil { return nil, err }
	}
	return &user, nil
}

func timePtr(t time.Time) *time.Time { return &t }
