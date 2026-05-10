package repository

import (
	"strings"
	"testing"
	"time"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/memory"
)

func setupRepoTest() IAgentRepository {
	config.Cfg = &config.Config{
		Storage: config.StorageConfig{Mode: "memory"},
	}
	return NewMemoryAgentRepository()
}

func TestSearchManualChunks_BasicMatch(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	docID := store.NextID()
	typeID := uint(100)
	store.ManualDocuments[docID] = &model.ManualDocument{
		BaseModel:       model.BaseModel{ID: docID},
		EquipmentTypeID: &typeID,
	}

	chunkID := store.NextID()
	store.ManualChunks[chunkID] = &model.ManualChunk{
		BaseModel:     model.BaseModel{ID: chunkID},
		DocumentID:    docID,
		Content:       "检查液压系统是否存在漏油现象",
		SectionTitle:  "液压系统维护",
	}

	results, err := repo.SearchManualChunks("漏油", &typeID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}
	if results[0].ID != chunkID {
		t.Errorf("Expected chunk ID %d, got %d", chunkID, results[0].ID)
	}
}

func TestSearchManualChunks_CaseInsensitive(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	chunkID := store.NextID()
	store.ManualChunks[chunkID] = &model.ManualChunk{
		BaseModel: model.BaseModel{ID: chunkID},
		Content:   "Motor Bearing Replacement",
	}

	results, _ := repo.SearchManualChunks("motor", nil)
	if len(results) != 1 {
		t.Errorf("Expected case-insensitive match, got %d results", len(results))
	}
}

func TestSearchManualChunks_Limit(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	for i := 0; i < 15; i++ {
		id := store.NextID()
		store.ManualChunks[id] = &model.ManualChunk{
			BaseModel: model.BaseModel{ID: id},
			Content:   "轴承润滑 轴承润滑 轴承润滑",
		}
	}

	results, _ := repo.SearchManualChunks("轴承润滑", nil)
	if len(results) > 10 {
		t.Errorf("Expected at most 10 results, got %d", len(results))
	}
}

func TestGetRepairCostByOrderID(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	costID := store.NextID()
	store.RepairCostDetails[costID] = &model.RepairCostDetail{
		ID:            costID,
		OrderID:       500,
		SparePartCost: 1000.0,
		LaborCost:     500.0,
	}

	cost, err := repo.GetRepairCostByOrderID(500)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cost.SparePartCost != 1000.0 {
		t.Errorf("Expected SparePartCost 1000, got %f", cost.SparePartCost)
	}
}

func TestGetRepairCostByOrderID_NotFound(t *testing.T) {
	repo := setupRepoTest()
	_, err := repo.GetRepairCostByOrderID(999)
	if err == nil {
		t.Error("Expected error for non-existent order")
	}
}

func TestGetRuntimeSnapshots_DateRange(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	snap1 := store.NextID()
	store.RuntimeSnapshots[snap1] = &model.EquipmentRuntimeSnapshot{
		BaseModel:     model.BaseModel{ID: snap1},
		EquipmentID:   1,
		SnapshotDate:  "2026-04-15",
	}
	snap2 := store.NextID()
	store.RuntimeSnapshots[snap2] = &model.EquipmentRuntimeSnapshot{
		BaseModel:     model.BaseModel{ID: snap2},
		EquipmentID:   1,
		SnapshotDate:  "2026-04-20",
	}
	snap3 := store.NextID()
	store.RuntimeSnapshots[snap3] = &model.EquipmentRuntimeSnapshot{
		BaseModel:     model.BaseModel{ID: snap3},
		EquipmentID:   2,
		SnapshotDate:  "2026-04-18",
	}

	results, _ := repo.GetRuntimeSnapshots(1, "2026-04-10", "2026-04-18")
	if len(results) != 1 {
		t.Errorf("Expected 1 result (snap1), got %d", len(results))
	}
	if len(results) > 0 && results[0].SnapshotDate != "2026-04-15" {
		t.Errorf("Expected snap1 date, got %s", results[0].SnapshotDate)
	}
}

func TestListSessionsByUserID_SortedDesc(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	s1 := store.NextID()
	store.AgentSessions[s1] = &model.AgentSession{
		BaseModel: model.BaseModel{ID: s1, CreatedAt: time.Now().Add(-2 * time.Hour)},
		UserID:    1,
		Scenario:  "first",
	}
	s2 := store.NextID()
	store.AgentSessions[s2] = &model.AgentSession{
		BaseModel: model.BaseModel{ID: s2, CreatedAt: time.Now()},
		UserID:    1,
		Scenario:  "second",
	}

	results, err := repo.ListSessionsByUserID(1, 10)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("Expected 2 sessions, got %d", len(results))
	}
	if results[0].Scenario != "second" {
		t.Error("Expected most recent session first (desc order)")
	}
}

func TestListSessionsByUserID_Limit(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	for i := 0; i < 5; i++ {
		id := store.NextID()
		store.AgentSessions[id] = &model.AgentSession{
			BaseModel: model.BaseModel{ID: id, CreatedAt: time.Now()},
			UserID:    1,
		}
	}

	results, _ := repo.ListSessionsByUserID(1, 3)
	if len(results) != 3 {
		t.Errorf("Expected 3 results (limit), got %d", len(results))
	}
}

func TestGetArtifactByID_WithEvidence(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	artID := store.NextID()
	store.AgentArtifacts[artID] = &model.AgentArtifact{
		BaseModel: model.BaseModel{ID: artID},
		RiskLevel: "high",
	}

	linkID := store.NextID()
	store.AgentEvidenceLinks[linkID] = &model.AgentEvidenceLink{
		BaseModel:    model.BaseModel{ID: linkID},
		ArtifactID:   artID,
		EvidenceType: "knowledge",
	}

	art, err := repo.GetArtifactByID(artID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(art.EvidenceLinks) != 1 {
		t.Errorf("Expected 1 evidence link, got %d", len(art.EvidenceLinks))
	}
}

func TestMatchSkills(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	s1 := store.NextID()
	store.AgentSkills[s1] = &model.AgentSkill{
		BaseModel:          model.BaseModel{ID: s1},
		Name:               "TCO分析技能",
		Status:             "active",
		ApplicableScenarios: "成本分析 财务评估 TCO",
	}
	s2 := store.NextID()
	store.AgentSkills[s2] = &model.AgentSkill{
		BaseModel: model.BaseModel{ID: s2},
		Name:      "停用技能",
		Status:    "inactive",
	}

	results, err := repo.MatchSkills("TCO", 10)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 active skill match, got %d", len(results))
	}
	if results[0].Name != "TCO分析技能" {
		t.Errorf("Expected TCO skill, got %s", results[0].Name)
	}
}

func TestApplyDecayToExperiences(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	e1 := store.NextID()
	store.AgentExperiences[e1] = &model.AgentExperience{
		BaseModel: model.BaseModel{ID: e1},
		Status:    "active",
		Weight:    1.0,
		DecayRate: 0.1,
	}
	e2 := store.NextID()
	store.AgentExperiences[e2] = &model.AgentExperience{
		BaseModel: model.BaseModel{ID: e2},
		Status:    "inactive",
		Weight:    1.0,
		DecayRate: 0.1,
	}

	repo.ApplyDecayToExperiences()

	if store.AgentExperiences[e1].Weight != 0.9 {
		t.Errorf("Expected active experience weight 0.9, got %f", store.AgentExperiences[e1].Weight)
	}
	if store.AgentExperiences[e2].Weight != 1.0 {
		t.Errorf("Expected inactive experience weight unchanged at 1.0, got %f", store.AgentExperiences[e2].Weight)
	}
}

func TestListKnowledges_StatusFilter(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	store.AgentKnowledges["k1"] = &model.AgentKnowledge{
		ID:     "k1",
		Title:  "轴承故障诊断",
		Status: "draft",
	}
	store.AgentKnowledges["k2"] = &model.AgentKnowledge{
		ID:     "k2",
		Title:  "电机过热处理",
		Status: "confirmed",
	}

	results, _ := repo.ListKnowledges("confirmed", "", 10)
	if len(results) != 1 {
		t.Fatalf("Expected 1 confirmed knowledge, got %d", len(results))
	}
	if results[0].ID != "k2" {
		t.Errorf("Expected k2, got %s", results[0].ID)
	}
}

func TestListKnowledges_QueryFilter(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	store.AgentKnowledges["k1"] = &model.AgentKnowledge{
		ID:     "k1",
		Title:  "轴承故障诊断",
		Status: "confirmed",
	}
	store.AgentKnowledges["k2"] = &model.AgentKnowledge{
		ID:     "k2",
		Title:  "电机过热处理",
		Status: "confirmed",
	}

	results, _ := repo.ListKnowledges("", "轴承", 10)
	if len(results) != 1 {
		t.Errorf("Expected 1 result matching 轴承, got %d", len(results))
	}
}

func TestGetConversationByID_WithMessages(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	convID := store.NextID()
	store.AgentConversations[convID] = &model.AgentConversation{
		BaseModel: model.BaseModel{ID: convID, CreatedAt: time.Now()},
		UserID:    1,
	}

	m1 := store.NextID()
	store.AgentMessages[m1] = &model.AgentMessage{
		ID:             m1,
		ConversationID: convID,
		Role:           "user",
		Content:        "first",
		CreatedAt:      time.Now().Add(-1 * time.Hour),
	}
	m2 := store.NextID()
	store.AgentMessages[m2] = &model.AgentMessage{
		ID:             m2,
		ConversationID: convID,
		Role:           "assistant",
		Content:        "second",
		CreatedAt:      time.Now(),
	}

	conv, err := repo.GetConversationByID(convID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(conv.Messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(conv.Messages))
	}
	if conv.Messages[0].Content != "first" {
		t.Errorf("Expected oldest message first, got %s", conv.Messages[0].Content)
	}
}

func TestCreatePushSubscription_Upsert(t *testing.T) {
	repo := setupRepoTest()

	sub1 := &model.AgentPushSubscription{
		UserID:     1,
		PushType:   "predictive",
		WebhookURL: "https://example.com/hook1",
	}
	repo.CreatePushSubscription(sub1)

	found, _ := repo.GetPushSubscription(1, "predictive")
	if found.WebhookURL != "https://example.com/hook1" {
		t.Errorf("Expected original webhook, got %s", found.WebhookURL)
	}

	sub2 := &model.AgentPushSubscription{
		BaseModel:  model.BaseModel{ID: sub1.ID},
		UserID:     1,
		PushType:   "predictive",
		WebhookURL: "https://example.com/hook2",
	}
	repo.CreatePushSubscription(sub2)

	found2, _ := repo.GetPushSubscription(1, "predictive")
	if found2.WebhookURL != "https://example.com/hook2" {
		t.Errorf("Expected updated webhook after upsert, got %s", found2.WebhookURL)
	}
}

func TestSearchManualChunks_NoResults(t *testing.T) {
	repo := setupRepoTest()
	results, err := repo.SearchManualChunks("nonexistent", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty store, got %d", len(results))
	}
}

func TestSearchManualChunks_TypeIDFilter(t *testing.T) {
	repo := setupRepoTest()
	store := memory.GetStore()

	typeA := uint(1)
	typeB := uint(2)

	docA := store.NextID()
	store.ManualDocuments[docA] = &model.ManualDocument{
		BaseModel:       model.BaseModel{ID: docA},
		EquipmentTypeID: &typeA,
	}
	docB := store.NextID()
	store.ManualDocuments[docB] = &model.ManualDocument{
		BaseModel:       model.BaseModel{ID: docB},
		EquipmentTypeID: &typeB,
	}

	chunkA := store.NextID()
	store.ManualChunks[chunkA] = &model.ManualChunk{
		BaseModel:  model.BaseModel{ID: chunkA},
		DocumentID: docA,
		Content:    strings.Repeat("保养流程", 10),
	}
	chunkB := store.NextID()
	store.ManualChunks[chunkB] = &model.ManualChunk{
		BaseModel:  model.BaseModel{ID: chunkB},
		DocumentID: docB,
		Content:    strings.Repeat("保养流程", 10),
	}

	results, _ := repo.SearchManualChunks("保养", &typeA)
	for _, r := range results {
		if r.DocumentID == docB {
			t.Error("Should not return chunks from different equipment type")
		}
	}
}
