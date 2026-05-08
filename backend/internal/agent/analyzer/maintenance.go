package analyzer

import (
	"fmt"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/tool"
	"github.com/ems/backend/internal/model"
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

func (a *MaintenanceAnalyzer) Analyze(req *dto.MaintenanceRecommendRequest, user model.User) (*dto.MaintenanceRecommendData, error) {
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
	evidence, err := a.retrievalTool.SearchManualKnowledge("保养周期 维护项", &req.EquipmentTypeID, user)
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

func (a *MaintenanceAnalyzer) Audit(req *dto.MaintenanceAuditRequest, user model.User) (*dto.MaintenanceAuditData, error) {
	data := &dto.MaintenanceAuditData{
		Anomalies: []dto.AnomalyItem{},
		Evidence:  []dto.EvidenceItem{},
	}

	// 1. Fetch recent tasks for auditing
	// For MVP, we'll use rule-based analysis on compliance and delays.
	
	stats := map[string]interface{}{
		"total_tasks_checked": 0,
		"delayed_tasks":       0,
		"compliance_rate":     1.0,
	}

	// Simulated logic for MVP audit
	data.AuditSummary = "保养计划执行基本合规，但发现部分任务存在延期风险。"
	
	data.Anomalies = append(data.Anomalies, dto.AnomalyItem{
		AnomalyType: "delayed_task",
		Severity:    "medium",
		Title:       "发现保养任务延期",
		Description: "当前工厂有 3 个保养任务超过预计开始时间 48 小时未启动。",
		SuggestedAction: "建议核查维护班组负荷，必要时调整排班或外协维护。",
	})

	stats["total_tasks_checked"] = 45
	stats["delayed_tasks"] = 3
	stats["compliance_rate"] = 0.93

	data.PlanComparisons = stats

	return data, nil
}
