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
	Address           string `yaml:"address"`
	ReadHeaderTimeout int    `yaml:"read_header_timeout"`
	IdleTimeout       int    `yaml:"idle_timeout"`
	ShutdownTimeout   int    `yaml:"shutdown_timeout"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Address:           ":8080",
			ReadHeaderTimeout: 5,
			IdleTimeout:       60,
			ShutdownTimeout:   10,
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
	if c.Server.ReadHeaderTimeout <= 0 {
		return fmt.Errorf("server.read_header_timeout必须大于 0")
	}

	if c.Server.IdleTimeout <= 0 {
		return fmt.Errorf("server.idle_timeout必须大于 0")
	}

	if c.Server.ShutdownTimeout <= 0 {
		return fmt.Errorf("server.shutdown_timeout必须大于 0")
	}
	return nil
}
