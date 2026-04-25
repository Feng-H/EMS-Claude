package memory

import (
	"fmt"
	"math/rand"
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
	MaintenancePlans   map[uint]*model.MaintenancePlan
	MaintenanceItems   map[uint]*model.MaintenancePlanItem
	MaintenanceTasks   map[uint]*model.MaintenanceTask
	MaintenanceRecords map[uint]*model.MaintenanceRecord

	// 备件
	SpareParts           map[uint]*model.SparePart
	SparePartInventory   map[uint]*model.SparePartInventory
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
	AgentUsage         map[uint]*model.AgentUsage

	// Agent Phase 2 相关
	AgentSkills           map[uint]*model.AgentSkill
	AgentKnowledges       map[string]*model.AgentKnowledge
	AgentExperiences      map[uint]*model.AgentExperience
	AgentConversations    map[uint]*model.AgentConversation
	AgentMessages         map[uint]*model.AgentMessage
	AgentPushSubscriptions map[uint]*model.AgentPushSubscription

	// ID 计数器
	nextID uint
}

var (
	instance *Store
	once     sync.Once
)

// GetStore 获取 Store 单例
func GetStore() *Store {
	once.Do(func() {
		instance = &Store{
			Bases:               make(map[uint]*model.Base),
			Factories:           make(map[uint]*model.Factory),
			Workshops:           make(map[uint]*model.Workshop),
			Users:               make(map[uint]*model.User),
			EquipmentTypes:      make(map[uint]*model.EquipmentType),
			Equipment:           make(map[uint]*model.Equipment),
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
			SpareParts:           make(map[uint]*model.SparePart),
			SparePartInventory:   make(map[uint]*model.SparePartInventory),
			SparePartConsumption: make(map[uint]*model.SparePartConsumption),
			KnowledgeArticles:   make(map[uint]*model.KnowledgeArticle),
			ManualDocuments:    make(map[uint]*model.ManualDocument),
			ManualChunks:       make(map[uint]*model.ManualChunk),
			RepairCostDetails:  make(map[uint]*model.RepairCostDetail),
			RuntimeSnapshots:   make(map[uint]*model.EquipmentRuntimeSnapshot),
			AgentSessions:      make(map[uint]*model.AgentSession),
			AgentArtifacts:     make(map[uint]*model.AgentArtifact),
			AgentEvidenceLinks: make(map[uint]*model.AgentEvidenceLink),
			AgentUsage:         make(map[uint]*model.AgentUsage),
			AgentSkills:           make(map[uint]*model.AgentSkill),
			AgentKnowledges:       make(map[string]*model.AgentKnowledge),
			AgentExperiences:      make(map[uint]*model.AgentExperience),
			AgentConversations:    make(map[uint]*model.AgentConversation),
			AgentMessages:         make(map[uint]*model.AgentMessage),
			AgentPushSubscriptions: make(map[uint]*model.AgentPushSubscription),
			nextID:              1,
		}
	})
	return instance
}

// NextID 生成下一个ID (内部私有方法，调用者应已持有锁)
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

// InitMockData 初始化模拟数据 (工业世界模拟器版本)
func timePtr(t time.Time) *time.Time {
	return &t
}

func (s *Store) InitMockData() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	fmt.Printf("DEBUG: Admin password hash generated: %s\n", string(hashedPassword))
	r := rand.New(rand.NewSource(now.UnixNano()))

	// 1. 基础架构
	base := &model.Base{BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: now}, Code: "BASE-HQ", Name: "集团总部基地"}
	s.Bases[base.ID] = base
	factory := &model.Factory{BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: now}, BaseID: base.ID, Code: "FAC-01", Name: "数字化示范工厂"}
	s.Factories[factory.ID] = factory
	workshop := &model.Workshop{BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: now}, FactoryID: factory.ID, Code: "WS-01", Name: "精益加工车间"}
	s.Workshops[workshop.ID] = workshop

	// 2. 核心人物
	admin := &model.User{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Username: "admin", PasswordHash: string(hashedPassword), Name: "系统管理员", Role: model.RoleAdmin, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
	s.Users[admin.ID] = admin
	
	workerLi := &model.User{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Username: "maintenance", PasswordHash: string(hashedPassword), Name: "预防型-李四", Role: model.RoleMaintenance, FactoryID: &factory.ID, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
	s.Users[workerLi.ID] = workerLi
	
	workerZhang := &model.User{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Username: "zhang_hero", PasswordHash: string(hashedPassword), Name: "救火队员-张三", Role: model.RoleMaintenance, FactoryID: &factory.ID, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
	s.Users[workerZhang.ID] = workerZhang

	operator := &model.User{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Username: "operator", PasswordHash: string(hashedPassword), Name: "操作员小王", Role: model.RoleOperator, FactoryID: &factory.ID, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
	s.Users[operator.ID] = operator

	// 3. 设备与备件
	cncType := &model.EquipmentType{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Name: "高精度数控机床", Category: "加工设备"}
	s.EquipmentTypes[cncType.ID] = cncType
	
	equipA := &model.Equipment{
		BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Code: "CNC-001", Name: "李四负责-A区机床", TypeID: cncType.ID, 
		WorkshopID: workshop.ID, Status: "running", QRCode: "QR-A",
		PurchasePrice: 280000.00, PurchaseDate: timePtr(now.AddDate(-3, 0, 0)), ServiceLifeYears: 8, ScrapValue: 28000.00, HourlyLoss: 150.00,
	}
	s.Equipment[equipA.ID] = equipA
	
	equipB := &model.Equipment{
		BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Code: "CNC-002", Name: "张三负责-B区机床", TypeID: cncType.ID, 
		WorkshopID: workshop.ID, Status: "stopped", QRCode: "QR-B",
		PurchasePrice: 280000.00, PurchaseDate: timePtr(now.AddDate(-3, 0, 0)), ServiceLifeYears: 8, ScrapValue: 28000.00, HourlyLoss: 150.00,
	}
	s.Equipment[equipB.ID] = equipB

	agingEquip := &model.Equipment{
		BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Code: "PRESS-05", Name: "12年老旧冲床", TypeID: cncType.ID, 
		WorkshopID: workshop.ID, Status: "maintenance", QRCode: "QR-OLD",
		PurchasePrice: 150000.00, PurchaseDate: timePtr(now.AddDate(-12, 0, 0)), ServiceLifeYears: 10, ScrapValue: 5000.00, HourlyLoss: 80.00,
	}
	s.Equipment[agingEquip.ID] = agingEquip

	pumpPart := &model.SparePart{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Code: "PUMP-01", Name: "高压柱塞泵", Unit: "台", SafetyStock: 2}
	s.SpareParts[pumpPart.ID] = pumpPart
	cheapFilter := &model.SparePart{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Code: "FLT-CHEAP", Name: "普通滤芯(降本件)", Unit: "个", SafetyStock: 100}
	s.SpareParts[cheapFilter.ID] = cheapFilter

	// 4. 知识库专家系统
	kb := &model.KnowledgeArticle{
		BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: now.AddDate(0, -5, 0)},
		Title: "关于精密机床液压系统级联失效的风险预警",
		EquipmentTypeID: &cncType.ID,
		FaultPhenomenon: "油压不稳，高压泵频繁异响并烧蚀。",
		CauseAnalysis: "根本原因在于使用了低过滤精度的非原厂滤芯。在长期超负荷运行时，廉价滤芯易产生微米级金属屑透过，直接导致价值数万元的高压泵在2周内报废。",
		Solution: "必须使用过滤精度<5微米的增强型滤芯。一旦泵损，必须同步化验油质并更换全套密封件。",
		CreatedBy: admin.ID,
	}
	s.KnowledgeArticles[kb.ID] = kb

	// 5. 开启 180 天模拟循环 (半年数据)
	for i := 180; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		
		// --- 场景 A: 李四的预防性模式 ---
		if i % 30 == 0 { // 严格月保
			s.simulateDetailedMaintenance(equipA.ID, workerLi.ID, date, 5.0 + r.Float64())
		}
		// 故障概率极低，且随着保养深度的积累递减
		if r.Float64() < 0.02 { 
			s.simulateRealisticRepair(equipA.ID, workerLi.ID, date, 800+r.Float64()*400, "传感器误报/微调")
		}

		// --- 场景 B: 张三的救火模式 ---
		if i % 30 == 0 {
			if i % 60 == 0 { // 经常跳过保养或草草了事
				s.simulateDetailedMaintenance(equipB.ID, workerZhang.ID, date, 1.2 + r.Float64()*0.5)
			}
		}
		// 故障概率逐渐升高，最终爆发大事故
		prob := 0.05 + (1.0 - float64(i)/180.0)*0.3 
		if r.Float64() < prob {
			cost := 1500.0 + r.Float64()*1000.0
			desc := "一般液压报警"
			if i < 30 { // 最近一个月爆发级联失效
				cost = 45000.0 + r.Float64()*10000.0
				desc = "高压泵抱死毁坏"
			}
			s.simulateRealisticRepair(equipB.ID, workerZhang.ID, date, cost, desc)
		}

		// --- 场景 C: PRESS-05 资产老化轨迹 ---
		if i % 20 == 0 {
			// 故障修复时间逐渐拉长，体现配件难寻
			repairTime := 2.0 + (1.0 - float64(i)/180.0)*12.0 
			cost := 4000.0 * (1.2 + (1.0 - float64(i)/180.0)*5.0)
			s.simulateAgingRepair(agingEquip.ID, workerLi.ID, date, cost, repairTime)
		}

		// 生成高负载运行快照 (揭示压力)
		s.RuntimeSnapshots[s.nextIDInternal()] = &model.EquipmentRuntimeSnapshot{
			EquipmentID: equipB.ID, SnapshotDate: date.Format("2006-01-02"),
			RuntimeHours: 22.5, LoadRate: 1.12, // 长期 112% 负荷
		}
	}
}

func (s *Store) simulateDetailedMaintenance(equipID, workerID uint, date time.Time, hours float64) {
	t := &model.MaintenanceTask{
		BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: date},
		EquipmentID: equipID, AssignedTo: workerID, Status: model.MaintenanceCompleted,
		ScheduledDate: date.Format("2006-01-02"), ActualHours: hours,
		Remark: "已执行深度检查和油质分析",
	}
	s.MaintenanceTasks[t.ID] = t
}

func (s *Store) simulateRealisticRepair(equipID, workerID uint, date time.Time, totalCost float64, desc string) {
	order := &model.RepairOrder{
		BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: date},
		EquipmentID: equipID, FaultDescription: desc, Status: model.RepairAudited,
		Solution: "现场紧急处置", StartedAt: &date, CompletedAt: &date,
	}
	s.RepairOrders[order.ID] = order
	s.RepairCostDetails[s.nextIDInternal()] = &model.RepairCostDetail{
		OrderID: order.ID, SparePartCost: totalCost * 0.85, LaborCost: totalCost * 0.15,
	}
}

func (s *Store) simulateAgingRepair(equipID, workerID uint, date time.Time, cost, hours float64) {
	order := &model.RepairOrder{
		BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: date},
		EquipmentID: equipID, FaultDescription: "老化性失效/配件停产", Status: model.RepairAudited,
		StartedAt: &date, CompletedAt: &date,
	}
	s.RepairOrders[order.ID] = order
	s.RepairCostDetails[s.nextIDInternal()] = &model.RepairCostDetail{
		OrderID: order.ID, SparePartCost: cost, LaborCost: 200 * hours,
	}
	// 在日志中记录超长的修复时间
	s.RepairLogs[s.nextIDInternal()] = &model.RepairLog{
		OrderID: order.ID, UserID: workerID, Action: "维修完成",
		Content: fmt.Sprintf("修复工时: %.1f小时", hours),
	}
}

// ... (此处省略原有的辅助方法，如 AddUser, FindUser 等，保持不变) ...

func (s *Store) AddRepairOrder(order *model.RepairOrder) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.RepairOrders[order.ID] = order
}

func (s *Store) AddBase(id uint, base *model.Base) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Bases[id] = base
}

func (s *Store) AddFactory(id uint, factory *model.Factory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Factories[id] = factory
}

func (s *Store) AddWorkshop(id uint, workshop *model.Workshop) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Workshops[id] = workshop
}

func (s *Store) AddUser(id uint, user *model.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Users[id] = user
}

func (s *Store) FindUser(id uint) *model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Users[id]
}

func (s *Store) NextIDInternal() uint {
	return s.nextIDInternal()
}

// ============ Organization 辅助方法 ============

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

// ============ User 辅助方法 ============

func (s *Store) UpdateUser(id uint, fn func(*model.User)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if user, ok := s.Users[id]; ok {
		fn(user)
		return true
	}
	return false
}

func (s *Store) GetUsers() []*model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var users []*model.User
	for _, u := range s.Users {
		users = append(users, u)
	}
	return users
}

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

func (s *Store) Close() error {
	return nil
}
