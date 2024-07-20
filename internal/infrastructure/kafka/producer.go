package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"

	"github.com/pkg/errors"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
}

func newAsyncProducer(ctx context.Context, brokers Brokers) (sarama.AsyncProducer, error) {
	asyncProducerConfig := sarama.NewConfig()

	asyncProducerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	asyncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll

	asyncProducerConfig.Producer.Return.Successes = true
	asyncProducerConfig.Producer.Return.Errors = true

	asyncProducer, err := sarama.NewAsyncProducer(brokers, asyncProducerConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error with async kafka-producer")
	}

	go func() {
		successesChan := asyncProducer.Successes()
		errorsChan := asyncProducer.Errors()
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-successesChan:
				if !ok {
					return
				}
			case e, ok := <-errorsChan:
				if !ok {
					return
				}
				if e == nil {
					continue
				}
				fmt.Println("Async error with key", e.Error())
			}
		}
	}()

	return asyncProducer, nil
}

func NewProducer(ctx context.Context, brokers Brokers) (*Producer, error) {
	asyncProducer, err := newAsyncProducer(ctx, brokers)
	if err != nil {
		return nil, errors.Wrap(err, "error with async kafka-producer")
	}

	producer := &Producer{
		asyncProducer: asyncProducer,
	}

	return producer, nil
}

func (k *Producer) SendAsyncMessage(message *sarama.ProducerMessage) {
	k.asyncProducer.Input() <- message
}

func (k *Producer) Close() error {
	err := k.asyncProducer.Close()
	if err != nil {
		return errors.Wrap(err, "kafka.Connector.Close")
	}

	return nil
}
