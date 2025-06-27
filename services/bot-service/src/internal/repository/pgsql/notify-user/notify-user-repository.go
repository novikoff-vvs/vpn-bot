package notify_user

import (
	"bot-service/internal/model"
	"gorm.io/gorm"
)

type NotifyUserRepository struct {
	db *gorm.DB
}

func (r NotifyUserRepository) All() ([]model.NotifyUser, error) {
	var users []model.NotifyUser
	err := r.db.Model(&model.NotifyUser{}).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r NotifyUserRepository) Create(chatId int64) error {
	return r.db.Model(&model.NotifyUser{}).Create(&model.NotifyUser{ChatId: chatId}).Error
}

func NewNotifyUserRepository(db *gorm.DB) *NotifyUserRepository {
	return &NotifyUserRepository{db: db}
}
