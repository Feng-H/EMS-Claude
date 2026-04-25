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
