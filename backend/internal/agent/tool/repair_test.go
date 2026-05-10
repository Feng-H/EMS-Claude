package tool

import (
	"testing"
	"time"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/memory"
)

func setupRepairToolTest() (*RepairTool, *memory.Store) {
	config.Cfg = &config.Config{
		Storage: config.StorageConfig{Mode: "memory"},
	}
	return NewRepairTool(), memory.GetStore()
}

func TestGetCostByEquipmentID_Aggregation(t *testing.T) {
	rt, store := setupRepairToolTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{BaseModel: model.BaseModel{ID: eqID}, WorkshopID: 1}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	o1 := store.NextID()
	store.RepairOrders[o1] = &model.RepairOrder{
		BaseModel:   model.BaseModel{ID: o1},
		EquipmentID: eqID,
		Status:      model.RepairClosed,
	}
	o2 := store.NextID()
	store.RepairOrders[o2] = &model.RepairOrder{
		BaseModel:   model.BaseModel{ID: o2},
		EquipmentID: eqID,
		Status:      model.RepairClosed,
	}

	c1 := store.NextID()
	store.RepairCostDetails[c1] = &model.RepairCostDetail{
		ID: c1, OrderID: o1,
		SparePartCost: 1000.0, LaborCost: 500.0, OtherCost: 100.0, DowntimeLoss: 200.0,
	}
	c2 := store.NextID()
	store.RepairCostDetails[c2] = &model.RepairCostDetail{
		ID: c2, OrderID: o2,
		SparePartCost: 2000.0, LaborCost: 300.0, OtherCost: 50.0, DowntimeLoss: 100.0,
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	result, err := rt.GetCostByEquipmentID(eqID, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedTotal := 3000.0 + 800.0 + 150.0 + 300.0 // 4250
	if result["total_cost"] != expectedTotal {
		t.Errorf("Expected total_cost %f, got %v", expectedTotal, result["total_cost"])
	}
	if result["spare_part_cost"] != 3000.0 {
		t.Errorf("Expected spare_part_cost 3000, got %v", result["spare_part_cost"])
	}
}

func TestGetCostByEquipmentID_Empty(t *testing.T) {
	rt, store := setupRepairToolTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{BaseModel: model.BaseModel{ID: eqID}, WorkshopID: 1}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	result, err := rt.GetCostByEquipmentID(eqID, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result["total_cost"] != 0.0 {
		t.Errorf("Expected zero total_cost for no repairs, got %v", result["total_cost"])
	}
}

func TestGetFailureStats_MTTR(t *testing.T) {
	rt, store := setupRepairToolTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{BaseModel: model.BaseModel{ID: eqID}, WorkshopID: 1}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	now := time.Now()
	start := now.Add(-4 * time.Hour)
	end := now

	o1 := store.NextID()
	store.RepairOrders[o1] = &model.RepairOrder{
		BaseModel:   model.BaseModel{ID: o1},
		EquipmentID: eqID,
		Status:      model.RepairClosed,
		StartedAt:   &start,
		CompletedAt: &end,
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	result, err := rt.GetFailureStats(eqID, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result["repair_count"] != 1 {
		t.Errorf("Expected repair_count 1, got %v", result["repair_count"])
	}
	if result["mttr"] == 0 {
		t.Error("Expected non-zero MTTR")
	}
}

func TestGetFailureStats_IgnoresOpenOrders(t *testing.T) {
	rt, store := setupRepairToolTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{BaseModel: model.BaseModel{ID: eqID}, WorkshopID: 1}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	o1 := store.NextID()
	store.RepairOrders[o1] = &model.RepairOrder{
		BaseModel:   model.BaseModel{ID: o1},
		EquipmentID: eqID,
		Status:      model.RepairPending,
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	result, _ := rt.GetFailureStats(eqID, user)
	if result["repair_count"] != 0 {
		t.Error("Expected 0 repair count for pending orders (not audited/closed)")
	}
}

func TestGetFailureStats_PermissionDenied(t *testing.T) {
	rt, store := setupRepairToolTest()

	eqID := store.NextID()
	factory1 := uint(1)
	factory2 := uint(2)
	store.Equipment[eqID] = &model.Equipment{BaseModel: model.BaseModel{ID: eqID}, WorkshopID: 1}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: factory1}

	user := model.User{
		BaseModel:  model.BaseModel{ID: 2},
		Role:       "engineer",
		FactoryID:  &factory2,
	}

	_, err := rt.GetFailureStats(eqID, user)
	if err == nil {
		t.Error("Expected permission denied for different factory")
	}
}
