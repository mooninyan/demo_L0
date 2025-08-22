package kfk

import (
	"demoL0/internal/utils"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
)

func InitializeTopics() {
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaDlqTopic := os.Getenv("KAFKA_DLQ_TOPIC")
	kafkaPartitionCount := int32(utils.AtoiMust(os.Getenv("KAFKA_TOPIC_PARTITION_COUNT")))
	kafkaDlqPartitionCount := int32(utils.AtoiMust(os.Getenv("KAFKA_DLQ_TOPIC_PARTITION_COUNT")))

	createTopicWithPartitions(kafkaTopic, kafkaPartitionCount)
	createTopicWithPartitions(kafkaDlqTopic, kafkaDlqPartitionCount)
}

func createTopicWithPartitions(kafkaTopic string, partitionsCount int32) {
	kafkaHostPort := os.Getenv("KAFKA_HOST_PORT")

	config := sarama.NewConfig()

	brokers := []string{kafkaHostPort}

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatalf("Error creating cluster admin: %v", err)
	}
	defer func() {
		if err := admin.Close(); err != nil {
			log.Printf("Error closing admin: %v", err)
		}
	}()

	topics, err := admin.ListTopics()
	if err != nil {
		log.Fatalf("Error listing topics: %v", err)
	}

	if _, exists := topics[kafkaTopic]; exists {
		currentPartitions := topics[kafkaTopic].NumPartitions

		if partitionsCount > currentPartitions {
			err = admin.CreatePartitions(kafkaTopic, partitionsCount, nil, false)
			if err != nil {
				log.Printf("Error increasing partitions: %v", err)
			} else {
				fmt.Printf("Партиции топика '%s' успешно увеличены до %d\n", kafkaTopic, partitionsCount)
			}
		} else {
			topicDetail := &sarama.TopicDetail{
				NumPartitions:     partitionsCount,
				ReplicationFactor: 1,
			}

			err = admin.CreateTopic(kafkaTopic, topicDetail, false)
			if err != nil {
				log.Printf("Error creating topic: %v", err)
			} else {
				fmt.Printf("Топик '%s' успешно создан\n", kafkaTopic)
			}
		}
	}
}

func CreateConsumerGroup() sarama.ConsumerGroup {
	kafkaHostPort := os.Getenv("KAFKA_HOST_PORT")
	kafkaGroupId := os.Getenv("KAFKA_GROUP_ID")

	brokers := []string{kafkaHostPort}
	groupID := kafkaGroupId

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	return consumerGroup
}

func CreateDlqConsumerGroup() sarama.ConsumerGroup {
	kafkaHostPort := os.Getenv("KAFKA_HOST_PORT")
	kafkaGroupId := os.Getenv("KAFKA_DLQ_GROUP_ID")

	brokers := []string{kafkaHostPort}
	groupID := kafkaGroupId

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	return consumerGroup
}

func CreateKafkaWriter() sarama.SyncProducer {
	brokers := []string{os.Getenv("KAFKA_HOST_PORT")}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	return producer
}
