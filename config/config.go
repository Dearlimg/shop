package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Server   ServerConfig   `yaml:"server"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`     // Redis地址，格式: host:port
	Password string `yaml:"password"` // Redis密码，空字符串表示无密码
	DB       int    `yaml:"db"`       // Redis数据库编号，默认0
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Database, c.Charset)
}

// GetAddr 获取服务器地址
func (c *ServerConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

var AppConfig *Config

// LoadConfig 从YAML文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 设置默认值
	if config.Database.Port == 0 {
		config.Database.Port = 3306
	}
	if config.Database.Charset == "" {
		config.Database.Charset = "utf8mb4"
	}
	if config.Redis.Addr == "" {
		config.Redis.Addr = "localhost:6379"
	}
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}

	AppConfig = &config
	return &config, nil
}
