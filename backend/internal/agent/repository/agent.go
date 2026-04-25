package repository

import (
	"github.com/ems/backend/internal/model"
	"gorm.io/gorm"
)

type IAgentRepository interface {
	CreateManualDocument(doc *model.ManualDocument) error
	GetManualDocumentByID(id uint) (*model.ManualDocument, error)
	CreateManualChunks(chunks []model.ManualChunk) error
	SearchManualChunks(query string, equipmentTypeID *uint) ([]model.ManualChunk, error)
	CreateOrUpdateRepairCost(cost *model.RepairCostDetail) error
	GetRepairCostByOrderID(orderID uint) (*model.RepairCostDetail, error)
	CreateOrUpdateRuntimeSnapshot(snapshot *model.EquipmentRuntimeSnapshot) error
	GetRuntimeSnapshots(equipmentID uint, startDate, endDate string) ([]model.EquipmentRuntimeSnapshot, error)
	CreateSession(session *model.AgentSession) error
	GetSessionByID(id uint) (*model.AgentSession, error)
	GetSessionByTraceID(traceID string) (*model.AgentSession, error)
	CreateArtifact(artifact *model.AgentArtifact) error
	GetArtifactByID(id uint) (*model.AgentArtifact, error)
	CreateEvidenceLinks(links []model.AgentEvidenceLink) error
	CreateUsage(usage *model.AgentUsage) error
	ListSessionsByUserID(userID uint, limit int) ([]model.AgentSession, error)

	// Phase 2: Conversations & Messages
	CreateConversation(conv *model.AgentConversation) error
	GetConversationByID(id uint) (*model.AgentConversation, error)
	ListConversationsByUserID(userID uint, limit int) ([]model.AgentConversation, error)
	CreateMessage(msg *model.AgentMessage) error
	GetMessagesByConversationID(convID uint) ([]model.AgentMessage, error)
	CreateKnowledge(knowledge *model.AgentKnowledge) error

	// Phase 2: Skills
	CreateSkill(skill *model.AgentSkill) error
	GetSkillByID(id uint) (*model.AgentSkill)
	UpdateSkill(skill *model.AgentSkill) error
	ListSkills(status string, limit int) ([]model.AgentSkill, error)
	MatchSkills(intent string, limit int) ([]model.AgentSkill, error)

	// Phase 2: Experience
	CreateExperience(exp *model.AgentExperience) error
	ListActiveExperiences(userID uint) ([]model.AgentExperience, error)
	UpdateExperience(exp *model.AgentExperience) error
	ApplyDecayToExperiences() error

	// Phase 2: Push Subscriptions
	CreatePushSubscription(sub *model.AgentPushSubscription) error
	GetPushSubscription(userID uint, pushType string) (*model.AgentPushSubscription, error)
	ListPushSubscriptions(userID uint) ([]model.AgentPushSubscription, error)
}

type DBAgentRepository struct {
	db *gorm.DB
}

func NewDBAgentRepository(db *gorm.DB) *DBAgentRepository {
	return &DBAgentRepository{db: db}
}

// =====================================================
// Manual & Retrieval Repositories
// =====================================================

func (r *DBAgentRepository) CreateManualDocument(doc *model.ManualDocument) error {
	return r.db.Create(doc).Error
}

func (r *DBAgentRepository) GetManualDocumentByID(id uint) (*model.ManualDocument, error) {
	var doc model.ManualDocument
	err := r.db.Preload("EquipmentType").Preload("Equipment").First(&doc, id).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *DBAgentRepository) CreateManualChunks(chunks []model.ManualChunk) error {
	return r.db.Create(&chunks).Error
}

// SearchManualChunks performs a full-text search on manual chunks
func (r *DBAgentRepository) SearchManualChunks(query string, equipmentTypeID *uint) ([]model.ManualChunk, error) {
	var chunks []model.ManualChunk
	db := r.db.Model(&model.ManualChunk{})
	
	if equipmentTypeID != nil {
		db = db.Joins("JOIN equipment_manual_documents ON equipment_manual_documents.id = equipment_manual_chunks.document_id").
			Where("equipment_manual_documents.equipment_type_id = ?", *equipmentTypeID)
	}
	
	// Phase 1: Simple LIKE search (will upgrade to tsvector later)
	err := db.Where("content LIKE ?", "%"+query+"%").Limit(10).Find(&chunks).Error
	return chunks, err
}

// =====================================================
// Runtime & Cost Analysis Repositories
// =====================================================

func (r *DBAgentRepository) CreateOrUpdateRepairCost(cost *model.RepairCostDetail) error {
	return r.db.Save(cost).Error
}

func (r *DBAgentRepository) GetRepairCostByOrderID(orderID uint) (*model.RepairCostDetail, error) {
	var cost model.RepairCostDetail
	err := r.db.Where("order_id = ?", orderID).First(&cost).Error
	if err != nil {
		return nil, err
	}
	return &cost, nil
}

func (r *DBAgentRepository) CreateOrUpdateRuntimeSnapshot(snapshot *model.EquipmentRuntimeSnapshot) error {
	return r.db.Save(snapshot).Error
}

func (r *DBAgentRepository) GetRuntimeSnapshots(equipmentID uint, startDate, endDate string) ([]model.EquipmentRuntimeSnapshot, error) {
	var snapshots []model.EquipmentRuntimeSnapshot
	err := r.db.Where("equipment_id = ? AND snapshot_date BETWEEN ? AND ?", equipmentID, startDate, endDate).
		Order("snapshot_date DESC").Find(&snapshots).Error
	return snapshots, err
}

// =====================================================
// Agent Session & Artifact Repositories
// =====================================================

func (r *DBAgentRepository) CreateSession(session *model.AgentSession) error {
	return r.db.Create(session).Error
}

func (r *DBAgentRepository) GetSessionByID(id uint) (*model.AgentSession, error) {
	var session model.AgentSession
	err := r.db.Preload("Artifacts").First(&session, id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *DBAgentRepository) GetSessionByTraceID(traceID string) (*model.AgentSession, error) {
	var session model.AgentSession
	err := r.db.Where("trace_id = ?", traceID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *DBAgentRepository) CreateArtifact(artifact *model.AgentArtifact) error {
	return r.db.Create(artifact).Error
}

func (r *DBAgentRepository) GetArtifactByID(id uint) (*model.AgentArtifact, error) {
	var artifact model.AgentArtifact
	err := r.db.Preload("EvidenceLinks").First(&artifact, id).Error
	if err != nil {
		return nil, err
	}
	return &artifact, nil
}

func (r *DBAgentRepository) CreateEvidenceLinks(links []model.AgentEvidenceLink) error {
	return r.db.Create(&links).Error
}

func (r *DBAgentRepository) CreateUsage(usage *model.AgentUsage) error {
	return r.db.Create(usage).Error
}

func (r *DBAgentRepository) ListSessionsByUserID(userID uint, limit int) ([]model.AgentSession, error) {
	var sessions []model.AgentSession
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&sessions).Error
	return sessions, err
}

// =====================================================
// Phase 2: Conversations & Messages
// =====================================================

func (r *DBAgentRepository) CreateConversation(conv *model.AgentConversation) error {
	return r.db.Create(conv).Error
}

func (r *DBAgentRepository) GetConversationByID(id uint) (*model.AgentConversation, error) {
	var conv model.AgentConversation
	err := r.db.Preload("Messages").First(&conv, id).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *DBAgentRepository) ListConversationsByUserID(userID uint, limit int) ([]model.AgentConversation, error) {
	var convs []model.AgentConversation
	err := r.db.Where("user_id = ?", userID).Order("updated_at DESC").Limit(limit).Find(&convs).Error
	return convs, err
}

func (r *DBAgentRepository) CreateMessage(msg *model.AgentMessage) error {
	return r.db.Create(msg).Error
}

func (r *DBAgentRepository) GetMessagesByConversationID(convID uint) ([]model.AgentMessage, error) {
	var msgs []model.AgentMessage
	err := r.db.Where("conversation_id = ?", convID).Order("created_at ASC").Find(&msgs).Error
	return msgs, err
}

func (r *DBAgentRepository) CreateKnowledge(k *model.AgentKnowledge) error {
	return r.db.Create(k).Error
}

// =====================================================
// Phase 2: Skills
// =====================================================

func (r *DBAgentRepository) CreateSkill(skill *model.AgentSkill) error {
	return r.db.Create(skill).Error
}

func (r *DBAgentRepository) GetSkillByID(id uint) *model.AgentSkill {
	var skill model.AgentSkill
	if err := r.db.First(&skill, id).Error; err != nil {
		return nil
	}
	return &skill
}

func (r *DBAgentRepository) UpdateSkill(skill *model.AgentSkill) error {
	return r.db.Save(skill).Error
}

func (r *DBAgentRepository) ListSkills(status string, limit int) ([]model.AgentSkill, error) {
	var skills []model.AgentSkill
	db := r.db.Model(&model.AgentSkill{})
	if status != "" {
		db = db.Where("status = ?", status)
	}
	err := db.Order("usage_count DESC, created_at DESC").Limit(limit).Find(&skills).Error
	return skills, err
}

func (r *DBAgentRepository) MatchSkills(intent string, limit int) ([]model.AgentSkill, error) {
	var skills []model.AgentSkill
	// 初步使用全文检索或简单的 LIKE 匹配场景描述
	err := r.db.Where("status = 'active' AND (name ILIKE ? OR applicable_scenarios::text ILIKE ?)", 
		"%"+intent+"%", "%"+intent+"%").
		Order("success_rate DESC").Limit(limit).Find(&skills).Error
	return skills, err
}

// =====================================================
// Phase 2: Experience Repositories
// =====================================================

func (r *DBAgentRepository) CreateExperience(exp *model.AgentExperience) error {
	return r.db.Create(exp).Error
}

func (r *DBAgentRepository) ListActiveExperiences(userID uint) ([]model.AgentExperience, error) {
	var exps []model.AgentExperience
	err := r.db.Where("user_id = ? AND status = 'active' AND weight > 0.1", userID).
		Order("weight DESC").Find(&exps).Error
	return exps, err
}

func (r *DBAgentRepository) UpdateExperience(exp *model.AgentExperience) error {
	return r.db.Save(exp).Error
}

func (r *DBAgentRepository) ApplyDecayToExperiences() error {
	// 执行衰减公式：weight = weight * (1 - decay_rate)
	// 在生产环境中通常通过 Cron 每天运行一次
	return r.db.Model(&model.AgentExperience{}).
		Where("status = 'active'").
		Update("weight", gorm.Expr("weight * (1 - decay_rate)")).Error
}

// =====================================================
// Phase 2: Push Subscription Repositories
// =====================================================

func (r *DBAgentRepository) CreatePushSubscription(sub *model.AgentPushSubscription) error {
	return r.db.Save(sub).Error // Save uses upsert logic
}

func (r *DBAgentRepository) GetPushSubscription(userID uint, pushType string) (*model.AgentPushSubscription, error) {
	var sub model.AgentPushSubscription
	err := r.db.Where("user_id = ? AND push_type = ?", userID, pushType).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *DBAgentRepository) ListPushSubscriptions(userID uint) ([]model.AgentPushSubscription, error) {
	var subs []model.AgentPushSubscription
	err := r.db.Where("user_id = ?", userID).Find(&subs).Error
	return subs, err
}
