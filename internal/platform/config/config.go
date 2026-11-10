package config

import (
	"encoding/json/v2"
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

func Load(filePath string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}
	if err := json.Unmarshal(file, &config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	if err := config.Validate(); err != nil {
		return config, fmt.Errorf("failed to validate config file: %w", err)
	}
	return config, nil
}

func (c Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return err
	}
	return nil
}
