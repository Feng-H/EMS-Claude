package service

import (
	"context"
	"encoding/json"
	"fmt"

	agentDto "github.com/ems/backend/internal/agent/dto"
	agentService "github.com/ems/backend/internal/agent/service"
	"github.com/ems/backend/internal/dto"
	"github.com/ems/backend/internal/repository"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/lark"
)

type LarkService struct {
	userRepo     *repository.UserRepository
	agentService *agentService.AgentService
}

func NewLarkService() *LarkService {
	return &LarkService{
		userRepo:     repository.NewUserRepository(),
		agentService: agentService.NewAgentService(),
	}
}

func (s *LarkService) SendAck(ctx context.Context, appID, openID string) error {
	user, err := s.userRepo.GetByLarkAppID(appID)
	if err != nil {
		return err
	}
	client := lark.NewClient(user.LarkAppID, user.LarkAppSecret)
	return client.SendTextMessage(ctx, "open_id", openID, "👍 收到，正在分析中...")
}

func (s *LarkService) HandleIncomingMessage(ctx context.Context, appID string, event dto.LarkMessageEvent) error {
	openID := event.Sender.SenderID.OpenID
	if openID == "" {
		return fmt.Errorf("missing openid in lark event")
	}

	// 1. Try to find user by AppID
	user, err := s.userRepo.GetByLarkAppID(appID)
	if err != nil {
		return fmt.Errorf("unrecognized lark app id: %s", appID)
	}

	client := lark.NewClient(user.LarkAppID, user.LarkAppSecret)

	// 2. Check if OpenID matches (ensure this message is for this user's bot)
	if user.LarkOpenID == nil || *user.LarkOpenID != openID {
		// If not bound to THIS user yet, maybe it's the first time?
		// We use the same binding guide but the user is already found by appID
		return s.sendBindingGuide(ctx, client, openID)
	}

	// 3. Parse text content
	var content dto.LarkMessageTextContent
	if err := json.Unmarshal([]byte(event.Message.Content), &content); err != nil {
		return err
	}

	// 4. Call Agent service
	chatReq := &agentDto.ChatRequest{
		Message: content.Text,
	}

	resp, err := s.agentService.Chat(user.ID, string(user.Role), chatReq)
	if err != nil {
		return client.SendTextMessage(ctx, "open_id", openID, "抱歉，分析过程中出现了点问题："+err.Error())
	}

	return client.SendTextMessage(ctx, "open_id", openID, resp.Reply)
}

func (s *LarkService) sendBindingGuide(ctx context.Context, client *lark.Client, openID string) error {
	baseURL := config.Cfg.App.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:5173" // Fallback
	}
	bindURL := fmt.Sprintf("%s/h5/bind-lark?openid=%s", baseURL, openID)
	text := fmt.Sprintf("您尚未完成 EMS 系统账号绑定（或此机器人属于其他用户）。请点击下方链接完成身份验证：\n%s", bindURL)
	return client.SendTextMessage(ctx, "open_id", openID, text)
}

func (s *LarkService) BindUser(userID uint, openID string) error {
	return s.userRepo.UpdateLarkOpenID(userID, openID)
}
