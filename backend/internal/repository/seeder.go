package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/ems/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// mustCreate wraps db.Create with error logging. Returns error for propagation.
func mustCreate(db *gorm.DB, name string, value interface{}) error {
	if err := db.Create(value).Error; err != nil {
		log.Printf("❌ Seeder: failed to create %s: %v", name, err)
		return fmt.Errorf("create %s: %w", name, err)
	}
	return nil
}

// SeedDatabase seeds comprehensive demo data for all modules.
// Part 1 (org + users) always runs with upsert logic.
// Part 2 (business data) only runs when equipment table is empty.
func SeedDatabase(db *gorm.DB) error {
	now := time.Now()

	// Dynamic bcrypt hash for demo password
	hp, err := bcryptHash("admin123")
	if err != nil {
		return fmt.Errorf("bcrypt hash: %w", err)
	}

	// ============================================================
	// Part 1: Organization + Users (always runs, upsert)
	// ============================================================

	base, err := findOrCreateBase(db, "BASE-HQ", "集团总部基地")
	if err != nil {
		return err
	}

	fac, err := findOrCreateFactory(db, "FAC-SZ", "苏州智能工厂", base.ID)
	if err != nil {
		return err
	}

	ws1, err := findOrCreateWorkshop(db, "WS-MCH", "精密机加车间", fac.ID)
	if err != nil {
		return err
	}

	ws2, err := findOrCreateWorkshop(db, "WS-ASM", "全自动装配车间", fac.ID)
	if err != nil {
		return err
	}

	// Users — force-reset admin password every startup
	admin, err := findOrCreateUser(db, "admin", hp, "系统管理员", model.RoleAdmin, nil, true)
	if err != nil {
		return err
	}
	// Force-reset admin password
	db.Model(&admin).Updates(map[string]interface{}{
		"password_hash":   hp,
		"is_active":       true,
		"approval_status": model.ApprovalStatusApproved,
	})

	liSi, err := findOrCreateUser(db, "maint_li", hp, "预防型-李四", model.RoleMaintenance, &fac.ID, true)
	if err != nil {
		return err
	}

	zs, err := findOrCreateUser(db, "maint_zhang", hp, "救火型-张三", model.RoleMaintenance, &fac.ID, true)
	if err != nil {
		return err
	}

	op, err := findOrCreateUser(db, "operator", hp, "操作员小王", model.RoleOperator, &fac.ID, true)
	if err != nil {
		return err
	}

	// ============================================================
	// Part 2: Business data (skip if equipment already exists)
	// ============================================================

	var equipCount int64
	db.Model(&model.Equipment{}).Count(&equipCount)
	if equipCount > 0 {
		log.Println("ℹ️  Business data already exists, skipping full seeding.")
		return nil
	}

	log.Println("🚀 Starting comprehensive business data seeding...")

	// --- Equipment Types ---
	cncType := model.EquipmentType{Name: "高精度数控机床", Category: "加工设备"}
	pressType := model.EquipmentType{Name: "全自动冲床", Category: "成型设备"}
	robotType := model.EquipmentType{Name: "工业机器人", Category: "智能装备"}
	if err := mustCreate(db, "EquipmentTypes", []*model.EquipmentType{&cncType, &pressType, &robotType}); err != nil {
		return err
	}

	// --- Inspection Templates & Items ---
	cncTemp := model.InspectionTemplate{Name: "CNC日常点检标准", EquipmentTypeID: cncType.ID}
	if err := mustCreate(db, "CNC InspectionTemplate", &cncTemp); err != nil {
		return err
	}
	cncItems := []model.InspectionItem{
		{TemplateID: cncTemp.ID, Name: "液压压力检查", Method: "观察压力表", Criteria: "4.0-5.0 MPa", SequenceOrder: 1},
		{TemplateID: cncTemp.ID, Name: "主轴温度", Method: "手感/红外测温", Criteria: "≤60℃，无烫感", SequenceOrder: 2},
	}
	if err := mustCreate(db, "CNC InspectionItems", &cncItems); err != nil {
		return err
	}

	pressTemp := model.InspectionTemplate{Name: "冲床周检标准", EquipmentTypeID: pressType.ID}
	if err := mustCreate(db, "Press InspectionTemplate", &pressTemp); err != nil {
		return err
	}
	pressItems := []model.InspectionItem{
		{TemplateID: pressTemp.ID, Name: "离合器间隙", Method: "塞尺测量", Criteria: "0.3-0.5mm", SequenceOrder: 1},
		{TemplateID: pressTemp.ID, Name: "制动器磨损", Method: "目视检查", Criteria: "无明显沟槽", SequenceOrder: 2},
		{TemplateID: pressTemp.ID, Name: "液压油位", Method: "观察油标", Criteria: "在上下限之间", SequenceOrder: 3},
	}
	if err := mustCreate(db, "Press InspectionItems", &pressItems); err != nil {
		return err
	}

	// --- Equipment ---
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
	if err := mustCreate(db, "Equipment", []*model.Equipment{&e1, &e2, &eP, &eR}); err != nil {
		return err
	}

	// --- Spare Parts & Inventory ---
	pump := model.SparePart{Code: "PUMP-01", Name: "高压柱塞泵", FactoryID: &fac.ID, Specification: "Rexroth A10V", Unit: "台", SafetyStock: 2}
	filter := model.SparePart{Code: "FLT-CHEAP", Name: "普通滤芯(降本件)", FactoryID: &fac.ID, Unit: "个", SafetyStock: 50}
	oil := model.SparePart{Code: "OIL-HM46", Name: "液压油 HM46", FactoryID: &fac.ID, Specification: "L-HM46 抗磨液压油", Unit: "桶", SafetyStock: 100}
	filterOEM := model.SparePart{Code: "FLT-OEM", Name: "原厂精密滤芯", FactoryID: &fac.ID, Specification: "Rexroth 原厂", Unit: "个", SafetyStock: 20}
	bearing := model.SparePart{Code: "BRG-6205", Name: "深沟球轴承", FactoryID: &fac.ID, Specification: "SKF 6205-2RS", Unit: "个", SafetyStock: 10}
	seal := model.SparePart{Code: "SEAL-NBR50", Name: "NBR密封圈套装", FactoryID: &fac.ID, Specification: "φ50 标准", Unit: "套", SafetyStock: 15}
	belt := model.SparePart{Code: "BELT-3M", Name: "同步带", FactoryID: &fac.ID, Specification: "3M-GT3-900", Unit: "根", SafetyStock: 10}
	sensor := model.SparePart{Code: "SEN-PROX", Name: "接近开关", FactoryID: &fac.ID, Specification: "OMRON TL-Q5MC1", Unit: "个", SafetyStock: 15}
	relay := model.SparePart{Code: "RLY-24V", Name: "中间继电器", FactoryID: &fac.ID, Specification: "Schneider RXM2AB2BD", Unit: "个", SafetyStock: 30}
	valve := model.SparePart{Code: "VLV-SOL", Name: "电磁阀", FactoryID: &fac.ID, Specification: "SMC SY5120-5G-C6", Unit: "个", SafetyStock: 5}
	tool := model.SparePart{Code: "CNC-TOOL-01", Name: "数控刀头", FactoryID: &fac.ID, Specification: "Sandvik Coromant", Unit: "把", SafetyStock: 10}

	if err := mustCreate(db, "SpareParts", []*model.SparePart{&pump, &filter, &oil, &filterOEM, &bearing, &seal, &belt, &sensor, &relay, &valve, &tool}); err != nil {
		return err
	}
	if err := mustCreate(db, "SparePartInventories", []*model.SparePartInventory{
		{SparePartID: pump.ID, FactoryID: fac.ID, Quantity: 3},
		{SparePartID: filter.ID, FactoryID: fac.ID, Quantity: 120},
		{SparePartID: oil.ID, FactoryID: fac.ID, Quantity: 250},
		{SparePartID: filterOEM.ID, FactoryID: fac.ID, Quantity: 35},
		{SparePartID: bearing.ID, FactoryID: fac.ID, Quantity: 8}, // Below safety stock!
		{SparePartID: seal.ID, FactoryID: fac.ID, Quantity: 40},
		{SparePartID: belt.ID, FactoryID: fac.ID, Quantity: 25},
		{SparePartID: sensor.ID, FactoryID: fac.ID, Quantity: 12},
		{SparePartID: relay.ID, FactoryID: fac.ID, Quantity: 55},
		{SparePartID: valve.ID, FactoryID: fac.ID, Quantity: 8},
		{SparePartID: tool.ID, FactoryID: fac.ID, Quantity: 18},
	}); err != nil {
		return err
	}

	// --- Maintenance Plan + Items ---
	mPlan := model.MaintenancePlan{Name: "CNC二级保养", EquipmentTypeID: cncType.ID, CycleDays: 30, Level: 2, WorkHours: 4.0}
	if err := mustCreate(db, "MaintenancePlan", &mPlan); err != nil {
		return err
	}
	planItems := []model.MaintenancePlanItem{
		{PlanID: mPlan.ID, Name: "检查液压油位", Method: "观察油标", Criteria: "在上下限之间", SequenceOrder: 1},
		{PlanID: mPlan.ID, Name: "检查导轨润滑", Method: "手动试运行", Criteria: "运行顺畅无异响", SequenceOrder: 2},
		{PlanID: mPlan.ID, Name: "检查冷却液浓度", Method: "折光仪测量", Criteria: "5-8%", SequenceOrder: 3},
		{PlanID: mPlan.ID, Name: "清理排屑器", Method: "目视检查", Criteria: "无残留铁屑", SequenceOrder: 4},
		{PlanID: mPlan.ID, Name: "检查电气接线端子", Method: "紧固检查", Criteria: "无松动发热", SequenceOrder: 5},
	}
	if err := mustCreate(db, "MaintenancePlanItems", &planItems); err != nil {
		return err
	}

	// 180 days of maintenance tasks (6 cycles × 2 equipment)
	for i := 1; i <= 6; i++ {
		date := now.AddDate(0, 0, -i*30)
		if err := mustCreate(db, "MaintenanceTask(liSi)", &model.MaintenanceTask{
			PlanID: mPlan.ID, EquipmentID: e1.ID, AssignedTo: liSi.ID, Status: "completed", CompletedAt: &date, ActualHours: 4.2, ScheduledDate: date.Format("2006-01-02"),
		}); err != nil {
			return err
		}
		if err := mustCreate(db, "MaintenanceTask(zs)", &model.MaintenanceTask{
			PlanID: mPlan.ID, EquipmentID: e2.ID, AssignedTo: zs.ID, Status: "completed", CompletedAt: &date, ActualHours: 0.2, ScheduledDate: date.Format("2006-01-02"),
		}); err != nil {
			return err
		}
	}

	// --- Repair Orders (full workflow coverage) ---

	// Closed order 1: CNC-001 主轴轴承异响 (90 days ago)
	cnc1Fail1 := now.AddDate(0, 0, -90)
	cnc1Order1 := model.RepairOrder{
		EquipmentID: e1.ID, Status: model.RepairClosed, Priority: 2, AssignedTo: &liSi.ID, ReporterID: op.ID,
		FaultDescription: "主轴轴承异响，运转时有明显金属摩擦声", Solution: "更换主轴轴承 NSK 7010C，重新调整预紧力",
		StartedAt: &cnc1Fail1, CompletedAt: timePtr(cnc1Fail1.Add(6 * time.Hour)), ClosedAt: timePtr(cnc1Fail1.Add(24 * time.Hour)),
	}
	if err := mustCreate(db, "RepairOrder(cnc1-1)", &cnc1Order1); err != nil {
		return err
	}
	mustCreate(db, "RepairCostDetail(cnc1-1)", &model.RepairCostDetail{OrderID: cnc1Order1.ID, SparePartCost: 6500.0, LaborCost: 2000.0})
	mustCreate(db, "SparePartConsumption(cnc1-1)", &model.SparePartConsumption{SparePartID: bearing.ID, OrderID: &cnc1Order1.ID, Quantity: 2, UserID: liSi.ID})
	createRepairLogs(db, cnc1Order1.ID, liSi.ID, cnc1Fail1)

	// Closed order 2: CNC-001 液压泵渗漏 (45 days ago)
	cnc1Fail2 := now.AddDate(0, 0, -45)
	cnc1Order2 := model.RepairOrder{
		EquipmentID: e1.ID, Status: model.RepairClosed, Priority: 2, AssignedTo: &liSi.ID, ReporterID: op.ID,
		FaultDescription: "液压泵接头处渗漏油液，地面有油渍", Solution: "更换液压泵出口密封圈，紧固管接头",
		StartedAt: &cnc1Fail2, CompletedAt: timePtr(cnc1Fail2.Add(3 * time.Hour)), ClosedAt: timePtr(cnc1Fail2.Add(24 * time.Hour)),
	}
	if err := mustCreate(db, "RepairOrder(cnc1-2)", &cnc1Order2); err != nil {
		return err
	}
	mustCreate(db, "RepairCostDetail(cnc1-2)", &model.RepairCostDetail{OrderID: cnc1Order2.ID, SparePartCost: 1800.0, LaborCost: 1400.0})
	mustCreate(db, "SparePartConsumption(cnc1-2)", &model.SparePartConsumption{SparePartID: seal.ID, OrderID: &cnc1Order2.ID, Quantity: 1, UserID: liSi.ID})
	createRepairLogs(db, cnc1Order2.ID, liSi.ID, cnc1Fail2)

	// In-progress order: CNC-001 冷却液泄漏 (today)
	cnc1Order3 := model.RepairOrder{
		EquipmentID: e1.ID, Status: model.RepairInProgress, Priority: 2, AssignedTo: &liSi.ID, ReporterID: op.ID,
		FaultDescription: "冷却液管路接头处泄漏，冷却效率下降", StartedAt: timePtr(now.Add(-2 * time.Hour)),
	}
	mustCreate(db, "RepairOrder(cnc1-3)", &cnc1Order3)

	// Closed order: CNC-002 泵轴承毁坏 (the original 5.5W incident, 10 days ago)
	cnc2Fail := now.AddDate(0, 0, -10)
	cnc2Order := model.RepairOrder{
		EquipmentID: e2.ID, Status: model.RepairClosed, Priority: 1, AssignedTo: &zs.ID, ReporterID: op.ID,
		FaultDescription: "液压泵轴承毁坏，泵体温度异常升高", Solution: "更换高压柱塞泵及滤芯",
		StartedAt: &cnc2Fail, CompletedAt: timePtr(cnc2Fail.Add(4 * time.Hour)), ClosedAt: timePtr(cnc2Fail.Add(24 * time.Hour)),
	}
	if err := mustCreate(db, "RepairOrder(cnc2)", &cnc2Order); err != nil {
		return err
	}
	mustCreate(db, "RepairCostDetail(cnc2)", &model.RepairCostDetail{OrderID: cnc2Order.ID, SparePartCost: 52000.0, LaborCost: 3000.0})
	mustCreate(db, "SparePartConsumption(cnc2-pump)", &model.SparePartConsumption{SparePartID: pump.ID, OrderID: &cnc2Order.ID, Quantity: 1, UserID: zs.ID})
	mustCreate(db, "SparePartConsumption(cnc2-filter)", &model.SparePartConsumption{SparePartID: filter.ID, OrderID: &cnc2Order.ID, Quantity: 2, UserID: zs.ID})
	createRepairLogs(db, cnc2Order.ID, zs.ID, cnc2Fail)

	// Testing order: CNC-002 主轴异响
	cnc2Testing := model.RepairOrder{
		EquipmentID: e2.ID, Status: model.RepairTesting, Priority: 2, AssignedTo: &liSi.ID, ReporterID: zs.ID,
		FaultDescription: "主轴高速运转时有间歇性异响", Solution: "更换主轴轴承并做动平衡校正",
		StartedAt: timePtr(now.Add(-5 * time.Hour)), CompletedAt: timePtr(now.Add(-1 * time.Hour)),
	}
	mustCreate(db, "RepairOrder(cnc2-testing)", &cnc2Testing)

	// Closed order: PRESS-05 离合器打滑 (30 days ago)
	pressFail1 := now.AddDate(0, 0, -30)
	pressOrder1 := model.RepairOrder{
		EquipmentID: eP.ID, Status: model.RepairClosed, Priority: 2, AssignedTo: &zs.ID, ReporterID: op.ID,
		FaultDescription: "离合器打滑，冲压行程无力", Solution: "更换离合器摩擦片，调整间隙至0.4mm",
		StartedAt: &pressFail1, CompletedAt: timePtr(pressFail1.Add(8 * time.Hour)), ClosedAt: timePtr(pressFail1.Add(48 * time.Hour)),
	}
	if err := mustCreate(db, "RepairOrder(press-1)", &pressOrder1); err != nil {
		return err
	}
	mustCreate(db, "RepairCostDetail(press-1)", &model.RepairCostDetail{OrderID: pressOrder1.ID, SparePartCost: 9500.0, LaborCost: 2500.0})
	createRepairLogs(db, pressOrder1.ID, zs.ID, pressFail1)

	// Closed order: PRESS-05 液压油泄漏 (60 days ago)
	pressFail2 := now.AddDate(0, 0, -60)
	pressOrder2 := model.RepairOrder{
		EquipmentID: eP.ID, Status: model.RepairClosed, Priority: 1, AssignedTo: &zs.ID, ReporterID: op.ID,
		FaultDescription: "液压系统油管破裂，大量液压油泄漏", Solution: "更换液压油管，补充液压油",
		StartedAt: &pressFail2, CompletedAt: timePtr(pressFail2.Add(5 * time.Hour)), ClosedAt: timePtr(pressFail2.Add(24 * time.Hour)),
	}
	if err := mustCreate(db, "RepairOrder(press-2)", &pressOrder2); err != nil {
		return err
	}
	mustCreate(db, "RepairCostDetail(press-2)", &model.RepairCostDetail{OrderID: pressOrder2.ID, SparePartCost: 4800.0, LaborCost: 2000.0})
	mustCreate(db, "SparePartConsumption(press-2)", &model.SparePartConsumption{SparePartID: oil.ID, OrderID: &pressOrder2.ID, Quantity: 2, UserID: zs.ID})
	createRepairLogs(db, pressOrder2.ID, zs.ID, pressFail2)

	// Assigned order: ROBOT-01 伺服报警
	robotOrder := model.RepairOrder{
		EquipmentID: eR.ID, Status: model.RepairAssigned, Priority: 3, AssignedTo: &liSi.ID, ReporterID: op.ID,
		FaultDescription: "伺服电机报警 E501，编码器通信异常",
	}
	mustCreate(db, "RepairOrder(robot)", &robotOrder)

	// --- Inspection Tasks & Records (180 days for CNC-001) ---
	for i := 1; i <= 180; i++ {
		date := now.AddDate(0, 0, -i)
		task := model.InspectionTask{
			EquipmentID: e1.ID, TemplateID: cncTemp.ID, AssignedTo: op.ID,
			ScheduledDate: date, Status: "completed", CompletedAt: timePtr(date.Add(10 * time.Minute)),
		}
		if err := mustCreate(db, "InspectionTask", &task); err != nil {
			return err
		}
		// Create records for each item using actual IDs
		for _, item := range cncItems {
			result := "OK"
			// Every 15 days, simulate an NG record for variety
			if i%15 == 0 && item.Name == "主轴温度" {
				result = "NG"
			}
			mustCreate(db, "InspectionRecord", &model.InspectionRecord{
				TaskID: task.ID, ItemID: item.ID, Result: result,
			})
		}
	}

	// --- Manual Document + Chunks (for Agent retrieval) ---
	manual := model.ManualDocument{
		Title: "CNC数控机床维护手册", EquipmentTypeID: &cncType.ID, FilePath: "/manuals/cnc_maintenance.pdf",
	}
	if err := mustCreate(db, "ManualDocument", &manual); err != nil {
		return err
	}
	chunks := []model.ManualChunk{
		{DocumentID: manual.ID, SectionTitle: "液压系统维护", PageNumber: 12,
			Content: "液压油每2000小时或6个月更换一次，以先到者为准。更换时必须同时更换回油滤芯和吸油滤芯。推荐使用L-HM46抗磨液压油，油温正常工作范围30-55℃。液压泵出口压力应保持在4.0-5.0MPa，低于3.5MPa时需检查泵体磨损情况。每月检查一次油位，油位低于油标下限时需补充同品牌同型号液压油。"},
		{DocumentID: manual.ID, SectionTitle: "主轴维护", PageNumber: 28,
			Content: "主轴轴承采用脂润滑，每运行2000小时补充一次润滑脂（推荐SKF LGMT2）。主轴温升不超过环境温度+25℃，正常运转温度范围30-60℃。当出现异常振动或金属摩擦声时，应立即停机检查轴承状态。主轴轴承设计寿命约15000小时，实际寿命受负载和维护质量影响。预紧力调整需使用专业工具，建议由专业维修人员操作。"},
		{DocumentID: manual.ID, SectionTitle: "冷却系统", PageNumber: 45,
			Content: "冷却液推荐浓度5-8%（折光仪测量），浓度过低会导致工件生锈，过高会刺激皮肤。每2周检查一次浓度并补充清水或原液。冷却液每月更换一次，更换时需彻底清洗水箱和管路。冷却液温度应保持在15-25℃，超过30℃需开启冷却机组。定期检查冷却液管路接头是否有泄漏，发现泄漏及时紧固或更换密封件。"},
	}
	if err := mustCreate(db, "ManualChunks", &chunks); err != nil {
		return err
	}

	// --- Knowledge Articles ---
	articles := []model.KnowledgeArticle{
		{Title: "Rexroth泵早期损坏特征", EquipmentTypeID: &cncType.ID,
			FaultPhenomenon: "压力异常波动，泵体温度升高", CauseAnalysis: "滤芯堵塞导致空化，润滑不足加速磨损",
			Solution: "更换原厂滤芯，定期检测油液清洁度", CreatedBy: admin.ID, SourceType: "expert_summary"},
		{Title: "冲床离合器打滑预防", EquipmentTypeID: &pressType.ID,
			FaultPhenomenon: "行程无力，冲压深度不足，有异响", CauseAnalysis: "摩擦片磨损超限，间隙过大",
			Solution: "定期检查离合器间隙(0.3-0.5mm)，及时更换摩擦片", CreatedBy: admin.ID, SourceType: "repair_record"},
		{Title: "机器人伺服报警处理", EquipmentTypeID: &robotType.ID,
			FaultPhenomenon: "伺服电机报警E501，定位精度下降", CauseAnalysis: "编码器线缆接触不良或编码器本身故障",
			Solution: "检查并重新插拔编码器连接器，必要时更换编码器线缆", CreatedBy: admin.ID, SourceType: "expert_summary"},
		{Title: "液压系统油品选择指南", EquipmentTypeID: nil,
			FaultPhenomenon: "油温异常升高，液压动作迟缓", CauseAnalysis: "油液粘度不匹配环境温度，油品劣化",
			Solution: "按设备手册选用HM46抗磨液压油，定期检测油液酸值和水分含量", CreatedBy: admin.ID, SourceType: "expert_summary"},
	}
	if err := mustCreate(db, "KnowledgeArticles", &articles); err != nil {
		return err
	}

	// --- Agent Brain ---
	skill := model.AgentSkill{
		Name: "级联失效审计", Description: "分析连锁损失风险，评估备件质量和维护合规性对设备健康的影响", Status: "active",
		Steps: `[{"step": 1, "tool": "get_cost_analysis"}, {"step": 2, "tool": "get_maintenance_compliance"}, {"step": 3, "tool": "search_knowledge"}]`,
	}
	if err := mustCreate(db, "AgentSkill", &skill); err != nil {
		return err
	}

	experiences := []model.AgentExperience{
		{UserID: admin.ID, Category: "reporting", Content: "该用户偏好简洁的财务摘要，重点关注 TCO 异常和维护成本占比。", Weight: 1.0, Status: "active"},
		{UserID: admin.ID, Category: "analysis_depth", Content: "该用户希望看到完整数据链路，包括原始记录和计算过程。", Weight: 0.8, Status: "active"},
		{UserID: liSi.ID, Category: "work_preference", Content: "该维修工偏好预防性维护，会主动检查关联部件。", Weight: 0.9, Status: "active"},
	}
	if err := mustCreate(db, "AgentExperiences", &experiences); err != nil {
		return err
	}

	// AgentKnowledge: one confirmed + one draft (for knowledge audit panel)
	nowPtr := timePtr(now)
	knowledges := []model.AgentKnowledge{
		{ID: "k_seed_001", Title: "高压泵级联失效预防", Type: "failure_pattern", Summary: "劣质滤芯诱因分析：使用非原厂滤芯导致过滤精度不足，杂质进入泵体加速轴承磨损，最终导致泵体毁坏。关联损失：停机4小时 + 备件费5.2万 + 人工费3000。",
			Confidence: 0.98, Status: "confirmed", CreatedBy: "admin", VerifiedBy: &admin.ID, VerifiedAt: nowPtr, ReferencedCount: 5},
		{ID: "k_seed_002", Title: "老旧冲床高频故障模式分析", Type: "failure_pattern", Summary: "服役超10年的冲床主要故障模式：离合器摩擦片磨损(35%)、液压系统泄漏(28%)、导轨磨损(20%)。建议制定专项保养计划或评估退役。",
			Confidence: 0.85, Status: "draft", CreatedBy: "system"},
	}
	if err := mustCreate(db, "AgentKnowledges", &knowledges); err != nil {
		return err
	}

	log.Println("✅ Comprehensive seeding complete. Login: admin / admin123")
	return nil
}

// ============================================================
// Helper functions
// ============================================================

func bcryptHash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func findOrCreateBase(db *gorm.DB, code, name string) (*model.Base, error) {
	var base model.Base
	if err := db.Where("code = ?", code).First(&base).Error; err != nil {
		base = model.Base{Code: code, Name: name}
		if err := db.Create(&base).Error; err != nil {
			return nil, fmt.Errorf("create base %s: %w", code, err)
		}
	}
	return &base, nil
}

func findOrCreateFactory(db *gorm.DB, code, name string, baseID uint) (*model.Factory, error) {
	var fac model.Factory
	if err := db.Where("code = ?", code).First(&fac).Error; err != nil {
		fac = model.Factory{BaseID: baseID, Code: code, Name: name}
		if err := db.Create(&fac).Error; err != nil {
			return nil, fmt.Errorf("create factory %s: %w", code, err)
		}
	}
	return &fac, nil
}

func findOrCreateWorkshop(db *gorm.DB, code, name string, factoryID uint) (*model.Workshop, error) {
	var ws model.Workshop
	if err := db.Where("code = ?", code).First(&ws).Error; err != nil {
		ws = model.Workshop{FactoryID: factoryID, Code: code, Name: name}
		if err := db.Create(&ws).Error; err != nil {
			return nil, fmt.Errorf("create workshop %s: %w", code, err)
		}
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
		if err := db.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("create user %s: %w", username, err)
		}
	}
	return &user, nil
}

func createRepairLogs(db *gorm.DB, orderID, userID uint, baseTime time.Time) {
	logs := []model.RepairLog{
		{OrderID: orderID, UserID: userID, Action: "create", Content: "创建维修工单"},
		{OrderID: orderID, UserID: userID, Action: "assign", Content: "分配维修任务"},
		{OrderID: orderID, UserID: userID, Action: "start", Content: "开始维修作业"},
		{OrderID: orderID, UserID: userID, Action: "complete", Content: "维修完成，提交测试"},
		{OrderID: orderID, UserID: userID, Action: "confirm", Content: "测试通过，确认修复"},
		{OrderID: orderID, UserID: userID, Action: "audit", Content: "审核通过，关闭工单"},
	}
	mustCreate(db, "RepairLogs", &logs)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
