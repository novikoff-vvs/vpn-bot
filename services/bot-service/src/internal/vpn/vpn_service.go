package vpn

import (
	"bot-service/config"
	"bot-service/internal/models"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/novikoff-vvs/logger"
	"github.com/novikoff-vvs/xui"
	"github.com/novikoff-vvs/xui/dto"
	"github.com/novikoff-vvs/xui/requests"
	"strconv"
	"time"
)

type ServiceInterface interface {
	UserExistsByChatId(chatId int64) bool
	UserRegisterByChatId(chatId int64, comment, uuid string) error
	UserGetByChatId(chatId int64) (models.User, error)
	UserGetByEmail(email string) (models.User, error)
	ResetClientTraffic(chatId int64) error
}

type Service struct {
	client *xui.Client
	cfg    config.VpnService
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

func (s Service) UserRegisterByChatId(chatId int64, comment, phone string) error {
	uuId := uuid.New()
	clients := requests.AddClientToInboundClientRequest{
		Clients: []dto.Client{
			{
				Comment:    comment,
				Email:      phone,
				Enable:     true,
				ExpiryTime: time.Now().AddDate(0, 0, 1).UnixMilli(),
				Flow:       "",
				Id:         uuId.String(),
				LimitIp:    0,
				Reset:      0,
				SubId:      uuId.String(),
				TgId:       json.Number(strconv.FormatInt(chatId, 10)),
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

func (s Service) UserGetByChatId(chatId int64) (models.User, error) {
	inbound, err := s.client.GetInbound(s.cfg.InboundID)
	if err != nil {
		s.lg.Error(err.Error())
		return models.User{}, err
	}

	for _, client := range inbound.GetSettings().Clients {
		id, _ := client.TgId.Int64()
		if id == chatId {
			return models.User{
				ChatId:         id,
				Email:          client.Email,
				UUID:           client.Id,
				SubscriptionId: client.SubId,
			}, nil
		}
	}

	return models.User{}, nil
}

func (s Service) UserGetByEmail(email string) (models.User, error) {
	_, err := s.client.GetUserByEmail(requests.GetUserByEmailRequest{Email: email})
	if err != nil {
		s.lg.Error(err.Error())
		return models.User{}, err
	}
	return models.User{}, nil
}

func NewVPNService(cfg config.VpnService, lg logger.Interface) (*Service, error) {
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
