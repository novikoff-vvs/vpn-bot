package user

import (
	"bot-service/config"
	"errors"
	"fmt"
	"pkg/exceptions"
	"resty.dev/v3"
)

type CreateUserRequest struct {
	ChatId int64  `json:"chat_id"`
	Email  string `json:"email"`
	UUID   string `gorm:"id"`
}

type GetUserByChatIdRequest struct {
	ChatId int64 `json:"chat_id"`
}

type Client struct {
	client *resty.Client
	cfg    config.UserService
}

func (c Client) Create(req CreateUserRequest) (*resty.Response, error) {
	//TODO добавить логирование
	response, err := c.client.R().SetBody(req).Post("user/create")
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		return response, errors.New(response.String())
	}
	return response, nil
}

func (c Client) GetByChatID(req GetUserByChatIdRequest) (*resty.Response, error) {
	// TODO: добавить логирование
	response, err := c.client.R().
		SetBody(req).
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("user/by-chat/%d", req.ChatId))

	if err != nil {
		return nil, err
	}

	if response.StatusCode() == 404 {
		return nil, exceptions.ErrModelNotFound
	}

	if response.IsError() {
		return nil, fmt.Errorf("unexpected status: %d, body: %s", response.StatusCode(), response.String())
	}

	return response, nil
}

func (c Client) Close() error {
	err := c.client.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewUserClient(cfg config.UserService) *Client {
	client := resty.New()

	client = client.SetBaseURL(cfg.Url)
	client = client.EnableTrace()

	return &Client{
		client: client,
		cfg:    cfg,
	}
}
