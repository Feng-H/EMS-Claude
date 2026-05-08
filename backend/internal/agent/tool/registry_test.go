package tool

import (
	"testing"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/model"
)

func TestToolRegistry_RegisterAndList(t *testing.T) {
	registry := NewToolRegistry()
	
	def := dto.ToolDefinition{
		Name: "test_tool",
		Description: "A test tool",
	}
	
	handler := func(user model.User, args map[string]interface{}) (interface{}, error) {
		return "ok", nil
	}
	
	registry.Register("test_tool", def, handler, []string{"test_scope"}, true)
	
	tools := registry.List(model.User{})
	if len(tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}
	
	if tools[0].Name != "test_tool" {
		t.Errorf("Expected tool name 'test_tool', got '%s'", tools[0].Name)
	}
}

func TestToolRegistry_Call(t *testing.T) {
	registry := NewToolRegistry()
	
	def := dto.ToolDefinition{Name: "test_tool"}
	handler := func(user model.User, args map[string]interface{}) (interface{}, error) {
		return args["input"], nil
	}
	
	registry.Register("test_tool", def, handler, []string{"read"}, true)
	
	user := model.User{BaseModel: model.BaseModel{ID: 1}}
	args := map[string]interface{}{"input": "hello"}
	
	result, err := registry.Call("test_tool", user, args, []string{"read"})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if result != "hello" {
		t.Errorf("Expected 'hello', got %v", result)
	}
}

func TestToolRegistry_ToolNotFound(t *testing.T) {
	registry := NewToolRegistry()
	
	_, err := registry.Call("non_existent", model.User{}, nil, nil)
	if err == nil {
		t.Error("Expected error for non-existent tool, got nil")
	}
}
