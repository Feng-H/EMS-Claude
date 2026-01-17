package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ems/backend/api/v1"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/database"
	"github.com/ems/backend/pkg/redis"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title EMS Equipment Management System API
// @version 1.0
// @description Equipment Management System API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load configuration
	if err := config.Load("config/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Initialize Redis (optional - log warning if unavailable)
	if err := redis.Init(); err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v (caching disabled)", err)
	} else {
		defer redis.Close()
		log.Println("Redis connected successfully")
	}

	// Auto migrate tables
	if err := database.GetDB().AutoMigrate(
		&model.Base{},
		&model.Factory{},
		&model.Workshop{},
		&model.User{},
		&model.EquipmentType{},
		&model.Equipment{},
		&model.InspectionTemplate{},
		&model.InspectionItem{},
		&model.InspectionTask{},
		&model.InspectionRecord{},
		&model.RepairOrder{},
		&model.RepairLog{},
		&model.MaintenancePlan{},
		&model.MaintenancePlanItem{},
		&model.MaintenanceTask{},
		&model.MaintenanceRecord{},
		&model.SparePart{},
		&model.SparePartInventory{},
		&model.SparePartConsumption{},
		&model.KnowledgeArticle{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize API handlers
	v1.InitAuth(database.GetDB())
	v1.InitEquipment(database.GetDB())
	v1.InitInspection()
	v1.InitRepair()
	v1.InitMaintenance()
	v1.InitSparePart()
	v1.InitAnalytics()
	v1.InitKnowledge()

	// Setup Gin
	gin.SetMode(config.Cfg.Server.Mode)
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", v1.Login)
			auth.POST("/refresh", v1.RefreshToken)
			auth.POST("/apply", v1.ApplyForAccount)
			auth.GET("/me", middleware.AuthMiddleware(), v1.GetCurrentUser)
			auth.POST("/change-password", middleware.AuthMiddleware(), v1.ChangePassword)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User management routes (admin only)
			users := protected.Group("/users")
			{
				users.GET("", v1.GetUsers)
				users.POST("", v1.CreateUser)
				users.GET("/applications", v1.GetPendingApplications)
				users.PUT("/:id/approve", v1.ApproveApplication)
				users.PUT("/:id", v1.UpdateUser)
			}

			// Organization routes
			org := protected.Group("/organization")
			{
				org.GET("/bases", v1.ListBases)
				org.POST("/bases", v1.CreateBase)
				org.PUT("/bases/:id", v1.UpdateBase)
				org.DELETE("/bases/:id", v1.DeleteBase)

				org.GET("/factories", v1.ListFactories)
				org.POST("/factories", v1.CreateFactory)
				org.PUT("/factories/:id", v1.UpdateFactory)
				org.DELETE("/factories/:id", v1.DeleteFactory)

				org.GET("/workshops", v1.ListWorkshops)
				org.POST("/workshops", v1.CreateWorkshop)
				org.PUT("/workshops/:id", v1.UpdateWorkshop)
				org.DELETE("/workshops/:id", v1.DeleteWorkshop)
			}

			// Equipment routes
			equipment := protected.Group("/equipment")
			{
				equipment.GET("", v1.ListEquipment)
				equipment.GET("/:id", v1.GetEquipment)
				equipment.GET("/qr/:code", v1.GetEquipmentByQRCode)
				equipment.POST("", v1.CreateEquipment)
				equipment.PUT("/:id", v1.UpdateEquipment)
				equipment.DELETE("/:id", v1.DeleteEquipment)
				equipment.GET("/statistics", v1.GetEquipmentStatistics)

				// Equipment types
				equipment.GET("/types", v1.ListEquipmentTypes)
				equipment.POST("/types", v1.CreateEquipmentType)
			}

			// Inspection routes
			inspection := protected.Group("/inspection")
			{
				inspection.GET("/templates", v1.ListInspectionTemplates)
				inspection.GET("/templates/:id", v1.GetInspectionTemplate)
				inspection.POST("/templates", v1.CreateInspectionTemplate)
				inspection.POST("/items", v1.CreateInspectionItem)

				inspection.GET("/tasks", v1.ListInspectionTasks)
				inspection.GET("/tasks/:id", v1.GetInspectionTask)
				inspection.GET("/my-tasks", v1.GetMyTasks)
				inspection.GET("/my-stats", v1.GetMyTaskStatistics)
				inspection.POST("/start", v1.StartInspection)
				inspection.POST("/complete", v1.CompleteInspection)
				inspection.GET("/statistics", v1.GetInspectionStatistics)
			}

			// Repair routes
			repair := protected.Group("/repair")
			{
				repair.GET("/orders", v1.ListRepairOrders)
				repair.POST("/orders", v1.CreateRepairOrder)
				repair.GET("/orders/:id", v1.GetRepairOrder)
				repair.POST("/orders/:id/assign", v1.AssignRepairOrder)
				repair.POST("/orders/:id/start", v1.StartRepair)
				repair.POST("/orders/:id/update", v1.UpdateRepair)
				repair.POST("/orders/:id/confirm", v1.ConfirmRepair)
				repair.POST("/orders/:id/audit", v1.AuditRepair)
				repair.GET("/my-tasks", v1.GetMyRepairTasks)
				repair.GET("/my-stats", v1.GetMyRepairStatistics)
				repair.GET("/statistics", v1.GetRepairStatistics)
			}

			// Maintenance routes
			maintenance := protected.Group("/maintenance")
			{
				// Plan management
				maintenance.GET("/plans", v1.ListMaintenancePlans)
				maintenance.POST("/plans", v1.CreateMaintenancePlan)
				maintenance.POST("/items", v1.CreateMaintenanceItem)

				// Task management
				maintenance.POST("/tasks/generate", v1.GenerateMaintenanceTasks)
				maintenance.GET("/tasks", v1.ListMaintenanceTasks)
				maintenance.GET("/tasks/:id", v1.GetMaintenanceTask)
				maintenance.GET("/my-tasks", v1.GetMyMaintenanceTasks)

				// Execution
				maintenance.POST("/start", v1.StartMaintenance)
				maintenance.POST("/complete", v1.CompleteMaintenance)

				// Statistics
				maintenance.GET("/statistics", v1.GetMaintenanceStatistics)
			}

			// Spare parts routes
			spareparts := protected.Group("/spareparts")
			{
				// Part management
				spareparts.GET("", v1.ListSpareParts)
				spareparts.POST("", v1.CreateSparePart)
				spareparts.PUT("/:id", v1.UpdateSparePart)
				spareparts.DELETE("/:id", v1.DeleteSparePart)

				// Inventory
				spareparts.GET("/inventory", v1.GetInventory)
				spareparts.POST("/stock-in", v1.StockIn)
				spareparts.POST("/stock-out", v1.StockOut)
				spareparts.GET("/alerts", v1.GetLowStockAlerts)

				// Consumption
				spareparts.GET("/consumptions", v1.GetConsumptionRecords)
				spareparts.POST("/consumptions", v1.CreateConsumption)

				// Statistics
				spareparts.GET("/statistics", v1.GetSparePartStatistics)
			}

			// Analytics routes
			analytics := protected.Group("/analytics")
			{
				analytics.GET("/dashboard", v1.GetDashboardOverview)
				analytics.GET("/mttr-mtbf", v1.GetMTTRMTBF)
				analytics.GET("/trends", v1.GetTrendData)
				analytics.GET("/failures", v1.GetFailureAnalysis)
				analytics.GET("/top-failures", v1.GetTopFailureEquipment)
			}

			// Knowledge base routes
			knowledge := protected.Group("/knowledge")
			{
				knowledge.GET("", v1.ListKnowledgeArticles)
				knowledge.GET("/:id", v1.GetKnowledgeArticle)
				knowledge.POST("", v1.CreateKnowledgeArticle)
				knowledge.PUT("/:id", v1.UpdateKnowledgeArticle)
				knowledge.DELETE("/:id", v1.DeleteKnowledgeArticle)
				knowledge.GET("/search", v1.SearchKnowledgeArticles)
				knowledge.POST("/convert-repair", v1.ConvertFromRepair)
			}
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ems-api"})
	})

	// Start server
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("Server starting on port %d", config.Cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
