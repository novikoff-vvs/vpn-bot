package config

import (
	pkg_config "pkg/config"
)

type Config struct {
	Logger   Logger   `mapstructure:"logger"`
	Base     Base     `mapstructure:"base"`
	Database Database `mapstructure:"database"`
}

type Logger struct {
	Path     string `mapstructure:"path"`
	Name     string `mapstructure:"name"`
	IsOutput bool   `mapstructure:"is_output"`
}

type Base struct {
	AppPort string `mapstructure:"APP_PORT"`
	Swagger bool   `mapstructure:"SWAGGER"`
}

type Database struct {
	Path string `mapstructure:"PATH"`
}

func LoadConfigs() (*Config, error) {
	v := pkg_config.GetViper()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
