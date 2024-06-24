package model

import "github.com/shopspring/decimal"

const (
	boxType     = WrapperType("box")
	packageType = WrapperType("package")
	stretchType = WrapperType("stretch")
)

type (
	CapacityInGram float64
	PriceInRub     decimal.Decimal

	WrapperType string

	Wrapper interface {
		GetCapacityInGram() CapacityInGram
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

func NewWrapper(t WrapperType, capacityInKg CapacityInGram, priceInRub PriceInRub) (Wrapper, error) {
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

func (p PriceInRub) Add(p2 PriceInRub) PriceInRub {
	return PriceInRub(decimal.Decimal(p).Add(decimal.Decimal(p2)))
}

func GetAllWrapperTypes() []WrapperType {
	return []WrapperType{
		boxType,
		packageType,
		stretchType,
	}
}
