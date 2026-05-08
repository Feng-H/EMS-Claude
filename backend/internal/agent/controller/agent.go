package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/service"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/database"
	"github.com/ems/backend/pkg/trace"
	"github.com/gin-gonic/gin"
)

type AgentController struct {
	agentService *service.AgentService
}

func NewAgentController() *AgentController {
	return &AgentController{
		agentService: service.NewAgentService(),
	}
}

// requireAuth extracts user ID and role, aborts with 401 if missing.
func requireAuth(c *gin.Context) (uint, string, bool) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "UNAUTHORIZED", Message: "User not authenticated"},
		})
		return 0, "", false
	}
	role, ok := middleware.GetUserRole(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "UNAUTHORIZED", Message: "User role not found"},
		})
		return 0, "", false
	}
	return userID, role, true
}

// RecommendMaintenance generates maintenance optimization recommendations
func (ctrl *AgentController) RecommendMaintenance(c *gin.Context) {
	var req dto.MaintenanceRecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: err.Error()},
		})
		return
	}
	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}
	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: "User not found"},
		})
		return
	}
	result, err := ctrl.agentService.RecommendMaintenance(user, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// AuditRepair audits repair reasonableness
func (ctrl *AgentController) AuditRepair(c *gin.Context) {
	var req dto.RepairAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: err.Error()},
		})
		return
	}
	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}
	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: "User not found"},
		})
		return
	}
	result, err := ctrl.agentService.AuditRepair(user, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// AuditMaintenance audits maintenance plan
func (ctrl *AgentController) AuditMaintenance(c *gin.Context) {
	var req dto.MaintenanceAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: err.Error()},
		})
		return
	}
	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}
	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: "User not found"},
		})
		return
	}
	result, err := ctrl.agentService.AuditMaintenance(user, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Analyze answer management questions
func (ctrl *AgentController) Analyze(c *gin.Context) {
	var req dto.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: err.Error()},
		})
		return
	}
	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}
	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: "User not found"},
		})
		return
	}
	result, err := ctrl.agentService.Analyze(user, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetSession returns session metadata
func (ctrl *AgentController) GetSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: "Invalid ID"},
		})
		return
	}
	userID, role, ok := requireAuth(c)
	if !ok {
		return
	}
	result, err := ctrl.agentService.GetSession(uint(id), userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetArtifact returns one artifact
func (ctrl *AgentController) GetArtifact(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: "Invalid ID"},
		})
		return
	}
	userID, role, ok := requireAuth(c)
	if !ok {
		return
	}
	result, err := ctrl.agentService.GetArtifact(uint(id), userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ListSessions returns user's recent agent sessions
func (ctrl *AgentController) ListSessions(c *gin.Context) {
	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}
	result, err := ctrl.agentService.ListSessions(userID, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// AuditKnowledge confirms or rejects a knowledge draft
func (ctrl *AgentController) AuditKnowledge(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"` // confirmed, rejected
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: err.Error()},
		})
		return
	}

	userID, role, ok := requireAuth(c)
	if !ok {
		return
	}
	// Only admin or manager can audit knowledge
	if role != "admin" && role != "manager" {
		c.JSON(http.StatusForbidden, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "FORBIDDEN", Message: "Only admin or manager can audit knowledge"},
		})
		return
	}

	err := ctrl.agentService.AuditKnowledge(id, req.Status, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Knowledge audit complete"})
}

// ListKnowledges returns all agent-generated knowledge
func (ctrl *AgentController) ListKnowledges(c *gin.Context) {
	_, role, ok := requireAuth(c)
	if !ok {
		return
	}
	// Only admin or manager can list knowledge drafts
	if role != "admin" && role != "manager" {
		c.JSON(http.StatusForbidden, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "FORBIDDEN", Message: "Only admin or manager can list knowledge drafts"},
		})
		return
	}

	status := c.Query("status")
	query := c.Query("query")
	eqTypeIDStr := c.Query("equipment_type_id")
	var eqTypeID *uint
	if eqTypeIDStr != "" {
		if id, err := strconv.ParseUint(eqTypeIDStr, 10, 32); err == nil {
			val := uint(id)
			eqTypeID = &val
		}
	}

	result, err := ctrl.agentService.ListKnowledges(status, query, eqTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// =====================================================
// Phase 2: Chat & Conversation Endpoints
// =====================================================

// Chat handles multi-turn conversation
func (ctrl *AgentController) Chat(c *gin.Context) {
	var req dto.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[AgentController] Chat bind JSON error: %v", err)
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: err.Error()},
		})
		return
	}
	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}

	log.Printf("[AgentController] Chat request from User:%d, Message: %s", userID, req.Message)

	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: "Failed to get user context"},
		})
		return
	}

	result, err := ctrl.agentService.Chat(user, &req)
	if err != nil {
		log.Printf("[AgentController] Chat service error: %v", err)
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	
	log.Printf("[AgentController] Chat success, reply length: %d", len(result.Reply))
	c.JSON(http.StatusOK, result)
}

// ListConversations returns user's chat history
func (ctrl *AgentController) ListConversations(c *gin.Context) {
	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}
	result, err := ctrl.agentService.ListConversations(userID, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetConversation returns a single conversation with messages
func (ctrl *AgentController) GetConversation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: "Invalid ID"},
		})
		return
	}
	userID, role, ok := requireAuth(c)
	if !ok {
		return
	}
	result, err := ctrl.agentService.GetConversation(uint(id), userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// =====================================================
// Phase 2: Skill Management Endpoints
// =====================================================

// ListSkills returns all available agent skills
func (ctrl *AgentController) ListSkills(c *gin.Context) {
	_, role, ok := requireAuth(c)
	if !ok {
		return
	}
	// Only admin or manager can list skills (drafts)
	if role != "admin" && role != "manager" {
		c.JSON(http.StatusForbidden, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "FORBIDDEN", Message: "Only admin or manager can list skills"},
		})
		return
	}

	status := c.Query("status")
	query := c.Query("query")
	result, err := ctrl.agentService.ListSkills(status, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// CreateSkill creates a new skill
func (ctrl *AgentController) CreateSkill(c *gin.Context) {
	_, role, ok := requireAuth(c)
	if !ok {
		return
	}
	// Only admin or manager can create skills
	if role != "admin" && role != "manager" {
		c.JSON(http.StatusForbidden, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "FORBIDDEN", Message: "Only admin or manager can create skills"},
		})
		return
	}

	var req dto.CreateSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: err.Error()},
		})
		return
	}
	result, err := ctrl.agentService.CreateSkill(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusCreated, result)
}

// GetSkill returns skill detail
func (ctrl *AgentController) GetSkill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: "Invalid ID"},
		})
		return
	}
	_, role, ok := requireAuth(c)
	if !ok {
		return
	}
	// Only admin or manager can see full skill details
	if role != "admin" && role != "manager" {
		c.JSON(http.StatusForbidden, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "FORBIDDEN", Message: "Only admin or manager can access skill details"},
		})
		return
	}

	result, err := ctrl.agentService.GetSkill(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// UpdateSkill updates a skill
func (ctrl *AgentController) UpdateSkill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: "Invalid ID"},
		})
		return
	}
	_, role, ok := requireAuth(c)
	if !ok {
		return
	}
	// Only admin or manager can update skills
	if role != "admin" && role != "manager" {
		c.JSON(http.StatusForbidden, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "FORBIDDEN", Message: "Only admin or manager can update skills"},
		})
		return
	}

	var req dto.UpdateSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: err.Error()},
		})
		return
	}
	result, err := ctrl.agentService.UpdateSkill(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// =====================================================
// Phase 2: Proactive Notification Endpoints
// =====================================================

// Subscribe configures user's push notifications
func (ctrl *AgentController) Subscribe(c *gin.Context) {
	var req struct {
		PushType   string `json:"push_type" binding:"required"`
		Enabled    bool   `json:"enabled"`
		Scope      any    `json:"scope"`
		WebhookURL string `json:"webhook_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: err.Error()},
		})
		return
	}
	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}
	err := ctrl.agentService.Subscribe(userID, req.PushType, req.Enabled, req.Scope, req.WebhookURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscription updated"})
}

// ListSubscriptions returns user's push configurations
func (ctrl *AgentController) ListSubscriptions(c *gin.Context) {
	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}
	subs, err := ctrl.agentService.ListSubscriptions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, subs)
}

// GetEquipmentPrediction returns RUL and TCO for a specific equipment
func (ctrl *AgentController) GetEquipmentPrediction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: "Invalid ID"},
		})
		return
	}

	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}

	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: "User not found"},
		})
		return
	}

	result, err := ctrl.agentService.GetEquipmentPrediction(uint(id), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ListTools returns available tools
func (ctrl *AgentController) ListTools(c *gin.Context) {
	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}
	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: "User not found"},
		})
		return
	}

	tools, err := ctrl.agentService.ListTools(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.ListToolsResponse{Tools: tools})
}

// CallTool executes a tool call
func (ctrl *AgentController) CallTool(c *gin.Context) {
	var req dto.CallToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: err.Error()},
		})
		return
	}

	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}
	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: "User not found"},
		})
		return
	}

	scopes, _ := middleware.GetAPIKeyScopes(c)
	result, err := ctrl.agentService.CallTool(user, &req, scopes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error:   dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
