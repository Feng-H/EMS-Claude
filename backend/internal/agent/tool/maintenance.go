package tool

import (
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/memory"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/internal/repository"
)

type MaintenanceTool struct {
	planRepo *repository.MaintenancePlanRepository
	taskRepo *repository.MaintenanceTaskRepository
}

func NewMaintenanceTool() *MaintenanceTool {
	var planRepo *repository.MaintenancePlanRepository
	var taskRepo *repository.MaintenanceTaskRepository
	
	if config.Cfg.Storage.Mode != "memory" {
		planRepo = repository.NewMaintenancePlanRepository()
		taskRepo = repository.NewMaintenanceTaskRepository()
	}
	
	return &MaintenanceTool{
		planRepo: planRepo,
		taskRepo: taskRepo,
	}
}

func (t *MaintenanceTool) GetPlanByEquipmentType(typeID uint) (*model.MaintenancePlan, error) {
	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		for _, plan := range store.MaintenancePlans {
			if plan.EquipmentTypeID == typeID {
				return plan, nil
			}
		}
		return nil, nil
	}
	
	// Database mode (simplified for MVP)
	plans, err := t.planRepo.List()
	if err != nil {
		return nil, err
	}
	for _, p := range plans {
		if p.EquipmentTypeID == typeID {
			return &p, nil
		}
	}
	return nil, nil
}
