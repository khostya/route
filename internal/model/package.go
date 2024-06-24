package model

import "github.com/shopspring/decimal"

var (
	packageCapacityInGram CapacityInGram = 10 * 1000
	packagePriceInRub                    = PriceInRub(decimal.NewFromInt(5))
)

type Package struct {
	capacityInGram CapacityInGram
	priceInRub     PriceInRub
}

func NewDefaultPackage() Package {
	return NewPackage(packageCapacityInGram, packagePriceInRub)
}

func NewPackage(capacity CapacityInGram, price PriceInRub) Package {
	return Package{
		capacityInGram: capacity,
		priceInRub:     price,
	}
}

func (p Package) GetCapacityInGram() CapacityInGram {
	return packageCapacityInGram
}

func (p Package) GetPriceInRub() PriceInRub {
	return packagePriceInRub
}

func (p Package) GetType() WrapperType {
	return packageType
}

func (p Package) WillFitKg(kg float64) bool {
	return kg < float64(p.capacityInGram)
}
