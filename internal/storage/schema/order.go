package schema

import (
	"github.com/shopspring/decimal"
	"homework/internal/model"
	"homework/internal/model/wrapper"
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

		Hash         string          `db:"hash"`
		CreatedAt    time.Time       `db:"created_at"`
		WeightInGram float64         `db:"weight_in_gram"`
		PriceInRub   decimal.Decimal `db:"orders_price_in_rub"`
	}
)

func NewOrder(order model.Order, hash string) Order {
	return Order{
		ID:              order.ID,
		RecipientID:     order.RecipientID,
		Status:          order.Status,
		StatusUpdatedAt: order.StatusUpdatedAt,
		ExpirationDate:  order.ExpirationDate,
		WeightInGram:    order.WeightInGram,
		Hash:            hash,
		PriceInRub:      decimal.Decimal(order.PriceInRub),
		CreatedAt:       time.Now(),
	}
}

func (o Order) Columns() []string {
	return []string{
		"id", "recipient_id", "status", "status_updated_at",
		"expiration_date", "hash", "created_at", "weight_in_gram", "price_in_rub",
	}
}

func (o Order) SelectColumns() []string {
	return []string{
		"id", "recipient_id", "status", "status_updated_at",
		"expiration_date", "hash", "created_at", "weight_in_gram", "orders.price_in_rub as orders_price_in_rub",
	}
}

func (o Order) Values() []any {
	return []any{
		o.ID, o.RecipientID, o.Status, o.StatusUpdatedAt,
		o.ExpirationDate, o.Hash, o.CreatedAt, o.WeightInGram, o.PriceInRub,
	}
}

func ExtractOrdersFromWrapperOrder(records []WrapperOrder) ([]model.Order, error) {
	return mapFuncErr(records, func(record WrapperOrder) (model.Order, error) {
		order := record.Order
		wrapperModel := extractWrapper(record.NullableWrapper)
		return model.Order{
			ID:              order.ID,
			RecipientID:     order.RecipientID,
			Status:          order.Status,
			StatusUpdatedAt: order.StatusUpdatedAt,
			ExpirationDate:  order.ExpirationDate,
			WeightInGram:    order.WeightInGram,
			PriceInRub:      wrapper.PriceInRub(order.PriceInRub),
			Wrapper:         wrapperModel,
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
