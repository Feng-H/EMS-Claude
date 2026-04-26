package v1

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/memory"
	"golang.org/x/crypto/bcrypt"
)

var (
	store *memory.Store
)

func InitMemory() {
	store = memory.GetStore()
	store.InitMockData()
}

// ============ User & Auth ============

func GetUsersMemory(c *gin.Context) {
	s := memory.GetStore()
	var res []gin.H
	for _, u := range s.Users {
		res = append(res, gin.H{
			"id": u.ID, "username": u.Username, "name": u.Name, "role": string(u.Role),
			"is_active": u.IsActive, "approval_status": string(u.ApprovalStatus),
			"factory_id": u.FactoryID, "created_at": u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func CreateUserMemory(c *gin.Context) {
	var req struct {
		Username string `json:"username"`; Password string `json:"password"`; Name string `json:"name"`; Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	u := &model.User{
		BaseModel: model.BaseModel{ID: memory.GetStore().NextID(), CreatedAt: time.Now()},
		Username: req.Username, PasswordHash: string(h), Name: req.Name, Role: model.UserRole(req.Role),
		IsActive: true, ApprovalStatus: model.ApprovalStatusApproved,
	}
	memory.GetStore().AddUser(u.ID, u)
	c.JSON(201, gin.H{"id": u.ID})
}

func GetPendingApplicationsMemory(c *gin.Context) {
	var res []gin.H
	for _, u := range memory.GetStore().Users {
		if u.ApprovalStatus == model.ApprovalStatusPending {
			res = append(res, gin.H{"id": u.ID, "username": u.Username, "name": u.Name})
		}
	}
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func UpdateUserMemory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	c.JSON(200, gin.H{"id": id, "message": "Updated"})
}

func ApproveApplicationMemory(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Approved"})
}

func LogoutMemory(c *gin.Context) { c.JSON(200, gin.H{"message": "ok"}) }

// ============ Organization ============

func ListBasesMemory(c *gin.Context) {
	var res []*model.Base
	for _, b := range memory.GetStore().Bases { res = append(res, b) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func ListFactoriesMemory(c *gin.Context) {
	var res []*model.Factory
	for _, f := range memory.GetStore().Factories { res = append(res, f) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func ListWorkshopsMemory(c *gin.Context) {
	var res []*model.Workshop
	for _, w := range memory.GetStore().Workshops { res = append(res, w) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

// ============ Equipment ============

func ListEquipmentMemory(c *gin.Context) {
	var res []*model.Equipment
	for _, e := range memory.GetStore().Equipment { res = append(res, e) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func GetEquipmentMemory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	e := memory.GetStore().FindEquipment(uint(id))
	if e == nil { c.JSON(404, gin.H{"error": "404"}); return }
	c.JSON(200, e)
}

func GetEquipmentByQRCodeMemory(c *gin.Context) {
	code := c.Param("code")
	for _, e := range memory.GetStore().Equipment {
		if e.QRCode == code { c.JSON(200, e); return }
	}
	c.JSON(404, gin.H{"error": "404"})
}

func GetEquipmentStatisticsMemory(c *gin.Context) {
	s := memory.GetStore()
	total := len(s.Equipment)
	running, stopped, maintenance, scrapped := 0, 0, 0, 0
	for _, e := range s.Equipment {
		switch e.Status {
		case "running": running++
		case "stopped": stopped++
		case "maintenance": maintenance++
		case "scrapped": scrapped++
		default: running++
		}
	}
	c.JSON(200, gin.H{
		"total": total, "running": running, "stopped": stopped, "maintenance": maintenance, "scrapped": scrapped,
	})
}

func ListEquipmentTypesMemory(c *gin.Context) {
	var res []*model.EquipmentType
	for _, t := range memory.GetStore().EquipmentTypes { res = append(res, t) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

// ============ Inspection ============

func ListInspectionTemplatesMemory(c *gin.Context) {
	var res []*model.InspectionTemplate
	for _, t := range memory.GetStore().InspectionTemplates { res = append(res, t) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func ListInspectionTasksMemory(c *gin.Context) {
	var res []*model.InspectionTask
	for _, t := range memory.GetStore().InspectionTasks { res = append(res, t) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func GetMyTasksMemory(c *gin.Context) {
	var res []*model.InspectionTask
	for _, t := range memory.GetStore().InspectionTasks { res = append(res, t) }
	c.JSON(200, res)
}

func GetMyTaskStatisticsMemory(c *gin.Context) {
	c.JSON(200, gin.H{
		"pending_count":     2,
		"in_progress_count": 1,
		"today_tasks":       5,
	})
}

func GetInspectionStatisticsMemory(c *gin.Context) {
	c.JSON(200, gin.H{
		"total_tasks":       100,
		"pending_tasks":     20,
		"in_progress_tasks": 10,
		"completed_tasks":   60,
		"overdue_tasks":     10,
		"today_completed":   5,
		"completion_rate":   60.0,
	})
}

// ============ Repair ============

func ListRepairOrdersMemory(c *gin.Context) {
	var res []*model.RepairOrder
	for _, o := range memory.GetStore().RepairOrders { res = append(res, o) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func GetMyRepairTasksMemory(c *gin.Context) {
	var res []*model.RepairOrder
	for _, o := range memory.GetStore().RepairOrders { res = append(res, o) }
	c.JSON(200, res)
}

func GetRepairStatisticsMemory(c *gin.Context) {
	c.JSON(200, gin.H{
		"total_orders":       10,
		"pending_orders":     2,
		"in_progress_orders": 3,
		"completed_orders":   5,
		"today_completed":    1,
		"today_created":      2,
		"avg_repair_time":    120.0,
		"avg_response_time":  30.0,
	})
}

func GetMyRepairStatisticsMemory(c *gin.Context) {
	c.JSON(200, gin.H{
		"total_orders":       5,
		"pending_orders":     1,
		"in_progress_orders": 2,
		"completed_orders":   2,
	})
}

// ============ Maintenance ============

func ListMaintenancePlansMemory(c *gin.Context) {
	var res []*model.MaintenancePlan
	for _, p := range memory.GetStore().MaintenancePlans { res = append(res, p) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func ListMaintenanceTasksMemory(c *gin.Context) {
	var res []*model.MaintenanceTask
	for _, t := range memory.GetStore().MaintenanceTasks { res = append(res, t) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func GetMyMaintenanceTasksMemory(c *gin.Context) {
	var res []*model.MaintenanceTask
	for _, t := range memory.GetStore().MaintenanceTasks { res = append(res, t) }
	c.JSON(200, res)
}

func GetMaintenanceStatisticsMemory(c *gin.Context) {
	c.JSON(200, gin.H{
		"total_plans":       5,
		"total_tasks":       20,
		"pending_tasks":     5,
		"in_progress_tasks": 2,
		"completed_tasks":   10,
		"overdue_tasks":     3,
		"today_completed":   2,
		"completion_rate":   50.0,
	})
}

// ============ Spare Parts ============

func ListSparePartsMemory(c *gin.Context) {
	var res []*model.SparePart
	for _, p := range memory.GetStore().SpareParts { res = append(res, p) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func GetSparePartStatisticsMemory(c *gin.Context) {
	c.JSON(200, gin.H{
		"total_parts":      50,
		"low_stock_count":  5,
		"total_stock_value": 150000.0,
	})
}

// ============ Analytics & Knowledge ============

func GetDashboardOverviewMemory(c *gin.Context) {
	c.JSON(200, gin.H{
		"equipment_total": 4,
		"health_index":    92.5,
		"uptime_rate":     98.2,
		"repair_trend":    []int{2, 3, 1, 4, 2, 1, 0},
	})
}

func ListKnowledgeArticlesMemory(c *gin.Context) {
	var res []*model.KnowledgeArticle
	for _, a := range memory.GetStore().KnowledgeArticles { 
		res = append(res, a) 
	}
	c.JSON(200, res)
}

func SearchKnowledgeArticlesMemory(c *gin.Context) {
	var res []gin.H
	for _, a := range memory.GetStore().KnowledgeArticles {
		res = append(res, gin.H{"id": a.ID, "title": a.Title})
	}
	c.JSON(200, res)
}

func HealthCheckMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "mode": "memory"}) }

// Stub placeholders for missing logic to prevent 404/500
func GetInspectionTemplateMemory(c *gin.Context) { c.JSON(200, gin.H{}) }
func CreateInspectionTemplateMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func CreateInspectionItemMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func GetInspectionTaskMemory(c *gin.Context) { c.JSON(200, gin.H{}) }
func StartInspectionMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "started"}) }
func CompleteInspectionMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "completed"}) }
func GetRepairOrderMemory(c *gin.Context) { c.JSON(200, gin.H{}) }
func CreateRepairOrderMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func AssignRepairOrderMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "assigned"}) }
func StartRepairMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "started"}) }
func UpdateRepairMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "updated"}) }
func ConfirmRepairMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "confirmed"}) }
func AuditRepairMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "audited"}) }
func CreateMaintenancePlanMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func CreateMaintenanceItemMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func GenerateMaintenanceTasksMemory(c *gin.Context) { c.JSON(201, gin.H{"count": 5}) }
func GetMaintenanceTaskMemory(c *gin.Context) { c.JSON(200, gin.H{}) }
func StartMaintenanceMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "started"}) }
func CompleteMaintenanceMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "completed"}) }
func UpdateEquipmentMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "updated"}) }
func CreateEquipmentMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func DeleteEquipmentMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "deleted"}) }
func CreateEquipmentTypeMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func UpdateEquipmentTypeMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "updated"}) }
func DeleteEquipmentTypeMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "deleted"}) }
func CreateBaseMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func UpdateBaseMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "updated"}) }
func DeleteBaseMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "deleted"}) }
func CreateFactoryMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func UpdateFactoryMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "updated"}) }
func DeleteFactoryMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "deleted"}) }
func CreateWorkshopMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func UpdateWorkshopMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "updated"}) }
func DeleteWorkshopMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "deleted"}) }
func UpdateSparePartMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "updated"}) }
func CreateSparePartMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func DeleteSparePartMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "deleted"}) }
func GetInventoryMemory(c *gin.Context) { c.JSON(200, gin.H{"items": []any{}}) }
func StockInMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "in"}) }
func StockOutMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "out"}) }
func GetLowStockAlertsMemory(c *gin.Context) { c.JSON(200, []any{}) }
func GetConsumptionsMemory(c *gin.Context) { c.JSON(200, gin.H{"items": []any{}}) }
func CreateConsumptionMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func GetMTTRMTBFMemory(c *gin.Context) { c.JSON(200, gin.H{"mttr": 120, "mtbf": 720}) }
func GetTrendDataMemory(c *gin.Context) { c.JSON(200, []any{}) }
func GetFailureAnalysisMemory(c *gin.Context) { c.JSON(200, []any{}) }
func GetTopFailureEquipmentMemory(c *gin.Context) { c.JSON(200, []any{}) }
func GetKnowledgeArticleMemory(c *gin.Context) { c.JSON(200, gin.H{}) }
func CreateKnowledgeArticleMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
func UpdateKnowledgeArticleMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "updated"}) }
func DeleteKnowledgeArticleMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "deleted"}) }
func ConvertFromRepairMemory(c *gin.Context) { c.JSON(201, gin.H{}) }
