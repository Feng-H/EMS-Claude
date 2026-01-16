package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
)

// RepairOrderService
type RepairOrderService struct {
	orderRepo *repository.RepairOrderRepository
	logRepo   *repository.RepairLogRepository
	equipRepo *repository.EquipmentRepository
	userRepo  *UserRepository
}

func NewRepairOrderService() *RepairOrderService {
	return &RepairOrderService{
		orderRepo: repository.NewRepairOrderRepository(),
		logRepo:   repository.NewRepairLogRepository(),
		equipRepo: repository.NewEquipmentRepo(),
		userRepo:  NewUserRepository(),
	}
}

// CreateOrderRequest creates a new repair order
type CreateOrderRequest struct {
	EquipmentID      uint
	FaultDescription string
	FaultCode        string
	Photos           []string
	Priority         int
	ReporterID       uint
}

func (s *RepairOrderService) CreateOrder(req *CreateOrderRequest) (*model.RepairOrder, error) {
	// Verify equipment exists
	equipment, err := s.equipRepo.GetByID(req.EquipmentID)
	if err != nil {
		return nil, ErrNotFound
	}

	// Set default priority
	priority := req.Priority
	if priority == 0 {
		priority = 2 // default to medium
	}

	order := &model.RepairOrder{
		EquipmentID:      req.EquipmentID,
		FaultDescription: req.FaultDescription,
		FaultCode:        req.FaultCode,
		Photos:           req.Photos,
		Priority:         priority,
		ReporterID:       req.ReporterID,
		Status:           model.RepairPending,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Update equipment status to maintenance
	equipment.Status = "maintenance"
	s.equipRepo.Update(equipment)

	// Create log
	s.createLog(order.ID, req.ReporterID, "created", "创建维修工单")

	return s.orderRepo.GetByID(order.ID)
}

// GetByID returns a repair order by ID
func (s *RepairOrderService) GetByID(id uint) (*model.RepairOrder, error) {
	return s.orderRepo.GetByID(id)
}

// List returns repair orders with filtering
type RepairOrderFilter struct {
	Status     string
	Priority   int
	AssignedTo uint
	DateFrom   time.Time
	DateTo     time.Time
	Page       int
	PageSize   int
}

func (s *RepairOrderService) List(filter *RepairOrderFilter) (*RepairOrderListResult, error) {
	repoFilter := repository.RepairOrderFilter{
		Status:     filter.Status,
		Priority:   filter.Priority,
		AssignedTo: filter.AssignedTo,
		DateFrom:   filter.DateFrom,
		DateTo:     filter.DateTo,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}

	orders, total, err := s.orderRepo.List(repoFilter)
	if err != nil {
		return nil, err
	}

	return &RepairOrderListResult{
		Items: orders,
		Total: total,
	}, nil
}

// AssignOrder assigns a repair order to a technician
func (s *RepairOrderService) AssignOrder(orderID uint, assignTo uint, assignedBy uint) (*model.RepairOrder, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, ErrNotFound
	}

	if order.Status != model.RepairPending {
		return nil, errors.New("only pending orders can be assigned")
	}

	// Verify assignee exists
	if _, err := s.userRepo.GetByID(assignTo); err != nil {
		return nil, ErrNotFound
	}

	order.AssignedTo = &assignTo
	order.Status = model.RepairAssigned
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}

	// Create log
	s.createLog(orderID, assignedBy, "assigned", fmt.Sprintf("指派给维修工 #%d", assignTo))

	return s.orderRepo.GetByID(orderID)
}

// AutoAssign automatically assigns pending orders to available technicians
func (s *RepairOrderService) AutoAssign(limit int) (*AutoAssignResult, error) {
	pendingOrders, err := s.orderRepo.GetPendingOrders()
	if err != nil {
		return nil, err
	}

	result := &AutoAssignResult{
		AssignedCount: 0,
		OrderIDs:      []uint{},
	}

	// Simple round-robin assignment based on current workload
	// In production, this should consider skills, location, etc.
	var technicians []uint
	// Get users with maintenance role (simplified)
	// TODO: Query users by role

	for i, order := range pendingOrders {
		if i >= limit {
			break
		}

		if len(technicians) == 0 {
			// No technicians available
			continue
		}

		// Assign to technician with least workload
		assignee := s.findLeastBusyTechnician(technicians)
		if assignee > 0 {
			_, err := s.AssignOrder(order.ID, assignee, 1) // 1 for system
			if err == nil {
				result.AssignedCount++
				result.OrderIDs = append(result.OrderIDs, order.ID)
			}
		}
	}

	return result, nil
}

// StartRepair starts a repair task
func (s *RepairOrderService) StartRepair(orderID uint, userID uint) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != model.RepairAssigned {
		return errors.New("order is not in assignable state")
	}

	now := time.Now()
	order.Status = model.RepairInProgress
	order.StartedAt = &now

	if err := s.orderRepo.Update(order); err != nil {
		return err
	}

	s.createLog(orderID, userID, "started", "开始维修")

	return nil
}

// UpdateRepair updates repair progress
type UpdateRepairRequest struct {
	Solution    string
	SpareParts  string
	ActualHours float64
	Photos      []string
	NextStatus  string // testing, confirmed
}

func (s *RepairOrderService) UpdateRepair(orderID uint, userID uint, req *UpdateRepairRequest) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != model.RepairInProgress && order.Status != model.RepairTesting {
		return errors.New("order is not in updatable state")
	}

	order.Solution = req.Solution
	if len(req.Photos) > 0 {
		order.Photos = req.Photos
	}

	// Handle status transition
	var statusChanged bool
	var newStatus model.RepairStatus
	var logContent string

	switch req.NextStatus {
	case "testing":
		if order.Status == model.RepairInProgress {
			newStatus = model.RepairTesting
			logContent = "等待报修人确认"
			statusChanged = true
		}
	case "confirmed":
		if order.Status == model.RepairInProgress || order.Status == model.RepairTesting {
			now := time.Now()
			newStatus = model.RepairConfirmed
			order.CompletedAt = &now
			logContent = "维修完成，等待审核"
			statusChanged = true
		}
	}

	if statusChanged {
		order.Status = newStatus
	}

	if err := s.orderRepo.Update(order); err != nil {
		return err
	}

	logAction := "updated"
	if statusChanged {
		logAction = "status_changed"
	}
	s.createLog(orderID, userID, logAction, logContent)

	return nil
}

// ConfirmRepair confirms or rejects a completed repair
func (s *RepairOrderService) ConfirmRepair(orderID uint, userID uint, accepted bool, comment string, photos []string) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != model.RepairTesting && order.Status != model.RepairConfirmed {
		return errors.New("order is not ready for confirmation")
	}

	if accepted {
		now := time.Now()
		order.Status = model.RepairAudited
		order.ConfirmedAt = &now
		if order.CompletedAt == nil {
			order.CompletedAt = &now
		}
		s.createLog(orderID, userID, "confirmed", "确认维修完成")
	} else {
		// Reject - send back to in_progress
		order.Status = model.RepairInProgress
		s.createLog(orderID, userID, "rejected", "确认不通过: "+comment)
	}

	if len(photos) > 0 {
		order.Photos = photos
	}

	return s.orderRepo.Update(order)
}

// AuditRepair audits a confirmed repair (supervisor/engineer)
func (s *RepairOrderService) AuditRepair(orderID uint, userID uint, approved bool, comment string, actualHours *float64) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != model.RepairAudited {
		return errors.New("order is not ready for audit")
	}

	if approved {
		now := time.Now()
		order.Status = model.RepairClosed
		order.AuditedAt = &now
		if actualHours != nil {
			// Store actual hours in solution field for now
			order.Solution = fmt.Sprintf("%s (Actual hours: %.1f)", order.Solution, *actualHours)
		}

		// Update equipment status back to running
		if equipment, err := s.equipRepo.GetByID(order.EquipmentID); err == nil {
			equipment.Status = "running"
			s.equipRepo.Update(equipment)
		}

		s.createLog(orderID, userID, "audited", "审核通过，工单关闭")
	} else {
		// Reject - send back to in_progress
		order.Status = model.RepairInProgress
		order.ConfirmedAt = nil
		s.createLog(orderID, userID, "audit_rejected", "审核不通过: "+comment)
	}

	return s.orderRepo.Update(order)
}

// GetStatistics returns repair statistics
func (s *RepairOrderService) GetStatistics() (*RepairStatistics, error) {
	stats, err := s.orderRepo.GetStatistics()
	if err != nil {
		return nil, err
	}

	avgRepairTime, _ := s.orderRepo.GetAvgRepairTime()
	avgResponseTime, _ := s.orderRepo.GetAvgResponseTime()

	return &RepairStatistics{
		TotalOrders:      stats["total"],
		PendingOrders:    stats["pending"],
		InProgressOrders: stats["in_progress"],
		CompletedOrders:  stats["completed"],
		TodayCompleted:   stats["today_completed"],
		TodayCreated:     stats["today_created"],
		AvgRepairTime:    avgRepairTime,
		AvgResponseTime:  avgResponseTime,
	}, nil
}

// GetUserStatistics returns current user's repair statistics
func (s *RepairOrderService) GetUserStatistics(userID uint) (*MyRepairStatistics, error) {
	stats, err := s.orderRepo.GetUserStatistics(userID)
	if err != nil {
		return nil, err
	}

	return &MyRepairStatistics{
		PendingCount:    int(stats["pending"]),
		InProgressCount: int(stats["in_progress"]),
		CompletedCount:  int(stats["completed"]),
		TodayCompleted:  int(stats["today_completed"]),
	}, nil
}

// GetMyTasks returns current user's active repair tasks
func (s *RepairOrderService) GetMyTasks(userID uint) ([]model.RepairOrder, error) {
	return s.orderRepo.GetByAssignee(userID)
}

// Helper functions
func (s *RepairOrderService) createLog(orderID uint, userID uint, action, content string) {
	log := &model.RepairLog{
		OrderID: orderID,
		UserID:  userID,
		Action:  action,
		Content: content,
	}
	s.logRepo.Create(log)
}

func (s *RepairOrderService) findLeastBusyTechnician(technicians []uint) uint {
	// Find technician with least active orders
	minCount := int64(999999)
	var selected uint

	for _, techID := range technicians {
		stats, _ := s.orderRepo.GetUserStatistics(techID)
		activeCount := stats["pending"] + stats["in_progress"]
		if activeCount < minCount {
			minCount = activeCount
			selected = techID
		}
	}

	return selected
}

// Types
type RepairOrderListResult struct {
	Items []model.RepairOrder
	Total int64
}

type RepairStatistics struct {
	TotalOrders      int64
	PendingOrders    int64
	InProgressOrders int64
	CompletedOrders  int64
	TodayCompleted   int64
	TodayCreated     int64
	AvgRepairTime    float64
	AvgResponseTime  float64
}

type MyRepairStatistics struct {
	PendingCount    int
	InProgressCount int
	CompletedCount  int
	TodayCompleted  int
}

type AutoAssignResult struct {
	AssignedCount int
	OrderIDs      []uint
}
