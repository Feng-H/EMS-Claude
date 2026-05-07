package tool

import (
	"fmt"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/repository"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/memory"
	internalRepo "github.com/ems/backend/internal/repository"
	"sort"
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
func (t *RetrievalTool) GetEquipmentProfile(id uint, user model.User) (map[string]interface{}, error) {
	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		e := store.FindEquipment(id)
		if e == nil { return nil, fmt.Errorf("equipment not found") }
		
		// Permission check
		if user.Role != "admin" && user.FactoryID != nil {
			workshop, ok := store.Workshops[e.WorkshopID]
			if !ok || workshop.FactoryID != *user.FactoryID {
				return nil, fmt.Errorf("access denied: equipment belongs to another factory")
			}
		}

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
	
	// Permission check
	if user.Role != "admin" && user.FactoryID != nil {
		if e.Workshop.FactoryID != *user.FactoryID {
			return nil, fmt.Errorf("access denied: equipment belongs to another factory")
		}
	}

	res := map[string]interface{}{
		"id": e.ID, "code": e.Code, "name": e.Name, "status": e.Status,
		"type_name": e.Type.Name,
		"workshop_name": e.Workshop.Name,
		"factory_name": e.Workshop.Factory.Name,
	}
	return res, nil
}

// SearchManualKnowledge searches both knowledge articles and manual chunks with weighted ranking
func (t *RetrievalTool) SearchManualKnowledge(query string, equipmentTypeID *uint, user model.User) ([]dto.EvidenceItem, error) {
	var results []dto.EvidenceItem
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return results, nil
	}

	// 1. 获取候选条目
	articles, _ := t.searchKnowledge(query, equipmentTypeID)
	chunks, _ := t.agentRepo.SearchManualChunks(query, equipmentTypeID)

	// 2. 混合检索打分逻辑
	// 知识库权重较高，因为是经过人工审核的专家经验
	for _, art := range articles {
		score := 0.70 // 基础分
		title := strings.ToLower(art.Title)
		phenom := strings.ToLower(art.FaultPhenomenon)
		solution := strings.ToLower(art.Solution)

		// 标题完全匹配或包含关键短语
		if title == query {
			score += 0.25
		} else if strings.Contains(title, query) {
			score += 0.15
		}

		// 现象描述匹配
		if strings.Contains(phenom, query) {
			score += 0.10
		}

		// 解决方案匹配 (权重略低)
		if strings.Contains(solution, query) {
			score += 0.05
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

	// 技术手册权重居中，覆盖面广但可能冗余
	for _, chunk := range chunks {
		score := 0.50 // 基础分
		title := strings.ToLower(chunk.SectionTitle)
		content := strings.ToLower(chunk.Content)

		if strings.Contains(title, query) {
			score += 0.30
		} else if strings.Contains(content, query) {
			score += 0.15
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
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// 4. 截断：仅返回前 5 条最相关的证据
	if len(results) > 5 {
		results = results[:5]
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
