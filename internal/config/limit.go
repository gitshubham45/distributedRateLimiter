package config

import (
	"os"
	"strconv"

	"github.com/goccy/go-yaml"
)

const (
	PER_UNIT_TIME = "per_unit_time"
	LEAKY_BUCKET  = "leaky_bucket"
)

type StrategyConfig struct {
	Name     string `json:"name"`
	Limit    int64  `json:"limit"`
	Interval int64  `json:"interval"`
}

type ServiceMap map[string]string

var (
	Services ServiceMap
	DB       DBConfig
)

type AppConfig struct {
	Services ServiceMap `yaml:"services"`
	DB       DBConfig   `yaml:"db"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// LoadFromEnv overrides configuration values with environment variables if present.
func (c *DBConfig) LoadFromEnv() {
	if host := os.Getenv("DB_HOST"); host != "" {
		c.Host = host
	}
	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			c.Port = port
		}
	}
	if user := os.Getenv("DB_USER"); user != "" {
		c.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		c.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		c.DBName = dbname
	}
	if sslmode := os.Getenv("DB_SSLMODE"); sslmode != "" {
		c.SSLMode = sslmode
	}
}

func InitConfig(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var appConfig AppConfig
	if err := yaml.Unmarshal(file, &appConfig); err != nil {
		return err
	}

	if appConfig.Services == nil {
		appConfig.Services = ServiceMap{}
	}

	Services = appConfig.Services
	DB = appConfig.DB
	DB.LoadFromEnv()
	return nil
}
