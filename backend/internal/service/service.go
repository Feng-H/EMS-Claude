package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrNotFound      = errors.New("record not found")
	ErrDuplicateCode = errors.New("code already exists")
	ErrInvalidInput  = errors.New("invalid input")
)

// Base Service
type BaseService struct {
	repo *repository.BaseEntityRepository
}

func NewBaseService() *BaseService {
	return &BaseService{
		repo: repository.NewBaseEntityRepository(),
	}
}

func (s *BaseService) Create(code, name string) (*model.Base, error) {
	// Check if code exists
	if _, err := s.repo.GetByCode(code); err == nil {
		return nil, ErrDuplicateCode
	}

	base := &model.Base{
		Code: code,
		Name: name,
	}

	if err := s.repo.Create(base); err != nil {
		return nil, err
	}

	return base, nil
}

func (s *BaseService) GetByID(id uint) (*model.Base, error) {
	return s.repo.GetByID(id)
}

func (s *BaseService) List() ([]model.Base, error) {
	return s.repo.List()
}

func (s *BaseService) Update(id uint, code, name string) (*model.Base, error) {
	base, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Check if new code conflicts
	if base.Code != code {
		if _, err := s.repo.GetByCode(code); err == nil {
			return nil, ErrDuplicateCode
		}
	}

	base.Code = code
	base.Name = name

	if err := s.repo.Update(base); err != nil {
		return nil, err
	}

	return base, nil
}

func (s *BaseService) Delete(id uint) error {
	// Check if has factories
	factoryRepo := repository.NewFactoryRepo()
	factories, _ := factoryRepo.ListByBaseID(id)
	if len(factories) > 0 {
		return errors.New("cannot delete base with existing factories")
	}

	return s.repo.Delete(id)
}

// Factory Service
type FactoryService struct {
	repo     *repository.FactoryRepository
	baseRepo *repository.BaseEntityRepository
}

func NewFactoryService() *FactoryService {
	return &FactoryService{
		repo:     repository.NewFactoryRepo(),
		baseRepo: repository.NewBaseEntityRepository(),
	}
}

func (s *FactoryService) Create(baseID uint, code, name string) (*model.Factory, error) {
	// Verify base exists
	if _, err := s.baseRepo.GetByID(baseID); err != nil {
		return nil, ErrNotFound
	}

	// Check if code exists in this base
	if _, err := s.repo.GetByCode(baseID, code); err == nil {
		return nil, ErrDuplicateCode
	}

	factory := &model.Factory{
		BaseID: baseID,
		Code:   code,
		Name:   name,
	}

	if err := s.repo.Create(factory); err != nil {
		return nil, err
	}

	return factory, nil
}

func (s *FactoryService) GetByID(id uint) (*model.Factory, error) {
	return s.repo.GetByID(id)
}

func (s *FactoryService) List() ([]model.Factory, error) {
	return s.repo.List()
}

func (s *FactoryService) ListByBaseID(baseID uint) ([]model.Factory, error) {
	return s.repo.ListByBaseID(baseID)
}

func (s *FactoryService) Update(id uint, baseID uint, code, name string) (*model.Factory, error) {
	factory, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Verify base exists
	if _, err := s.baseRepo.GetByID(baseID); err != nil {
		return nil, ErrNotFound
	}

	// Check if new code conflicts in the new base
	if factory.BaseID != baseID || factory.Code != code {
		if _, err := s.repo.GetByCode(baseID, code); err == nil {
			return nil, ErrDuplicateCode
		}
	}

	factory.BaseID = baseID
	factory.Code = code
	factory.Name = name

	if err := s.repo.Update(factory); err != nil {
		return nil, err
	}

	return factory, nil
}

func (s *FactoryService) Delete(id uint) error {
	// Check if has workshops
	workshopRepo := repository.NewWorkshopRepo()
	workshops, _ := workshopRepo.ListByFactoryID(id)
	if len(workshops) > 0 {
		return errors.New("cannot delete factory with existing workshops")
	}

	return s.repo.Delete(id)
}

// Workshop Service
type WorkshopService struct {
	repo        *repository.WorkshopRepository
	factoryRepo *repository.FactoryRepository
}

func NewWorkshopService() *WorkshopService {
	return &WorkshopService{
		repo:        repository.NewWorkshopRepo(),
		factoryRepo: repository.NewFactoryRepo(),
	}
}

func (s *WorkshopService) Create(factoryID uint, code, name string) (*model.Workshop, error) {
	// Verify factory exists
	if _, err := s.factoryRepo.GetByID(factoryID); err != nil {
		return nil, ErrNotFound
	}

	// Check if code exists in this factory
	if _, err := s.repo.GetByCode(factoryID, code); err == nil {
		return nil, ErrDuplicateCode
	}

	workshop := &model.Workshop{
		FactoryID: factoryID,
		Code:      code,
		Name:      name,
	}

	if err := s.repo.Create(workshop); err != nil {
		return nil, err
	}

	return workshop, nil
}

func (s *WorkshopService) GetByID(id uint) (*model.Workshop, error) {
	return s.repo.GetByID(id)
}

func (s *WorkshopService) List() ([]model.Workshop, error) {
	return s.repo.List()
}

func (s *WorkshopService) ListByFactoryID(factoryID uint) ([]model.Workshop, error) {
	return s.repo.ListByFactoryID(factoryID)
}

func (s *WorkshopService) Update(id uint, factoryID uint, code, name string) (*model.Workshop, error) {
	workshop, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Verify factory exists
	if _, err := s.factoryRepo.GetByID(factoryID); err != nil {
		return nil, ErrNotFound
	}

	// Check if new code conflicts in the new factory
	if workshop.FactoryID != factoryID || workshop.Code != code {
		if _, err := s.repo.GetByCode(factoryID, code); err == nil {
			return nil, ErrDuplicateCode
		}
	}

	workshop.FactoryID = factoryID
	workshop.Code = code
	workshop.Name = name

	if err := s.repo.Update(workshop); err != nil {
		return nil, err
	}

	return workshop, nil
}

func (s *WorkshopService) Delete(id uint) error {
	// Check if has equipment
	equipRepo := repository.NewEquipmentRepo()
	filter := repository.EquipmentFilter{
		WorkshopID: &id,
		Page:       1,
		PageSize:   1,
	}
	equipments, _, _ := equipRepo.List(filter)
	if len(equipments) > 0 {
		return errors.New("cannot delete workshop with existing equipment")
	}

	return s.repo.Delete(id)
}

// Equipment Type Service
type EquipmentTypeService struct {
	repo *repository.EquipmentTypeRepository
}

func NewEquipmentTypeService() *EquipmentTypeService {
	return &EquipmentTypeService{
		repo: repository.NewEquipmentTypeRepo(),
	}
}

func (s *EquipmentTypeService) Create(name, category string) (*model.EquipmentType, error) {
	equipmentType := &model.EquipmentType{
		Name:     name,
		Category: category,
	}

	if err := s.repo.Create(equipmentType); err != nil {
		return nil, err
	}

	return equipmentType, nil
}

func (s *EquipmentTypeService) GetByID(id uint) (*model.EquipmentType, error) {
	return s.repo.GetByID(id)
}

func (s *EquipmentTypeService) List() ([]model.EquipmentType, error) {
	return s.repo.List()
}

func (s *EquipmentTypeService) Update(id uint, name, category string) (*model.EquipmentType, error) {
	equipmentType, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	equipmentType.Name = name
	equipmentType.Category = category

	if err := s.repo.Update(equipmentType); err != nil {
		return nil, err
	}

	return equipmentType, nil
}

func (s *EquipmentTypeService) Delete(id uint) error {
	// Check if has equipment
	equipRepo := repository.NewEquipmentRepo()
	filter := repository.EquipmentFilter{
		TypeID:   &id,
		Page:     1,
		PageSize: 1,
	}
	equipments, _, _ := equipRepo.List(filter)
	if len(equipments) > 0 {
		return errors.New("cannot delete equipment type with existing equipment")
	}

	return s.repo.Delete(id)
}

// Equipment Service
type EquipmentService struct {
	repo       *repository.EquipmentRepository
	typeRepo   *repository.EquipmentTypeRepository
	workshopRepo *repository.WorkshopRepository
	userRepo   *UserRepository
}

func NewEquipmentService() *EquipmentService {
	return &EquipmentService{
		repo:         repository.NewEquipmentRepo(),
		typeRepo:     repository.NewEquipmentTypeRepo(),
		workshopRepo: repository.NewWorkshopRepo(),
		userRepo:     NewUserRepository(),
	}
}

func (s *EquipmentService) Create(req *CreateEquipmentRequest) (*model.Equipment, error) {
	// Validate type exists
	if _, err := s.typeRepo.GetByID(req.TypeID); err != nil {
		return nil, ErrNotFound
	}

	// Validate workshop exists
	if _, err := s.workshopRepo.GetByID(req.WorkshopID); err != nil {
		return nil, ErrNotFound
	}

	// Check if code exists
	if _, err := s.repo.GetByCode(req.Code); err == nil {
		return nil, ErrDuplicateCode
	}

	// Generate QR code
	qrCode := fmt.Sprintf("QR_%s", req.Code)

	// Validate maintenance user if provided
	if req.DedicatedMaintenanceID != nil {
		if _, err := s.userRepo.GetByID(*req.DedicatedMaintenanceID); err != nil {
			return nil, ErrNotFound
		}
	}

	equipment := &model.Equipment{
		Code:                   req.Code,
		Name:                   req.Name,
		TypeID:                 req.TypeID,
		WorkshopID:             req.WorkshopID,
		QRCode:                 qrCode,
		Spec:                   req.Spec,
		PurchaseDate:           req.PurchaseDate,
		Status:                 req.Status,
		DedicatedMaintenanceID:  req.DedicatedMaintenanceID,
	}

	if equipment.Status == "" {
		equipment.Status = "running"
	}

	if err := s.repo.Create(equipment); err != nil {
		return nil, err
	}

	return s.repo.GetByID(equipment.ID)
}

func (s *EquipmentService) GetByID(id uint) (*model.Equipment, error) {
	return s.repo.GetByID(id)
}

func (s *EquipmentService) GetByQRCode(qrCode string) (*model.Equipment, error) {
	return s.repo.GetByQRCode(qrCode)
}

func (s *EquipmentService) List(filter *EquipmentFilter) (*EquipmentListResult, error) {
	repoFilter := repository.EquipmentFilter{
		Code:       filter.Code,
		Name:       filter.Name,
		TypeID:     filter.TypeID,
		FactoryID:  filter.FactoryID,
		WorkshopID: filter.WorkshopID,
		Status:     filter.Status,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}

	equipments, total, err := s.repo.List(repoFilter)
	if err != nil {
		return nil, err
	}

	return &EquipmentListResult{
		Items: equipments,
		Total: total,
	}, nil
}

func (s *EquipmentService) Update(id uint, req *UpdateEquipmentRequest) (*model.Equipment, error) {
	equipment, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Validate type exists
	if _, err := s.typeRepo.GetByID(req.TypeID); err != nil {
		return nil, ErrNotFound
	}

	// Validate workshop exists
	if _, err := s.workshopRepo.GetByID(req.WorkshopID); err != nil {
		return nil, ErrNotFound
	}

	// Check if new code conflicts
	if equipment.Code != req.Code {
		if _, err := s.repo.GetByCode(req.Code); err == nil {
			return nil, ErrDuplicateCode
		}
	}

	// Validate maintenance user if provided
	if req.DedicatedMaintenanceID != nil {
		if _, err := s.userRepo.GetByID(*req.DedicatedMaintenanceID); err != nil {
			return nil, ErrNotFound
		}
	}

	equipment.Code = req.Code
	equipment.Name = req.Name
	equipment.TypeID = req.TypeID
	equipment.WorkshopID = req.WorkshopID
	equipment.Spec = req.Spec
	equipment.PurchaseDate = req.PurchaseDate
	equipment.Status = req.Status
	equipment.DedicatedMaintenanceID = req.DedicatedMaintenanceID

	if err := s.repo.Update(equipment); err != nil {
		return nil, err
	}

	return s.repo.GetByID(id)
}

func (s *EquipmentService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *EquipmentService) GetStatistics() (*EquipmentStatistics, error) {
	stats, err := s.repo.GetStatistics()
	if err != nil {
		return nil, err
	}

	return &EquipmentStatistics{
		Total:      stats["total"],
		Running:    stats["running"],
		Stopped:    stats["stopped"],
		Maintenance: stats["maintenance"],
		Scrapped:   stats["scrapped"],
	}, nil
}

func (s *EquipmentService) GetByTypeStatistics() ([]TypeStatistics, error) {
	results, err := s.repo.GetByTypeStatistics()
	if err != nil {
		return nil, err
	}

	typeStats := make([]TypeStatistics, len(results))
	for i, r := range results {
		typeStats[i] = TypeStatistics{
			TypeName:     r["type_name"].(string),
			Count:        r["count"].(int64),
			Running:      r["running"].(int64),
			Maintenance:  r["maintenance"].(int64),
		}
	}

	return typeStats, nil
}

// Request/Response types for service
type CreateEquipmentRequest struct {
	Code                  string     `json:"code"`
	Name                  string     `json:"name"`
	TypeID                uint       `json:"type_id"`
	WorkshopID            uint       `json:"workshop_id"`
	Spec                  string     `json:"spec"`
	PurchaseDate          *time.Time `json:"purchase_date"`
	Status                string     `json:"status"`
	DedicatedMaintenanceID *uint      `json:"dedicated_maintenance_id"`
}

type UpdateEquipmentRequest struct {
	Code                  string     `json:"code"`
	Name                  string     `json:"name"`
	TypeID                uint       `json:"type_id"`
	WorkshopID            uint       `json:"workshop_id"`
	Spec                  string     `json:"spec"`
	PurchaseDate          *time.Time `json:"purchase_date"`
	Status                string     `json:"status"`
	DedicatedMaintenanceID *uint      `json:"dedicated_maintenance_id"`
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

type EquipmentListResult struct {
	Items []model.Equipment
	Total int64
}

type EquipmentStatistics struct {
	Total      int64
	Running    int64
	Stopped    int64
	Maintenance int64
	Scrapped   int64
}

type TypeStatistics struct {
	TypeName    string
	Count       int64
	Running     int64
	Maintenance int64
}
