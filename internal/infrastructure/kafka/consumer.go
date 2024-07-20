package kafka

import (
	"github.com/IBM/sarama"
	"time"
)

type Consumer struct {
	SingleConsumer sarama.Consumer
}

func NewConsumer(brokers Brokers) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second

	consumer, err := sarama.NewConsumer(brokers, config)
	return &Consumer{
		SingleConsumer: consumer,
	}, err
}
