package tool

import (
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
	if len(entry.Scopes) > 0 {
		hasScope := false
		scopeMap := make(map[string]bool)
		for _, s := range userScopes {
			scopeMap[s] = true
		}

		for _, required := range entry.Scopes {
			if scopeMap[required] {
				hasScope = true
				break
			}
		}

		if !hasScope && len(userScopes) > 0 {
			// If user has scopes (API Key), but none match
			return nil, fmt.Errorf("permission denied: missing required scope(s) %v", entry.Scopes)
		}
		// If userScopes is empty, it might be a Web user (JWT), 
		// who relies on Role/Factory checks inside the tool handler.
		// We could strictly enforce that API Key users MUST have scopes.
	}

	return entry.Handler(user, args)
}

func (r *ToolRegistry) GetTool(name string) (ToolEntry, bool) {
	t, ok := r.tools[name]
	return t, ok
}
