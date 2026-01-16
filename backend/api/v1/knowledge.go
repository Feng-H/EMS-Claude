package v1

import (
	"net/http"
	"strconv"

	"github.com/ems/backend/internal/dto"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/internal/repository"
	"github.com/ems/backend/internal/service"
	"github.com/gin-gonic/gin"
)

var (
	knowledgeService *service.KnowledgeService
)

func InitKnowledge() {
	knowledgeService = service.NewKnowledgeService()
}

// =====================================================
// Knowledge Base APIs
// =====================================================

// ListKnowledgeArticles returns all knowledge articles
// @Summary Get knowledge articles
// @Tags knowledge
// @Router /knowledge [get]
func ListKnowledgeArticles(c *gin.Context) {
	var query dto.KnowledgeArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	filter := repository.KnowledgeArticleFilter{
		Keyword:          query.Keyword,
		EquipmentTypeID:   query.EquipmentTypeID,
		Tag:              query.Tag,
		SourceType:       query.SourceType,
		Page:             query.Page,
		PageSize:         query.PageSize,
	}

	result, err := knowledgeService.ListArticles(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.KnowledgeArticleResponse, len(result.Items))
	for i, a := range result.Items {
		response[i] = dto.KnowledgeArticleResponse{
			ID:              a.ID,
			Title:           a.Title,
			EquipmentTypeID: a.EquipmentTypeID,
			FaultPhenomenon: a.FaultPhenomenon,
			CauseAnalysis:   a.CauseAnalysis,
			Solution:        a.Solution,
			SourceType:      a.SourceType,
			SourceID:        a.SourceID,
			Tags:            a.Tags,
			CreatedBy:       a.CreatedBy,
			CreatedAt:       a.CreatedAt,
			UpdatedAt:       a.UpdatedAt,
		}
		if a.EquipmentType != nil {
			response[i].EquipmentTypeName = a.EquipmentType.Name
		}
		if a.Creator.ID > 0 {
			response[i].CreatorName = a.Creator.Name
		}
	}

	c.JSON(http.StatusOK, dto.KnowledgeArticleListResponse{
		Total: result.Total,
		Items: response,
	})
}

// GetKnowledgeArticle returns a single knowledge article
// @Summary Get knowledge article
// @Tags knowledge
// @Router /knowledge/:id [get]
func GetKnowledgeArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	article, err := knowledgeService.GetArticleByID(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.KnowledgeArticleResponse{
		ID:              article.ID,
		Title:           article.Title,
		EquipmentTypeID: article.EquipmentTypeID,
		FaultPhenomenon: article.FaultPhenomenon,
		CauseAnalysis:   article.CauseAnalysis,
		Solution:        article.Solution,
		SourceType:      article.SourceType,
		SourceID:        article.SourceID,
		Tags:            article.Tags,
		CreatedBy:       article.CreatedBy,
		CreatedAt:       article.CreatedAt,
		UpdatedAt:       article.UpdatedAt,
	}
	if article.EquipmentType != nil {
		response.EquipmentTypeName = article.EquipmentType.Name
	}
	if article.Creator.ID > 0 {
		response.CreatorName = article.Creator.Name
	}

	c.JSON(http.StatusOK, response)
}

// CreateKnowledgeArticle creates a new knowledge article
// @Summary Create knowledge article
// @Tags knowledge
// @Router /knowledge [post]
func CreateKnowledgeArticle(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.CreateKnowledgeArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article, err := knowledgeService.CreateArticle(
		req.Title,
		req.EquipmentTypeID,
		req.FaultPhenomenon,
		req.CauseAnalysis,
		req.Solution,
		req.SourceType,
		req.SourceID,
		req.Tags,
		userID,
	)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.KnowledgeArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		Solution:  article.Solution,
		SourceType: article.SourceType,
		Tags:      article.Tags,
		CreatedBy: article.CreatedBy,
		CreatedAt: article.CreatedAt,
	})
}

// UpdateKnowledgeArticle updates a knowledge article
// @Summary Update knowledge article
// @Tags knowledge
// @Router /knowledge/:id [put]
func UpdateKnowledgeArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.UpdateKnowledgeArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := knowledgeService.UpdateArticle(uint(id), req.Title, req.EquipmentTypeID, req.FaultPhenomenon, req.CauseAnalysis, req.Solution, req.Tags); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

// DeleteKnowledgeArticle deletes a knowledge article
// @Summary Delete knowledge article
// @Tags knowledge
// @Router /knowledge/:id [delete]
func DeleteKnowledgeArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if role, _ := middleware.GetUserRole(c); role != "admin" && role != "engineer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	if err := knowledgeService.DeleteArticle(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// SearchKnowledgeArticles searches knowledge articles
// @Summary Search knowledge articles
// @Tags knowledge
// @Router /knowledge/search [get]
func SearchKnowledgeArticles(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keyword is required"})
		return
	}

	articles, err := knowledgeService.SearchArticles(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.KnowledgeArticleResponse, len(articles))
	for i, a := range articles {
		response[i] = dto.KnowledgeArticleResponse{
			ID:              a.ID,
			Title:           a.Title,
			EquipmentTypeID: a.EquipmentTypeID,
			FaultPhenomenon: a.FaultPhenomenon,
			CauseAnalysis:   a.CauseAnalysis,
			Solution:        a.Solution,
			SourceType:      a.SourceType,
			SourceID:        a.SourceID,
			Tags:            a.Tags,
			CreatedBy:       a.CreatedBy,
			CreatedAt:       a.CreatedAt,
			UpdatedAt:       a.UpdatedAt,
		}
		if a.EquipmentType != nil {
			response[i].EquipmentTypeName = a.EquipmentType.Name
		}
		if a.Creator.ID > 0 {
			response[i].CreatorName = a.Creator.Name
		}
	}

	c.JSON(http.StatusOK, response)
}

// ConvertFromRepair converts a repair order to a knowledge article
// @Summary Convert from repair order
// @Tags knowledge
// @Router /knowledge/convert-repair [post]
func ConvertFromRepair(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.ConvertFromRepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article, err := knowledgeService.ConvertFromRepair(
		req.OrderID,
		req.Title,
		req.FaultPhenomenon,
		req.CauseAnalysis,
		req.Tags,
		userID,
	)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.KnowledgeArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		Solution:  article.Solution,
		SourceType: article.SourceType,
		SourceID:  article.SourceID,
		Tags:      article.Tags,
		CreatedBy: article.CreatedBy,
		CreatedAt: article.CreatedAt,
	})
}
