package kafka

import (
	"demoL0/internal/utils"
	"log"
	"os"

	"github.com/IBM/sarama"
)

// InitializeTopics создает или обновляет Kafka топики
func InitializeTopics() {
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaDlqTopic := os.Getenv("KAFKA_DLQ_TOPIC")
	kafkaPartitionCount := int32(utils.AtoiMust(os.Getenv("KAFKA_TOPIC_PARTITION_COUNT")))
	kafkaDlqPartitionCount := int32(utils.AtoiMust(os.Getenv("KAFKA_DLQ_TOPIC_PARTITION_COUNT")))

	createTopicWithPartitions(kafkaTopic, kafkaPartitionCount)
	createTopicWithPartitions(kafkaDlqTopic, kafkaDlqPartitionCount)
	log.Println("kafka topics initialized")
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
				log.Printf("Партиции топика '%s' успешно увеличены до %d\n", kafkaTopic, partitionsCount)
			}
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
			log.Printf("Топик '%s' успешно создан\n", kafkaTopic)
		}
	}
}
