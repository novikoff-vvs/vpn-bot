package user

import (
	"encoding/json"
	"errors"
	"pkg/exceptions"
	user_client "pkg/infrastructure/client/user"
	"pkg/models"
)

type RepositoryInterface interface {
	GetUser(id string) (models.VpnUser, error)
	GetUserByChatId(chatId int64) (models.VpnUser, error)
	CreateUser(user *models.VpnUser) error
	All() ([]models.VpnUser, error)
}

type HTTPUserRepository struct {
	client *user_client.Client
}

func (r *HTTPUserRepository) GetUserByChatId(chatId int64) (models.VpnUser, error) {
	var request = user_client.GetUserByChatIdRequest{ChatId: chatId}
	response, err := r.client.GetByChatID(request)
	if errors.Is(err, exceptions.ErrModelNotFound) {
		return models.VpnUser{}, err
	}
	if err != nil {
		return models.VpnUser{}, err
	}
	var user = models.VpnUser{}
	err = json.Unmarshal(response.Bytes(), &user)
	if err != nil {
		return models.VpnUser{}, err
	}
	return user, nil
}

func NewHTTPUserRepository(client *user_client.Client) *HTTPUserRepository {
	return &HTTPUserRepository{
		client: client,
	}
}

func (r *HTTPUserRepository) GetUser(id string) (models.VpnUser, error) {
	return models.VpnUser{}, nil
}

func (r *HTTPUserRepository) CreateUser(user *models.VpnUser) error {
	//TODO добавить логирование
	var req = user_client.CreateUserRequest{
		ChatId: user.ChatId,
		Email:  user.Email,
		UUID:   user.UUID,
	}
	_, err := r.client.Create(req)
	if err != nil {
		return err
	}

	return nil
}

func (r *HTTPUserRepository) All() ([]models.VpnUser, error) {
	var req, err = r.client.All()
	var result []models.VpnUser
	if err != nil {
		return result, err
	}
	for _, user := range req.Result {
		result = append(result, models.VpnUser{
			ChatId: int64(user.ChatId),
			Email:  user.Email,
			UUID:   user.UUID,
		})
	}

	return result, nil
}
