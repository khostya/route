package oncall

import (
	"github.com/IBM/sarama"
	"homework/internal/infrastructure/kafka"
)

type HandleFunc func(message *sarama.ConsumerMessage)

type KafkaReceiver struct {
	consumer *kafka.Consumer
}

func NewKafkaReceiver(consumer *kafka.Consumer) *KafkaReceiver {
	return &KafkaReceiver{
		consumer: consumer,
	}
}

func (r *KafkaReceiver) Subscribe(topic kafka.Topic, handler HandleFunc) error {
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

		go func(pc sarama.PartitionConsumer, partition int32) {
			for message := range pc.Messages() {
				handler(message)
			}
		}(pc, partition)
	}

	return nil
}

func (r *KafkaReceiver) Close() error {
	return r.consumer.SingleConsumer.Close()
}
