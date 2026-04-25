package service

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/ems/backend/internal/agent/analyzer"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/policy"
	"github.com/ems/backend/internal/agent/prompt"
	"github.com/ems/backend/internal/agent/repository"
	"github.com/ems/backend/internal/agent/tool"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/database"
	"github.com/ems/backend/pkg/llm"
	"github.com/ems/backend/pkg/trace"
)

type AgentService struct {
	repo   repository.IAgentRepository
	policy *policy.PolicyService
	
	// Tools
	retrievalTool   *tool.RetrievalTool
	maintenanceTool *tool.MaintenanceTool
	repairTool      *tool.RepairTool
	promptTool      *prompt.PromptTool
	
	// LLM
	llmClient llm.LLMClient
	
	// Analyzers
	maintenanceAnalyzer *analyzer.MaintenanceAnalyzer
	repairAuditAnalyzer *analyzer.RepairAuditAnalyzer
	predictiveAnalyzer  *analyzer.PredictiveAnalyzer
}

func NewAgentService() *AgentService {
	var repo repository.IAgentRepository
	if config.Cfg.Storage.Mode == "memory" {
		repo = repository.NewMemoryAgentRepository()
	} else {
		repo = repository.NewDBAgentRepository(database.GetDB())
	}
	
	retrievalTool := tool.NewRetrievalTool(repo)
	maintenanceTool := tool.NewMaintenanceTool()
	repairTool := tool.NewRepairTool()
	
	var llmClient llm.LLMClient
	if config.Cfg.LLM.APIKey != "" {
		llmClient = llm.NewOpenAIClient(config.Cfg.LLM.BaseURL, config.Cfg.LLM.APIKey, config.Cfg.LLM.Model)
	}
	
	return &AgentService{
		repo:   repo,
		policy: policy.NewPolicyService(),
		retrievalTool:   retrievalTool,
		maintenanceTool: maintenanceTool,
		repairTool:      repairTool,
		promptTool:      prompt.NewPromptTool(),
		llmClient:       llmClient,
		maintenanceAnalyzer: analyzer.NewMaintenanceAnalyzer(retrievalTool, maintenanceTool),
		repairAuditAnalyzer: analyzer.NewRepairAuditAnalyzer(retrievalTool, repairTool),
		predictiveAnalyzer:  analyzer.NewPredictiveAnalyzer(repairTool, maintenanceTool, retrievalTool),
	}
}

func (s *AgentService) RecommendMaintenance(userID uint, role string, req *dto.MaintenanceRecommendRequest) (*dto.AgentResponseEnvelope, error) {
	startTime := time.Now()
	traceID := trace.GenerateTraceID()
	
	agentCtx, err := s.policy.DeriveAgentContext(userID, role, req.Language)
	if err != nil { return nil, err }

	targetFactoryID := req.FactoryID
	if targetFactoryID == 0 && agentCtx.FactoryID != nil {
		targetFactoryID = *agentCtx.FactoryID
	}
	
	if err := s.policy.ValidateScope(agentCtx, &targetFactoryID); err != nil {
		return nil, err
	}

	analysisResult, err := s.maintenanceAnalyzer.Analyze(req)
	if err != nil { return nil, err }

	summary := "建议缩短保养周期，以提高设备可用性。"
	if s.llmClient != nil {
		p := req.SystemPrompt
		if p == "" {
			p = s.promptTool.BuildMaintenanceRecommendPrompt(analysisResult.CurrentPlan, analysisResult.Evidence)
		} else {
			p = fmt.Sprintf("%s\n\n### 原始数据参考\n当前计划: %v\n参考证据: %v", p, analysisResult.CurrentPlan, analysisResult.Evidence)
		}
		resp, _ := s.llmClient.ChatCompletion([]llm.Message{
			{Role: "system", Content: "你是一个专业的工业设备管理助手。"},
			{Role: "user", Content: p},
		})
		if resp != "" { summary = resp }
	} else if len(analysisResult.Recommendations) > 0 {
		summary = analysisResult.Recommendations[0].Description + "。" + analysisResult.Recommendations[0].Reason
	}

	inputSnap, _ := json.Marshal(req)
	resultJSON, _ := json.Marshal(analysisResult)
	session := &model.AgentSession{
		UserID: userID, Scenario: "maintenance_recommendation", FactoryID: &targetFactoryID,
		Language: agentCtx.Language, InputSnapshot: string(inputSnap), TraceID: traceID, Status: "completed",
	}
	_ = s.repo.CreateSession(session)

	artifact := &model.AgentArtifact{
		SessionID: session.ID, ArtifactType: "recommendation", Title: "设备保养优化建议",
		Summary: summary, ResultJSON: string(resultJSON), RiskLevel: "medium",
	}
	_ = s.repo.CreateArtifact(artifact)

	for _, ev := range analysisResult.Evidence {
		link := model.AgentEvidenceLink{
			ArtifactID: artifact.ID, EvidenceType: ev.EvidenceType,
			SourceTable: ev.SourceTable, SourceID: ev.SourceID, Excerpt: ev.Excerpt, Score: ev.Score,
		}
		_ = s.repo.CreateEvidenceLinks([]model.AgentEvidenceLink{link})
	}

	res := &dto.AgentResponseEnvelope{
		Success: true, TraceID: traceID, Language: agentCtx.Language, Scenario: "maintenance_recommendation",
		ScopeSummary: map[string]interface{}{"factory_id": targetFactoryID},
		Summary: summary, RiskLevel: "medium", ArtifactID: artifact.ID,
		EvidenceCount: len(analysisResult.Evidence), Data: analysisResult,
	}
	s.logUsage(session.ID, userID, "maintenance_recommendation", startTime)
	return res, nil
}

func (s *AgentService) AuditRepair(userID uint, role string, req *dto.RepairAuditRequest) (*dto.AgentResponseEnvelope, error) {
	startTime := time.Now()
	traceID := trace.GenerateTraceID()
	
	agentCtx, err := s.policy.DeriveAgentContext(userID, role, req.Language)
	if err != nil { return nil, err }

	targetFactoryID := req.FactoryID
	if targetFactoryID == 0 && agentCtx.FactoryID != nil {
		targetFactoryID = *agentCtx.FactoryID
	}
	
	if err := s.policy.ValidateScope(agentCtx, &targetFactoryID); err != nil {
		return nil, err
	}

	analysisResult, err := s.repairAuditAnalyzer.Analyze(req)
	if err != nil { return nil, err }

	summary := "发现维修异常，建议复核维修质量。"
	if s.llmClient != nil {
		p := req.SystemPrompt
		if p == "" {
			p = s.promptTool.BuildRepairAuditPrompt(analysisResult.Anomalies, analysisResult.Evidence)
		} else {
			p = fmt.Sprintf("%s\n\n### 原始数据参考\n异常项: %v\n参考证据: %v", p, analysisResult.Anomalies, analysisResult.Evidence)
		}
		resp, _ := s.llmClient.ChatCompletion([]llm.Message{
			{Role: "system", Content: "你是一个设备维修审计助手。"},
			{Role: "user", Content: p},
		})
		if resp != "" { summary = resp }
	} else if stats, ok := analysisResult.Stats.(map[string]interface{}); ok {
		if val, exists := stats["anomaly_summary"]; exists {
			summary = val.(string)
		}
	}

	inputSnap, _ := json.Marshal(req)
	resultJSON, _ := json.Marshal(analysisResult)
	session := &model.AgentSession{
		UserID: userID, Scenario: "repair_audit", FactoryID: &targetFactoryID,
		Language: agentCtx.Language, InputSnapshot: string(inputSnap), TraceID: traceID, Status: "completed",
	}
	_ = s.repo.CreateSession(session)

	artifact := &model.AgentArtifact{
		SessionID: session.ID, ArtifactType: "audit_report", Title: "设备维修审计报告",
		Summary: summary, ResultJSON: string(resultJSON), RiskLevel: "high",
	}
	_ = s.repo.CreateArtifact(artifact)

	for _, ev := range analysisResult.Evidence {
		link := model.AgentEvidenceLink{
			ArtifactID: artifact.ID, EvidenceType: ev.EvidenceType,
			SourceTable: ev.SourceTable, SourceID: ev.SourceID, Excerpt: ev.Excerpt, Score: ev.Score,
		}
		_ = s.repo.CreateEvidenceLinks([]model.AgentEvidenceLink{link})
	}

	res := &dto.AgentResponseEnvelope{
		Success: true, TraceID: traceID, Language: agentCtx.Language, Scenario: "repair_audit",
		ScopeSummary: map[string]interface{}{"factory_id": targetFactoryID},
		Summary: summary, RiskLevel: "high", ArtifactID: artifact.ID,
		EvidenceCount: len(analysisResult.Evidence), Data: analysisResult,
	}
	s.logUsage(session.ID, userID, "repair_audit", startTime)
	return res, nil
}

func (s *AgentService) AuditMaintenance(userID uint, role string, req *dto.MaintenanceAuditRequest) (*dto.AgentResponseEnvelope, error) {
	traceID := trace.GenerateTraceID()
	return &dto.AgentResponseEnvelope{
		Success: true, TraceID: traceID, Language: "zh-CN", Scenario: "maintenance_audit",
		Summary: "这是保养审计的占位符响应（开发中）", Data: dto.MaintenanceAuditData{},
	}, nil
}

func (s *AgentService) Analyze(userID uint, role string, req *dto.AnalyzeRequest) (*dto.AgentResponseEnvelope, error) {
	traceID := trace.GenerateTraceID()
	return &dto.AgentResponseEnvelope{
		Success: true, TraceID: traceID, Language: "zh-CN", Scenario: "analysis",
		Summary: "这是分析助手的占位符响应（开发中）", Data: dto.AnalyzeData{},
	}, nil
}

func (s *AgentService) ListSessions(userID uint, limit int) ([]dto.AgentSessionResponse, error) {
	sessions, err := s.repo.ListSessionsByUserID(userID, limit)
	if err != nil { return nil, err }
	results := make([]dto.AgentSessionResponse, len(sessions))
	for i, s := range sessions {
		results[i] = dto.AgentSessionResponse{
			ID: s.ID, UserID: s.UserID, Scenario: s.Scenario, Language: s.Language,
			Status: s.Status, TraceID: s.TraceID, CreatedAt: s.CreatedAt,
		}
		if s.FactoryID != nil { results[i].FactoryID = *s.FactoryID }
	}
	return results, nil
}

func (s *AgentService) GetSession(id uint) (*dto.AgentSessionResponse, error) {
	session, err := s.repo.GetSessionByID(id)
	if err != nil { return nil, err }
	res := &dto.AgentSessionResponse{
		ID: session.ID, UserID: session.UserID, Scenario: session.Scenario, Language: session.Language,
		Status: session.Status, TraceID: session.TraceID, CreatedAt: session.CreatedAt,
	}
	if session.FactoryID != nil { res.FactoryID = *session.FactoryID }
	for _, a := range session.Artifacts { res.Artifacts = append(res.Artifacts, a.ID) }
	return res, nil
}

func (s *AgentService) GetArtifact(id uint) (*dto.AgentArtifactResponse, error) {
	artifact, err := s.repo.GetArtifactByID(id)
	if err != nil { return nil, err }
	res := &dto.AgentArtifactResponse{
		ID: artifact.ID, SessionID: artifact.SessionID, ArtifactType: artifact.ArtifactType,
		Title: artifact.Title, Summary: artifact.Summary, RiskLevel: artifact.RiskLevel, CreatedAt: artifact.CreatedAt,
	}
	_ = json.Unmarshal([]byte(artifact.ResultJSON), &res.ResultJSON)
	for _, ev := range artifact.EvidenceLinks {
		res.Evidence = append(res.Evidence, dto.EvidenceItem{
			EvidenceType: ev.EvidenceType, SourceTable: ev.SourceTable, SourceID: ev.SourceID, Excerpt: ev.Excerpt, Score: ev.Score,
		})
	}
	return res, nil
}

func (s *AgentService) AuditKnowledge(id string, status string, verifierID uint) error {
	return s.repo.UpdateKnowledgeStatus(id, status, &verifierID)
}

// =====================================================
// Phase 2: Chat & Conversational Logic
// =====================================================

func (s *AgentService) Chat(userID uint, role string, req *dto.ChatRequest) (*dto.ChatResponse, error) {
	startTime := time.Now()
	traceID := trace.GenerateTraceID()

	// 1. 获取或创建对话
	var convID = req.ConversationID
	if convID == 0 {
		newConv := &model.AgentConversation{
			UserID: userID,
			Title:  req.Message,
		}
		if len(newConv.Title) > 50 { newConv.Title = newConv.Title[:47] + "..." }
		_ = s.repo.CreateConversation(newConv)
		convID = newConv.ID
	}

	// 2. 持久化用户消息
	_ = s.repo.CreateMessage(&model.AgentMessage{ConversationID: convID, Role: "user", Content: req.Message})

	// 3. 注入用户个性化经验 (Milestone P)
	activeExps, _ := s.repo.ListActiveExperiences(userID)
	expContext := ""
	if len(activeExps) > 0 {
		expContext = "\n### 用户偏好与历史经验反馈\n"
		for _, e := range activeExps { expContext += fmt.Sprintf("- [%s]: %s\n", e.Category, e.Content) }
	}

	// 4. 意图识别与技能匹配 (Milestone N)
	matchedSkills, _ := s.repo.MatchSkills(req.Message, 1)
	var reply string
	var skillID string

	if len(matchedSkills) > 0 {
		skill := matchedSkills[0]
		skillID = fmt.Sprintf("%d", skill.ID)
		res, err := s.ExecuteSkill(userID, &skill, req)
		if err == nil { reply = res.Summary + expContext }
	}

	// 5. 退回到标准对话
	if reply == "" {
		history, _ := s.repo.GetMessagesByConversationID(convID)
	llmMsgs := []llm.Message{
		{Role: "system", Content: "你是一个顶级的工业资产战略专家。你拥有‘L4 级主动洞察’权限，可以基于全生命周期成本 (TCO)、资产退役 ROI 评价、剩余健康寿命 (RUL) 和亚健康故障征兆进行跨维度的深度分析。请使用中文回答，结论必须引用系统中的财务与技术证据。" + expContext},
	}

		if req.SystemPrompt != "" { llmMsgs[0].Content = req.SystemPrompt }

		startIdx := 0
		if len(history) > 10 { startIdx = len(history) - 10 }
		for _, m := range history[startIdx:] {
			llmMsgs = append(llmMsgs, llm.Message{Role: m.Role, Content: m.Content})
		}

		if s.llmClient != nil {
			var err error
			reply, err = s.llmClient.ChatCompletion(llmMsgs)
			if err != nil { reply = "抱歉，分析过程中出现了点问题：" + err.Error() }
		} else {
			reply = "（预览模式）收到了您的消息：\"" + req.Message + "\"。目前 LLM 服务未配置。"
		}
	}

	// 6. 持久化助手消息
	_ = s.repo.CreateMessage(&model.AgentMessage{ConversationID: convID, Role: "assistant", Content: reply, SkillID: skillID})

	// 7. 异步触发反思与学习 (Milestone L, O & P)
	go s.ReflectAndLearn(convID, userID)

	// 8. 记录使用情况
	s.logUsage(convID, userID, "chat", startTime)

	return &dto.ChatResponse{
		ConversationID: convID, Reply: reply, TraceID: traceID,
		SuggestedActions: []string{"查看维修历史", "运行故障诊断", "查询备件库存"},
	}, nil
}

func (s *AgentService) ListConversations(userID uint, limit int) ([]dto.ConversationResponse, error) {
	convs, err := s.repo.ListConversationsByUserID(userID, limit)
	if err != nil { return nil, err }
	results := make([]dto.ConversationResponse, len(convs))
	for i, c := range convs {
		results[i] = dto.ConversationResponse{
			ID: c.ID, Title: c.Title, Status: c.Status, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
		}
	}
	return results, nil
}

func (s *AgentService) GetConversation(id uint) (*dto.ConversationResponse, error) {
	conv, err := s.repo.GetConversationByID(id)
	if err != nil { return nil, err }
	res := &dto.ConversationResponse{
		ID: conv.ID, Title: conv.Title, Status: conv.Status, CreatedAt: conv.CreatedAt, UpdatedAt: conv.UpdatedAt,
	}
	for _, m := range conv.Messages {
		res.Messages = append(res.Messages, dto.MessageItem{
			ID: m.ID, Role: m.Role, Content: m.Content, CreatedAt: m.CreatedAt,
		})
	}
	return res, nil
}

// =====================================================
// Phase 2: Skill Management
// =====================================================

func (s *AgentService) CreateSkill(req *dto.CreateSkillRequest) (*dto.SkillResponse, error) {
	appTo, _ := json.Marshal(req.ApplicableTo)
	appSce, _ := json.Marshal(req.ApplicableScenarios)
	steps, _ := json.Marshal(req.Steps)
	skill := &model.AgentSkill{
		Name: req.Name, Description: req.Description, ApplicableTo: string(appTo),
		ApplicableScenarios: string(appSce), Steps: string(steps), Status: "draft",
	}
	if err := s.repo.CreateSkill(skill); err != nil { return nil, err }
	return s.mapSkillToResponse(skill), nil
}

func (s *AgentService) ListSkills(status string) ([]dto.SkillResponse, error) {
	skills, err := s.repo.ListSkills(status, 100)
	if err != nil { return nil, err }
	results := make([]dto.SkillResponse, len(skills))
	for i, sk := range skills { results[i] = *s.mapSkillToResponse(&sk) }
	return results, nil
}

func (s *AgentService) ExecuteSkill(userID uint, skill *model.AgentSkill, req *dto.ChatRequest) (*dto.AgentResponseEnvelope, error) {
	var steps []struct {
		Step   int    `json:"step"`
		Action string `json:"action"`
		Tool   string `json:"tool"`
	}
	_ = json.Unmarshal([]byte(skill.Steps), &steps)
	evidence := []dto.EvidenceItem{}
	for _, step := range steps {
		switch step.Tool {
		case "get_equipment_profile":
			if res, err := s.retrievalTool.GetEquipmentProfile(1); err == nil {
				resJSON, _ := json.Marshal(res)
				evidence = append(evidence, dto.EvidenceItem{
					EvidenceType: "equipment_profile", Title: "设备基础信息", Excerpt: string(resJSON), Score: 1.0,
				})
			}
		case "get_failure_stats":
			if res, err := s.repairTool.GetFailureStats(1); err == nil {
				resJSON, _ := json.Marshal(res)
				evidence = append(evidence, dto.EvidenceItem{
					EvidenceType: "failure_stats", Title: "故障统计分析", Excerpt: string(resJSON), Score: 0.95,
				})
			}
		case "get_cost_analysis":
			if res, err := s.repairTool.GetCostAnalysis(1); err == nil {
				resJSON, _ := json.Marshal(res)
				evidence = append(evidence, dto.EvidenceItem{
					EvidenceType: "cost_analysis", Title: "维修成本分析", Excerpt: string(resJSON), Score: 0.9,
				})
			}
		case "get_maintenance_compliance":
			if res, err := s.maintenanceTool.GetMaintenanceCompliance(1); err == nil {
				resJSON, _ := json.Marshal(res)
				evidence = append(evidence, dto.EvidenceItem{
					EvidenceType: "maintenance_compliance", Title: "保养合规性评估", Excerpt: string(resJSON), Score: 0.85,
				})
			}
		case "get_failure_distribution":
			auditReq := &dto.RepairAuditRequest{EquipmentTypeID: 12}
			res, err := s.repairAuditAnalyzer.Analyze(auditReq)
			if err == nil { evidence = append(evidence, res.Evidence...) }
		case "search_manual_knowledge":
			res, err := s.retrievalTool.SearchManualKnowledge(req.Message, nil)
			if err == nil {
				evidence = append(evidence, res...)
			}
		case "predict_remaining_life":
			// 默认针对 ID 1 进行预测 (Demo 逻辑)
			if res, err := s.predictiveAnalyzer.PredictRUL(1); err == nil {
				resJSON, _ := json.Marshal(res)
				evidence = append(evidence, dto.EvidenceItem{
					EvidenceType: "prediction", Title: "RUL 剩余健康寿命预测",
					Excerpt: string(resJSON), Score: 0.98,
				})
				}
				case "detect_symptoms":
					if res, err := s.predictiveAnalyzer.DetectSymptoms(1); err == nil {
						resJSON, _ := json.Marshal(res)
						evidence = append(evidence, dto.EvidenceItem{
							EvidenceType: "symptoms", Title: "设备亚健康征兆识别",
							Excerpt: string(resJSON), Score: 0.92,
						})
					}
				case "get_tco_analysis":
					if res, err := s.predictiveAnalyzer.CalculateTCO(1); err == nil {
						resJSON, _ := json.Marshal(res)
						evidence = append(evidence, dto.EvidenceItem{
							EvidenceType: "tco", Title: "资产总持有成本(TCO)分析",
							Excerpt: string(resJSON), Score: 0.95,
						})
					}
				case "get_retirement_recommendation":
					if res, err := s.predictiveAnalyzer.EvaluateRetirement(1); err == nil {
						resJSON, _ := json.Marshal(res)
						evidence = append(evidence, dto.EvidenceItem{
							EvidenceType: "retirement", Title: "资产退役与投资决策建议",
							Excerpt: string(resJSON), Score: 0.99,
						})
					}
				}
				}

	summary := fmt.Sprintf("执行技能【%s】: 已完成 %d 个分析步骤，收集到 %d 条证据。", skill.Name, len(steps), len(evidence))
	if s.llmClient != nil {
		prompt := fmt.Sprintf("用户意图: %s\n执行技能: %s\n技能描述: %s\n收集到的证据链: %v\n\n请根据以上信息给出一份专业的中文分析摘要。", req.Message, skill.Name, skill.Description, evidence)
		resp, _ := s.llmClient.ChatCompletion([]llm.Message{
			{Role: "system", Content: "你是一个专业的工业设备管理助手，正在执行预定义的分析技能。"},
			{Role: "user", Content: prompt},
		})
		if resp != "" { summary = resp }
	}
	return &dto.AgentResponseEnvelope{
		Success: true, Scenario: "skill_execution", Summary: summary, EvidenceCount: len(evidence),
		Data: map[string]interface{}{"skill_id": skill.ID, "skill_name": skill.Name, "executed_steps": len(steps), "evidence": evidence},
	}, nil
}

func (s *AgentService) GetSkill(id uint) (*dto.SkillResponse, error) {
	skill := s.repo.GetSkillByID(id)
	if skill == nil { return nil, fmt.Errorf("skill not found") }
	return s.mapSkillToResponse(skill), nil
}

func (s *AgentService) UpdateSkill(id uint, req *dto.UpdateSkillRequest) (*dto.SkillResponse, error) {
	skill := s.repo.GetSkillByID(id)
	if skill == nil { return nil, fmt.Errorf("skill not found") }
	if req.Name != "" { skill.Name = req.Name }
	if req.Description != "" { skill.Description = req.Description }
	if req.Status != "" { skill.Status = req.Status }
	if req.ApplicableTo != nil {
		appTo, _ := json.Marshal(req.ApplicableTo)
		skill.ApplicableTo = string(appTo)
	}
	if req.ApplicableScenarios != nil {
		appSce, _ := json.Marshal(req.ApplicableScenarios)
		skill.ApplicableScenarios = string(appSce)
	}
	if req.Steps != nil {
		steps, _ := json.Marshal(req.Steps)
		skill.Steps = string(steps)
	}
	if err := s.repo.UpdateSkill(skill); err != nil { return nil, err }
	return s.mapSkillToResponse(skill), nil
}

// =====================================================
// Phase 2: Proactive Notification Logic
// =====================================================

func (s *AgentService) Subscribe(userID uint, pushType string, enabled bool, scope any) error {
	scopeJSON, _ := json.Marshal(scope)
	sub := &model.AgentPushSubscription{
		UserID: userID, PushType: pushType, Enabled: enabled, Scope: string(scopeJSON),
	}
	return s.repo.CreatePushSubscription(sub)
}

func (s *AgentService) NotifyEvent(eventType string, targetID uint, context map[string]interface{}) { }

func (s *AgentService) mapSkillToResponse(sk *model.AgentSkill) *dto.SkillResponse {
	var appTo, appSce []string
	var steps []any
	_ = json.Unmarshal([]byte(sk.ApplicableTo), &appTo)
	_ = json.Unmarshal([]byte(sk.ApplicableScenarios), &appSce)
	_ = json.Unmarshal([]byte(sk.Steps), &steps)
	return &dto.SkillResponse{
		ID: sk.ID, Name: sk.Name, Description: sk.Description, ApplicableTo: appTo, ApplicableScenarios: appSce,
		Steps: steps, Version: sk.Version, Status: sk.Status, UsageCount: sk.UsageCount, SuccessRate: sk.SuccessRate, CreatedAt: sk.CreatedAt,
	}
}

func (s *AgentService) logUsage(sessionID, userID uint, scenario string, startTime time.Time) {
	duration := time.Since(startTime).Milliseconds()
	modelName := config.Cfg.LLM.Model
	if modelName == "" { modelName = "rule-based" }
	usage := &model.AgentUsage{
		SessionID: sessionID, UserID: userID, Scenario: scenario, Model: modelName, ResponseTimeMs: duration,
	}
	_ = s.repo.CreateUsage(usage)
}

func (s *AgentService) ReflectAndLearn(convID uint, userID uint) {
	if s.llmClient == nil { return }
	history, err := s.repo.GetMessagesByConversationID(convID)
	if err != nil || len(history) < 2 { return }
	s.asyncExtractKnowledge(history, convID)
	s.asyncExtractSkill(history, convID)
	s.asyncCollectExperience(history, userID)
}

func (s *AgentService) asyncCollectExperience(history []model.AgentMessage, userID uint) { }

func (s *AgentService) asyncExtractKnowledge(history []model.AgentMessage, convID uint) {
	p := s.promptTool.BuildKnowledgeExtractionPrompt(history)
	resp, _ := s.llmClient.ChatCompletion([]llm.Message{
		{Role: "system", Content: "你是一个专业的工业设备知识专家。"},
		{Role: "user", Content: p},
	})
	if resp == "" { return }
	var extracted struct {
		Title string `json:"title"`; Type string `json:"type"`; Summary string `json:"summary"`; Details any `json:"details"`; Confidence float64 `json:"confidence"`
	}
	if err := json.Unmarshal([]byte(resp), &extracted); err == nil && extracted.Title != "" {
		detailsJSON, _ := json.Marshal(extracted.Details)
		knowledge := &model.AgentKnowledge{
			ID: fmt.Sprintf("k_%d_%d", convID, time.Now().Unix()), Title: extracted.Title, Type: extracted.Type, Summary: extracted.Summary,
			Details: string(detailsJSON), Confidence: extracted.Confidence, Status: "draft", CreatedBy: fmt.Sprintf("agent:conv_%d", convID),
		}
		_ = s.repo.CreateKnowledge(knowledge)
	}
}

func (s *AgentService) asyncExtractSkill(history []model.AgentMessage, convID uint) {
	p := s.promptTool.BuildSkillExtractionPrompt(history)
	resp, _ := s.llmClient.ChatCompletion([]llm.Message{
		{Role: "system", Content: "你是一个资深的工业诊断专家。"},
		{Role: "user", Content: p},
	})
	if resp == "" { return }
	var extracted struct {
		Name string `json:"name"`; Description string `json:"description"`; ApplicableScenarios []string `json:"applicable_scenarios"`; Steps []any `json:"steps"`
	}
	if err := json.Unmarshal([]byte(resp), &extracted); err == nil && extracted.Name != "" {
		appSce, _ := json.Marshal(extracted.ApplicableScenarios)
		steps, _ := json.Marshal(extracted.Steps)
		skill := &model.AgentSkill{
			Name: extracted.Name, Description: extracted.Description, ApplicableScenarios: string(appSce), Steps: string(steps),
			Status: "draft", CreatedBy: fmt.Sprintf("agent:conv_%d", convID),
		}
		_ = s.repo.CreateSkill(skill)
	}
}
