package v1

import (
	"net/http"
	"time"

	"github.com/ems/backend/internal/dto"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/pkg/jwt"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/memory"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitAuth(database *gorm.DB) {
	db = database
}

// Login handles user login
func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		u := store.FindUserByUsername(req.Username)
		if u == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		user = *u
	} else {
		if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is disabled"})
		return
	}

	if user.ApprovalStatus != model.ApprovalStatusApproved {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is not approved yet", "status": string(user.ApprovalStatus)})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
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

// GetCurrentUser returns current user info
func GetCurrentUser(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var user model.User
	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		u := store.FindUser(userID)
		if u == nil {
			c.JSON(404, gin.H{"error": "用户不存在"})
			return
		}
		user = *u
	} else {
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(404, gin.H{"error": "用户不存在"})
			return
		}
	}
	c.JSON(200, dto.UserInfo{
		ID: user.ID, Username: user.Username, Name: user.Name, Role: string(user.Role),
		FactoryID: user.FactoryID, ApprovalStatus: string(user.ApprovalStatus), MustChangePassword: user.MustChangePassword,
	})
}

// RefreshToken refreshes an expired token
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
		Token: token, ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
	})
}

// ChangePassword allows user to change password
func ChangePassword(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if config.Cfg.Storage.Mode == "memory" {
		u := memory.GetStore().FindUser(userID)
		if u == nil { c.JSON(404, gin.H{"error": "用户不存在"}); return }
		user = *u
	} else {
		if err := db.First(&user, userID).Error; err != nil { c.JSON(404, gin.H{"error": "用户不存在"}); return }
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		c.JSON(400, gin.H{"error": "原密码错误"})
		return
	}

	h, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if config.Cfg.Storage.Mode == "memory" {
		memory.GetStore().UpdateUser(userID, func(u *model.User) {
			u.PasswordHash = string(h)
			u.MustChangePassword = false
			u.FirstLogin = false
		})
	} else {
		user.PasswordHash = string(h)
		user.MustChangePassword = false
		user.FirstLogin = false
		db.Save(&user)
	}
	c.JSON(200, gin.H{"message": "密码修改成功"})
}

// ApplyForAccount handles account application
func ApplyForAccount(c *gin.Context) {
	var req dto.UserApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := model.User{
		Username: req.Username, PasswordHash: string(h), Name: req.Name,
		Role: model.UserRole(req.Role), FactoryID: req.FactoryID, Phone: req.Phone,
		IsActive: true, ApprovalStatus: model.ApprovalStatusPending, MustChangePassword: true, FirstLogin: true,
	}
	if config.Cfg.Storage.Mode == "memory" {
		memory.GetStore().AddUser(memory.GetStore().NextID(), &user)
	} else {
		db.Create(&user)
	}
	c.JSON(201, gin.H{"message": "申请提交成功"})
}

// GetPendingApplications (Memory proxy in main.go)
func GetPendingApplications(c *gin.Context) {
	var users []model.User
	db.Where("approval_status = ?", model.ApprovalStatusPending).Find(&users)
	res := make([]dto.UserListResponse, len(users))
	for i, u := range users {
		res[i] = dto.UserListResponse{
			ID: u.ID, Username: u.Username, Name: u.Name, Role: string(u.Role),
			Phone: u.Phone, IsActive: u.IsActive, ApprovalStatus: string(u.ApprovalStatus),
			FactoryID: u.FactoryID, CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	c.JSON(200, res)
}

// ApproveApplication (Memory proxy in main.go)
func ApproveApplication(c *gin.Context) {
	targetID := c.Param("id")
	var req dto.ApproveUserRequest
	c.ShouldBindJSON(&req)
	var u model.User
	if db.First(&u, targetID).Error != nil { c.JSON(404, gin.H{"error": "用户不存在"}); return }
	if req.Approve { u.ApprovalStatus = model.ApprovalStatusApproved } else { u.ApprovalStatus = model.ApprovalStatusRejected }
	db.Save(&u)
	c.JSON(200, gin.H{"message": "审核完成"})
}

// GetUsers (Memory proxy in main.go)
func GetUsers(c *gin.Context) {
	var users []model.User
	db.Find(&users)
	res := make([]dto.UserListResponse, len(users))
	for i, u := range users {
		res[i] = dto.UserListResponse{
			ID: u.ID, Username: u.Username, Name: u.Name, Role: string(u.Role),
			Phone: u.Phone, IsActive: u.IsActive, ApprovalStatus: string(u.ApprovalStatus),
			FactoryID: u.FactoryID, CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	c.JSON(200, res)
}

// CreateUser (Memory proxy in main.go)
func CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	c.ShouldBindJSON(&req)
	h, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	u := model.User{
		Username: req.Username, PasswordHash: string(h), Name: req.Name,
		Role: model.UserRole(req.Role), FactoryID: req.FactoryID, Phone: req.Phone,
		IsActive: true, ApprovalStatus: model.ApprovalStatusApproved, MustChangePassword: true, FirstLogin: true,
	}
	db.Create(&u)
	c.JSON(201, gin.H{"message": "用户创建成功"})
}

// UpdateUser (Memory proxy in main.go)
func UpdateUser(c *gin.Context) {
	targetID := c.Param("id")
	var req dto.UpdateUserRequest
	c.ShouldBindJSON(&req)
	var u model.User
	if db.First(&u, targetID).Error != nil { c.JSON(404, gin.H{"error": "用户不存在"}); return }
	if req.Name != "" { u.Name = req.Name }
	if req.Role != "" { u.Role = model.UserRole(req.Role) }
	u.FactoryID = req.FactoryID
	db.Save(&u)
	c.JSON(200, gin.H{"message": "用户更新成功"})
}
