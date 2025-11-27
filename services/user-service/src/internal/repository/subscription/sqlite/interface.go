package sqlite

import (
	"user-service/internal/models"
)

type SubscriptionRepositoryInterface interface {
	Create(subscription *models.Subscription) (uint, error)
	GetActiveByUserUUID(userUUID string) (*models.Subscription, error)
	GetByUserUUID(userUUID string) (*models.Subscription, error)
	Extend(subscription *models.Subscription) error
	Deactivate(subscriptionID uint) error
	GetExpiringInDays(days int) ([]models.Subscription, error)

	BeginTransaction()
	CommitTransaction() error
	RollbackTransaction() error
}
