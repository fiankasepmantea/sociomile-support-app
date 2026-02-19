package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"backend/internal/cache"
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository/pg"
	"backend/internal/service"
	"backend/pkg/database"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresDB(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("❌ Database initialization failed: %v", err)
	}
	defer db.Close()

	if err := db.RunSeeder(); err != nil {
		log.Printf("⚠️ Seeder failed: %v", err)
	}

	repo := pg.NewPgRepository(db.DB)

	// Redis optional
	redisClient := cache.NewRedis(cfg)

	// SERVICES
	authService := service.NewAuthService(repo, cfg.JWTSecret, cfg.JWTExpiryHours, redisClient)
	convService := service.NewConversationService(repo)
	messageService := service.NewMessageService(repo)
	customerService := service.NewCustomerService(repo)
	ticketService := service.NewTicketService(repo)

	// HANDLERS
	authHandler := handler.NewAuthHandler(authService)
	webhookHandler := handler.NewWebhookHandler(convService)
	conversationHandler := handler.NewConversationHandler(convService)
	messageHandler := handler.NewMessageHandler(messageService)
	customerHandler := handler.NewCustomerHandler(customerService)
	ticketHandler := handler.NewTicketHandler(ticketService)

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), CORSMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().UTC(),
		})
	})

	api := router.Group("/api")

	// ---------- AUTH ----------
	authRoutes := api.Group("/auth")
	authRoutes.Use(middleware.TenantMiddleware(db.DB))
	authRoutes.POST("/login", authHandler.Login)

	// ---------- WEBHOOK ----------
	webhookRoutes := api.Group("/webhooks")
	webhookRoutes.Use(middleware.TenantMiddleware(db.DB))
	webhookRoutes.Use(middleware.WebhookRateLimit(redisClient, cfg.WebhookRateLimit))
	webhookRoutes.POST("/channel", webhookHandler.HandleChannel)

	// ---------- PROTECTED ----------
	protected := api.Group("")
	protected.Use(middleware.TenantMiddleware(db.DB))
	protected.Use(middleware.AuthMiddleware(db.DB, cfg.JWTSecret, redisClient))

	protected.POST("/logout", authHandler.Logout)
	protected.GET("/me", authHandler.Me)

	// Conversations
	conv := protected.Group("/conversations")
	conv.GET("", conversationHandler.ListConversations)
	conv.GET("/:id", conversationHandler.GetDetail)
	conv.PATCH("/:id/assign", conversationHandler.AssignAgent)
	conv.PATCH("/:id/status", conversationHandler.UpdateStatus)
	conv.POST("/:id/messages", messageHandler.SendAgentMessage)
	conv.GET("/:id/messages", messageHandler.ListMessages)

	// Customers
	cust := protected.Group("/customers")
	cust.GET("", customerHandler.List)
	cust.GET("/:id", customerHandler.Get)
	cust.POST("", customerHandler.Create)

	// Tickets
	ticket := protected.Group("/tickets")
	ticket.GET("", ticketHandler.List)
	ticket.GET("/:id", ticketHandler.Get)
	ticket.POST("", ticketHandler.Create)
	ticket.PATCH("/:id/status", ticketHandler.UpdateStatus)

	addr := ":" + cfg.ServerPort
	log.Println("🚀 Server running on", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Authorization, Content-Type, X-Tenant-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods",
			"GET, POST, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
