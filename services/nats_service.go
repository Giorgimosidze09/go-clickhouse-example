package services

import (
	"encoding/json"
	"log"

	"go-clickhouse-example/models"

	"github.com/nats-io/nats.go"
)

type NATSService struct {
	js          nats.JetStreamContext
	streamName  string
	subjectName string
}

func NewNATSService(natsURL, streamName, subjectName string) *NATSService {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Failed to create JetStream context: %v", err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subjectName},
		Storage:  nats.FileStorage,
	})
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	return &NATSService{
		js:          js,
		streamName:  streamName,
		subjectName: subjectName,
	}
}

func (n *NATSService) PublishItem(item models.ItemResponse) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	_, err = n.js.Publish(n.subjectName, data)
	return err
}

func (n *NATSService) SubscribeItem(handler func(models.ItemResponse)) error {
	_, err := n.js.Subscribe(n.subjectName, func(msg *nats.Msg) {
		var item models.ItemResponse
		if err := json.Unmarshal(msg.Data, &item); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return
		}
		handler(item)
	})
	return err
}
