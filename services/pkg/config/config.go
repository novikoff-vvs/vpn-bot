package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type UserService struct {
	Url string `mapstructure:"url"`
}

type VpnService struct {
	Url string `mapstructure:"url"`
}

func GetViper() *viper.Viper {
	// 1. Загружаем .env (опционально, если не используется в Docker/K8s)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS environment variables")
	}
	v := viper.New()
	// Настраиваем Viper
	v.SetConfigName("config")     // имя файла без расширения
	v.SetConfigType("yaml")       // формат конфига
	v.AddConfigPath("./configs/") // путь к файлу (текущая директория)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return v
}
