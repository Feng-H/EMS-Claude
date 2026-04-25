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

func (t *PromptTool) BuildKnowledgeExtractionPrompt(history interface{}) string {
	return fmt.Sprintf(`你是一个资深的工业设备知识专家。请仔细阅读下面这段工程师与 AI 助手的对话记录，判断其中是否包含有价值的设备管理知识（如故障根因、预防措施、操作经验等）。

### 对话记录
%v

### 提取任务
如果包含有价值的结论，请将其提取为以下格式的 JSON（注意：如果不值得提取，请只返回 {}）：
{
  "title": "简短的知识标题",
  "type": "root_cause_analysis 或 pattern 或 equipment_profile 等",
  "summary": "一句话核心结论",
  "details": {
    "evidence": ["证据1", "证据2"],
    "root_cause": "根本原因说明",
    "solution": "解决建议",
    "prevention": "预防措施"
  },
  "confidence": 0.0到1.0的置信度分数
}

### 要求
1. 只返回纯 JSON 字符串，不要包含任何 Markdown 代码块包裹（即不要有 ` + "```" + `json 等）。
2. 确保结论是基于对话事实提取的。
3. 语言必须是中文。`, history)
}

func (t *PromptTool) BuildSkillExtractionPrompt(history interface{}) string {
	return fmt.Sprintf(`你是一个高级工业诊断专家。请分析以下工程师与 AI 助手的对话记录，看其中是否隐藏了一套通用的“故障排查或数据分析逻辑”。

### 对话记录
%v

### 任务
如果工程师引导你完成了一次成功的、具有代表性的深度排查，请将这套排查套路提炼为一个“技能草稿” JSON。
要求：
{
  "name": "技能名称（如：液压泵内泄排查）",
  "description": "简述该技能解决什么问题",
  "applicable_scenarios": ["场景1", "场景2"],
  "steps": [
    { "step": 1, "action": "具体动作描述", "tool": "建议使用的工具名" },
    { "step": 2, "action": "...", "tool": "..." }
  ]
}

### 可用工具参考
- get_failure_distribution: 统计故障分布
- search_manual_knowledge: 搜索手册与知识库
- get_maintenance_profile: 获取保养计划与执行情况
- get_equipment_runtime: 获取运行快照

### 要求
1. 只返回纯 JSON 字符串，不要包含任何 Markdown 代码块包裹。
2. 提炼的步骤应具有通用性，不局限于本次对话的具体设备。
3. 如果不值得提炼，只返回 {}。`, history)
}
