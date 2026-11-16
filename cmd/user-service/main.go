package main

import (
	"context"
	"flag"
	"fmt"
	userservice "go-backend-template/internal/service/user"
	usermodule "go-backend-template/internal/service/user/user"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-backend-template/internal/platform/config"
	"go-backend-template/internal/platform/database"
	"go-backend-template/internal/platform/logger"
	"go-backend-template/internal/platform/middleware"
	"go-backend-template/internal/platform/router"
	"go-backend-template/internal/platform/server"
)

func main() {
	if err := run(); err != nil {
		log.Printf("user-service 退出失败: %v", err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String(
		"config",
		"configs/user-service.yaml",
		"user-service 配置文件路径",
	)
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		return fmt.Errorf("加载 user-service 配置失败: %w", err)
	}

	appLogger, err := logger.New(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		return fmt.Errorf("初始化日志系统失败: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	databases, err := database.Open(ctx, cfg.Database, appLogger)
	if err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}

	defer func() {
		if err = databases.Close(); err != nil {
			appLogger.Error("关闭数据库连接失败", "error", err)
		}
	}()

	userDB, err := databases.Get("user")
	if err != nil {
		return fmt.Errorf("获取 user-service 数据库失败: %w", err)
	}

	userRepository := usermodule.NewRepository(userDB)
	userService := usermodule.NewService(userRepository)
	userHandler := usermodule.NewHandler(userService)

	appRouter := router.New(appLogger)
	userservice.RegisterRoutes(appRouter, userHandler)

	httpHandler := middleware.Chain(
		appRouter,
		middleware.RequestID(),
		middleware.AccessLog(appLogger),
		middleware.Recovery(appLogger),
	)

	appLogger.Info("user-service 初始化完成")

	if err = server.Run(
		ctx,
		cfg.Server,
		httpHandler,
		appLogger,
	); err != nil {
		return fmt.Errorf("user-service 运行失败: %w", err)
	}

	appLogger.Info("user-service 已停止")
	return nil
}
