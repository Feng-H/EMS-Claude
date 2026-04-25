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
