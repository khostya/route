package schema

import (
	"github.com/shopspring/decimal"
	"homework/internal/model/wrapper"
)

type (
	Wrapper struct {
		OrderID        string                 `db:"order_id"`
		Type           wrapper.WrapperType    `db:"type"`
		PriceInRub     decimal.Decimal        `db:"wrappers_price_in_rub"`
		CapacityInGram wrapper.CapacityInGram `db:"capacity_in_gram"`
	}

	NullableWrapper struct {
		OrderID        *string                 `db:"order_id"`
		Type           *wrapper.WrapperType    `db:"type"`
		PriceInRub     *decimal.Decimal        `db:"wrappers_price_in_rub"`
		CapacityInGram *wrapper.CapacityInGram `db:"capacity_in_gram"`
	}
)

func NewWrapper(wrapper wrapper.Wrapper, orderID string) Wrapper {
	return Wrapper{
		OrderID:        orderID,
		Type:           wrapper.GetType(),
		PriceInRub:     decimal.Decimal(wrapper.GetPriceInRub()),
		CapacityInGram: wrapper.GetCapacityInGram(),
	}
}

func (w Wrapper) Columns() []string {
	return []string{"order_id", "type", "capacity_in_gram", "price_in_rub"}
}

func (w Wrapper) SelectColumns() []string {
	return []string{"order_id", "type", "capacity_in_gram", "wrappers.price_in_rub as wrappers_price_in_rub"}
}

func (w Wrapper) Values() []any {
	return []any{w.OrderID, w.Type, w.CapacityInGram, w.PriceInRub}
}

func extractWrapper(nullableWrapper NullableWrapper) *wrapper.Wrapper {
	if nullableWrapper.OrderID == nil {
		return nil
	}
	return wrapper.NewWrapper(*nullableWrapper.Type, *nullableWrapper.CapacityInGram, wrapper.PriceInRub(*nullableWrapper.PriceInRub))
}
