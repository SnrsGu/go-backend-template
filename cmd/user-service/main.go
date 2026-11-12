package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-backend-template/internal/platform/apperror"
	"go-backend-template/internal/platform/config"
	"go-backend-template/internal/platform/database"
	"go-backend-template/internal/platform/logger"
	"go-backend-template/internal/platform/middleware"
	"go-backend-template/internal/platform/response"
	"go-backend-template/internal/platform/router"
)

func main() {
	configPath := flag.String("config", "configs/user-service.json", "path to the configuration file")
	flag.Parse()

	appLogger := logger.New(os.Stdout, slog.LevelInfo)

	if err := run(appLogger, *configPath); err != nil {
		appLogger.Error("user service stopped with error", "error", err)
		os.Exit(1)
	}
}

func run(appLogger *slog.Logger, configPath string) error {
	configuration, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	databaseURL := os.Getenv("DATABASE_URL") //root:devpassword@tcp(127.0.0.1:3306)/backend_template?parseTime=true&loc=UTC

	startupContext, cancelStartup := context.WithTimeout(context.Background(), configuration.Database.ConnectTimeout())
	defer cancelStartup()

	dbConn, err := database.Open(startupContext, databaseURL, configuration.Database)
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	sqlDB, err := dbConn.DB()
	defer func() {
		if err := sqlDB.Close(); err != nil {
			appLogger.Error("failed to close database", "error", err)
		}
	}()

	appLogger.Info("database connection established", "driver", configuration.Database.Driver)

	appRouter := router.New(appLogger)

	appRouter.Use(middleware.Recovery(appLogger))

	appRouter.GET("/health", func(w http.ResponseWriter, req *http.Request) error {
		return response.Success(w, map[string]string{"status": "ok"})
	})

	appRouter.GET("/ready", func(w http.ResponseWriter, req *http.Request) error {
		return checkDatabaseReady(w, req, sqlDB)
	})

	httpServer := &http.Server{
		Addr:              configuration.Server.Address,
		Handler:           appRouter.Handler(),
		ReadHeaderTimeout: configuration.Server.ReadHeaderTimeout(),
	}

	listener, err := net.Listen("tcp", configuration.Server.Address)
	if err != nil {
		return fmt.Errorf("listen HTTP address: %w", err)
	}
	defer listener.Close()

	signalContext, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- httpServer.Serve(listener)
	}()

	appLogger.Info(
		"user service started",
		"address", listener.Addr().String(),
	)

	select {
	case <-signalContext.Done():
		appLogger.Info("shutdown signal received")
		stop()

	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("HTTP server failed: %w", err)
	}

	shutdownContext, cancelShutdown := context.WithTimeout(
		context.Background(),
		configuration.Server.ShutdownTimeout(),
	)
	defer cancelShutdown()

	if err := httpServer.Shutdown(shutdownContext); err != nil {
		_ = httpServer.Close()
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	appLogger.Info("user service stopped")

	return nil
}

func checkDatabaseReady(w http.ResponseWriter, req *http.Request, sqlDB *sql.DB) error {
	checkContext, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(checkContext); err != nil {
		return apperror.Wrap(
			apperror.CodeInternal,
			"",
			fmt.Errorf("database readiness check: %w", err),
		)
	}
	return response.Success(w, map[string]string{
		"status": "ready",
	})
}
