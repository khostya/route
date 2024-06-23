package schema

import (
	"github.com/shopspring/decimal"
	"homework/internal/model"
)

type (
	Wrapper struct {
		OrderID      string             `db:"order_id"`
		Type         model.WrapperType  `db:"type"`
		PriceInRub   decimal.Decimal    `db:"price_in_rub"`
		CapacityInKg model.CapacityInKg `db:"capacity_in_kg"`
	}

	NullableWrapper struct {
		OrderID      *string             `db:"order_id"`
		Type         *model.WrapperType  `db:"type"`
		PriceInRub   *decimal.Decimal    `db:"price_in_rub"`
		CapacityInKg *model.CapacityInKg `db:"capacity_in_kg"`
	}
)

func NewWrapper(wrapper model.Wrapper, orderID string) Wrapper {
	return Wrapper{
		OrderID:      orderID,
		Type:         wrapper.GetType(),
		PriceInRub:   decimal.Decimal(wrapper.GetPriceInRub()),
		CapacityInKg: wrapper.GetCapacityInKg(),
	}
}

func (w Wrapper) Columns() []string {
	return []string{"order_id", "type", "capacity_in_kg", "price_in_rub"}
}

func (w Wrapper) Values() []any {
	return []any{w.OrderID, w.Type, w.CapacityInKg, w.PriceInRub}
}

func extractWrapper(wrapper NullableWrapper) (model.Wrapper, error) {
	if wrapper.OrderID == nil {
		return nil, nil
	}
	return model.NewWrapper(*wrapper.Type, *wrapper.CapacityInKg, model.PriceInRub(*wrapper.PriceInRub))
}
