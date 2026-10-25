package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"

	"fleet/internal/adapter/mqtt"
	"fleet/internal/adapter/postgres"
	"fleet/internal/adapter/rabbitmq"
	"fleet/internal/adapter/rest"
	"fleet/internal/config"
	"fleet/internal/core/service"
)

func main() {
	cfg := config.Load()

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}

	publisher, err := rabbitmq.NewPublisher(conn)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewLocationService(postgres.NewLocationRepository(db), publisher, cfg.Geofence)

	log.Fatal(err)
	if err := mqtt.NewSubscriber(cfg.MQTTBroker, "fleet-server", svc).Start(); err != nil {
	}

	r := gin.Default()
	rest.NewHandler(svc).Routes(r)
	log.Fatal(r.Run(cfg.HTTPAddr))
}
