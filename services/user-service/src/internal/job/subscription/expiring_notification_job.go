package subscription

import (
	"encoding/json"
	"fmt"
	"log"
	"pkg/events"
	"pkg/singleton"
	"user-service/internal/repository/subscription/sqlite"

	"github.com/novikoff-vvs/logger"
)

type ExpiringNotificationJob struct {
	repo *sqlite.SubscriptionRepository
	lg   logger.Interface
}

func (j *ExpiringNotificationJob) Run() {
	j.lg.Info("Starting subscription expiring notification job")

	// Проверяем подписки, истекающие через 7, 3 и 1 день
	daysToCheck := []int{7, 3, 1}

	for _, days := range daysToCheck {
		j.lg.Info(fmt.Sprintf("Checking subscriptions expiring in %d days", days))
		subs, err := j.repo.GetExpiringInDays(days)
		if err != nil {
			j.lg.Error(fmt.Sprintf("Error getting subscriptions expiring in %d days: %s", days, err.Error()))
			continue
		}

		j.lg.Info(fmt.Sprintf("Found %d subscriptions expiring in %d days", len(subs), days))

		for _, sub := range subs {
			j.lg.Info(fmt.Sprintf("Processing subscription ID=%d, UserUUID=%s", sub.ID, sub.UserUUID))

			if sub.User == nil {
				j.lg.Error(fmt.Sprintf("Subscription %d (UserUUID: %s) has no user - skipping", sub.ID, sub.UserUUID))
				continue
			}

			if sub.User.ChatId == 0 {
				j.lg.Error(fmt.Sprintf("Subscription %d (UserUUID: %s) has invalid ChatId (0) - skipping", sub.ID, sub.UserUUID))
				continue
			}

			event := events.SubscriptionExpiring{
				UserUUID:      sub.UserUUID,
				ChatId:        sub.User.ChatId,
				DaysRemaining: days,
			}

			j.lg.Info(fmt.Sprintf("Creating event: UserUUID=%s, ChatId=%d, DaysRemaining=%d",
				event.UserUUID, event.ChatId, event.DaysRemaining))

			eventData, err := json.Marshal(event)
			if err != nil {
				j.lg.Error(fmt.Sprintf("Error marshalling event for subscription %d: %s", sub.ID, err.Error()))
				continue
			}

			j.lg.Info(fmt.Sprintf("Publishing event to NATS: %s", string(eventData)))

			err = singleton.NatsPublisher().Publish("events.subscription.expiring", eventData)
			if err != nil {
				j.lg.Error(fmt.Sprintf("Error publishing event for subscription %d: %s", sub.ID, err.Error()))
				log.Printf("Error publishing event: %s", err.Error())
				continue
			}

			j.lg.Info(fmt.Sprintf("Successfully published expiring notification event for user %s (ChatId: %d, Days: %d)",
				sub.UserUUID, sub.User.ChatId, days))
		}
	}

	j.lg.Info("Subscription expiring notification job completed")
}

func NewExpiringNotificationJob(repo *sqlite.SubscriptionRepository, lg logger.Interface) *ExpiringNotificationJob {
	return &ExpiringNotificationJob{
		repo: repo,
		lg:   lg,
	}
}
