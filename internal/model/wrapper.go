package model

const (
	boxType     = WrapperType("box")
	packageType = WrapperType("package")
	stretchType = WrapperType("stretch")
)

type (
	CapacityInKg float64
	PriceInRub   float64

	WrapperType string

	Wrapper interface {
		GetCapacityInKg() CapacityInKg
		GetType() WrapperType
		GetPriceInRub() PriceInRub
		WillFitKg(kg float64) bool
	}
)

func NewDefaultWrapper(t WrapperType) (Wrapper, error) {
	switch t {
	case boxType:
		return NewDefaultBox(), nil
	case packageType:
		return NewDefaultPackage(), nil
	case stretchType:
		return NewDefaultStretch(), nil
	default:
		return nil, ErrUnknownWrapperType
	}
}

func NewWrapper(t WrapperType, capacityInKg CapacityInKg, priceInRub PriceInRub) (Wrapper, error) {
	switch t {
	case boxType:
		return NewBox(capacityInKg, priceInRub), nil
	case packageType:
		return NewPackage(capacityInKg, priceInRub), nil
	case stretchType:
		return NewStretch(capacityInKg, priceInRub), nil
	default:
		return nil, ErrUnknownWrapperType
	}
}

func GetAllWrapperTypes() []WrapperType {
	return []WrapperType{
		boxType,
		packageType,
		stretchType,
	}
}
