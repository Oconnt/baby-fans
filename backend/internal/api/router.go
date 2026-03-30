package api

import (
	"baby-fans/internal/api/handler"
	"baby-fans/internal/api/middleware"
	"baby-fans/internal/model"
	"baby-fans/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 优化的 CORS 中间件
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// JSON 日志中间件（输出到 stdout，供 Docker log 收集）
	r.Use(middleware.JSONLogger())
	// Panic recovery 中间件（带日志输出）
	r.Use(middleware.RecoveryWithLogger())

	// 1. 先注册具体的静态资源路径
	r.Static("/storage/uploads", "./storage/uploads")
	r.Static("/uploads", "./storage/uploads") // 增加兼容性

	// 2. 使用 StaticFile 精确托管主页，避免根路径通配符冲突
	r.StaticFile("/", "./web/index.html")

	authService := &service.AuthService{}
	pointsService := &service.PointsService{}
	shopService := &service.ShopService{}
	taskHandler := handler.NewTaskHandler()

	authHandler := &handler.AuthHandler{Service: authService}
	pointsHandler := &handler.PointsHandler{Service: pointsService}
	shopHandler := &handler.ShopHandler{Service: shopService}
	versionHandler := &handler.VersionHandler{}
	adminHandler := &handler.AdminHandler{}

	// Health check
	r.Any("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	public := r.Group("")
	public.Use(middleware.IdempotencyMiddleware())
	{
		public.GET("/version", versionHandler.GetVersion)
		public.POST("/login/face", authHandler.LoginFace)
		public.GET("/login/code", authHandler.LoginCode)
		public.POST("/register", authHandler.Register)
		public.POST("/api/v1/auth/wechat/login", authHandler.WeChatLogin)
	}

	// Admin routes (should be protected in production, e.g., via IP whitelist)
	admin := r.Group("/admin")
	{
		admin.POST("/reload", adminHandler.ReloadConfig)
	}

	// Global / Shop access for both
	r.GET("/parent/items", shopHandler.GetItems)

	// Parent routes
	parent := r.Group("/parent")
	parent.Use(middleware.AuthMiddleware(model.RoleParent))
	parent.Use(middleware.IdempotencyMiddleware())
	{
		parent.GET("/children", authHandler.GetChildren)
		parent.POST("/children/bind", authHandler.BindChildByCode)
		parent.DELETE("/children/:id", authHandler.UnbindChild)
		parent.POST("/binding/code", authHandler.GenerateBindingCode)
		parent.GET("/templates", pointsHandler.GetTemplates)
		parent.POST("/templates", pointsHandler.SaveTemplate)
		parent.DELETE("/templates/:id", pointsHandler.DeleteTemplate)
		parent.POST("/items", shopHandler.SaveItem)
		parent.PUT("/items/:id/stock", shopHandler.UpdateStock)
		parent.DELETE("/items/:id", shopHandler.DeleteItem)
		parent.GET("/redemptions", shopHandler.GetRedemptions)
		parent.POST("/points/manage", pointsHandler.ManagePoints)
		parent.GET("/points/records", pointsHandler.GetPointsRecords)
		parent.POST("/redemption/confirm/:id", shopHandler.Confirm)
		parent.POST("/redemption/cancel/:id", shopHandler.Cancel)
		// Task routes
		parent.GET("/task-templates", taskHandler.GetTaskTemplates)
		parent.POST("/task-templates", taskHandler.CreateTaskTemplate)
		parent.DELETE("/task-templates/:id", taskHandler.DeleteTaskTemplate)
		parent.GET("/tasks", taskHandler.GetParentTasks)
		parent.POST("/tasks", taskHandler.CreateTask)
		parent.PUT("/tasks/:id/status", taskHandler.UpdateTaskStatus)
	}

	// Child routes
	child := r.Group("/child")
	child.Use(middleware.AuthMiddleware(model.RoleChild))
	child.Use(middleware.IdempotencyMiddleware())
	{
		child.GET("/overview", authHandler.GetOverview)
		child.POST("/binding/accept", authHandler.AcceptBinding)
		child.POST("/exchange", shopHandler.Exchange)
		child.POST("/profile", authHandler.UpdateProfile)
		child.GET("/points/history", pointsHandler.GetPointsHistory)
		child.GET("/redemptions", shopHandler.GetRedemptions)
		// Task routes
		child.GET("/tasks/today", taskHandler.GetTodayTasks)
		child.GET("/tasks", taskHandler.GetChildTasks)
		child.GET("/tasks/:id", taskHandler.GetTaskDetail)
		child.PUT("/tasks/:id/complete", taskHandler.CompleteTask)
	}

	// Common routes (requires auth)
	parent.POST("/profile", authHandler.UpdateProfile)
	parent.POST("/avatar", authHandler.UploadAvatar)
	child.POST("/avatar", authHandler.UploadAvatar)

	return r
}
