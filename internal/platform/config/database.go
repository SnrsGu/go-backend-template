package config

import (
	"errors"
	"fmt"
	"time"
)

type DatabaseConfig struct {
	Driver                string `json:"driver"`
	MaxOpenConnections    int    `json:"max_open_connections"`
	MaxIdleConnections    int    `json:"max_idle_connections"`
	ConnectionMaxLifetime int    `json:"connection_max_lifetime_seconds"`
	ConnectionMaxIdleTime int    `json:"connection_max_idle_time_seconds"`
	ConnectTimeoutSeconds int    `json:"connect_timeout_seconds"`
}

func (d DatabaseConfig) Validate() error {
	switch d.Driver {
	case "mysql":
	default:
		return fmt.Errorf(`invalid driver "%s"`, d.Driver)
	}
	if d.MaxOpenConnections <= 0 {
		return errors.New("max open connections must be greater than zero")
	}
	if d.MaxIdleConnections < 0 || d.MaxIdleConnections > d.MaxOpenConnections {
		return errors.New("max idle connections must be greater than or equal to max idle connections")
	}
	if d.ConnectionMaxLifetime <= 0 {
		return errors.New("connection max lifetime must be greater than zero")
	}
	if d.ConnectionMaxIdleTime <= 0 {
		return errors.New("connection max idle time must be greater than zero")
	}
	if d.ConnectTimeoutSeconds <= 0 {
		return errors.New("connect timeout seconds must be greater than zero")
	}
	return nil
}

func (d DatabaseConfig) ConnectionMaxLifetimeDuration() time.Duration {
	return time.Duration(d.ConnectionMaxLifetime) * time.Second
}

func (d DatabaseConfig) ConnectionMaxIdleTimeDuration() time.Duration {
	return time.Duration(d.ConnectionMaxIdleTime) * time.Second
}

func (d DatabaseConfig) ConnectTimeout() time.Duration {
	return time.Duration(d.ConnectTimeoutSeconds) * time.Second
}
