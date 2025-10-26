package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"route256/cart/pkg/logger"
	"route256/notifier/internal/domain"

	"github.com/IBM/sarama"
)

// OrderEventTopicSubKafka реализует подписку на события заказа через Kafka.
type OrderEventTopicSubKafka struct {
	consumerGroup  sarama.ConsumerGroup
	topics         []string
	eventProcessor orderEventProcessor
}

type orderEventProcessor interface {
	Process(event *domain.OrderEvent)
}

// NewOrderEventTopicSubKafka создает новый экземпляр OrderEventTopicSubKafka.
func NewOrderEventTopicSubKafka(groupID string, topics []string, brokers []string, eventProcessor orderEventProcessor) (*OrderEventTopicSubKafka, error) {
	o := &OrderEventTopicSubKafka{
		topics:         topics,
		eventProcessor: eventProcessor,
	}

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true

	var err error
	o.consumerGroup, err = sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewConsumerGroup: %w", err)
	}

	return o, nil
}

// Start запускает обработку событий заказа.
func (o *OrderEventTopicSubKafka) Start(ctx context.Context) {
	cons := &consumer{
		eventProcessor: o.eventProcessor,
	}

	go func() {
		for err := range o.consumerGroup.Errors() {
			logger.Errorw("ConsumerGroup error", "err", err)
		}
	}()

	go func() {
		defer o.consumerGroup.Close()

		for {
			err := o.consumerGroup.Consume(ctx, o.topics, cons)
			if err != nil {
				logger.Errorw("consumerGroup.Consume error", "err", err)
			}

			if ctx.Err() != nil {
				break
			}
		}
	}()
}

type consumer struct {
	eventProcessor orderEventProcessor
}

func (c *consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	if sess.Context().Err() != nil {
		return nil
	}

	for msg := range claim.Messages() {
		var orderEvent OrderEventKafka
		err := json.Unmarshal(msg.Value, &orderEvent)
		if err != nil {
			logger.Errorw("can't unmarshall kafka message to json", "err", err, "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset)
			continue
		}

		c.eventProcessor.Process(&domain.OrderEvent{
			OrderID: orderEvent.OrderID,
			Status:  orderEvent.Status,
			Moment:  orderEvent.Moment,
		})

		sess.MarkMessage(msg, "")
	}
	return nil
}
