package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"marketplace/config"
	"marketplace/internal/logger"
	"marketplace/migrations"
	"marketplace/pkg/database"
	"marketplace/pkg/server"
)

func main() {
	cfg := config.MustLoad()
	slogger.SetLogger(cfg.AppEnv)

	db := database.New(cfg.GetDBURL())
	migrations.Migrate()

	srv := server.NewGinServer(cfg, db)
	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Gin server error", slog.Any("error", err))
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Error during server shutdown", slog.Any("error", err))
	}
	db.Pool.Close()
	slog.Info("App gracefully stopped")
}
