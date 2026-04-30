package v1

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ems/backend/internal/dto"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/internal/repository"
	"github.com/ems/backend/internal/service"
	"github.com/ems/backend/pkg/config"
	"github.com/gin-gonic/gin"
)

var (
	larkService *service.LarkService
)

func InitLark() {
	larkService = service.NewLarkService()
}

// verifyLarkSignature ensures the request is coming from Lark using user's specific encrypt_key
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
// @Summary Lark webhook
// @Tags lark
// @Router /lark/webhook [post]
func LarkWebhook(c *gin.Context) {
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

	// 0. Identify user by AppID
	appID := req.Header.AppID
	if appID == "" {
		// Fallback for some events
		if eventAppID, ok := req.Event["app_id"].(string); ok {
			appID = eventAppID
		}
	}
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing app_id"})
		return
	}

	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByLarkAppID(appID)
	if err != nil {
		log.Printf("[Lark] Unrecognized AppID: %s", appID)
		c.JSON(http.StatusForbidden, gin.H{"error": "unrecognized app_id"})
		return
	}

	// 1. Handle URL Challenge FIRST
	if req.Type == "url_verification" || req.Header.EventType == "url_verification" {
		token := req.Token
		if token == "" {
			token = req.Header.Token
		}
		if token != user.LarkVerifyToken {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid verification token"})
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

	// 2. Verify Signature
	if !verifyLarkSignature(c, body, user.LarkEncryptKey) {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
		return
	}

	// 3. Handle Events
	if req.Header.EventType != "" {
		c.Status(http.StatusOK)
		go handleLarkEvent(appID, req)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func handleLarkEvent(appID string, req dto.LarkWebhookRequest) {
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
		// Send quick ack immediately
		larkService.SendAck(ctx, appID, openID)

		// Process the message
		if err := larkService.HandleIncomingMessage(ctx, appID, event); err != nil {
			fmt.Printf("failed to handle incoming lark message: %v\n", err)
		}
	default:
		fmt.Printf("unhandled lark event type: %s\n", req.Header.EventType)
	}
}

// UpdateLarkConfig updates the user's private bot credentials
// @Summary Update Lark bot config
// @Tags auth
// @Router /auth/lark-config [post]
func UpdateLarkConfig(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		AppID       string `json:"lark_app_id" binding:"required"`
		AppSecret   string `json:"lark_app_secret" binding:"required"`
		VerifyToken string `json:"lark_verify_token"`
		EncryptKey  string `json:"lark_encrypt_key"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRepo := repository.NewUserRepository()
	updates := map[string]interface{}{
		"lark_app_id":       req.AppID,
		"lark_app_secret":   req.AppSecret,
		"lark_verify_token": req.VerifyToken,
		"lark_encrypt_key":  req.EncryptKey,
	}

	if err := userRepo.UpdateLarkCredentials(userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lark configuration updated successfully"})
}

// BindLark handles account binding from H5 (OpenID)
// @Summary Bind Lark account (OpenID)
// @Tags auth
// @Router /auth/bind-lark [post]
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := larkService.BindUser(userID, req.OpenID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lark account bound successfully"})
}
