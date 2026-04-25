package model

import (
	"time"
)

// =====================================================
// Phase 2: Skills, Knowledge, Experience & Conversations
// =====================================================

// AgentSkill 技能库：一套可执行的分析方法
type AgentSkill struct {
	ID                  uint      `json:"id" gorm:"primarykey"`
	Name                string    `json:"name" gorm:"size:200;not null"`
	Description         string    `json:"description" gorm:"type:text"`
	ApplicableTo        string    `json:"applicable_to" gorm:"type:jsonb;default:'[]'"` // 适用对象类型
	ApplicableScenarios string    `json:"applicable_scenarios" gorm:"type:jsonb;default:'[]'"` // 适用场景描述
	Steps               string    `json:"steps" gorm:"type:jsonb;not null"` // 执行步骤序列
	Scope               string    `json:"scope" gorm:"type:jsonb;default:'{}'"` // 权限/范围限制
	Version             int       `json:"version" gorm:"default:1"`
	Status              string    `json:"status" gorm:"size:20;default:'draft';index"`
	CreatedBy           string    `json:"created_by" gorm:"size:100"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	UsageCount          int       `json:"usage_count" gorm:"default:0"`
	SuccessRate         float64   `json:"success_rate" gorm:"type:decimal(5,4);default:0"`
	LastUsed            *time.Time `json:"last_used"`
	Changelog           string    `json:"changelog" gorm:"type:jsonb;default:'[]'"`
}

// AgentKnowledge 智能体知识库：分析过程中产出的有价值结论
type AgentKnowledge struct {
	ID               string    `json:"id" gorm:"primarykey;size:50"` // 字符串ID (e.g. k_20251215_001)
	Title            string    `json:"title" gorm:"size:500;not null"`
	Type             string    `json:"type" gorm:"size:50;not null;index"` // root_cause_analysis, pattern, etc.
	Scope            string    `json:"scope" gorm:"type:jsonb;default:'{}'"` // 关联的工厂/车间/设备类型
	Summary          string    `json:"summary" gorm:"type:text;not null"`
	Details          string    `json:"details" gorm:"type:jsonb"`
	RelatedSkillID   string    `json:"related_skill_id" gorm:"size:100"`
	Confidence       float64   `json:"confidence" gorm:"type:decimal(5,4);default:0"`
	Status           string    `json:"status" gorm:"size:20;default:'draft';index"` // draft, confirmed, rejected, archived
	CreatedBy        string    `json:"created_by" gorm:"size:100"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	VerifiedBy       *uint     `json:"verified_by"`
	Verifier         *User     `json:"verifier,omitempty" gorm:"foreignKey:VerifiedBy"`
	VerifiedAt       *time.Time `json:"verified_at"`
	ReferencedCount  int       `json:"referenced_count" gorm:"default:0"`
	LastReferenced   *time.Time `json:"last_referenced"`
	ExpireAt         *time.Time `json:"expire_at"`
	SearchVector     string    `json:"-" gorm:"type:tsvector"` // 用于全文搜索
}

// AgentExperience 经验库：用户偏好与行为校准信息
type AgentExperience struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Type        string    `json:"type" gorm:"size:50;not null"` // preference, correction, boundary, quality_feedback
	Category    string    `json:"category" gorm:"size:100"` // analysis_depth, display_style, etc.
	Content     string    `json:"content" gorm:"type:jsonb;not null"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	User        *User     `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Weight      float64   `json:"weight" gorm:"type:decimal(5,4);default:0.8"`
	DecayRate   float64   `json:"decay_rate" gorm:"type:decimal(5,4);default:0.05"`
	Status      string    `json:"status" gorm:"size:20;default:'active';index"`
	CreatedAt   time.Time `json:"created_at"`
	LastApplied *time.Time `json:"last_applied"`
	ExpireAt    *time.Time `json:"expire_at"`
}

// AgentConversation Agent 对话会话
type AgentConversation struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Title     string    `json:"title" gorm:"size:500"`
	Status    string    `json:"status" gorm:"size:20;default:'active'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Messages  []AgentMessage `json:"messages,omitempty" gorm:"foreignKey:ConversationID"`
}

// AgentMessage 对话历史消息
type AgentMessage struct {
	ID             uint      `json:"id" gorm:"primarykey"`
	ConversationID uint      `json:"conversation_id" gorm:"not null;index"`
	Conversation   *AgentConversation `json:"-" gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE;"`
	Role           string    `json:"role" gorm:"size:20;not null"` // system, user, assistant, tool
	Content        string    `json:"content" gorm:"type:text;not null"`
	ToolCalls      string    `json:"tool_calls" gorm:"type:jsonb"`
	SkillID        string    `json:"skill_id" gorm:"size:100"`
	KnowledgeIDs   string    `json:"knowledge_ids" gorm:"type:jsonb;default:'[]'"`
	CreatedAt      time.Time `json:"created_at"`
}

// AgentPushSubscription 推送订阅管理
type AgentPushSubscription struct {
	ID              uint      `json:"id" gorm:"primarykey"`
	UserID          uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_push"`
	User            *User     `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	PushType        string    `json:"push_type" gorm:"size:50;not null;uniqueIndex:idx_user_push"` // daily_report, alert, cycle_report, knowledge
	Scope           string    `json:"scope" gorm:"type:jsonb;default:'{}'"` // 关注的车间/设备类型范围
	Frequency       string    `json:"frequency" gorm:"size:20;default:'daily'"`
	QuietHoursStart string    `json:"quiet_hours_start" gorm:"size:8"` // HH:MM:SS
	QuietHoursEnd   string    `json:"quiet_hours_end" gorm:"size:8"`   // HH:MM:SS
	Enabled         bool      `json:"enabled" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
