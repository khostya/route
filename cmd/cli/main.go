package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework/config"
	"homework/internal/app"
	"homework/internal/cli"
	"homework/internal/infrastructure/app/oncall"
	"homework/internal/service"
	"homework/internal/storage"
	"homework/internal/storage/transactor"
	"homework/pkg/output"
	pool "homework/pkg/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	numJobs    = 2
	numWorkers = 2
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	kafkaCFG := config.NewMustKafkaConfig()
	outputCFG := config.NewMustOutputConfig()

	controller := output.NewController[output.Message[string]]()
	jobs := getJobs(ctx, getLines())
	commands, closePG := getCommands(ctx)

	onCallProducer := getOnCallKafkaSender(ctx, kafkaCFG)
	defer onCallProducer.Close()

	app := app.NewApp(commands, jobs, onCallProducer)
	err := app.Start(ctx, numWorkers)
	if err != nil {
		log.Fatalln(err)
	}

	if outputCFG.Filter == output.Kafka {
		kafkaMessages, handler := oncall.NewTopicHandler()
		onCallConsumer := getOnCallKafkaReceiver(kafkaCFG, handler)

		controller.Add(output.BuildMessageChan[string](output.Kafka, kafkaMessages))
		defer onCallConsumer.Close()
	} else {
		controller.Add(output.BuildMessageChan[string](output.CLI, app.GetOutput()))
	}

	go run(ctx, cancel, app, controller.Subscribe())

	app.Wait()
	controller.Close()
	commands.Close()
	closePG()
	_, _ = fmt.Fprintln(os.Stdout, "done")
}

func getCommands(ctx context.Context) (*cli.CLI, func()) {
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
	return cli.NewCLI(cli.Deps{
		Service: &orderService,
	}), pool.Close
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
