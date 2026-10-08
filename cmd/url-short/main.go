package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-short/internal/config"
	"url-short/internal/handler"
	"url-short/internal/service"
	"url-short/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// 1. config: cfg := config.MustLoad()
	// 2. logger: slog.New(slog.NewJSONHandler(os.Stdout, ...)) + slog.SetDefault(logger)
	// 3. storage: pgxpool.New(ctx, cfg.Postgres.DSN()) + Ping
	// 4. service: service.New(repo, logger)
	// 5. handler: handler.New(svc) -> h.Router(cfg.Auth, logger)
	// 6. http.Server + graceful shutdown

	_ = godotenv.Load(".env")

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

	err = pool.Ping(ctxTime)
	if err != nil {
		logger.Error("failed to ping to db", slog.Any("error", err))
		os.Exit(1)
	}

	defer pool.Close()

	repo := storage.NewPostgresPool(pool)
	svc := service.New(repo, logger)
	h := handler.New(svc)

	r := h.Router(cfg.Auth, logger)

	srv := &http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: r, // gin.Engine реализует http.Handler
	}

	go func() {
		logger.Info("server starting", slog.String("addr", cfg.HTTPServer.Addr))
		err := srv.ListenAndServe()
		if err != nil && (!errors.Is(err, http.ErrServerClosed)) {
			logger.Error("server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", slog.Any("error", err))
	}

	logger.Info("server stopped")

}
