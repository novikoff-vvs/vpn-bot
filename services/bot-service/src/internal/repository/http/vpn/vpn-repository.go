package vpn

import (
	"errors"
	"fmt"
	"pkg/exceptions"
	vpnclient "pkg/infrastructure/client/vpn"
	"pkg/models"
)

type RepositoryInterface interface {
	GetUserByChatId(chatId int64) (models.VpnUser, error)
	CreateUser(user *models.VpnUser) error
	ResetTraffic(chatId int64) error
	Exists(chatId int64) (bool, error)
	GetUserByEmail(email string) (models.VpnUser, error)
}

type HTTPVPNUserRepository struct {
	client *vpnclient.Client
}

func NewHTTPVPNUserRepository(client *vpnclient.Client) *HTTPVPNUserRepository {
	return &HTTPVPNUserRepository{client: client}
}

func (r *HTTPVPNUserRepository) GetUserByChatId(chatId int64) (models.VpnUser, error) {
	resp, err := r.client.GetByChatID(chatId)
	if err != nil {
		if errors.Is(err, exceptions.ErrModelNotFound) {
			return models.VpnUser{}, err
		}
		return models.VpnUser{}, fmt.Errorf("vpn repo: failed to get user: %w", err)
	}

	return models.VpnUser{
		ChatId: chatId,
		Email:  resp.Email,
		UUID:   resp.UUID,
	}, nil
}

func (r *HTTPVPNUserRepository) CreateUser(user *models.VpnUser) error {
	req := vpnclient.RegisterUserRequest{
		ChatId: user.ChatId,
		Email:  user.Email,
		UUID:   user.UUID,
	}

	_, err := r.client.RegisterUser(req)
	if err != nil {
		return fmt.Errorf("vpn repo: failed to create user: %w", err)
	}

	return nil
}

func (r *HTTPVPNUserRepository) ResetTraffic(chatId int64) error {
	if err := r.client.ResetTraffic(chatId); err != nil {
		return fmt.Errorf("vpn repo: failed to reset traffic: %w", err)
	}
	return nil
}

func (r *HTTPVPNUserRepository) Exists(chatId int64) (bool, error) {
	exists, err := r.client.ExistsByChatID(chatId)
	if err != nil {
		return false, fmt.Errorf("vpn repo: failed to check existence: %w", err)
	}
	return exists, nil
}

func (r *HTTPVPNUserRepository) GetUserByEmail(email string) (models.VpnUser, error) {
	resp, err := r.client.GetByEmail(email)
	if err != nil {
		if errors.Is(err, exceptions.ErrModelNotFound) {
			return models.VpnUser{}, err
		}
		return models.VpnUser{}, fmt.Errorf("vpn repo: failed to get user by email: %w", err)
	}

	return models.VpnUser{
		ChatId: resp.ChatId,
		Email:  resp.Email,
		UUID:   resp.UUID,
	}, nil
}
