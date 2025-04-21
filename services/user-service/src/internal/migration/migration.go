package migration

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"user-service/config"
	"user-service/internal/models"
)

func InitDBConnection(cfg config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.DB, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Subscription{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
