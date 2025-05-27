package nats

import (
	"encoding/json"
	"log"
	"pkg/singleton"

	"github.com/nats-io/nats.go"

	"bot-service/internal/bot"
	"pkg/events"
)

type SubscriptionHandler struct {
	botService *bot.Service
}

func NewSubscriptionHandler(botService *bot.Service) *SubscriptionHandler {
	return &SubscriptionHandler{botService: botService}
}

func (h *SubscriptionHandler) SubscribeToEvents() ([]*nats.Subscription, error) {
	var subs []*nats.Subscription

	sub1, err := singleton.NatsPublisher().Subscribe("events.subscription.refreshed", "bot_service_subscription_refreshed_consumer", h.handleSubscriptionRefreshed)
	if err != nil {
		return nil, err
	}
	subs = append(subs, sub1)

	sub2, err := singleton.NatsPublisher().Subscribe("events.user.deactivated", "bot_service_user_deactivate_consumer", h.handleUserDeactivated)
	if err != nil {
		return nil, err
	}
	subs = append(subs, sub2)

	return subs, nil
}

func (h *SubscriptionHandler) handleSubscriptionRefreshed(msg *nats.Msg) {
	var event events.SubscriptionRefreshed
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Error unmarshalling event: %v", err)
		return
	}

	log.Println("Получено событие обновления подписки")

	if err := h.botService.NotifySubscriptionRefreshed(event); err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}

	if err := msg.Ack(); err != nil {
		log.Println("Ошибка при подтверждении сообщения:", err)
	}
}

func (h *SubscriptionHandler) handleUserDeactivated(msg *nats.Msg) {
	var event events.UserDeactivated
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Error unmarshalling event: %v", err)
		return
	}

	log.Println("Получено событие деактивации пользователя ")

	if err := h.botService.NotifyDeactivatedUser(event); err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}

	if err := msg.Ack(); err != nil {
		log.Println("Ошибка при подтверждении сообщения:", err)
	}
}
