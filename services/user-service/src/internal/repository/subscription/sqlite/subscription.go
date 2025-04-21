package sqlite

import (
	grm "gorm.io/gorm"
	"pkg/infrastructure/DB/gorm"
	"time"
	"user-service/internal/models"
)

type SubscriptionRepository struct {
	dbService *gorm.DBService
}

// Транзакции

func (r *SubscriptionRepository) BeginTransaction() {
	r.dbService.Begin()
}

func (r *SubscriptionRepository) CommitTransaction() error {
	return r.dbService.Commit()
}

func (r *SubscriptionRepository) RollbackTransaction() error {
	return r.dbService.Rollback()
}

// Создание новой подписки
func (r *SubscriptionRepository) Create(subscription *models.Subscription) (uint, error) {
	tx := r.dbService.ActiveTx().Create(subscription)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return subscription.ID, nil
}

// Получение активной подписки по UUID пользователя
func (r *SubscriptionRepository) GetActiveByUserUUID(userUUID string) (*models.Subscription, error) {
	var sub models.Subscription
	tx := r.dbService.ActiveTx().Where("user_uuid = ? AND is_active = ? AND expires_at > ?", userUUID, true, time.Now()).
		Order("expires_at DESC").
		First(&sub)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &sub, nil
}

// Продлить подписку (обновить expires_at)
func (r *SubscriptionRepository) Extend(subscriptionID uint, duration time.Duration) error {
	tx := r.dbService.ActiveTx().Model(&models.Subscription{}).
		Where("id = ?", subscriptionID).
		Update("expires_at", grm.Expr("expires_at + ?", duration))

	return tx.Error
}

// Деактивировать подписку
func (r *SubscriptionRepository) Deactivate(subscriptionID uint) error {
	tx := r.dbService.ActiveTx().Model(&models.Subscription{}).
		Where("id = ?", subscriptionID).
		Update("is_active", false)

	return tx.Error
}

func NewSubscriptionRepository(dbService *gorm.DBService) *SubscriptionRepository {
	return &SubscriptionRepository{
		dbService: dbService,
	}
}
