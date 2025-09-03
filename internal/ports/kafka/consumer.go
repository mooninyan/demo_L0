package kafka

import "context"

// Consumer определяет интерфейс для получения сообщений из Kafka
type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
}
