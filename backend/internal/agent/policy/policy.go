package policy

import (
	"fmt"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
	"github.com/ems/backend/pkg/memory"
	"github.com/ems/backend/pkg/config"
)

type AgentContext struct {
	UserID     uint
	Role       string
	FactoryID  *uint
	WorkshopID *uint
	Language   string
}

type PolicyService struct {
}

func NewPolicyService() *PolicyService {
	return &PolicyService{}
}

// DeriveAgentContext derives the authorized scope for an agent request
func (s *PolicyService) DeriveAgentContext(userID uint, role string, language string) (*AgentContext, error) {
	ctx := &AgentContext{
		UserID:   userID,
		Role:     role,
		Language: language,
	}
	
	if language == "" {
		ctx.Language = "zh-CN"
	}

	// Fetch user to get factory_id
	var user *model.User
	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		user = store.FindUser(userID)
	} else {
		userRepo := repository.NewUserRepository()
		user, _ = userRepo.GetByID(userID)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	ctx.FactoryID = user.FactoryID
	
	return ctx, nil
}

// ValidateScope checks if requested identifiers are within the authorized context
func (s *PolicyService) ValidateScope(ctx *AgentContext, requestedFactoryID *uint) error {
	// Admin can see everything
	if ctx.Role == string(model.RoleAdmin) {
		return nil
	}

	// Others are restricted to their factory
	if requestedFactoryID != nil && ctx.FactoryID != nil && *requestedFactoryID != *ctx.FactoryID {
		return fmt.Errorf("FORBIDDEN_SCOPE: access to other factory is denied")
	}
	
	// If user has no factory_id but is not admin, they might be restricted
	if ctx.Role != string(model.RoleAdmin) && ctx.FactoryID == nil {
		// This depends on business logic, for now assume they can't see anything
		return fmt.Errorf("FORBIDDEN_SCOPE: user has no authorized factory")
	}

	return nil
}
