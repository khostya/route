package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework/internal/service"
	"homework/internal/storage"
	"homework/internal/storage/transactor"
	pool "homework/pkg/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	pool, err := pool.PoolFromEnv(ctx, "DATABASE_URL")
	if err != nil {
		log.Fatalln(err)
	}

	startGrpcServer(ctx, getService(pool)).Wait()
	pool.Close()
	_, _ = fmt.Fprintln(os.Stdout, "done")
}

func getService(pool *pgxpool.Pool) *service.Order {
	transactionManager := transactor.NewTransactionManager(pool)

	orderStorage := storage.NewOrderStorage(&transactionManager)
	wrapperStorage := storage.NewWrapperStorage(&transactionManager)

	var orderService = service.NewOrder(service.Deps{
		Storage:            orderStorage,
		WrapperStorage:     wrapperStorage,
		TransactionManager: &transactionManager,
	})
	return &orderService
}
