package model

import "github.com/shopspring/decimal"

var (
	boxCapacityInGram CapacityInGram = 30 * 1000
	boxPriceInRub                    = PriceInRub(decimal.NewFromInt(20))
)

type Box struct {
	capacityInGram CapacityInGram
	priceInRub     PriceInRub
}

func NewDefaultBox() Box {
	return NewBox(boxCapacityInGram, boxPriceInRub)
}

func NewBox(capacity CapacityInGram, price PriceInRub) Box {
	return Box{
		capacityInGram: capacity,
		priceInRub:     price,
	}
}

func (b Box) GetCapacityInGram() CapacityInGram {
	return b.capacityInGram
}

func (b Box) GetPriceInRub() PriceInRub {
	return b.priceInRub
}

func (b Box) GetType() WrapperType {
	return boxType
}

func (b Box) WillFitKg(kg float64) bool {
	return kg*1000 < float64(b.capacityInGram)
}
