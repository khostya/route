package model

import (
	"github.com/shopspring/decimal"
	"math"
)

var (
	stretchCapacityInKg = CapacityInKg(math.Inf(1))
	stretchPriceInRub   = PriceInRub(decimal.NewFromInt(0))
)

type Stretch struct {
	capacityInKg CapacityInKg
	priceInRub   PriceInRub
}

func NewDefaultStretch() Stretch {
	return NewStretch(stretchCapacityInKg, stretchPriceInRub)
}

func NewStretch(capacity CapacityInKg, price PriceInRub) Stretch {
	return Stretch{
		capacityInKg: capacity,
		priceInRub:   price,
	}
}

func (s Stretch) GetCapacityInKg() CapacityInKg {
	return s.capacityInKg
}

func (s Stretch) GetPriceInRub() PriceInRub {
	return s.priceInRub
}

func (s Stretch) GetType() WrapperType {
	return stretchType
}

func (s Stretch) WillFitKg(kg float64) bool {
	return kg < float64(s.capacityInKg)
}
