package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"

	"github.com/pkg/errors"
)

type Producer struct {
	brokers       Brokers
	syncProducer  sarama.SyncProducer
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

func newSyncProducer(brokers Brokers) (sarama.SyncProducer, error) {
	syncProducerConfig := sarama.NewConfig()

	syncProducerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	syncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll

	syncProducerConfig.Producer.Idempotent = true
	syncProducerConfig.Net.MaxOpenRequests = 1

	syncProducerConfig.Producer.CompressionLevel = sarama.CompressionLevelDefault
	syncProducerConfig.Producer.Compression = sarama.CompressionGZIP

	syncProducerConfig.Producer.Return.Successes = true
	syncProducerConfig.Producer.Return.Errors = true

	syncProducer, err := sarama.NewSyncProducer(brokers, syncProducerConfig)

	if err != nil {
		return nil, errors.Wrap(err, "error with sync kafka-producer")
	}

	return syncProducer, nil
}

func NewProducer(ctx context.Context, brokers Brokers) (*Producer, error) {
	syncProducer, err := newSyncProducer(brokers)
	if err != nil {
		return nil, errors.Wrap(err, "error with sync kafka-producer")
	}

	asyncProducer, err := newAsyncProducer(ctx, brokers)
	if err != nil {
		return nil, errors.Wrap(err, "error with async kafka-producer")
	}

	producer := &Producer{
		brokers:       brokers,
		syncProducer:  syncProducer,
		asyncProducer: asyncProducer,
	}

	return producer, nil
}

func (k *Producer) SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return k.syncProducer.SendMessage(message)
}

func (k *Producer) SendSyncMessages(messages []*sarama.ProducerMessage) error {
	return k.syncProducer.SendMessages(messages)
}

func (k *Producer) SendAsyncMessage(message *sarama.ProducerMessage) {
	k.asyncProducer.Input() <- message
}

func (k *Producer) Close() error {
	err := k.syncProducer.Close()
	if err != nil {
		return errors.Wrap(err, "kafka.Connector.Close")
	}

	err = k.asyncProducer.Close()
	if err != nil {
		return errors.Wrap(err, "kafka.Connector.Close")
	}

	return nil
}
