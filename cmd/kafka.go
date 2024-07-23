package cmd

import (
	"context"
	"homework/config"
	"homework/internal/infrastructure/app/oncall"
	"homework/internal/infrastructure/kafka"
	"log"
)

func GetOnCallKafkaSender(ctx context.Context) *oncall.KafkaProducer {
	cfg := config.MustNewKafkaConfig()

	kafkaProducer, err := kafka.NewProducer(ctx, cfg.Brokers)
	if err != nil {
		log.Fatalln(err)
	}

	onCallConsumer := oncall.NewKafkaProducer(kafkaProducer, kafka.Topic(cfg.OnCallTopic))
	return onCallConsumer
}

func GetOnCallKafkaReceiver(handler oncall.HandleFunc) *oncall.KafkaConsumer {
	cfg := config.MustNewKafkaConfig()

	kafkaConsumer, err := kafka.NewConsumer(cfg.Brokers)
	if err != nil {
		log.Fatalln(err)
	}

	onCallConsumer := oncall.NewKafkaReceiver(kafkaConsumer)
	err = onCallConsumer.Subscribe(kafka.Topic(cfg.OnCallTopic), handler)
	if err != nil {
		_ = onCallConsumer.Close()
		log.Fatalln(err)
	}

	return onCallConsumer
}
