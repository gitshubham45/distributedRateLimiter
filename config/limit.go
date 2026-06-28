package config

import (
	"os"

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

var Services ServiceMap

type AppConfig struct {
	Services ServiceMap `yaml:"services"`
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
	return nil
}
