package analyzer

import (
	"testing"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/agent/tool"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
)

func TestMaintenanceAnalyzer_Audit(t *testing.T) {
	config.Cfg = &config.Config{
		Storage: config.StorageConfig{Mode: "memory"},
	}

	analyzer := NewMaintenanceAnalyzer(nil, tool.NewMaintenanceTool())

	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "admin"}
	req := &dto.MaintenanceAuditRequest{
		FactoryID: 1,
		TimeRange: dto.TimeRange{StartDate: "2026-05-01"},
	}

	data, err := analyzer.Audit(req, user)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if data == nil {
		t.Fatal("Expected non-nil data")
	}

	if data.PlanComparisons == nil {
		t.Error("Expected plan comparisons stats")
	}
}

func TestMaintenanceAnalyzer_Analyze(t *testing.T) {
	// Analyze currently requires tools to be non-nil for fetching plans
	// We'll skip deep testing until we have a mock repo/tool setup
}
