package main

import (
	"context"
	"homework/config"
	"homework/internal/infrastructure/app/oncall"
	"homework/internal/infrastructure/kafka"
)

func getOnCallKafkaSender(ctx context.Context, cfg config.KafkaConfig) (*oncall.KafkaSender, error) {
	kafkaProducer, err := kafka.NewProducer(ctx, cfg.Brokers)
	if err != nil {
		return nil, err
	}

	onCallSender := oncall.NewKafkaSender(kafkaProducer, kafka.Topic(cfg.OnCallTopic))
	return onCallSender, nil
}

func getOnCallKafkaReceiver(cfg config.KafkaConfig, handler oncall.HandleFunc) (*oncall.KafkaReceiver, error) {
	kafkaConsumer, err := kafka.NewConsumer(cfg.Brokers)
	if err != nil {
		return nil, err
	}

	onCallReceiver := oncall.NewKafkaReceiver(kafkaConsumer)
	err = onCallReceiver.Subscribe(kafka.Topic(cfg.OnCallTopic), handler)
	if err != nil {
		_ = onCallReceiver.Close()
		return nil, err
	}

	return onCallReceiver, nil
}
