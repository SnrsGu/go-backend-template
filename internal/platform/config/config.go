package config

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Server ServerConfig `json:"server"`
}

type ServerConfig struct {
	Address                  string `json:"address"`
	ReadHeaderTimeoutSeconds int    `json:"read_header_timeout_seconds"`
	ShutdownTimeoutSeconds   int    `json:"shutdown_timeout_seconds"`
}

func Load(filePath string) (Config, error) {
	var configuration Config
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return configuration, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	if err := json.Unmarshal(fileContents, &configuration); err != nil {
		return configuration, fmt.Errorf("failed to unmarshal file %s: %w", filePath, err)
	}
	if err := configuration.Validate(); err != nil {
		return configuration, fmt.Errorf("failed to validate file %s: %w", filePath, err)
	}
	return configuration, nil
}

func (c Config) Validate() error {
	if c.Server.Address == "" {
		return errors.New("server address cannot be empty")
	}
	if c.Server.ReadHeaderTimeoutSeconds <= 0 {
		return errors.New("read header timeout must be positive")
	}
	if c.Server.ShutdownTimeoutSeconds <= 0 {
		return errors.New("shutdown timeout must be positive")
	}
	return nil
}

func (s ServerConfig) ReadHeaderTimeout() time.Duration {
	return time.Duration(s.ReadHeaderTimeoutSeconds) * time.Second
}

func (s ServerConfig) ShutdownTimeout() time.Duration {
	return time.Duration(s.ShutdownTimeoutSeconds) * time.Second
}
