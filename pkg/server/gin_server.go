package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"marketplace/config"
	"marketplace/internal/handlers/http/v1"
	"marketplace/internal/repositories/postgres"
	"marketplace/internal/services"
	"marketplace/pkg/database"

	"github.com/gin-gonic/gin"
)

type GinServer struct {
	router     *gin.Engine
	db         *database.PostgresDatabase
	cfg        *config.Config
	httpServer *http.Server
}

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
	authService := services.NewAuthServiceImpl(userRepository, cfg)
	authHandlers := v1.NewAuthHTTPHandlers(authService)

	router := gin.Default()
	v1Routes := router.Group(
		"/api/v1",
		v1.RequestIDMiddleware(),
		v1.SetLoggerMiddleware(),
	)

	authRoutes := v1Routes.Group("/auth")
	authRoutes.POST("/register", authHandlers.Register)
	authRoutes.POST("/login", authHandlers.Login)
	authRoutes.GET("/me", authHandlers.GetMe)

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
