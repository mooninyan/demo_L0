package kafka

import (
	"context"
	"demoL0/internal/ports/kafka"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

// producer реализует интерфейс kafka.Producer
type producer struct {
	producer sarama.SyncProducer
}

// NewProducer создает новый экземпляр Kafka producer
func NewProducer() (kafka.Producer, error) {
	brokers := []string{os.Getenv("KAFKA_HOST_PORT")}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	saramaProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("error creating producer: %v", err)
	}

	return &producer{producer: saramaProducer}, nil
}

func (p *producer) SendMessage(_ context.Context, topic string, key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Message sent to topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

func (p *producer) Close() error {
	return p.producer.Close()
}
