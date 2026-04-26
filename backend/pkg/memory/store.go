package memory

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

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
	
	SpareParts           map[uint]*model.SparePart
	SparePartInventory   map[uint]*model.SparePartInventory
	SparePartConsumption map[uint]*model.SparePartConsumption
	
	KnowledgeArticles    map[uint]*model.KnowledgeArticle
	
	AgentSkills           map[uint]*model.AgentSkill
	AgentKnowledges       map[string]*model.AgentKnowledge
	AgentExperiences      map[uint]*model.AgentExperience
	AgentConversations    map[uint]*model.AgentConversation
	AgentMessages         map[uint]*model.AgentMessage
	AgentPushSubscriptions map[uint]*model.AgentPushSubscription
	AgentUsages           map[uint]*model.AgentUsage
	AgentArtifacts        map[uint]*model.AgentArtifact
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
			KnowledgeArticles:    make(map[uint]*model.KnowledgeArticle),
			AgentSkills:           make(map[uint]*model.AgentSkill),
			AgentKnowledges:       make(map[string]*model.AgentKnowledge),
			AgentExperiences:      make(map[uint]*model.AgentExperience),
			AgentConversations:    make(map[uint]*model.AgentConversation),
			AgentMessages:         make(map[uint]*model.AgentMessage),
			AgentPushSubscriptions: make(map[uint]*model.AgentPushSubscription),
			AgentUsages:           make(map[uint]*model.AgentUsage),
			AgentArtifacts:        make(map[uint]*model.AgentArtifact),
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
	hp, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	base := &model.Base{BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: now}, Code: "BASE-HQ", Name: "集团总部基地"}
	s.Bases[base.ID] = base
	fac := &model.Factory{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, BaseID: base.ID, Code: "FAC-SZ", Name: "苏州智能工厂"}
	s.Factories[fac.ID] = fac
	ws := &model.Workshop{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, FactoryID: fac.ID, Code: "WS-MCH", Name: "精密机加车间"}
	s.Workshops[ws.ID] = ws

	admin := &model.User{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Username: "admin", PasswordHash: string(hp), Name: "系统管理员", Role: model.RoleAdmin, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
	s.Users[admin.ID] = admin
	liSi := &model.User{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Username: "maint_li", PasswordHash: string(hp), Name: "预防型-李四", Role: model.RoleMaintenance, FactoryID: &fac.ID, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
	s.Users[liSi.ID] = liSi
	zs := &model.User{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Username: "maint_zhang", PasswordHash: string(hp), Name: "救火型-张三", Role: model.RoleMaintenance, FactoryID: &fac.ID, IsActive: true, ApprovalStatus: model.ApprovalStatusApproved}
	s.Users[zs.ID] = zs

	cnc := &model.EquipmentType{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Name: "高精度数控机床", Category: "加工设备"}
	s.EquipmentTypes[cnc.ID] = cnc

	eA := &model.Equipment{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Code: "CNC-001", Name: "李四维护-A区机床", TypeID: cnc.ID, WorkshopID: ws.ID, PurchasePrice: 280000.0, PurchaseDate: timePtr(now.AddDate(-3, 0, 0)), ServiceLifeYears: 8, ScrapValue: 28000.0, HourlyLoss: 150.0, Status: "running", DedicatedMaintenanceID: &liSi.ID}
	s.Equipment[eA.ID] = eA
	eB := &model.Equipment{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Code: "CNC-002", Name: "张三维护-B区机床", TypeID: cnc.ID, WorkshopID: ws.ID, PurchasePrice: 280000.0, PurchaseDate: timePtr(now.AddDate(-3, 0, 0)), ServiceLifeYears: 8, ScrapValue: 28000.0, HourlyLoss: 150.0, Status: "stopped", DedicatedMaintenanceID: &zs.ID}
	s.Equipment[eB.ID] = eB
	eP := &model.Equipment{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Code: "PRESS-05", Name: "12年老旧冲床", TypeID: cnc.ID, WorkshopID: ws.ID, PurchasePrice: 150000.0, PurchaseDate: timePtr(now.AddDate(-12, 0, 0)), ServiceLifeYears: 10, ScrapValue: 5000.0, HourlyLoss: 80.0, Status: "maintenance"}
	s.Equipment[eP.ID] = eP

	pump := &model.SparePart{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Code: "PUMP-01", Name: "高压柱塞泵", FactoryID: &fac.ID, SafetyStock: 2}
	s.SpareParts[pump.ID] = pump
	s.SparePartInventory[s.nextIDInternal()] = &model.SparePartInventory{SparePartID: pump.ID, FactoryID: fac.ID, Quantity: 3}

	plan := &model.MaintenancePlan{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Name: "CNC二级保养", EquipmentTypeID: cnc.ID, CycleDays: 30, WorkHours: 4.0}
	s.MaintenancePlans[plan.ID] = plan

	for i := 1; i <= 6; i++ {
		dAgo := i * 30
		tA := &model.MaintenanceTask{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, PlanID: plan.ID, EquipmentID: eA.ID, AssignedTo: liSi.ID, Status: "completed", CompletedAt: timePtr(now.AddDate(0, 0, -dAgo)), ActualHours: 4.2}
		s.MaintenanceTasks[tA.ID] = tA
		tB := &model.MaintenanceTask{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, PlanID: plan.ID, EquipmentID: eB.ID, AssignedTo: zs.ID, Status: "completed", CompletedAt: timePtr(now.AddDate(0, 0, -dAgo)), ActualHours: 0.2}
		s.MaintenanceTasks[tB.ID] = tB
	}

	order := &model.RepairOrder{BaseModel: model.BaseModel{ID: s.nextIDInternal(), CreatedAt: now.AddDate(0, 0, -10)}, EquipmentID: eB.ID, Status: model.RepairClosed, Priority: 1, AssignedTo: &zs.ID, FaultDescription: "泵轴承毁坏", Solution: "更换高压泵", StartedAt: timePtr(now.AddDate(0, 0, -10)), CompletedAt: timePtr(now.AddDate(0, 0, -10).Add(4 * time.Hour)), ClosedAt: timePtr(now.AddDate(0, 0, -8))}
	s.RepairOrders[order.ID] = order
	s.RepairCostDetails[s.nextIDInternal()] = &model.RepairCostDetail{OrderID: order.ID, SparePartCost: 52000.0, LaborCost: 3000.0}

	skill := &model.AgentSkill{BaseModel: model.BaseModel{ID: s.nextIDInternal()}, Name: "级联失效审计", Status: "active", Description: "分析连锁损失风险", Steps: `[{"step": 1, "tool": "get_cost_analysis"}]`}
	s.AgentSkills[skill.ID] = skill
	k := &model.AgentKnowledge{ID: "k_seed_001", Title: "高压泵级联失效预防", Status: "confirmed", Summary: "劣质滤芯诱因", Confidence: 0.98}
	s.AgentKnowledges[k.ID] = k
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
