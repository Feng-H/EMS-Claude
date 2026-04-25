package tool

import (
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/memory"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/internal/repository"
)

type RepairTool struct {
	orderRepo *repository.RepairRepository
}

func NewRepairTool() *RepairTool {
	var orderRepo *repository.RepairRepository
	if config.Cfg.Storage.Mode != "memory" {
		orderRepo = repository.NewRepairRepo()
	}
	
	return &RepairTool{
		orderRepo: orderRepo,
	}
}

func (t *RepairTool) GetRecentOrdersByEquipment(equipmentID uint, limit int) ([]model.RepairOrder, error) {
	if config.Cfg.Storage.Mode == "memory" {
		var results []model.RepairOrder
		store := memory.GetStore()
		count := 0
		// In memory, we might need a more efficient way or just iterate (this is mock)
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
	filter := repository.RepairFilter{
		EquipmentID: &equipmentID,
		Page:        1,
		PageSize:    limit,
	}
	result, err := t.orderRepo.List(filter)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
