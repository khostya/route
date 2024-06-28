package postgresql

import (
	"context"
	"homework/internal/model"
	"homework/internal/storage/schema"
)

func (d *DBPool) CreateOrder(ctx context.Context, order model.Order, hash string) error {
	record := schema.NewOrder(order, hash)

	sql := `insert into ozon.orders (
                         id, recipient_id, status, status_updated_at, hash,
                         created_at, expiration_date, weight_in_gram, price_in_rub) 
			values ($1, $2,$3,$4, $5, $6, $7, $8, $9)`

	_, err := d.pool.Exec(ctx, sql, record.ID, record.RecipientID, record.Status, record.StatusUpdatedAt, record.Hash,
		record.CreatedAt, record.ExpirationDate, record.WeightInGram, record.PriceInRub)

	return err
}
