package tool

import (
	"testing"

	"github.com/ems/backend/internal/agent/repository"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/memory"
)

func setupRetrievalTest() (*RetrievalTool, *memory.Store) {
	config.Cfg = &config.Config{
		Storage: config.StorageConfig{Mode: "memory"},
	}
	agentRepo := repository.NewMemoryAgentRepository()
	return NewRetrievalTool(agentRepo), memory.GetStore()
}

func TestGetEquipmentProfile_Found(t *testing.T) {
	rt, store := setupRetrievalTest()

	typeID := store.NextID()
	store.EquipmentTypes[typeID] = &model.EquipmentType{
		BaseModel: model.BaseModel{ID: typeID},
		Name:      "注塑机",
	}
	wsID := store.NextID()
	store.Workshops[wsID] = &model.Workshop{
		BaseModel: model.BaseModel{ID: wsID},
		Name:      "A车间",
		FactoryID: 1,
	}
	store.Factories[1] = &model.Factory{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "总厂",
	}

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{
		BaseModel:     model.BaseModel{ID: eqID},
		Code:          "CNC-001",
		Name:          "一号注塑机",
		TypeID:        typeID,
		WorkshopID:    wsID,
		PurchasePrice: 50000.0,
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	result, err := rt.GetEquipmentProfile(eqID, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result["code"] != "CNC-001" {
		t.Errorf("Expected code CNC-001, got %v", result["code"])
	}
	if result["name"] != "一号注塑机" {
		t.Errorf("Expected name 一号注塑机, got %v", result["name"])
	}
	if result["type_name"] != "注塑机" {
		t.Errorf("Expected type_name 注塑机, got %v", result["type_name"])
	}
	if result["workshop_name"] != "A车间" {
		t.Errorf("Expected workshop_name A车间, got %v", result["workshop_name"])
	}
	if result["factory_name"] != "总厂" {
		t.Errorf("Expected factory_name 总厂, got %v", result["factory_name"])
	}
	if result["purchase_price"] != 50000.0 {
		t.Errorf("Expected purchase_price 50000, got %v", result["purchase_price"])
	}
}

func TestGetEquipmentProfile_NotFound(t *testing.T) {
	rt, _ := setupRetrievalTest()

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	_, err := rt.GetEquipmentProfile(9999, user)
	if err == nil {
		t.Error("Expected error for non-existent equipment")
	}
}

func TestGetEquipmentProfile_PermissionDenied(t *testing.T) {
	rt, store := setupRetrievalTest()

	wsID := store.NextID()
	factory1 := uint(1)
	store.Workshops[wsID] = &model.Workshop{BaseModel: model.BaseModel{ID: wsID}, FactoryID: factory1}

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{
		BaseModel:  model.BaseModel{ID: eqID},
		WorkshopID: wsID,
	}

	factory2 := uint(2)
	user := model.User{
		BaseModel:  model.BaseModel{ID: 2},
		Role:       "engineer",
		FactoryID:  &factory2,
	}

	_, err := rt.GetEquipmentProfile(eqID, user)
	if err == nil {
		t.Error("Expected permission denied for different factory")
	}
}

func TestSearchManualKnowledge_EmptyQuery(t *testing.T) {
	rt, _ := setupRetrievalTest()

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	results, err := rt.SearchManualKnowledge("", nil, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty query, got %d", len(results))
	}
}

func TestSearchManualKnowledge_KnowledgeWeight(t *testing.T) {
	rt, store := setupRetrievalTest()

	artID := store.NextID()
	store.KnowledgeArticles[artID] = &model.KnowledgeArticle{
		BaseModel:        model.BaseModel{ID: artID},
		Title:            "轴承异响诊断",
		FaultPhenomenon:  "设备运行时出现间歇性异响",
		Solution:         "检查轴承润滑状态",
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	results, err := rt.SearchManualKnowledge("轴承异响", nil, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("Expected at least 1 result")
	}
	if results[0].EvidenceType != "knowledge" {
		t.Errorf("Expected evidence type 'knowledge', got %s", results[0].EvidenceType)
	}
	if results[0].Score < 0.70 {
		t.Errorf("Expected base score >= 0.70 for knowledge article, got %f", results[0].Score)
	}
}

func TestSearchManualKnowledge_LimitToFive(t *testing.T) {
	rt, store := setupRetrievalTest()

	for i := 0; i < 8; i++ {
		artID := store.NextID()
		store.KnowledgeArticles[artID] = &model.KnowledgeArticle{
			BaseModel:       model.BaseModel{ID: artID},
			Title:           "电机保养",
			FaultPhenomenon: "电机温度过高",
			Solution:        "检查散热系统",
		}
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	results, _ := rt.SearchManualKnowledge("电机", nil, user)
	if len(results) > 5 {
		t.Errorf("Expected at most 5 results, got %d", len(results))
	}
}
