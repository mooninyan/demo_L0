package listener

import (
	"context"
	"demoL0/internal/dto"
	"demoL0/internal/service"
	"demoL0/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"strings"
	"time"
)

type ConsumerGroupHandler struct {
	isDlq       bool
	dlqProducer *service.KafkaProducer
	service     *service.OrderService
}

func NewConsumerGroupHandler(isDlq bool,
	dlqProducer *service.KafkaProducer,
	service *service.OrderService) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		isDlq:       isDlq,
		dlqProducer: dlqProducer,
		service:     service,
	}
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.handleMessage(session, message)
	}
	return nil
}

func (h *ConsumerGroupHandler) handleMessage(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	session.MarkMessage(msg, "")
	fmt.Printf("Partition: %d, Offset: %d, Key: %s, Value: %s\n",
		msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	var order dto.MainOrder
	err := json.Unmarshal(msg.Value, &order)
	if err != nil {
		log.Printf("[isDlq?:%v] serrialization error: %v\n", h.isDlq, err)
		if h.isDlq {
			return
		}
		h.dlqProducer.SendToDlq(msg)
		return
	}
	if err = h.createOrUpdateOrder(order); err != nil && !h.isDlq {
		h.dlqProducer.SendToDlq(msg)
	}
	log.Printf("[isDlq?:%v] end of message handling", h.isDlq)
}

func (h *ConsumerGroupHandler) createOrUpdateOrder(order dto.MainOrder) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = h.service.Create(timeout, &order)
	if err != nil {
		condition := strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") ||
			strings.Contains(err.Error(), "SQLSTATE 23505")
		if condition {
			err = h.service.Update(context.TODO(), order.OrderUID, &order)
			if err != nil {
				return utils.WrapAndLog(err)
			}
			return nil
		} else {
			return utils.WrapAndLog(err)
		}
	}
	return nil
}
