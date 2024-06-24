package model

import (
	"github.com/shopspring/decimal"
	"math"
)

var (
	stretchCapacityInGram = CapacityInGram(math.Inf(1))
	stretchPriceInRub     = PriceInRub(decimal.NewFromInt(0))
)

type Stretch struct {
	capacityInGram CapacityInGram
	priceInRub     PriceInRub
}

func NewDefaultStretch() Stretch {
	return NewStretch(stretchCapacityInGram, stretchPriceInRub)
}

func NewStretch(capacity CapacityInGram, price PriceInRub) Stretch {
	return Stretch{
		capacityInGram: capacity,
		priceInRub:     price,
	}
}

func (s Stretch) GetCapacityInGram() CapacityInGram {
	return s.capacityInGram
}

func (s Stretch) GetPriceInRub() PriceInRub {
	return s.priceInRub
}

func (s Stretch) GetType() WrapperType {
	return stretchType
}

func (s Stretch) WillFitKg(kg float64) bool {
	return kg < float64(s.capacityInGram)
}
