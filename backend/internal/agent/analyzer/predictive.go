package analyzer

import (
	"fmt"
	"math"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/tool"
)

type PredictiveAnalyzer struct {
	repairTool      *tool.RepairTool
	maintenanceTool *tool.MaintenanceTool
	retrievalTool   *tool.RetrievalTool
}

func NewPredictiveAnalyzer(rt *tool.RepairTool, mt *tool.MaintenanceTool, ret *tool.RetrievalTool) *PredictiveAnalyzer {
	return &PredictiveAnalyzer{
		repairTool:      rt,
		maintenanceTool: mt,
		retrievalTool:   ret,
	}
}

// PredictRUL 预测剩余健康寿命
func (a *PredictiveAnalyzer) PredictRUL(equipmentID uint) (*dto.RULPrediction, error) {
	// 1. 获取基础统计 (MTBF 等)
	stats, err := a.repairTool.GetFailureStats(equipmentID)
	if err != nil { return nil, err }
	
	// 2. 获取最近 30 天工况 (从证据链逻辑中简化)
	// 在生产中这会调用 get_equipment_runtime 工具
	loadFactor := 1.15 // 默认使用 Demo 中的超负荷系数
	avgMTBFHours := 300.0 // 假设该型号标准 MTBF 为 300 小时
	
	repairCount := 0.0
	if val, ok := stats["repair_count"]; ok {
		repairCount = float64(val.(int64))
	}
	
	// 如果故障频繁，动态下调预期 MTBF
	if repairCount > 3 {
		avgMTBFHours = avgMTBFHours * 0.7
	}

	// 3. 计算预计剩余时间
	// 简化模型：RUL_hours = Max(0, MTBF - (TimeSinceLastRepair * Load))
	currentUsedHours := 240.0 // 模拟已连续运行时间
	rulHours := (avgMTBFHours - currentUsedHours) / loadFactor
	if rulHours < 0 { rulHours = 0 }
	
	rulDays := int(math.Ceil(rulHours / 24.0))
	
	// 4. 计算健康分 (0-100)
	healthScore := (rulHours / avgMTBFHours) * 100
	if healthScore > 100 { healthScore = 100 }

	prediction := &dto.RULPrediction{
		EquipmentID:      equipmentID,
		HealthScore:      healthScore,
		EstimatedRULDays: rulDays,
		Reliability:      0.85,
		RiskFactors:      []string{},
	}

	// 5. 风险因子识别
	if loadFactor > 1.1 {
		prediction.RiskFactors = append(prediction.RiskFactors, "超负荷运行 (115%)")
	}
	if repairCount > 2 {
		prediction.RiskFactors = append(prediction.RiskFactors, "近期发生重复故障，硬件疲劳加速")
	}

	// 6. 给出建议
	if rulDays <= 3 {
		prediction.Recommendation = "紧急建议：设备已进入高风险故障期，预计 3 天内可能发生停机，请立即安排预防性检修。"
	} else if rulDays <= 7 {
		prediction.Recommendation = "预警：健康分较低，建议在下周内安排油质监测和同心度校准。"
	} else {
		prediction.Recommendation = "状态良好：建议继续保持当前的预防性保养节奏。"
	}

	return prediction, nil
}

// SymptomFinding 代表一个识别出的故障征兆
type SymptomFinding struct {
	Type        string   `json:"type"` // micro_stop, pm_failure, mttr_drift
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	Evidence    []string `json:"evidence"`
}

// DetectSymptoms 识别设备亚健康征兆
func (a *PredictiveAnalyzer) DetectSymptoms(equipmentID uint) ([]SymptomFinding, error) {
	findings := []SymptomFinding{}
	
	// 1. 获取最近 30 天的维修记录
	orders, err := a.repairTool.GetRecentOrdersByEquipment(equipmentID, 20)
	if err != nil { return nil, err }

	// 2. 分析“频发微停” (Micro-stops)
	shortRepairs := 0
	for _, o := range orders {
		if o.StartedAt != nil && o.CompletedAt != nil {
			duration := o.CompletedAt.Sub(*o.StartedAt).Minutes()
			if duration < 30 { shortRepairs++ }
		}
	}
	if shortRepairs >= 3 {
		findings = append(findings, SymptomFinding{
			Type: "micro_stop", Title: "频发微停预警", Severity: "medium",
			Description: fmt.Sprintf("设备在近期出现 %d 次短时停机，这通常是核心部件即将彻底失效的前兆。", shortRepairs),
			Evidence: []string{"近期 30 分钟以内的紧急处置记录超过 3 次"},
		})
	}

	// 3. 分析“保养无效性” (PM vs CM)
	// 检查最近一次保养后 48 小时内是否有报修
	// (此处为 Demo 逻辑简化)
	findings = append(findings, SymptomFinding{
		Type: "pm_ineffective", Title: "保养有效性质疑", Severity: "high",
		Description: "检测到设备在执行『二级保养』后 72 小时内即发生了液压系统报警。",
		Evidence: []string{"保养单 ID: 202, 维修单 ID: 105"},
	})

	return findings, nil
}

// TCOResult 代表设备的总持有成本分析结果
type TCOResult struct {
	EquipmentID      uint    `json:"equipment_id"`
	AccumulatedRepair float64 `json:"accumulated_repair_cost"`  // 累计维修费
	DowntimeLoss      float64 `json:"downtime_loss"`           // 累计停机损失
	DepreciatedValue  float64 `json:"depreciated_value"`       // 当前账面净值
	TCO               float64 `json:"total_cost_of_ownership"` // TCO 总额
	MaintenanceRatio  float64 `json:"maintenance_to_asset_ratio"` // 维护原值比
}

// CalculateTCO 计算设备全生命周期总成本
func (a *PredictiveAnalyzer) CalculateTCO(equipmentID uint) (*TCOResult, error) {
	// 1. 获取设备财务档案
	profile, err := a.retrievalTool.GetEquipmentProfile(equipmentID)
	if err != nil { return nil, err }
	
	purchasePrice := profile["purchase_price"].(float64)
	scrapValue := profile["scrap_value"].(float64)
	serviceLife := float64(profile["service_life_years"].(int))
	hourlyLoss := profile["hourly_loss"].(float64)
	
	// 2. 获取累计维修费与停机时长
	costStats, _ := a.repairTool.GetCostAnalysis(equipmentID)
	failureStats, _ := a.repairTool.GetFailureStats(equipmentID)
	
	repairCost := costStats["total_cost"].(float64)
	downtimeHours := failureStats["total_downtime"].(float64)
	
	// 3. 计算逻辑
	// 停机损失 = 累计停机小时 * 产值损失单价
	downtimeLoss := downtimeHours * hourlyLoss
	
	// 计算折旧 (直线折旧法)
	// 假设已使用年限 (从故事线推断: CNC 3年, PRESS 12年)
	yearsUsed := 3.0
	if profile["code"] == "PRESS-05" { yearsUsed = 12.0 }
	
	annualDepreciation := (purchasePrice - scrapValue) / serviceLife
	currentNetValue := purchasePrice - (annualDepreciation * yearsUsed)
	if currentNetValue < scrapValue { currentNetValue = scrapValue }
	
	tco := repairCost + downtimeLoss + (purchasePrice - currentNetValue)
	
	return &TCOResult{
		EquipmentID:      equipmentID,
		AccumulatedRepair: repairCost,
		DowntimeLoss:      downtimeLoss,
		DepreciatedValue:  currentNetValue,
		TCO:               tco,
		MaintenanceRatio:  repairCost / purchasePrice,
	}, nil
}

// EvaluateRetirement 评估设备是否建议退役
func (a *PredictiveAnalyzer) EvaluateRetirement(equipmentID uint) (map[string]interface{}, error) {
	tco, err := a.CalculateTCO(equipmentID)
	if err != nil { return nil, err }
	
	decision := "continue"
	reason := "设备维护成本处于合理区间，运行ROI良好。"
	
	// 触发退役的科学条件：
	// 1. 累计维修费 > 资产原值的 60%
	// 2. 或是维护费增长速率过快 (在Demo中简化为比值判断)
	if tco.MaintenanceRatio > 0.6 {
		decision = "retire"
		reason = fmt.Sprintf("强烈建议退役。累计维修费已达原值的 %.1f%%，继续持有将产生负向财务现金流。", tco.MaintenanceRatio*100)
	} else if tco.MaintenanceRatio > 0.4 {
		decision = "evaluate"
		reason = "建议列入观察名单。维护成本偏高，建议对比新一代机型的能效比后再做决策。"
	}
	
	return map[string]interface{}{
		"decision": decision,
		"reason":   reason,
		"tco_data": tco,
	}, nil
}
