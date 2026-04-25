package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/memory"
)

type MemoryAgentRepository struct {
	store *memory.Store
}

func NewMemoryAgentRepository() *MemoryAgentRepository {
	return &MemoryAgentRepository{
		store: memory.GetStore(),
	}
}

var _ IAgentRepository = (*MemoryAgentRepository)(nil)

// =====================================================
// Manual & Retrieval Repositories
// =====================================================

func (r *MemoryAgentRepository) CreateManualDocument(doc *model.ManualDocument) error {
	doc.ID = r.store.NextID()
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()
	r.store.ManualDocuments[doc.ID] = doc
	return nil
}

func (r *MemoryAgentRepository) GetManualDocumentByID(id uint) (*model.ManualDocument, error) {
	doc, ok := r.store.ManualDocuments[id]
	if !ok {
		return nil, fmt.Errorf("document not found")
	}
	return doc, nil
}

func (r *MemoryAgentRepository) CreateManualChunks(chunks []model.ManualChunk) error {
	for i := range chunks {
		chunks[i].ID = r.store.NextID()
		chunks[i].CreatedAt = time.Now()
		r.store.ManualChunks[chunks[i].ID] = &chunks[i]
	}
	return nil
}

func (r *MemoryAgentRepository) SearchManualChunks(query string, equipmentTypeID *uint) ([]model.ManualChunk, error) {
	var results []model.ManualChunk
	count := 0
	for _, chunk := range r.store.ManualChunks {
		if count >= 10 {
			break
		}
		
		match := strings.Contains(strings.ToLower(chunk.Content), strings.ToLower(query))
		if !match {
			continue
		}
		
		if equipmentTypeID != nil {
			doc, ok := r.store.ManualDocuments[chunk.ID]
			if ok && doc.EquipmentTypeID != nil && *doc.EquipmentTypeID != *equipmentTypeID {
				continue
			}
		}
		
		results = append(results, *chunk)
		count++
	}
	return results, nil
}

// =====================================================
// Runtime & Cost Analysis Repositories
// =====================================================

func (r *MemoryAgentRepository) CreateOrUpdateRepairCost(cost *model.RepairCostDetail) error {
	if cost.ID == 0 {
		cost.ID = r.store.NextID()
		cost.CreatedAt = time.Now()
	}
	cost.UpdatedAt = time.Now()
	r.store.RepairCostDetails[cost.ID] = cost
	return nil
}

func (r *MemoryAgentRepository) GetRepairCostByOrderID(orderID uint) (*model.RepairCostDetail, error) {
	for _, cost := range r.store.RepairCostDetails {
		if cost.OrderID == orderID {
			return cost, nil
		}
	}
	return nil, fmt.Errorf("cost detail not found")
}

func (r *MemoryAgentRepository) CreateOrUpdateRuntimeSnapshot(snapshot *model.EquipmentRuntimeSnapshot) error {
	if snapshot.ID == 0 {
		snapshot.ID = r.store.NextID()
		snapshot.CreatedAt = time.Now()
	}
	r.store.RuntimeSnapshots[snapshot.ID] = snapshot
	return nil
}

func (r *MemoryAgentRepository) GetRuntimeSnapshots(equipmentID uint, startDate, endDate string) ([]model.EquipmentRuntimeSnapshot, error) {
	var results []model.EquipmentRuntimeSnapshot
	for _, s := range r.store.RuntimeSnapshots {
		if s.EquipmentID == equipmentID && s.SnapshotDate >= startDate && s.SnapshotDate <= endDate {
			results = append(results, *s)
		}
	}
	return results, nil
}

// =====================================================
// Agent Session & Artifact Repositories
// =====================================================

func (r *MemoryAgentRepository) CreateSession(session *model.AgentSession) error {
	session.ID = r.store.NextID()
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	r.store.AgentSessions[session.ID] = session
	return nil
}

func (r *MemoryAgentRepository) GetSessionByID(id uint) (*model.AgentSession, error) {
	session, ok := r.store.AgentSessions[id]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}

func (r *MemoryAgentRepository) GetSessionByTraceID(traceID string) (*model.AgentSession, error) {
	for _, s := range r.store.AgentSessions {
		if s.TraceID == traceID {
			return s, nil
		}
	}
	return nil, fmt.Errorf("session not found")
}

func (r *MemoryAgentRepository) CreateArtifact(artifact *model.AgentArtifact) error {
	artifact.ID = r.store.NextID()
	artifact.CreatedAt = time.Now()
	r.store.AgentArtifacts[artifact.ID] = artifact
	return nil
}

func (r *MemoryAgentRepository) GetArtifactByID(id uint) (*model.AgentArtifact, error) {
	artifact, ok := r.store.AgentArtifacts[id]
	if !ok {
		return nil, fmt.Errorf("artifact not found")
	}
	return artifact, nil
}

func (r *MemoryAgentRepository) CreateEvidenceLinks(links []model.AgentEvidenceLink) error {
	for i := range links {
		links[i].ID = r.store.NextID()
		links[i].CreatedAt = time.Now()
		r.store.AgentEvidenceLinks[links[i].ID] = &links[i]
	}
	return nil
}

func (r *MemoryAgentRepository) ListSessionsByUserID(userID uint, limit int) ([]model.AgentSession, error) {
	var results []model.AgentSession
	count := 0
	for _, s := range r.store.AgentSessions {
		if count >= limit {
			break
		}
		if s.UserID == userID {
			results = append(results, *s)
			count++
		}
	}
	return results, nil
}
