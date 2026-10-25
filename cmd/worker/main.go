package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"fleet/internal/adapter/rabbitmq"
	"fleet/internal/config"
	"fleet/internal/core/domain"
)

func main() {
	cfg := config.Load()

	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := rabbitmq.NewConsumer(conn)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(consumer.Consume(func(e domain.GeofenceEvent) {
		log.Printf("client entered the area: %+v", e)
	}))
}
