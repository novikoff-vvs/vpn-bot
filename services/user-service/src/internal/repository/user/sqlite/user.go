package sqlite

import (
	"errors"
	grm "gorm.io/gorm"
	"pkg/infrastructure/DB/gorm"
	"user-service/internal/models"
)

type UserRepository struct {
	dbService *gorm.DBService
}

// Обёртки над транзакциями через dbService

func (r *UserRepository) BeginTransaction() {
	r.dbService.Begin()
}

func (r *UserRepository) CommitTransaction() error {
	return r.dbService.Commit()
}

func (r *UserRepository) RollbackTransaction() error {
	return r.dbService.Rollback()
}

// Методы репозитория

func (r *UserRepository) Create(user *models.User) (string, error) {
	tx := r.dbService.ActiveTx()
	if tx == nil {
		tx = r.dbService.DB()
	}

	if err := tx.Create(user).Error; err != nil {
		return "", err
	}
	return user.UUID, nil
}

func (r *UserRepository) Activate(user *models.User) (string, error) {
	tx := r.dbService.ActiveTx()
	if tx == nil {
		tx = r.dbService.DB()
	}

	if err := tx.Unscoped().Model(user).Update("deleted_at", nil).Error; err != nil {
		return "", err
	}
	return user.UUID, nil
}

func (r *UserRepository) GetByUUID(uuid string) (*models.User, error) {
	var user models.User
	user.UUID = uuid

	tx := r.dbService.ActiveTx()
	if tx == nil {
		tx = r.dbService.DB()
	}

	if err := tx.
		Preload("Subscription").
		Preload("Subscription.Plan").
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByChatId(chatId int64) (*models.User, error) {
	var user models.User
	user.ChatId = chatId

	tx := r.dbService.ActiveTx()
	if tx == nil {
		tx = r.dbService.DB()
	}

	err := tx.Where("chat_id = ?", chatId).First(&user).Error
	if errors.Is(err, grm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAllUUIDs(uuids []string) ([]string, error) {
	var existingUUIDs []string
	err := r.dbService.DB().
		Model(&models.User{}).
		Where("uuid NOT IN ?", uuids).
		Pluck("uuid", &existingUUIDs).Error
	return existingUUIDs, err
}

func (r *UserRepository) DeleteByUUID(uuid string) error {
	tx := r.dbService.ActiveTx()
	if tx == nil {
		tx = r.dbService.DB()
	}

	result := tx.Unscoped().Where("uuid = ?", uuid).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func NewUserRepository(dbService *gorm.DBService) *UserRepository {
	return &UserRepository{
		dbService: dbService,
	}
}
