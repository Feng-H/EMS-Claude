package dto

import "time"

// Organization DTOs

type BaseRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type BaseResponse struct {
	ID        uint      `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FactoryRequest struct {
	BaseID uint   `json:"base_id" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

type FactoryResponse struct {
	ID        uint      `json:"id"`
	BaseID    uint      `json:"base_id"`
	BaseName  string    `json:"base_name,omitempty"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WorkshopRequest struct {
	FactoryID uint   `json:"factory_id" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

type WorkshopResponse struct {
	ID         uint      `json:"id"`
	FactoryID  uint      `json:"factory_id"`
	FactoryName string   `json:"factory_name,omitempty"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Equipment Type DTOs

type EquipmentTypeRequest struct {
	Name string `json:"name" binding:"required"`
	Category string `json:"category"`
	InspectionTemplateID *uint `json:"inspection_template_id"`
}

type EquipmentTypeResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	TemplateID  *uint  `json:"inspection_template_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// Equipment DTOs

type EquipmentRequest struct {
	Code                   string     `json:"code" binding:"required"`
	Name                   string     `json:"name" binding:"required"`
	TypeID                 uint       `json:"type_id" binding:"required"`
	WorkshopID             uint       `json:"workshop_id" binding:"required"`
	Spec                   string     `json:"spec"`
	PurchaseDate           *time.Time `json:"purchase_date"`
	Status                 string     `json:"status"`
	DedicatedMaintenanceID *uint      `json:"dedicated_maintenance_id"`
}

type EquipmentResponse struct {
	ID                      uint       `json:"id"`
	Code                    string     `json:"code"`
	Name                    string     `json:"name"`
	TypeID                  uint       `json:"type_id"`
	TypeName                string     `json:"type_name,omitempty"`
	WorkshopID              uint       `json:"workshop_id"`
	WorkshopName            string     `json:"workshop_name,omitempty"`
	FactoryID               uint       `json:"factory_id,omitempty"`
	FactoryName             string     `json:"factory_name,omitempty"`
	QRCode                  string     `json:"qr_code"`
	Spec                    string     `json:"spec"`
	PurchaseDate            *time.Time `json:"purchase_date"`
	Status                  string     `json:"status"`
	DedicatedMaintenanceID  *uint      `json:"dedicated_maintenance_id"`
	DedicatedMaintenanceName string    `json:"dedicated_maintenance_name,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

type EquipmentListResponse struct {
	Total int64                `json:"total"`
	Items []EquipmentResponse  `json:"items"`
}

type EquipmentQuery struct {
	Page       int    `form:"page" binding:"min=1"`
	PageSize   int    `form:"page_size" binding:"min=1,max=100"`
	Code       string `form:"code"`
	Name       string `form:"name"`
	TypeID     *uint  `form:"type_id"`
	FactoryID  *uint  `form:"factory_id"`
	WorkshopID *uint  `form:"workshop_id"`
	Status     string `form:"status"`
}

type QRCodeRequest struct {
	EquipmentID []uint `json:"equipment_ids" binding:"required"`
}

type QRCodeResponse struct {
	EquipmentID uint   `json:"equipment_id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	QRCodeURL   string `json:"qr_code_url"`
	QRCodeData  string `json:"qr_code_data"` // Base64 encoded image
}

// Batch Import DTOs

type EquipmentImportItem struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	TypeCode    string `json:"type_code" binding:"required"`
	FactoryCode string `json:"factory_code" binding:"required"`
	WorkshopCode string `json:"workshop_code" binding:"required"`
	Spec        string `json:"spec"`
}

type BatchImportRequest struct {
	Items []EquipmentImportItem `json:"items" binding:"required"`
}

type BatchImportResponse struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	Errors       []string `json:"errors"`
}

// Organization Tree Response

type OrganizationTree struct {
	ID       uint                `json:"id"`
	Code     string              `json:"code"`
	Name     string              `json:"name"`
	Type     string              `json:"type"` // base, factory, workshop
	Children []OrganizationTree  `json:"children"`
}

// Equipment Statistics

type EquipmentStatistics struct {
	Total      int64 `json:"total"`
	Running    int64 `json:"running"`
	Stopped    int64 `json:"stopped"`
	Maintenance int64 `json:"maintenance"`
	Scrapped   int64 `json:"scrapped"`
}

type EquipmentTypeStatistics struct {
	TypeName  string `json:"type_name"`
	Count     int64  `json:"count"`
	Running   int64  `json:"running"`
	Maintenance int64 `json:"maintenance"`
}
