package singleton

import (
	"pkg/config"
	"pkg/infrastructure/client/subscription"
	"sync"
)

var subscriptionClient *subscription.Client
var onceSubscriptions sync.Once

func SubscriptionClientBoot(cfg config.UserService) {
	onceSubscriptions.Do(func() {
		subscriptionClient = subscription.NewSubscriptionClient(cfg)
	})
}

func SubscriptionClient() *subscription.Client {
	return subscriptionClient
}
