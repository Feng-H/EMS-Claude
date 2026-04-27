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
	// Read body for signature verification
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	// 1. Verify Signature
	if !verifyLarkSignature(c, body) {
		// During debug, we might want to log this but not block if AppSecret is empty
		if config.Cfg.Lark.AppSecret != "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
			return
		}
	}

	var req dto.LarkWebhookRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Handle URL Challenge
	if req.Type == "url_verification" {
		if req.Token != config.Cfg.Lark.VerificationToken {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"challenge": req.Challenge})
		return
	}

	// 2. Handle Events (V2 Schema)
	if req.Header.EventType != "" {
		// Quick response to Lark
		c.Status(http.StatusOK)

		// Async processing
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
