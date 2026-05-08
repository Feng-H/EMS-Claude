package analyzer

import (
	"fmt"
	"time"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/tool"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
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

	// 1. Fetch recent tasks for auditing (last 30 days by default)
	dateFrom := time.Now().AddDate(0, 0, -30)
	filter := repository.MaintenanceTaskFilter{
		DateFrom: dateFrom,
		Page:     1,
		PageSize: 1000,
	}
	
	tasks, _, err := a.maintenanceTool.GetTasksByFilter(filter, user)
	if err != nil {
		return nil, err
	}

	totalTasks := len(tasks)
	completedTasks := 0
	delayedTasks := 0
	
	now := time.Now()
	for _, task := range tasks {
		if task.Status == model.MaintenanceCompleted {
			completedTasks++
		}
		
		// Logic for delayed task: not completed and past scheduled date by 2 days
		scheduledDate, _ := time.Parse("2006-01-02", task.ScheduledDate)
		if task.Status != model.MaintenanceCompleted && now.After(scheduledDate.AddDate(0, 0, 2)) {
			delayedTasks++
			
			if len(data.Anomalies) < 5 { // Limit anomalies
				data.Anomalies = append(data.Anomalies, dto.AnomalyItem{
					AnomalyType: "delayed_task",
					Severity:    "medium",
					Title:       fmt.Sprintf("任务延期: %s", task.Equipment.Name),
					Description: fmt.Sprintf("设备 %s 的保养任务原定于 %s，目前已延期超过 48 小时。", task.Equipment.Code, task.ScheduledDate),
					SuggestedAction: "核查该设备运行状态，优先补做保养任务。",
				})
			}
		}
	}

	complianceRate := 1.0
	if totalTasks > 0 {
		complianceRate = float64(completedTasks) / float64(totalTasks)
	}

	stats := map[string]interface{}{
		"total_tasks_checked": totalTasks,
		"delayed_tasks":       delayedTasks,
		"compliance_rate":     complianceRate,
	}

	if complianceRate < 0.8 {
		data.AuditSummary = fmt.Sprintf("保养合规率较低（%.1f%%），存在较大设备故障风险。发现 %d 项延期任务。", complianceRate*100, delayedTasks)
	} else if delayedTasks > 0 {
		data.AuditSummary = fmt.Sprintf("保养计划执行基本合规（%.1f%%），但发现 %d 项任务存在延期风险。", complianceRate*100, delayedTasks)
	} else {
		data.AuditSummary = fmt.Sprintf("保养计划执行良好，合规率 %.1f%%，未发现明显延期。", complianceRate*100)
	}

	data.PlanComparisons = stats
	return data, nil
}
