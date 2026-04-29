package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequireRole_AllowsMatchingRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(ContextKeyRole, "admin")
		c.Next()
	})
	router.Use(RequireRole("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRequireRole_RejectsNonMatchingRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(ContextKeyRole, "operator")
		c.Next()
	})
	router.Use(RequireRole("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestRequireRole_RejectsMissingRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	// Don't set any role in context
	router.Use(RequireRole("admin"))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetUserID_Exists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ContextKeyUserID, uint(42))

	uid, ok := GetUserID(c)
	if !ok {
		t.Error("expected ok=true")
	}
	if uid != 42 {
		t.Errorf("expected 42, got %d", uid)
	}
}

func TestGetUserID_NotExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	uid, ok := GetUserID(c)
	if ok {
		t.Error("expected ok=false")
	}
	if uid != 0 {
		t.Errorf("expected 0, got %d", uid)
	}
}

func TestGetUserRole_Exists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(ContextKeyRole, "admin")

	role, ok := GetUserRole(c)
	if !ok {
		t.Error("expected ok=true")
	}
	if role != "admin" {
		t.Errorf("expected 'admin', got '%s'", role)
	}
}
