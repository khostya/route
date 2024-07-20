package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type IMDBConfig struct {
	TTL      time.Duration `yaml:"ttl"`
	Capacity uint          `yaml:"capacity"`
}

func NewIMDBConfig() (IMDBConfig, error) {
	path := os.Getenv("IMDB_CONFIG_PATH")
	if path == "" {
		return IMDBConfig{}, ErrIMDBConfigPathIsEmpty
	}
	var cfg IMDBConfig
	err := cleanenv.ReadConfig(path, &cfg)
	return cfg, err
}

func MustNewIMDBConfig() IMDBConfig {
	cfg, err := NewIMDBConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
