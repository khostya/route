package cmd

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework/internal/service"
	"homework/internal/storage"
	"homework/internal/storage/transactor"
	pool "homework/pkg/postgres"
	"log"
	"os"
)

func GetOrderService(ctx context.Context) (*service.OrderService, func()) {
	pool, err := getPool(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	transactionManager := transactor.NewTransactionManager(pool)

	orderStorage := storage.NewOrderStorage(&transactionManager)
	wrapperStorage := storage.NewWrapperStorage(&transactionManager)

	var orderService = service.NewOrder(service.Deps{
		Storage:            orderStorage,
		WrapperStorage:     wrapperStorage,
		TransactionManager: &transactionManager,
	})
	return &orderService, pool.Close
}

func getPool(ctx context.Context) (*pgxpool.Pool, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, errors.New("unable to parse DATABASE_URL")
	}

	pool, err := pool.Pool(ctx, url)
	if err != nil {
		return nil, err
	}
	return pool, err
}
