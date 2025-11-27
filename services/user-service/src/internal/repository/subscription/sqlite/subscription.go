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

// GetExpiringInDays возвращает активные подписки, которые истекают через указанное количество дней
func (r *SubscriptionRepository) GetExpiringInDays(days int) ([]models.Subscription, error) {
	var subs []models.Subscription
	
	now := time.Now()
	// Вычисляем дату через N дней (начало дня)
	targetDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, days)
	// Начало и конец дня для точного сравнения
	startOfDay := targetDate
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	// Ищем активные подписки, которые истекают в указанный день
	// Проверяем, что expires_at попадает в диапазон указанного дня (включительно начало, исключительно конец)
	tx := r.dbService.DB().
		Where("deleted_at IS NULL").
		Where("is_active = ?", true).
		Where("expires_at >= ? AND expires_at < ?", startOfDay, endOfDay).
		Preload("Plan").
		Preload("User").
		Find(&subs)
	
	if tx.Error != nil {
		return nil, tx.Error
	}
	
	return subs, nil
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
