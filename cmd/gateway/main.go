package main

import (
	"context"
	"flag"
	"fmt"
	"go-backend-template/internal/platform/config"
	"go-backend-template/internal/platform/logger"
	"go-backend-template/internal/platform/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Gateway服务黜退失败: %v", err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	appLogger, err := logger.New(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("OK"))
	})
	appLogger.Info("Gateway服务初始化完成")
	if err = server.Run(ctx, cfg.Server, mux, appLogger); err != nil {
		return fmt.Errorf("gateway HTTP 服务运行失败: %w", err)
	}
	appLogger.Info("Gateway 服务已退出")
	return nil
}
