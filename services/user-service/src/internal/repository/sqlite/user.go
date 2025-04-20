package sqlite

import (
	"gorm.io/gorm"
	"user-service/internal/models"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r UserRepository) Create(user *models.User) (string, error) {
	tx := r.DB.Create(user)

	if tx.Error != nil {
		return "", tx.Error
	}
	return user.UUID, nil
}

func (r UserRepository) GetByUUID(uuid string) (*models.User, error) {
	var user models.User
	user.UUID = uuid
	tx := r.DB.First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func (r UserRepository) GetByChatId(chatId int64) (*models.User, error) {
	var user models.User
	user.ChatId = chatId
	tx := r.DB.First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}
