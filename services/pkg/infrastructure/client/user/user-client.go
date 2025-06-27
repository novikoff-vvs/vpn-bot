package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"pkg/config"
	"pkg/exceptions"
	response2 "pkg/infrastructure/client/user/response"
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

type GetUserByUUIDRequest struct {
	UUID string `json:"uuid"`
}

type SyncUsersRequest struct {
	UUIDs []string `json:"uuids"`
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
func (c Client) GetUserByUUID(req GetUserByUUIDRequest) (*resty.Response, error) {
	// TODO: добавить логирование
	response, err := c.client.R().
		SetBody(req).
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("user/%s", req.UUID))

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
func (c Client) SyncUsers(req SyncUsersRequest) (*resty.Response, error) {
	// TODO: добавить логирование
	response, err := c.client.R().SetBody(req).Post("user/sync-users")

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
func (c Client) All() (response2.AllResponse, error) {
	// TODO: добавить логирование
	response, err := c.client.R().Get("user/all")

	if err != nil {
		return response2.AllResponse{}, err
	}

	if response.IsError() {
		return response2.AllResponse{}, fmt.Errorf("unexpected status: %d, body: %s", response.StatusCode(), response.String())
	}
	var resp response2.AllResponse
	err = json.Unmarshal(response.Bytes(), &resp)

	if err != nil {
		return response2.AllResponse{}, err
	}

	return resp, nil
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
