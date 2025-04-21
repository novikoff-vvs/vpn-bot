package user

import (
	"bot-service/internal/models"
	usrRepo "bot-service/internal/repository/http/user"
	"bot-service/internal/vpn"
	"errors"
	"github.com/google/uuid"
	"pkg/exceptions"
)

type ServiceInterface interface {
	UserExistsByChatId(chatId int64) bool
	UserRegisterByChatId(chatId int64, comment, phone string) (models.User, error)
	UserGetByChatId(chatId int64) (models.User, error)
	UserGetByEmail(email string) (models.User, error)
	ResetClientTraffic(chatId int64) error
}

type Service struct {
	vpnService vpn.ServiceInterface
	userRepo   usrRepo.RepositoryInterface
}

func (u Service) UserExistsByChatId(chatId int64) bool {
	//TODO implement me
	panic("implement me")
}

func (u Service) UserRegisterByChatId(chatId int64, comment, phone string) (models.User, error) {
	uuId := uuid.New()
	var user = models.User{
		ChatId:         chatId,
		Email:          phone,
		UUID:           uuId.String(),
		SubscriptionId: "",
	}

	err := u.userRepo.CreateUser(&user)
	if err != nil {
		return models.User{}, err
	}

	err = u.vpnService.UserRegisterByChatId(&user, comment)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u Service) UserGetByChatId(chatId int64) (models.User, error) {
	user, err := u.userRepo.GetUserByChatId(chatId)
	if err == nil {
		return user, nil
	}

	user, err = u.vpnService.UserGetByChatId(chatId)
	if err == nil {
		return user, nil
	}

	if errors.Is(err, exceptions.ErrModelNotFound) {
		return models.User{}, err
	}

	return models.User{}, err
}

func (u Service) UserGetByEmail(email string) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u Service) ResetClientTraffic(chatId int64) error {
	//TODO implement me
	panic("implement me")
}

func NewUserService(vpnService vpn.ServiceInterface, userRepo usrRepo.RepositoryInterface) *Service {
	return &Service{
		vpnService: vpnService,
		userRepo:   userRepo,
	}
}
