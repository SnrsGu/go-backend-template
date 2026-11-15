package server

import (
	"context"
	"errors"
	"fmt"
	"go-backend-template/internal/platform/config"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func Run(ctx context.Context, config config.ServerConfig, handler http.Handler, logger *slog.Logger) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("HTTP 服务启动已取消: %w", err)
	}

	httpServer := &http.Server{
		Addr:              config.Address,
		Handler:           handler,
		ReadHeaderTimeout: time.Duration(config.ReadHeaderTimeout) * time.Second,
		IdleTimeout:       time.Duration(config.IdleTimeout) * time.Second,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return fmt.Errorf("监听HTTP地址 %s 失败: %w", config.Address, err)
	}
	defer listener.Close()

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- httpServer.Serve(listener)
	}()

	logger.Info("HTTP服务已启动", "address", listener.Addr().String())

	select {
	case err = <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP服务运行失败: %w", err)
		}
		return nil
	case <-ctx.Done():
		logger.Info("HTTP服务开始关闭")
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(config.ShutdownTimeout)*time.Second,
		)
		defer cancel()
		if err = httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("HTTP服务退出失败，开始强制关闭", "error", err)
			if closeErr := httpServer.Close(); closeErr != nil {
				logger.Error("HTTP服务强制关闭失败", "error", closeErr)
			}
			<-serverErrors
			return fmt.Errorf("HTTP服务优雅退出失败: %w", err)
		}
		if err = <-serverErrors; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP服务关闭时发生错误: %w", err)
		}
		logger.Info("HTTP服务已关闭")
		return nil
	}
}
