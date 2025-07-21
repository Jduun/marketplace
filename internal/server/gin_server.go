package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"marketplace/config"
	"marketplace/internal/database"
	"marketplace/internal/handlers/http/v1"
	"marketplace/internal/repositories/postgres"
	"marketplace/internal/services"
)

type GinServer struct {
	router     *gin.Engine
	db         *database.PostgresDatabase
	cfg        *config.Config
	httpServer *http.Server
}

// @title           Marketplace
// @version         1.0
// @description     Simple marketplace implementation.
// @schemes         http
// @host            localhost:8089
// @BasePath        /api/v1
func NewGinServer(cfg *config.Config, db *database.PostgresDatabase) *GinServer {
	switch cfg.AppEnv {
	case config.Local, config.Dev:
		gin.SetMode(gin.DebugMode)
	case config.Prod:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	userRepository := postgres.NewUserPostgresRepository(db)
	authService := services.NewAuthServiceImpl(userRepository, time.Duration(cfg.TokenTTLMinutes), cfg.JWTSecret)
	authHandlers := v1.NewAuthHTTPHandlers(authService)

	advertisementRepository := postgres.NewAdvertisementPostgresRepository(db)
	advertisementService := services.NewAdvertisementServiceImpl(advertisementRepository)
	advertisementHandlers := v1.NewAdvertisementHTTPHandlers(advertisementService)

	router := gin.Default()
	v1Routes := router.Group(
		"/api/v1",
		v1.RequestIDMiddleware(),
		v1.SetLoggerMiddleware(),
	)

	authRoutes := v1Routes.Group("/auth")
	authRoutes.POST("/register", authHandlers.Register)
	authRoutes.POST("/login", authHandlers.Login)
	authRoutes.GET("/me", v1.AuthMiddleware(cfg.JWTSecret), authHandlers.GetMe)

	advertisementRoutes := v1Routes.Group("/advertisements")
	advertisementRoutes.POST("/", v1.AuthMiddleware(cfg.JWTSecret), advertisementHandlers.CreateAdvertisement)
	advertisementRoutes.GET("/", v1.SetUserInfoMiddleware(cfg.JWTSecret), advertisementHandlers.GetAdvertisements)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.AppPort),
		Handler: router,
	}

	return &GinServer{
		router:     router,
		db:         db,
		cfg:        cfg,
		httpServer: httpServer,
	}
}

func (s *GinServer) Run() error {
	slog.Info("Starting Gin server")
	return s.httpServer.ListenAndServe()
}

func (s *GinServer) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down Gin server...")
	return s.httpServer.Shutdown(ctx)
}
