package config

import (
	"errors"
	"time"
)

type ServerConfig struct {
	Address                  string `json:"address"`
	ReadHeaderTimeoutSeconds int    `json:"read_header_timeout_seconds"`
	ShutdownTimeoutSeconds   int    `json:"shutdown_timeout_seconds"`
}

func (s ServerConfig) Validate() error {
	if s.Address == "" {
		return errors.New("server address cannot be empty")
	}
	if s.ReadHeaderTimeoutSeconds <= 0 {
		return errors.New("server read_header_timeout_seconds must be greater than zero")
	}
	if s.ShutdownTimeoutSeconds <= 0 {
		return errors.New("server shutdown_timeout_seconds must be greater than zero")
	}
	return nil
}

func (s ServerConfig) ReadHeaderTimeout() time.Duration {
	return time.Duration(s.ReadHeaderTimeoutSeconds) * time.Second
}

func (s ServerConfig) ShutdownTimeout() time.Duration {
	return time.Duration(s.ShutdownTimeoutSeconds) * time.Second
}
