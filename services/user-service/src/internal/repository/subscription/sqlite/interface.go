package sqlite

import (
	"time"
	"user-service/internal/models"
)

type SubscriptionRepositoryInterface interface {
	Create(subscription *models.Subscription) (uint, error)
	GetActiveByUserUUID(userUUID string) (*models.Subscription, error)
	Extend(subscriptionID uint, duration time.Duration) error
	Deactivate(subscriptionID uint) error

	BeginTransaction()
	CommitTransaction() error
	RollbackTransaction() error
}
