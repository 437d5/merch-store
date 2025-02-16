package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/437d5/merch-store/internal/config"
	"github.com/437d5/merch-store/internal/handler"
	"github.com/437d5/merch-store/internal/repository"
	"github.com/437d5/merch-store/internal/service"
	"github.com/437d5/merch-store/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.MustLoad()

	logger := logger.NewLogger(cfg.Log.LogMode, slog.LevelError)

	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Db.DbUser, cfg.Db.DbPass, cfg.Db.DbHost, cfg.Db.DbPort, cfg.Db.DbName,
	))
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	config.MaxConns = 100
	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	userRepo := repository.NewUserRepo(dbpool, logger)
	itemRepo := repository.NewItemRepo(dbpool, logger)
	transactionRepo := repository.NewTransRepo(dbpool, logger)

	userService := service.NewUserService(userRepo, logger)
	marketService := service.NewMarketService(userRepo, logger, itemRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo, logger)

	h := handler.NewHandler(
		userService, marketService, transactionService, logger,
	)

	router := gin.Default()
	h.SetupRoutes(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Srv.SrvPort),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()
	logger.Info("server started", "port", cfg.Srv.SrvPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
	}

	logger.Info("server exited")
}
