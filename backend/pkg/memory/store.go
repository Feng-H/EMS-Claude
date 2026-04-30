package memory

import (
	"fmt"
	"sync"
	"time"
	"math/rand"

	"github.com/ems/backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	mu     sync.RWMutex
	nextID uint

	Bases            map[uint]*model.Base
	Factories        map[uint]*model.Factory
	Workshops        map[uint]*model.Workshop
	Users            map[uint]*model.User
	EquipmentTypes   map[uint]*model.EquipmentType
	Equipment        map[uint]*model.Equipment
	
	InspectionTemplates map[uint]*model.InspectionTemplate
	InspectionItems     map[uint]*model.InspectionItem
	InspectionTasks     map[uint]*model.InspectionTask
	InspectionRecords   map[uint]*model.InspectionRecord
	
	RepairOrders        map[uint]*model.RepairOrder
	RepairLogs          map[uint]*model.RepairLog
	RepairCostDetails   map[uint]*model.RepairCostDetail
	
	MaintenancePlans    map[uint]*model.MaintenancePlan
	MaintenancePlanItems map[uint]*model.MaintenancePlanItem
	MaintenanceTasks    map[uint]*model.MaintenanceTask
	MaintenanceRecords  map[uint]*model.MaintenanceRecord
	
	SpareParts            map[uint]*model.SparePart
	SparePartInventory    map[uint]*model.SparePartInventory
	SparePartConsumption  map[uint]*model.SparePartConsumption
	SparePartTransactions map[uint]*model.SparePartTransaction
	
	KnowledgeArticles    map[uint]*model.KnowledgeArticle
	ManualDocuments      map[uint]*model.ManualDocument
	ManualChunks         map[uint]*model.ManualChunk
	
	AgentSkills           map[uint]*model.AgentSkill
	AgentKnowledges       map[string]*model.AgentKnowledge
	AgentExperiences      map[uint]*model.AgentExperience
	AgentConversations    map[uint]*model.AgentConversation
	AgentMessages         map[uint]*model.AgentMessage
	AgentPushSubscriptions map[uint]*model.AgentPushSubscription
	AgentUsages           map[uint]*model.AgentUsage
	AgentArtifacts        map[uint]*model.AgentArtifact
	AgentEvidenceLinks    map[uint]*model.AgentEvidenceLink
	AgentSessions         map[uint]*model.AgentSession
	
	RuntimeSnapshots      map[uint]*model.EquipmentRuntimeSnapshot
}

var (
	instance *Store
	once     sync.Once
)

func GetStore() *Store {
	once.Do(func() {
		instance = &Store{
			nextID: 1,
			Bases:            make(map[uint]*model.Base),
			Factories:        make(map[uint]*model.Factory),
			Workshops:        make(map[uint]*model.Workshop),
			Users:            make(map[uint]*model.User),
			EquipmentTypes:   make(map[uint]*model.EquipmentType),
			Equipment:        make(map[uint]*model.Equipment),
			InspectionTemplates: make(map[uint]*model.InspectionTemplate),
			InspectionItems:     make(map[uint]*model.InspectionItem),
			InspectionTasks:     make(map[uint]*model.InspectionTask),
			InspectionRecords:   make(map[uint]*model.InspectionRecord),
			RepairOrders:        make(map[uint]*model.RepairOrder),
			RepairLogs:          make(map[uint]*model.RepairLog),
			RepairCostDetails:   make(map[uint]*model.RepairCostDetail),
			MaintenancePlans:    make(map[uint]*model.MaintenancePlan),
			MaintenancePlanItems: make(map[uint]*model.MaintenancePlanItem),
			MaintenanceTasks:    make(map[uint]*model.MaintenanceTask),
			MaintenanceRecords:  make(map[uint]*model.MaintenanceRecord),
			SpareParts:           make(map[uint]*model.SparePart),
			SparePartInventory:   make(map[uint]*model.SparePartInventory),
			SparePartConsumption: make(map[uint]*model.SparePartConsumption),
			SparePartTransactions: make(map[uint]*model.SparePartTransaction),
			KnowledgeArticles:    make(map[uint]*model.KnowledgeArticle),
			ManualDocuments:      make(map[uint]*model.ManualDocument),
			ManualChunks:         make(map[uint]*model.ManualChunk),
			AgentSkills:           make(map[uint]*model.AgentSkill),
			AgentKnowledges:       make(map[string]*model.AgentKnowledge),
			AgentExperiences:      make(map[uint]*model.AgentExperience),
			AgentConversations:    make(map[uint]*model.AgentConversation),
			AgentMessages:         make(map[uint]*model.AgentMessage),
			AgentPushSubscriptions: make(map[uint]*model.AgentPushSubscription),
			AgentUsages:           make(map[uint]*model.AgentUsage),
			AgentArtifacts:        make(map[uint]*model.AgentArtifact),
			AgentEvidenceLinks:    make(map[uint]*model.AgentEvidenceLink),
			AgentSessions:         make(map[uint]*model.AgentSession),
			RuntimeSnapshots:      make(map[uint]*model.EquipmentRuntimeSnapshot),
		}
	})
	return instance
}

func (s *Store) nextIDInternal() uint {
	id := s.nextID
	s.nextID++
	return id
}

func timePtr(t time.Time) *time.Time { return &t }

func (s *Store) InitMockData() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	rnd := rand.New(rand.NewSource(now.UnixNano()))
	hp, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	// Organization
	base := &model.Base{BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: now}, Code: "BASE-HQ", Name: "集团总部基地"}
	s.Bases[base.ID] = base
	fac := &model.Factory{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, BaseID: base.ID, Code: "FAC-SZ", Name: "苏州智能工厂"}
	s.Factories[fac.ID] = fac
	ws1 := &model.Workshop{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, FactoryID: fac.ID, Code: "WS-MCH", Name: "精密机加车间"}
	s.Workshops[ws1.ID] = ws1
	ws2 := &model.Workshop{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, FactoryID: fac.ID, Code: "WS-ASM", Name: "全自动装配车间"}
	s.Workshops[ws2.ID] = ws2
	workshops := []*model.Workshop{ws1, ws2}

	// Users
	admin := &model.User{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Username: "admin", PasswordHash: string(hp), Name: "系统管理员", Role: model.RoleAdmin, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
	s.Users[admin.ID] = admin

	var operators []*model.User
	for i := 1; i <= 10; i++ {
		u := &model.User{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Username: fmt.Sprintf("operator_%d", i), PasswordHash: string(hp), Name: fmt.Sprintf("操作工%d号", i), Role: model.RoleOperator, FactoryID: &fac.ID, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
		s.Users[u.ID] = u
		operators = append(operators, u)
	}

	var maints []*model.User
	for i := 1; i <= 5; i++ {
		u := &model.User{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Username: fmt.Sprintf("maint_%d", i), PasswordHash: string(hp), Name: fmt.Sprintf("维修工%d号", i), Role: model.RoleMaintenance, FactoryID: &fac.ID, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
		s.Users[u.ID] = u
		maints = append(maints, u)
	}

	for i := 1; i <= 2; i++ {
		u := &model.User{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Username: fmt.Sprintf("eng_%d", i), PasswordHash: string(hp), Name: fmt.Sprintf("工程师%d号", i), Role: model.RoleEngineer, FactoryID: &fac.ID, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
		s.Users[u.ID] = u
	}

	// 5 Equipment Types and Templates
	typeNames := []string{"数控机床", "全自动冲床", "工业机器人", "注塑机", "自动化流水线"}
	var equipTypes []*model.EquipmentType
	tempMap := make(map[uint]*model.InspectionTemplate)
	planMap := make(map[uint]*model.MaintenancePlan)

	for _, name := range typeNames {
		et := &model.EquipmentType{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Name: name, Category: "生产设备"}
		s.EquipmentTypes[et.ID] = et
		equipTypes = append(equipTypes, et)

		inspTemp := &model.InspectionTemplate{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Name: name + "日常点检", EquipmentTypeID: et.ID, EquipmentType: et}
		s.InspectionTemplates[inspTemp.ID] = inspTemp
		tempMap[et.ID] = inspTemp

		it1 := &model.InspectionItem{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, TemplateID: inspTemp.ID, Name: "外观检查", Method: "目视", Criteria: "无破损", SequenceOrder: 1}
		it2 := &model.InspectionItem{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, TemplateID: inspTemp.ID, Name: "运行状态", Method: "耳听/目视", Criteria: "平稳无异响", SequenceOrder: 2}
		it3 := &model.InspectionItem{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, TemplateID: inspTemp.ID, Name: "关键参数", Method: "仪表读取", Criteria: "在标准范围内", SequenceOrder: 3}
		s.InspectionItems[it1.ID] = it1
		s.InspectionItems[it2.ID] = it2
		s.InspectionItems[it3.ID] = it3
		inspTemp.Items = []model.InspectionItem{*it1, *it2, *it3}

		maintPlan := &model.MaintenancePlan{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Name: name + "月度保养", EquipmentTypeID: et.ID, EquipmentType: et, CycleDays: 30, Level: 1, WorkHours: 2.0}
		s.MaintenancePlans[maintPlan.ID] = maintPlan
		planMap[et.ID] = maintPlan

		mi1 := &model.MaintenancePlanItem{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, PlanID: maintPlan.ID, Name: "全面清洁", Method: "人工", Criteria: "无油污、灰尘", SequenceOrder: 1}
		mi2 := &model.MaintenancePlanItem{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, PlanID: maintPlan.ID, Name: "加注润滑油", Method: "注油", Criteria: "油位正常", SequenceOrder: 2}
		mi3 := &model.MaintenancePlanItem{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, PlanID: maintPlan.ID, Name: "紧固螺栓", Method: "工具", Criteria: "无松动", SequenceOrder: 3}
		s.MaintenancePlanItems[mi1.ID] = mi1
		s.MaintenancePlanItems[mi2.ID] = mi2
		s.MaintenancePlanItems[mi3.ID] = mi3
		maintPlan.Items = []model.MaintenancePlanItem{*mi1, *mi2, *mi3}
	}

	// 20+ Spare Parts
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
	var spareParts []*model.SparePart
	for i, name := range partNames {
		sp := &model.SparePart{
			BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Code: fmt.Sprintf("SP-%03d", i+1), Name: name, 
			FactoryID: &fac.ID, Unit: "件", SafetyStock: rnd.Intn(20) + 5,
		}
		s.SpareParts[sp.ID] = sp
		spareParts = append(spareParts, sp)

		qty := rnd.Intn(100) + 50
		inv := &model.SparePartInventory{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, SparePartID: sp.ID, SparePart: sp, FactoryID: fac.ID, Factory: fac, Quantity: qty}
		s.SparePartInventory[inv.ID] = inv

		tx := &model.SparePartTransaction{BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: now.AddDate(0, -3, 0)}, SparePartID: sp.ID, SparePart: sp, FactoryID: fac.ID, Factory: fac, Type: "in", Quantity: qty, OperatorID: admin.ID, Operator: admin, Remark: "初始大批量入库"}
		s.SparePartTransactions[tx.ID] = tx
	}

	// 100 Equipments
	var equipments []*model.Equipment
	for i := 1; i <= 100; i++ {
		eType := equipTypes[rnd.Intn(len(equipTypes))]
		maint := maints[rnd.Intn(len(maints))]
		ws := workshops[rnd.Intn(len(workshops))]
		status := "running"
		if rnd.Float32() < 0.05 { status = "stopped" } else if rnd.Float32() < 0.1 { status = "maintenance" }

		eq := &model.Equipment{
			BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Code: fmt.Sprintf("EQ-%04d", i), Name: fmt.Sprintf("%s-%03d", eType.Name, i),
			TypeID: eType.ID, Type: eType, WorkshopID: ws.ID, Workshop: ws, QRCode: fmt.Sprintf("EQ-%04d", i),
			PurchasePrice: float64(rnd.Intn(500000) + 50000), PurchaseDate: timePtr(now.AddDate(-rnd.Intn(5), -rnd.Intn(12), 0)),
			ServiceLifeYears: 10, ScrapValue: 5000.0, HourlyLoss: float64(rnd.Intn(400) + 50),
			Status: status, DedicatedMaintenanceID: &maint.ID,
		}
		s.Equipment[eq.ID] = eq
		equipments = append(equipments, eq)
	}

	// History
	historyDays := 15
	for d := historyDays; d >= 0; d-- {
		date := now.AddDate(0, 0, -d)
		dateStr := date.Format("2006-01-02")
		
		for _, eq := range equipments {
			op := operators[rnd.Intn(len(operators))]
			maint := maints[rnd.Intn(len(maints))]
// Inspection
temp := tempMap[eq.TypeID]
iTask := &model.InspectionTask{
	BaseModel: model.BaseModel{ID: s.nextIDInternal()}, EquipmentID: eq.ID, Equipment: eq, TemplateID: temp.ID, Template: temp, AssignedTo: op.ID, Assignee: op,
	ScheduledDate: date, Status: "completed", CompletedAt: timePtr(date.Add(time.Minute * time.Duration(rnd.Intn(30)+5))),
}
s.InspectionTasks[iTask.ID] = iTask

			hasNG := false
			var taskRecords []model.InspectionRecord
			for _, item := range temp.Items {
				res := "OK"
				if rnd.Float32() < 0.03 { res = "NG"; hasNG = true }
				rec := &model.InspectionRecord{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, TaskID: iTask.ID, ItemID: item.ID, Item: &item, Result: res}
				s.InspectionRecords[rec.ID] = rec
				taskRecords = append(taskRecords, *rec)
			}
			iTask.Records = taskRecords

			if d == 15 || d == 0 {
				mPlan := planMap[eq.TypeID]
				mTask := &model.MaintenanceTask{
					BaseModel: model.BaseModel{ID: s.nextIDInternal()}, PlanID: mPlan.ID, Plan: mPlan, EquipmentID: eq.ID, Equipment: eq, AssignedTo: maint.ID, Assignee: maint,
					Status: "completed", CompletedAt: timePtr(date.Add(time.Hour * 2)), ActualHours: float64(rnd.Intn(3)+1), ScheduledDate: dateStr,
				}
				s.MaintenanceTasks[mTask.ID] = mTask
				var mTaskRecords []model.MaintenanceRecord
				for _, item := range mPlan.Items {
					rec := &model.MaintenanceRecord{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, TaskID: mTask.ID, ItemID: item.ID, Item: &item, Result: "OK"}
					s.MaintenanceRecords[rec.ID] = rec
					mTaskRecords = append(mTaskRecords, *rec)
				}
				mTask.Records = mTaskRecords
			}

			if hasNG {
				ro := &model.RepairOrder{
					BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: date}, EquipmentID: eq.ID, Equipment: eq, Status: model.RepairClosed, Priority: rnd.Intn(3)+1,
					AssignedTo: &maint.ID, Assignee: maint, ReporterID: op.ID, Reporter: op, FaultDescription: "日常点检发现异常，需要紧急修复",
					StartedAt: timePtr(date.Add(time.Hour)), Solution: "修复完成并已恢复生产",
					CompletedAt: timePtr(date.Add(time.Hour * 4)), ClosedAt: timePtr(date.Add(time.Hour * 6)),
				}
				s.RepairOrders[ro.ID] = ro

				cd := &model.RepairCostDetail{ID: s.nextIDInternal(), OrderID: ro.ID, SparePartCost: float64(rnd.Intn(2000)), LaborCost: float64(rnd.Intn(1000))}
				s.RepairCostDetails[cd.ID] = cd

				for i := 0; i < rnd.Intn(2)+1; i++ {
					sp := spareParts[rnd.Intn(len(spareParts))]
					qty := rnd.Intn(2)+1
					cons := &model.SparePartConsumption{BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: *ro.CompletedAt}, SparePartID: sp.ID, SparePart: sp, OrderID: &ro.ID, Quantity: qty, UserID: maint.ID, User: maint}
					s.SparePartConsumption[cons.ID] = cons
					tx := &model.SparePartTransaction{BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: *ro.CompletedAt}, SparePartID: sp.ID, SparePart: sp, FactoryID: fac.ID, Factory: fac, Type: "out", Quantity: -qty, OperatorID: maint.ID, Operator: maint, RelatedID: &ro.ID, Remark: "维修消耗"}
					s.SparePartTransactions[tx.ID] = tx
				}
				
				l1 := &model.RepairLog{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, OrderID: ro.ID, UserID: maint.ID, Action: "create", Content: "自动创建"}
				l2 := &model.RepairLog{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, OrderID: ro.ID, UserID: maint.ID, Action: "close", Content: "关闭工单"}
				s.RepairLogs[l1.ID] = l1
				s.RepairLogs[l2.ID] = l2
			}
		}
	}

	// Knowledge Articles
	k1 := &model.KnowledgeArticle{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Title: "高频故障：传感器失效分析", EquipmentTypeID: &equipTypes[0].ID, EquipmentType: equipTypes[0], FaultPhenomenon: "信号丢失", CauseAnalysis: "油污干扰", Solution: "加装防护罩", CreatedBy: admin.ID, SourceType: "expert_summary"}
	s.KnowledgeArticles[k1.ID] = k1

	manual := &model.ManualDocument{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Title: "通用生产设备维护手册", EquipmentTypeID: &equipTypes[0].ID}
	s.ManualDocuments[manual.ID] = manual
	chunk := &model.ManualChunk{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, DocumentID: manual.ID, SectionTitle: "通用维护", Content: "定期检查和润滑是保持设备长寿命的关键...", PageNumber: 1}
	s.ManualChunks[chunk.ID] = chunk

	sk1 := &model.AgentSkill{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Name: "生产效率分析", Status: "active", Description: "分析停机损失", Steps: `[{"step": 1, "tool": "get_cost_analysis"}]`}
	s.AgentSkills[sk1.ID] = sk1
}

// ============ Repository Helpers ============
func (s *Store) NextID() uint { return s.nextIDInternal() }
func (s *Store) AddUser(id uint, u *model.User) { s.mu.Lock(); defer s.mu.Unlock(); s.Users[id] = u }
func (s *Store) FindUser(id uint) *model.User { s.mu.RLock(); defer s.mu.RUnlock(); return s.Users[id] }
func (s *Store) FindUserByUsername(u string) *model.User {
	s.mu.RLock(); defer s.mu.RUnlock()
	for _, user := range s.Users { if user.Username == u { return user } }; return nil
}
func (s *Store) FindEquipment(id uint) *model.Equipment { s.mu.RLock(); defer s.mu.RUnlock(); return s.Equipment[id] }
func (s *Store) AddEquipment(id uint, e *model.Equipment) { s.mu.Lock(); defer s.mu.Unlock(); s.Equipment[id] = e }
func (s *Store) AddKnowledgeArticle(id uint, a *model.KnowledgeArticle) { s.mu.Lock(); defer s.mu.Unlock(); s.KnowledgeArticles[id] = a }
func (s *Store) UpdateUser(id uint, fn func(*model.User)) bool {
	s.mu.Lock(); defer s.mu.Unlock(); if u, ok := s.Users[id]; ok { fn(u); return true }; return false
}
func (s *Store) CreateConversation(c *model.AgentConversation) error {
	s.mu.Lock(); defer s.mu.Unlock(); c.ID = s.nextIDInternal(); c.CreatedAt = time.Now(); s.AgentConversations[c.ID] = c; return nil
}
func (s *Store) CreateMessage(m *model.AgentMessage) error {
	s.mu.Lock(); defer s.mu.Unlock(); m.ID = s.nextIDInternal(); m.CreatedAt = time.Now(); s.AgentMessages[m.ID] = m; return nil
}
func (s *Store) CreateArtifact(a *model.AgentArtifact) error {
	s.mu.Lock(); defer s.mu.Unlock(); a.ID = s.nextIDInternal(); a.CreatedAt = time.Now(); s.AgentArtifacts[a.ID] = a; return nil
}
func (s *Store) CreateUsage(u *model.AgentUsage) error {
	s.mu.Lock(); defer s.mu.Unlock(); u.ID = s.nextIDInternal(); u.CreatedAt = time.Now(); s.AgentUsages[u.ID] = u; return nil
}
func (s *Store) CreateKnowledge(k *model.AgentKnowledge) error {
	s.mu.Lock(); defer s.mu.Unlock(); if k.ID == "" { k.ID = fmt.Sprintf("k_%d", s.nextIDInternal()) }; k.CreatedAt = time.Now(); s.AgentKnowledges[k.ID] = k; return nil
}
func (s *Store) CreateSkill(sk *model.AgentSkill) error {
	s.mu.Lock(); defer s.mu.Unlock(); sk.ID = s.nextIDInternal(); sk.CreatedAt = time.Now(); s.AgentSkills[sk.ID] = sk; return nil
}
func (s *Store) Close() error { return nil }
