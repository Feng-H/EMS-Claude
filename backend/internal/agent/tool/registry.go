package tool

import (
	"context"
	"fmt"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/model"
)

// ToolFunc is the signature for a tool implementation
type ToolFunc func(user model.User, args map[string]interface{}) (interface{}, error)

// ToolEntry represents a registered tool with its metadata
type ToolEntry struct {
	Definition dto.ToolDefinition
	Handler    ToolFunc
	Scopes     []string // Required scopes for this tool
	IsReadOnly bool
}

// ToolRegistry manages a collection of agent tools
type ToolRegistry struct {
	tools map[string]ToolEntry
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]ToolEntry),
	}
}

func (r *ToolRegistry) Register(name string, def dto.ToolDefinition, handler ToolFunc, scopes []string, isReadOnly bool) {
	r.tools[name] = ToolEntry{
		Definition: def,
		Handler:    handler,
		Scopes:     scopes,
		IsReadOnly: isReadOnly,
	}
}

func (r *ToolRegistry) List(user model.User) []dto.ToolDefinition {
	var defs []dto.ToolDefinition
	for _, t := range r.tools {
		// In a more advanced implementation, filter by user scopes here
		defs = append(defs, t.Definition)
	}
	return defs
}

func (r *ToolRegistry) Call(name string, user model.User, args map[string]interface{}, userScopes []string) (interface{}, error) {
	entry, ok := r.tools[name]
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}

	// Scope check
	if !entry.IsReadOnly {
		// Basic scope check for write operations
		// In production, check if userScopes contains at least one of entry.Scopes
	}

	return entry.Handler(user, args)
}

func (r *ToolRegistry) GetTool(name string) (ToolEntry, bool) {
	t, ok := r.tools[name]
	return t, ok
}
