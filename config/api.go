package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type ApiConfig struct {
	GrpcPort     uint   `yaml:"grpc_port"`
	GrpcENDPOINT string `yaml:"grpc_endpoint"`
	HttpPort     uint   `yaml:"http_port"`
	HttpENDPOINT string `yaml:"http_endpoint"`
	SwaggerPort  uint   `yaml:"swagger_port"`
}

func NewApiConfig() (ApiConfig, error) {
	path := os.Getenv("API_CONFIG_PATH")
	if path == "" {
		return ApiConfig{}, ErrApiConfigPathIsEmpty
	}
	var cfg ApiConfig
	err := cleanenv.ReadConfig(path, &cfg)
	return cfg, err
}

func MustNewApiConfig() ApiConfig {
	cfg, err := NewApiConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
