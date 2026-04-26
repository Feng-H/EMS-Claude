package controller

import (
	"net/http"
	"strconv"

	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/service"
	"github.com/ems/backend/internal/middleware"
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

// RecommendMaintenance generates maintenance optimization recommendations
func (ctrl *AgentController) RecommendMaintenance(c *gin.Context) {
	var req dto.MaintenanceRecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error: dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: err.Error()},
		})
		return
	}
	userID, _ := middleware.GetUserID(c)
	role, _ := middleware.GetUserRole(c)
	result, err := ctrl.agentService.RecommendMaintenance(userID, role, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error: dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
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
			Error: dto.AgentErrDetail{Code: "INVALID_ARGUMENT", Message: err.Error()},
		})
		return
	}
	userID, _ := middleware.GetUserID(c)
	role, _ := middleware.GetUserRole(c)
	result, err := ctrl.agentService.AuditRepair(userID, role, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error: dto.AgentErrDetail{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

// AuditMaintenance audits maintenance plan
func (ctrl *AgentController) AuditMaintenance(c *gin.Context) {
	var req dto.MaintenanceAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := middleware.GetUserID(c)
	role, _ := middleware.GetUserRole(c)
	result, err := ctrl.agentService.AuditMaintenance(userID, role, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Analyze answer management questions
func (ctrl *AgentController) Analyze(c *gin.Context) {
	var req dto.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := middleware.GetUserID(c)
	role, _ := middleware.GetUserRole(c)
	result, err := ctrl.agentService.Analyze(userID, role, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetSession returns session metadata
func (ctrl *AgentController) GetSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result, err := ctrl.agentService.GetSession(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetArtifact returns one artifact
func (ctrl *AgentController) GetArtifact(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result, err := ctrl.agentService.GetArtifact(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ListSessions returns user's recent agent sessions
func (ctrl *AgentController) ListSessions(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	result, err := ctrl.agentService.ListSessions(userID, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := middleware.GetUserID(c)
	err := ctrl.agentService.AuditKnowledge(id, req.Status, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Knowledge audit complete"})
}

// =====================================================
// Phase 2: Chat & Conversation Endpoints
// =====================================================

// Chat handles multi-turn conversation
func (ctrl *AgentController) Chat(c *gin.Context) {
	var req dto.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := middleware.GetUserID(c)
	role, _ := middleware.GetUserRole(c)
	result, err := ctrl.agentService.Chat(userID, role, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// ListConversations returns user's chat history
func (ctrl *AgentController) ListConversations(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	result, err := ctrl.agentService.ListConversations(userID, 20)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetConversation returns a single conversation with messages
func (ctrl *AgentController) GetConversation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result, err := ctrl.agentService.GetConversation(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// =====================================================
// Phase 2: Skill Management Endpoints
// =====================================================

// ListSkills returns all available agent skills
func (ctrl *AgentController) ListSkills(c *gin.Context) {
	status := c.Query("status")
	result, err := ctrl.agentService.ListSkills(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// CreateSkill creates a new skill
func (ctrl *AgentController) CreateSkill(c *gin.Context) {
	var req dto.CreateSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := ctrl.agentService.CreateSkill(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

// GetSkill returns skill detail
func (ctrl *AgentController) GetSkill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	result, err := ctrl.agentService.GetSkill(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// UpdateSkill updates skill status
func (ctrl *AgentController) UpdateSkill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var req dto.UpdateSkillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := ctrl.agentService.UpdateSkill(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		PushType string `json:"push_type" binding:"required"`
		Enabled  bool   `json:"enabled"`
		Scope    any    `json:"scope"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := middleware.GetUserID(c)
	err := ctrl.agentService.Subscribe(userID, req.PushType, req.Enabled, req.Scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscription updated"})
}

// GetEquipmentPrediction returns RUL and TCO for a specific equipment
func (ctrl *AgentController) GetEquipmentPrediction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := ctrl.agentService.GetEquipmentPrediction(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
