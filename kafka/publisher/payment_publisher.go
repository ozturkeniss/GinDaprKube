package publisher

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"daprps/api/proto/events"

	"github.com/Shopify/sarama"
)

type PaymentPublisher struct {
	producer sarama.SyncProducer
}

func NewPaymentPublisher() (*PaymentPublisher, error) {
	// Get Kafka brokers from environment variable
	brokersStr := getEnv("KAFKA_BROKERS", "localhost:9092")
	brokers := strings.Split(brokersStr, ",")

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &PaymentPublisher{
		producer: producer,
	}, nil
}

func (p *PaymentPublisher) PublishPaymentCompleted(ctx context.Context, event *events.PaymentCompletedEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "payment-completed",
		Key:   sarama.StringEncoder(event.PaymentId),
		Value: sarama.ByteEncoder(eventBytes),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Payment completed event published to partition %d at offset %d", partition, offset)
	return nil
}

func (p *PaymentPublisher) Close() error {
	return p.producer.Close()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
