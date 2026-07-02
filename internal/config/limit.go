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
	DbConfig DBConfig
)

type AppConfig struct {
	Services ServiceMap `yaml:"services"`
	DB       DBConfig   `yaml:"db"`
}

type DBConfig struct {
	Host                   string `yaml:"host"`
	Port                   int    `yaml:"port"`
	User                   string `yaml:"user"`
	Password               string `yaml:"password"`
	DBName                 string `yaml:"dbname"`
	SSLMode                string `yaml:"sslmode"`
	MaxOpenConns           int    `yaml:"max_open_conns"`
	MaxIdleConns           int    `yaml:"max_idle_conns"`
	ConnMaxLifetimeSeconds int    `yaml:"conn_max_lifetime_seconds"`
	ConnMaxIdleTimeSeconds int    `yaml:"conn_max_idle_time_seconds"`
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
	if maxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS"); maxOpenConnsStr != "" {
		if maxOpenConns, err := strconv.Atoi(maxOpenConnsStr); err == nil {
			c.MaxOpenConns = maxOpenConns
		}
	}
	if maxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS"); maxIdleConnsStr != "" {
		if maxIdleConns, err := strconv.Atoi(maxIdleConnsStr); err == nil {
			c.MaxIdleConns = maxIdleConns
		}
	}
	if connMaxLifetimeStr := os.Getenv("DB_CONN_MAX_LIFETIME_SECONDS"); connMaxLifetimeStr != "" {
		if connMaxLifetime, err := strconv.Atoi(connMaxLifetimeStr); err == nil {
			c.ConnMaxLifetimeSeconds = connMaxLifetime
		}
	}
	if connMaxIdleTimeStr := os.Getenv("DB_CONN_MAX_IDLE_TIME_SECONDS"); connMaxIdleTimeStr != "" {
		if connMaxIdleTime, err := strconv.Atoi(connMaxIdleTimeStr); err == nil {
			c.ConnMaxIdleTimeSeconds = connMaxIdleTime
		}
	}
}

func (c *DBConfig) ApplyDefaults() {
	if c.Port == 0 {
		c.Port = 5432
	}
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 10
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 5
	}
	if c.ConnMaxLifetimeSeconds == 0 {
		c.ConnMaxLifetimeSeconds = 300
	}
	if c.ConnMaxIdleTimeSeconds == 0 {
		c.ConnMaxIdleTimeSeconds = 60
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
	DbConfig = appConfig.DB
	DbConfig.ApplyDefaults()
	DbConfig.LoadFromEnv()

	return nil
}
