package config

import (
	pkg_config "pkg/config"
)

type Config struct {
	Logger      Logger                 `mapstructure:"logger"`
	Xui         Xui                    `mapstructure:"xui"`
	Base        Base                   `mapstructure:"base"`
	UserService pkg_config.UserService `mapstructure:"user_service"`
}

type Logger struct {
	Path     string `mapstructure:"path"`
	Name     string `mapstructure:"name"`
	IsOutput bool   `mapstructure:"is_output"`
}

type Xui struct {
	BaseURL   string `mapstructure:"base_url"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	InboundID int    `mapstructure:"inbound_id"`
}

type Base struct {
	AppPort string `mapstructure:"APP_PORT"`
	Swagger bool   `mapstructure:"SWAGGER"`
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
