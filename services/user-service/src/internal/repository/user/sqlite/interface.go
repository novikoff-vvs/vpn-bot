package sqlite

import (
	"user-service/internal/models"
)

type UserRepositoryInterface interface {
	Create(user *models.User) (string, error)
	GetByUUID(uuid string) (*models.User, error)
	GetByChatId(chatId int64) (*models.User, error)

	BeginTransaction()
	CommitTransaction() error
	RollbackTransaction() error
}
