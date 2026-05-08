package analyzer

import (
	"testing"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/model"
)

func TestMaintenanceAnalyzer_Audit(t *testing.T) {
	// Note: In a full test, we would mock tools. For now, we test the basic logic flow.
	analyzer := NewMaintenanceAnalyzer(nil, nil)
	
	user := model.User{BaseModel: model.BaseModel{ID: 1}, Role: "engineer"}
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
	
	if len(data.Anomalies) == 0 {
		t.Error("Expected anomalies in audit report")
	}
	
	if data.PlanComparisons == nil {
		t.Error("Expected plan comparisons stats")
	}
}

func TestMaintenanceAnalyzer_Analyze(t *testing.T) {
	// Analyze currently requires tools to be non-nil for fetching plans
	// We'll skip deep testing until we have a mock repo/tool setup
}
