package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Log    LogConfig    `yaml:"log"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Address: ":8080",
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
	}
}

func Load(path string) (*Config, error) {
	cfg := defaultConfig()

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开配置文件失败: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("解码配置文件失败: %w", err)
	}

	applyEnv(&cfg)

	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("验证配置错误: %w", err)
	}

	return &cfg, nil
}

func applyEnv(cfg *Config) {
	if value, exists := os.LookupEnv("GBT_SERVER_ADDRESS"); exists {
		cfg.Server.Address = value
	}
	if value, exists := os.LookupEnv("GBT_SERVER_LOG_LEVEL"); exists {
		cfg.Log.Level = value
	}
	if value, exists := os.LookupEnv("GBT_SERVER_LOG_FORMAT"); exists {
		cfg.Log.Format = value
	}
}

func (c *Config) Validate() error {
	address := strings.TrimSpace(c.Server.Address)
	if address == "" {
		return fmt.Errorf("server.address不能为空")
	}

	_, portString, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("无效的server.address: %q %w", address, err)
	}

	port, err := strconv.Atoi(portString)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("无效的server.address: port %q", portString)
	}
	return nil
}
