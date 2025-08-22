package listener

import (
	"context"
	"github.com/IBM/sarama"
	"log"
	"os"
	"sync"
)

type Listener struct {
	ctx             context.Context
	groupHandlerMap map[string]*GroupAndHandler
}

type GroupAndHandler struct {
	Group   sarama.ConsumerGroup
	Handler *ConsumerGroupHandler
}

func NewListener(
	ctx context.Context,
	groupHandlerMap map[string]*GroupAndHandler,
) *Listener {
	return &Listener{
		ctx:             ctx,
		groupHandlerMap: groupHandlerMap,
	}
}

func (lr *Listener) Listen(wg *sync.WaitGroup) {
	topic := os.Getenv("KAFKA_TOPIC")
	lr.listenInternal(lr.ctx, lr.groupHandlerMap["main_group"], topic)
	wg.Done()
}

func (lr *Listener) ListenDlq(wg *sync.WaitGroup) {
	topic := os.Getenv("KAFKA_DLQ_TOPIC")
	lr.listenInternal(lr.ctx, lr.groupHandlerMap["dlq_group"], topic)
	wg.Done()
}

func (lr *Listener) listenInternal(ctx context.Context, gh *GroupAndHandler, topic string) {
	defer func() {
		if err := gh.Group.Close(); err != nil {
			log.Printf("Error closing consumer group: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := gh.Group.Consume(ctx, []string{topic}, gh.Handler); err != nil {
				log.Printf("Error from consumer: %v", err)
				return
			}
			// Проверяем, что не завершили работу
			if ctx.Err() != nil {
				return
			}
		}
	}
}
