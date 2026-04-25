package analyzer

import (
	"fmt"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/tool"
)

type MaintenanceAnalyzer struct {
	retrievalTool   *tool.RetrievalTool
	maintenanceTool *tool.MaintenanceTool
}

func NewMaintenanceAnalyzer(retrievalTool *tool.RetrievalTool, maintenanceTool *tool.MaintenanceTool) *MaintenanceAnalyzer {
	return &MaintenanceAnalyzer{
		retrievalTool:   retrievalTool,
		maintenanceTool: maintenanceTool,
	}
}

func (a *MaintenanceAnalyzer) Analyze(req *dto.MaintenanceRecommendRequest) (*dto.MaintenanceRecommendData, error) {
	data := &dto.MaintenanceRecommendData{
		Recommendations: []dto.RecommendationItem{},
		ExpectedBenefits: []string{
			"提高设备可用性",
			"降低突发故障率",
		},
		Evidence: []dto.EvidenceItem{},
	}

	// 1. Fetch current plan
	plan, err := a.maintenanceTool.GetPlanByEquipmentType(req.EquipmentTypeID)
	if err == nil && plan != nil {
		data.CurrentPlan = map[string]interface{}{
			"plan_id":    plan.ID,
			"plan_name":  plan.Name,
			"cycle_days": plan.CycleDays,
			"level":      plan.Level,
			"item_count": len(plan.Items),
		}
	}

	// 2. Search for evidence in knowledge base/manuals
	evidence, err := a.retrievalTool.SearchManualKnowledge("保养周期 维护项", &req.EquipmentTypeID)
	if err == nil {
		data.Evidence = evidence
	}

	// 3. Simple rule-based recommendation (Mock for MVP)
	// In a real implementation, this would compare current plan vs evidence
	if plan != nil && plan.CycleDays > 30 {
		data.Recommendations = append(data.Recommendations, dto.RecommendationItem{
			Type:        "cycle_adjustment",
			Target:      "maintenance_plan",
			TargetID:    plan.ID,
			Title:       "缩短保养周期",
			Description: fmt.Sprintf("建议将 %d 天调整为 30 天", plan.CycleDays),
			Reason:      "根据同类设备最佳实践，30天是更稳健的维护频率。",
			Impact:      "预计降低 15% 的重复故障率。",
		})
	}

	return data, nil
}
