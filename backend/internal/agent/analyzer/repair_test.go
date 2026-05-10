package analyzer

import (
	"testing"
	"time"

	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/tool"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/memory"
)

func setupRepairAuditTest() (*RepairAuditAnalyzer, *memory.Store) {
	config.Cfg = &config.Config{
		Storage: config.StorageConfig{Mode: "memory"},
	}
	store := memory.GetStore()
	rt := tool.NewRepairTool()
	ret := tool.NewRetrievalTool(nil)
	return NewRepairAuditAnalyzer(ret, rt), store
}

func TestAnalyze_NoAnomalies(t *testing.T) {
	a, store := setupRepairAuditTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{
		BaseModel:  model.BaseModel{ID: eqID},
		WorkshopID: 1,
	}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	req := &dto.RepairAuditRequest{}

	data, err := a.Analyze(req, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if data == nil {
		t.Fatal("Expected non-nil data")
	}
}

func TestAnalyze_RepeatFailure(t *testing.T) {
	a, store := setupRepairAuditTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{
		BaseModel:  model.BaseModel{ID: eqID},
		WorkshopID: 1,
	}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	now := time.Now()
	o1 := store.NextID()
	store.RepairOrders[o1] = &model.RepairOrder{
		BaseModel:        model.BaseModel{ID: o1, CreatedAt: now},
		EquipmentID:      eqID,
		Status:           model.RepairClosed,
		FaultDescription: "电机异响严重",
	}
	o2 := store.NextID()
	store.RepairOrders[o2] = &model.RepairOrder{
		BaseModel:        model.BaseModel{ID: o2, CreatedAt: now.Add(-24 * time.Hour)},
		EquipmentID:      eqID,
		Status:           model.RepairClosed,
		FaultDescription: "电机异响频繁",
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	req := &dto.RepairAuditRequest{}

	data, err := a.Analyze(req, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	found := false
	for _, anomaly := range data.Anomalies {
		if anomaly.AnomalyType == "repeat_failure" {
			found = true
		}
	}
	if !found {
		t.Error("Expected repeat_failure anomaly for similar faults within 72h")
	}
}

func TestAnalyze_HighCost(t *testing.T) {
	a, store := setupRepairAuditTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{
		BaseModel:  model.BaseModel{ID: eqID},
		WorkshopID: 1,
	}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	o1 := store.NextID()
	store.RepairOrders[o1] = &model.RepairOrder{
		BaseModel:        model.BaseModel{ID: o1, CreatedAt: time.Now()},
		EquipmentID:      eqID,
		Status:           model.RepairClosed,
		FaultDescription: "液压系统故障",
		CostDetail: &model.RepairCostDetail{
			SparePartCost: 4000.0,
			LaborCost:     2000.0,
		},
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	req := &dto.RepairAuditRequest{}

	data, err := a.Analyze(req, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	found := false
	for _, anomaly := range data.Anomalies {
		if anomaly.AnomalyType == "high_cost" {
			found = true
		}
	}
	if !found {
		t.Error("Expected high_cost anomaly for repair > 5000")
	}
}

func TestIsSimilar_MatchingKeywords(t *testing.T) {
	a := &RepairAuditAnalyzer{}
	if !a.isSimilar("电机异响", "电机异响频繁") {
		t.Error("Expected similar for matching keyword '电机'")
	}
	if !a.isSimilar("轴承温度过高", "轴承温度异常") {
		t.Error("Expected similar for matching keyword '轴承'")
	}
}

func TestIsSimilar_NoMatch(t *testing.T) {
	a := &RepairAuditAnalyzer{}
	if a.isSimilar("设备外观磨损", "液压油需要更换") {
		t.Error("Expected not similar for unrelated descriptions")
	}
}
