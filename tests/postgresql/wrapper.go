package postgresql

import (
	"context"
	"github.com/jackc/pgx/v4"
	"homework/internal/model"
	"homework/internal/storage/schema"
)

func (d *DBPool) CreateWrapper(ctx context.Context, order model.Order, hash string) error {
	record := schema.NewOrder(order, hash)
	recordWrapper := schema.NewWrapper(*order.Wrapper, order.ID)
	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
insert into ozon.orders (
                         id, recipient_id, status, status_updated_at, hash,
                         created_at, expiration_date, weight_in_gram, price_in_rub)
values ($1, $2,$3,$4, $5, $6, $7, $8, $9)`,
		record.ID, record.RecipientID, record.Status, record.StatusUpdatedAt, record.Hash,
		record.CreatedAt, record.ExpirationDate, record.WeightInGram, record.PriceInRub)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
insert into ozon.wrappers (
                           order_id, type, price_in_rub, capacity_in_gram) 
values ($1,$2,$3, $4)`,
		recordWrapper.OrderID, recordWrapper.Type, recordWrapper.PriceInRub, recordWrapper.CapacityInGram)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
