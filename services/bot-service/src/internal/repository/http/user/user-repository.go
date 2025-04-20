package user

import (
	user_client "bot-service/internal/infrastructure/client/user"
	"bot-service/internal/models"
	"encoding/json"
)

type RepositoryInterface interface {
	GetUser(id string) (models.User, error)
	GetUserByChatId(chatId int64) (models.User, error)
	CreateUser(user *models.User) error
}

type HTTPUserRepository struct {
	client *user_client.Client
}

func (r *HTTPUserRepository) GetUserByChatId(chatId int64) (models.User, error) {
	var request = user_client.GetUserByChatIdRequest{ChatId: chatId}
	response, err := r.client.GetByChatID(request)
	if err != nil {
		return models.User{}, err
	}
	var user = models.User{}
	err = json.Unmarshal(response.Bytes(), &user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func NewHTTPUserRepository(client *user_client.Client) *HTTPUserRepository {
	return &HTTPUserRepository{
		client: client,
	}
}

// GetUser получает пользователя по ID
func (r *HTTPUserRepository) GetUser(id string) (models.User, error) {
	return models.User{}, nil
}

func (r *HTTPUserRepository) CreateUser(user *models.User) error {
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
