package seeds

import (
	"gorm.io/gorm"
	"user-service/internal/models"
)

func SeedPlans(db *gorm.DB) error {
	// Проверим, есть ли уже тарифы
	var count int64
	if err := db.Model(&models.Plan{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil // Уже есть тарифы — не дублируем
	}

	plans := []models.Plan{
		{
			Code:         models.SubscriptionTestCode,
			Name:         "Тестовый",
			Price:        100,
			DurationDays: 1,
			MaxDevices:   1,
			Description:  "Первый тариф после авторизации",
		},
		{
			Code:         models.SubscriptionBaseCode,
			Name:         "Базовый",
			Price:        110,
			DurationDays: 30,
			MaxDevices:   1,
			Description:  "Один пользователь, базовые функции",
		},
	}

	// Вставка
	return db.Create(&plans).Error
}
