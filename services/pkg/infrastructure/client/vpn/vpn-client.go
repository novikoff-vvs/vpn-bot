package vpn

import (
	"errors"
	"fmt"
	"github.com/novikoff-vvs/logger"
	"pkg/config"
	"pkg/exceptions"
	"resty.dev/v3"
)

type RegisterUserRequest struct {
	ChatId int64  `json:"chat_id"`
	Email  string `json:"email"`
	UUID   string `json:"uuid"`
}

type RegisterUserResponse struct {
	Message string `json:"message"`
}

type GetVpnUserResponse struct {
	UUID   string `json:"uuid"`
	Email  string `json:"email"`
	ChatId int64  `json:"chat_id"`
}

type ExistsResponse struct {
	Exists bool `json:"exists"`
}

type Client struct {
	client *resty.Client
	cfg    config.VpnService
	lg     logger.Interface
}

func NewVpnClient(cfg config.VpnService, lg logger.Interface) *Client {
	client := resty.New().
		SetBaseURL(cfg.Url).
		EnableTrace()

	return &Client{
		client: client,
		cfg:    cfg,
		lg:     lg,
	}
}

func (c Client) RegisterUser(req RegisterUserRequest) (*RegisterUserResponse, error) {
	c.lg.Debug(fmt.Sprintf("Registering VPN user: %+v", req))

	resp, err := c.client.R().
		SetBody(req).
		SetResult(&RegisterUserResponse{}).
		Post("/vpn/register")

	if err != nil {
		c.lg.Error("register request failed: " + err.Error())
		return nil, err
	}

	if resp.IsError() {
		errMsg := fmt.Sprintf("register response error: %s", resp.String())
		c.lg.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	return resp.Result().(*RegisterUserResponse), nil
}

func (c Client) GetByChatID(chatId int64) (*GetVpnUserResponse, error) {
	c.lg.Debug(fmt.Sprintf("Getting VPN user by chatId: %d", chatId))

	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		SetResult(&GetVpnUserResponse{}).
		Get(fmt.Sprintf("/vpn/by-chat/%d", chatId))

	if err != nil {
		c.lg.Error("get user request failed: " + err.Error())
		return nil, err
	}

	if resp.StatusCode() == 404 {
		c.lg.Error(fmt.Sprintf("User not found: chatId=%d", chatId))
		return nil, exceptions.ErrModelNotFound
	}

	if resp.IsError() {
		errMsg := fmt.Sprintf("get user response error: %d %s", resp.StatusCode(), resp.String())
		c.lg.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	return resp.Result().(*GetVpnUserResponse), nil
}

func (c Client) ResetTraffic(chatId int64) error {
	c.lg.Info(fmt.Sprintf("Resetting traffic for chatId=%d", chatId))

	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		Post(fmt.Sprintf("/vpn/reset-traffic/%d", chatId))

	if err != nil {
		c.lg.Error("reset traffic request failed: " + err.Error())
		return err
	}

	if resp.IsError() {
		errMsg := fmt.Sprintf("reset traffic failed: %s", resp.String())
		c.lg.Error(errMsg)
		return errors.New(errMsg)
	}

	c.lg.Info(fmt.Sprintf("Traffic successfully reset for chatId=%d", chatId))
	return nil
}

func (c Client) ExistsByChatID(chatId int64) (bool, error) {
	c.lg.Debug(fmt.Sprintf("Checking if VPN user exists by chatId: %d", chatId))

	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		SetResult(&ExistsResponse{}).
		Get(fmt.Sprintf("/vpn/exists/%d", chatId))

	if err != nil {
		c.lg.Error("exists request failed: " + err.Error())
		return false, err
	}

	if resp.IsError() {
		errMsg := fmt.Sprintf("exists response error: %s", resp.String())
		c.lg.Error(errMsg)
		return false, errors.New(errMsg)
	}

	return resp.Result().(*ExistsResponse).Exists, nil
}

func (c Client) GetByEmail(email string) (*GetVpnUserResponse, error) {
	c.lg.Debug(fmt.Sprintf("Getting VPN user by email: %s", email))

	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		SetResult(&GetVpnUserResponse{}).
		Get(fmt.Sprintf("/vpn/by-email/%s", email))

	if err != nil {
		c.lg.Error("get by email request failed: " + err.Error())
		return nil, err
	}

	if resp.StatusCode() == 404 {
		c.lg.Info(fmt.Sprintf("User not found by email: %s", email))
		return nil, exceptions.ErrModelNotFound
	}

	if resp.IsError() {
		errMsg := fmt.Sprintf("get by email error: %s", resp.String())
		c.lg.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	return resp.Result().(*GetVpnUserResponse), nil
}

func (c Client) Close() error {
	return c.client.Close()
}
