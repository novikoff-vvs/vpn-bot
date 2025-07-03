package config

import (
	pkg_config "pkg/config"
)

type Config struct {
	Logger      Logger                 `mapstructure:"logger"`
	Base        Base                   `mapstructure:"base"`
	Crypto      Crypto                 `mapstructure:"crypto"`
	UserService pkg_config.UserService `mapstructure:"user_service"`
	Database    pkg_config.Database    `mapstructure:"database"`
}

type Logger struct {
	Path     string `mapstructure:"path"`
	Name     string `mapstructure:"name"`
	IsOutput bool   `mapstructure:"is_output"`
}

type Base struct {
	AppPort       string `mapstructure:"app_port"`
	PaymentSecret string `mapstructure:"payment_secret"`
}

type Crypto struct {
	Key string `mapstructure:"key"`
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
