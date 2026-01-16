package repository

import (
	"time"

	"github.com/ems/backend/internal/model"
	"gorm.io/gorm"
)

// SparePart Repository
type SparePartRepository struct {
	db *gorm.DB
}

func NewSparePartRepository() *SparePartRepository {
	return &SparePartRepository{db: DB}
}

func (r *SparePartRepository) Create(part *model.SparePart) error {
	return r.db.Create(part).Error
}

func (r *SparePartRepository) GetByID(id uint) (*model.SparePart, error) {
	var part model.SparePart
	err := r.db.Preload("Factory").First(&part, id).Error
	if err != nil {
		return nil, err
	}
	return &part, nil
}

func (r *SparePartRepository) GetByCode(code string) (*model.SparePart, error) {
	var part model.SparePart
	err := r.db.Preload("Factory").Where("code = ?", code).First(&part).Error
	if err != nil {
		return nil, err
	}
	return &part, nil
}

func (r *SparePartRepository) List(filter SparePartFilter) ([]model.SparePart, int64, error) {
	var parts []model.SparePart
	var total int64

	query := r.db.Model(&model.SparePart{})

	if filter.Code != "" {
		query = query.Where("code LIKE ?", "%"+filter.Code+"%")
	}
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (filter.Page - 1) * filter.PageSize
	err := query.Preload("Factory").
		Offset(offset).Limit(filter.PageSize).
		Find(&parts).Error

	return parts, total, err
}

func (r *SparePartRepository) Update(part *model.SparePart) error {
	return r.db.Save(part).Error
}

func (r *SparePartRepository) Delete(id uint) error {
	return r.db.Delete(&model.SparePart{}, id).Error
}

// Inventory Repository
type SparePartInventoryRepository struct {
	db *gorm.DB
}

func NewSparePartInventoryRepository() *SparePartInventoryRepository {
	return &SparePartInventoryRepository{db: DB}
}

func (r *SparePartInventoryRepository) GetByPartAndFactory(partID, factoryID uint) (*model.SparePartInventory, error) {
	var inv model.SparePartInventory
	err := r.db.Where("spare_part_id = ? AND factory_id = ?", partID, factoryID).
		Preload("SparePart").
		Preload("Factory").
		First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *SparePartInventoryRepository) List(filter InventoryFilter) ([]model.SparePartInventory, int64, error) {
	var inventories []model.SparePartInventory
	var total int64

	query := r.db.Model(&model.SparePartInventory{})

	if filter.SparePartID != nil {
		query = query.Where("spare_part_id = ?", *filter.SparePartID)
	}
	if filter.FactoryID != nil {
		query = query.Where("factory_id = ?", *filter.FactoryID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (filter.Page - 1) * filter.PageSize
	err := query.Preload("SparePart").
		Preload("Factory").
		Offset(offset).Limit(filter.PageSize).
		Find(&inventories).Error

	return inventories, total, err
}

func (r *SparePartInventoryRepository) Upsert(partID, factoryID uint, quantity int) error {
	var inv model.SparePartInventory
	err := r.db.Where("spare_part_id = ? AND factory_id = ?", partID, factoryID).
		First(&inv).Error

	if err == gorm.ErrRecordNotFound {
		// Create new
		inv = model.SparePartInventory{
			SparePartID: partID,
			FactoryID:   factoryID,
			Quantity:    quantity,
		}
		return r.db.Create(&inv).Error
	}

	// Update existing
	inv.Quantity = quantity
	return r.db.Save(&inv).Error
}

func (r *SparePartInventoryRepository) UpdateQuantity(partID, factoryID uint, delta int) error {
	var inv model.SparePartInventory
	err := r.db.Where("spare_part_id = ? AND factory_id = ?", partID, factoryID).
		First(&inv).Error

	if err == gorm.ErrRecordNotFound {
		// Create new with delta
		inv = model.SparePartInventory{
			SparePartID: partID,
			FactoryID:   factoryID,
			Quantity:    delta,
		}
		return r.db.Create(&inv).Error
	}

	if err != nil {
		return err
	}

	inv.Quantity += delta
	if inv.Quantity < 0 {
		inv.Quantity = 0
	}
	return r.db.Save(&inv).Error
}

func (r *SparePartInventoryRepository) GetLowStockAlerts() ([]LowStockAlert, error) {
	var alerts []LowStockAlert

	query := `
		SELECT
			sp.id as spare_part_id,
			sp.code as spare_part_code,
			sp.name as spare_part_name,
			spi.factory_id,
			f.name as factory_name,
			spi.quantity as current_stock,
			COALESCE(sp.safety_stock, 0) as safety_stock,
			COALESCE(sp.safety_stock, 0) - spi.quantity as shortage
		FROM spare_part_inventories spi
		INNER JOIN spare_parts sp ON spi.spare_part_id = sp.id
		INNER JOIN factories f ON spi.factory_id = f.id
		WHERE spi.quantity < COALESCE(sp.safety_stock, 0)
		ORDER BY shortage DESC
	`

	err := r.db.Raw(query).Scan(&alerts).Error
	return alerts, err
}

// Consumption Repository
type SparePartConsumptionRepository struct {
	db *gorm.DB
}

func NewSparePartConsumptionRepository() *SparePartConsumptionRepository {
	return &SparePartConsumptionRepository{db: DB}
}

func (r *SparePartConsumptionRepository) Create(consumption *model.SparePartConsumption) error {
	return r.db.Create(consumption).Error
}

func (r *SparePartConsumptionRepository) List(filter ConsumptionFilter) ([]model.SparePartConsumption, int64, error) {
	var consumptions []model.SparePartConsumption
	var total int64

	query := r.db.Model(&model.SparePartConsumption{})

	if filter.SparePartID != nil {
		query = query.Where("spare_part_id = ?", *filter.SparePartID)
	}
	if filter.OrderID != nil {
		query = query.Where("order_id = ?", *filter.OrderID)
	}
	if filter.TaskID != nil {
		query = query.Where("task_id = ?", *filter.TaskID)
	}
	if filter.DateFrom != "" {
		if t, err := time.Parse("2006-01-02", filter.DateFrom); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if filter.DateTo != "" {
		if t, err := time.Parse("2006-01-02", filter.DateTo); err == nil {
			query = query.Where("created_at < ?", t.AddDate(0, 0, 1))
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (filter.Page - 1) * filter.PageSize
	err := query.Preload("SparePart").
		Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(filter.PageSize).
		Find(&consumptions).Error

	return consumptions, total, err
}

// Get statistics
func (r *SparePartConsumptionRepository) GetStatistics() (map[string]int64, error) {
	stats := make(map[string]int64)

	// Total parts
	var totalParts int64
	r.db.Model(&model.SparePart{}).Count(&totalParts)
	stats["total_parts"] = totalParts

	// Low stock count
	var lowStockCount int64
	r.db.Raw(`
		SELECT COUNT(*)
		FROM spare_part_inventories spi
		INNER JOIN spare_parts sp ON spi.spare_part_id = sp.id
		WHERE spi.quantity < COALESCE(sp.safety_stock, 0)
	`).Scan(&lowStockCount)
	stats["low_stock_count"] = lowStockCount

	// Monthly consumption
	thisMonth := time.Now().Format("2006-01")
	var monthlyConsumption int64
	r.db.Model(&model.SparePartConsumption{}).
		Where("created_at >= ?", thisMonth+"-01").
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&monthlyConsumption)
	stats["monthly_consumption"] = monthlyConsumption

	return stats, nil
}

// Filter types
type SparePartFilter struct {
	Code     string
	Name     string
	Page     int
	PageSize int
}

type InventoryFilter struct {
	SparePartID *uint
	FactoryID   *uint
	LowStock    *bool
	Page        int
	PageSize    int
}

type ConsumptionFilter struct {
	SparePartID *uint
	OrderID     *uint
	TaskID      *uint
	DateFrom    string
	DateTo      string
	Page        int
	PageSize    int
}

type LowStockAlert struct {
	SparePartID   uint
	SparePartCode string
	SparePartName string
	FactoryID     uint
	FactoryName   string
	CurrentStock  int
	SafetyStock   int
	Shortage      int
}
