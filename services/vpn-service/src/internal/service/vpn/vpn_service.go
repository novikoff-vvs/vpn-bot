package vpn

import (
	"encoding/json"
	"github.com/novikoff-vvs/logger"
	"github.com/novikoff-vvs/xui"
	"github.com/novikoff-vvs/xui/dto"
	"github.com/novikoff-vvs/xui/requests"
	"pkg/exceptions"
	"pkg/models"
	"strconv"
	"time"
	"vpn-service/config"
)

type ServiceInterface interface {
	UserExistsByChatId(chatId int64) bool
	UserRegisterByChatId(user *models.VpnUser, comment string) error
	UserGetByChatId(chatId int64) (models.VpnUser, error)
	UserGetByEmail(email string) (models.VpnUser, error)
	ResetClientTraffic(chatId int64) error
}

type Service struct {
	client *xui.Client
	cfg    config.Xui
	lg     logger.Interface
}

func (s Service) ResetClientTraffic(chatId int64) error {
	user, err := s.UserGetByChatId(chatId)
	if err != nil {
		return err
	}

	err = s.client.ResetClientTraffic(s.cfg.InboundID, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) UserExistsByChatId(chatId int64) bool {
	inbound, err := s.client.GetInbound(s.cfg.InboundID)
	if err != nil { //TODO добавить логирование
		return false
	}
	for _, client := range inbound.GetSettings().Clients {
		id, _ := client.TgId.Int64()
		if id == chatId {
			return true
		}
	}

	return false
}

func (s Service) UserRegisterByChatId(user *models.VpnUser, comment string) error {
	clients := requests.AddClientToInboundClientRequest{
		Clients: []dto.Client{
			{
				Comment:    comment,
				Email:      user.Email,
				Enable:     true,
				ExpiryTime: time.Now().AddDate(0, 0, 1).UnixMilli(),
				Flow:       "",
				Id:         user.UUID,
				LimitIp:    0,
				Reset:      0,
				SubId:      user.UUID,
				TgId:       json.Number(strconv.FormatInt(user.ChatId, 10)),
				TotalGB:    2 * 1024 * 1024 * 1024,
			},
		},
	}

	clientsString, err := json.Marshal(clients)
	if err != nil {
		s.lg.Error(err.Error())
		return err
	}

	request := requests.AddClientToInboundRequest{
		InboundId: 2,
		Settings:  string(clientsString),
	}

	err = s.client.AddClientToInbound(request)
	if err != nil {
		s.lg.Error(err.Error())
		return err
	}
	return nil
}

func (s Service) UserGetByChatId(chatId int64) (models.VpnUser, error) {
	inbound, err := s.client.GetInbound(s.cfg.InboundID)
	if err != nil {
		s.lg.Error(err.Error())
		return models.VpnUser{}, err
	}

	for _, client := range inbound.GetSettings().Clients {
		id, _ := client.TgId.Int64()
		if id == chatId {
			return models.VpnUser{
				ChatId:         id,
				Email:          client.Email,
				UUID:           client.Id,
				SubscriptionId: client.SubId,
			}, nil
		}
	}

	return models.VpnUser{}, exceptions.ErrModelNotFound
}

func (s Service) UserGetByEmail(email string) (models.VpnUser, error) {
	_, err := s.client.GetUserByEmail(requests.GetUserByEmailRequest{Email: email})
	if err != nil {
		s.lg.Error(err.Error())
		return models.VpnUser{}, err
	}
	return models.VpnUser{}, nil
}

func NewVPNService(cfg config.Xui, lg logger.Interface) (*Service, error) {
	client := xui.NewClient(cfg.BaseURL, cfg.Username, cfg.Password)
	err := client.Login()
	if err != nil {
		return nil, err
	}

	return &Service{
		client: client,
		lg:     lg,
		cfg:    cfg,
	}, nil
}
