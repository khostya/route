package oncall

import (
	"github.com/IBM/sarama"
	"homework/internal/infrastructure/kafka"
	"sync"
)

type HandleFunc func(message *sarama.ConsumerMessage)

type KafkaConsumer struct {
	consumer *kafka.Consumer

	closeWG     sync.WaitGroup
	closeNotify chan struct{}
}

func NewKafkaReceiver(consumer *kafka.Consumer) *KafkaConsumer {
	return &KafkaConsumer{
		consumer:    consumer,
		closeNotify: make(chan struct{}),
	}
}

func (r *KafkaConsumer) Subscribe(topic kafka.Topic, handler HandleFunc) error {
	partitionList, err := r.consumer.SingleConsumer.Partitions(string(topic))
	if err != nil {
		return err
	}

	initialOffset := sarama.OffsetNewest
	for _, partition := range partitionList {
		pc, err := r.consumer.SingleConsumer.ConsumePartition(string(topic), partition, initialOffset)
		if err != nil {
			return err
		}

		r.closeWG.Add(1)
		go func() {
			<-r.closeNotify
			pc.Close()
			r.closeWG.Done()
		}()

		go func(pc sarama.PartitionConsumer, partition int32) {
			for message := range pc.Messages() {
				handler(message)
			}
		}(pc, partition)
	}

	return nil
}

func (r *KafkaConsumer) Close() error {
	close(r.closeNotify)
	r.closeWG.Wait()
	return r.consumer.SingleConsumer.Close()
}
