package user

import (
	"bot-service/internal/models"
	usrRepo "bot-service/internal/repository/http/user"
	"bot-service/internal/vpn"
	"errors"
)

type ServiceInterface interface {
	UserExistsByChatId(chatId int64) bool
	UserRegisterByChatId(chatId int64, comment, uuid string) error
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

func (u Service) UserRegisterByChatId(chatId int64, comment, uuid string) error {
	err := u.vpnService.UserRegisterByChatId(chatId, comment, uuid)
	if err != nil {
		return err
	}
	user, err := u.vpnService.UserGetByChatId(chatId)
	if err != nil {
		return err
	}
	err = u.userRepo.CreateUser(&user)
	if err != nil {
		return err
	}
	return nil
}

func (u Service) UserGetByChatId(chatId int64) (models.User, error) {
	var user models.User
	user, err := u.userRepo.GetUserByChatId(chatId)

	if err != nil || user.ChatId == 0 {
		user, err = u.vpnService.UserGetByChatId(chatId)
		if err != nil {
			return models.User{}, err
		}
	}
	if user.ChatId == 0 {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
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
