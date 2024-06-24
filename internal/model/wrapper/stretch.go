package wrapper

import (
	"github.com/shopspring/decimal"
	"math"
)

var (
	stretchCapacityInGram = CapacityInGram(math.Inf(1))
	stretchPriceInRub     = PriceInRub(decimal.NewFromInt(0))
)

type stretch struct {
	capacityInGram CapacityInGram
	priceInRub     PriceInRub
}

func newDefaultStretch() stretch {
	return newStretch(stretchCapacityInGram, stretchPriceInRub)
}

func newStretch(capacity CapacityInGram, price PriceInRub) stretch {
	return stretch{
		capacityInGram: capacity,
		priceInRub:     price,
	}
}

func (s stretch) GetCapacityInGram() CapacityInGram {
	return s.capacityInGram
}

func (s stretch) GetPriceInRub() PriceInRub {
	return s.priceInRub
}

func (s stretch) GetType() WrapperType {
	return stretchType
}

func (s stretch) WillFitKg(kg float64) bool {
	return kg < float64(s.capacityInGram)
}
