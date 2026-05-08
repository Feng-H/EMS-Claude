package middleware
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ems/backend/pkg/jwt"
	"github.com/ems/backend/pkg/redis"
	"github.com/ems/backend/internal/model"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

const (
	ContextKeyUserID           = "user_id"
	ContextKeyRole             = "role"
	ContextKeyAPIKeyID         = "api_key_id"
	ContextKeyAPIKeyScopes     = "api_key_scopes"
	ContextKeyAPIKeyRateLimit  = "api_key_rate_limit"
)

var db *gorm.DB

func Init(database *gorm.DB) {
	db = database
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Try API Key first (for external Agents)
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey != "" {
			// Hash incoming key for comparison
			hash := sha256.Sum256([]byte(apiKey))
			hashedKey := hex.EncodeToString(hash[:])

			var keyRecord model.UserAPIKey
			err := db.Preload("User").Where("key = ? AND is_active = ?", hashedKey, true).First(&keyRecord).Error
			if err == nil {
				// Check expiration
				if keyRecord.ExpiresAt == nil || keyRecord.ExpiresAt.After(time.Now()) {
					// Enforce Rate Limit
					if keyRecord.RateLimit > 0 {
						if !checkRateLimit(keyRecord.ID, keyRecord.RateLimit) {
							c.JSON(http.StatusTooManyRequests, gin.H{"error": "API rate limit exceeded"})
							c.Abort()
							return
						}
					}

					// Update last used
					now := time.Now()
					db.Model(&keyRecord).Update("last_used_at", &now)

					c.Set(ContextKeyUserID, keyRecord.UserID)
					c.Set(ContextKeyRole, string(keyRecord.User.Role))
					c.Set(ContextKeyAPIKeyID, keyRecord.ID)
					
					scopes := strings.Split(keyRecord.Scopes, ",")
					trimmedScopes := make([]string, 0, len(scopes))
					for _, s := range scopes {
						if ts := strings.TrimSpace(s); ts != "" {
							trimmedScopes = append(trimmedScopes, ts)
						}
					}
					c.Set(ContextKeyAPIKeyScopes, trimmedScopes)
					c.Next()
					return
				}
			}
		}

		// 2. Try JWT (for Web/H5 users)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyRole, claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	roleMap := make(map[string]bool)
	for _, role := range roles {
		roleMap[role] = true
	}

	return func(c *gin.Context) {
		userRole, exists := c.Get(ContextKeyRole)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok || !roleMap[roleStr] {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return 0, false
	}
	uid, ok := userID.(uint)
	if !ok {
		return 0, false
	}
	return uid, true
}

func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(ContextKeyRole)
	if !exists {
		return "", false
	}
	roleStr, ok := role.(string)
	if !ok {
		return "", false
	}
	return roleStr, true
}

func checkRateLimit(apiKeyID uint, limit int) bool {
	if redis.Client == nil {
		return true // Skip if redis is not available
	}

	key := fmt.Sprintf("ratelimit:apikey:%d:%d", apiKeyID, time.Now().Unix()/60)
	count, err := redis.Client.Incr(redis.Ctx, key).Result()
	if err != nil {
		return true // Fail open on redis error
	}

	if count == 1 {
		redis.Client.Expire(redis.Ctx, key, time.Minute)
	}

	return int(count) <= limit
}

func GetAPIKeyScopes(c *gin.Context) ([]string, bool) {
	scopes, exists := c.Get(ContextKeyAPIKeyScopes)
	if !exists {
		return nil, false
	}
	scopesList, ok := scopes.([]string)
	if !ok {
		return nil, false
	}
	return scopesList, true
}
