package v1

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ems/backend/internal/dto"
	"github.com/ems/backend/internal/middleware"
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

// verifyLarkSignature ensures the request is coming from Lark
func verifyLarkSignature(c *gin.Context, body []byte) bool {
	timestamp := c.GetHeader("X-Lark-Request-Timestamp")
	nonce := c.GetHeader("X-Lark-Request-Nonce")
	signature := c.GetHeader("X-Lark-Signature")
	appSecret := config.Cfg.Lark.AppSecret

	// Signature = sha256(timestamp + nonce + appSecret + body)
	content := timestamp + nonce + appSecret + string(body)
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

	// 1. Handle URL Challenge FIRST (skip signature check — the challenge IS the verification)
	if req.Type == "url_verification" || req.Header.EventType == "url_verification" {
		token := req.Token
		if token == "" {
			token = req.Header.Token
		}
		if token != config.Cfg.Lark.VerificationToken {
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
	if config.Cfg.Lark.AppSecret != "" && !verifyLarkSignature(c, body) {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
		return
	}

	// 3. Handle Events (V2 Schema)
	if req.Header.EventType != "" {
		c.Status(http.StatusOK)
		go handleLarkEvent(req)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func handleLarkEvent(req dto.LarkWebhookRequest) {
	ctx := context.Background()
	
	switch req.Header.EventType {
	case "im.message.receive_v1":
		eventBody, _ := json.Marshal(req.Event)
		var event dto.LarkMessageEvent
		if err := json.Unmarshal(eventBody, &event); err != nil {
			fmt.Printf("failed to unmarshal lark message event: %v\n", err)
			return
		}
		
		if err := larkService.HandleIncomingMessage(ctx, event); err != nil {
			fmt.Printf("failed to handle incoming lark message: %v\n", err)
		}
	default:
		fmt.Printf("unhandled lark event type: %s\n", req.Header.EventType)
	}
}

// BindLark handles account binding from H5
// @Summary Bind Lark account
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
