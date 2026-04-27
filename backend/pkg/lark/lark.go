package lark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/redis"
)

const (
	tokenURL       = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	sendMessageURL = "https://open.feishu.cn/open-apis/im/v1/messages"
	redisTokenKey  = "ems:lark:tenant_access_token"
)

type Client struct {
	appID     string
	appSecret string
	mu        sync.RWMutex
}

var (
	defaultClient *Client
	once          sync.Once
)

func GetClient() *Client {
	once.Do(func() {
		defaultClient = &Client{
			appID:     config.Cfg.Lark.AppID,
			appSecret: config.Cfg.Lark.AppSecret,
		}
	})
	return defaultClient
}

type TokenResponse struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

func (c *Client) GetTenantAccessToken(ctx context.Context) (string, error) {
	// 1. Try to get from Redis
	if redis.Client != nil {
		token, err := redis.Client.Get(ctx, redisTokenKey).Result()
		if err == nil && token != "" {
			return token, nil
		}
	}

	// 2. Fetch from Lark API
	reqBody, _ := json.Marshal(map[string]string{
		"app_id":     c.appID,
		"app_secret": c.appSecret,
	})

	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Code != 0 {
		return "", fmt.Errorf("lark api error: %s (code: %d)", result.Msg, result.Code)
	}

	// 3. Cache to Redis if possible
	if redis.Client != nil {
		// Set expiration slightly earlier than the actual one
		expire := time.Duration(result.Expire-60) * time.Second
		if expire < 0 {
			expire = 1 * time.Minute
		}
		redis.Client.Set(ctx, redisTokenKey, result.TenantAccessToken, expire)
	}

	return result.TenantAccessToken, nil
}

type MessageRequest struct {
	ReceiveID string `json:"receive_id"`
	MsgType   string `json:"msg_type"`
	Content   string `json:"content"`
	UUID      string `json:"uuid,omitempty"`
}

func (c *Client) SendTextMessage(ctx context.Context, receiveIDType, receiveID, text string) error {
	content, _ := json.Marshal(map[string]string{
		"text": text,
	})
	return c.SendMessage(ctx, receiveIDType, receiveID, "text", string(content))
}

func (c *Client) SendCardMessage(ctx context.Context, receiveIDType, receiveID, cardJSON string) error {
	return c.SendMessage(ctx, receiveIDType, receiveID, "interactive", cardJSON)
}

func (c *Client) SendMessage(ctx context.Context, receiveIDType, receiveID, msgType, content string) error {
	token, err := c.GetTenantAccessToken(ctx)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s?receive_id_type=%s", sendMessageURL, receiveIDType)
	
	reqBody, _ := json.Marshal(MessageRequest{
		ReceiveID: receiveID,
		MsgType:   msgType,
		Content:   content,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Code != 0 {
		return fmt.Errorf("lark send message error: %s (code: %d)", result.Msg, result.Code)
	}

	return nil
}
