package migration

import (
	"gorm.io/gorm"
	"user-service/config"
	"user-service/internal/models"
)
import "gorm.io/driver/sqlite"

func InitDBConnection(cfg config.Database) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
