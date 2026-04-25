package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/ems/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// Store 内存存储
type Store struct {
	mu sync.RWMutex

	// 组织架构
	Bases     map[uint]*model.Base
	Factories map[uint]*model.Factory
	Workshops map[uint]*model.Workshop

	// 用户
	Users map[uint]*model.User

	// 设备
	EquipmentTypes map[uint]*model.EquipmentType
	Equipment      map[uint]*model.Equipment

	// 点检
	InspectionTemplates map[uint]*model.InspectionTemplate
	InspectionItems     map[uint]*model.InspectionItem
	InspectionTasks     map[uint]*model.InspectionTask
	InspectionRecords   map[uint]*model.InspectionRecord

	// 维修
	RepairOrders map[uint]*model.RepairOrder
	RepairLogs   map[uint]*model.RepairLog

	// 保养
	MaintenancePlans    map[uint]*model.MaintenancePlan
	MaintenanceItems    map[uint]*model.MaintenancePlanItem
	MaintenanceTasks    map[uint]*model.MaintenanceTask
	MaintenanceRecords  map[uint]*model.MaintenanceRecord

	// 备件
	SpareParts          map[uint]*model.SparePart
	SparePartInventory  map[uint]*model.SparePartInventory
	SparePartConsumption map[uint]*model.SparePartConsumption

	// 知识库
	KnowledgeArticles map[uint]*model.KnowledgeArticle

	// Agent 相关
	ManualDocuments    map[uint]*model.ManualDocument
	ManualChunks       map[uint]*model.ManualChunk
	RepairCostDetails  map[uint]*model.RepairCostDetail
	RuntimeSnapshots   map[uint]*model.EquipmentRuntimeSnapshot
	AgentSessions      map[uint]*model.AgentSession
	AgentArtifacts     map[uint]*model.AgentArtifact
	AgentEvidenceLinks map[uint]*model.AgentEvidenceLink

	// ID 计数器
	nextID uint
}

var (
	instance *Store
	once     sync.Once
)

// GetStore 获取存储实例
func GetStore() *Store {
	once.Do(func() {
		instance = &Store{
			Bases:       make(map[uint]*model.Base),
			Factories:   make(map[uint]*model.Factory),
			Workshops:   make(map[uint]*model.Workshop),
			Users:       make(map[uint]*model.User),
			EquipmentTypes: make(map[uint]*model.EquipmentType),
			Equipment:   make(map[uint]*model.Equipment),
			InspectionTemplates: make(map[uint]*model.InspectionTemplate),
			InspectionItems:     make(map[uint]*model.InspectionItem),
			InspectionTasks:     make(map[uint]*model.InspectionTask),
			InspectionRecords:   make(map[uint]*model.InspectionRecord),
			RepairOrders:        make(map[uint]*model.RepairOrder),
			RepairLogs:          make(map[uint]*model.RepairLog),
			MaintenancePlans:    make(map[uint]*model.MaintenancePlan),
			MaintenanceItems:    make(map[uint]*model.MaintenancePlanItem),
			MaintenanceTasks:    make(map[uint]*model.MaintenanceTask),
			MaintenanceRecords:  make(map[uint]*model.MaintenanceRecord),
			SpareParts:          make(map[uint]*model.SparePart),
			SparePartInventory:  make(map[uint]*model.SparePartInventory),
			SparePartConsumption: make(map[uint]*model.SparePartConsumption),
			KnowledgeArticles:   make(map[uint]*model.KnowledgeArticle),
			ManualDocuments:    make(map[uint]*model.ManualDocument),
			ManualChunks:       make(map[uint]*model.ManualChunk),
			RepairCostDetails:  make(map[uint]*model.RepairCostDetail),
			RuntimeSnapshots:   make(map[uint]*model.EquipmentRuntimeSnapshot),
			AgentSessions:      make(map[uint]*model.AgentSession),
			AgentArtifacts:     make(map[uint]*model.AgentArtifact),
			AgentEvidenceLinks: make(map[uint]*model.AgentEvidenceLink),
			nextID:              1,
		}
	})
	return instance
}

// nextIDInternal 生成下一个ID (内部使用，调用者需持有锁)
func (s *Store) nextIDInternal() uint {
	id := s.nextID
	s.nextID++
	return id
}

// NextID 生成下一个ID (公开方法，线程安全)
func (s *Store) NextID() uint {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.nextID
	s.nextID++
	return id
}

// InitMockData 初始化模拟数据
func (s *Store) InitMockData() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	// 创建基地
	base1 := &model.Base{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Code: "BASE-HD",
		Name: "华东基地",
	}
	s.Bases[base1.ID] = base1

	// 创建工厂
	factory1 := &model.Factory{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		BaseID: base1.ID,
		Code:   "FACTORY-01",
		Name:   "第一工厂",
	}
	s.Factories[factory1.ID] = factory1

	// 创建车间
	workshop1 := &model.Workshop{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		FactoryID: factory1.ID,
		Code:      "WORKSHOP-JJ",
		Name:      "机加工车间",
	}
	s.Workshops[workshop1.ID] = workshop1

	workshop2 := &model.Workshop{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		FactoryID: factory1.ID,
		Code:      "WORKSHOP-HJ",
		Name:      "焊接车间",
	}
	s.Workshops[workshop2.ID] = workshop2

	// 创建管理员用户
	adminUser := &model.User{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Username:           "admin",
		PasswordHash:       string(hashedPassword),
		Name:               "系统管理员",
		Role:               model.RoleAdmin,
		Phone:              "13800138000",
		IsActive:           true,
		ApprovalStatus:     model.ApprovalStatusApproved,
		MustChangePassword: false,
		FirstLogin:         false,
	}
	s.Users[adminUser.ID] = adminUser

	// 创建工程师用户
	engineerUser := &model.User{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Username:           "engineer",
		PasswordHash:       string(hashedPassword),
		Name:               "设备工程师",
		Role:               model.RoleEngineer,
		Phone:              "13800138001",
		FactoryID:          &factory1.ID,
		IsActive:           true,
		ApprovalStatus:     model.ApprovalStatusApproved,
		MustChangePassword: false,
		FirstLogin:         false,
	}
	s.Users[engineerUser.ID] = engineerUser

	// 创建维修工
	workerUser := &model.User{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Username:           "worker",
		PasswordHash:       string(hashedPassword),
		Name:               "维修工张三",
		Role:               model.RoleMaintenance,
		Phone:              "13800138002",
		FactoryID:          &factory1.ID,
		IsActive:           true,
		ApprovalStatus:     model.ApprovalStatusApproved,
		MustChangePassword: false,
		FirstLogin:         false,
	}
	s.Users[workerUser.ID] = workerUser

	// 创建操作员
	operatorUser := &model.User{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Username:           "operator",
		PasswordHash:       string(hashedPassword),
		Name:               "操作员李四",
		Role:               model.RoleOperator,
		Phone:              "13800138003",
		FactoryID:          &factory1.ID,
		IsActive:           true,
		ApprovalStatus:     model.ApprovalStatusApproved,
		MustChangePassword: false,
		FirstLogin:         false,
	}
	s.Users[operatorUser.ID] = operatorUser

	// 创建 3 个待审核申请
	pendingApplications := []struct {
		Username string
		Name     string
		Role     model.UserRole
	}{
		{"applicant1", "申请人 A", model.RoleOperator},
		{"applicant2", "申请人 B", model.RoleMaintenance},
		{"applicant3", "申请人 C", model.RoleEngineer},
	}

	for _, app := range pendingApplications {
		pendingUser := &model.User{
			BaseModel: model.BaseModel{
				ID:        s.nextIDInternal(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			Username:           app.Username,
			PasswordHash:       string(hashedPassword),
			Name:               app.Name,
			Role:               app.Role,
			Phone:              "13900139000",
			IsActive:           false,
			ApprovalStatus:     model.ApprovalStatusPending,
			MustChangePassword: true,
			FirstLogin:         true,
		}
		s.Users[pendingUser.ID] = pendingUser
	}

	// 创建设备类型
	cncType := &model.EquipmentType{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:     "数控机床",
		Category: "机加工设备",
	}
	s.EquipmentTypes[cncType.ID] = cncType

	welderType := &model.EquipmentType{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:     "焊接机器人",
		Category: "焊接设备",
	}
	s.EquipmentTypes[welderType.ID] = welderType

	// 创建设备
	for i := 1; i <= 10; i++ {
		equipment := &model.Equipment{
			BaseModel: model.BaseModel{
				ID:        s.nextIDInternal(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:       fmt.Sprintf("EQ-JJ-%03d", i),
			Name:       fmt.Sprintf("数控机床-%d", i),
			TypeID:     cncType.ID,
			WorkshopID: workshop1.ID,
			QRCode:     fmt.Sprintf("QR-EQ-JJ-%03d", i),
			Spec:       "型号:CK6150, 功率:15kW",
			Status:     "running",
		}
		purchaseDate := now.AddDate(-2, 0, 0)
		equipment.PurchaseDate = &purchaseDate
		s.Equipment[equipment.ID] = equipment
	}

	for i := 1; i <= 5; i++ {
		equipment := &model.Equipment{
			BaseModel: model.BaseModel{
				ID:        s.nextIDInternal(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:       fmt.Sprintf("EQ-HJ-%03d", i),
			Name:       fmt.Sprintf("焊接机器人-%d", i),
			TypeID:     welderType.ID,
			WorkshopID: workshop2.ID,
			QRCode:     fmt.Sprintf("QR-EQ-HJ-%03d", i),
			Spec:       "型号:ABB-IRB, 功率:8kW",
			Status:     "running",
		}
		purchaseDate := now.AddDate(-1, 0, 0)
		equipment.PurchaseDate = &purchaseDate
		s.Equipment[equipment.ID] = equipment
	}

	// 创建点检模板
	dailyTemplate := &model.InspectionTemplate{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:            "数控机床日常点检",
		EquipmentTypeID: cncType.ID,
	}
	s.InspectionTemplates[dailyTemplate.ID] = dailyTemplate

	// 创建点检项目
	items := []string{"检查液压油位", "检查冷却液液位", "检查主轴运转", "检查安全防护装置", "清洁设备表面"}
	for _, item := range items {
		inspectionItem := &model.InspectionItem{
			BaseModel: model.BaseModel{
				ID:        s.nextIDInternal(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			TemplateID: dailyTemplate.ID,
			Name:       item,
			Method:     "目视检查",
			Criteria:   "正常/无异常",
		}
		s.InspectionItems[inspectionItem.ID] = inspectionItem
	}

	// 创建一些模拟的点检任务和记录
	for i := 1; i <= 5; i++ {
		scheduledDate := now.AddDate(0, 0, -i)
		task := &model.InspectionTask{
			BaseModel: model.BaseModel{
				ID:        s.nextIDInternal(),
				CreatedAt: scheduledDate,
				UpdatedAt: scheduledDate,
			},
			EquipmentID:    getFirstEquipmentKey(s.Equipment),
			TemplateID:     dailyTemplate.ID,
			AssignedTo:     operatorUser.ID,
			ScheduledDate:  scheduledDate,
			Status:         model.InspectionCompleted,
		}
		startedAt := scheduledDate.Add(10 * time.Minute)
		completedAt := scheduledDate.Add(25 * time.Minute)
		task.StartedAt = &startedAt
		task.CompletedAt = &completedAt
		s.InspectionTasks[task.ID] = task

		// 对应的完成记录
		record := &model.InspectionRecord{
			BaseModel: model.BaseModel{
				ID:        s.nextIDInternal(),
				CreatedAt: scheduledDate,
				UpdatedAt: scheduledDate,
			},
			TaskID: task.ID,
			ItemID: getFirstItemKey(s.InspectionItems),
			Result: "OK",
			Remark: "点检正常",
		}
		s.InspectionRecords[record.ID] = record
	}

	// 创建维修工单
	reportedAt := now.AddDate(0, 0, -2)
	startedAt := reportedAt.Add(30 * time.Minute)
	completedAt := reportedAt.Add(4 * time.Hour)

	repairOrder := &model.RepairOrder{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: reportedAt,
			UpdatedAt: completedAt,
		},
		EquipmentID:       getFirstEquipmentKey(s.Equipment),
		FaultDescription:  "主轴异响，需检查轴承",
		ReporterID:        operatorUser.ID,
		AssignedTo:        &workerUser.ID,
		Status:            model.RepairAudited,
		Priority:          1,
		Solution:          "更换主轴轴承，调整预紧力",
		FaultCode:         "FLT-SPINDLE-001",
		StartedAt:         &startedAt,
		CompletedAt:       &completedAt,
	}
	s.RepairOrders[repairOrder.ID] = repairOrder

	// 创建保养计划
	plan := &model.MaintenancePlan{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:            "数控机床月度保养",
		EquipmentTypeID: cncType.ID,
		Level:           2,
		CycleDays:       30,
		FlexibleDays:    3,
		WorkHours:       4.0,
	}
	s.MaintenancePlans[plan.ID] = plan

	// 创建保养项目
	planItems := []string{"更换液压油", "检查导轨润滑", "检查电气系统连接", "清洁冷却系统", "检查防护装置"}
	for _, item := range planItems {
		planItem := &model.MaintenancePlanItem{
			BaseModel: model.BaseModel{
				ID:        s.nextIDInternal(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			PlanID: plan.ID,
			Name:   item,
			Method: "按操作规程执行",
			Criteria: "功能正常，无异常",
		}
		s.MaintenanceItems[planItem.ID] = planItem
	}

	// 创建备件
	spareParts := []*model.SparePart{
		{
			BaseModel: model.BaseModel{
				ID:        s.nextIDInternal(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:         "SP-001",
			Name:         "主轴轴承",
			Specification: "NSK 6208",
			Unit:         "个",
			SafetyStock:  20,
		},
		{
			BaseModel: model.BaseModel{
				ID:        s.nextIDInternal(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			Code:         "SP-002",
			Name:         "液压油",
			Specification: "ISO VG 46",
			Unit:         "升",
			SafetyStock:  50,
		},
	}
	for _, part := range spareParts {
		s.SpareParts[part.ID] = part

		// 库存
		inventory := &model.SparePartInventory{
			BaseModel: model.BaseModel{
				ID:        s.nextIDInternal(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			SparePartID: part.ID,
			FactoryID:   factory1.ID,
			Quantity:    100,
		}
		s.SparePartInventory[inventory.ID] = inventory
	}

	// 创建知识库文章
	article := &model.KnowledgeArticle{
		BaseModel: model.BaseModel{
			ID:        s.nextIDInternal(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Title:           "主轴异响处理方法",
		EquipmentTypeID: &cncType.ID,
		FaultPhenomenon: "主轴运转时出现异常响声",
		CauseAnalysis:   "轴承润滑不良或轴承损坏",
		Solution:        "1. 检查润滑系统\n2. 清洁轴承座\n3. 必要时更换轴承\n4. 调整预紧力",
		SourceType:      "repair",
		Tags:            []string{"主轴", "轴承", "异响"},
		CreatedBy:       engineerUser.ID,
	}
	s.KnowledgeArticles[article.ID] = article

	// 生成历史数据用于统计展示（过去30天）
	s.generateHistoricalData(30)
}

// generateHistoricalData 生成历史数据用于统计分析
func (s *Store) generateHistoricalData(days int) {
	// 获取操作员和维修工用户ID
	var operatorID, workerID uint
	for _, u := range s.Users {
		if u.Role == model.RoleOperator && operatorID == 0 {
			operatorID = u.ID
		}
		if u.Role == model.RoleMaintenance && workerID == 0 {
			workerID = u.ID
		}
	}
	if operatorID == 0 || workerID == 0 {
		return
	}

	// 获取第一个保养计划ID
	var planID uint
	for id := range s.MaintenancePlans {
		planID = id
		break
	}

	// 获取所有设备ID
	var equipmentIDs []uint
	for id := range s.Equipment {
		equipmentIDs = append(equipmentIDs, id)
	}
	if len(equipmentIDs) == 0 {
		return
	}

	// 获取点检模板ID
	var templateIDs []uint
	for id := range s.InspectionTemplates {
		templateIDs = append(templateIDs, id)
	}
	if len(templateIDs) == 0 {
		return
	}

	// 获取点检项目ID
	var itemIDs []uint
	for id := range s.InspectionItems {
		itemIDs = append(itemIDs, id)
	}
	if len(itemIDs) == 0 {
		return
	}

	now := time.Now()

	// 生成过去N天的数据
	for day := 0; day < days; day++ {
		date := now.AddDate(0, 0, -day)

		// 每天生成点检任务和记录（约80%完成率）
		dailyTaskCount := 8 + (day % 5)
		for i := 0; i < dailyTaskCount; i++ {
			equipmentID := equipmentIDs[i%len(equipmentIDs)]
			templateID := templateIDs[0]

			scheduledDate := date
			var startedAt, completedAt *time.Time
			status := model.InspectionCompleted

			sat := scheduledDate.Add(8*time.Hour + time.Duration(i*10)*time.Minute)
			cat := sat.Add(15*time.Minute + time.Duration(i%5)*time.Minute)
			startedAt = &sat
			completedAt = &cat

			// 20%未完成
			if day%5 == 0 && i >= dailyTaskCount-2 {
				status = model.InspectionPending
				startedAt = nil
				completedAt = nil
			}

			task := &model.InspectionTask{
				BaseModel: model.BaseModel{
					ID:        s.nextIDInternal(),
					CreatedAt: scheduledDate,
					UpdatedAt: scheduledDate,
				},
				EquipmentID:    equipmentID,
				TemplateID:     templateID,
				AssignedTo:     operatorID,
				ScheduledDate:  scheduledDate,
				Status:         status,
				StartedAt:      startedAt,
				CompletedAt:    completedAt,
			}
			s.InspectionTasks[task.ID] = task

			// 已完成的任务生成记录
			if status == model.InspectionCompleted {
				record := &model.InspectionRecord{
					BaseModel: model.BaseModel{
						ID:        s.nextIDInternal(),
						CreatedAt: scheduledDate,
						UpdatedAt: scheduledDate,
					},
					TaskID: task.ID,
					ItemID: itemIDs[i%len(itemIDs)],
					Result: "OK",
					Remark:  "点检正常",
				}
				s.InspectionRecords[record.ID] = record
			}
		}

		// 生成维修工单（约30%的天数有故障）
		repairCount := 0
		if day%3 != 0 {
			repairCount = 1 + (day % 3)
		}
		for i := 0; i < repairCount; i++ {
			equipmentID := equipmentIDs[i%len(equipmentIDs)]
			reportedAt := date.Add(10*time.Hour + time.Duration(i*2)*time.Hour)
			startedAt := reportedAt.Add(30 * time.Minute)
			completedAt := startedAt.Add(time.Duration(2+i) * time.Hour)

			// 计算停机时长（小时）
			downtimeHours := completedAt.Sub(startedAt).Hours()

			repairOrder := &model.RepairOrder{
				BaseModel: model.BaseModel{
					ID:        s.nextIDInternal(),
					CreatedAt: reportedAt,
					UpdatedAt: completedAt,
				},
				EquipmentID:       equipmentID,
				FaultDescription:  "设备故障，需要维修",
				ReporterID:        operatorID,
				AssignedTo:        &workerID,
				Status:            model.RepairAudited,
				Priority:          int(day%3) + 1,
				Solution:          "更换故障部件，恢复正常运行",
				FaultCode:         fmt.Sprintf("FLT-%03d", day%10),
				StartedAt:         &startedAt,
				CompletedAt:       &completedAt,
			}
			s.RepairOrders[repairOrder.ID] = repairOrder

			// 维修日志（在Content字段中记录停机时长）
			repairLog := &model.RepairLog{
				BaseModel: model.BaseModel{
					ID:        s.nextIDInternal(),
					CreatedAt: startedAt,
					UpdatedAt: startedAt,
				},
				OrderID:  repairOrder.ID,
				UserID:   workerID,
				Action:   "开始维修",
				Content:  fmt.Sprintf("停机时长: %.2f小时", downtimeHours),
			}
			s.RepairLogs[repairLog.ID] = repairLog
		}

		// 生成保养任务（约60%的天数有保养）
		if day%5 != 0 && planID > 0 {
			maintenanceCount := 2 + (day % 4)
			for i := 0; i < maintenanceCount; i++ {
				equipmentID := equipmentIDs[i%len(equipmentIDs)]
				scheduledDate := date.Add(14 * time.Hour)

				var startedAt, completedAt *time.Time
				status := model.MaintenanceCompleted

				sat := scheduledDate
				cat := scheduledDate.Add(time.Duration(2+i%3) * time.Hour)
				startedAt = &sat
				completedAt = &cat

				if day%7 == 0 && i == 0 {
					status = model.MaintenancePending
					startedAt = nil
					completedAt = nil
				}

				task := &model.MaintenanceTask{
					BaseModel: model.BaseModel{
						ID:        s.nextIDInternal(),
						CreatedAt: scheduledDate,
						UpdatedAt: scheduledDate,
					},
					EquipmentID:    equipmentID,
					PlanID:         planID,
					ScheduledDate:  scheduledDate.Format("2006-01-02"),
					AssignedTo:     workerID,
					Status:         status,
					StartedAt:      startedAt,
					CompletedAt:    completedAt,
				}
				s.MaintenanceTasks[task.ID] = task

				if status == model.MaintenanceCompleted {
					// 获取保养项目ID
					var maintenanceItemID uint
					for id := range s.MaintenanceItems {
						maintenanceItemID = id
						break
					}

					record := &model.MaintenanceRecord{
						BaseModel: model.BaseModel{
							ID:        s.nextIDInternal(),
							CreatedAt: scheduledDate,
							UpdatedAt: scheduledDate,
						},
						TaskID:     task.ID,
						ItemID:     maintenanceItemID,
						Result:     "合格",
						Remark:     "保养完成",
					}
					s.MaintenanceRecords[record.ID] = record
				}
			}
		}
	}
}

func getFirstEquipmentKey(m map[uint]*model.Equipment) uint {
	for k := range m {
		return k
	}
	return 1
}

func getFirstItemKey(m map[uint]*model.InspectionItem) uint {
	for k := range m {
		return k
	}
	return 1
}

// Close 关闭存储（内存模式无需操作）
func (s *Store) Close() error {
	return nil
}

// AddRepairOrder 添加维修工单
func (s *Store) AddRepairOrder(order *model.RepairOrder) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.RepairOrders[order.ID] = order
}

// ============ Organization 辅助方法 ============

func (s *Store) AddBase(id uint, base *model.Base) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Bases[id] = base
}

func (s *Store) UpdateBase(id uint, fn func(*model.Base)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if base, ok := s.Bases[id]; ok {
		fn(base)
		return true
	}
	return false
}

func (s *Store) DeleteBase(id uint) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Bases[id]; ok {
		delete(s.Bases, id)
		return true
	}
	return false
}

func (s *Store) AddFactory(id uint, factory *model.Factory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Factories[id] = factory
}

func (s *Store) UpdateFactory(id uint, fn func(*model.Factory)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if factory, ok := s.Factories[id]; ok {
		fn(factory)
		return true
	}
	return false
}

func (s *Store) DeleteFactory(id uint) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Factories[id]; ok {
		delete(s.Factories, id)
		return true
	}
	return false
}

func (s *Store) AddWorkshop(id uint, workshop *model.Workshop) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Workshops[id] = workshop
}

func (s *Store) UpdateWorkshop(id uint, fn func(*model.Workshop)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if workshop, ok := s.Workshops[id]; ok {
		fn(workshop)
		return true
	}
	return false
}

func (s *Store) DeleteWorkshop(id uint) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Workshops[id]; ok {
		delete(s.Workshops, id)
		return true
	}
	return false
}

// ============ Equipment 辅助方法 ============

func (s *Store) AddEquipment(id uint, equipment *model.Equipment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Equipment[id] = equipment
}

func (s *Store) FindEquipment(id uint) *model.Equipment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Equipment[id]
}

func (s *Store) UpdateEquipment(id uint, fn func(*model.Equipment)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if equipment, ok := s.Equipment[id]; ok {
		fn(equipment)
		return true
	}
	return false
}

func (s *Store) DeleteEquipment(id uint) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Equipment[id]; ok {
		delete(s.Equipment, id)
		return true
	}
	return false
}

func (s *Store) AddEquipmentType(id uint, et *model.EquipmentType) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.EquipmentTypes[id] = et
}

func (s *Store) UpdateEquipmentType(id uint, fn func(*model.EquipmentType)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if et, ok := s.EquipmentTypes[id]; ok {
		fn(et)
		return true
	}
	return false
}

func (s *Store) DeleteEquipmentType(id uint) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.EquipmentTypes[id]; ok {
		delete(s.EquipmentTypes, id)
		return true
	}
	return false
}

// ============ Inspection 辅助方法 ============

func (s *Store) AddInspectionTemplate(id uint, template *model.InspectionTemplate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.InspectionTemplates[id] = template
}

func (s *Store) FindInspectionTemplate(id uint) *model.InspectionTemplate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.InspectionTemplates[id]
}

func (s *Store) AddInspectionItem(id uint, item *model.InspectionItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.InspectionItems[id] = item
}

func (s *Store) GetInspectionItemsByTemplate(templateID uint) []*model.InspectionItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var items []*model.InspectionItem
	for _, item := range s.InspectionItems {
		if item.TemplateID == templateID {
			items = append(items, item)
		}
	}
	return items
}

func (s *Store) AddInspectionTask(id uint, task *model.InspectionTask) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.InspectionTasks[id] = task
}

func (s *Store) FindInspectionTask(id uint) *model.InspectionTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.InspectionTasks[id]
}

func (s *Store) UpdateInspectionTask(id uint, fn func(*model.InspectionTask)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task, ok := s.InspectionTasks[id]; ok {
		fn(task)
		return true
	}
	return false
}

func (s *Store) AddInspectionRecord(id uint, record *model.InspectionRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.InspectionRecords[id] = record
}

// ============ Repair 辅助方法 ============

func (s *Store) FindRepairOrder(id uint) *model.RepairOrder {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.RepairOrders[id]
}

func (s *Store) UpdateRepairOrder(id uint, fn func(*model.RepairOrder)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if order, ok := s.RepairOrders[id]; ok {
		fn(order)
		return true
	}
	return false
}

// ============ Maintenance 辅助方法 ============

func (s *Store) AddMaintenancePlan(id uint, plan *model.MaintenancePlan) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.MaintenancePlans[id] = plan
}

func (s *Store) AddMaintenanceItem(id uint, item *model.MaintenancePlanItem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.MaintenanceItems[id] = item
}

func (s *Store) AddMaintenanceTask(id uint, task *model.MaintenanceTask) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.MaintenanceTasks[id] = task
}

func (s *Store) FindMaintenanceTask(id uint) *model.MaintenanceTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.MaintenanceTasks[id]
}

func (s *Store) UpdateMaintenanceTask(id uint, fn func(*model.MaintenanceTask)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task, ok := s.MaintenanceTasks[id]; ok {
		fn(task)
		return true
	}
	return false
}

func (s *Store) AddMaintenanceRecord(id uint, record *model.MaintenanceRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.MaintenanceRecords[id] = record
}

// ============ Spare Parts 辅助方法 ============

func (s *Store) AddSparePart(id uint, part *model.SparePart) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SpareParts[id] = part
}

func (s *Store) FindSparePart(id uint) *model.SparePart {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.SpareParts[id]
}

func (s *Store) UpdateSparePart(id uint, fn func(*model.SparePart)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if part, ok := s.SpareParts[id]; ok {
		fn(part)
		return true
	}
	return false
}

func (s *Store) DeleteSparePart(id uint) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.SpareParts[id]; ok {
		delete(s.SpareParts, id)
		return true
	}
	return false
}

func (s *Store) AddSparePartInventory(id uint, inv *model.SparePartInventory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SparePartInventory[id] = inv
}

func (s *Store) AddSparePartConsumption(id uint, cons *model.SparePartConsumption) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SparePartConsumption[id] = cons
}

// ============ Knowledge 辅助方法 ============

func (s *Store) AddKnowledgeArticle(id uint, article *model.KnowledgeArticle) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.KnowledgeArticles[id] = article
}

func (s *Store) FindKnowledgeArticle(id uint) *model.KnowledgeArticle {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.KnowledgeArticles[id]
}

func (s *Store) UpdateKnowledgeArticle(id uint, fn func(*model.KnowledgeArticle)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if article, ok := s.KnowledgeArticles[id]; ok {
		fn(article)
		return true
	}
	return false
}

func (s *Store) DeleteKnowledgeArticle(id uint) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.KnowledgeArticles[id]; ok {
		delete(s.KnowledgeArticles, id)
		return true
	}
	return false
}

// AddUser 添加用户
func (s *Store) AddUser(id uint, user *model.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Users[id] = user
}

// UpdateUser 更新用户（回调函数持有锁）
func (s *Store) UpdateUser(id uint, fn func(*model.User)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if user, ok := s.Users[id]; ok {
		fn(user)
		return true
	}
	return false
}

// GetUsers 获取所有用户
func (s *Store) GetUsers() []*model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var users []*model.User
	for _, u := range s.Users {
		users = append(users, u)
	}
	return users
}

// FindUserByUsername 根据用户名查找用户
func (s *Store) FindUserByUsername(username string) *model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.Users {
		if u.Username == username {
			return u
		}
	}
	return nil
}

// FindUser 根据ID查找用户
func (s *Store) FindUser(id uint) *model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Users[id]
}
