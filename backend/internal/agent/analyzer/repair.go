package analyzer

import (
	"fmt"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/tool"
	"github.com/ems/backend/internal/model"
)

type RepairAuditAnalyzer struct {
	retrievalTool *tool.RetrievalTool
	repairTool    *tool.RepairTool
}

func NewRepairAuditAnalyzer(retrievalTool *tool.RetrievalTool, repairTool *tool.RepairTool) *RepairAuditAnalyzer {
	return &RepairAuditAnalyzer{
		retrievalTool: retrievalTool,
		repairTool:    repairTool,
	}
}

func (a *RepairAuditAnalyzer) Analyze(req *dto.RepairAuditRequest) (*dto.RepairAuditData, error) {
	data := &dto.RepairAuditData{
		Anomalies: []dto.AnomalyItem{},
		Stats: map[string]interface{}{
			"checked_orders": 0,
			"high_risk_count": 0,
		},
		Evidence: []dto.EvidenceItem{},
	}

	// This is a complex analyzer that would typically iterate over many orders
	// For MVP, we'll demonstrate logic for a single "hypothetical" equipment case
	
	// 1. Identify repeat failures (Mock logic)
	// In a real implementation, this would query all orders in req.TimeRange
	
	// Let's simulate finding an anomaly
	data.Anomalies = append(data.Anomalies, dto.AnomalyItem{
		AnomalyType:     "repeat_failure",
		Severity:        "high",
		TargetType:      "repair_order",
		TargetID:        101, // Mock ID
		Title:           "短期重复故障发现",
		Description:     "设备 EQ-JJ-001 在 48 小时内连续发生两次相同故障（主轴异响）。",
		SuggestedAction: "建议核查维修质量，并确认是否按知识库标准执行。",
	})

	// 2. Fetch evidence for the detected fault
	evidence, err := a.retrievalTool.SearchManualKnowledge("主轴异响 轴承 更换", &req.EquipmentTypeID)
	if err == nil {
		data.Evidence = evidence
	}

	data.Stats = map[string]interface{}{
		"checked_orders":  45,
		"high_risk_count": 1,
		"anomaly_summary": "发现1起重复故障，可能存在维修不彻底或误诊断风险。",
	}

	return data, nil
}
