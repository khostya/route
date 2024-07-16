package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type (
	KafkaConfig struct {
		Brokers     []string `yaml:"brokers"`
		OnCallTopic string   `yaml:"on_call_topic"`
	}
)

func NewKafkaConfig() (KafkaConfig, error) {
	path := os.Getenv("KAFKA_CONFIG_PATH")
	if path == "" {
		return KafkaConfig{}, ErrKafkaConfigPathIsEmpty
	}
	var cfg KafkaConfig
	err := cleanenv.ReadConfig(path, &cfg)
	return cfg, err
}

func MustNewKafkaConfig() KafkaConfig {
	config, err := NewKafkaConfig()
	if err != nil {
		panic(err)
	}
	return config
}
