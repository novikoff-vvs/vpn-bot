package singleton

import "bot-service/config"

func Boot(cfg *config.Config) {
	messageBuilderBoot(cfg.PaymentService)
}
