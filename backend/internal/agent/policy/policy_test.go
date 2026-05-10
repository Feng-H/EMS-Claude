package policy

import (
	"testing"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/memory"
)

func setupPolicyTest() {
	config.Cfg = &config.Config{
		Storage: config.StorageConfig{Mode: "memory"},
	}
	store := memory.GetStore()
	factoryID := uint(10)
	store.Users[1] = &model.User{
		BaseModel:  model.BaseModel{ID: 1},
		Role:       model.RoleAdmin,
		FactoryID:  &factoryID,
	}
	store.Users[2] = &model.User{
		BaseModel:  model.BaseModel{ID: 2},
		Role:       model.RoleEngineer,
		FactoryID:  &factoryID,
	}
	store.Users[3] = &model.User{
		BaseModel: model.BaseModel{ID: 3},
		Role:      model.RoleOperator,
	}
}

func TestDeriveAgentContext_DefaultLanguage(t *testing.T) {
	setupPolicyTest()
	svc := NewPolicyService()

	ctx, err := svc.DeriveAgentContext(1, "admin", "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if ctx.Language != "zh-CN" {
		t.Errorf("Expected default language zh-CN, got %s", ctx.Language)
	}
}

func TestDeriveAgentContext_CustomLanguage(t *testing.T) {
	setupPolicyTest()
	svc := NewPolicyService()

	ctx, err := svc.DeriveAgentContext(1, "admin", "en-US")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if ctx.Language != "en-US" {
		t.Errorf("Expected language en-US, got %s", ctx.Language)
	}
}

func TestDeriveAgentContext_FactoryID(t *testing.T) {
	setupPolicyTest()
	svc := NewPolicyService()

	ctx, err := svc.DeriveAgentContext(2, "engineer", "")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if ctx.FactoryID == nil || *ctx.FactoryID != 10 {
		t.Errorf("Expected FactoryID 10, got %v", ctx.FactoryID)
	}
}

func TestDeriveAgentContext_UserNotFound(t *testing.T) {
	setupPolicyTest()
	svc := NewPolicyService()

	_, err := svc.DeriveAgentContext(999, "admin", "")
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
}

func TestValidateScope_AdminBypass(t *testing.T) {
	svc := NewPolicyService()
	factoryID := uint(99)
	ctx := &AgentContext{
		UserID:    1,
		Role:      "admin",
		FactoryID: &factoryID,
	}

	otherFactory := uint(50)
	err := svc.ValidateScope(ctx, &otherFactory)
	if err != nil {
		t.Errorf("Admin should bypass scope check, got: %v", err)
	}
}

func TestValidateScope_SameFactory(t *testing.T) {
	svc := NewPolicyService()
	factoryID := uint(10)
	ctx := &AgentContext{
		UserID:    2,
		Role:      "engineer",
		FactoryID: &factoryID,
	}

	err := svc.ValidateScope(ctx, &factoryID)
	if err != nil {
		t.Errorf("Same factory should be allowed, got: %v", err)
	}
}

func TestValidateScope_DifferentFactory(t *testing.T) {
	svc := NewPolicyService()
	factoryID := uint(10)
	ctx := &AgentContext{
		UserID:    2,
		Role:      "engineer",
		FactoryID: &factoryID,
	}

	otherFactory := uint(20)
	err := svc.ValidateScope(ctx, &otherFactory)
	if err == nil {
		t.Error("Expected error for different factory access")
	}
}

func TestValidateScope_NoFactory(t *testing.T) {
	svc := NewPolicyService()
	ctx := &AgentContext{
		UserID:    3,
		Role:      "operator",
		FactoryID: nil,
	}

	err := svc.ValidateScope(ctx, nil)
	if err == nil {
		t.Error("Expected error for user with no factory")
	}
}

func TestValidateScope_NilRequested(t *testing.T) {
	svc := NewPolicyService()
	factoryID := uint(10)
	ctx := &AgentContext{
		UserID:    2,
		Role:      "engineer",
		FactoryID: &factoryID,
	}

	err := svc.ValidateScope(ctx, nil)
	if err != nil {
		t.Errorf("Nil requested factory should be allowed, got: %v", err)
	}
}
