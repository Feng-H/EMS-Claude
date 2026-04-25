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
// @Summary Recommend maintenance optimization
// @Tags agent
// @Router /agent/maintenance/recommend [post]
func (ctrl *AgentController) RecommendMaintenance(c *gin.Context) {
	var req dto.MaintenanceRecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error: dto.AgentErrDetail{
				Code:    "INVALID_ARGUMENT",
				Message: err.Error(),
			},
		})
		return
	}

	userID, _ := middleware.GetUserID(c)
	role, _ := middleware.GetUserRole(c)

	// Inject context
	result, err := ctrl.agentService.RecommendMaintenance(userID, role, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error: dto.AgentErrDetail{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AuditRepair audits repair reasonableness
// @Summary Audit repair
// @Tags agent
// @Router /agent/audit/repair [post]
func (ctrl *AgentController) AuditRepair(c *gin.Context) {
	var req dto.RepairAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error: dto.AgentErrDetail{
				Code:    "INVALID_ARGUMENT",
				Message: err.Error(),
			},
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
			Error: dto.AgentErrDetail{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AuditMaintenance audits maintenance plan density and quality
// @Summary Audit maintenance
// @Tags agent
// @Router /agent/audit/maintenance [post]
func (ctrl *AgentController) AuditMaintenance(c *gin.Context) {
	var req dto.MaintenanceAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error: dto.AgentErrDetail{
				Code:    "INVALID_ARGUMENT",
				Message: err.Error(),
			},
		})
		return
	}

	userID, _ := middleware.GetUserID(c)
	role, _ := middleware.GetUserRole(c)

	result, err := ctrl.agentService.AuditMaintenance(userID, role, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error: dto.AgentErrDetail{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Analyze answer structured management questions
// @Summary Analyze management questions
// @Tags agent
// @Router /agent/analyze [post]
func (ctrl *AgentController) Analyze(c *gin.Context) {
	var req dto.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error: dto.AgentErrDetail{
				Code:    "INVALID_ARGUMENT",
				Message: err.Error(),
			},
		})
		return
	}

	userID, _ := middleware.GetUserID(c)
	role, _ := middleware.GetUserRole(c)

	result, err := ctrl.agentService.Analyze(userID, role, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.AgentErrorEnvelope{
			Success: false,
			TraceID: trace.GenerateTraceID(),
			Error: dto.AgentErrDetail{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetSession returns session metadata
// @Summary Get agent session
// @Tags agent
// @Router /agent/sessions/:id [get]
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
// @Summary Get agent artifact
// @Tags agent
// @Router /agent/artifacts/:id [get]
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
