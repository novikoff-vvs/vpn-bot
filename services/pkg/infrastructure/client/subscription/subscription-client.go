package subscription

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pkg/config"
	"pkg/exceptions"
	"pkg/infrastructure/client/subscription/responses"
	"resty.dev/v3"
)

type Client struct {
	client *resty.Client
	cfg    config.UserService
}

type GetSubscriptionByUUIDRequest struct {
	UUID string
}

func (c Client) GetSubscriptionByUUID(req GetSubscriptionByUUIDRequest) (responses.GetSubscriptionResponse, error) {
	// TODO: добавить логирование
	response, err := c.client.R().
		SetBody(req).
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("subscription/by-user/%s", req.UUID))

	if err != nil {
		return responses.GetSubscriptionResponse{}, err
	}

	if response.StatusCode() == 404 {
		return responses.GetSubscriptionResponse{}, exceptions.ErrModelNotFound
	}

	if response.IsError() {
		return responses.GetSubscriptionResponse{}, fmt.Errorf("unexpected status: %d, body: %s", response.StatusCode(), response.String())
	}
	var result responses.GetSubscriptionResponse

	err = json.Unmarshal(response.Bytes(), &result)
	if err != nil {
		return responses.GetSubscriptionResponse{}, err
	}
	return result, nil
}

type RefreshRequest struct {
	UserUUID      string `json:"user_uuid"`
	AmountPeriods int    `json:"amount_periods"`
}

func (c Client) RefreshSubscription(req RefreshRequest) (responses.GetSubscriptionResponse, error) {

	response, err := c.client.R().
		SetBody(req).
		SetHeader("Accept", "application/json").
		SetResult(&responses.GetSubscriptionResponse{}).
		Post("subscription/refresh") // Предполагаемый эндпоинт

	if err != nil {
		return responses.GetSubscriptionResponse{}, err
	}

	if response.StatusCode() == http.StatusNotFound {
		return responses.GetSubscriptionResponse{}, exceptions.ErrModelNotFound
	}

	if response.IsError() {
		errMsg := fmt.Sprintf("refresh failed with status %d: %s",
			response.StatusCode(), response.String())
		return responses.GetSubscriptionResponse{}, fmt.Errorf(errMsg)
	}

	return *response.Result().(*responses.GetSubscriptionResponse), nil
}

func NewSubscriptionClient(cfg config.UserService) *Client {
	client := resty.New()

	client = client.SetBaseURL(cfg.Url)
	client = client.EnableTrace()

	return &Client{
		client: client,
		cfg:    cfg,
	}
}
