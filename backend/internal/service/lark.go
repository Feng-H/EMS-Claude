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

func (s *LarkService) getClient(user model.User) (*lark.Client, error) {
	if user.LarkAppID == nil || user.LarkAppSecret == nil || *user.LarkAppID == "" || *user.LarkAppSecret == "" {
		return nil, fmt.Errorf("lark bot not configured for this user")
	}
	return lark.NewClient(*user.LarkAppID, *user.LarkAppSecret), nil
}

func (s *LarkService) SendAck(ctx context.Context, botUser model.User, openID string) error {
	client, err := s.getClient(botUser)
	if err != nil {
		return err
	}
	return client.SendTextMessage(ctx, "open_id", openID, "👍 收到，正在分析中...")
}

func (s *LarkService) HandleIncomingMessage(ctx context.Context, botUser model.User, event dto.LarkMessageEvent) error {
	openID := event.Sender.SenderID.OpenID
	if openID == "" {
		return fmt.Errorf("missing openid in lark event")
	}

	client, err := s.getClient(botUser)
	if err != nil {
		return err
	}

	// 1. Try to find bound user (the sender)
	user, err := s.userRepo.GetByLarkOpenID(openID)
	if err != nil {
		// Not bound, send binding link
		return s.sendBindingGuide(ctx, client, openID)
	}

	// 2. Parse text content
	var content dto.LarkMessageTextContent
	if err := json.Unmarshal([]byte(event.Message.Content), &content); err != nil {
		return err
	}

	// 3. Call Agent service
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
	text := fmt.Sprintf("您尚未绑定 EMS 系统账号。请点击下方链接完成身份验证后，即可在飞书中使用智能助手：\n%s", bindURL)
	return client.SendTextMessage(ctx, "open_id", openID, text)
}

func (s *LarkService) BindUser(userID uint, openID string) error {
	return s.userRepo.UpdateLarkOpenID(userID, openID)
}
