package service

import (
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
)

// KnowledgeService
type KnowledgeService struct {
	knowledgeRepo *repository.KnowledgeArticleRepository
	userRepo      *UserRepository
	repairRepo    *repository.RepairOrderRepository
}

func NewKnowledgeService() *KnowledgeService {
	return &KnowledgeService{
		knowledgeRepo: repository.NewKnowledgeArticleRepository(),
		userRepo:      NewUserRepository(),
		repairRepo:    repository.NewRepairOrderRepository(),
	}
}

func (s *KnowledgeService) CreateArticle(title string, equipmentTypeID *uint, faultPhenomenon, causeAnalysis, solution, sourceType string, sourceID *uint, tags []string, createdBy uint) (*model.KnowledgeArticle, error) {
	article := &model.KnowledgeArticle{
		Title:           title,
		EquipmentTypeID: equipmentTypeID,
		FaultPhenomenon: faultPhenomenon,
		CauseAnalysis:   causeAnalysis,
		Solution:        solution,
		SourceType:      sourceType,
		SourceID:        sourceID,
		Tags:            tags,
		CreatedBy:       createdBy,
	}

	if err := s.knowledgeRepo.Create(article); err != nil {
		return nil, err
	}

	return s.knowledgeRepo.GetByID(article.ID)
}

func (s *KnowledgeService) GetArticleByID(id uint) (*model.KnowledgeArticle, error) {
	return s.knowledgeRepo.GetByID(id)
}

func (s *KnowledgeService) ListArticles(filter repository.KnowledgeArticleFilter) (*KnowledgeArticleListResult, error) {
	articles, total, err := s.knowledgeRepo.List(filter)
	if err != nil {
		return nil, err
	}

	return &KnowledgeArticleListResult{
		Items: articles,
		Total: total,
	}, nil
}

func (s *KnowledgeService) UpdateArticle(id uint, title string, equipmentTypeID *uint, faultPhenomenon, causeAnalysis, solution string, tags []string) error {
	article, err := s.knowledgeRepo.GetByID(id)
	if err != nil {
		return ErrNotFound
	}

	article.Title = title
	article.EquipmentTypeID = equipmentTypeID
	article.FaultPhenomenon = faultPhenomenon
	article.CauseAnalysis = causeAnalysis
	article.Solution = solution
	article.Tags = tags

	return s.knowledgeRepo.Update(article)
}

func (s *KnowledgeService) DeleteArticle(id uint) error {
	return s.knowledgeRepo.Delete(id)
}

func (s *KnowledgeService) SearchArticles(keyword string) ([]model.KnowledgeArticle, error) {
	return s.knowledgeRepo.Search(keyword)
}

func (s *KnowledgeService) ConvertFromRepair(orderID uint, title, faultPhenomenon, causeAnalysis string, tags []string, createdBy uint) (*model.KnowledgeArticle, error) {
	// Get repair order
	repair, err := s.repairRepo.GetByID(orderID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Build solution from repair order
	solution := "故障现象: " + repair.FaultDescription + "\n"
	if repair.Solution != "" {
		solution += "解决方法: " + repair.Solution + "\n"
	}
	if repair.FaultCode != "" {
		solution += "故障代码: " + repair.FaultCode
	}

	article := &model.KnowledgeArticle{
		Title:           title,
		EquipmentTypeID: nil, // Can be derived from equipment
		FaultPhenomenon: faultPhenomenon,
		CauseAnalysis:   causeAnalysis,
		Solution:        solution,
		SourceType:      "repair",
		SourceID:        &orderID,
		Tags:            tags,
		CreatedBy:       createdBy,
	}

	if err := s.knowledgeRepo.Create(article); err != nil {
		return nil, err
	}

	return s.knowledgeRepo.GetByID(article.ID)
}

// Types
type KnowledgeArticleListResult struct {
	Items []model.KnowledgeArticle
	Total int64
}
