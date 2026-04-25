package prompt

import (
	"fmt"
)

type PromptTool struct {
}

func NewPromptTool() *PromptTool {
	return &PromptTool{}
}

func (t *PromptTool) BuildMaintenanceRecommendPrompt(data interface{}, evidence interface{}) string {
	return fmt.Sprintf(`你是一个专业的工业设备管理助手。请根据以下设备当前的保养计划和相关的参考证据（手册或最佳实践），为工程师生成一份中文保养优化建议。

### 当前计划
%v

### 参考证据
%v

### 输出要求
1. 语言：中文
2. 风格：专业、严谨、客观
3. 重点：评估当前周期的合理性，是否需要增加或删除维护项，并给出理由。
4. 格式：简短的摘要（50-100字），随后是具体的建议项。`, data, evidence)
}

func (t *PromptTool) BuildRepairAuditPrompt(data interface{}, evidence interface{}) string {
	return fmt.Sprintf(`你是一个设备维修审计助手。请根据以下维修记录异常分析和相关的标准证据，生成一份中文审计报告结论。

### 异常分析结果
%v

### 参考标准/知识
%v

### 输出要求
1. 语言：中文
2. 重点：指出风险点（如重复故障、费用异常），解释为什么这被认为是异常，并给出核查建议。
3. 风格：批判性思维但保持专业。`, data, evidence)
}
