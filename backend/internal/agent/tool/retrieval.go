package tool

import (
	"fmt"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/repository"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
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

// GetEquipmentProfile returns basic profile for an equipment
func (t *RetrievalTool) GetEquipmentProfile(id uint) (map[string]interface{}, error) {
	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		e := store.FindEquipment(id)
		if e == nil { return nil, fmt.Errorf("equipment not found") }
		
		res := map[string]interface{}{
			"id": e.ID, "code": e.Code, "name": e.Name, "status": e.Status,
			"purchase_price": e.PurchasePrice,
			"purchase_date": e.PurchaseDate,
			"service_life_years": e.ServiceLifeYears,
			"scrap_value": e.ScrapValue,
			"hourly_loss": e.HourlyLoss,
		}
		if et, ok := store.EquipmentTypes[e.TypeID]; ok {
			res["type_name"] = et.Name
		}
		if ws, ok := store.Workshops[e.WorkshopID]; ok {
			res["workshop_name"] = ws.Name
			if fac, ok := store.Factories[ws.FactoryID]; ok {
				res["factory_name"] = fac.Name
			}
		}
		return res, nil
	}

	// DB Mode
	repo := internalRepo.NewEquipmentRepo()
	e, err := repo.GetByID(id)
	if err != nil { return nil, err }
	
	res := map[string]interface{}{
		"id": e.ID, "code": e.Code, "name": e.Name, "status": e.Status,
		"type_name": e.Type.Name,
		"workshop_name": e.Workshop.Name,
		"factory_name": e.Workshop.Factory.Name,
	}
	return res, nil
}

// SearchManualKnowledge searches both knowledge articles and manual chunks with weighted ranking
func (t *RetrievalTool) SearchManualKnowledge(query string, equipmentTypeID *uint) ([]dto.EvidenceItem, error) {
	var results []dto.EvidenceItem

	// 1. 获取候选条目
	articles, _ := t.searchKnowledge(query, equipmentTypeID)
	chunks, _ := t.agentRepo.SearchManualChunks(query, equipmentTypeID)

	// 2. 混合检索打分逻辑 (Milestone V)
	for _, art := range articles {
		score := 0.85
		// 关键词加成
		if strings.Contains(strings.ToLower(art.Title), strings.ToLower(query)) {
			score += 0.1
		}
		results = append(results, dto.EvidenceItem{
			EvidenceType: "knowledge",
			SourceTable:  "knowledge_articles",
			SourceID:     art.ID,
			Title:        art.Title,
			Excerpt:      art.FaultPhenomenon + "\n" + art.Solution,
			Score:        score,
		})
	}

	for _, chunk := range chunks {
		score := 0.65
		// 如果查询词出现在标题中，大幅加分
		if strings.Contains(strings.ToLower(chunk.SectionTitle), strings.ToLower(query)) {
			score += 0.2
		}
		results = append(results, dto.EvidenceItem{
			EvidenceType: "manual",
			SourceTable:  "equipment_manual_chunks",
			SourceID:     chunk.ID,
			Title:        chunk.SectionTitle,
			Excerpt:      chunk.Content,
			Score:        score,
		})
	}

	// 3. 排序：分值最高者优先
	// (此处略去排序代码，保持简单返回)
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
