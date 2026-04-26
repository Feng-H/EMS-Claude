package repository

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/memory"
)

type MemoryAgentRepository struct {
	store *memory.Store
}

func NewMemoryAgentRepository() IAgentRepository {
	return &MemoryAgentRepository{
		store: memory.GetStore(),
	}
}

// =====================================================
// Manual & Retrieval Repositories
// =====================================================

func (r *MemoryAgentRepository) CreateManualDocument(doc *model.ManualDocument) error {
	doc.ID = r.store.NextID()
	doc.CreatedAt = time.Now()
	r.store.ManualDocuments[doc.ID] = doc
	return nil
}

func (r *MemoryAgentRepository) GetManualDocumentByID(id uint) (*model.ManualDocument, error) {
	if doc, ok := r.store.ManualDocuments[id]; ok {
		return doc, nil
	}
	return nil, fmt.Errorf("document not found")
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
	limit := 10

	for _, chunk := range r.store.ManualChunks {
		if count >= limit {
			break
		}
		
		match := strings.Contains(strings.ToLower(chunk.Content), strings.ToLower(query))
		if !match {
			continue
		}
		
		if equipmentTypeID != nil {
			doc, ok := r.store.ManualDocuments[chunk.DocumentID]
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
// Session & Artifact Repositories
// =====================================================

func (r *MemoryAgentRepository) CreateSession(session *model.AgentSession) error {
	session.ID = r.store.NextID()
	session.CreatedAt = time.Now()
	r.store.AgentSessions[session.ID] = session
	return nil
}

func (r *MemoryAgentRepository) GetSessionByID(id uint) (*model.AgentSession, error) {
	if s, ok := r.store.AgentSessions[id]; ok {
		// Populate artifacts
		var artifacts []model.AgentArtifact
		for _, a := range r.store.AgentArtifacts {
			if a.SessionID == id {
				artifacts = append(artifacts, *a)
			}
		}
		s.Artifacts = artifacts
		return s, nil
	}
	return nil, fmt.Errorf("session not found")
}

func (r *MemoryAgentRepository) GetSessionByTraceID(traceID string) (*model.AgentSession, error) {
	for _, s := range r.store.AgentSessions {
		if s.TraceID == traceID {
			return s, nil
		}
	}
	return nil, fmt.Errorf("session not found")
}

func (r *MemoryAgentRepository) ListSessionsByUserID(userID uint, limit int) ([]model.AgentSession, error) {
	var results []model.AgentSession
	for _, s := range r.store.AgentSessions {
		if s.UserID == userID {
			results = append(results, *s)
		}
	}
	// Sort by CreatedAt desc
	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.After(results[j].CreatedAt)
	})
	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

func (r *MemoryAgentRepository) CreateArtifact(artifact *model.AgentArtifact) error {
	artifact.ID = r.store.NextID()
	artifact.CreatedAt = time.Now()
	r.store.AgentArtifacts[artifact.ID] = artifact
	return nil
}

func (r *MemoryAgentRepository) GetArtifactByID(id uint) (*model.AgentArtifact, error) {
	if a, ok := r.store.AgentArtifacts[id]; ok {
		// Populate evidence links
		var links []model.AgentEvidenceLink
		for _, l := range r.store.AgentEvidenceLinks {
			if l.ArtifactID == id {
				links = append(links, *l)
			}
		}
		a.EvidenceLinks = links
		return a, nil
	}
	return nil, fmt.Errorf("artifact not found")
}

func (r *MemoryAgentRepository) CreateEvidenceLinks(links []model.AgentEvidenceLink) error {
	for i := range links {
		links[i].ID = r.store.NextID()
		links[i].CreatedAt = time.Now()
		r.store.AgentEvidenceLinks[links[i].ID] = &links[i]
	}
	return nil
}

func (r *MemoryAgentRepository) CreateUsage(usage *model.AgentUsage) error {
	usage.ID = r.store.NextID()
	usage.CreatedAt = time.Now()
	r.store.AgentUsages[usage.ID] = usage
	return nil
}

// =====================================================
// Conversation & Message Repositories
// =====================================================

func (r *MemoryAgentRepository) CreateConversation(conv *model.AgentConversation) error {
	conv.ID = r.store.NextID()
	conv.CreatedAt = time.Now()
	conv.UpdatedAt = time.Now()
	r.store.AgentConversations[conv.ID] = conv
	return nil
}

func (r *MemoryAgentRepository) GetConversationByID(id uint) (*model.AgentConversation, error) {
	if c, ok := r.store.AgentConversations[id]; ok {
		// Load messages
		var msgs []model.AgentMessage
		for _, m := range r.store.AgentMessages {
			if m.ConversationID == id {
				msgs = append(msgs, *m)
			}
		}
		sort.Slice(msgs, func(i, j int) bool {
			return msgs[i].CreatedAt.Before(msgs[j].CreatedAt)
		})
		c.Messages = msgs
		return c, nil
	}
	return nil, fmt.Errorf("conversation not found")
}

func (r *MemoryAgentRepository) ListConversationsByUserID(userID uint, limit int) ([]model.AgentConversation, error) {
	var results []model.AgentConversation
	for _, c := range r.store.AgentConversations {
		if c.UserID == userID {
			results = append(results, *c)
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].UpdatedAt.After(results[j].UpdatedAt)
	})
	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

func (r *MemoryAgentRepository) CreateMessage(msg *model.AgentMessage) error {
	msg.ID = r.store.NextID()
	msg.CreatedAt = time.Now()
	r.store.AgentMessages[msg.ID] = msg
	
	// Update conversation updated_at
	if c, ok := r.store.AgentConversations[msg.ConversationID]; ok {
		c.UpdatedAt = time.Now()
	}
	return nil
}

func (r *MemoryAgentRepository) GetMessagesByConversationID(convID uint) ([]model.AgentMessage, error) {
	var results []model.AgentMessage
	for _, m := range r.store.AgentMessages {
		if m.ConversationID == convID {
			results = append(results, *m)
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].CreatedAt.Before(results[j].CreatedAt)
	})
	return results, nil
}

// =====================================================
// Knowledge & Skill Repositories
// =====================================================

func (r *MemoryAgentRepository) CreateKnowledge(k *model.AgentKnowledge) error {
	if k.ID == "" {
		k.ID = fmt.Sprintf("k_%d", r.store.NextID())
	}
	k.CreatedAt = time.Now()
	r.store.AgentKnowledges[k.ID] = k
	return nil
}

func (r *MemoryAgentRepository) UpdateKnowledgeStatus(id string, status string, verifierID *uint) error {
	if k, ok := r.store.AgentKnowledges[id]; ok {
		k.Status = status
		k.UpdatedAt = time.Now()
		if verifierID != nil {
			k.VerifiedBy = verifierID
			now := time.Now()
			k.VerifiedAt = &now
		}
		return nil
	}
	return fmt.Errorf("knowledge not found")
}

func (r *MemoryAgentRepository) ListKnowledges(status string, limit int) ([]model.AgentKnowledge, error) {
	var results []model.AgentKnowledge
	count := 0
	for _, k := range r.store.AgentKnowledges {
		if count >= limit {
			break
		}
		if status == "" || k.Status == status {
			results = append(results, *k)
			count++
		}
	}
	return results, nil
}

func (r *MemoryAgentRepository) CreateSkill(skill *model.AgentSkill) error {
	skill.ID = r.store.NextID()
	skill.CreatedAt = time.Now()
	r.store.AgentSkills[skill.ID] = skill
	return nil
}

func (r *MemoryAgentRepository) GetSkillByID(id uint) *model.AgentSkill {
	return r.store.AgentSkills[id]
}

func (r *MemoryAgentRepository) UpdateSkill(skill *model.AgentSkill) error {
	skill.UpdatedAt = time.Now()
	r.store.AgentSkills[skill.ID] = skill
	return nil
}

func (r *MemoryAgentRepository) ListSkills(status string, limit int) ([]model.AgentSkill, error) {
	var results []model.AgentSkill
	count := 0
	for _, s := range r.store.AgentSkills {
		if count >= limit {
			break
		}
		if status == "" || s.Status == status {
			results = append(results, *s)
			count++
		}
	}
	return results, nil
}

func (r *MemoryAgentRepository) MatchSkills(intent string, limit int) ([]model.AgentSkill, error) {
	var results []model.AgentSkill
	count := 0
	for _, s := range r.store.AgentSkills {
		if count >= limit {
			break
		}
		if s.Status == "active" && (strings.Contains(s.Name, intent) || strings.Contains(s.ApplicableScenarios, intent)) {
			results = append(results, *s)
			count++
		}
	}
	return results, nil
}

// =====================================================
// Experience Repositories
// =====================================================

func (r *MemoryAgentRepository) CreateExperience(exp *model.AgentExperience) error {
	exp.ID = r.store.NextID()
	exp.CreatedAt = time.Now()
	r.store.AgentExperiences[exp.ID] = exp
	return nil
}

func (r *MemoryAgentRepository) ListActiveExperiences(userID uint) ([]model.AgentExperience, error) {
	var results []model.AgentExperience
	for _, e := range r.store.AgentExperiences {
		if e.UserID == userID && e.Status == "active" && e.Weight > 0.1 {
			results = append(results, *e)
		}
	}
	return results, nil
}

func (r *MemoryAgentRepository) UpdateExperience(exp *model.AgentExperience) error {
	r.store.AgentExperiences[exp.ID] = exp
	return nil
}

func (r *MemoryAgentRepository) ApplyDecayToExperiences() error {
	for _, e := range r.store.AgentExperiences {
		if e.Status == "active" {
			e.Weight = e.Weight * (1.0 - e.DecayRate)
		}
	}
	return nil
}

// =====================================================
// Push Subscription Repositories
// =====================================================

func (r *MemoryAgentRepository) CreatePushSubscription(sub *model.AgentPushSubscription) error {
	if sub.ID == 0 {
		sub.ID = r.store.NextID()
		sub.CreatedAt = time.Now()
	}
	sub.UpdatedAt = time.Now()
	r.store.AgentPushSubscriptions[sub.ID] = sub
	return nil
}

func (r *MemoryAgentRepository) GetPushSubscription(userID uint, pushType string) (*model.AgentPushSubscription, error) {
	for _, s := range r.store.AgentPushSubscriptions {
		if s.UserID == userID && s.PushType == pushType {
			return s, nil
		}
	}
	return nil, fmt.Errorf("subscription not found")
}

func (r *MemoryAgentRepository) ListPushSubscriptions(userID uint) ([]model.AgentPushSubscription, error) {
	var results []model.AgentPushSubscription
	for _, s := range r.store.AgentPushSubscriptions {
		if s.UserID == userID {
			results = append(results, *s)
		}
	}
	return results, nil
}
