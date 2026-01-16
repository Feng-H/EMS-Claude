package repository

import (
	"strings"

	"github.com/ems/backend/internal/model"
	"gorm.io/gorm"
)

// KnowledgeArticle Repository
type KnowledgeArticleRepository struct {
	db *gorm.DB
}

func NewKnowledgeArticleRepository() *KnowledgeArticleRepository {
	return &KnowledgeArticleRepository{db: DB}
}

func (r *KnowledgeArticleRepository) Create(article *model.KnowledgeArticle) error {
	return r.db.Create(article).Error
}

func (r *KnowledgeArticleRepository) GetByID(id uint) (*model.KnowledgeArticle, error) {
	var article model.KnowledgeArticle
	err := r.db.Preload("EquipmentType").
		Preload("Creator").
		First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *KnowledgeArticleRepository) List(filter KnowledgeArticleFilter) ([]model.KnowledgeArticle, int64, error) {
	var articles []model.KnowledgeArticle
	var total int64

	query := r.db.Model(&model.KnowledgeArticle{})

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("title LIKE ? OR fault_phenomenon LIKE ? OR solution LIKE ?", keyword, keyword, keyword)
	}
	if filter.EquipmentTypeID != nil {
		query = query.Where("equipment_type_id = ?", *filter.EquipmentTypeID)
	}
	if filter.Tag != "" {
		query = query.Where("? = ANY(tags)", filter.Tag)
	}
	if filter.SourceType != "" {
		query = query.Where("source_type = ?", filter.SourceType)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (filter.Page - 1) * filter.PageSize
	err := query.Preload("EquipmentType").
		Preload("Creator").
		Order("created_at DESC").
		Offset(offset).Limit(filter.PageSize).
		Find(&articles).Error

	return articles, total, err
}

func (r *KnowledgeArticleRepository) Update(article *model.KnowledgeArticle) error {
	return r.db.Save(article).Error
}

func (r *KnowledgeArticleRepository) Delete(id uint) error {
	return r.db.Delete(&model.KnowledgeArticle{}, id).Error
}

// Search performs full-text search on articles
func (r *KnowledgeArticleRepository) Search(keyword string) ([]model.KnowledgeArticle, error) {
	var articles []model.KnowledgeArticle

	// Use GORM LIKE for simpler search (PostgreSQL FTS can be added later)
	likeKeyword := "%" + strings.ToUpper(keyword) + "%"
	err := r.db.Preload("EquipmentType").
		Preload("Creator").
		Where("UPPER(title) LIKE ? OR UPPER(fault_phenomenon) LIKE ? OR UPPER(solution) LIKE ?",
			likeKeyword, likeKeyword, likeKeyword).
		Order("created_at DESC").
		Limit(50).
		Find(&articles).Error

	return articles, err
}

// GetByTags returns articles with specific tags
func (r *KnowledgeArticleRepository) GetByTags(tags []string) ([]model.KnowledgeArticle, error) {
	var articles []model.KnowledgeArticle

	err := r.db.Preload("EquipmentType").
		Preload("Creator").
		Where("? = ANY(tags)", tags[0]).
		Order("created_at DESC").
		Find(&articles).Error

	return articles, err
}

// KnowledgeArticleFilter represents filter parameters for knowledge articles
type KnowledgeArticleFilter struct {
	Keyword          string
	EquipmentTypeID *uint
	Tag              string
	SourceType       string
	Page             int
	PageSize         int
}
