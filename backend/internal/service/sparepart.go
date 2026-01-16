package service

import (
	"errors"
	"fmt"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
)

// SparePartService
type SparePartService struct {
	partRepo       *repository.SparePartRepository
	inventoryRepo  *repository.SparePartInventoryRepository
	consumptionRepo *repository.SparePartConsumptionRepository
	userRepo       *UserRepository
}

func NewSparePartService() *SparePartService {
	return &SparePartService{
		partRepo:       repository.NewSparePartRepository(),
		inventoryRepo:  repository.NewSparePartInventoryRepository(),
		consumptionRepo: repository.NewSparePartConsumptionRepository(),
		userRepo:       NewUserRepository(),
	}
}

func (s *SparePartService) CreatePart(code, name, specification, unit string, factoryID *uint, safetyStock int) (*model.SparePart, error) {
	// Check if code exists
	if _, err := s.partRepo.GetByCode(code); err == nil {
		return nil, errors.New("code already exists")
	}

	part := &model.SparePart{
		Code:          code,
		Name:          name,
		Specification: specification,
		Unit:          unit,
		FactoryID:     factoryID,
		SafetyStock:   safetyStock,
	}

	if err := s.partRepo.Create(part); err != nil {
		return nil, err
	}

	return s.partRepo.GetByID(part.ID)
}

func (s *SparePartService) GetPartByID(id uint) (*model.SparePart, error) {
	return s.partRepo.GetByID(id)
}

func (s *SparePartService) ListParts(filter repository.SparePartFilter) (*SparePartListResult, error) {
	parts, total, err := s.partRepo.List(filter)
	if err != nil {
		return nil, err
	}

	return &SparePartListResult{
		Items: parts,
		Total: total,
	}, nil
}

func (s *SparePartService) UpdatePart(id uint, code, name, specification, unit string, factoryID *uint, safetyStock int) error {
	part, err := s.partRepo.GetByID(id)
	if err != nil {
		return ErrNotFound
	}

	// Check if new code conflicts
	if code != part.Code {
		if existing, _ := s.partRepo.GetByCode(code); existing != nil {
			return errors.New("code already exists")
		}
	}

	part.Code = code
	part.Name = name
	part.Specification = specification
	part.Unit = unit
	part.FactoryID = factoryID
	part.SafetyStock = safetyStock

	return s.partRepo.Update(part)
}

func (s *SparePartService) DeletePart(id uint) error {
	return s.partRepo.Delete(id)
}

// Inventory operations
func (s *SparePartService) StockIn(partID, factoryID uint, quantity int, remark string, userID uint) error {
	// Verify part exists
	part, err := s.partRepo.GetByID(partID)
	if err != nil {
		return ErrNotFound
	}

	// Update inventory
	if err := s.inventoryRepo.UpdateQuantity(partID, factoryID, quantity); err != nil {
		return err
	}

	// TODO: Create stock in transaction record
	_ = part
	_ = remark
	_ = userID

	return nil
}

func (s *SparePartService) StockOut(partID, factoryID uint, quantity int, orderID, taskID *uint, remark string, userID uint) error {
	// Verify part exists
	part, err := s.partRepo.GetByID(partID)
	if err != nil {
		return ErrNotFound
	}

	// Check stock
	inv, err := s.inventoryRepo.GetByPartAndFactory(partID, factoryID)
	if err != nil {
		return errors.New("inventory not found")
	}

	if inv.Quantity < quantity {
		return fmt.Errorf("insufficient stock: available %d, requested %d", inv.Quantity, quantity)
	}

	// Update inventory
	if err := s.inventoryRepo.UpdateQuantity(partID, factoryID, -quantity); err != nil {
		return err
	}

	// Create consumption record
	consumption := &model.SparePartConsumption{
		SparePartID: partID,
		OrderID:     orderID,
		TaskID:      taskID,
		Quantity:    quantity,
		UserID:      userID,
	}

	if err := s.consumptionRepo.Create(consumption); err != nil {
		return err
	}

	_ = part
	_ = remark

	return nil
}

func (s *SparePartService) GetInventory(filter repository.InventoryFilter) (*InventoryListResult, error) {
	inventories, total, err := s.inventoryRepo.List(filter)
	if err != nil {
		return nil, err
	}

	return &InventoryListResult{
		Items: inventories,
		Total: total,
	}, nil
}

func (s *SparePartService) GetLowStockAlerts() ([]repository.LowStockAlert, error) {
	return s.inventoryRepo.GetLowStockAlerts()
}

// Consumption operations
func (s *SparePartService) RecordConsumption(partID uint, quantity int, orderID, taskID *uint, userID uint, remark string) error {
	// Verify part exists
	_, err := s.partRepo.GetByID(partID)
	if err != nil {
		return ErrNotFound
	}

	consumption := &model.SparePartConsumption{
		SparePartID: partID,
		OrderID:     orderID,
		TaskID:      taskID,
		Quantity:    quantity,
		UserID:      userID,
	}

	if err := s.consumptionRepo.Create(consumption); err != nil {
		return err
	}

	_ = remark
	return nil
}

func (s *SparePartService) GetConsumption(filter repository.ConsumptionFilter) (*ConsumptionListResult, error) {
	consumptions, total, err := s.consumptionRepo.List(filter)
	if err != nil {
		return nil, err
	}

	return &ConsumptionListResult{
		Items: consumptions,
		Total: total,
	}, nil
}

func (s *SparePartService) GetStatistics() (*SparePartStatistics, error) {
	stats, err := s.consumptionRepo.GetStatistics()
	if err != nil {
		return nil, err
	}

	return &SparePartStatistics{
		TotalParts:        stats["total_parts"],
		LowStockCount:     stats["low_stock_count"],
		TotalStockValue:   0, // Calculate if needed
		MonthlyConsumption: stats["monthly_consumption"],
	}, nil
}

// Types
type SparePartListResult struct {
	Items []model.SparePart
	Total int64
}

type InventoryListResult struct {
	Items []model.SparePartInventory
	Total int64
}

type ConsumptionListResult struct {
	Items []model.SparePartConsumption
	Total int64
}

type SparePartStatistics struct {
	TotalParts        int64
	LowStockCount     int64
	TotalStockValue   float64
	MonthlyConsumption int64
}
