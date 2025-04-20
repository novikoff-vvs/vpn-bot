package config

import (
	"log"
	pkg_config "pkg/config"
)

type Config struct {
	Logger      Logger      `mapstructure:"logger"`
	VpnService  VpnService  `mapstructure:"vpn_service"`
	BotSettings BotSettings `mapstructure:"bot_settings"`
	UserService UserService `mapstructure:"user_service"`
}

type Logger struct {
	Path     string `mapstructure:"path"`
	Name     string `mapstructure:"name"`
	IsOutput bool   `mapstructure:"is_output"`
}

type VpnService struct {
	BaseURL   string `mapstructure:"base_url"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
	InboundID int    `mapstructure:"inbound_id"`
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

	if config.BotSettings.Token == "" {
		log.Fatal("Bot token is required! Set it in config.yaml or BOT_TOKEN in .env")
	}

	return &config, nil
}
