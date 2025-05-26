package subscription

import (
	"encoding/json"
	"fmt"
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

func NewSubscriptionClient(cfg config.UserService) *Client {
	client := resty.New()

	client = client.SetBaseURL(cfg.Url)
	client = client.EnableTrace()

	return &Client{
		client: client,
		cfg:    cfg,
	}
}
