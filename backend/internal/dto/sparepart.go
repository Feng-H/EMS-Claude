package dto

import "time"

// =====================================================
// Spare Part DTOs
// =====================================================

// CreateSparePartRequest creates a new spare part
type CreateSparePartRequest struct {
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Specification string  `json:"specification"`
	Unit          string  `json:"unit"`
	FactoryID     *uint   `json:"factory_id"`
	SafetyStock   int     `json:"safety_stock"`
}

// UpdateSparePartRequest updates a spare part
type UpdateSparePartRequest struct {
	Code          string  `json:"code" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Specification string  `json:"specification"`
	Unit          string  `json:"unit"`
	FactoryID     *uint   `json:"factory_id"`
	SafetyStock   int     `json:"safety_stock"`
}

// SparePartResponse represents a spare part in API responses
type SparePartResponse struct {
	ID            uint   `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Specification string `json:"specification,omitempty"`
	Unit          string `json:"unit,omitempty"`
	FactoryID     *uint  `json:"factory_id,omitempty"`
	FactoryName   string `json:"factory_name,omitempty"`
	SafetyStock   int    `json:"safety_stock"`
	CurrentStock  int    `json:"current_stock"`
	CreatedAt     time.Time `json:"created_at"`
}

// SparePartListResponse represents a paginated list of spare parts
type SparePartListResponse struct {
	Total int64                 `json:"total"`
	Items []SparePartResponse `json:"items"`
}

// =====================================================
// Inventory DTOs
// =====================================================

// StockInRequest adds stock to inventory
type StockInRequest struct {
	SparePartID uint   `json:"spare_part_id" binding:"required"`
	FactoryID   uint   `json:"factory_id" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
	Remark      string `json:"remark"`
}

// StockOutRequest removes stock from inventory
type StockOutRequest struct {
	SparePartID uint   `json:"spare_part_id" binding:"required"`
	FactoryID   uint   `json:"factory_id" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
	OrderID     *uint  `json:"order_id"`     // repair order
	TaskID      *uint  `json:"task_id"`      // maintenance task
	Remark      string `json:"remark"`
}

// InventoryResponse represents inventory status
type InventoryResponse struct {
	ID            uint   `json:"id"`
	SparePartID   uint   `json:"spare_part_id"`
	SparePartCode string `json:"spare_part_code"`
	SparePartName string `json:"spare_part_name"`
	FactoryID     uint   `json:"factory_id"`
	FactoryName   string `json:"factory_name"`
	Quantity      int    `json:"quantity"`
	IsLowStock    bool   `json:"is_low_stock"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// InventoryListResponse represents a paginated list of inventory
type InventoryListResponse struct {
	Total int64               `json:"total"`
	Items []InventoryResponse `json:"items"`
}

// =====================================================
// Consumption DTOs
// =====================================================

// CreateConsumptionRequest creates a consumption record
type CreateConsumptionRequest struct {
	SparePartID uint   `json:"spare_part_id" binding:"required"`
	OrderID     *uint  `json:"order_id"`
	TaskID      *uint  `json:"task_id"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
	Remark      string `json:"remark"`
}

// ConsumptionResponse represents a consumption record
type ConsumptionResponse struct {
	ID            uint   `json:"id"`
	SparePartID   uint   `json:"spare_part_id"`
	SparePartCode string `json:"spare_part_code"`
	SparePartName string `json:"spare_part_name"`
	OrderID       *uint  `json:"order_id,omitempty"`
	TaskID        *uint  `json:"task_id,omitempty"`
	Quantity      int    `json:"quantity"`
	UserID        uint   `json:"user_id"`
	UserName      string `json:"user_name"`
	Remark        string `json:"remark,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// ConsumptionListResponse represents a paginated list of consumption records
type ConsumptionListResponse struct {
	Total int64                 `json:"total"`
	Items []ConsumptionResponse `json:"items"`
}

// =====================================================
// Statistics DTOs
// =====================================================

// SparePartStatistics represents spare part statistics
type SparePartStatistics struct {
	TotalParts       int64 `json:"total_parts"`
	LowStockCount    int64 `json:"low_stock_count"`
	TotalStockValue  float64 `json:"total_stock_value"`
	MonthlyConsumption int64 `json:"monthly_consumption"`
}

// LowStockAlert represents a low stock alert
type LowStockAlert struct {
	SparePartID   uint   `json:"spare_part_id"`
	SparePartCode string `json:"spare_part_code"`
	SparePartName string `json:"spare_part_name"`
	FactoryID     uint   `json:"factory_id"`
	FactoryName   string `json:"factory_name"`
	CurrentStock  int    `json:"current_stock"`
	SafetyStock   int    `json:"safety_stock"`
	Shortage      int    `json:"shortage"`
}

// =====================================================
// Query DTOs
// =====================================================

// SparePartQuery represents query parameters for spare parts
type SparePartQuery struct {
	Code     string `form:"code"`
	Name     string `form:"name"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

// InventoryQuery represents query parameters for inventory
type InventoryQuery struct {
	SparePartID *uint  `form:"spare_part_id"`
	FactoryID   *uint  `form:"factory_id"`
	LowStock    *bool  `form:"low_stock"`
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
}

// ConsumptionQuery represents query parameters for consumption records
type ConsumptionQuery struct {
	SparePartID *uint  `form:"spare_part_id"`
	OrderID     *uint  `form:"order_id"`
	TaskID      *uint  `form:"task_id"`
	DateFrom    string `form:"date_from"`
	DateTo      string `form:"date_to"`
	Page        int    `form:"page"`
	PageSize    int    `form:"page_size"`
}
