package service

import (
	"encoding/json"
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
	}
}

func (s *AgentService) RecommendMaintenance(userID uint, role string, req *dto.MaintenanceRecommendRequest) (*dto.AgentResponseEnvelope, error) {
	startTime := time.Now()
	traceID := trace.GenerateTraceID()
	
	// 1. Derive authorized context
	agentCtx, err := s.policy.DeriveAgentContext(userID, role, req.Language)
	if err != nil {
		return nil, err
	}

	// 2. Validate requested scope
	targetFactoryID := req.FactoryID
	if targetFactoryID == 0 && agentCtx.FactoryID != nil {
		targetFactoryID = *agentCtx.FactoryID
	}
	
	if err := s.policy.ValidateScope(agentCtx, &targetFactoryID); err != nil {
		return nil, err
	}

	// 3. Trigger analyzer
	analysisResult, err := s.maintenanceAnalyzer.Analyze(req)
	if err != nil {
		return nil, err
	}

	// 4. Call LLM for summarization
	summary := "建议缩短保养周期，以提高设备可用性。" // Default
	if s.llmClient != nil {
		prompt := s.promptTool.BuildMaintenanceRecommendPrompt(analysisResult.CurrentPlan, analysisResult.Evidence)
		resp, err := s.llmClient.ChatCompletion([]llm.Message{
			{Role: "system", Content: "你是一个专业的工业设备管理助手。"},
			{Role: "user", Content: prompt},
		})
		if err == nil {
			summary = resp
		}
	} else if len(analysisResult.Recommendations) > 0 {
		summary = analysisResult.Recommendations[0].Description + "。" + analysisResult.Recommendations[0].Reason
	}

	// 5. Persist session and artifact
	inputSnap, _ := json.Marshal(req)
	resultJSON, _ := json.Marshal(analysisResult)
	
	session := &model.AgentSession{
		UserID:        userID,
		Scenario:      "maintenance_recommendation",
		FactoryID:     &targetFactoryID,
		Language:      agentCtx.Language,
		InputSnapshot: string(inputSnap),
		TraceID:       traceID,
		Status:        "completed",
	}
	_ = s.repo.CreateSession(session)

	artifact := &model.AgentArtifact{
		SessionID:    session.ID,
		ArtifactType: "recommendation",
		Title:        "设备保养优化建议",
		Summary:      summary,
		ResultJSON:   string(resultJSON),
		RiskLevel:    "medium",
	}
	_ = s.repo.CreateArtifact(artifact)

	// Add evidence links
	for _, ev := range analysisResult.Evidence {
		link := model.AgentEvidenceLink{
			ArtifactID:   artifact.ID,
			EvidenceType: ev.EvidenceType,
			SourceTable:  ev.SourceTable,
			SourceID:     ev.SourceID,
			Excerpt:      ev.Excerpt,
			Score:        ev.Score,
		}
		_ = s.repo.CreateEvidenceLinks([]model.AgentEvidenceLink{link})
	}

	res := &dto.AgentResponseEnvelope{
		Success:  true,
		TraceID:  traceID,
		Language: agentCtx.Language,
		Scenario: "maintenance_recommendation",
		ScopeSummary: map[string]interface{}{
			"factory_id": targetFactoryID,
		},
		Summary:       summary,
		RiskLevel:     "medium",
		ArtifactID:    artifact.ID,
		EvidenceCount: len(analysisResult.Evidence),
		Data:          analysisResult,
	}

	// 6. Log usage
	s.logUsage(session.ID, userID, "maintenance_recommendation", startTime)

	return res, nil
}

func (s *AgentService) AuditRepair(userID uint, role string, req *dto.RepairAuditRequest) (*dto.AgentResponseEnvelope, error) {
	startTime := time.Now()
	traceID := trace.GenerateTraceID()
	
	// 1. Derive authorized context
	agentCtx, err := s.policy.DeriveAgentContext(userID, role, req.Language)
	if err != nil {
		return nil, err
	}

	// 2. Validate requested scope
	targetFactoryID := req.FactoryID
	if targetFactoryID == 0 && agentCtx.FactoryID != nil {
		targetFactoryID = *agentCtx.FactoryID
	}
	
	if err := s.policy.ValidateScope(agentCtx, &targetFactoryID); err != nil {
		return nil, err
	}

	// 3. Trigger analyzer
	analysisResult, err := s.repairAuditAnalyzer.Analyze(req)
	if err != nil {
		return nil, err
	}

	// 4. Call LLM for summarization
	summary := "发现维修异常，建议复核维修质量。"
	if s.llmClient != nil {
		prompt := s.promptTool.BuildRepairAuditPrompt(analysisResult.Anomalies, analysisResult.Evidence)
		resp, err := s.llmClient.ChatCompletion([]llm.Message{
			{Role: "system", Content: "你是一个设备维修审计助手。"},
			{Role: "user", Content: prompt},
		})
		if err == nil {
			summary = resp
		}
	} else if stats, ok := analysisResult.Stats.(map[string]interface{}); ok {
		if val, exists := stats["anomaly_summary"]; exists {
			summary = val.(string)
		}
	}

	// 5. Persist session and artifact
	inputSnap, _ := json.Marshal(req)
	resultJSON, _ := json.Marshal(analysisResult)

	session := &model.AgentSession{
		UserID:        userID,
		Scenario:      "repair_audit",
		FactoryID:     &targetFactoryID,
		Language:      agentCtx.Language,
		InputSnapshot: string(inputSnap),
		TraceID:       traceID,
		Status:        "completed",
	}
	_ = s.repo.CreateSession(session)

	artifact := &model.AgentArtifact{
		SessionID:    session.ID,
		ArtifactType: "audit_report",
		Title:        "设备维修审计报告",
		Summary:      summary,
		ResultJSON:   string(resultJSON),
		RiskLevel:    "high",
	}
	_ = s.repo.CreateArtifact(artifact)

	// Add evidence links
	for _, ev := range analysisResult.Evidence {
		link := model.AgentEvidenceLink{
			ArtifactID:   artifact.ID,
			EvidenceType: ev.EvidenceType,
			SourceTable:  ev.SourceTable,
			SourceID:     ev.SourceID,
			Excerpt:      ev.Excerpt,
			Score:        ev.Score,
		}
		_ = s.repo.CreateEvidenceLinks([]model.AgentEvidenceLink{link})
	}

	res := &dto.AgentResponseEnvelope{
		Success:  true,
		TraceID:  traceID,
		Language: agentCtx.Language,
		Scenario: "repair_audit",
		ScopeSummary: map[string]interface{}{
			"factory_id": targetFactoryID,
		},
		Summary:       summary,
		RiskLevel:     "high",
		ArtifactID:    artifact.ID,
		EvidenceCount: len(analysisResult.Evidence),
		Data:          analysisResult,
	}

	// 6. Log usage
	s.logUsage(session.ID, userID, "repair_audit", startTime)

	return res, nil
}

func (s *AgentService) AuditMaintenance(userID uint, role string, req *dto.MaintenanceAuditRequest) (*dto.AgentResponseEnvelope, error) {
	traceID := trace.GenerateTraceID()
	
	return &dto.AgentResponseEnvelope{
		Success:  true,
		TraceID:  traceID,
		Language: "zh-CN",
		Scenario: "maintenance_audit",
		Summary:  "这是保养审计的占位符响应（开发中）",
		Data:     dto.MaintenanceAuditData{},
	}, nil
}

func (s *AgentService) Analyze(userID uint, role string, req *dto.AnalyzeRequest) (*dto.AgentResponseEnvelope, error) {
	traceID := trace.GenerateTraceID()
	
	return &dto.AgentResponseEnvelope{
		Success:  true,
		TraceID:  traceID,
		Language: "zh-CN",
		Scenario: "analysis",
		Summary:  "这是分析助手的占位符响应（开发中）",
		Data:     dto.AnalyzeData{},
	}, nil
}

func (s *AgentService) ListSessions(userID uint, limit int) ([]dto.AgentSessionResponse, error) {
	sessions, err := s.repo.ListSessionsByUserID(userID, limit)
	if err != nil {
		return nil, err
	}

	results := make([]dto.AgentSessionResponse, len(sessions))
	for i, s := range sessions {
		results[i] = dto.AgentSessionResponse{
			ID:         s.ID,
			UserID:     s.UserID,
			Scenario:   s.Scenario,
			Language:   s.Language,
			Status:     s.Status,
			TraceID:    s.TraceID,
			CreatedAt:  s.CreatedAt,
		}
		if s.FactoryID != nil {
			results[i].FactoryID = *s.FactoryID
		}
	}
	return results, nil
}

func (s *AgentService) GetSession(id uint) (*dto.AgentSessionResponse, error) {
	session, err := s.repo.GetSessionByID(id)
	if err != nil {
		return nil, err
	}

	res := &dto.AgentSessionResponse{
		ID:         session.ID,
		UserID:     session.UserID,
		Scenario:   session.Scenario,
		Language:   session.Language,
		Status:     session.Status,
		TraceID:    session.TraceID,
		CreatedAt:  session.CreatedAt,
	}
	if session.FactoryID != nil {
		res.FactoryID = *session.FactoryID
	}
	
	for _, a := range session.Artifacts {
		res.Artifacts = append(res.Artifacts, a.ID)
	}

	return res, nil
}

func (s *AgentService) GetArtifact(id uint) (*dto.AgentArtifactResponse, error) {
	artifact, err := s.repo.GetArtifactByID(id)
	if err != nil {
		return nil, err
	}

	res := &dto.AgentArtifactResponse{
		ID:           artifact.ID,
		SessionID:    artifact.SessionID,
		ArtifactType: artifact.ArtifactType,
		Title:        artifact.Title,
		Summary:      artifact.Summary,
		RiskLevel:    artifact.RiskLevel,
		CreatedAt:    artifact.CreatedAt,
	}

	_ = json.Unmarshal([]byte(artifact.ResultJSON), &res.ResultJSON)

	return res, nil
}

func (s *AgentService) logUsage(sessionID, userID uint, scenario string, startTime time.Time) {
	duration := time.Since(startTime).Milliseconds()
	modelName := config.Cfg.LLM.Model
	if modelName == "" {
		modelName = "rule-based"
	}

	usage := &model.AgentUsage{
		SessionID:      sessionID,
		UserID:         userID,
		Scenario:       scenario,
		Model:          modelName,
		ResponseTimeMs: duration,
	}
	_ = s.repo.CreateUsage(usage)
}
