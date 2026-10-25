package config

import (
	"os"
	"strconv"

	"fleet/internal/core/domain"
)

type Config struct {
	HTTPAddr      string
	PublisherAddr string
	DatabaseURL   string
	MQTTBroker    string
	RabbitMQURL   string
	Geofence      domain.Geofence
}

func Load() Config {
	return Config{
		HTTPAddr:      env("HTTP_ADDR", ":8080"),
		PublisherAddr: env("PUBLISHER_ADDR", ":8081"),
		DatabaseURL:   env("DATABASE_URL", "postgres://fleet:fleet@localhost:5432/fleet?sslmode=disable"),
		MQTTBroker:    env("MQTT_BROKER", "tcp://localhost:1883"),
		RabbitMQURL:   env("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		Geofence: domain.Geofence{
			Center: domain.Point{
				Latitude:  envFloat("GEOFENCE_LAT", -6.236737656443795),
				Longitude: envFloat("GEOFENCE_LNG", 106.85633139688976),
			},
			Radius: envFloat("GEOFENCE_RADIUS", 50),
		},
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envFloat(key string, def float64) float64 {
	v, err := strconv.ParseFloat(os.Getenv(key), 64)
	if err != nil {
		return def
	}
	return v
}
