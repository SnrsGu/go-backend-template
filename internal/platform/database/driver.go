package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-backend-template/internal/platform/config"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openEndpoint(
	ctx context.Context,
	driver string,
	endpoint config.DatabaseEndpointConfig,
	pool config.DatabasePoolConfig,
	connectTimeout int,
) (*sql.DB, gorm.Dialector, error) {
	dsn, err := resolveDSN(endpoint)
	if err != nil {
		return nil, nil, err
	}
	driverName := ""
	var dialectFactory func(*sql.DB) gorm.Dialector

	switch driver {
	case "mysql":
		driverName = "mysql"
		dialectFactory = func(db *sql.DB) gorm.Dialector {
			return mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true})
		}
	case "postgres":
		driverName = "pgx"
		dialectFactory = func(db *sql.DB) gorm.Dialector {
			return postgres.New(postgres.Config{Conn: db})
		}
	default:
		return nil, nil, fmt.Errorf("不支持的数据库驱动: %q", driver)
	}
	sqlDB, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("创建数据库连接池失败: %w", err)
	}
	sqlDB.SetMaxOpenConns(pool.MaxOpenConns)
	sqlDB.SetMaxIdleConns(pool.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(pool.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(pool.ConnMaxIdleTime) * time.Second)

	pingCtx, cancel := context.WithTimeout(ctx, time.Duration(connectTimeout)*time.Second)
	defer cancel()
	if err = sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("数据库连接检测失败: %w", err)
	}
	return sqlDB, dialectFactory(sqlDB), nil
}

func resolveDSN(endpoint config.DatabaseEndpointConfig) (string, error) {
	if err := endpoint.Validate(); err != nil {
		return "", err
	}
	if endpoint.DSN != "" {
		return endpoint.DSN, nil
	}
	dsn, exists := os.LookupEnv(endpoint.DSNEnv)
	if !exists || strings.TrimSpace(dsn) == "" {
		return "", fmt.Errorf("数据库环境变量 %q 为空", endpoint.DSNEnv)
	}
	return dsn, nil
}
