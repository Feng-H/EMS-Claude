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
	c.ShouldBindJSON(&req)
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
	c.JSON(200, gin.H{"total": len(memory.GetStore().Equipment), "running": len(memory.GetStore().Equipment)})
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

func GetMyTaskStatisticsMemory(c *gin.Context) { c.JSON(200, gin.H{"total": 5}) }
func GetInspectionStatisticsMemory(c *gin.Context) { c.JSON(200, gin.H{"total": 100}) }

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

func GetRepairStatisticsMemory(c *gin.Context) { c.JSON(200, gin.H{"total": 10}) }

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

func GetMaintenanceStatisticsMemory(c *gin.Context) { c.JSON(200, gin.H{"total": 20}) }

// ============ Spare Parts ============

func ListSparePartsMemory(c *gin.Context) {
	var res []*model.SparePart
	for _, p := range memory.GetStore().SpareParts { res = append(res, p) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func GetSparePartStatisticsMemory(c *gin.Context) { c.JSON(200, gin.H{"total": 50}) }

// ============ Analytics & Knowledge ============

func GetDashboardOverviewMemory(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) }

func ListKnowledgeArticlesMemory(c *gin.Context) {
	var res []*model.KnowledgeArticle
	for _, a := range memory.GetStore().KnowledgeArticles { res = append(res, a) }
	c.JSON(200, gin.H{"items": res, "total": len(res)})
}

func SearchKnowledgeArticlesMemory(c *gin.Context) {
	var res []gin.H
	for _, a := range memory.GetStore().KnowledgeArticles {
		res = append(res, gin.H{"id": a.ID, "title": a.Title})
	}
	c.JSON(200, res)
}

func HealthCheckMemory(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "mode": "memory"}) }

// Stub placeholders for routes used in main.go
func GetInspectionTemplateMemory(c *gin.Context) {}
func CreateInspectionTemplateMemory(c *gin.Context) {}
func CreateInspectionItemMemory(c *gin.Context) {}
func GetInspectionTaskMemory(c *gin.Context) {}
func StartInspectionMemory(c *gin.Context) {}
func CompleteInspectionMemory(c *gin.Context) {}
func GetRepairOrderMemory(c *gin.Context) {}
func CreateRepairOrderMemory(c *gin.Context) {}
func AssignRepairOrderMemory(c *gin.Context) {}
func StartRepairMemory(c *gin.Context) {}
func UpdateRepairMemory(c *gin.Context) {}
func ConfirmRepairMemory(c *gin.Context) {}
func AuditRepairMemory(c *gin.Context) {}
func GetMyRepairStatisticsMemory(c *gin.Context) {}
func CreateMaintenancePlanMemory(c *gin.Context) {}
func CreateMaintenanceItemMemory(c *gin.Context) {}
func GenerateMaintenanceTasksMemory(c *gin.Context) {}
func GetMaintenanceTaskMemory(c *gin.Context) {}
func StartMaintenanceMemory(c *gin.Context) {}
func CompleteMaintenanceMemory(c *gin.Context) {}
func UpdateEquipmentMemory(c *gin.Context) {}
func CreateEquipmentMemory(c *gin.Context) {}
func DeleteEquipmentMemory(c *gin.Context) {}
func CreateEquipmentTypeMemory(c *gin.Context) {}
func UpdateEquipmentTypeMemory(c *gin.Context) {}
func DeleteEquipmentTypeMemory(c *gin.Context) {}
func CreateBaseMemory(c *gin.Context) {}
func UpdateBaseMemory(c *gin.Context) {}
func DeleteBaseMemory(c *gin.Context) {}
func CreateFactoryMemory(c *gin.Context) {}
func UpdateFactoryMemory(c *gin.Context) {}
func DeleteFactoryMemory(c *gin.Context) {}
func CreateWorkshopMemory(c *gin.Context) {}
func UpdateWorkshopMemory(c *gin.Context) {}
func DeleteWorkshopMemory(c *gin.Context) {}
func UpdateSparePartMemory(c *gin.Context) {}
func CreateSparePartMemory(c *gin.Context) {}
func DeleteSparePartMemory(c *gin.Context) {}
func GetInventoryMemory(c *gin.Context) {}
func StockInMemory(c *gin.Context) {}
func StockOutMemory(c *gin.Context) {}
func GetLowStockAlertsMemory(c *gin.Context) {}
func GetConsumptionsMemory(c *gin.Context) {}
func CreateConsumptionMemory(c *gin.Context) {}
func GetMTTRMTBFMemory(c *gin.Context) {}
func GetTrendDataMemory(c *gin.Context) {}
func GetFailureAnalysisMemory(c *gin.Context) {}
func GetTopFailureEquipmentMemory(c *gin.Context) {}
func GetKnowledgeArticleMemory(c *gin.Context) {}
func CreateKnowledgeArticleMemory(c *gin.Context) {}
func UpdateKnowledgeArticleMemory(c *gin.Context) {}
func DeleteKnowledgeArticleMemory(c *gin.Context) {}
func ConvertFromRepairMemory(c *gin.Context) {}
