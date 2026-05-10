package qrcode

import (
	"strings"
	"testing"

	"github.com/ems/backend/pkg/config"
)

func TestValidateQRCode_ShortCodes(t *testing.T) {
	if ValidateQRCode("ab") {
		t.Error("Expected 2-char code to be invalid (too short)")
	}
	if ValidateQRCode("") {
		t.Error("Expected empty code to be invalid")
	}
	if ValidateQRCode("abc") {
		t.Error("Expected 3-char code to be invalid (too short)")
	}
}

func TestValidateQRCode_ValidURL(t *testing.T) {
	if !ValidateQRCode("https://ems.example.com/equipment/CNC-001") {
		t.Error("Expected valid URL to pass validation")
	}
	if !ValidateQRCode("http://localhost/eq/ABC") {
		t.Error("Expected valid HTTP URL to pass validation")
	}
}

func TestValidateQRCode_TooLong(t *testing.T) {
	longCode := strings.Repeat("x", 256)
	if ValidateQRCode(longCode) {
		t.Error("Expected 256-char code to be invalid (too long)")
	}
}

func TestValidateQRCode_InvalidURL(t *testing.T) {
	if ValidateQRCode("    ") {
		t.Error("Expected whitespace-only string to fail URL parsing")
	}
}

func TestGetEquipmentURL(t *testing.T) {
	config.Cfg = &config.Config{
		App: config.AppConfig{
			QRCodeBaseURL: "https://ems.example.com",
		},
	}

	url := GetEquipmentURL("CNC-001")
	expected := "https://ems.example.com/equipment/CNC-001"
	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}

func TestGetEquipmentURL_EmptyBase(t *testing.T) {
	config.Cfg = &config.Config{
		App: config.AppConfig{
			QRCodeBaseURL: "",
		},
	}

	url := GetEquipmentURL("CNC-001")
	expected := "/equipment/CNC-001"
	if url != expected {
		t.Errorf("Expected %s, got %s", expected, url)
	}
}
