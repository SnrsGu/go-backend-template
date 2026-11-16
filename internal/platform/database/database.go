package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-backend-template/internal/platform/config"
	"log/slog"
	"sort"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type databaseConnection struct {
	db    *gorm.DB
	pools []*sql.DB
}

type Manager struct {
	connections map[string]*databaseConnection
}

func Open(
	ctx context.Context,
	cfg config.DatabaseConfig,
	logger *slog.Logger,
) (*Manager, error) {
	if logger == nil {
		return nil, fmt.Errorf("数据库日志实例为空")
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("数据库配置校验失败: %w", err)
	}
	manager := &Manager{
		connections: make(map[string]*databaseConnection),
	}
	names := make([]string, 0, len(cfg.Connections))
	for name := range cfg.Connections {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		connectionCfg := cfg.Connections[name]
		connection, err := openConnection(ctx, connectionCfg, logger)
		if err != nil {
			_ = manager.Close()
			return nil, fmt.Errorf("初始化数据库 %q 失败: %w", name, err)
		}
		manager.connections[name] = connection
		logger.Info("数据库初始化成功", "name", name, "driver", connectionCfg.Driver, "replicas", len(connectionCfg.Replicas))
	}
	return manager, nil
}

func openConnection(ctx context.Context, cfg config.DatabaseConnectionConfig, logger *slog.Logger) (_ *databaseConnection, err error) {
	connection := &databaseConnection{}
	defer func() {
		if err != nil {
			_ = connection.close()
		}
	}()
	primaryPool, primaryDialect, err := openEndpoint(ctx, cfg.Driver, cfg.Primary, cfg.Pool, cfg.ConnectTimeout)
	if err != nil {
		return nil, fmt.Errorf("初始化主数据库失败: %w", err)
	}
	connection.pools = append(connection.pools, primaryPool)
	gormDB, err := gorm.Open(
		primaryDialect,
		&gorm.Config{
			DisableAutomaticPing: true,
			Logger:               newGormLogger(logger, time.Duration(cfg.SlowThreshold)*time.Millisecond),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("初始化 GORM 失败: %w", err)
	}
	replicaDialects := make([]gorm.Dialector, 0, len(cfg.Replicas))
	for idx, replica := range cfg.Replicas {
		replicaPool, replicaDialect, err := openEndpoint(
			ctx, cfg.Driver, replica, cfg.Pool, cfg.ConnectTimeout,
		)
		if err != nil {
			return nil, fmt.Errorf("初始化从库 %d 失败: %w", idx+1, err)
		}
		connection.pools = append(connection.pools, replicaPool)
		replicaDialects = append(replicaDialects, replicaDialect)
	}
	if len(replicaDialects) > 0 {
		resolver := dbresolver.Register(dbresolver.Config{
			Replicas: replicaDialects, Policy: dbresolver.RandomPolicy{},
		})
		if err := gormDB.Use(resolver); err != nil {
			return nil, fmt.Errorf("注册数据库读写分离失败: %w", err)
		}
	}
	connection.db = gormDB
	return connection, nil
}

func (m *Manager) Get(name string) (*gorm.DB, error) {
	connection, exists := m.connections[name]
	if !exists {
		return nil, fmt.Errorf("数据库 %q 未配置", name)
	}
	return connection.db, nil
}

func (m *Manager) Ping(ctx context.Context, name string) error {
	connection, exists := m.connections[name]
	if !exists {
		return fmt.Errorf("数据库 %q 未配置", name)
	}
	for idx, pool := range connection.pools {
		if err := pool.PingContext(ctx); err != nil {
			return fmt.Errorf("数据库 %q 节点 %d 健康检查失败: %w", name, idx, err)
		}
	}
	return nil
}

func (m *Manager) Stats(name string) ([]sql.DBStats, error) {
	connection, exists := m.connections[name]
	if !exists {
		return nil, fmt.Errorf("数据库 %q 未配置", name)
	}

	stats := make([]sql.DBStats, 0, len(connection.pools))
	for _, pool := range connection.pools {
		stats = append(stats, pool.Stats())
	}

	return stats, nil
}

func (m *Manager) Close() error {
	var errs []error
	for name, connection := range m.connections {
		if err := connection.close(); err != nil {
			errs = append(errs, fmt.Errorf("关闭数据库 %q 失败: %w", name, err))
		}
	}
	return errors.Join(errs...)
}

func (c *databaseConnection) close() error {
	var errs []error
	for _, pool := range c.pools {
		if err := pool.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
