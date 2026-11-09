package main

import (
	"context"
	"errors"
	"flag"
	"go-backend-template/internal/platform/config"
	"go-backend-template/internal/platform/middleware"
	"go-backend-template/internal/platform/response"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	configuration, err := config.Load(*configPath)
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}
	appRouter := router.New(logger)
	appRouter.GET("/health", func(writer http.ResponseWriter, request *http.Request) error {
		return response.Success(writer, map[string]string{"status": "ok"})
	})
	httpHandler := middleware.Recovery(logger)(
		appRouter.Handler(),
	)
	httpServer := &http.Server{
		Addr:              configuration.Server.Address,
		Handler:           httpHandler,
		ReadHeaderTimeout: configuration.Server.ReadHeaderTimeout(),
	}

	//router := http.NewServeMux()
	//router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
	//	w.Header().Set("Content-Type", "application/json")
	//	if err := json.NewEncoder(w).Encode(map[string]string{
	//		"status": "ok",
	//	}); err != nil {
	//		logger.Error("failed to encode health response", "error", err)
	//	}
	//})

	//httpServer := &http.Server{
	//	Addr:              configuration.Server.Address,
	//	Handler:           router,
	//	ReadHeaderTimeout: configuration.Server.ReadHeaderTimeout(),
	//}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("starting http server", "address", configuration.Server.Address)
		serverErrors <- httpServer.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		logger.Info("shutting down http server")
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start http server", "error", err)
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
		logger.Error("failed to shutdown http server", "error", err)
		os.Exit(1)
	}
	logger.Info("HTTP server stopped")
}
