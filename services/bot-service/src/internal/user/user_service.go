package user

import (
	usrRepo "bot-service/internal/repository/http/user"
	"bot-service/internal/repository/http/vpn"
	"bot-service/internal/singleton"
	"errors"
	"github.com/google/uuid"
	"pkg/exceptions"
	"pkg/models"
)

type ServiceInterface interface {
	UserExistsByChatId(chatId int64) bool
	UserRegisterByChatId(chatId int64, comment, phone string) (models.VpnUser, error)
	UserGetByChatId(chatId int64) (models.VpnUser, error)
	UserGetByEmail(email string) (models.VpnUser, error)
	ResetClientTraffic(chatId int64) error
}

type Service struct {
	vpnService vpn.RepositoryInterface
	userRepo   usrRepo.RepositoryInterface
}

func (u Service) UserExistsByChatId(chatId int64) bool {
	//TODO implement me
	panic("implement me")
}

func (u Service) UserRegisterByChatId(chatId int64, comment, phone string) (models.VpnUser, error) {
	uuId := uuid.New()
	var user = models.VpnUser{
		ChatId:         chatId,
		Email:          phone,
		UUID:           uuId.String(),
		SubscriptionId: "",
	}

	err := u.userRepo.CreateUser(&user)
	if err != nil {
		return models.VpnUser{}, err
	}

	err = u.vpnService.CreateUser(&user)
	if err != nil {
		return models.VpnUser{}, err
	}

	return user, nil
}

func (u Service) UserGetByChatId(chatId int64) (models.VpnUser, error) {
	var err error
	var user models.VpnUser
	cachedUser, ok := singleton.UserContainer().Get(chatId)
	user = cachedUser.User
	if !ok {
		user, err = u.userRepo.GetUserByChatId(chatId)
		if err == nil {
			return user, nil
		}

		user, err = u.vpnService.GetUserByChatId(chatId)
		if err == nil {
			return user, nil
		}
		if errors.Is(err, exceptions.ErrModelNotFound) {
			return models.VpnUser{}, err
		}
	}

	return models.VpnUser{}, err
}

func (u Service) UserGetByEmail(email string) (models.VpnUser, error) {
	//TODO implement me
	panic("implement me")
}

func (u Service) ResetClientTraffic(chatId int64) error {
	//TODO implement me
	panic("implement me")
}

func NewUserService(vpnService vpn.RepositoryInterface, userRepo usrRepo.RepositoryInterface) *Service {
	return &Service{
		vpnService: vpnService,
		userRepo:   userRepo,
	}
}
