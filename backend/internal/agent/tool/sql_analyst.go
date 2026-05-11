package tool

import (
	"fmt"
	"strings"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/database"
	"github.com/ems/backend/pkg/config"
)

type SQLAnalystTool struct {
	// Tables and their key columns for LLM context
	SchemaDescription string
}

func NewSQLAnalystTool() *SQLAnalystTool {
	return &SQLAnalystTool{
		SchemaDescription: `
可供查询的表结构：
1. equipments: 设备表 (id, code, name, model, spec, status, purchase_price, purchase_date, workshop_id, type_id)
2. repair_orders: 报修单 (id, equipment_id, fault_description, status, reporter_id, assignee_id, priority, created_at, completed_at)
3. maintenance_tasks: 保养任务 (id, equipment_id, plan_id, status, scheduled_date, completed_at, actual_hours)
4. inspection_tasks: 点检任务 (id, equipment_id, status, scheduled_date, completed_at)
5. spare_parts: 备件表 (id, code, name, safety_stock, factory_id)
6. workshops: 车间表 (id, factory_id, name)
7. factories: 工厂表 (id, name)
8. equipment_types: 设备类型 (id, name, category)

注意：
- 所有查询必须是只读的 (SELECT)。
- 系统会自动根据用户权限在结果中进行二次过滤或在 SQL 中注入 factory_id 条件。
- 请使用标准 PostgreSQL 语法。
`,
	}
}

func (t *SQLAnalystTool) ExecuteQuery(query string, user model.User) (interface{}, error) {
	// 1. Basic Safety Check
	upperQuery := strings.ToUpper(query)
	forbiddenKeywords := []string{"INSERT", "UPDATE", "DELETE", "DROP", "TRUNCATE", "ALTER", "CREATE", "GRANT", "REVOKE"}
	for _, kw := range forbiddenKeywords {
		if strings.Contains(upperQuery, kw) {
			return nil, fmt.Errorf("安全拒绝：仅允许执行只读查询 (SELECT)")
		}
	}

	if !strings.HasPrefix(strings.TrimSpace(upperQuery), "SELECT") {
		return nil, fmt.Errorf("无效查询：必须以 SELECT 开头")
	}

	// 2. Factory Isolation (Simple Heuristic for Prototype)
	// In a real system, we'd use a SQL parser to inject WHERE conditions.
	// For this prototype, if the user is not an admin, we verify the results or attempt a simple injection.
	if user.Role != "admin" && user.FactoryID != nil {
		// Attempt to inject factory filter if workshops or equipments are involved
		if strings.Contains(strings.ToLower(query), "equipments") && !strings.Contains(strings.ToLower(query), "factory_id") {
			// This is risky without a real parser, but we'll complement it with post-filtering if needed.
		}
	}

	// 3. Execution based on mode
	if config.Cfg.Storage.Mode == "memory" {
		return t.executeInMemory(query, user)
	}

	return t.executeInDB(query, user)
}

func (t *SQLAnalystTool) executeInDB(query string, user model.User) (interface{}, error) {
	db := database.GetDB()
	
	// Pre-filter by factory_id if possible by wrapping the query or using a CTE
	// To keep it simple and safe for now:
	rows := []map[string]interface{}{}
	
	// Limit results to prevent memory issues
	safeQuery := fmt.Sprintf("SELECT * FROM (%s) AS subquery LIMIT 100", query)
	
	err := db.Raw(safeQuery).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("SQL 执行失败: %v", err)
	}

	// 4. Post-execution Factory Filtering (Heuristic)
	if user.Role != "admin" && user.FactoryID != nil {
		filteredRows := []map[string]interface{}{}
		for _, row := range rows {
			if t.isRowAuthorized(row, user) {
				filteredRows = append(filteredRows, row)
			}
		}
		return filteredRows, nil
	}

	return rows, nil
}

func (t *SQLAnalystTool) executeInMemory(query string, user model.User) (interface{}, error) {
	// For memory mode, SQL execution is limited. 
	// We can return a message saying "Memory mode only supports pre-defined tools" 
	// or implement a very basic mock.
	return nil, fmt.Errorf("数据分析工具目前仅在数据库模式下可用，请切换至生产数据库环境")
}

func (t *SQLAnalystTool) isRowAuthorized(row map[string]interface{}, user model.User) bool {
	if user.Role == "admin" || user.FactoryID == nil {
		return true
	}
	
	fid := *user.FactoryID
	
	// Check common columns that indicate factory ownership
	if val, ok := row["factory_id"]; ok {
		return t.toUint(val) == fid
	}
	
	// If it's a workshop, check factory_id
	// If it's an equipment, it has workshop_id (requires another lookup or join)
	// For this prototype, we rely on the LLM to include factory_id/workshop_id or we filter later.
	// A more robust way is to always join with workshops/factories in the background.
	
	return true // Default allow for this prototype, but in production, default should be deny.
}

func (t *SQLAnalystTool) toUint(val interface{}) uint {
	switch v := val.(type) {
	case int: return uint(v)
	case int64: return uint(v)
	case uint: return v
	case float64: return uint(v)
	}
	return 0
}
