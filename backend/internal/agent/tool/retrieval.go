package tool

import (
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/repository"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/database"
	"github.com/ems/backend/pkg/memory"
	internalRepo "github.com/ems/backend/internal/repository"
	"strings"
)

type RetrievalTool struct {
	agentRepo     repository.IAgentRepository
	knowledgeRepo *internalRepo.KnowledgeArticleRepository
}

func NewRetrievalTool(agentRepo repository.IAgentRepository) *RetrievalTool {
	var knowledgeRepo *internalRepo.KnowledgeArticleRepository
	if config.Cfg.Storage.Mode != "memory" {
		knowledgeRepo = internalRepo.NewKnowledgeArticleRepository()
	}
	
	return &RetrievalTool{
		agentRepo:     agentRepo,
		knowledgeRepo: knowledgeRepo,
	}
}

// SearchManualKnowledge searches both knowledge articles and manual chunks
func (t *RetrievalTool) SearchManualKnowledge(query string, equipmentTypeID *uint) ([]dto.EvidenceItem, error) {
	var results []dto.EvidenceItem

	// 1. Search knowledge articles (higher priority)
	articles, err := t.searchKnowledge(query, equipmentTypeID)
	if err == nil {
		for _, art := range articles {
			results = append(results, dto.EvidenceItem{
				EvidenceType: "knowledge",
				SourceTable:  "knowledge_articles",
				SourceID:     art.ID,
				Title:        art.Title,
				Excerpt:      art.FaultPhenomenon + "\n" + art.Solution,
				Score:        0.9, // Higher score for knowledge base
			})
		}
	}

	// 2. Search manual chunks
	chunks, err := t.agentRepo.SearchManualChunks(query, equipmentTypeID)
	if err == nil {
		for _, chunk := range chunks {
			results = append(results, dto.EvidenceItem{
				EvidenceType: "manual",
				SourceTable:  "equipment_manual_chunks",
				SourceID:     chunk.ID,
				Title:        chunk.SectionTitle,
				Excerpt:      chunk.Content,
				Score:        0.7,
			})
		}
	}

	return results, nil
}

func (t *RetrievalTool) searchKnowledge(query string, equipmentTypeID *uint) ([]model.KnowledgeArticle, error) {
	if config.Cfg.Storage.Mode == "memory" {
		var results []model.KnowledgeArticle
		store := memory.GetStore()
		count := 0
		for _, art := range store.KnowledgeArticles {
			if count >= 5 {
				break
			}
			match := strings.Contains(strings.ToLower(art.Title), strings.ToLower(query)) ||
				strings.Contains(strings.ToLower(art.FaultPhenomenon), strings.ToLower(query)) ||
				strings.Contains(strings.ToLower(art.Solution), strings.ToLower(query))
			
			if match {
				if equipmentTypeID != nil && art.EquipmentTypeID != nil && *art.EquipmentTypeID != *equipmentTypeID {
					continue
				}
				results = append(results, *art)
				count++
			}
		}
		return results, nil
	}
	
	// Database mode
	return t.knowledgeRepo.Search(query)
}
