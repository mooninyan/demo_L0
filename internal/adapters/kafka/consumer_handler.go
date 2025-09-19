package kafka

import (
	"context"
	"demoL0/internal/ports/kafka"
	"demoL0/internal/usecase/order"
	"demoL0/internal/utils"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

// ConsumerGroupHandler обрабатывает сообщения из Kafka
type ConsumerGroupHandler struct {
	isDlq        bool
	producer     kafka.Producer
	orderService order.Service
	validator    *Validator
}

// NewConsumerGroupHandler создает новый обработчик consumer group
func NewConsumerGroupHandler(isDlq bool, producer kafka.Producer, orderService order.Service, validator *Validator) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		isDlq:        isDlq,
		producer:     producer,
		orderService: orderService,
		validator:    validator,
	}
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			h.handleMessage(session, message)
		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *ConsumerGroupHandler) handleMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	session.MarkMessage(msg, "")
	log.Printf("Partition: %d, Offset: %d, Key: %s, Value: %s\n",
		msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))

	var orderDto order.MainOrder
	err := json.Unmarshal(msg.Value, &orderDto)
	if err != nil {
		log.Printf("[isDlq?:%v] serialization error: %v\n", h.isDlq, err)
		if h.isDlq {
			return
		}
		h.sendToDlq(msg)
		return
	}

	if err = h.validator.Validate(orderDto); err != nil {
		log.Printf("[isDlq?:%v] validation error: %v\n", h.isDlq, err)
		if !h.isDlq {
			h.sendToDlq(msg)
		}
		return
	}

	if err = h.createOrUpdateOrder(orderDto); err != nil && !h.isDlq {
		h.sendToDlq(msg)
	}
	log.Printf("[isDlq?:%v] end of message handling", h.isDlq)
}

func (h *ConsumerGroupHandler) createOrUpdateOrder(orderData order.MainOrder) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = h.orderService.Create(timeout, &orderData)
	if err != nil {
		condition := strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") ||
			strings.Contains(err.Error(), "SQLSTATE 23505")
		if condition {
			err = h.orderService.Update(context.TODO(), orderData.OrderUID, &orderData)
			if err != nil {
				return utils.WrapAndLog(err)
			}
			return nil
		}
		return utils.WrapAndLog(err)
	}
	return nil
}

func (h *ConsumerGroupHandler) sendToDlq(msg *sarama.ConsumerMessage) {
	dlqTopic := os.Getenv("KAFKA_DLQ_TOPIC")
	if dlqTopic == "" {
		log.Println("DLQ topic not configured")
		return
	}

	err := h.producer.SendMessage(context.Background(), dlqTopic, string(msg.Key), msg.Value)
	if err != nil {
		log.Printf("Failed to send message to DLQ: %v", err)
	}
}
