package vpn

import (
	"encoding/json"
	"github.com/novikoff-vvs/logger"
	"github.com/novikoff-vvs/xui"
	"github.com/novikoff-vvs/xui/dto"
	"github.com/novikoff-vvs/xui/requests"
	"log"
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
	GetAllUsers() ([]models.VpnUser, error)
	UpdateClient(userUUID string, dto UpdateClientDTO) error
	UserGetByChatUUID(uuid string) (models.VpnUser, error)
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

type UpdateClientsDTO struct {
	Clients []UpdateClientDTO `json:"clients"`
}

type UpdateClientDTO struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	TotalGB        int64  `json:"totalGB"`
	ExpiryTime     int64  `json:"expiryTime"`
	Enable         bool   `json:"enable"`
	TgID           string `json:"tgId"`
	Comment        string `json:"comment"`
	SubscriptionId string `json:"subId"`
}

func (s Service) UpdateClient(userUUID string, dto UpdateClientDTO) error {
	marshal, err := json.Marshal(UpdateClientsDTO{Clients: []UpdateClientDTO{dto}})
	if err != nil {
		return err
	}
	log.Println(string(marshal))
	err = s.client.UpdateClient(userUUID, s.cfg.InboundID, string(marshal))
	if err != nil {
		return err
	}
	return nil
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

func (s Service) UserGetByChatUUID(uuid string) (models.VpnUser, error) {
	inbound, err := s.client.GetInbound(s.cfg.InboundID)
	if err != nil {
		s.lg.Error(err.Error())
		return models.VpnUser{}, err
	}

	for _, client := range inbound.GetSettings().Clients {
		id, _ := client.TgId.Int64()
		if client.Id == uuid {
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

func (s Service) GetAllUsers() ([]models.VpnUser, error) {
	var result []models.VpnUser
	inbound, err := s.client.GetInbound(s.cfg.InboundID)
	if err != nil {
		s.lg.Error(err.Error())
		return result, err
	}

	for _, client := range inbound.GetSettings().Clients {
		id, err := client.TgId.Int64()
		if err != nil {
			s.lg.Error(err.Error())
		}
		user := models.VpnUser{
			ChatId:         id,
			Email:          client.Email,
			UUID:           client.Id,
			SubscriptionId: client.SubId,
		}
		result = append(result, user)
	}

	return result, nil
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
