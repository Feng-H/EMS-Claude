package redis

import (
	"testing"
)

func TestBuildUserKey(t *testing.T) {
	key := BuildUserKey(42)
	expected := "ems:user:42"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
}

func TestBuildEquipmentKey(t *testing.T) {
	key := BuildEquipmentKey(100)
	expected := "ems:equipment:100"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
}

func TestBuildEquipmentQRKey(t *testing.T) {
	key := BuildEquipmentQRKey("CNC-001")
	expected := "ems:equipment:qr:CNC-001"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
}

func TestBuildEquipmentTypeKey(t *testing.T) {
	key := BuildEquipmentTypeKey(5)
	expected := "ems:equipment_type:5"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
}

func TestBuildWorkshopKey(t *testing.T) {
	key := BuildWorkshopKey(10)
	expected := "ems:workshop:10"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
}

func TestBuildFactoryKey(t *testing.T) {
	key := BuildFactoryKey(3)
	expected := "ems:factory:3"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
}

func TestBuildInspectionTemplateKey(t *testing.T) {
	key := BuildInspectionTemplateKey(7)
	expected := "ems:inspection_template:7"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
}

func TestBuildStatisticsKey(t *testing.T) {
	key := BuildStatisticsKey("mttr")
	expected := "ems:stats:mttr"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
}

func TestBuildKeys_ZeroID(t *testing.T) {
	if BuildUserKey(0) != "ems:user:0" {
		t.Error("Should handle zero ID")
	}
}
