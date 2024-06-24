package wrapper

import (
	"github.com/shopspring/decimal"
)

var (
	packageCapacityInGram CapacityInGram = 10 * 1000
	packagePriceInRub                    = PriceInRub(decimal.NewFromInt(5))
)

type packageModel struct {
	capacityInGram CapacityInGram
	priceInRub     PriceInRub
}

func newDefaultPackage() packageModel {
	return newPackage(packageCapacityInGram, packagePriceInRub)
}

func newPackage(capacity CapacityInGram, price PriceInRub) packageModel {
	return packageModel{
		capacityInGram: capacity,
		priceInRub:     price,
	}
}

func (p packageModel) GetCapacityInGram() CapacityInGram {
	return packageCapacityInGram
}

func (p packageModel) GetPriceInRub() PriceInRub {
	return packagePriceInRub
}

func (p packageModel) GetType() WrapperType {
	return packageType
}

func (p packageModel) WillFitKg(kg float64) bool {
	return kg < float64(p.capacityInGram)
}
