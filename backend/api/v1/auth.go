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
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is disabled"})
		return
	}

	// Check approval status
	if user.ApprovalStatus != model.ApprovalStatusApproved {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is not approved yet", "status": string(user.ApprovalStatus)})
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
		Token:             token,
		ExpireAt:          time.Now().Add(24 * time.Hour).Unix(),
		UserInfo:          dto.UserInfo{
			ID:                user.ID,
			Username:          user.Username,
			Name:              user.Name,
			Role:              string(user.Role),
			FactoryID:         user.FactoryID,
			ApprovalStatus:    string(user.ApprovalStatus),
			MustChangePassword: user.MustChangePassword,
		},
		MustChangePassword: user.MustChangePassword,
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
		ID:                user.ID,
		Username:          user.Username,
		Name:              user.Name,
		Role:              string(user.Role),
		FactoryID:         user.FactoryID,
		ApprovalStatus:    string(user.ApprovalStatus),
		MustChangePassword: user.MustChangePassword,
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

// ChangePassword allows user to change their password
// @Summary Change password
// @Description Change the current user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.ChangePasswordRequest true "Password change request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/change-password [post]
func ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update password and reset flags
	user.PasswordHash = string(hashedPassword)
	user.MustChangePassword = false
	user.FirstLogin = false

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ApplyForAccount handles account application
// @Summary Apply for account
// @Description Submit an application for a new account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UserApplicationRequest true "Application details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/auth/apply [post]
func ApplyForAccount(c *gin.Context) {
	var req dto.UserApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var existingUser model.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user with pending status
	user := model.User{
		Username:           req.Username,
		PasswordHash:       string(hashedPassword),
		Name:               req.Name,
		Role:               model.UserRole(req.Role),
		FactoryID:          req.FactoryID,
		Phone:              req.Phone,
		IsActive:           true,
		ApprovalStatus:     model.ApprovalStatusPending,
		MustChangePassword: true,
		FirstLogin:         true,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Application submitted successfully", "id": user.ID})
}

// GetPendingApplications returns list of pending user applications (admin only)
// @Summary Get pending applications
// @Description Get list of pending user applications (admin only)
// @Tags users
// @Produce json
// @Security Bearer
// @Success 200 {array} dto.UserListResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/v1/users/applications [get]
func GetPendingApplications(c *gin.Context) {
	// Check if user is admin
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var currentUser model.User
	if err := db.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != model.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can view applications"})
		return
	}

	var users []model.User
	if err := db.Where("approval_status = ?", model.ApprovalStatusPending).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	response := make([]dto.UserListResponse, len(users))
	for i, user := range users {
		response[i] = dto.UserListResponse{
			ID:                user.ID,
			Username:          user.Username,
			Name:              user.Name,
			Role:              string(user.Role),
			Phone:             user.Phone,
			IsActive:          user.IsActive,
			ApprovalStatus:    string(user.ApprovalStatus),
			MustChangePassword: user.MustChangePassword,
			FactoryID:         user.FactoryID,
			CreatedAt:         user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, response)
}

// ApproveApplication approves or rejects a user application (admin only)
// @Summary Approve/reject application
// @Description Approve or reject a user application (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param request body dto.ApproveUserRequest true "Approval decision"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/{id}/approve [put]
func ApproveApplication(c *gin.Context) {
	// Check if user is admin
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var currentUser model.User
	if err := db.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != model.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can approve applications"})
		return
	}

	// Get user ID from path
	targetUserID := c.Param("id")
	var targetUser model.User
	if err := db.First(&targetUser, targetUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if user is in pending status
	if targetUser.ApprovalStatus != model.ApprovalStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User application is not in pending status"})
		return
	}

	var req dto.ApproveUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Approve {
		targetUser.ApprovalStatus = model.ApprovalStatusApproved
	} else {
		targetUser.ApprovalStatus = model.ApprovalStatusRejected
		targetUser.RejectionReason = req.Reason
	}

	if err := db.Save(&targetUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application"})
		return
	}

	status := "approved"
	if !req.Approve {
		status = "rejected"
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application " + status + " successfully"})
}

// GetUsers returns list of users (admin only)
// @Summary Get users
// @Description Get list of users (admin only)
// @Tags users
// @Produce json
// @Security Bearer
// @Success 200 {array} dto.UserListResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/v1/users [get]
func GetUsers(c *gin.Context) {
	// Check if user is admin
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var currentUser model.User
	if err := db.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != model.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can view users"})
		return
	}

	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	response := make([]dto.UserListResponse, len(users))
	for i, user := range users {
		response[i] = dto.UserListResponse{
			ID:                user.ID,
			Username:          user.Username,
			Name:              user.Name,
			Role:              string(user.Role),
			Phone:             user.Phone,
			IsActive:          user.IsActive,
			ApprovalStatus:    string(user.ApprovalStatus),
			MustChangePassword: user.MustChangePassword,
			FactoryID:         user.FactoryID,
			CreatedAt:         user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	c.JSON(http.StatusOK, response)
}

// CreateUser creates a new user (admin only)
// @Summary Create user
// @Description Create a new user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateUserRequest true "User details"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/v1/users [post]
func CreateUser(c *gin.Context) {
	// Check if user is admin
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var currentUser model.User
	if err := db.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != model.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can create users"})
		return
	}

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var existingUser model.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user with approved status
	user := model.User{
		Username:           req.Username,
		PasswordHash:       string(hashedPassword),
		Name:               req.Name,
		Role:               model.UserRole(req.Role),
		FactoryID:          req.FactoryID,
		Phone:              req.Phone,
		IsActive:           true,
		ApprovalStatus:     model.ApprovalStatusApproved,
		MustChangePassword: true, // Admin-created users also need to change password
		FirstLogin:         true,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "id": user.ID})
}

// UpdateUser updates a user (admin only)
// @Summary Update user
// @Description Update a user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserRequest true "User details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/{id} [put]
func UpdateUser(c *gin.Context) {
	// Check if user is admin
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var currentUser model.User
	if err := db.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if currentUser.Role != model.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can update users"})
		return
	}

	// Get user ID from path
	targetUserID := c.Param("id")
	var targetUser model.User
	if err := db.First(&targetUser, targetUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if req.Name != "" {
		targetUser.Name = req.Name
	}
	if req.Role != "" {
		targetUser.Role = model.UserRole(req.Role)
	}
	targetUser.FactoryID = req.FactoryID
	if req.Phone != "" {
		targetUser.Phone = req.Phone
	}
	if req.IsActive != nil {
		targetUser.IsActive = *req.IsActive
	}

	if err := db.Save(&targetUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
