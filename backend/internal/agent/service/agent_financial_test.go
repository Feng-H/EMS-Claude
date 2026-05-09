package service

import (
	"testing"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/memory"
)

func TestAgentService_HandleGetEquipmentFinancials(t *testing.T) {
        // Setup config
        config.Cfg = &config.Config{
                Storage: config.StorageConfig{Mode: "memory"},
        }
        store := memory.GetStore()
	// Create mock data
	eq := &model.Equipment{
		BaseModel: model.BaseModel{ID: 1001},
		Name: "Test Equipment",
		PurchasePrice: 50000.0,
		ScrapValue: 5000.0,
		HourlyLoss: 100.0,
		ServiceLifeYears: 10,
	}
	store.Equipment[eq.ID] = eq
	
	user := model.User{
		BaseModel: model.BaseModel{ID: 1},
		Role: model.RoleAdmin,
	}
	
	svc := NewAgentService()
	args := map[string]interface{}{"equipment_id": 1001}
	
	result, err := svc.handleGetEquipmentFinancials(user, args)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	financials, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map[string]interface{} result")
	}
	
	if financials["purchase_price"] != 50000.0 {
		t.Errorf("Expected purchase_price 50000.0, got %v", financials["purchase_price"])
	}
	if financials["scrap_value"] != 5000.0 {
		t.Errorf("Expected scrap_value 5000.0, got %v", financials["scrap_value"])
	}
	if financials["hourly_loss"] != 100.0 {
		t.Errorf("Expected hourly_loss 100.0, got %v", financials["hourly_loss"])
	}
}

func TestAgentService_HandleGetRepairCosts(t *testing.T) {
        // Setup config
        config.Cfg = &config.Config{
                Storage: config.StorageConfig{Mode: "memory"},
        }
        store := memory.GetStore()
	// Create mock data
	eqID := uint(2001)
	store.Equipment[eqID] = &model.Equipment{BaseModel: model.BaseModel{ID: eqID}, WorkshopID: 1}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	orderID := uint(3001)
	store.RepairOrders[orderID] = &model.RepairOrder{
		BaseModel: model.BaseModel{ID: orderID},
		EquipmentID: eqID,
		Status: model.RepairClosed,
	}
	
	costID := uint(4001)
	store.RepairCostDetails[costID] = &model.RepairCostDetail{
		ID: costID,
		OrderID: orderID,
		SparePartCost: 1000.0,
		LaborCost: 500.0,
		OtherCost: 100.0,
		DowntimeLoss: 200.0,
	}
	
	user := model.User{
		BaseModel: model.BaseModel{ID: 1},
		Role: model.RoleAdmin,
	}
	
	svc := NewAgentService()
	args := map[string]interface{}{"equipment_id": eqID}
	
	result, err := svc.handleGetRepairCosts(user, args)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	costs, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map[string]interface{} result")
	}
	
	if costs["total_cost"] != 1800.0 {
		t.Errorf("Expected total_cost 1800.0, got %v", costs["total_cost"])
	}
	if costs["downtime_loss"] != 200.0 {
		t.Errorf("Expected downtime_loss 200.0, got %v", costs["downtime_loss"])
	}
}
