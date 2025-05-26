package subscription

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"user-service/internal/models"
	"user-service/internal/plan"
	"user-service/internal/repository/subscription/sqlite"
)

// todo вынести
type RefreshDTO struct {
	UserUUID      string
	AmountPeriods int
}

type Service struct {
	repo        *sqlite.SubscriptionRepository
	planService *plan.Service
}

func (s Service) GetSubscriptionByUser(uuid string) (*models.Subscription, error) {
	return s.repo.GetByUserUUID(uuid)
}

func (s Service) GetActiveSubscriptionByUser(uuid string) (*models.Subscription, error) {
	return s.repo.GetActiveByUserUUID(uuid)
}

func (s Service) Refresh(dto RefreshDTO) (*models.Subscription, error) {
	s.repo.BeginTransaction()

	activeSubscription, err := s.repo.GetActiveByUserUUID(dto.UserUUID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		err = s.repo.RollbackTransaction()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	if activeSubscription == nil {
		activeSubscription, err = s.repo.GetByUserUUID(dto.UserUUID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			err = s.repo.RollbackTransaction()
			if err != nil {
				return nil, err
			}
			return nil, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = s.repo.RollbackTransaction()
			return nil, err
		}
	}

	if activeSubscription == nil {
		err = s.repo.RollbackTransaction()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	upgradedPlan, err := s.planService.GetUpgradedPlan(activeSubscription.Plan)
	if err != nil {
		err = s.repo.RollbackTransaction()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	now := time.Now()
	expiresAt := activeSubscription.ExpiresAt
	if expiresAt.Before(now) {
		expiresAt = now
	}
	expiresAt = expiresAt.Add(time.Duration(upgradedPlan.DurationDays * dto.AmountPeriods))

	activeSubscription.Plan = upgradedPlan
	activeSubscription.IsActive = true
	activeSubscription.ExpiresAt = expiresAt //TODO возможно нужно обновлять started_at

	err = s.repo.Extend(activeSubscription)
	if err != nil {
		err = s.repo.RollbackTransaction()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	err = s.repo.CommitTransaction()
	if err != nil {
		return nil, err
	}
	return activeSubscription, nil
}

func NewSubscriptionService(repo *sqlite.SubscriptionRepository, planService *plan.Service) *Service {
	return &Service{repo: repo, planService: planService}
}
