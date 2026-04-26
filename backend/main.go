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
	agentController "github.com/ems/backend/internal/agent/controller"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/internal/repository"
	"github.com/ems/backend/pkg/config"

	"github.com/ems/backend/pkg/database"
	"github.com/ems/backend/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	// 1. Load .env file (Try current and parent directory)
	_ = godotenv.Load() // Try current directory (e.g. if running from root)
	_ = godotenv.Load("../.env") // Try parent directory (e.g. if running from backend/)

	// 2. Load configuration
	if err := config.Load("config/config.yaml"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 根据配置选择存储模式
	storageMode := "memory"
	if config.Cfg.Storage.Mode != "" {
		storageMode = config.Cfg.Storage.Mode
	}

	log.Printf("Starting EMS in %s mode", storageMode)

	if storageMode == "memory" {
		runMemoryMode()
	} else {
		runDatabaseMode()
	}
}

// runMemoryMode 内存模式运行
func runMemoryMode() {
	log.Println("Running in MEMORY mode - data will be lost on restart")

	// 初始化内存存储
	v1.InitMemory()
	log.Println("Memory store initialized with mock data")

	// Setup Gin
	gin.SetMode(config.Cfg.Server.Mode)
	router := gin.Default()

	// CORS middleware
	router.Use(CORSMiddleware())

	// Setup routes
	setupMemoryRoutes(router)

	// Start server
	startServer(router)
}

// runDatabaseMode 数据库模式运行
func runDatabaseMode() {
	log.Println("Running in DATABASE mode")

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
		&model.ManualDocument{},
		&model.ManualChunk{},
		&model.RepairCostDetail{},
		&model.EquipmentRuntimeSnapshot{},
		&model.AgentSession{},
		&model.AgentArtifact{},
		&model.AgentEvidenceLink{},
		&model.AgentUsage{},
		&model.AgentSkill{},
		&model.AgentKnowledge{},
		&model.AgentExperience{},
		&model.AgentConversation{},
		&model.AgentMessage{},
		&model.AgentPushSubscription{},
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

	// 补种演示数据 (Milestone: Data Parity)
	if err := repository.SeedDatabase(database.GetDB()); err != nil {
		log.Printf("Warning: Seeding failed: %v", err)
	}

	// Setup Gin
	gin.SetMode(config.Cfg.Server.Mode)
	router := gin.Default()

	// CORS middleware
	router.Use(CORSMiddleware())

	// Setup routes
	setupDatabaseRoutes(router)

	// Start server
	startServer(router)
}

// CORSMiddleware CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// setupMemoryRoutes 设置内存模式路由
func setupMemoryRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", v1.Login)
			auth.POST("/logout", v1.LogoutMemory)
			auth.POST("/refresh", v1.RefreshToken)
			auth.POST("/change-password", middleware.AuthMiddleware(), v1.ChangePassword)
			auth.POST("/apply", v1.ApplyForAccount)
			auth.GET("/me", middleware.AuthMiddleware(), v1.GetCurrentUser)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User management routes (admin only)
			users := protected.Group("/users")
			{
				users.GET("", v1.GetUsersMemory)
				users.POST("", v1.CreateUserMemory)
				users.GET("/applications", v1.GetPendingApplicationsMemory)
				users.PUT("/:id", v1.UpdateUserMemory)
				users.PUT("/:id/approve", v1.ApproveApplicationMemory)
			}

			// Organization routes
			org := protected.Group("/organization")
			{
				// Bases
				org.GET("/bases", v1.ListBasesMemory)
				org.POST("/bases", v1.CreateBaseMemory)
				org.PUT("/bases/:id", v1.UpdateBaseMemory)
				org.DELETE("/bases/:id", v1.DeleteBaseMemory)

				// Factories
				org.GET("/factories", v1.ListFactoriesMemory)
				org.POST("/factories", v1.CreateFactoryMemory)
				org.PUT("/factories/:id", v1.UpdateFactoryMemory)
				org.DELETE("/factories/:id", v1.DeleteFactoryMemory)

				// Workshops
				org.GET("/workshops", v1.ListWorkshopsMemory)
				org.POST("/workshops", v1.CreateWorkshopMemory)
				org.PUT("/workshops/:id", v1.UpdateWorkshopMemory)
				org.DELETE("/workshops/:id", v1.DeleteWorkshopMemory)
			}

			// Equipment routes
			equipment := protected.Group("/equipment")
			{
				equipment.GET("", v1.ListEquipmentMemory)
				equipment.GET("/:id", v1.GetEquipmentMemory)
				equipment.GET("/qr/:code", v1.GetEquipmentByQRCodeMemory)
				equipment.POST("", v1.CreateEquipmentMemory)
				equipment.PUT("/:id", v1.UpdateEquipmentMemory)
				equipment.DELETE("/:id", v1.DeleteEquipmentMemory)
				equipment.GET("/statistics", v1.GetEquipmentStatisticsMemory)

				// Equipment types
				equipment.GET("/types", v1.ListEquipmentTypesMemory)
				equipment.POST("/types", v1.CreateEquipmentTypeMemory)
				equipment.PUT("/types/:id", v1.UpdateEquipmentTypeMemory)
				equipment.DELETE("/types/:id", v1.DeleteEquipmentTypeMemory)
			}

			// Inspection routes
			inspection := protected.Group("/inspection")
			{
				// Templates
				inspection.GET("/templates", v1.ListInspectionTemplatesMemory)
				inspection.GET("/templates/:id", v1.GetInspectionTemplateMemory)
				inspection.POST("/templates", v1.CreateInspectionTemplateMemory)
				inspection.POST("/items", v1.CreateInspectionItemMemory)

				// Tasks
				inspection.GET("/tasks", v1.ListInspectionTasksMemory)
				inspection.GET("/tasks/:id", v1.GetInspectionTaskMemory)
				inspection.GET("/my-tasks", v1.GetMyTasksMemory)
				inspection.GET("/my-stats", v1.GetMyTaskStatisticsMemory)
				inspection.POST("/start", v1.StartInspectionMemory)
				inspection.POST("/complete", v1.CompleteInspectionMemory)
				inspection.GET("/statistics", v1.GetInspectionStatisticsMemory)
			}

			// Repair routes
			repair := protected.Group("/repair")
			{
				repair.GET("/orders", v1.ListRepairOrdersMemory)
				repair.GET("/orders/:id", v1.GetRepairOrderMemory)
				repair.POST("/orders", v1.CreateRepairOrderMemory)
				repair.POST("/orders/:id/assign", v1.AssignRepairOrderMemory)
				repair.POST("/orders/:id/start", v1.StartRepairMemory)
				repair.POST("/orders/:id/update", v1.UpdateRepairMemory)
				repair.POST("/orders/:id/confirm", v1.ConfirmRepairMemory)
				repair.POST("/orders/:id/audit", v1.AuditRepairMemory)
				repair.GET("/my-tasks", v1.GetMyRepairTasksMemory)
				repair.GET("/my-stats", v1.GetMyRepairStatisticsMemory)
				repair.GET("/statistics", v1.GetRepairStatisticsMemory)
			}

			// Maintenance routes
			maintenance := protected.Group("/maintenance")
			{
				// Plans
				maintenance.GET("/plans", v1.ListMaintenancePlansMemory)
				maintenance.POST("/plans", v1.CreateMaintenancePlanMemory)
				maintenance.POST("/items", v1.CreateMaintenanceItemMemory)

				// Tasks
				maintenance.POST("/tasks/generate", v1.GenerateMaintenanceTasksMemory)
				maintenance.GET("/tasks", v1.ListMaintenanceTasksMemory)
				maintenance.GET("/tasks/:id", v1.GetMaintenanceTaskMemory)
				maintenance.GET("/my-tasks", v1.GetMyMaintenanceTasksMemory)

				// Execution
				maintenance.POST("/start", v1.StartMaintenanceMemory)
				maintenance.POST("/complete", v1.CompleteMaintenanceMemory)

				// Statistics
				maintenance.GET("/statistics", v1.GetMaintenanceStatisticsMemory)
			}

			// Spare parts routes
			spareparts := protected.Group("/spareparts")
			{
				// Parts
				spareparts.GET("", v1.ListSparePartsMemory)
				spareparts.POST("", v1.CreateSparePartMemory)
				spareparts.PUT("/:id", v1.UpdateSparePartMemory)
				spareparts.DELETE("/:id", v1.DeleteSparePartMemory)

				// Inventory
				spareparts.GET("/inventory", v1.GetInventoryMemory)
				spareparts.POST("/stock-in", v1.StockInMemory)
				spareparts.POST("/stock-out", v1.StockOutMemory)
				spareparts.GET("/alerts", v1.GetLowStockAlertsMemory)

				// Consumption
				spareparts.GET("/consumptions", v1.GetConsumptionsMemory)
				spareparts.POST("/consumptions", v1.CreateConsumptionMemory)

				// Statistics
				spareparts.GET("/statistics", v1.GetSparePartStatisticsMemory)
			}

			// Analytics routes
			analytics := protected.Group("/analytics")
			{
				analytics.GET("/dashboard", v1.GetDashboardOverviewMemory)
				analytics.GET("/mttr-mtbf", v1.GetMTTRMTBFMemory)
				analytics.GET("/trends", v1.GetTrendDataMemory)
				analytics.GET("/failures", v1.GetFailureAnalysisMemory)
				analytics.GET("/top-failures", v1.GetTopFailureEquipmentMemory)
			}

			// Knowledge base routes
			knowledge := protected.Group("/knowledge")
			{
				knowledge.GET("", v1.ListKnowledgeArticlesMemory)
				knowledge.GET("/:id", v1.GetKnowledgeArticleMemory)
				knowledge.POST("", v1.CreateKnowledgeArticleMemory)
				knowledge.PUT("/:id", v1.UpdateKnowledgeArticleMemory)
				knowledge.DELETE("/:id", v1.DeleteKnowledgeArticleMemory)
				knowledge.GET("/search", v1.SearchKnowledgeArticlesMemory)
				knowledge.POST("/convert-repair", v1.ConvertFromRepairMemory)
			}

			// Agent routes
			agent := protected.Group("/agent")
			agentCtrl := agentController.NewAgentController()
			{
				agent.POST("/maintenance/recommend", agentCtrl.RecommendMaintenance)
				agent.POST("/audit/repair", agentCtrl.AuditRepair)
				agent.POST("/audit/maintenance", agentCtrl.AuditMaintenance)
				agent.POST("/analyze", agentCtrl.Analyze)
				agent.POST("/chat", agentCtrl.Chat)
				agent.GET("/knowledges", agentCtrl.ListKnowledges)
				agent.PUT("/knowledge/:id/status", agentCtrl.AuditKnowledge)
				agent.GET("/conversations", agentCtrl.ListConversations)
				agent.GET("/conversations/:id", agentCtrl.GetConversation)
				agent.GET("/skills", agentCtrl.ListSkills)
				agent.POST("/skills", agentCtrl.CreateSkill)
				agent.GET("/skills/:id", agentCtrl.GetSkill)
				agent.PUT("/skills/:id", agentCtrl.UpdateSkill)
				agent.GET("/equipment/:id/prediction", agentCtrl.GetEquipmentPrediction)
				agent.POST("/subscribe", agentCtrl.Subscribe)
				agent.GET("/sessions", agentCtrl.ListSessions)
				agent.GET("/sessions/:id", agentCtrl.GetSession)
				agent.GET("/artifacts/:id", agentCtrl.GetArtifact)
			}
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/health", v1.HealthCheckMemory)
}

// setupDatabaseRoutes 设置数据库模式路由
func setupDatabaseRoutes(router *gin.Engine) {
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

			// Agent routes
			agent := protected.Group("/agent")
			agentCtrl := agentController.NewAgentController()
			{
				agent.POST("/maintenance/recommend", agentCtrl.RecommendMaintenance)
				agent.POST("/audit/repair", agentCtrl.AuditRepair)
				agent.POST("/audit/maintenance", agentCtrl.AuditMaintenance)
				agent.POST("/analyze", agentCtrl.Analyze)
				agent.POST("/chat", agentCtrl.Chat)
				agent.GET("/knowledges", agentCtrl.ListKnowledges)
				agent.PUT("/knowledge/:id/status", agentCtrl.AuditKnowledge)
				agent.GET("/conversations", agentCtrl.ListConversations)
				agent.GET("/conversations/:id", agentCtrl.GetConversation)
				agent.GET("/skills", agentCtrl.ListSkills)
				agent.POST("/skills", agentCtrl.CreateSkill)
				agent.GET("/skills/:id", agentCtrl.GetSkill)
				agent.PUT("/skills/:id", agentCtrl.UpdateSkill)
				agent.GET("/equipment/:id/prediction", agentCtrl.GetEquipmentPrediction)
				agent.POST("/subscribe", agentCtrl.Subscribe)
				agent.GET("/sessions", agentCtrl.ListSessions)
				agent.GET("/sessions/:id", agentCtrl.GetSession)
				agent.GET("/artifacts/:id", agentCtrl.GetArtifact)
			}
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ems-api"})
	})
}

// startServer 启动 HTTP 服务器
func startServer(router *gin.Engine) {
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
