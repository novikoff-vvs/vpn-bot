package user

import (
	"errors"
	"time"
	"user-service/internal/models"
	sqliteSubscription "user-service/internal/repository/subscription/sqlite"
	sqliteUser "user-service/internal/repository/user/sqlite"
)

type Service struct {
	userRepo         sqliteUser.UserRepositoryInterface
	subscriptionRepo sqliteSubscription.SubscriptionRepositoryInterface
}

func NewUserService(userRepo sqliteUser.UserRepositoryInterface, subscriptionRepo sqliteSubscription.SubscriptionRepositoryInterface) *Service {
	return &Service{
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s Service) CreateUser(user *models.User) (string, error) {
	s.userRepo.BeginTransaction()
	s.subscriptionRepo.BeginTransaction()

	uuid, err := s.userRepo.Create(user)
	if err != nil {
		_ = s.userRepo.RollbackTransaction()
		_ = s.subscriptionRepo.RollbackTransaction()
		return "", err
	}
	var subscription = models.Subscription{
		UserUUID:  uuid,
		PlanCode:  models.SubscriptionTestCode,
		StartedAt: time.Now(),
		ExpiresAt: time.Now().AddDate(0, 0, 1),
		IsActive:  true,
		AutoRenew: false,
	}

	id, err := s.subscriptionRepo.Create(&subscription)
	if err != nil {
		_ = s.userRepo.RollbackTransaction()
		_ = s.subscriptionRepo.RollbackTransaction()
		return "", err
	}

	if id == 0 {
		_ = s.userRepo.RollbackTransaction()
		_ = s.subscriptionRepo.RollbackTransaction()
		return "", errors.New("cannot create subscription")
	}
	err = s.userRepo.CommitTransaction()
	if err != nil {
		return "", err
	}

	err = s.subscriptionRepo.CommitTransaction()
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (s Service) GetByChatId(id int64) (*models.User, error) {
	user, err := s.userRepo.GetByChatId(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s Service) GetByUUID(uuid string) (*models.User, error) {
	user, err := s.userRepo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
