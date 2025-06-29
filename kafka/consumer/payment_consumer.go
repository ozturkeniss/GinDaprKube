package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"daprps/api/proto/events"

	"github.com/Shopify/sarama"
)

type PaymentConsumer struct {
	consumer sarama.ConsumerGroup
	handler  PaymentEventHandler
}

type PaymentEventHandler interface {
	HandlePaymentCompleted(ctx context.Context, event *events.PaymentCompletedEvent) error
}

func NewPaymentConsumer(handler PaymentEventHandler) (*PaymentConsumer, error) {
	// Get Kafka brokers from environment variable
	brokersStr := getEnv("KAFKA_BROKERS", "localhost:9092")
	brokers := strings.Split(brokersStr, ",")

	groupID := getEnv("KAFKA_CONSUMER_GROUP", "dapr-consumer-group")

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &PaymentConsumer{
		consumer: consumer,
		handler:  handler,
	}, nil
}

func (c *PaymentConsumer) Start(ctx context.Context, topics []string) error {
	for {
		err := c.consumer.Consume(ctx, topics, c)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *PaymentConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var event events.PaymentCompletedEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		if err := c.handler.HandlePaymentCompleted(session.Context(), &event); err != nil {
			log.Printf("Error handling payment completed event: %v", err)
		}

		session.MarkMessage(message, "")
	}
	return nil
}

func (c *PaymentConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *PaymentConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *PaymentConsumer) Close() error {
	return c.consumer.Close()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
