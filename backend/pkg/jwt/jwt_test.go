package jwt

import (
	"testing"

	"github.com/ems/backend/pkg/config"
)

func setupTestJWT() {
	config.Cfg = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret-key-for-unit-tests",
			Issuer:      "ems-test",
			ExpireHours: 24,
		},
	}
}

func TestGenerateAndParseToken(t *testing.T) {
	setupTestJWT()

	tokenStr, err := GenerateToken(42, "admin")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("Expected non-empty token")
	}

	claims, err := ParseToken(tokenStr)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("Expected UserID 42, got %d", claims.UserID)
	}
	if claims.Role != "admin" {
		t.Errorf("Expected Role admin, got %s", claims.Role)
	}
	if claims.Issuer != "ems-test" {
		t.Errorf("Expected Issuer ems-test, got %s", claims.Issuer)
	}
}

func TestParseToken_Invalid(t *testing.T) {
	setupTestJWT()

	_, err := ParseToken("invalid.token.string")
	if err != ErrTokenInvalid {
		t.Errorf("Expected ErrTokenInvalid, got %v", err)
	}
}

func TestParseToken_Expired(t *testing.T) {
	config.Cfg = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			Issuer:      "ems-test",
			ExpireHours: -1,
		},
	}

	tokenStr, err := GenerateToken(1, "operator")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = ParseToken(tokenStr)
	if err != ErrTokenExpired {
		t.Errorf("Expected ErrTokenExpired, got %v", err)
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	setupTestJWT()
	tokenStr, _ := GenerateToken(1, "admin")

	config.Cfg.JWT.Secret = "wrong-secret"
	_, err := ParseToken(tokenStr)
	if err == nil {
		t.Error("Expected error for wrong secret, got nil")
	}
}

func TestRefreshToken_TooEarly(t *testing.T) {
	setupTestJWT()
	tokenStr, _ := GenerateToken(1, "admin")

	_, err := RefreshToken(tokenStr)
	if err == nil {
		t.Error("Expected error when token is not close to expiry")
	}
}

func TestRefreshToken_WithinWindow(t *testing.T) {
	config.Cfg = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			Issuer:      "ems-test",
			ExpireHours: 1,
		},
	}
	// Token expires in ~1 hour, which is within the 1-hour refresh window
	tokenStr, _ := GenerateToken(1, "engineer")

	newToken, err := RefreshToken(tokenStr)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}
	if newToken == "" {
		t.Fatal("Expected non-empty refreshed token")
	}

	claims, err := ParseToken(newToken)
	if err != nil {
		t.Fatalf("ParseToken on refreshed token failed: %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", claims.UserID)
	}
}

func TestRefreshToken_ExpiredToken(t *testing.T) {
	config.Cfg = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			Issuer:      "ems-test",
			ExpireHours: -1,
		},
	}
	tokenStr, _ := GenerateToken(1, "admin")

	_, err := RefreshToken(tokenStr)
	if err == nil {
		t.Error("Expected error when refreshing expired token")
	}
}
