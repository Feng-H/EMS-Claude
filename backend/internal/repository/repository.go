package repository

import (
	"github.com/ems/backend/internal/model"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(database *gorm.DB) {
	DB = database
}

// Base Repository for model.Base entities
type BaseEntityRepository struct {
	db *gorm.DB
}

func NewBaseEntityRepository() *BaseEntityRepository {
	return &BaseEntityRepository{db: DB}
}

func (r *BaseEntityRepository) Create(base *model.Base) error {
	return r.db.Create(base).Error
}

func (r *BaseEntityRepository) GetByID(id uint) (*model.Base, error) {
	var base model.Base
	err := r.db.First(&base, id).Error
	if err != nil {
		return nil, err
	}
	return &base, nil
}

func (r *BaseEntityRepository) List() ([]model.Base, error) {
	var bases []model.Base
	err := r.db.Order("code").Find(&bases).Error
	return bases, err
}

func (r *BaseEntityRepository) Update(base *model.Base) error {
	return r.db.Save(base).Error
}

func (r *BaseEntityRepository) Delete(id uint) error {
	return r.db.Delete(&model.Base{}, id).Error
}

func (r *BaseEntityRepository) GetByCode(code string) (*model.Base, error) {
	var base model.Base
	err := r.db.Where("code = ?", code).First(&base).Error
	if err != nil {
		return nil, err
	}
	return &base, nil
}

// Factory Repository
type FactoryRepository struct {
	db *gorm.DB
}

func NewFactoryRepo() *FactoryRepository {
	return &FactoryRepository{db: DB}
}

func (r *FactoryRepository) Create(factory *model.Factory) error {
	return r.db.Create(factory).Error
}

func (r *FactoryRepository) GetByID(id uint) (*model.Factory, error) {
	var factory model.Factory
	err := r.db.Preload("Base").First(&factory, id).Error
	if err != nil {
		return nil, err
	}
	return &factory, nil
}

func (r *FactoryRepository) List() ([]model.Factory, error) {
	var factories []model.Factory
	err := r.db.Preload("Base").Order("base_id, code").Find(&factories).Error
	return factories, err
}

func (r *FactoryRepository) ListByBaseID(baseID uint) ([]model.Factory, error) {
	var factories []model.Factory
	err := r.db.Where("base_id = ?", baseID).Order("code").Find(&factories).Error
	return factories, err
}

func (r *FactoryRepository) Update(factory *model.Factory) error {
	return r.db.Save(factory).Error
}

func (r *FactoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Factory{}, id).Error
}

func (r *FactoryRepository) GetByCode(baseID uint, code string) (*model.Factory, error) {
	var factory model.Factory
	err := r.db.Where("base_id = ? AND code = ?", baseID, code).First(&factory).Error
	if err != nil {
		return nil, err
	}
	return &factory, nil
}

// Workshop Repository
type WorkshopRepository struct {
	db *gorm.DB
}

func NewWorkshopRepo() *WorkshopRepository {
	return &WorkshopRepository{db: DB}
}

func (r *WorkshopRepository) Create(workshop *model.Workshop) error {
	return r.db.Create(workshop).Error
}

func (r *WorkshopRepository) GetByID(id uint) (*model.Workshop, error) {
	var workshop model.Workshop
	err := r.db.Preload("Factory").Preload("Factory.Base").First(&workshop, id).Error
	if err != nil {
		return nil, err
	}
	return &workshop, nil
}

func (r *WorkshopRepository) List() ([]model.Workshop, error) {
	var workshops []model.Workshop
	err := r.db.Preload("Factory").Preload("Factory.Base").Order("factory_id, code").Find(&workshops).Error
	return workshops, err
}

func (r *WorkshopRepository) ListByFactoryID(factoryID uint) ([]model.Workshop, error) {
	var workshops []model.Workshop
	err := r.db.Where("factory_id = ?", factoryID).Order("code").Find(&workshops).Error
	return workshops, err
}

func (r *WorkshopRepository) Update(workshop *model.Workshop) error {
	return r.db.Save(workshop).Error
}

func (r *WorkshopRepository) Delete(id uint) error {
	return r.db.Delete(&model.Workshop{}, id).Error
}

func (r *WorkshopRepository) GetByCode(factoryID uint, code string) (*model.Workshop, error) {
	var workshop model.Workshop
	err := r.db.Where("factory_id = ? AND code = ?", factoryID, code).First(&workshop).Error
	if err != nil {
		return nil, err
	}
	return &workshop, nil
}

// Equipment Type Repository
type EquipmentTypeRepository struct {
	db *gorm.DB
}

func NewEquipmentTypeRepo() *EquipmentTypeRepository {
	return &EquipmentTypeRepository{db: DB}
}

func (r *EquipmentTypeRepository) Create(equipmentType *model.EquipmentType) error {
	return r.db.Create(equipmentType).Error
}

func (r *EquipmentTypeRepository) GetByID(id uint) (*model.EquipmentType, error) {
	var equipmentType model.EquipmentType
	err := r.db.First(&equipmentType, id).Error
	if err != nil {
		return nil, err
	}
	return &equipmentType, nil
}

func (r *EquipmentTypeRepository) List() ([]model.EquipmentType, error) {
	var types []model.EquipmentType
	err := r.db.Order("category, name").Find(&types).Error
	return types, err
}

func (r *EquipmentTypeRepository) Update(equipmentType *model.EquipmentType) error {
	return r.db.Save(equipmentType).Error
}

func (r *EquipmentTypeRepository) Delete(id uint) error {
	return r.db.Delete(&model.EquipmentType{}, id).Error
}

// Equipment Repository
type EquipmentRepository struct {
	db *gorm.DB
}

func NewEquipmentRepo() *EquipmentRepository {
	return &EquipmentRepository{db: DB}
}

func (r *EquipmentRepository) Create(equipment *model.Equipment) error {
	return r.db.Create(equipment).Error
}

func (r *EquipmentRepository) GetByID(id uint) (*model.Equipment, error) {
	var equipment model.Equipment
	err := r.db.Preload("Type").Preload("Workshop").Preload("Workshop.Factory").Preload("Workshop.Factory.Base").
		Preload("DedicatedMaintenance").First(&equipment, id).Error
	if err != nil {
		return nil, err
	}
	return &equipment, nil
}

func (r *EquipmentRepository) GetByQRCode(qrCode string) (*model.Equipment, error) {
	var equipment model.Equipment
	err := r.db.Preload("Type").Where("qr_code = ? OR code = ?", qrCode, qrCode).First(&equipment).Error
	if err != nil {
		return nil, err
	}
	return &equipment, nil
}

func (r *EquipmentRepository) GetByCode(code string) (*model.Equipment, error) {
	var equipment model.Equipment
	err := r.db.Preload("Type").Where("code = ?", code).First(&equipment).Error
	if err != nil {
		return nil, err
	}
	return &equipment, nil
}

type EquipmentFilter struct {
	Code       string
	Name       string
	TypeID     *uint
	FactoryID  *uint
	WorkshopID *uint
	Status     string
	Page       int
	PageSize   int
}

func (r *EquipmentRepository) List(filter EquipmentFilter) ([]model.Equipment, int64, error) {
	var equipments []model.Equipment
	var total int64

	query := r.db.Model(&model.Equipment{})

	if filter.Code != "" {
		query = query.Where("code LIKE ?", "%"+filter.Code+"%")
	}
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.TypeID != nil {
		query = query.Where("type_id = ?", *filter.TypeID)
	}
	if filter.WorkshopID != nil {
		query = query.Where("workshop_id = ?", *filter.WorkshopID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Filter by factory_id through workshop
	if filter.FactoryID != nil {
		var workshopIDs []uint
		r.db.Model(&model.Workshop{}).Where("factory_id = ?", *filter.FactoryID).Pluck("id", &workshopIDs)
		if len(workshopIDs) > 0 {
			query = query.Where("workshop_id IN ?", workshopIDs)
		} else {
			return []model.Equipment{}, 0, nil
		}
	}

	// Pagination
	offset := (filter.Page - 1) * filter.PageSize
	err := query.Preload("Type").Preload("Workshop").Preload("Workshop.Factory").
		Preload("DedicatedMaintenance").
		Offset(offset).Limit(filter.PageSize).Order("code").Find(&equipments).Error

	return equipments, total, err
}

func (r *EquipmentRepository) Update(equipment *model.Equipment) error {
	return r.db.Save(equipment).Error
}

func (r *EquipmentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Equipment{}, id).Error
}

func (r *EquipmentRepository) GetStatistics() (map[string]int64, error) {
	stats := make(map[string]int64)

	var total, running, stopped, maintenance, scrapped int64

	r.db.Model(&model.Equipment{}).Count(&total)
	r.db.Model(&model.Equipment{}).Where("status = ?", "running").Count(&running)
	r.db.Model(&model.Equipment{}).Where("status = ?", "stopped").Count(&stopped)
	r.db.Model(&model.Equipment{}).Where("status = ?", "maintenance").Count(&maintenance)
	r.db.Model(&model.Equipment{}).Where("status = ?", "scrapped").Count(&scrapped)

	stats["total"] = total
	stats["running"] = running
	stats["stopped"] = stopped
	stats["maintenance"] = maintenance
	stats["scrapped"] = scrapped

	return stats, nil
}

func (r *EquipmentRepository) GetByTypeStatistics() ([]map[string]interface{}, error) {
	type Result struct {
		TypeName  string
		Count     int64
		Running   int64
		Maintenance int64
	}

	var results []Result
	err := r.db.Model(&model.Equipment{}).
		Select("equipment_types.name as type_name, COUNT(*) as count, "+
			"SUM(CASE WHEN equipment.status = 'running' THEN 1 ELSE 0 END) as running, "+
			"SUM(CASE WHEN equipment.status = 'maintenance' THEN 1 ELSE 0 END) as maintenance").
		Joins("LEFT JOIN equipment_types ON equipment.type_id = equipment_types.id").
		Group("equipment_types.name").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	output := make([]map[string]interface{}, len(results))
	for i, r := range results {
		output[i] = map[string]interface{}{
			"type_name": r.TypeName,
			"count":     r.Count,
			"running":   r.Running,
			"maintenance": r.Maintenance,
		}
	}

	return output, nil
}
