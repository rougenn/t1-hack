package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

// Отправка сообщения в Kafka
func SendToKafka(topic, assistantID string, content []string) error {
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		return err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(fmt.Sprintf("%s: %s", assistantID, content)),
	}

	_, _, err = producer.SendMessage(msg)
	return err
}
