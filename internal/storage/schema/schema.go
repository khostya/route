package schema

import (
	"homework/internal/model"
	"time"
)

type (
	Record struct {
		ID          string `db:"id"`
		RecipientID string `db:"recipient_id"`

		Status          model.Status `db:"status"`
		StatusUpdatedAt time.Time    `db:"status_updated_at"`

		ExpirationDate time.Time `db:"expiration_date"`

		Hash      string    `db:"hash"`
		CreatedAt time.Time `db:"created_at"`
	}

	PageParam struct {
		Size int
		Page int
	}

	GetParam struct {
		Ids         []string
		Status      *model.Status
		Order       *string
		Limit       *int
		RecipientId *string
		Offset      *int
	}
)

func NewRecord(order model.Order, hash string) Record {
	return Record{
		ID:              order.ID,
		RecipientID:     order.RecipientID,
		Status:          order.Status,
		StatusUpdatedAt: order.StatusUpdatedAt,
		ExpirationDate:  order.ExpirationDate,
		Hash:            hash,
		CreatedAt:       time.Now(),
	}
}

func (r Record) Columns() []string {
	return []string{"id", "recipient_id", "status", "status_updated_at", "expiration_date", "hash", "created_at"}
}

func (r Record) Values() []any {
	return []any{r.ID, r.RecipientID, r.Status, r.StatusUpdatedAt, r.ExpirationDate, r.Hash, r.CreatedAt}
}

func ExtractOrders(records []Record) []model.Order {
	return mapFunc(records, func(record Record) model.Order {
		return model.Order{
			ID:              record.ID,
			RecipientID:     record.RecipientID,
			Status:          record.Status,
			StatusUpdatedAt: record.StatusUpdatedAt,
			ExpirationDate:  record.ExpirationDate,
		}
	})
}

func mapFunc[IN any, OUT any](in []IN, m func(IN) OUT) []OUT {
	var out []OUT

	for _, i := range in {
		out = append(out, m(i))
	}

	return out
}
