package service

import (
	"github.com/IBM/sarama"
	"log"
	"os"
)

type KafkaProducer struct {
	writer sarama.SyncProducer
}

func NewProducer(writer sarama.SyncProducer) *KafkaProducer {
	return &KafkaProducer{writer: writer}
}

func (pr *KafkaProducer) SendToDlq(consumerMsg *sarama.ConsumerMessage) {
	kafkaDlqTopic := os.Getenv("KAFKA_DLQ_TOPIC")

	msg := &sarama.ProducerMessage{
		Topic: kafkaDlqTopic,
		Value: sarama.ByteEncoder(consumerMsg.Value),
		Key:   sarama.ByteEncoder(consumerMsg.Key),
	}

	_, _, err := pr.writer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
}
