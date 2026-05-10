package trace

import (
	"strings"
	"testing"
)

func TestGenerateTraceID_Format(t *testing.T) {
	id := GenerateTraceID()

	if !strings.HasPrefix(id, "agt_") {
		t.Errorf("Expected prefix 'agt_', got %s", id)
	}

	parts := strings.Split(id, "_")
	if len(parts) != 3 {
		t.Fatalf("Expected 3 parts separated by '_', got %d: %s", len(parts), id)
	}

	// Date part should be 8 digits (YYYYMMDD)
	if len(parts[1]) != 8 {
		t.Errorf("Expected date part length 8, got %d: %s", len(parts[1]), parts[1])
	}

	// Random part should be 6 digits
	if len(parts[2]) != 6 {
		t.Errorf("Expected random part length 6, got %d: %s", len(parts[2]), parts[2])
	}
}

func TestGenerateTraceID_Unique(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := GenerateTraceID()
		if ids[id] {
			t.Errorf("Duplicate trace ID generated: %s", id)
		}
		ids[id] = true
	}
}

func TestGenerateTraceID_NonEmpty(t *testing.T) {
	id := GenerateTraceID()
	if id == "" {
		t.Error("Expected non-empty trace ID")
	}
	if len(id) < 15 {
		t.Errorf("Trace ID seems too short: %s (len=%d)", id, len(id))
	}
}
