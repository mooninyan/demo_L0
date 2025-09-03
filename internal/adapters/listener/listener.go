package listener

import (
	"context"
	kafkaAdapter "demoL0/internal/adapters/kafka"
	kafkaPort "demoL0/internal/ports/kafka"
	"demoL0/internal/usecase/order"
	"demoL0/internal/utils"
	"log"
	"sync"
)

// Listener управляет Kafka consumers
type Listener struct {
	ctx          context.Context
	mainConsumer kafkaPort.Consumer
	dlqConsumer  kafkaPort.Consumer
	mainHandler  *kafkaAdapter.ConsumerGroupHandler
	dlqHandler   *kafkaAdapter.ConsumerGroupHandler
}

// NewListener создает новый экземпляр listener
func NewListener(ctx context.Context, producer kafkaPort.Producer, orderService order.Service, validator *kafkaAdapter.Validator) *Listener {
	// Создаем handlers
	mainHandler := kafkaAdapter.NewConsumerGroupHandler(false, producer, orderService, validator)
	dlqHandler := kafkaAdapter.NewConsumerGroupHandler(true, producer, orderService, validator)

	// Создаем consumers
	mainConsumer, err := kafkaAdapter.NewConsumer(mainHandler)
	if err != nil {
		log.Fatalf("Failed to create main consumer: %v", err)
	}

	var dlqConsumer kafkaPort.Consumer
	if utils.GetDlqFlagOrDefault() {
		dlqConsumer, err = kafkaAdapter.NewDlqConsumer(dlqHandler)
		if err != nil {
			log.Fatalf("Failed to create DLQ consumer: %v", err)
		}
	}

	return &Listener{
		ctx:          ctx,
		mainConsumer: mainConsumer,
		dlqConsumer:  dlqConsumer,
		mainHandler:  mainHandler,
		dlqHandler:   dlqHandler,
	}
}

// Start запускает все consumers
func (l *Listener) Start(wg *sync.WaitGroup) {
	// Запускаем основной consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := l.mainConsumer.Start(l.ctx); err != nil {
			log.Printf("Main consumer error: %v", err)
		}
	}()

	// Запускаем DLQ consumer если включен
	if l.dlqConsumer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := l.dlqConsumer.Start(l.ctx); err != nil {
				log.Printf("DLQ consumer error: %v", err)
			}
		}()
	}
}

// Stop останавливает все consumers
func (l *Listener) Stop() {
	if err := l.mainConsumer.Stop(); err != nil {
		log.Printf("Error stopping main consumer: %v", err)
	}

	if l.dlqConsumer != nil {
		if err := l.dlqConsumer.Stop(); err != nil {
			log.Printf("Error stopping DLQ consumer: %v", err)
		}
	}
}
