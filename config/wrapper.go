package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"math"
	"os"
	"strconv"
	"strings"
)

type (
	Wrapper struct {
		Type           string
		CapacityInGram float64
		PriceInRub     float64
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
	return cfg, cleanenv.ReadConfig(path, &cfg)
}

func (w *Wrapper) UnmarshalText(text []byte) error {
	tokens := strings.Split(string(text), "|")
	if len(tokens) != 3 {
		return ErrTokensLengthIsNotValid
	}
	for _, v := range tokens {
		if v == "" {
			return ErrValueIsEmpty
		}
	}
	priceInRub, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		return err
	}
	if priceInRub < 0 {
		return ErrPriceInRubIsNotValid
	}

	capacityInGram, err := strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		return err
	}
	if capacityInGram < -1 {
		return ErrCapacityInGramInNotValid
	}

	w.Type = tokens[0]
	w.PriceInRub = priceInRub

	if capacityInGram == -1 {
		capacityInGram = math.Inf(1)
	}
	w.CapacityInGram = capacityInGram
	return nil
}
