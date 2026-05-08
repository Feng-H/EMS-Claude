package service

import (
	"fmt"
	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/database"
)

// ListTools returns the list of available tools for external Agents
func (s *AgentService) ListTools(user model.User) ([]dto.ToolDefinition, error) {
	tools := []dto.ToolDefinition{
		{
			Name:        "search_equipment",
			Description: "Search for equipment by name, code or model",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"keyword": map[string]interface{}{
						"type":        "string",
						"description": "Search keyword (name, code, or model)",
					},
				},
				"required": []string{"keyword"},
			},
		},
		{
			Name:        "get_equipment_health",
			Description: "Get real-time health analysis and remaining useful life (RUL) prediction for a specific equipment",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"equipment_id": map[string]interface{}{
						"type":        "integer",
						"description": "The unique ID of the equipment",
					},
				},
				"required": []string{"equipment_id"},
			},
		},
		{
			Name:        "get_spare_part_inventory",
			Description: "Check stock levels of spare parts across different factories",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"spare_part_id": map[string]interface{}{
						"type":        "integer",
						"description": "The unique ID of the spare part",
					},
					"factory_id": map[string]interface{}{
						"type":        "integer",
						"description": "Optional: Filter by factory ID",
					},
				},
				"required": []string{"spare_part_id"},
			},
		},
		{
			Name:        "report_repair",
			Description: "Submit a new repair request for a faulty equipment",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"equipment_id": map[string]interface{}{
						"type":        "integer",
						"description": "The unique ID of the equipment",
					},
					"fault_description": map[string]interface{}{
						"type":        "string",
						"description": "Detailed description of the fault",
					},
					"priority": map[string]interface{}{
						"type":        "integer",
						"description": "Priority: 1=High, 2=Medium, 3=Low",
					},
				},
				"required": []string{"equipment_id", "fault_description"},
			},
		},
	}

	return tools, nil
}

// CallTool executes a tool call for an external Agent
func (s *AgentService) CallTool(user model.User, req *dto.CallToolRequest) (*dto.CallToolResponse, error) {
	db := database.GetDB()

	switch req.Name {
	case "search_equipment":
		keyword, ok := req.Arguments["keyword"].(string)
		if !ok {
			return nil, fmt.Errorf("missing or invalid 'keyword' argument")
		}
		
		var equipments []model.Equipment
		query := db.Preload("Workshop").Preload("Workshop.Factory")
		if user.Role != "admin" && user.FactoryID != nil {
			query = query.Joins("JOIN workshops ON equipments.workshop_id = workshops.id").
				Where("workshops.factory_id = ?", *user.FactoryID)
		}
		
		err := query.Where("equipments.name ILIKE ? OR equipments.code ILIKE ?", "%"+keyword+"%", "%"+keyword+"%").
			Limit(10).Find(&equipments).Error
		if err != nil {
			return &dto.CallToolResponse{Content: err.Error(), IsError: true}, nil
		}
		
		return &dto.CallToolResponse{Content: equipments, IsError: false}, nil

	case "get_equipment_health":
		idVal, ok := req.Arguments["equipment_id"]
		if !ok {
			return nil, fmt.Errorf("missing 'equipment_id' argument")
		}
		var id uint
		switch v := idVal.(type) {
		case float64: id = uint(v)
		case int: id = uint(v)
		default: return nil, fmt.Errorf("invalid 'equipment_id' type")
		}
		
		result, err := s.GetEquipmentPrediction(id, user)
		if err != nil {
			return &dto.CallToolResponse{Content: err.Error(), IsError: true}, nil
		}
		return &dto.CallToolResponse{Content: result, IsError: false}, nil

	case "get_spare_part_inventory":
		idVal, ok := req.Arguments["spare_part_id"]
		if !ok {
			return nil, fmt.Errorf("missing 'spare_part_id' argument")
		}
		var partID uint
		switch v := idVal.(type) {
		case float64: partID = uint(v)
		case int: partID = uint(v)
		default: return nil, fmt.Errorf("invalid 'spare_part_id' type")
		}

		var inventories []model.SparePartInventory
		query := db.Preload("Factory").Preload("SparePart").Where("spare_part_id = ?", partID)
		
		// If user is not admin, they can only see inventory in their own factory or factories they have access to.
		// For MVP, we'll restrict to their assigned FactoryID if they are not admin.
		if user.Role != "admin" && user.FactoryID != nil {
			query = query.Where("factory_id = ?", *user.FactoryID)
		} else if fid, ok := req.Arguments["factory_id"]; ok {
			query = query.Where("factory_id = ?", fid)
		}
		
		if err := query.Find(&inventories).Error; err != nil {
			return &dto.CallToolResponse{Content: err.Error(), IsError: true}, nil
		}
		return &dto.CallToolResponse{Content: inventories, IsError: false}, nil

	case "report_repair":
		equipIDVal, ok := req.Arguments["equipment_id"]
		desc, ok2 := req.Arguments["fault_description"].(string)
		if !ok || !ok2 {
			return nil, fmt.Errorf("missing required arguments for report_repair")
		}
		
		var equipID uint
		switch v := equipIDVal.(type) {
		case float64: equipID = uint(v)
		case int: equipID = uint(v)
		default: return nil, fmt.Errorf("invalid 'equipment_id' type")
		}

		// Ownership check: Ensure equipment belongs to user's factory
		var equipment model.Equipment
		if err := db.Joins("JOIN workshops ON workshops.id = equipments.workshop_id").
			First(&equipment, equipID).Error; err != nil {
			return &dto.CallToolResponse{Content: "Equipment not found", IsError: true}, nil
		}
		
		if user.Role != "admin" && user.FactoryID != nil {
			var workshop model.Workshop
			db.First(&workshop, equipment.WorkshopID)
			if workshop.FactoryID != *user.FactoryID {
				return &dto.CallToolResponse{Content: "Permission denied: Equipment belongs to another factory", IsError: true}, nil
			}
		}

		priority := 2
		if p, ok := req.Arguments["priority"].(float64); ok {
			priority = int(p)
		}

		order := model.RepairOrder{
			EquipmentID:      equipID,
			FaultDescription: desc,
			Priority:         priority,
			Status:           model.RepairPending,
			ReporterID:       user.ID,
		}

		if err := db.Create(&order).Error; err != nil {
			return &dto.CallToolResponse{Content: err.Error(), IsError: true}, nil
		}

		return &dto.CallToolResponse{
			Content: fmt.Sprintf("Repair order #%d created successfully", order.ID),
			IsError: false,
		}, nil

	default:
		return nil, fmt.Errorf("unknown tool: %s", req.Name)
	}
}
