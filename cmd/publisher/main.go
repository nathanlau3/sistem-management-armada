package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"

	"fleet/internal/config"
	"fleet/internal/core/domain"
)

func main() {
	cfg := config.Load()

	ids := []string{"B1234ABC"}

	opts := paho.NewClientOptions().
		AddBroker(cfg.MQTTBroker).
		SetClientID(fmt.Sprintf("publisher-%d", time.Now().UnixNano())).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(2 * time.Second)

	client := paho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	data, err := os.ReadFile("movement.json")
	if err != nil {
		data, err = os.ReadFile("../../movement.json")
		if err != nil {
			log.Fatal(err)
		}
	}

	var route []domain.Point
	if err := json.Unmarshal(data, &route); err != nil {
		log.Fatal(err)
	}

	var running atomic.Bool

	http.HandleFunc("POST /trigger", func(w http.ResponseWriter, r *http.Request) {
		if !running.CompareAndSwap(false, true) {
			http.Error(w, "publisher is already running", http.StatusConflict)
			return
		}
		go func() {
			defer running.Store(false)
			publish(client, ids, route)
		}()
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintln(w, "started")
	})

	log.Println("publisher listening on", cfg.PublisherAddr)
	log.Fatal(http.ListenAndServe(cfg.PublisherAddr, nil))
}

func publish(client paho.Client, ids []string, route []domain.Point) {
	for i, p := range route {
		if i > 0 {
			time.Sleep(5 * time.Second)
		}
		for _, id := range ids {
			body, _ := json.Marshal(domain.Location{
				VehicleID: id,
				Latitude:  p.Latitude,
				Longitude: p.Longitude,
				Timestamp: time.Now().Unix(),
			})

			topic := "/fleet/vehicle/" + id + "/location"
			client.Publish(topic, 1, false, body).Wait()
			log.Println(topic, string(body))
		}
	}
}
