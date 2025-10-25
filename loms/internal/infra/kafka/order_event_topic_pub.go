package kafka

import (
	"encoding/json"
	"fmt"
	"route256/loms/internal/domain"

	"github.com/IBM/sarama"
)

// OrderEventTopicKafka реализует публикацию событий заказов в Kafka.
type OrderEventTopicKafka struct {
	producer sarama.SyncProducer
	topic    string
}

// NewOrderEventTopicKafka создает новый экземпляр OrderEventTopicKafka.
func NewOrderEventTopicKafka(brokers []string, topic string) (*OrderEventTopicKafka, error) {
	o := &OrderEventTopicKafka{
		topic: topic,
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = false
	config.Producer.Return.Successes = true

	var err error
	o.producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewSyncProducer: %w", err)
	}

	return o, nil
}

// Send публикует сообщение в Kafka с заданным ключом и событием заказа.
func (o *OrderEventTopicKafka) Send(key string, value *domain.OrderEvent) error {
	valueBytes, err := json.Marshal(&OrderEventKafka{
		OrderID: value.OrderID,
		Status:  value.Status,
		Moment:  value.Moment,
	})
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: o.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(valueBytes),
	}

	_, _, err = o.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("producer.SendMessage: %w", err)
	}

	return nil
}
