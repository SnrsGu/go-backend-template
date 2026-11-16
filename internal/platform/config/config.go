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
	Server   ServerConfig   `yaml:"server"`
	Log      LogConfig      `yaml:"log"`
	Database DatabaseConfig `yaml:"database"`
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

type DatabaseConfig struct {
	Connections map[string]DatabaseConnectionConfig `yaml:"connections"`
}

type DatabaseConnectionConfig struct {
	Driver         string                   `yaml:"driver"`
	Primary        DatabaseEndpointConfig   `yaml:"primary"`
	Replicas       []DatabaseEndpointConfig `yaml:"replicas"`
	Pool           DatabasePoolConfig       `yaml:"pool"`
	ConnectTimeout int                      `yaml:"connect_timeout"`
	SlowThreshold  int                      `yaml:"slow_threshold"`
}

type DatabaseEndpointConfig struct {
	DSN    string `yaml:"dsn"`
	DSNEnv string `yaml:"dsn_env"`
}

type DatabasePoolConfig struct {
	MaxOpenConns    int `yaml:"max_open_conns"`
	MaxIdleConns    int `yaml:"max_idle_conns"`
	ConnMaxLifetime int `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int `yaml:"conn_max_idle_time"`
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
		}, Database: DatabaseConfig{
			Connections: map[string]DatabaseConnectionConfig{},
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
		return fmt.Errorf("server.address 不能为空")
	}

	_, portString, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("无效的 server.address: %q %w", address, err)
	}

	port, err := strconv.Atoi(portString)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("无效的 server.address: port %q", portString)
	}
	if c.Server.ReadHeaderTimeout <= 0 {
		return fmt.Errorf("server.read_header_timeout 必须大于 0")
	}

	if c.Server.IdleTimeout <= 0 {
		return fmt.Errorf("server.idle_timeout 必须大于 0")
	}

	if c.Server.ShutdownTimeout <= 0 {
		return fmt.Errorf("server.shutdown_timeout 必须大于 0")
	}
	if err = c.Database.Validate(); err != nil {
		return fmt.Errorf("数据库配置校验失败: %w", err)
	}
	return nil
}

func (c DatabaseConfig) Validate() error {
	for name, conn := range c.Connections {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("数据库名称不可为空")
		}
		if err := conn.Validate(); err != nil {
			return fmt.Errorf("数据库 %s 配置错误: %w", name, err)
		}
	}
	return nil
}

func (c DatabaseConnectionConfig) Validate() error {
	switch c.Driver {
	case "mysql", "postgres":
	default:
		return fmt.Errorf("不支持的数据库驱动： %s", c.Driver)
	}
	if err := c.Primary.Validate(); err != nil {
		return fmt.Errorf("主数据库配置错误: %w", err)
	}
	for idx, replica := range c.Replicas {
		if err := replica.Validate(); err != nil {
			return fmt.Errorf("从库 %d 配置错误: %w", idx+1, err)
		}
	}
	pool := c.Pool
	if pool.MaxOpenConns <= 0 {
		return fmt.Errorf("max_open_conns 必须大于 0")
	}
	if pool.MaxIdleConns < 0 ||
		pool.MaxIdleConns > pool.MaxOpenConns {
		return fmt.Errorf("max_idle_conns 超出合法范围")
	}
	if pool.ConnMaxLifetime <= 0 {
		return fmt.Errorf("conn_max_lifetime 必须大于 0")
	}
	if pool.ConnMaxIdleTime <= 0 {
		return fmt.Errorf("conn_max_idle_time 必须大于 0")
	}
	if c.ConnectTimeout <= 0 {
		return fmt.Errorf("connect_timeout 必须大于 0")
	}
	if c.SlowThreshold <= 0 {
		return fmt.Errorf("slow_threshold 必须大于 0")
	}
	return nil
}

func (c DatabaseEndpointConfig) Validate() error {
	hasDSN := strings.TrimSpace(c.DSN) != ""
	hasEnv := strings.TrimSpace(c.DSNEnv) != ""
	if hasDSN == hasEnv {
		return fmt.Errorf("dsn 和 dsn_env 必须且只能配置一个")
	}
	return nil
}
