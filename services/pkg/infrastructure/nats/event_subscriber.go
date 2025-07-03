package nats

import (
	"log"

	"github.com/nats-io/nats.go"
)

type EventHandler func(data []byte)

func StartSubscriber(natsURL, subject string, handler EventHandler) error {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return err
	}

	_, err = conn.Subscribe(subject, func(msg *nats.Msg) {
		log.Printf("Received event on [%s]: %s", subject, string(msg.Data))
		handler(msg.Data)
	})
	return err
}

func HandleUserCreated(data []byte) {
	log.Printf("Handling user created event: %s", string(data))
}
