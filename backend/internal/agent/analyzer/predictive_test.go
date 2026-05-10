package analyzer

import (
	"testing"
	"time"

	"github.com/ems/backend/internal/agent/tool"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/memory"
)

func setupPredictiveTest() (*PredictiveAnalyzer, *memory.Store) {
	config.Cfg = &config.Config{
		Storage: config.StorageConfig{Mode: "memory"},
	}
	store := memory.GetStore()
	rt := tool.NewRepairTool()
	mt := tool.NewMaintenanceTool()
	ret := tool.NewRetrievalTool(nil)
	return NewPredictiveAnalyzer(rt, mt, ret), store
}

func TestPredictRUL_NoRepairs(t *testing.T) {
	a, store := setupPredictiveTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{
		BaseModel: model.BaseModel{ID: eqID},
		WorkshopID: 1,
	}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	pred, err := a.PredictRUL(eqID, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if pred.HealthScore < 0 || pred.HealthScore > 100 {
		t.Errorf("Health score out of range: %f", pred.HealthScore)
	}
	if pred.EstimatedRULDays < 0 {
		t.Errorf("Negative RUL days: %d", pred.EstimatedRULDays)
	}
	if pred.Recommendation == "" {
		t.Error("Expected non-empty recommendation")
	}
}

func TestPredictRUL_HighRepairCount(t *testing.T) {
	a, store := setupPredictiveTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{
		BaseModel: model.BaseModel{ID: eqID},
		WorkshopID: 1,
	}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	for i := 0; i < 5; i++ {
		id := store.NextID()
		store.RepairOrders[id] = &model.RepairOrder{
			BaseModel:   model.BaseModel{ID: id},
			EquipmentID: eqID,
			Status:      model.RepairClosed,
		}
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	pred, err := a.PredictRUL(eqID, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(pred.RiskFactors) == 0 {
		t.Error("Expected risk factors for high repair count (>3)")
	}
}

func TestDetectSymptoms_MicroStop(t *testing.T) {
	a, store := setupPredictiveTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{
		BaseModel: model.BaseModel{ID: eqID},
		WorkshopID: 1,
	}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	now := time.Now()
	for i := 0; i < 3; i++ {
		id := store.NextID()
		start := now.Add(time.Duration(-i*10) * time.Minute)
		end := start.Add(15 * time.Minute)
		store.RepairOrders[id] = &model.RepairOrder{
			BaseModel:      model.BaseModel{ID: id, CreatedAt: now},
			EquipmentID:    eqID,
			Status:         model.RepairClosed,
			StartedAt:      &start,
			CompletedAt:    &end,
			FaultDescription: "微停",
		}
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	findings, err := a.DetectSymptoms(eqID, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	found := false
	for _, f := range findings {
		if f.Type == "micro_stop" {
			found = true
		}
	}
	if !found {
		t.Error("Expected micro_stop symptom for 3+ short repairs")
	}
}

func TestDetectSymptoms_NoSymptoms(t *testing.T) {
	a, store := setupPredictiveTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{
		BaseModel: model.BaseModel{ID: eqID},
		WorkshopID: 1,
	}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	findings, err := a.DetectSymptoms(eqID, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	for _, f := range findings {
		if f.Type != "pm_ineffective" || !contains(f.Title, "Demo") {
			t.Errorf("Expected no real findings for clean equipment, got: %s", f.Title)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestCalculateTCO(t *testing.T) {
	a, store := setupPredictiveTest()

	eqID := store.NextID()
	purchaseDate := time.Now().AddDate(-2, 0, 0)
	store.Equipment[eqID] = &model.Equipment{
		BaseModel:       model.BaseModel{ID: eqID},
		WorkshopID:      1,
		PurchasePrice:   100000.0,
		ScrapValue:      10000.0,
		HourlyLoss:      500.0,
		ServiceLifeYears: 10,
		PurchaseDate:    &purchaseDate,
	}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	orderID := store.NextID()
	store.RepairOrders[orderID] = &model.RepairOrder{
		BaseModel:   model.BaseModel{ID: orderID},
		EquipmentID: eqID,
		Status:      model.RepairClosed,
		StartedAt:   ptrTime(time.Now().Add(-5 * time.Hour)),
		CompletedAt: ptrTime(time.Now()),
	}
	costID := store.NextID()
	store.RepairCostDetails[costID] = &model.RepairCostDetail{
		ID:            costID,
		OrderID:       orderID,
		SparePartCost: 5000.0,
		LaborCost:     2000.0,
		OtherCost:     500.0,
		DowntimeLoss:  2500.0,
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	tco, err := a.CalculateTCO(eqID, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if tco.AccumulatedRepair == 0 {
		t.Error("Expected non-zero accumulated repair cost")
	}
	if tco.DepreciatedValue <= 0 {
		t.Errorf("Expected positive depreciated value, got %f", tco.DepreciatedValue)
	}
	if tco.TCO <= 0 {
		t.Errorf("Expected positive TCO, got %f", tco.TCO)
	}
}

func TestEvaluateRetirement_LowRatio(t *testing.T) {
	a, store := setupPredictiveTest()

	eqID := store.NextID()
	purchaseDate := time.Now().AddDate(-1, 0, 0)
	store.Equipment[eqID] = &model.Equipment{
		BaseModel:       model.BaseModel{ID: eqID},
		WorkshopID:      1,
		PurchasePrice:   100000.0,
		ScrapValue:      10000.0,
		HourlyLoss:      100.0,
		ServiceLifeYears: 10,
		PurchaseDate:    &purchaseDate,
	}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	orderID := store.NextID()
	store.RepairOrders[orderID] = &model.RepairOrder{
		BaseModel:   model.BaseModel{ID: orderID},
		EquipmentID: eqID,
		Status:      model.RepairClosed,
	}
	costID := store.NextID()
	store.RepairCostDetails[costID] = &model.RepairCostDetail{
		ID:            costID,
		OrderID:       orderID,
		SparePartCost: 1000.0,
		LaborCost:     500.0,
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	result, err := a.EvaluateRetirement(eqID, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result["decision"] != "continue" {
		t.Errorf("Expected 'continue' for low maintenance ratio, got %v", result["decision"])
	}
}

func TestEvaluateRetirement_HighRatio(t *testing.T) {
	a, store := setupPredictiveTest()

	eqID := store.NextID()
	purchaseDate := time.Now().AddDate(-5, 0, 0)
	store.Equipment[eqID] = &model.Equipment{
		BaseModel:       model.BaseModel{ID: eqID},
		WorkshopID:      1,
		PurchasePrice:   10000.0,
		ScrapValue:      1000.0,
		HourlyLoss:      50.0,
		ServiceLifeYears: 10,
		PurchaseDate:    &purchaseDate,
	}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	orderID := store.NextID()
	store.RepairOrders[orderID] = &model.RepairOrder{
		BaseModel:   model.BaseModel{ID: orderID},
		EquipmentID: eqID,
		Status:      model.RepairClosed,
	}
	costID := store.NextID()
	store.RepairCostDetails[costID] = &model.RepairCostDetail{
		ID:            costID,
		OrderID:       orderID,
		SparePartCost: 5000.0,
		LaborCost:     2000.0,
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	result, err := a.EvaluateRetirement(eqID, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	decision := result["decision"].(string)
	if decision != "retire" && decision != "evaluate" {
		t.Errorf("Expected 'retire' or 'evaluate' for high maintenance ratio, got %s", decision)
	}
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
