package wrapper

import (
	"github.com/shopspring/decimal"
	"homework/config"
)

type (
	CapacityInGram float64
	PriceInRub     decimal.Decimal
	WrapperType    string

	Wrapper struct {
		t              WrapperType
		capacityInGram CapacityInGram
		priceInRub     PriceInRub
	}
)

var (
	wrappers Wrappers = make(map[WrapperType]*Wrapper)
)

const (
	BoxWrapper     WrapperType = "box"
	NoneWrapper    WrapperType = ""
	StretchWrapper WrapperType = "stretch"
	PackageWrapper WrapperType = "package"
)

func init() {
	wrappersCFG, err := config.NewWrappersConfig()
	if err != nil {
		panic(err)
	}

	for _, wrapper := range wrappersCFG.Wrappers {
		wrapperType := WrapperType(wrapper.Type)
		capacityInGram := CapacityInGram(wrapper.CapacityInGram)
		priceInRub := PriceInRub(decimal.NewFromFloat(wrapper.PriceInRub))
		wrappers.Add(wrapperType, NewWrapper(wrapperType, capacityInGram, priceInRub))
	}
}

func NewDefaultWrapper(t WrapperType) (*Wrapper, error) {
	wrapper, ok := wrappers.Get(t)
	if !ok {
		return nil, ErrUnknownWrapperType
	}
	return wrapper, nil
}

func NewWrapper(t WrapperType, capacityInGram CapacityInGram, priceInRub PriceInRub) *Wrapper {
	return &Wrapper{t: t, capacityInGram: capacityInGram, priceInRub: priceInRub}
}

func (p PriceInRub) Add(p2 PriceInRub) PriceInRub {
	return PriceInRub(decimal.Decimal(p).Add(decimal.Decimal(p2)))
}

func (w Wrapper) GetCapacityInGram() CapacityInGram {
	return w.capacityInGram
}

func (w Wrapper) GetType() WrapperType {
	return w.t
}

func (w Wrapper) GetPriceInRub() PriceInRub {
	return w.priceInRub
}

func (w Wrapper) WillFitGram(gram float64) bool {
	return gram < float64(w.capacityInGram)
}

func GetAllWrapperTypes() []WrapperType {
	var types []WrapperType
	for key := range wrappers {
		types = append(types, key)
	}
	return types
}
