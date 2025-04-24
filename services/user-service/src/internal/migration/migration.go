package migration

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"user-service/config"
	"user-service/internal/models"
	"user-service/internal/seeds"
)

func InitDBConnection(cfg config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.DB, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Subscription{}, &models.Plan{})
	if err != nil {
		return nil, err
	}

	// Добавим тарифы, если они ещё не добавлены
	if err := seeds.SeedPlans(db); err != nil {
		return nil, err
	}

	return db, nil
}
