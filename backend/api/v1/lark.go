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
			fmt.Printf("[verifyLarkSignature] Warning: Received signature from Lark but Encrypt Key is NOT configured in EMS. Allowing request for now. Please configure Encrypt Key for security.\n")
		}
		return true
	}

	if signature == "" {
		fmt.Printf("[verifyLarkSignature] Error: Missing signature from Lark but Encrypt Key IS configured in EMS.\n")
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
	var user model.User
	var err error

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("[LarkWebhook] Failed to read body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	var req dto.LarkWebhookRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Printf("[LarkWebhook] Failed to unmarshal request: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Identify User
	if userIDStr != "" {
		// Use ID from URL if present
		userID, _ := strconv.ParseUint(userIDStr, 10, 64)
		if config.Cfg.Storage.Mode == "memory" {
			u := memory.GetStore().FindUser(uint(userID))
			if u != nil { user = *u } else { err = fmt.Errorf("user not found") }
		} else {
			err = db.First(&user, userID).Error
		}
	} else {
		// Try to match by AppID from headers/body
		appID := req.Header.AppID
		if appID == "" {
			// Fallback: check nested event for challenge
			if req.Event != nil {
				if aid, ok := req.Event["app_id"].(string); ok { appID = aid }
			}
		}

		if appID != "" {
			if config.Cfg.Storage.Mode == "memory" {
				found := false
				for _, u := range memory.GetStore().Users {
					if u.LarkAppID != nil && *u.LarkAppID == appID {
						user = *u; found = true; break
					}
				}
				if !found { err = fmt.Errorf("no user configured with app_id: %s", appID) }
			} else {
				err = db.Where("lark_app_id = ?", appID).First(&user).Error
			}
		} else {
			err = fmt.Errorf("missing app_id in request and no user_id in URL")
		}
	}

	if err != nil {
		fmt.Printf("[LarkWebhook] User identification failed: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
		return
	}

	fmt.Printf("[LarkWebhook] Processing event for user %d (AppID: %s, Event: %s)\n", 
		user.ID, req.Header.AppID, req.Header.EventType)

	// 2. Handle URL Challenge FIRST
	if req.Type == "url_verification" || req.Header.EventType == "url_verification" {
		token := req.Token
		if token == "" { token = req.Header.Token }
		
		expectedToken := ""
		if user.LarkVerificationToken != nil { expectedToken = *user.LarkVerificationToken }

		if token != expectedToken {
			fmt.Printf("[LarkWebhook] Verification token mismatch. Got: %s, Expected: %s\n", token, expectedToken)
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}
		challenge := req.Challenge
		if challenge == "" && req.Event != nil {
			if ch, ok := req.Event["challenge"].(string); ok { challenge = ch }
		}
		c.JSON(http.StatusOK, gin.H{"challenge": challenge})
		return
	}

	// 3. Verify Signature
	encryptKey := ""
	if user.LarkEncryptKey != nil { encryptKey = *user.LarkEncryptKey }
	if !verifyLarkSignature(c, body, encryptKey) {
		fmt.Printf("[LarkWebhook] Signature verification failed for user %d\n", user.ID)
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
		return
	}

	// 4. Handle Events
	if req.Header.EventType != "" {
		if larkService == nil {
			fmt.Printf("[LarkWebhook] larkService is nil!\n")
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
