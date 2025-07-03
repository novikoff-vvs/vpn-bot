package subscription

import (
	"errors"
	"gorm.io/gorm"
	"pkg/infrastructure/client/vpn"
	"strconv"
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
	vpnClient   *vpn.Client
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
		er := s.repo.RollbackTransaction()
		if er != nil {
			return nil, er
		}
		return nil, err
	}
	if activeSubscription == nil {
		activeSubscription, err = s.repo.GetByUserUUID(dto.UserUUID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			er := s.repo.RollbackTransaction()
			if er != nil {
				return nil, er
			}
			return nil, err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = s.repo.RollbackTransaction()
			return nil, err
		}
	}

	if activeSubscription == nil {
		er := s.repo.RollbackTransaction()
		if er != nil {
			return nil, er
		}
		return nil, err
	}

	upgradedPlan, err := s.planService.GetUpgradedPlan(activeSubscription.Plan)
	if err != nil {
		er := s.repo.RollbackTransaction()
		if er != nil {
			return nil, er
		}
		return nil, err
	}

	now := time.Now()
	expiresAt := activeSubscription.ExpiresAt
	if expiresAt.Before(now) {
		expiresAt = now
	}
	expiresAt = expiresAt.AddDate(0, 0, upgradedPlan.DurationDays*dto.AmountPeriods)

	activeSubscription.Plan = upgradedPlan
	activeSubscription.IsActive = true
	activeSubscription.ExpiresAt = expiresAt //TODO возможно нужно обновлять started_at
	oldSubscriptionId := activeSubscription.ID
	activeSubscription.ID = 0
	err = s.repo.CommitTransaction()
	if err != nil {
		return nil, err
	}
	s.repo.BeginTransaction()

	err = s.repo.Deactivate(oldSubscriptionId)
	if err != nil {
		er := s.repo.RollbackTransaction()
		if er != nil {
			return nil, er
		}
		return nil, err
	}

	err = s.repo.Extend(activeSubscription)
	if err != nil {
		er := s.repo.RollbackTransaction()
		if er != nil {
			return nil, er
		}
		return nil, err
	}

	_, err = s.vpnClient.UpdateClient(activeSubscription.UserUUID, vpn.UpdateClientRequest{
		Email:          activeSubscription.User.Email,
		TotalGB:        0,
		ExpiryTimeUnix: activeSubscription.ExpiresAt.UnixMilli(),
		Enable:         true,
		TgId:           strconv.FormatInt(activeSubscription.User.ChatId, 10),
	})

	if err != nil {
		er := s.repo.RollbackTransaction()
		if er != nil {
			return nil, er
		}
		return nil, err
	}

	err = s.repo.CommitTransaction()
	if err != nil {
		return nil, err
	}

	return activeSubscription, nil
}

func NewSubscriptionService(repo *sqlite.SubscriptionRepository, planService *plan.Service, vpnClient *vpn.Client) *Service {
	return &Service{repo: repo, planService: planService, vpnClient: vpnClient}
}
