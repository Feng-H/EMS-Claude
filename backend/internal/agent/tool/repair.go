package tool

import (
	"fmt"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/memory"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/internal/repository"
)

type RepairTool struct {
	orderRepo *repository.RepairOrderRepository
}

func NewRepairTool() *RepairTool {
	var orderRepo *repository.RepairOrderRepository
	if config.Cfg.Storage.Mode != "memory" {
		orderRepo = repository.NewRepairOrderRepository()
	}
	
	return &RepairTool{
		orderRepo: orderRepo,
	}
}

func (t *RepairTool) checkPermission(equipmentID uint, user model.User) error {
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

func (t *RepairTool) GetFailureStats(equipmentID uint, user model.User) (map[string]interface{}, error) {
	if err := t.checkPermission(equipmentID, user); err != nil { return nil, err }
	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		var totalDowntime float64
		repairCount := 0
		
		for _, order := range store.RepairOrders {
			if order.EquipmentID == equipmentID && (order.Status == model.RepairAudited || order.Status == model.RepairClosed) {
				repairCount++
				if order.StartedAt != nil && order.CompletedAt != nil {
					totalDowntime += order.CompletedAt.Sub(*order.StartedAt).Hours()
				}
			}
		}

		mttr := 0.0
		if repairCount > 0 { mttr = totalDowntime / float64(repairCount) }

		return map[string]interface{}{
			"repair_count":    repairCount,
			"total_downtime":  totalDowntime,
			"mttr":           mttr,
		}, nil
	}

	// DB Mode
	stats, err := t.orderRepo.GetStatisticsByEquipmentID(equipmentID)
	if err != nil { return nil, err }
	return stats, nil
}

func (t *RepairTool) GetCostByEquipmentID(equipmentID uint, user model.User) (map[string]interface{}, error) {
	if err := t.checkPermission(equipmentID, user); err != nil { return nil, err }
	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		var sparePartCost, laborCost float64
		for _, cost := range store.RepairCostDetails {
			if order, ok := store.RepairOrders[cost.OrderID]; ok && order.EquipmentID == equipmentID {
				sparePartCost += cost.SparePartCost
				laborCost += cost.LaborCost
			}
		}
		return map[string]interface{}{
			"total_cost":       sparePartCost + laborCost,
			"spare_part_cost": sparePartCost,
			"labor_cost":      laborCost,
		}, nil
	}

	// DB Mode
	return t.orderRepo.GetCostByEquipmentID(equipmentID)
}

func (t *RepairTool) GetRecentOrdersByEquipment(equipmentID uint, limit int, user model.User) ([]model.RepairOrder, error) {
	if err := t.checkPermission(equipmentID, user); err != nil { return nil, err }
	if config.Cfg.Storage.Mode == "memory" {
		var results []model.RepairOrder
		store := memory.GetStore()
		count := 0
		for _, order := range store.RepairOrders {
			if order.EquipmentID == equipmentID {
				results = append(results, *order)
				count++
				if count >= limit {
					break
				}
			}
		}
		return results, nil
	}
	
	// Database mode
	return t.orderRepo.GetByEquipmentID(equipmentID, limit)
}

func (t *RepairTool) GetOrdersByFilter(filter repository.RepairOrderFilter, user model.User) ([]model.RepairOrder, int64, error) {
	// If user is not admin, restrict to their factory
	if user.Role != "admin" && user.FactoryID != nil {
		filter.FactoryID = *user.FactoryID
	}

	if config.Cfg.Storage.Mode == "memory" {
		// Basic memory implementation
		var results []model.RepairOrder
		store := memory.GetStore()
		for _, order := range store.RepairOrders {
			// Factory isolation in memory
			if user.Role != "admin" && user.FactoryID != nil {
				e := store.FindEquipment(order.EquipmentID)
				if e == nil { continue }
				w := store.Workshops[e.WorkshopID]
				if w.FactoryID != *user.FactoryID { continue }
			}

			// Apply basic status filter if present
			if filter.Status != "" && order.Status != model.RepairStatus(filter.Status) {
				continue
			}
			results = append(results, *order)
		}
		return results, int64(len(results)), nil
	}

	return t.orderRepo.List(filter)
}

