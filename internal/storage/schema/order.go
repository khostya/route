package schema

import (
	"homework/internal/model"
	"time"
)

type (
	WrapperOrder struct {
		NullableWrapper
		Order
	}

	Order struct {
		ID          string `db:"id"`
		RecipientID string `db:"recipient_id"`

		Status          model.Status `db:"status"`
		StatusUpdatedAt time.Time    `db:"status_updated_at"`

		ExpirationDate time.Time `db:"expiration_date"`

		Hash       string    `db:"hash"`
		CreatedAt  time.Time `db:"created_at"`
		WeightInKg float64   `db:"weight_in_kg"`
	}
)

func NewOrder(order model.Order, hash string) Order {
	return Order{
		ID:              order.ID,
		RecipientID:     order.RecipientID,
		Status:          order.Status,
		StatusUpdatedAt: order.StatusUpdatedAt,
		ExpirationDate:  order.ExpirationDate,
		WeightInKg:      order.WeightInKg,
		Hash:            hash,
		CreatedAt:       time.Now(),
	}
}

func (o Order) Columns() []string {
	return []string{
		"id", "recipient_id", "status", "status_updated_at",
		"expiration_date", "hash", "created_at", "weight_in_kg",
	}
}

func (o Order) Values() []any {
	return []any{
		o.ID, o.RecipientID, o.Status, o.StatusUpdatedAt,
		o.ExpirationDate, o.Hash, o.CreatedAt, o.WeightInKg,
	}
}

func ExtractOrders(records []Order) []model.Order {
	return mapFunc(records, func(record Order) model.Order {
		return model.Order{
			ID:              record.ID,
			RecipientID:     record.RecipientID,
			Status:          record.Status,
			StatusUpdatedAt: record.StatusUpdatedAt,
			ExpirationDate:  record.ExpirationDate,
		}
	})
}

func ExtractOrdersFromWrapperOrder(records []WrapperOrder) ([]model.Order, error) {
	return mapFuncErr(records, func(record WrapperOrder) (model.Order, error) {
		order := record.Order

		wrapper, err := extractWrapper(record.NullableWrapper)
		if err != nil {
			return model.Order{}, err
		}

		var priceInRub model.PriceInRub
		if wrapper != nil {
			priceInRub += wrapper.GetPriceInRub()
		}

		return model.Order{
			ID:              order.ID,
			RecipientID:     order.RecipientID,
			Status:          order.Status,
			StatusUpdatedAt: order.StatusUpdatedAt,
			ExpirationDate:  order.ExpirationDate,
			WeightInKg:      order.WeightInKg,
			PriceInRub:      priceInRub,
			Wrapper:         wrapper,
		}, nil
	})
}

func mapFunc[IN any, OUT any](in []IN, m func(IN) OUT) []OUT {
	out, _ := mapFuncErr(in, func(i IN) (OUT, error) {
		return m(i), nil
	})
	return out
}

func mapFuncErr[IN any, OUT any](in []IN, m func(IN) (OUT, error)) ([]OUT, error) {
	var out []OUT

	for _, i := range in {
		res, err := m(i)
		if err != nil {
			return nil, err
		}
		out = append(out, res)
	}

	return out, nil
}
