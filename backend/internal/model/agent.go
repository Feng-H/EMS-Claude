package model

import (
	"time"
)

// =====================================================
// Manual & Retrieval Models
// =====================================================

type ManualDocument struct {
	BaseModel
	EquipmentTypeID *uint          `json:"equipment_type_id"`
	EquipmentType   *EquipmentType `json:"equipment_type,omitempty" gorm:"foreignKey:EquipmentTypeID"`
	EquipmentID     *uint          `json:"equipment_id"`
	Equipment       *Equipment     `json:"equipment,omitempty" gorm:"foreignKey:EquipmentID"`
	Name            string         `json:"name" gorm:"size:255;not null"`
	FileURL         string         `json:"file_url" gorm:"size:500;not null"`
	Version         string         `json:"version" gorm:"size:50"`
	Source          string         `json:"source" gorm:"size:100"`
	UploadedBy      uint           `json:"uploaded_by"`
	Uploader        *User          `json:"uploader,omitempty" gorm:"foreignKey:UploadedBy"`
	ParsingStatus   string         `json:"parsing_status" gorm:"size:20;default:'pending'"`
	Chunks          []ManualChunk  `json:"chunks,omitempty" gorm:"foreignKey:DocumentID"`
}

type ManualChunk struct {
	ID           uint    `json:"id" gorm:"primarykey"`
	DocumentID   uint    `json:"document_id" gorm:"not null"`
	SectionTitle string  `json:"section_title" gorm:"size:255"`
	Content      string  `json:"content" gorm:"type:text;not null"`
	ChunkIndex   int     `json:"chunk_index" gorm:"default:0"`
	EmbeddingRef string  `json:"embedding_ref" gorm:"size:255"`
	SearchVector string  `json:"-" gorm:"type:tsvector"` // PostgreSQL tsvector for FTS
	CreatedAt    time.Time `json:"created_at"`
}

// =====================================================
// Runtime & Cost Analysis Models
// =====================================================

type RepairCostDetail struct {
	BaseModel
	OrderID       uint         `json:"order_id" gorm:"uniqueIndex;not null"`
	Order         *RepairOrder `json:"-" gorm:"foreignKey:OrderID"`
	LaborCost     float64      `json:"labor_cost" gorm:"type:decimal(12,2);default:0"`
	SparePartCost float64      `json:"spare_part_cost" gorm:"type:decimal(12,2);default:0"`
	OutsourceCost float64      `json:"outsource_cost" gorm:"type:decimal(12,2);default:0"`
	OtherCost     float64      `json:"other_cost" gorm:"type:decimal(12,2);default:0"`
	Currency      string       `json:"currency" gorm:"size:10;default:'CNY'"`
}

type EquipmentRuntimeSnapshot struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	EquipmentID  uint      `json:"equipment_id" gorm:"uniqueIndex:idx_equip_date;not null"`
	Equipment    *Equipment `json:"-" gorm:"foreignKey:EquipmentID"`
	SnapshotDate string    `json:"snapshot_date" gorm:"uniqueIndex:idx_equip_date;size:10;not null"` // YYYY-MM-DD
	RuntimeHours float64   `json:"runtime_hours" gorm:"type:decimal(10,2);default:0"`
	DowntimeHours float64  `json:"downtime_hours" gorm:"type:decimal(10,2);default:0"`
	LoadRate     float64   `json:"load_rate" gorm:"type:decimal(5,2)"`
	OutputQty    float64   `json:"output_qty" gorm:"type:decimal(12,2)"`
	CreatedAt    time.Time `json:"created_at"`
}

// =====================================================
// Agent Session & Artifact Models
// =====================================================

type AgentSession struct {
	BaseModel
	UserID        uint      `json:"user_id" gorm:"not null"`
	User          *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Scenario      string    `json:"scenario" gorm:"size:50;not null"`
	FactoryID     *uint     `json:"factory_id"`
	WorkshopID    *uint     `json:"workshop_id"`
	Language      string    `json:"language" gorm:"size:20;default:'zh-CN'"`
	QueryText     string    `json:"query_text" gorm:"type:text"`
	InputSnapshot string    `json:"input_snapshot" gorm:"type:jsonb"`
	Status        string    `json:"status" gorm:"size:20;default:'completed'"`
	TraceID       string    `json:"trace_id" gorm:"uniqueIndex;size:100;not null"`
	Artifacts     []AgentArtifact `json:"artifacts,omitempty" gorm:"foreignKey:SessionID"`
}

type AgentArtifact struct {
	BaseModel
	SessionID     uint          `json:"session_id" gorm:"not null"`
	ArtifactType  string        `json:"artifact_type" gorm:"size:30;not null"`
	Title         string        `json:"title" gorm:"size:255"`
	Summary       string        `json:"summary" gorm:"type:text"`
	InputSnapshot string        `json:"input_snapshot" gorm:"type:jsonb"`
	ResultJSON    string        `json:"result_json" gorm:"type:jsonb;not null"`
	RiskLevel     string        `json:"risk_level" gorm:"size:20"`
	ConfirmedBy   *uint         `json:"confirmed_by"`
	Confirmer     *User         `json:"confirmer,omitempty" gorm:"foreignKey:ConfirmedBy"`
	ConfirmedAt   *time.Time    `json:"confirmed_at"`
	EvidenceLinks []AgentEvidenceLink `json:"evidence_links,omitempty" gorm:"foreignKey:ArtifactID"`
}

type AgentEvidenceLink struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	ArtifactID  uint      `json:"artifact_id" gorm:"not null"`
	EvidenceType string   `json:"evidence_type" gorm:"size:30;not null"`
	SourceTable string    `json:"source_table" gorm:"size:100"`
	SourceID    uint      `json:"source_id"`
	Excerpt     string    `json:"excerpt" gorm:"type:text"`
	Score       float64   `json:"score" gorm:"type:decimal(5,4)"`
	CreatedAt   time.Time `json:"created_at"`
}
