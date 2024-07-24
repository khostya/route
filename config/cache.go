package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type CacheConfig struct {
	TTL      time.Duration `yaml:"ttl"`
	Capacity uint          `yaml:"capacity"`
}

func NewCacheConfig() (CacheConfig, error) {
	path := os.Getenv("CACHE_CONFIG_PATH")
	if path == "" {
		return CacheConfig{}, ErrCacheConfigPathIsEmpty
	}
	var cfg CacheConfig
	err := cleanenv.ReadConfig(path, &cfg)
	return cfg, err
}

func MustNewCacheConfig() CacheConfig {
	cfg, err := NewCacheConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
