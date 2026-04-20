package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lusiker/clipper/internal/config"
	"github.com/lusiker/clipper/internal/handler"
	"github.com/lusiker/clipper/internal/middleware"
	"github.com/lusiker/clipper/internal/pkg/storage"
	"github.com/lusiker/clipper/internal/repository"
	"github.com/lusiker/clipper/internal/service"
	"github.com/lusiker/clipper/internal/ws"
)

func main() {
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Ensure uploads directory exists
	if err := storage.EnsureUploadDir(); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	db, err := repository.InitDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	clipRepo := repository.NewClipRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := service.NewAuthService(userRepo)
	clipService := service.NewClipService(clipRepo, deviceRepo)
	deviceService := service.NewDeviceService(deviceRepo)

	hub := ws.NewHub(clipService, deviceService)
	go hub.Run()

	middleware.InitSessionStore(cfg.Server.SessionSecret)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(corsMiddleware())

	api := r.Group("/api/v1")

	authHandler := handler.NewAuthHandler(authService)
	authHandler.RegisterRoutes(api.Group("/auth"))

	clipHandler := handler.NewClipHandler(clipService, hub)
	clipHandler.RegisterRoutes(api.Group("/clips"))

	deviceHandler := handler.NewDeviceHandler(deviceService, hub)
	deviceHandler.RegisterRoutes(api.Group("/devices"))

	// Static file routes for uploaded images
	r.Static("/uploads", "./uploads")

	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/favicon.svg", "./web/dist/favicon.svg")
	r.StaticFile("/", "./web/dist/index.html")
	r.StaticFile("/login", "./web/dist/index.html")
	r.StaticFile("/register", "./web/dist/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}