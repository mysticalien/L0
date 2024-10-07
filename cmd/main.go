package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/handlers"
	"L0/internal/kafka"
	"L0/internal/logger"
	"L0/internal/storage"
)

const configPath = "config/config.yaml"

func main() {
	cfg, err := config.MustLoad(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	serverLogger := logger.NewLogger(cfg.Logger.Level)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.Timeout)
	defer cancel()

	dbStorage, err := storage.NewStorage(ctx, cfg.Database.DatabaseURL, serverLogger)
	if err != nil {
		serverLogger.Error("Failed to initialize storage", slog.Any("error", err))
		return
	}
	defer dbStorage.Close()

	orderCache := cache.NewCache()

	err = orderCache.RestoreFromDB(dbStorage, ctx)
	if err != nil {
		serverLogger.Error("Failed to restore cache from DB", slog.Any("error", err))
		return
	}

	go kafka.ConsumeOrders(&cfg.Kafka, orderCache, dbStorage)

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: handlers.SetupRoutes(orderCache, dbStorage),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverLogger.Error("HTTP server ListenAndServe", slog.Any("error", err))
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	<-sigint
	serverLogger.Info("Shutting down server")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		serverLogger.Error("Server Shutdown failed", slog.Any("error", err))
	}

	cancel()
}
