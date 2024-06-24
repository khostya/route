package wrapper

import (
	"github.com/shopspring/decimal"
)

var (
	boxCapacityInGram CapacityInGram = 30 * 1000
	boxPriceInRub                    = PriceInRub(decimal.NewFromInt(20))
)

type box struct {
	capacityInGram CapacityInGram
	priceInRub     PriceInRub
}

func newDefaultBox() box {
	return newBox(boxCapacityInGram, boxPriceInRub)
}

func newBox(capacity CapacityInGram, price PriceInRub) box {
	return box{
		capacityInGram: capacity,
		priceInRub:     price,
	}
}

func (b box) GetCapacityInGram() CapacityInGram {
	return b.capacityInGram
}

func (b box) GetPriceInRub() PriceInRub {
	return b.priceInRub
}

func (b box) GetType() WrapperType {
	return boxType
}

func (b box) WillFitKg(kg float64) bool {
	return kg*1000 < float64(b.capacityInGram)
}
