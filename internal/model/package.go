package model

const (
	packageCapacityInKg CapacityInKg = 10
	packagePriceInRub   PriceInRub   = 5
)

type Package struct {
	capacityInKg CapacityInKg
	priceInRub   PriceInRub
}

func NewDefaultPackage() Package {
	return NewPackage(packageCapacityInKg, packagePriceInRub)
}

func NewPackage(capacity CapacityInKg, price PriceInRub) Package {
	return Package{
		capacityInKg: capacity,
		priceInRub:   price,
	}
}

func (p Package) GetCapacityInKg() CapacityInKg {
	return packageCapacityInKg
}

func (p Package) GetPriceInRub() PriceInRub {
	return packagePriceInRub
}

func (p Package) GetType() WrapperType {
	return packageType
}

func (p Package) WillFitKg(kg float64) bool {
	return kg < float64(p.capacityInKg)
}
