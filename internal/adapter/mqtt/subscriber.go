package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"

	"fleet/internal/core/domain"
	"fleet/internal/core/port"
)

const LocationTopic = "/fleet/vehicle/+/location"

type Subscriber struct {
	client paho.Client
	svc    port.LocationService
}

func NewSubscriber(broker, clientID string, svc port.LocationService) *Subscriber {
	s := &Subscriber{svc: svc}

	opts := paho.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(2 * time.Second).
		SetOnConnectHandler(func(c paho.Client) {
			c.Subscribe(LocationTopic, 1, s.handle)
		})

	s.client = paho.NewClient(opts)
	return s
}

func (s *Subscriber) Start() error {
	token := s.client.Connect()
	token.Wait()
	return token.Error()
}

func (s *Subscriber) handle(_ paho.Client, msg paho.Message) {
	loc, err := parseMessage(msg.Topic(), msg.Payload())
	if err != nil {
		log.Printf("invalid message on %s: %v", msg.Topic(), err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.svc.Record(ctx, loc); err != nil {
		log.Printf("record %s: %v", loc.VehicleID, err)
	}
}

func parseMessage(topic string, payload []byte) (domain.Location, error) {
	parts := strings.Split(topic, "/")
	if len(parts) != 5 || parts[3] == "" {
		return domain.Location{}, errors.New("invalid topic")
	}

	var loc domain.Location
	if err := json.Unmarshal(payload, &loc); err != nil {
		return loc, err
	}
	if loc.VehicleID != parts[3] {
		return loc, errors.New("vehicle_id does not match topic")
	}
	return loc, loc.Validate()
}
