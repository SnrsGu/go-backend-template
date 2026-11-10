package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-backend-template/internal/platform/config"

	_ "github.com/go-sql-driver/mysql"
)

// Open 根据配置创建并验证数据库连接池。
func Open(ctx context.Context, databaseDSN string, cfg config.DatabaseConfig) (*sql.DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	if databaseDSN == "" {
		return nil, errors.New("database DSN is empty")
	}

	var driverName string

	switch cfg.Driver {
	case "mysql":
		driverName = "mysql"
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	databaseConnection, err := sql.Open(driverName, databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("open %s database: %w", cfg.Driver, err)
	}

	databaseConnection.SetMaxOpenConns(cfg.MaxOpenConnections)

	databaseConnection.SetMaxIdleConns(cfg.MaxIdleConnections)

	databaseConnection.SetConnMaxLifetime(cfg.ConnectionMaxLifetimeDuration())

	databaseConnection.SetConnMaxIdleTime(cfg.ConnectionMaxIdleTimeDuration())

	connectContext, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout())
	defer cancel()

	if err := databaseConnection.PingContext(connectContext); err != nil {
		_ = databaseConnection.Close()
		return nil, fmt.Errorf("ping %s database: %w", cfg.Driver, err)
	}
	return databaseConnection, nil
}
