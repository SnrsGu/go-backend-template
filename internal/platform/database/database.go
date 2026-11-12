package database

import (
	"context"
	"fmt"
	"go-backend-template/internal/platform/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open(ctx context.Context, databaseDSN string, cfg config.DatabaseConfig) (*gorm.DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate database config: %w", err)
	}

	if databaseDSN == "" {
		return nil, fmt.Errorf("database DSN is empty")
	}
	var dialector gorm.Dialector
	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(databaseDSN)
	default:
		return nil, fmt.Errorf("unknown driver: %s", cfg.Driver)
	}
	connectContext, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout())
	defer cancel()
	dbConnection, err := gorm.Open(dialector, &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("open %s database: %w", cfg.Driver, err)
	}

	sqlDatabase, err := dbConnection.DB()
	if err != nil {
		return nil, fmt.Errorf("open %s database: %w", cfg.Driver, err)
	}
	sqlDatabase.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDatabase.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDatabase.SetConnMaxLifetime(cfg.ConnectionMaxLifetimeDuration())
	sqlDatabase.SetConnMaxIdleTime(cfg.ConnectionMaxIdleTimeDuration())

	if err := sqlDatabase.PingContext(connectContext); err != nil {
		_ = sqlDatabase.Close()
		return nil, fmt.Errorf("open %s database: %w", cfg.Driver, err)
	}
	return dbConnection, nil
}
