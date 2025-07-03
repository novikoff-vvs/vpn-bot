package nats

import (
	"github.com/nats-io/nats.go"
	"strings"
)

type Publisher struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func NewPublisher(natsURL string) (*Publisher, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, err
	}

	streamName := "EVENTS"

	_, err = js.AddStream(&nats.StreamConfig{
		Name:      streamName,
		Subjects:  []string{"events.>"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
	})

	if err != nil && !strings.Contains(err.Error(), "stream name already in use") {
		conn.Close()
		return nil, err
	}

	return &Publisher{
		conn: conn,
		js:   js,
	}, nil
}

// Publish публикует сообщение через JetStream
func (p *Publisher) Publish(subject string, data []byte) error {
	_, err := p.js.Publish(subject, data)
	return err
}

// Subscribe подписывается на subject через JetStream с ручным подтверждением
func (p *Publisher) Subscribe(subject string, durable string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return p.js.Subscribe(subject, handler, nats.Durable(durable), nats.ManualAck())
}

// Close закрывает соединение
func (p *Publisher) Close() {
	p.conn.Drain()
	p.conn.Close()
}
