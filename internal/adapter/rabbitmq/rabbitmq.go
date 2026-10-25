package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"

	"fleet/internal/core/domain"
)

const (
	Exchange      = "fleet.events"
	GeofenceQueue = "geofence_alerts"
)

func setup(ch *amqp.Channel) error {
	if err := ch.ExchangeDeclare(Exchange, "topic", true, false, false, false, nil); err != nil {
		return err
	}
	if _, err := ch.QueueDeclare(GeofenceQueue, true, false, false, false, nil); err != nil {
		return err
	}
	return ch.QueueBind(GeofenceQueue, "geofence.*", Exchange, false, nil)
}

type Publisher struct {
	ch *amqp.Channel
}

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{ch: ch}, setup(ch)
}

func (p *Publisher) PublishGeofenceEvent(ctx context.Context, event domain.GeofenceEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.ch.PublishWithContext(ctx, Exchange, "geofence.entry", false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

type Consumer struct {
	ch *amqp.Channel
}

func NewConsumer(conn *amqp.Connection) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Consumer{ch: ch}, setup(ch)
}

func (c *Consumer) Consume(handle func(domain.GeofenceEvent)) error {
	msgs, err := c.ch.Consume(GeofenceQueue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for d := range msgs {
		var event domain.GeofenceEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			d.Nack(false, false)
			continue
		}
		handle(event)
		d.Ack(false)
	}
	return errors.New("channel closed")
}
