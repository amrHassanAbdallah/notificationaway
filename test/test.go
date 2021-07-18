package main

import (
	"context"
	"encoding/json"
	"fmt"
	client2 "github.com/amrHassanAbdallah/notificationaway/client"
	"github.com/amrHassanAbdallah/notificationaway/logger"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"net/http"
	"os"
	"time"
)

var expectedMessage string

func handler(w http.ResponseWriter, r *http.Request) {
	var payload map[string]string

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		logger.Fatalf("failed to decode payload, error: %v", err)
	}

	if payload["message"] != expectedMessage {
		logger.Fatalw(fmt.Sprintf("invalid message check, expected (%v) got (%)", expectedMessage, payload["message"]))
	} else {
		logger.Infow("received a message and it match the expected")
	}
	os.Exit(0)

}
func main() {
	//check the svc some times then fail....
	for i := 0; i < 5; i++ {
		logger.Infow("checking the svc")
		resp, err := http.Get("http://localhost:7981/health")
		if err != nil || resp != nil && resp.StatusCode != 200 {
			code := 0
			if resp != nil {
				code = resp.StatusCode
			}
			logger.Infow("failed to connect to the svc", "status", code, "err", err)
			time.Sleep(1 * time.Second)
			continue
		}
		logger.Infow("service is up and running on http://localhost:7981")
		break
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := client2.NewClient("http://localhost:7981/api/v1")
	if err != nil {
		logger.Fatalw("failed to initiate the client")
	}
	expectedMessage = `Hi Hamda
								Welcome to the jungle,
								Enjoy,
											`
	resp, err := client.AddMessage(ctx, client2.AddMessageJSONRequestBody{
		Language:     "en",
		ProviderType: "webhook",
		Template:     expectedMessage,
		Type:         "greetings",
	})
	if err != nil || resp != nil && (resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusConflict) {
		code := 0
		if resp != nil {
			code = resp.StatusCode
		}
		logger.Fatalw("failed to add a new message", "err", err, "code", code)
	}
	logger.Infow("created a message")

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}
	defer producer.Close()
	ob := map[string]string{
		"lang": "en",
		"to":   "http://localhost:5000/webhook",
	}
	data, err := json.Marshal(ob)
	if err != nil {
		logger.Fatalf("failed to marshal event data, error: %v", err)
	}
	eventType := "notifications"
	stype := "greetings"
	msg := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &eventType,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(data),
		Headers: []kafka.Header{
			{Key: "type", Value: []byte(stype)},
		},
	}
	err = producer.Produce(&msg, nil)
	if err != nil {
		logger.Fatalf("failed to push event, error:%v", err)
	}
	logger.Infow("sent the notification event.")

	http.HandleFunc("/webhook", handler)
	logger.Fatal(http.ListenAndServe(":5000", nil))
}
