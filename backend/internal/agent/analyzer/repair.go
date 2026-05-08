package analyzer

import (
	"fmt"
	"time"
	"strings"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/tool"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
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

func (a *RepairAuditAnalyzer) Analyze(req *dto.RepairAuditRequest, user model.User) (*dto.RepairAuditData, error) {
	statsMap := map[string]interface{}{
		"checked_orders":  0,
		"high_risk_count": 0,
	}
	data := &dto.RepairAuditData{
		Anomalies: []dto.AnomalyItem{},
		Stats:     statsMap,
		Evidence:  []dto.EvidenceItem{},
	}

	// Fetch orders for the last 30 days if no range provided
	now := time.Now()
	dateFrom := now.AddDate(0, 0, -30)
	if req.TimeRange.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.TimeRange.StartDate); err == nil {
			dateFrom = t
		}
	}

	filter := repository.RepairOrderFilter{
		DateFrom: dateFrom,
		Page:     1,
		PageSize: 100,
	}

	orders, total, err := a.repairTool.GetOrdersByFilter(filter, user)
	if err != nil {
		return nil, err
	}

	statsMap["checked_orders"] = total

	// 1. Identify repeat failures
	// Group orders by equipment
	eqOrders := make(map[uint][]model.RepairOrder)
	for _, o := range orders {
		eqOrders[o.EquipmentID] = append(eqOrders[o.EquipmentID], o)
	}

	for eqID, oList := range eqOrders {
		if len(oList) < 2 { continue }
		
		// Sort by creation date (desc)
		for i := 0; i < len(oList)-1; i++ {
			for j := i + 1; j < len(oList); j++ {
				// If two orders on same equipment within 72 hours
				diff := oList[i].CreatedAt.Sub(oList[j].CreatedAt)
				if diff > 0 && diff < 72*time.Hour {
					// Check if descriptions are similar (simple keyword check)
					if a.isSimilar(oList[i].FaultDescription, oList[j].FaultDescription) {
						data.Anomalies = append(data.Anomalies, dto.AnomalyItem{
							AnomalyType:     "repeat_failure",
							Severity:        "high",
							TargetType:      "repair_order",
							TargetID:        oList[i].ID,
							Title:           "发现短期重复故障",
							Description:     fmt.Sprintf("设备 ID:%d 在 72 小时内连续报修：\"%s\" 和 \"%s\"。", eqID, oList[j].FaultDescription, oList[i].FaultDescription),
							SuggestedAction: "建议核查前次维修是否彻底，或是否存在误诊断。请参考知识库标准流程。",
						})
						statsMap["high_risk_count"] = statsMap["high_risk_count"].(int) + 1
					}
				}
			}
		}
	}

	// 2. Identify high cost anomalies
	for _, o := range orders {
		if o.CostDetail != nil {
			totalCost := o.CostDetail.SparePartCost + o.CostDetail.LaborCost
			if totalCost > 5000 { // Threshold for anomaly
				data.Anomalies = append(data.Anomalies, dto.AnomalyItem{
					AnomalyType:     "high_cost",
					Severity:        "medium",
					TargetType:      "repair_order",
					TargetID:        o.ID,
					Title:           "高额维修成本预警",
					Description:     fmt.Sprintf("维修单 #%d 产生总费用 %.2f，超过预算阈值。", o.ID, totalCost),
					SuggestedAction: "建议核实更换备件的必要性及备件领用记录。",
				})
			}
		}
	}

	// 3. Fetch evidence for the first anomaly if found
	if len(data.Anomalies) > 0 && req.EquipmentTypeID != 0 {
		evidence, err := a.retrievalTool.SearchManualKnowledge(data.Anomalies[0].Description, &req.EquipmentTypeID, user)
		if err == nil {
			data.Evidence = evidence
		}
	}

	statsMap["anomaly_summary"] = fmt.Sprintf("在 %d 个工单中发现 %d 个异常点。", total, len(data.Anomalies))

	return data, nil
}

func (a *RepairAuditAnalyzer) isSimilar(d1, d2 string) bool {
	// Simple keyword overlap check
	words := []string{"异响", "漏油", "报错", "停机", "温度", "压力", "轴承", "电机", "皮带"}
	matchCount := 0
	for _, w := range words {
		if strings.Contains(d1, w) && strings.Contains(d2, w) {
			matchCount++
		}
	}
	return matchCount >= 1
}

