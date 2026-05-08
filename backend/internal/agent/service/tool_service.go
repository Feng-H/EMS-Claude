package service

import (
	"fmt"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/database"
)

// ListTools returns the list of available tools for external Agents
func (s *AgentService) ListTools(user model.User) ([]dto.ToolDefinition, error) {
	return s.toolRegistry.List(user), nil
}

// CallTool executes a tool call for an external Agent
func (s *AgentService) CallTool(user model.User, req *dto.CallToolRequest) (*dto.CallToolResponse, error) {
	// Retrieve user's API Key scopes if applicable
	// For simplicity, we'll pass nil for now.
	result, err := s.toolRegistry.Call(req.Name, user, req.Arguments, nil)
	if err != nil {
		return &dto.CallToolResponse{Content: err.Error(), IsError: true}, nil
	}
	return &dto.CallToolResponse{Content: result, IsError: false}, nil
}
