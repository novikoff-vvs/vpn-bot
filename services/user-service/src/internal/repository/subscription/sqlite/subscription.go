package sqlite

import (
	"fmt"
	gorm2 "gorm.io/gorm"
	"log"
	"pkg/infrastructure/DB/gorm"
	"time"
	"user-service/internal/models"
)

type SubscriptionRepository struct {
	dbService *gorm.DBService
}

func (r *SubscriptionRepository) BeginTransaction() {
	r.dbService.Begin()
}

func (r *SubscriptionRepository) CommitTransaction() error {
	return r.dbService.Commit()
}

func (r *SubscriptionRepository) RollbackTransaction() error {
	return r.dbService.Rollback()
}

func (r *SubscriptionRepository) Create(subscription *models.Subscription) (uint, error) {
	tx := r.dbService.ActiveTx().Create(subscription)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return subscription.ID, nil
}

func (r *SubscriptionRepository) GetActiveByUserUUID(userUUID string) (*models.Subscription, error) {
	var sub models.Subscription
	tx := r.dbService.
		ActiveTx().
		Scopes(ActiveSubscription).
		Where("user_uuid = ?", userUUID).
		Preload("Plan").
		Preload("User").
		Order("expires_at DESC").
		First(&sub)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &sub, nil
}

func (r *SubscriptionRepository) GetByUserUUID(userUUID string) (*models.Subscription, error) {
	var sub models.Subscription
	tx := r.dbService.
		DB().
		Unscoped().
		Where("user_uuid = ?", userUUID).
		Preload("Plan").
		Preload("User").
		Order("expires_at DESC").
		First(&sub)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &sub, nil
}

func (r *SubscriptionRepository) Extend(subscription *models.Subscription) error {
	tx := r.dbService.ActiveTx().
		Create(subscription)

	return tx.Error
}

func (r *SubscriptionRepository) Deactivate(subscriptionID uint) error {
	tx := r.dbService.ActiveTx().Model(&models.Subscription{}).
		Unscoped().
		Where("id = ?", subscriptionID).
		Update("is_active", false)

	log.Println(fmt.Sprintf("Sub deactivate Row affected: %d", tx.RowsAffected))

	return tx.Error
}

func NewSubscriptionRepository(dbService *gorm.DBService) *SubscriptionRepository {
	return &SubscriptionRepository{
		dbService: dbService,
	}
}

func ActiveSubscription(db *gorm2.DB) *gorm2.DB {
	return db.
		Where("expires_at > ? AND deleted_at is NULL AND is_active = ? ", time.Now(), true)
}
