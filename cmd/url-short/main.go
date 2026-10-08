package main

import (
	"context"
	"log/slog"
	"os"
	"time"
	"url-short/internal/config"
	"url-short/internal/handler"
	"url-short/internal/service"
	"url-short/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. config: cfg := config.MustLoad()
	// 2. logger: slog.New(slog.NewJSONHandler(os.Stdout, ...)) + slog.SetDefault(logger)
	// 3. storage: pgxpool.New(ctx, cfg.Postgres.DSN()) + Ping
	// 4. service: service.New(repo, logger)
	// 5. handler: handler.New(svc) -> h.Router(cfg.Auth, logger)
	// 6. http.Server + graceful shutdown

	cfg := config.MustLoad()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx := context.Background()

	ctxTime, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctxTime, cfg.Postgres.DSN())

	if err != nil {
		logger.Error("failed to connect to db", slog.Any("error", err))
		os.Exit(1)
	}

	defer pool.Close()

	repo := storage.NewPostgresPool(pool)
	svc := service.New(repo, logger)
	h := handler.New(svc)

	r := h.Router(cfg.Auth, logger)

	logger.Info("server starting", slog.String("addr", cfg.HTTPServer.Addr))
	if err := r.Run(cfg.HTTPServer.Addr); err != nil {
		logger.Error("server failed", slog.Any("error", err))
		os.Exit(1)
	}

}
