package tool

import (
	"testing"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/memory"
)

func setupMaintenanceToolTest() (*MaintenanceTool, *memory.Store) {
	config.Cfg = &config.Config{
		Storage: config.StorageConfig{Mode: "memory"},
	}
	return NewMaintenanceTool(), memory.GetStore()
}

func TestGetMaintenanceCompliance_AllCompleted(t *testing.T) {
	mt, store := setupMaintenanceToolTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{BaseModel: model.BaseModel{ID: eqID}, WorkshopID: 1}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	for i := 0; i < 3; i++ {
		id := store.NextID()
		store.MaintenanceTasks[id] = &model.MaintenanceTask{
			BaseModel:   model.BaseModel{ID: id},
			EquipmentID: eqID,
			Status:      model.MaintenanceCompleted,
		}
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	result, err := mt.GetMaintenanceCompliance(eqID, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result["total_tasks"] != 3 {
		t.Errorf("Expected total_tasks 3, got %v", result["total_tasks"])
	}
	if result["completed_tasks"] != 3 {
		t.Errorf("Expected completed_tasks 3, got %v", result["completed_tasks"])
	}
	if result["compliance_rate"] != 1.0 {
		t.Errorf("Expected compliance_rate 1.0, got %v", result["compliance_rate"])
	}
}

func TestGetMaintenanceCompliance_Partial(t *testing.T) {
	mt, store := setupMaintenanceToolTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{BaseModel: model.BaseModel{ID: eqID}, WorkshopID: 1}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	id1 := store.NextID()
	store.MaintenanceTasks[id1] = &model.MaintenanceTask{
		BaseModel:   model.BaseModel{ID: id1},
		EquipmentID: eqID,
		Status:      model.MaintenanceCompleted,
	}
	id2 := store.NextID()
	store.MaintenanceTasks[id2] = &model.MaintenanceTask{
		BaseModel:   model.BaseModel{ID: id2},
		EquipmentID: eqID,
		Status:      model.MaintenancePending,
	}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	result, _ := mt.GetMaintenanceCompliance(eqID, user)

	if result["total_tasks"] != 2 {
		t.Errorf("Expected total_tasks 2, got %v", result["total_tasks"])
	}
	if result["completed_tasks"] != 1 {
		t.Errorf("Expected completed_tasks 1, got %v", result["completed_tasks"])
	}
}

func TestGetMaintenanceCompliance_NoTasks(t *testing.T) {
	mt, store := setupMaintenanceToolTest()

	eqID := store.NextID()
	store.Equipment[eqID] = &model.Equipment{BaseModel: model.BaseModel{ID: eqID}, WorkshopID: 1}
	store.Workshops[1] = &model.Workshop{BaseModel: model.BaseModel{ID: 1}, FactoryID: 1}

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	result, _ := mt.GetMaintenanceCompliance(eqID, user)

	if result["total_tasks"] != 0 {
		t.Errorf("Expected 0 tasks, got %v", result["total_tasks"])
	}
	if result["compliance_rate"] != 0.0 {
		t.Errorf("Expected 0 compliance rate for no tasks, got %v", result["compliance_rate"])
	}
}

func TestGetTasksByFilter_FactoryIsolation(t *testing.T) {
	mt, store := setupMaintenanceToolTest()

	factory1 := uint(1)
	factory2 := uint(2)

	ws1 := store.NextID()
	store.Workshops[ws1] = &model.Workshop{BaseModel: model.BaseModel{ID: ws1}, FactoryID: factory1}
	ws2 := store.NextID()
	store.Workshops[ws2] = &model.Workshop{BaseModel: model.BaseModel{ID: ws2}, FactoryID: factory2}

	eq1 := store.NextID()
	store.Equipment[eq1] = &model.Equipment{BaseModel: model.BaseModel{ID: eq1}, WorkshopID: ws1}
	eq2 := store.NextID()
	store.Equipment[eq2] = &model.Equipment{BaseModel: model.BaseModel{ID: eq2}, WorkshopID: ws2}

	t1 := store.NextID()
	store.MaintenanceTasks[t1] = &model.MaintenanceTask{
		BaseModel:   model.BaseModel{ID: t1},
		EquipmentID: eq1,
	}
	t2 := store.NextID()
	store.MaintenanceTasks[t2] = &model.MaintenanceTask{
		BaseModel:   model.BaseModel{ID: t2},
		EquipmentID: eq2,
	}

	user := model.User{
		BaseModel: model.BaseModel{ID: 1},
		Role:      "engineer",
		FactoryID: &factory1,
	}

	tasks, _, err := mt.GetTasksByFilter(
		repository.MaintenanceTaskFilter{Page: 1, PageSize: 100},
		user,
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	for _, task := range tasks {
		eq := store.FindEquipment(task.EquipmentID)
		ws := store.Workshops[eq.WorkshopID]
		if ws.FactoryID != factory1 {
			t.Errorf("Expected only factory1 tasks, got factory %d", ws.FactoryID)
		}
	}
}

func TestGetPlanByEquipmentType_Found(t *testing.T) {
	mt, store := setupMaintenanceToolTest()

	typeID := store.NextID()
	planID := store.NextID()
	store.MaintenancePlans[planID] = &model.MaintenancePlan{
		BaseModel:        model.BaseModel{ID: planID},
		EquipmentTypeID: typeID,
		Name:             "二级保养计划",
		CycleDays:        30,
	}

	plan, err := mt.GetPlanByEquipmentType(typeID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if plan == nil {
		t.Fatal("Expected plan to be found")
	}
	if plan.Name != "二级保养计划" {
		t.Errorf("Expected plan name '二级保养计划', got %s", plan.Name)
	}
}

func TestGetPlanByEquipmentType_NotFound(t *testing.T) {
	mt, _ := setupMaintenanceToolTest()

	plan, err := mt.GetPlanByEquipmentType(9999)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if plan != nil {
		t.Error("Expected nil plan for non-existent type")
	}
}
