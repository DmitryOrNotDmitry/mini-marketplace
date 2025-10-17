package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"route256/cart/pkg/logger"
	"route256/notifier/internal/domain"
	"time"

	"github.com/IBM/sarama"
)

// OrderEventTopicSubKafka реализует подписку на события заказа через Kafka.
type OrderEventTopicSubKafka struct {
	consumerGroup sarama.ConsumerGroup
	topics        []string
}

// NewOrderEventTopicSubKafka создает новый экземпляр OrderEventTopicSubKafka.
func NewOrderEventTopicSubKafka(groupID string, topics []string, brokers []string) (*OrderEventTopicSubKafka, error) {
	o := &OrderEventTopicSubKafka{
		topics: topics,
	}

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	var err error
	o.consumerGroup, err = sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewConsumerGroup: %w", err)
	}

	return o, nil
}

// Start запускает обработку событий заказа и возвращает канал для получения событий.
func (o *OrderEventTopicSubKafka) Start(ctx context.Context) <-chan *domain.OrderEvent {
	msgOut := make(chan *domain.OrderEvent)

	cons := &consumer{
		msgOut: msgOut,
	}

	go func() {
		for err := range o.consumerGroup.Errors() {
			logger.Warnw("ConsumerGroup error", "err", err)
		}
	}()

	go func() {
		defer o.consumerGroup.Close()
		defer close(msgOut)

		for {
			err := o.consumerGroup.Consume(ctx, o.topics, cons)
			if err != nil {
				logger.Warnw("consumerGroup.Consume error", "err", err)
				time.Sleep(1 * time.Second)
			}

			if ctx.Err() != nil {
				break
			}
		}
	}()

	return msgOut
}

type consumer struct {
	msgOut chan<- *domain.OrderEvent
}

func (c *consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var orderEvent domain.OrderEvent
		err := json.Unmarshal(msg.Value, &orderEvent)
		if err != nil {
			logger.Warnw("can't unmarshall kafka message to json", "err", err, "topic", msg.Topic, "partition", msg.Partition, "offset", msg.Offset)
			continue
		}

		c.msgOut <- &orderEvent

		sess.MarkMessage(msg, "")
		sess.Commit()
	}
	return nil
}
