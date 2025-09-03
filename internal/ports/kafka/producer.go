package kafka

import "context"

// Producer определяет интерфейс для отправки сообщений в Kafka
type Producer interface {
	SendMessage(ctx context.Context, topic string, key string, value []byte) error
	Close() error
}
