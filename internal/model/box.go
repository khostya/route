package model

import "github.com/shopspring/decimal"

var (
	boxCapacityInKg CapacityInKg = 30
	boxPriceInRub                = PriceInRub(decimal.NewFromInt(20))
)

type Box struct {
	capacityInKg CapacityInKg
	priceInRub   PriceInRub
}

func NewDefaultBox() Box {
	return NewBox(boxCapacityInKg, boxPriceInRub)
}

func NewBox(capacity CapacityInKg, price PriceInRub) Box {
	return Box{
		capacityInKg: capacity,
		priceInRub:   price,
	}
}

func (b Box) GetCapacityInKg() CapacityInKg {
	return b.capacityInKg
}

func (b Box) GetPriceInRub() PriceInRub {
	return b.priceInRub
}

func (b Box) GetType() WrapperType {
	return boxType
}

func (b Box) WillFitKg(kg float64) bool {
	return kg < float64(b.capacityInKg)
}
