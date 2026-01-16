package dto

import "time"

// =====================================================
// Knowledge Base DTOs
// =====================================================

// CreateKnowledgeArticleRequest creates a new knowledge article
type CreateKnowledgeArticleRequest struct {
	Title             string   `json:"title" binding:"required"`
	EquipmentTypeID   *uint    `json:"equipment_type_id"`
	FaultPhenomenon   string  `json:"fault_phenomenon"`
	CauseAnalysis     string  `json:"cause_analysis"`
	Solution          string  `json:"solution" binding:"required"`
	SourceType        string  `json:"source_type"` // repair, manual, other
	SourceID          *uint    `json:"source_id"`
	Tags              []string `json:"tags"`
}

// UpdateKnowledgeArticleRequest updates a knowledge article
type UpdateKnowledgeArticleRequest struct {
	Title             string   `json:"title" binding:"required"`
	EquipmentTypeID   *uint    `json:"equipment_type_id"`
	FaultPhenomenon   string  `json:"fault_phenomenon"`
	CauseAnalysis     string  `json:"cause_analysis"`
	Solution          string  `json:"solution" binding:"required"`
	Tags              []string `json:"tags"`
}

// KnowledgeArticleResponse represents a knowledge article in API responses
type KnowledgeArticleResponse struct {
	ID                uint      `json:"id"`
	Title             string    `json:"title"`
	EquipmentTypeID   *uint     `json:"equipment_type_id"`
	EquipmentTypeName string    `json:"equipment_type_name,omitempty"`
	FaultPhenomenon   string    `json:"fault_phenomenon,omitempty"`
	CauseAnalysis     string    `json:"cause_analysis,omitempty"`
	Solution          string    `json:"solution"`
	SourceType        string    `json:"source_type"`
	SourceID          *uint     `json:"source_id,omitempty"`
	Tags              []string  `json:"tags"`
	CreatedBy         uint      `json:"created_by"`
	CreatorName       string    `json:"creator_name,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// KnowledgeArticleListResponse represents a paginated list of knowledge articles
type KnowledgeArticleListResponse struct {
	Total int64                       `json:"total"`
	Items []KnowledgeArticleResponse `json:"items"`
}

// ConvertFromRepairRequest represents a request to convert a repair order to knowledge
type ConvertFromRepairRequest struct {
	OrderID           uint   `json:"order_id" binding:"required"`
	Title             string `json:"title" binding:"required"`
	FaultPhenomenon   string `json:"fault_phenomenon"`
	CauseAnalysis     string `json:"cause_analysis"`
	Tags              []string `json:"tags"`
}

// =====================================================
// Query DTOs
// =====================================================

// KnowledgeArticleQuery represents query parameters for knowledge articles
type KnowledgeArticleQuery struct {
	Keyword          string `form:"keyword"`
	EquipmentTypeID *uint  `form:"equipment_type_id"`
	Tag              string `form:"tag"`
	SourceType       string `form:"source_type"`
	Page             int    `form:"page"`
	PageSize         int    `form:"page_size"`
}
