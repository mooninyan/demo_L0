package kafka

import (
	"context"
	"demoL0/internal/ports/kafka"
	"log"
	"os"

	"github.com/IBM/sarama"
)

// consumer реализует интерфейс kafka.Consumer
type consumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       *ConsumerGroupHandler
}

// NewConsumer создает новый экземпляр Kafka consumer
func NewConsumer(handler *ConsumerGroupHandler) (kafka.Consumer, error) {
	kafkaHostPort := os.Getenv("KAFKA_HOST_PORT")
	kafkaGroupId := os.Getenv("KAFKA_GROUP_ID")

	brokers := []string{kafkaHostPort}
	groupID := kafkaGroupId

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &consumer{
		consumerGroup: consumerGroup,
		handler:       handler,
	}, nil
}

// NewDlqConsumer создает новый экземпляр DLQ consumer
func NewDlqConsumer(handler *ConsumerGroupHandler) (kafka.Consumer, error) {
	kafkaHostPort := os.Getenv("KAFKA_HOST_PORT")
	kafkaGroupId := os.Getenv("KAFKA_DLQ_GROUP_ID")

	brokers := []string{kafkaHostPort}
	groupID := kafkaGroupId

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &consumer{
		consumerGroup: consumerGroup,
		handler:       handler,
	}, nil
}

func (c *consumer) Start(ctx context.Context) error {
	topic := os.Getenv("KAFKA_TOPIC")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := c.consumerGroup.Consume(ctx, []string{topic}, c.handler); err != nil {
				log.Printf("Error from consumer: %v", err)
				return err
			}

			if ctx.Err() != nil {
				return ctx.Err()
			}
		}
	}
}

func (c *consumer) Stop() error {
	return c.consumerGroup.Close()
}
