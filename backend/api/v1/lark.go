package v1

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ems/backend/internal/dto"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/service"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/memory"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	larkService *service.LarkService
)

func InitLark(database *gorm.DB) {
	// db is already declared in auth.go (same package v1)
	db = database
	larkService = service.NewLarkService()
}

// GetLarkConfig retrieves user's bot config (masked)
func GetLarkConfig(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user model.User
	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		u := store.FindUser(userID)
		if u == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		user = *u
	} else {
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}
	}

	appID, token := "", ""
	if user.LarkAppID != nil {
		appID = *user.LarkAppID
	}
	if user.LarkVerificationToken != nil {
		token = *user.LarkVerificationToken
	}

	// Use domain from config
	baseURL := config.Cfg.App.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:8080" // fallback
	}
	webhookURL := fmt.Sprintf("%s/api/v1/lark/webhook/%d", baseURL, userID)

	c.JSON(http.StatusOK, dto.LarkConfigResp{
		AppID:             appID,
		HasAppSecret:      user.LarkAppSecret != nil && *user.LarkAppSecret != "",
		VerificationToken: token,
		HasEncryptKey:     user.LarkEncryptKey != nil && *user.LarkEncryptKey != "",
		WebhookURL:        webhookURL,
	})
}

// UpdateLarkConfig updates user's bot config
func UpdateLarkConfig(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.LarkConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		u := store.FindUser(userID)
		if u == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if req.AppID != "" { u.LarkAppID = &req.AppID }
		if req.AppSecret != "" { u.LarkAppSecret = &req.AppSecret }
		if req.VerificationToken != "" { u.LarkVerificationToken = &req.VerificationToken }
		if req.EncryptKey != "" { u.LarkEncryptKey = &req.EncryptKey }
		
		c.JSON(http.StatusOK, gin.H{"message": "success"})
		return
	}

	updates := make(map[string]interface{})
	if req.AppID != "" {
		updates["LarkAppID"] = req.AppID
	}
	if req.AppSecret != "" {
		updates["LarkAppSecret"] = req.AppSecret
	}
	if req.VerificationToken != "" {
		updates["LarkVerificationToken"] = req.VerificationToken
	}
	if req.EncryptKey != "" {
		updates["LarkEncryptKey"] = req.EncryptKey
	}

	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no changes"})
		return
	}

	if err := db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// verifyLarkSignature ensures the request is coming from Lark
func verifyLarkSignature(c *gin.Context, body []byte, encryptKey string) bool {
	timestamp := c.GetHeader("X-Lark-Request-Timestamp")
	nonce := c.GetHeader("X-Lark-Request-Nonce")
	signature := c.GetHeader("X-Lark-Signature")

	if encryptKey == "" {
		if signature != "" {
			return false
		}
		return true
	}

	if signature == "" {
		return false
	}

	// V2 signature: sha256(timestamp + "\n" + nonce + "\n" + encrypt_key + "\n" + body)
	content := timestamp + "\n" + nonce + "\n" + encryptKey + "\n" + string(body)
	h := sha256.New()
	h.Write([]byte(content))
	sum := fmt.Sprintf("%x", h.Sum(nil))

	return sum == signature
}

// LarkWebhook handles Lark events
func LarkWebhook(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var user model.User
	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		u := store.FindUser(uint(userID))
		if u == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		user = *u
	} else {
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	var req dto.LarkWebhookRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Handle URL Challenge FIRST
	if req.Type == "url_verification" || req.Header.EventType == "url_verification" {
		token := req.Token
		if token == "" {
			token = req.Header.Token
		}
		
		expectedToken := ""
		if user.LarkVerificationToken != nil {
			expectedToken = *user.LarkVerificationToken
		}

		if token != expectedToken {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}
		challenge := req.Challenge
		if challenge == "" && req.Event != nil {
			if ch, ok := req.Event["challenge"].(string); ok {
				challenge = ch
			}
		}
		c.JSON(http.StatusOK, gin.H{"challenge": challenge})
		return
	}

	// 2. Verify Signature for event callbacks
	encryptKey := ""
	if user.LarkEncryptKey != nil {
		encryptKey = *user.LarkEncryptKey
	}
	if !verifyLarkSignature(c, body, encryptKey) {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
		return
	}

	// 3. Handle Events
	if req.Header.EventType != "" {
		if larkService == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lark service not initialized"})
			return
		}
		c.Status(http.StatusOK)
		go handleLarkEvent(req, user)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func handleLarkEvent(req dto.LarkWebhookRequest, user model.User) {
	ctx := context.Background()

	switch req.Header.EventType {
	case "im.message.receive_v1":
		eventBody, _ := json.Marshal(req.Event)
		var event dto.LarkMessageEvent
		if err := json.Unmarshal(eventBody, &event); err != nil {
			fmt.Printf("failed to unmarshal lark message event: %v\n", err)
			return
		}

		openID := event.Sender.SenderID.OpenID
		// Send quick ack immediately so user knows bot is alive
		larkService.SendAck(ctx, user, openID)

		// Process the message and send full response
		if err := larkService.HandleIncomingMessage(ctx, user, event); err != nil {
			fmt.Printf("failed to handle incoming lark message: %v\n", err)
		}
	default:
		fmt.Printf("unhandled lark event type: %s\n", req.Header.EventType)
	}
}

// BindLark handles manual Lark binding from H5
func BindLark(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		OpenID string `json:"openid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[BindLark] Invalid request body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	fmt.Printf("[BindLark] User %d attempting to bind OpenID: %s\n", userID, req.OpenID)

	if config.Cfg.Storage.Mode == "memory" {
		store := memory.GetStore()
		u := store.FindUser(userID)
		if u == nil {
			fmt.Printf("[BindLark] User %d not found in memory store\n", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found in memory store"})
			return
		}
		u.LarkOpenID = &req.OpenID
		fmt.Printf("[BindLark] User %d bound successfully (Memory Mode)\n", userID)
		c.JSON(http.StatusOK, gin.H{"message": "Lark account bound successfully"})
		return
	}

	if larkService == nil {
		fmt.Printf("[BindLark] larkService is nil! Ensure InitLark was called.\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lark service not initialized on server"})
		return
	}

	if err := larkService.BindUser(userID, req.OpenID); err != nil {
		fmt.Printf("[BindLark] Failed to bind user %d: %v\n", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	fmt.Printf("[BindLark] User %d bound successfully (Database Mode)\n", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Lark account bound successfully"})
}
