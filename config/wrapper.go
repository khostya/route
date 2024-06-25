package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"math"
	"os"
)

type (
	Wrapper struct {
		Type           string  `json:"type"`
		CapacityInGram float64 `json:"capacity_in_gram"`
		PriceInRub     float64 `json:"price_in_rub"`
	}

	WrappersConfig struct {
		Wrappers []Wrapper `yaml:"wrappers"`
	}
)

func NewWrappersConfig() (WrappersConfig, error) {
	path := os.Getenv("WRAPPERS_CONFIG_PATH")
	if path == "" {
		return WrappersConfig{}, ErrWrappersConfigPathIsEmpty
	}
	var cfg WrappersConfig
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return WrappersConfig{}, err
	}

	for i, w := range cfg.Wrappers {
		if w.PriceInRub < 0 {
			return WrappersConfig{}, ErrPriceInRubIsNotValid
		}
		if w.CapacityInGram < -1 {
			return WrappersConfig{}, ErrCapacityInGramInNotValid
		}

		if w.CapacityInGram == -1 {
			cfg.Wrappers[i].CapacityInGram = math.Inf(1)
		}
	}

	return cfg, nil
}
