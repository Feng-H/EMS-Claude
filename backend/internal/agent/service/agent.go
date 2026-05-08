package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"strings"
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
	
	// Tool Registry
	toolRegistry *tool.ToolRegistry

	// Tools (Legacy & Internal)
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
		log.Printf("[AgentService] LLM client initialized (Provider: %s, Model: %s)", config.Cfg.LLM.Provider, config.Cfg.LLM.Model)
	} else {
		log.Printf("[AgentService] Warning: LLM API key is empty, AI features will be disabled")
	}
	
	svc := &AgentService{
		repo:   repo,
		policy: policy.NewPolicyService(),
		toolRegistry:    tool.NewToolRegistry(),
		retrievalTool:   retrievalTool,
		maintenanceTool: maintenanceTool,
		repairTool:      repairTool,
		promptTool:      prompt.NewPromptTool(),
		llmClient:       llmClient,
		maintenanceAnalyzer: analyzer.NewMaintenanceAnalyzer(retrievalTool, maintenanceTool),
		repairAuditAnalyzer: analyzer.NewRepairAuditAnalyzer(retrievalTool, repairTool),
		predictiveAnalyzer:  analyzer.NewPredictiveAnalyzer(repairTool, maintenanceTool, retrievalTool),
	}

	svc.initToolRegistry()
	return svc
}

func (s *AgentService) initToolRegistry() {
	// Register search_equipment
	s.toolRegistry.Register("search_equipment", dto.ToolDefinition{
		Name: "search_equipment", Description: "Search for equipment by name, code or model",
		InputSchema: map[string]interface{}{
			"type": "object", "properties": map[string]interface{}{
				"keyword": map[string]interface{}{"type": "string", "description": "Search keyword"},
			}, "required": []string{"keyword"},
		},
	}, s.handleSearchEquipment, []string{"read:equipment"}, true)

	// Register get_equipment_health
	s.toolRegistry.Register("get_equipment_health", dto.ToolDefinition{
		Name: "get_equipment_health", Description: "Get real-time health analysis and RUL prediction",
		InputSchema: map[string]interface{}{
			"type": "object", "properties": map[string]interface{}{
				"equipment_id": map[string]interface{}{"type": "integer"},
			}, "required": []string{"equipment_id"},
		},
	}, s.handleGetEquipmentHealth, []string{"read:equipment", "read:prediction"}, true)
	
	// Register get_spare_part_inventory
	s.toolRegistry.Register("get_spare_part_inventory", dto.ToolDefinition{
		Name: "get_spare_part_inventory", Description: "Check stock levels of spare parts",
		InputSchema: map[string]interface{}{
			"type": "object", "properties": map[string]interface{}{
				"spare_part_id": map[string]interface{}{"type": "integer"},
				"factory_id":    map[string]interface{}{"type": "integer"},
			}, "required": []string{"spare_part_id"},
		},
	}, s.handleGetSparePartInventory, []string{"read:sparepart"}, true)

	// Register report_repair
	s.toolRegistry.Register("report_repair", dto.ToolDefinition{
		Name: "report_repair", Description: "Submit a new repair request",
		InputSchema: map[string]interface{}{
			"type": "object", "properties": map[string]interface{}{
				"equipment_id":      map[string]interface{}{"type": "integer"},
				"fault_description": map[string]interface{}{"type": "string"},
				"priority":          map[string]interface{}{"type": "integer", "description": "1=High, 2=Medium, 3=Low"},
			}, "required": []string{"equipment_id", "fault_description"},
		},
	}, s.handleReportRepair, []string{"write:repair"}, false)
}

func (s *AgentService) handleSearchEquipment(user model.User, args map[string]interface{}) (interface{}, error) {
	keyword, _ := args["keyword"].(string)
	db := database.GetDB()
	var equipments []model.Equipment
	query := db.Preload("Workshop").Preload("Workshop.Factory")
	if user.Role != "admin" && user.FactoryID != nil {
		query = query.Joins("JOIN workshops ON equipments.workshop_id = workshops.id").
			Where("workshops.factory_id = ?", *user.FactoryID)
	}
	err := query.Where("equipments.name ILIKE ? OR equipments.code ILIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Limit(10).Find(&equipments).Error
	return equipments, err
}

func (s *AgentService) handleGetEquipmentHealth(user model.User, args map[string]interface{}) (interface{}, error) {
	var id uint
	if v, ok := args["equipment_id"].(float64); ok { id = uint(v) } else if v, ok := args["equipment_id"].(int); ok { id = uint(v) }
	return s.GetEquipmentPrediction(id, user)
}

func (s *AgentService) handleGetSparePartInventory(user model.User, args map[string]interface{}) (interface{}, error) {
	var partID uint
	if v, ok := args["spare_part_id"].(float64); ok { partID = uint(v) } else if v, ok := args["spare_part_id"].(int); ok { partID = uint(v) }
	
	db := database.GetDB()
	var inventories []model.SparePartInventory
	query := db.Preload("Factory").Preload("SparePart").Where("spare_part_id = ?", partID)
	if user.Role != "admin" && user.FactoryID != nil {
		query = query.Where("factory_id = ?", *user.FactoryID)
	} else if fid, ok := args["factory_id"]; ok {
		query = query.Where("factory_id = ?", fid)
	}
	err := query.Find(&inventories).Error
	return inventories, err
}

func (s *AgentService) handleReportRepair(user model.User, args map[string]interface{}) (interface{}, error) {
	var equipID uint
	if v, ok := args["equipment_id"].(float64); ok { equipID = uint(v) } else if v, ok := args["equipment_id"].(int); ok { equipID = uint(v) }
	desc, _ := args["fault_description"].(string)
	priority := 2
	if p, ok := args["priority"].(float64); ok { priority = int(p) }
	
	db := database.GetDB()
	var equipment model.Equipment
	if err := db.Joins("JOIN workshops ON workshops.id = equipments.workshop_id").First(&equipment, equipID).Error; err != nil {
		return nil, fmt.Errorf("equipment not found")
	}
	if user.Role != "admin" && user.FactoryID != nil {
		var workshop model.Workshop
		db.First(&workshop, equipment.WorkshopID)
		if workshop.FactoryID != *user.FactoryID {
			return nil, fmt.Errorf("permission denied: equipment belongs to another factory")
		}
	}

	order := model.RepairOrder{
		EquipmentID: equipID, FaultDescription: desc, Priority: priority,
		Status: model.RepairPending, ReporterID: user.ID,
	}
	err := db.Create(&order).Error
	if err != nil { return nil, err }
	return fmt.Sprintf("Repair order #%d created successfully", order.ID), nil
}

func (s *AgentService) RecommendMaintenance(user model.User, req *dto.MaintenanceRecommendRequest) (*dto.AgentResponseEnvelope, error) {
	startTime := time.Now()
	traceID := trace.GenerateTraceID()
	
	agentCtx, err := s.policy.DeriveAgentContext(user.ID, string(user.Role), req.Language)
	if err != nil { return nil, err }

	// Prevent system prompt override for non-admin users
	if req.SystemPrompt != "" && user.Role != "admin" {
		log.Printf("[AgentService] Security warning: Non-admin user %d tried to override system prompt", user.ID)
		req.SystemPrompt = ""
	}

	targetFactoryID := req.FactoryID
	if targetFactoryID == 0 && agentCtx.FactoryID != nil {
		targetFactoryID = *agentCtx.FactoryID
	}
	
	if err := s.policy.ValidateScope(agentCtx, &targetFactoryID); err != nil {
		return nil, err
	}

	analysisResult, err := s.maintenanceAnalyzer.Analyze(req, user)
	if err != nil { return nil, err }

	summary := "建议缩短保养周期，以提高设备可用性。"
	if s.llmClient != nil {
		p := req.SystemPrompt
		if p == "" {
			p = s.promptTool.BuildMaintenanceRecommendPrompt(analysisResult.CurrentPlan, analysisResult.Evidence)
		} else {
			p = fmt.Sprintf("%s\n\n### 原始数据参考\n当前计划: %v\n参考证据: %v", p, analysisResult.CurrentPlan, analysisResult.Evidence)
		}
		resp, err := s.llmClient.ChatCompletion([]llm.Message{
			{Role: "system", Content: "你是一个专业的工业设备管理助手。"},
			{Role: "user", Content: p},
		})
		if err != nil {
			log.Printf("[AgentService] LLM request failed in RecommendMaintenance: %v", err)
		} else if resp != "" {
			summary = resp
		}
	} else if len(analysisResult.Recommendations) > 0 {
		summary = analysisResult.Recommendations[0].Description + "。" + analysisResult.Recommendations[0].Reason
	}

	inputSnap, _ := json.Marshal(req)
	resultJSON, _ := json.Marshal(analysisResult)
	session := &model.AgentSession{
		UserID: user.ID, Scenario: "maintenance_recommendation", FactoryID: &targetFactoryID,
		Language: agentCtx.Language, InputSnapshot: string(inputSnap), TraceID: traceID, Status: "completed",
	}
	if err := s.repo.CreateSession(session); err != nil {
		log.Printf("[AgentService] Failed to create session: %v", err)
		return nil, err
	}

	artifact := &model.AgentArtifact{
		SessionID: session.ID, ArtifactType: "recommendation", Title: "设备保养优化建议",
		Summary: summary, ResultJSON: string(resultJSON), RiskLevel: "medium",
	}
	if err := s.repo.CreateArtifact(artifact); err != nil {
		log.Printf("[AgentService] Failed to create artifact: %v", err)
	}

	for _, ev := range analysisResult.Evidence {
		link := model.AgentEvidenceLink{
			ArtifactID: artifact.ID, EvidenceType: ev.EvidenceType,
			SourceTable: ev.SourceTable, SourceID: ev.SourceID, Excerpt: ev.Excerpt, Score: ev.Score,
		}
		if err := s.repo.CreateEvidenceLinks([]model.AgentEvidenceLink{link}); err != nil {
			log.Printf("[AgentService] Failed to create evidence link: %v", err)
		}
	}

	res := &dto.AgentResponseEnvelope{
		Success: true, TraceID: traceID, Language: agentCtx.Language, Scenario: "maintenance_recommendation",
		ScopeSummary: map[string]interface{}{"factory_id": targetFactoryID},
		Summary: summary, RiskLevel: "medium", ArtifactID: artifact.ID,
		EvidenceCount: len(analysisResult.Evidence), Data: analysisResult,
	}
	s.logUsage(session.ID, user.ID, "maintenance_recommendation", startTime)
	return res, nil
}

func (s *AgentService) AuditRepair(user model.User, req *dto.RepairAuditRequest) (*dto.AgentResponseEnvelope, error) {
	startTime := time.Now()
	traceID := trace.GenerateTraceID()
	
	agentCtx, err := s.policy.DeriveAgentContext(user.ID, string(user.Role), req.Language)
	if err != nil { return nil, err }

	// Prevent system prompt override for non-admin users
	if req.SystemPrompt != "" && user.Role != "admin" {
		log.Printf("[AgentService] Security warning: Non-admin user %d tried to override system prompt", user.ID)
		req.SystemPrompt = ""
	}

	targetFactoryID := req.FactoryID
	if targetFactoryID == 0 && agentCtx.FactoryID != nil {
		targetFactoryID = *agentCtx.FactoryID
	}
	
	if err := s.policy.ValidateScope(agentCtx, &targetFactoryID); err != nil {
		return nil, err
	}

	analysisResult, err := s.repairAuditAnalyzer.Analyze(req, user)
	if err != nil { return nil, err }

	summary := "发现维修异常，建议复核维修质量。"
	if s.llmClient != nil {
		p := req.SystemPrompt
		if p == "" {
			p = s.promptTool.BuildRepairAuditPrompt(analysisResult.Anomalies, analysisResult.Evidence)
		} else {
			p = fmt.Sprintf("%s\n\n### 原始数据参考\n异常项: %v\n参考证据: %v", p, analysisResult.Anomalies, analysisResult.Evidence)
		}
		resp, err := s.llmClient.ChatCompletion([]llm.Message{
			{Role: "system", Content: "你是一个设备维修审计助手。"},
			{Role: "user", Content: p},
		})
		if err != nil {
			log.Printf("[AgentService] LLM request failed in AuditRepair: %v", err)
		} else if resp != "" {
			summary = resp
		}
	} else if stats, ok := analysisResult.Stats.(map[string]interface{}); ok {
		if val, exists := stats["anomaly_summary"]; exists {
			summary = val.(string)
		}
	}

	inputSnap, _ := json.Marshal(req)
	resultJSON, _ := json.Marshal(analysisResult)
	session := &model.AgentSession{
		UserID: user.ID, Scenario: "repair_audit", FactoryID: &targetFactoryID,
		Language: agentCtx.Language, InputSnapshot: string(inputSnap), TraceID: traceID, Status: "completed",
	}
	if err := s.repo.CreateSession(session); err != nil {
		log.Printf("[AgentService] Failed to create session: %v", err)
		return nil, err
	}

	artifact := &model.AgentArtifact{
		SessionID: session.ID, ArtifactType: "audit_report", Title: "设备维修审计报告",
		Summary: summary, ResultJSON: string(resultJSON), RiskLevel: "high",
	}
	if err := s.repo.CreateArtifact(artifact); err != nil {
		log.Printf("[AgentService] Failed to create artifact: %v", err)
	}

	for _, ev := range analysisResult.Evidence {
		link := model.AgentEvidenceLink{
			ArtifactID: artifact.ID, EvidenceType: ev.EvidenceType,
			SourceTable: ev.SourceTable, SourceID: ev.SourceID, Excerpt: ev.Excerpt, Score: ev.Score,
		}
		if err := s.repo.CreateEvidenceLinks([]model.AgentEvidenceLink{link}); err != nil {
			log.Printf("[AgentService] Failed to create evidence link: %v", err)
		}
	}

	res := &dto.AgentResponseEnvelope{
		Success: true, TraceID: traceID, Language: agentCtx.Language, Scenario: "repair_audit",
		ScopeSummary: map[string]interface{}{"factory_id": targetFactoryID},
		Summary: summary, RiskLevel: "high", ArtifactID: artifact.ID,
		EvidenceCount: len(analysisResult.Evidence), Data: analysisResult,
	}
	s.logUsage(session.ID, user.ID, "repair_audit", startTime)
	return res, nil
}

func (s *AgentService) AuditMaintenance(user model.User, req *dto.MaintenanceAuditRequest) (*dto.AgentResponseEnvelope, error) {
	startTime := time.Now()
	traceID := trace.GenerateTraceID()
	
	agentCtx, err := s.policy.DeriveAgentContext(user.ID, string(user.Role), req.Language)
	if err != nil { return nil, err }

	// Prevent system prompt override for non-admin users
	if req.SystemPrompt != "" && user.Role != "admin" {
		log.Printf("[AgentService] Security warning: Non-admin user %d tried to override system prompt", user.ID)
		req.SystemPrompt = ""
	}

	targetFactoryID := req.FactoryID
	if targetFactoryID == 0 && agentCtx.FactoryID != nil {
		targetFactoryID = *agentCtx.FactoryID
	}
	
	if err := s.policy.ValidateScope(agentCtx, &targetFactoryID); err != nil {
		return nil, err
	}

	analysisResult, err := s.maintenanceAnalyzer.Audit(req, user)
	if err != nil { return nil, err }

	summary := analysisResult.AuditSummary
	if s.llmClient != nil {
		p := req.SystemPrompt
		if p == "" {
			p = s.promptTool.BuildMaintenanceAuditPrompt(analysisResult.Anomalies, analysisResult.Evidence)
		} else {
			p = fmt.Sprintf("%s\n\n### 审计发现\n异常: %v\n证据: %v", p, analysisResult.Anomalies, analysisResult.Evidence)
		}
		resp, err := s.llmClient.ChatCompletion([]llm.Message{
			{Role: "system", Content: "你是一个专业的设备保养审计专家。"},
			{Role: "user", Content: p},
		})
		if err == nil && resp != "" {
			summary = resp
		}
	}

	inputSnap, _ := json.Marshal(req)
	resultJSON, _ := json.Marshal(analysisResult)
	session := &model.AgentSession{
		UserID: user.ID, Scenario: "maintenance_audit", FactoryID: &targetFactoryID,
		Language: agentCtx.Language, InputSnapshot: string(inputSnap), TraceID: traceID, Status: "completed",
	}
	_ = s.repo.CreateSession(session)

	artifact := &model.AgentArtifact{
		SessionID: session.ID, ArtifactType: "audit_report", Title: "设备保养合规审计报告",
		Summary: summary, ResultJSON: string(resultJSON), RiskLevel: "medium",
	}
	_ = s.repo.CreateArtifact(artifact)

	res := &dto.AgentResponseEnvelope{
		Success: true, TraceID: traceID, Language: agentCtx.Language, Scenario: "maintenance_audit",
		ScopeSummary: map[string]interface{}{"factory_id": targetFactoryID},
		Summary: summary, RiskLevel: "medium", ArtifactID: artifact.ID,
		EvidenceCount: len(analysisResult.Evidence), Data: analysisResult,
	}
	s.logUsage(session.ID, user.ID, "maintenance_audit", startTime)
	return res, nil
}

func (s *AgentService) Analyze(user model.User, req *dto.AnalyzeRequest) (*dto.AgentResponseEnvelope, error) {
	startTime := time.Now()
	traceID := trace.GenerateTraceID()
	
	agentCtx, err := s.policy.DeriveAgentContext(user.ID, string(user.Role), req.Language)
	if err != nil { return nil, err }

	// Prevent system prompt override for non-admin users
	if req.SystemPrompt != "" && user.Role != "admin" {
		log.Printf("[AgentService] Security warning: Non-admin user %d tried to override system prompt", user.ID)
		req.SystemPrompt = ""
	}

	targetFactoryID := req.FactoryID
	if targetFactoryID == 0 && agentCtx.FactoryID != nil {
		targetFactoryID = *agentCtx.FactoryID
	}
	
	if err := s.policy.ValidateScope(agentCtx, &targetFactoryID); err != nil {
		return nil, err
	}

	// 1. Gather context based on the question (Entity extraction simplified for MVP)
	eqID := s.extractEquipmentID(req.Question, user)
	contextMap := make(map[string]interface{})
	
	if eqID != 1 {
		profile, _ := s.retrievalTool.GetEquipmentProfile(eqID, user)
		health, _ := s.GetEquipmentPrediction(eqID, user)
		failureStats, _ := s.repairTool.GetFailureStats(eqID, user)
		contextMap["equipment_profile"] = profile
		contextMap["equipment_health"] = health
		contextMap["failure_stats"] = failureStats
	}

	// 2. Generate summary via LLM
	summary := "已为您完成多维度分析。建议关注设备的 RUL 变化及维护成本趋势。"
	if s.llmClient != nil {
		p := req.SystemPrompt
		if p == "" {
			p = s.promptTool.BuildGenericAnalysisPrompt(req.Question, contextMap)
		} else {
			p = fmt.Sprintf("%s\n\n### 补充背景\n%v", p, contextMap)
		}
		
		resp, err := s.llmClient.ChatCompletion([]llm.Message{
			{Role: "system", Content: "你是一个顶级的工业资产战略分析师。"},
			{Role: "user", Content: p},
		})
		if err == nil && resp != "" {
			summary = resp
		}
	}

	analysisData := dto.AnalyzeData{
		KeyFindings: []string{},
		Evidence:    []dto.EvidenceItem{},
	}

	if health, ok := contextMap["equipment_health"].(map[string]interface{}); ok {
		if rul, ok := health["rul"].(*dto.RULPrediction); ok && rul.EstimatedRULDays < 10 {
			analysisData.KeyFindings = append(analysisData.KeyFindings, fmt.Sprintf("设备剩余寿命仅剩 %d 天，存在停机风险", rul.EstimatedRULDays))
			analysisData.Evidence = append(analysisData.Evidence, dto.EvidenceItem{
				EvidenceType: "prediction", Title: "RUL 预测", Excerpt: fmt.Sprintf("预计剩余寿命: %d天", rul.EstimatedRULDays), Score: 0.95,
			})
		}
	}

	if stats, ok := contextMap["failure_stats"].(map[string]interface{}); ok {
		if count, ok := stats["repair_count"].(int); ok && count > 5 {
			analysisData.KeyFindings = append(analysisData.KeyFindings, "近期维修频率较高，建议核查根本原因")
		}
	}

	if len(analysisData.KeyFindings) == 0 {
		analysisData.KeyFindings = append(analysisData.KeyFindings, "设备运行状况平稳", "未发现近期异常趋势")
	}

	inputSnap, _ := json.Marshal(req)
	resultJSON, _ := json.Marshal(analysisData)
	session := &model.AgentSession{
		UserID: user.ID, Scenario: "analysis", FactoryID: &targetFactoryID,
		Language: agentCtx.Language, InputSnapshot: string(inputSnap), TraceID: traceID, Status: "completed",
	}
	_ = s.repo.CreateSession(session)

	artifact := &model.AgentArtifact{
		SessionID: session.ID, ArtifactType: "analysis_result", Title: "深度业务分析报告",
		Summary: summary, ResultJSON: string(resultJSON), RiskLevel: "medium",
	}
	_ = s.repo.CreateArtifact(artifact)

	res := &dto.AgentResponseEnvelope{
		Success: true, TraceID: traceID, Language: agentCtx.Language, Scenario: "analysis",
		ScopeSummary: map[string]interface{}{"factory_id": targetFactoryID, "equipment_id": eqID},
		Summary: summary, RiskLevel: "medium", ArtifactID: artifact.ID,
		EvidenceCount: len(analysisData.Evidence), Data: analysisData,
	}
	s.logUsage(session.ID, user.ID, "analysis", startTime)
	return res, nil
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

func (s *AgentService) GetSession(id uint, userID uint, role string) (*dto.AgentSessionResponse, error) {
	session, err := s.repo.GetSessionByID(id)
	if err != nil { return nil, err }
	
	// Ownership check
	if role != "admin" && session.UserID != userID {
		return nil, fmt.Errorf("permission denied: unauthorized access to session")
	}

	res := &dto.AgentSessionResponse{
		ID: session.ID, UserID: session.UserID, Scenario: session.Scenario, Language: session.Language,
		Status: session.Status, TraceID: session.TraceID, CreatedAt: session.CreatedAt,
	}
	if session.FactoryID != nil { res.FactoryID = *session.FactoryID }
	for _, a := range session.Artifacts { res.Artifacts = append(res.Artifacts, a.ID) }
	return res, nil
}

func (s *AgentService) GetArtifact(id uint, userID uint, role string) (*dto.AgentArtifactResponse, error) {
	artifact, err := s.repo.GetArtifactByID(id)
	if err != nil { return nil, err }
	
	// Artifact belongs to a session, check session ownership
	session, err := s.repo.GetSessionByID(artifact.SessionID)
	if err != nil { return nil, err }
	if role != "admin" && session.UserID != userID {
		return nil, fmt.Errorf("permission denied: unauthorized access to artifact")
	}

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

func (s *AgentService) ListKnowledges(status string, query string, eqTypeID *uint) ([]model.AgentKnowledge, error) {
	return s.repo.ListKnowledges(status, query, 100)
}

// =====================================================
// Phase 2: Chat & Conversational Logic
// =====================================================

func (s *AgentService) Chat(user model.User, req *dto.ChatRequest) (*dto.ChatResponse, error) {
	startTime := time.Now()
	traceID := trace.GenerateTraceID()

	// 1. 获取或创建对话
	var convID = req.ConversationID
	if convID == 0 {
		title := req.Message
		runes := []rune(title)
		if len(runes) > 30 {
			title = string(runes[:27]) + "..."
		}
		newConv := &model.AgentConversation{
			UserID: user.ID,
			Title:  title,
		}
		if err := s.repo.CreateConversation(newConv); err != nil {
			log.Printf("[AgentService] Failed to create conversation: %v", err)
			return nil, err
		}
		convID = newConv.ID
	}

	// 2. 持久化用户消息
	_ = s.repo.CreateMessage(&model.AgentMessage{ConversationID: convID, Role: "user", Content: req.Message})

	// 3. 注入用户个性化经验 (Milestone P)
	activeExps, _ := s.repo.ListActiveExperiences(user.ID)
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
		res, err := s.ExecuteSkill(user, &skill, req)
		if err == nil { reply = res.Summary + expContext }
	}

	// 5. 退回到标准对话
	if reply == "" {
		history, _ := s.repo.GetMessagesByConversationID(convID)
		
		// Context retrieval: Find relevant equipment or knowledge
		eqID := s.extractEquipmentID(req.Message, user)
		businessContext := ""
		if eqID != 1 { // If a specific equipment was found
			profile, _ := s.retrievalTool.GetEquipmentProfile(eqID, user)
			health, _ := s.GetEquipmentPrediction(eqID, user)
			profileJSON, _ := json.Marshal(profile)
			healthJSON, _ := json.Marshal(health)
			businessContext += fmt.Sprintf("\n### 当前讨论的设备上下文\n基础信息: %s\n健康分析: %s\n", profileJSON, healthJSON)
		}
		
		// Retrieve relevant knowledge
		knowledge, _ := s.retrievalTool.SearchManualKnowledge(req.Message, nil, user)
		if len(knowledge) > 0 {
			businessContext += "\n### 相关知识参考\n"
			for i, k := range knowledge {
				if i >= 2 { break }
				businessContext += fmt.Sprintf("- [%s]: %s\n", k.Title, k.Excerpt)
			}
		}

		llmMsgs := []llm.Message{
			{Role: "system", Content: "你是一个顶级的工业资产战略专家。你拥有‘L4 级主动洞察’权限，可以基于全生命周期成本 (TCO)、资产退役 ROI 评价、剩余健康寿命 (RUL) 和亚健康故障征兆进行跨维度的深度分析。请使用中文回答，结论必须引用系统中的财务与技术证据。" + expContext + businessContext},
		}

		// Prevent system prompt override for non-admin users
		if req.SystemPrompt != "" {
			if user.Role == "admin" {
				llmMsgs[0].Content = req.SystemPrompt
			} else {
				log.Printf("[AgentService] Security warning: Non-admin user %d tried to override chat system prompt", user.ID)
			}
		}

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
	go s.ReflectAndLearn(convID, user.ID)

	// 8. 记录使用情况
	s.logUsage(convID, user.ID, "chat", startTime)

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

func (s *AgentService) extractEquipmentID(message string, user model.User) uint {
	if config.Cfg.Storage.Mode == "memory" {
		return 1 // Fallback for memory mode demo
	}
	
	var equipments []model.Equipment
	// Load all equipment for context matching. In a real system, use semantic search or NER.
	// Filter by factory at query level for efficiency
	query := database.GetDB().Joins("JOIN workshops ON workshops.id = equipments.workshop_id")
	if user.Role != "admin" && user.FactoryID != nil {
		query = query.Where("workshops.factory_id = ?", *user.FactoryID)
	}
	
	if err := query.Find(&equipments).Error; err != nil {
		return 1
	}

	for _, eq := range equipments {
		if strings.Contains(message, eq.Code) || (eq.Name != "" && strings.Contains(message, eq.Name)) {
			return eq.ID
		}
	}
	
	return 1 // Default fallback
}

func (s *AgentService) GetConversation(id uint, userID uint, role string) (*dto.ConversationResponse, error) {
	conv, err := s.repo.GetConversationByID(id)
	if err != nil { return nil, err }

	// Ownership check
	if role != "admin" && conv.UserID != userID {
		return nil, fmt.Errorf("permission denied: unauthorized access to conversation")
	}

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

func (s *AgentService) ListSkills(status string, query string) ([]dto.SkillResponse, error) {
	skills, err := s.repo.ListSkills(status, query, 100)
	if err != nil { return nil, err }
	results := make([]dto.SkillResponse, len(skills))
	for i, sk := range skills { results[i] = *s.mapSkillToResponse(&sk) }
	return results, nil
}

func (s *AgentService) ExecuteSkill(user model.User, skill *model.AgentSkill, req *dto.ChatRequest) (*dto.AgentResponseEnvelope, error) {
	var steps []struct {
		Step   int    `json:"step"`
		Action string `json:"action"`
		Tool   string `json:"tool"`
	}
	_ = json.Unmarshal([]byte(skill.Steps), &steps)
	
	eqID := s.extractEquipmentID(req.Message, user)
	
	evidence := []dto.EvidenceItem{}
	for _, step := range steps {
		switch step.Tool {
		case "get_equipment_profile":
			if res, err := s.retrievalTool.GetEquipmentProfile(eqID, user); err == nil {
				resJSON, _ := json.Marshal(res)
				evidence = append(evidence, dto.EvidenceItem{
					EvidenceType: "equipment_profile", Title: "设备基础信息", Excerpt: string(resJSON), Score: 1.0,
				})
			}
		case "get_failure_stats":
			if res, err := s.repairTool.GetFailureStats(eqID, user); err == nil {
				resJSON, _ := json.Marshal(res)
				evidence = append(evidence, dto.EvidenceItem{
					EvidenceType: "failure_stats", Title: "故障统计分析", Excerpt: string(resJSON), Score: 0.95,
				})
			}
		case "get_cost_analysis":
			if res, err := s.repairTool.GetCostByEquipmentID(eqID, user); err == nil {
				resJSON, _ := json.Marshal(res)
				evidence = append(evidence, dto.EvidenceItem{
					EvidenceType: "cost_analysis", Title: "维修成本分析", Excerpt: string(resJSON), Score: 0.9,
				})
			}
		case "get_maintenance_compliance":
			if res, err := s.maintenanceTool.GetMaintenanceCompliance(eqID, user); err == nil {
				resJSON, _ := json.Marshal(res)
				evidence = append(evidence, dto.EvidenceItem{
					EvidenceType: "maintenance_compliance", Title: "保养合规性评估", Excerpt: string(resJSON), Score: 0.85,
				})
			}
		case "get_failure_distribution":
			auditReq := &dto.RepairAuditRequest{EquipmentTypeID: 12}
			res, err := s.repairAuditAnalyzer.Analyze(auditReq, user)
			if err == nil { evidence = append(evidence, res.Evidence...) }
		case "search_manual_knowledge":
			res, err := s.retrievalTool.SearchManualKnowledge(req.Message, nil, user)
			if err == nil {
				evidence = append(evidence, res...)
			}
		case "predict_remaining_life":
			if res, err := s.predictiveAnalyzer.PredictRUL(eqID, user); err == nil {
				resJSON, _ := json.Marshal(res)
				evidence = append(evidence, dto.EvidenceItem{
					EvidenceType: "prediction", Title: "RUL 剩余健康寿命预测",
					Excerpt: string(resJSON), Score: 0.98,
				})
			}
		case "detect_symptoms":
			if res, err := s.predictiveAnalyzer.DetectSymptoms(eqID, user); err == nil {
				resJSON, _ := json.Marshal(res)
				evidence = append(evidence, dto.EvidenceItem{
					EvidenceType: "symptoms", Title: "设备亚健康征兆识别",
					Excerpt: string(resJSON), Score: 0.92,
				})
			}
		case "get_tco_analysis":
			if res, err := s.predictiveAnalyzer.CalculateTCO(eqID, user); err == nil {
				resJSON, _ := json.Marshal(res)
				evidence = append(evidence, dto.EvidenceItem{
					EvidenceType: "tco", Title: "资产总持有成本(TCO)分析",
					Excerpt: string(resJSON), Score: 0.95,
				})
			}
		case "get_retirement_recommendation":
			if res, err := s.predictiveAnalyzer.EvaluateRetirement(eqID, user); err == nil {
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
		resp, err := s.llmClient.ChatCompletion([]llm.Message{
			{Role: "system", Content: "你是一个专业的工业设备管理助手，正在执行预定义的分析技能。"},
			{Role: "user", Content: prompt},
		})
		if err != nil {
			log.Printf("[AgentService] LLM request failed in SkillExecution: %v", err)
		} else if resp != "" {
			summary = resp
		}
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

func (s *AgentService) Subscribe(userID uint, pushType string, enabled bool, scope any, webhookURL string) error {
	scopeJSON, _ := json.Marshal(scope)

	// Check if already exists
	db := database.GetDB()
	var sub model.AgentPushSubscription
	err := db.Where("user_id = ? AND push_type = ?", userID, pushType).First(&sub).Error

	if err == nil {
		// Update
		sub.Enabled = enabled
		sub.Scope = string(scopeJSON)
		sub.WebhookURL = webhookURL
		return db.Save(&sub).Error
	}

	// Create
	sub = model.AgentPushSubscription{
		UserID: userID, PushType: pushType, Enabled: enabled, Scope: string(scopeJSON),
		WebhookURL: webhookURL,
	}
	return db.Create(&sub).Error
}

func (s *AgentService) ListSubscriptions(userID uint) ([]model.AgentPushSubscription, error) {
	db := database.GetDB()
	var subs []model.AgentPushSubscription
	err := db.Where("user_id = ?", userID).Find(&subs).Error
	return subs, err
}

func (s *AgentService) NotifyEvent(eventType string, targetID uint, context map[string]interface{}) {
	// 1. 查找所有启用且类型匹配的订阅
	var subs []model.AgentPushSubscription
	database.GetDB().Preload("User").Where("push_type = ? AND enabled = ?", eventType, true).Find(&subs)

	for _, sub := range subs {
		// 2. 为每个订阅用户执行权限与范围校验
		user := model.User{}
		if sub.UserID > 0 {
			database.GetDB().First(&user, sub.UserID)
		} else {
			continue
		}

		// 检查设备是否属于该用户的工厂范围
		var equipment model.Equipment
		if err := database.GetDB().Joins("JOIN workshops ON workshops.id = equipments.workshop_id").First(&equipment, targetID).Error; err != nil {
			continue
		}
		if user.Role != "admin" && user.FactoryID != nil {
			var workshop model.Workshop
			database.GetDB().First(&workshop, equipment.WorkshopID)
			if workshop.FactoryID != *user.FactoryID {
				continue // 跨工厂，跳过此订阅者
			}
		}

		// 3. 在用户上下文中执行分析 (确保 RUL/TCO 等逻辑应用了正确的工厂参数，且报错能被捕捉)
		prediction, err := s.predictiveAnalyzer.PredictRUL(targetID, user)
		if err != nil || prediction.EstimatedRULDays >= 7 {
			continue
		}

		// 4. 创建 Artifact (如果尚未为此事件创建，或需要个性化)
		// 注意：实际生产中可能需要一个缓存避免重复创建完全一样的 Artifact，这里简化为每人一个或共享
		artifact := &model.AgentArtifact{
			ArtifactType: "proactive_push",
			Title:        "设备停机风险预警",
			Summary:      fmt.Sprintf("Agent 自动巡检发现风险：设备【%s】预计剩余健康寿命仅剩 %d 天，建议立即干预。", equipment.Name, prediction.EstimatedRULDays),
			ResultJSON:   "{\"prediction\": \"high_risk\"}",
			RiskLevel:    "high",
		}
		_ = s.repo.CreateArtifact(artifact)

		// 5. 投递
		go s.deliverPush(sub, artifact)
	}
}

func (s *AgentService) deliverPush(sub model.AgentPushSubscription, artifact *model.AgentArtifact) {
	payload := map[string]interface{}{
		"event":       "agent_alert",
		"artifact_id": artifact.ID,
		"title":       artifact.Title,
		"summary":     artifact.Summary,
		"risk_level":  artifact.RiskLevel,
		"timestamp":   time.Now().Unix(),
	}
	body, _ := json.Marshal(payload)
	
	logRecord := &model.AgentPushLog{
		SubscriptionID: sub.ID,
		ArtifactID:     artifact.ID,
		Payload:        string(body),
		Status:         "pending",
	}
	database.GetDB().Create(logRecord)

	if sub.WebhookURL != "" {
		req, _ := http.NewRequest("POST", sub.WebhookURL, strings.NewReader(string(body)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-EMS-Event", "agent_alert")
		// Future: Sign payload with sub.Secret

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		
		now := time.Now()
		if err != nil {
			logRecord.Status = "failed"
			logRecord.ErrorMessage = err.Error()
		} else {
			defer resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				logRecord.Status = "success"
				logRecord.DeliveredAt = &now
			} else {
				logRecord.Status = "failed"
				logRecord.ErrorMessage = fmt.Sprintf("HTTP %d", resp.StatusCode)
			}
		}
		database.GetDB().Save(logRecord)
	}
}
func (s *AgentService) GetEquipmentPrediction(equipmentID uint, user model.User) (map[string]interface{}, error) {
	// Web端调用，使用真实用户权限进行隔离校验
	rul, _ := s.predictiveAnalyzer.PredictRUL(equipmentID, user)
	tco, _ := s.predictiveAnalyzer.CalculateTCO(equipmentID, user)
	symptoms, _ := s.predictiveAnalyzer.DetectSymptoms(equipmentID, user)

	return map[string]interface{}{
		"rul":      rul,
		"tco":      tco,
		"symptoms": symptoms,
	}, nil
}

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
	if err := s.repo.CreateUsage(usage); err != nil {
		log.Printf("[AgentService] Failed to create usage record: %v", err)
	}
}

func (s *AgentService) ReflectAndLearn(convID uint, userID uint) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[AgentService] PANIC in ReflectAndLearn (conv=%d, user=%d): %v", convID, userID, r)
		}
	}()

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
	resp, err := s.llmClient.ChatCompletion([]llm.Message{
		{Role: "system", Content: "你是一个专业的工业设备知识专家。"},
		{Role: "user", Content: p},
	})
	if err != nil {
		log.Printf("[AgentService] LLM request failed in asyncExtractKnowledge: %v", err)
		return
	}
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
	resp, err := s.llmClient.ChatCompletion([]llm.Message{
		{Role: "system", Content: "你是一个资深的工业诊断专家。"},
		{Role: "user", Content: p},
	})
	if err != nil {
		log.Printf("[AgentService] LLM request failed in asyncExtractSkill: %v", err)
		return
	}
	if resp == "" { return }
	var extracted struct {
		Name string `json:"name"`; Description string `json:"description"`; ApplicableScenarios []string `json:"applicable_scenarios"`; Steps []any `json:"steps"`
	}
	if err := json.Unmarshal([]byte(resp), &extracted); err == nil && extracted.Name != "" {
		appSce, _ := json.Marshal(extracted.ApplicableScenarios)
		steps, _ := json.Marshal(extracted.Steps)
		skill := &model.AgentSkill{
			Name: extracted.Name, Description: extracted.Description, ApplicableScenarios: string(appSce), Steps: string(steps),
			Status: "draft",
		}
		skill.CreatedBy = fmt.Sprintf("agent:conv_%d", convID)
		_ = s.repo.CreateSkill(skill)
	}
}
