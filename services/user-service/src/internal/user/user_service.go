package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"log"
	"pkg/events"
	"pkg/singleton"
	"time"
	"user-service/internal/models"
	sqliteSubscription "user-service/internal/repository/subscription/sqlite"
	sqliteUser "user-service/internal/repository/user/sqlite"
)

type Service struct {
	userRepo         sqliteUser.UserRepositoryInterface
	subscriptionRepo sqliteSubscription.SubscriptionRepositoryInterface
}

func (s Service) CreateUser(user *models.User) (string, error) {
	s.userRepo.BeginTransaction()
	s.subscriptionRepo.BeginTransaction()

	uuid, err := s.userRepo.Create(user)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {

		_ = s.subscriptionRepo.RollbackTransaction()
		user.DeletedAt = gorm.DeletedAt(sql.NullTime{})
		activated, err := s.userRepo.Activate(user)
		if err != nil {
			log.Println(err)
			return "", err
		}

		err = s.userRepo.CommitTransaction()
		if err != nil {
			log.Println(err)
			return "", err
		}
		log.Println("Activated!")
		return activated, nil
	}

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

func (s Service) DeleteByUUID(uuid string) error {
	s.userRepo.BeginTransaction()
	err := s.userRepo.DeleteByUUID(uuid)
	if err != nil {
		errs := s.userRepo.RollbackTransaction()
		if errs != nil {
			return errs
		}
		return err
	}
	err = s.userRepo.CommitTransaction()
	if err != nil {
		return err
	}
	return nil
}

func (s Service) SyncUsers(UUIDs []string) ([]string, error) {
	err, exiting := s.userRepo.GetAllUUIDs(UUIDs)
	if err != nil {
		return nil, err
	}

	if len(exiting) == 0 {
		return nil, nil
	}

	for _, uuid := range exiting {
		marshal, err := json.Marshal(events.UserDeactivated{
			UserUUID: uuid,
			ChatId:   0,
		})
		if err != nil {
			return nil, err
		}
		err = singleton.NatsPublisher().Publish("events.user.deactivated", marshal)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		err = s.DeleteByUUID(uuid)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}

	return exiting, nil
}

func NewUserService(userRepo sqliteUser.UserRepositoryInterface, subscriptionRepo sqliteSubscription.SubscriptionRepositoryInterface) *Service {
	return &Service{
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
	}
}
