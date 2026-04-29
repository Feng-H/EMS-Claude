package config

import (
	"os"
	"testing"
)

func TestValidate_RejectsPlaceholderInRelease(t *testing.T) {
	Cfg = &Config{
		Server: ServerConfig{Mode: "release"},
		JWT:    JWTConfig{Secret: "your-secret-key-change-in-production"},
	}
	err := Validate()
	if err == nil {
		t.Error("expected error for placeholder secret in release mode, got nil")
	}
}

func TestValidate_AllowsCustomSecret(t *testing.T) {
	Cfg = &Config{
		Server: ServerConfig{Mode: "release"},
		JWT:    JWTConfig{Secret: "my-strong-random-secret-key-123456"},
	}
	err := Validate()
	if err != nil {
		t.Errorf("expected nil error for custom secret, got: %v", err)
	}
}

func TestValidate_WarnsInDebugMode(t *testing.T) {
	Cfg = &Config{
		Server: ServerConfig{Mode: "debug"},
		JWT:    JWTConfig{Secret: "your-secret-key-change-in-production"},
	}
	err := Validate()
	if err != nil {
		t.Errorf("expected nil error in debug mode (warning only), got: %v", err)
	}
}

func TestValidate_RejectsEmptySecretInRelease(t *testing.T) {
	Cfg = &Config{
		Server: ServerConfig{Mode: "release"},
		JWT:    JWTConfig{Secret: ""},
	}
	err := Validate()
	if err == nil {
		t.Error("expected error for empty secret in release mode, got nil")
	}
}

func TestValidate_NilConfig(t *testing.T) {
	origCfg := Cfg
	Cfg = nil
	defer func() { Cfg = origCfg }()

	err := Validate()
	if err == nil {
		t.Error("expected error for nil config, got nil")
	}
}

func TestValidate_SecretFromEnv(t *testing.T) {
	os.Setenv("EMS_JWT_SECRET", "env-secret-key")
	defer os.Unsetenv("EMS_JWT_SECRET")

	// Simulate what Load + env override does
	Cfg = &Config{
		Server: ServerConfig{Mode: "release"},
		JWT:    JWTConfig{Secret: "your-secret-key-change-in-production"},
	}
	applyEnvOverrides(Cfg)

	err := Validate()
	if err != nil {
		t.Errorf("expected nil error when secret set via env, got: %v", err)
	}
	if Cfg.JWT.Secret != "env-secret-key" {
		t.Errorf("expected secret from env, got: %s", Cfg.JWT.Secret)
	}
}
