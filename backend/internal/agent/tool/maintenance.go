package tool

import (
	"fmt"
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

func (t *MaintenanceTool) checkPermission(equipmentID uint, user model.User) error {
	if user.Role == "admin" || user.FactoryID == nil {
		return nil
	}

	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		e := store.FindEquipment(equipmentID)
		if e == nil { return fmt.Errorf("equipment not found") }
		w, ok := store.Workshops[e.WorkshopID]
		if !ok || w.FactoryID != *user.FactoryID {
			return fmt.Errorf("access denied: equipment belongs to another factory")
		}
		return nil
	}

	repo := repository.NewEquipmentRepo()
	e, err := repo.GetByID(equipmentID)
	if err != nil { return err }
	if e.Workshop.FactoryID != *user.FactoryID {
		return fmt.Errorf("access denied: equipment belongs to another factory")
	}
	return nil
}

func (t *MaintenanceTool) GetMaintenanceCompliance(equipmentID uint, user model.User) (map[string]interface{}, error) {
	if err := t.checkPermission(equipmentID, user); err != nil { return nil, err }
	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		completed := 0
		total := 0
		for _, task := range store.MaintenanceTasks {
			if task.EquipmentID == equipmentID {
				total++
				if task.Status == model.MaintenanceCompleted {
					completed++
				}
			}
		}
		rate := 0.0
		if total > 0 { rate = float64(completed) / float64(total) }
		return map[string]interface{}{
			"total_tasks":     total,
			"completed_tasks": completed,
			"compliance_rate": rate,
		}, nil
	}

	// DB Mode
	return t.taskRepo.GetComplianceByEquipmentID(equipmentID)
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

func (t *MaintenanceTool) GetRecentTasksByEquipment(equipmentID uint, limit int, user model.User) ([]model.MaintenanceTask, error) {
	if err := t.checkPermission(equipmentID, user); err != nil { return nil, err }
	if config.Cfg.Storage.Mode == "memory" {
		var results []model.MaintenanceTask
		store := memory.GetStore()
		count := 0
		for _, task := range store.MaintenanceTasks {
			if task.EquipmentID == equipmentID {
				results = append(results, *task)
				count++
				if count >= limit { break }
			}
		}
		return results, nil
	}
	
	// DB Mode
	return t.taskRepo.GetByEquipmentID(equipmentID, limit)
}

