package config

import (
	pkg_config "pkg/config"
)

type Config struct {
	Logger      Logger      `mapstructure:"logger"`
	BotSettings BotSettings `mapstructure:"bot_settings"`
	UserService UserService `mapstructure:"user_service"`
}

type Logger struct {
	Path     string `mapstructure:"path"`
	Name     string `mapstructure:"name"`
	IsOutput bool   `mapstructure:"is_output"`
}

type UserService struct {
	Url string `mapstructure:"url"`
}

type BotSettings struct {
	Token string `mapstructure:"token"`
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
