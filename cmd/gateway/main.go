package main

import (
	"context"
	"errors"
	"flag"
	gatewayroutes "go-backend-template/internal/gateway/routes"
	"go-backend-template/internal/platform/config"
	"go-backend-template/internal/platform/logger"
	"go-backend-template/internal/platform/middleware"
	"go-backend-template/internal/platform/router"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configPath := flag.String("config", "configs/gateway.json", "path to config file")
	flag.Parse()
	//logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	appLogger := logger.New(os.Stdout, slog.LevelInfo)

	configuration, err := config.Load(*configPath)
	if err != nil {
		appLogger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}
	appRouter := router.New(appLogger)

	appRouter.Use(
		middleware.Recovery(appLogger),
		middleware.CORS(middleware.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type"},
		}),
	)
	gatewayroutes.Register(appRouter)

	httpServer := &http.Server{
		Addr:              configuration.Server.Address,
		Handler:           appRouter.Handler(),
		ReadHeaderTimeout: configuration.Server.ReadHeaderTimeout(),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		appLogger.Info("starting http server", "address", configuration.Server.Address)
		serverErrors <- httpServer.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		appLogger.Info("shutting down http server")
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			appLogger.Error("failed to start http server", "error", err)
			os.Exit(1)
		}
		return
	}

	shutdownContext, cancel := context.WithTimeout(
		context.Background(),
		configuration.Server.ShutdownTimeout(),
	)
	defer cancel()

	if err := httpServer.Shutdown(shutdownContext); err != nil {
		appLogger.Error("failed to shutdown http server", "error", err)
		os.Exit(1)
	}
	appLogger.Info("HTTP server stopped")
}
