package v1

import (
	"net/http"
	"time"

	"github.com/ems/backend/internal/dto"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/pkg/jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitAuth(database *gorm.DB) {
	db = database
}

// Login handles user login
// @Summary User login
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := db.Where("username = ? AND is_active = ?", req.Username, true).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := jwt.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Token:    token,
		ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
		UserInfo: dto.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Name:      user.Name,
			Role:      string(user.Role),
			FactoryID: user.FactoryID,
		},
	})
}

// GetCurrentUser returns the current authenticated user info
// @Summary Get current user
// @Description Get information about the currently authenticated user
// @Tags auth
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.UserInfo
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/me [get]
func GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, dto.UserInfo{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Role:      string(user.Role),
		FactoryID: user.FactoryID,
	})
}

// RefreshToken refreshes an expired token
// @Summary Refresh token
// @Description Refresh an authentication token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Token to refresh"
// @Success 200 {object} dto.RefreshTokenResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.RefreshToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.RefreshTokenResponse{
		Token:    token,
		ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}
