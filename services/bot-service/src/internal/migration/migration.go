package migration

import (
	"bot-service/internal/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	pkg_config "pkg/config"
)

func InitDBConnection(cfg pkg_config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.DB, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.NotifyUser{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
