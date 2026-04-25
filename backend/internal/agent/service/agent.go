package service

import (
	"github.com/ems/backend/internal/agent/analyzer"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/policy"
	"github.com/ems/backend/internal/agent/prompt"
	"github.com/ems/backend/internal/agent/repository"
	"github.com/ems/backend/internal/agent/tool"
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

	return &dto.AgentResponseEnvelope{
		Success:  true,
		TraceID:  traceID,
		Language: agentCtx.Language,
		Scenario: "maintenance_recommendation",
		ScopeSummary: map[string]interface{}{
			"factory_id": targetFactoryID,
		},
		Summary:       summary,
		RiskLevel:     "medium",
		EvidenceCount: len(analysisResult.Evidence),
		Data:          analysisResult,
	}, nil
}

func (s *AgentService) AuditRepair(userID uint, role string, req *dto.RepairAuditRequest) (*dto.AgentResponseEnvelope, error) {
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

	return &dto.AgentResponseEnvelope{
		Success:  true,
		TraceID:  traceID,
		Language: agentCtx.Language,
		Scenario: "repair_audit",
		ScopeSummary: map[string]interface{}{
			"factory_id": targetFactoryID,
		},
		Summary:       summary,
		RiskLevel:     "high",
		EvidenceCount: len(analysisResult.Evidence),
		Data:          analysisResult,
	}, nil
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

func (s *AgentService) GetSession(id uint) (*dto.AgentSessionResponse, error) {
	// TODO: Implement logic
	return &dto.AgentSessionResponse{ID: id}, nil
}

func (s *AgentService) GetArtifact(id uint) (*dto.AgentArtifactResponse, error) {
	// TODO: Implement logic
	return &dto.AgentArtifactResponse{ID: id}, nil
}
