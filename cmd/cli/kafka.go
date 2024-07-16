package main

import (
	"context"
	"homework/config"
	"homework/internal/infrastructure/app/oncall"
	"homework/internal/infrastructure/kafka"
	"log"
)

func getOnCallKafkaSender(ctx context.Context, cfg config.KafkaConfig) *oncall.KafkaProducer {
	kafkaProducer, err := kafka.NewProducer(ctx, cfg.Brokers)
	if err != nil {
		log.Fatalln(err)
	}

	onCallConsumer := oncall.NewKafkaProducer(kafkaProducer, kafka.Topic(cfg.OnCallTopic))
	return onCallConsumer
}

func getOnCallKafkaReceiver(cfg config.KafkaConfig, handler oncall.HandleFunc) *oncall.KafkaConsumer {
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
