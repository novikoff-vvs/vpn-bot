package singleton

import (
	"log"
	"pkg/config"
	"pkg/infrastructure/nats"
	"sync"
)

var natsPublisher *nats.Publisher
var onceNatsPublisher sync.Once

func NatsPublisherBoot(cfg config.NatsPublisher) {
	onceNatsPublisher.Do(func() {
		var err error
		natsPublisher, err = nats.NewPublisher(cfg.Url)
		if err != nil {
			log.Println(err.Error())
		}
	})
}

func NatsPublisher() *nats.Publisher {
	return natsPublisher
}
