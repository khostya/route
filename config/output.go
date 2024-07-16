package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"homework/pkg/output"
	"os"
	"slices"
)

type (
	OutputConfig struct {
		Filter string `yaml:"filter"`
	}
)

func NewOutputConfig() (OutputConfig, error) {
	path := os.Getenv("OUTPUT_CONFIG_PATH")
	if path == "" {
		return OutputConfig{}, ErrOutputConfigPathIsEmpty
	}
	var cfg OutputConfig
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return cfg, err
	}
	if !slices.Contains([]string{output.CLI, output.Kafka}, cfg.Filter) {
		return cfg, ErrOutputFilterDoesNotExits
	}
	return cfg, err
}

func MustNewOutputConfig() OutputConfig {
	config, err := NewOutputConfig()
	if err != nil {
		panic(err)
	}
	return config
}
