package singleton

import (
	"bot-service/config"
	"bot-service/internal/bot/message"
	"sync"
)

var builder *message.Builder
var once sync.Once

func messageBuilderBoot(cfg config.PaymentService) {
	once.Do(func() {
		builder = message.NewSendMessageCallBuilder(cfg)
	})
}

func MessageBuilder() *message.Builder {
	return builder
}
