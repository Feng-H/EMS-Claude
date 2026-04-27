package dto

// LarkWebhookRequest represents the base structure of Lark event requests
type LarkWebhookRequest struct {
	Schema string                 `json:"schema"`
	Header LarkWebhookHeader      `json:"header"`
	Event  map[string]interface{} `json:"event"`
	
	// Challenge fields
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

type LarkWebhookHeader struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	CreateTime string `json:"create_time"`
	Token      string `json:"token"`
	AppID      string `json:"app_id"`
	TenantKey  string `json:"tenant_key"`
}

type LarkMessageEvent struct {
	Sender  LarkSender  `json:"sender"`
	Message LarkMessage `json:"message"`
}

type LarkSender struct {
	SenderID   LarkSenderID `json:"sender_id"`
	SenderType string       `json:"sender_type"`
	TenantKey  string       `json:"tenant_key"`
}

type LarkSenderID struct {
	UnionID string `json:"union_id"`
	UserID  string `json:"user_id"`
	OpenID  string `json:"open_id"`
}

type LarkMessage struct {
	MessageID   string `json:"message_id"`
	RootID      string `json:"root_id"`
	ParentID    string `json:"parent_id"`
	CreateTime  string `json:"create_time"`
	ChatID      string `json:"chat_id"`
	ChatType    string `json:"chat_type"`
	MessageType string `json:"message_type"`
	Content     string `json:"content"` // JSON string
}

type LarkMessageTextContent struct {
	Text string `json:"text"`
}
