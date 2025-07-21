package slogger

import (
	"context"
	"log/slog"
	"os"

	"github.com/sytallax/prettylog"

	"marketplace/config"
)

func SetLogger(env config.AppEnv) {
	var defaultLogger *slog.Logger
	switch env {
	case config.Local:
		prettyHandler := prettylog.NewHandler(&slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		defaultLogger = slog.New(prettyHandler)
	case config.Dev:
		defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case config.Prod:
		defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	slog.SetDefault(defaultLogger)
}

func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value("Logger").(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
