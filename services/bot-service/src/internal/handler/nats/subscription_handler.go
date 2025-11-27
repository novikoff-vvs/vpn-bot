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

	sub3, err := singleton.NatsPublisher().Subscribe("events.notify.new_message", "bot_service_notify_consumer", h.handleNotifyNewMessage)
	if err != nil {
		return nil, err
	}
	subs = append(subs, sub3)

	sub4, err := singleton.NatsPublisher().Subscribe("events.subscription.expiring", "bot_service_subscription_expiring_consumer", h.handleSubscriptionExpiring)
	if err != nil {
		return nil, err
	}
	subs = append(subs, sub4)

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

func (h *SubscriptionHandler) handleNotifyNewMessage(msg *nats.Msg) {
	var event events.NewMessage
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Error unmarshalling event: %v", err)
		return
	}

	log.Println("Получено событие отправки рассылки")

	if err := h.botService.NotifyNewMessage(event); err != nil {
		log.Println("Ошибка при отправке сообщения:", err)
	}

	if err := msg.Ack(); err != nil {
		log.Println("Ошибка при подтверждении сообщения:", err)
	}
}

func (h *SubscriptionHandler) handleSubscriptionExpiring(msg *nats.Msg) {
	log.Printf("Received raw message on events.subscription.expiring: %s", string(msg.Data))
	
	var event events.SubscriptionExpiring
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Error unmarshalling event: %v, raw data: %s", err, string(msg.Data))
		return
	}

	log.Printf("Получено событие об истекающей подписке для пользователя %s (ChatId: %d, Days: %d)", 
		event.UserUUID, event.ChatId, event.DaysRemaining)

	// Проверяем валидность данных
	if event.ChatId == 0 {
		log.Printf("ERROR: Invalid ChatId (0) in event for user %s", event.UserUUID)
		if err := msg.Ack(); err != nil {
			log.Printf("Error acknowledging message: %v", err)
		}
		return
	}

	if event.UserUUID == "" {
		log.Printf("ERROR: Empty UserUUID in event (ChatId: %d)", event.ChatId)
		if err := msg.Ack(); err != nil {
			log.Printf("Error acknowledging message: %v", err)
		}
		return
	}

	if err := h.botService.NotifySubscriptionExpiring(event); err != nil {
		log.Printf("Ошибка при отправке сообщения для пользователя %s (ChatId: %d): %v", 
			event.UserUUID, event.ChatId, err)
		if err := msg.Ack(); err != nil {
			log.Printf("Error acknowledging message after send error: %v", err)
		}
		return
	}

	log.Printf("Успешно отправлено уведомление пользователю %s (ChatId: %d, Days: %d)", 
		event.UserUUID, event.ChatId, event.DaysRemaining)

	if err := msg.Ack(); err != nil {
		log.Printf("Ошибка при подтверждении сообщения: %v", err)
	}
}
