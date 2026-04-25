package tool

import (
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

func (t *RepairTool) GetFailureStats(equipmentID uint) (map[string]interface{}, error) {
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

func (t *RepairTool) GetCostAnalysis(equipmentID uint) (map[string]interface{}, error) {
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

func (t *RepairTool) GetRecentOrdersByEquipment(equipmentID uint, limit int) ([]model.RepairOrder, error) {
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
	// Note: RepairOrderFilter doesn't have EquipmentID in this repo version
	// We'll use GetByEquipmentID instead
	return t.orderRepo.GetByEquipmentID(equipmentID, limit)
}
