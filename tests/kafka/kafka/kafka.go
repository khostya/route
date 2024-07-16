package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"homework/internal/infrastructure/app/oncall"
	"homework/internal/infrastructure/kafka"
	"os"
)

type Kafka struct {
	OnCallSender   *oncall.KafkaProducer
	OnCallReceiver *oncall.KafkaConsumer
	Topic          kafka.Topic
	broker         string
}

func NewOnCallFromEnv(ctx context.Context, topic kafka.Topic) *Kafka {
	broker := os.Getenv("TEST_KAFKA_BROKER")
	if broker == "" {
		panic("TEST_KAFKA_BROKER isn`t set")
	}
	consumer, err := kafka.NewConsumer([]string{broker})
	if err != nil {
		panic(err)
	}

	producer, err := kafka.NewProducer(ctx, []string{broker})
	if err != nil {
		panic(err)
	}
	onCallSender := oncall.NewKafkaProducer(producer, topic)
	onCallReceiver := oncall.NewKafkaReceiver(consumer)

	return &Kafka{
		onCallSender,
		onCallReceiver,
		topic,
		broker,
	}
}

func (k *Kafka) SetUp() {
	broker := k.newBroker(k.broker)
	defer broker.Close()

	_, err := broker.CreateTopics(&sarama.CreateTopicsRequest{
		TopicDetails: map[string]*sarama.TopicDetail{
			string(k.Topic): {
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		},
	})
	if err != nil {
		panic(err)
	}
}

func (k *Kafka) TearDown() {
	broker := k.newBroker(k.broker)
	defer broker.Close()

	_, err := broker.DeleteTopics(&sarama.DeleteTopicsRequest{
		Topics: []string{string(k.Topic)},
	})
	if err != nil {
		panic(err)
	}
}

func (k *Kafka) newBroker(brokerUrl string) *sarama.Broker {
	cfg := sarama.NewConfig()
	broker := sarama.NewBroker(brokerUrl)
	err := broker.Open(cfg)
	if err != nil {
		panic(err)
	}
	_, err = broker.Connected()
	if err != nil {
		panic(err)
	}
	return broker
}

func (k *Kafka) Close() {
	err := k.OnCallSender.Close()
	if err != nil {
		panic(err)
	}

	err = k.OnCallReceiver.Close()
	if err != nil {
		panic(err)
	}
}
