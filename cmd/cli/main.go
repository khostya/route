package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework/config"
	"homework/internal/app"
	"homework/internal/cli"
	"homework/internal/infrastructure/app/oncall"
	"homework/internal/service"
	"homework/internal/storage"
	"homework/internal/storage/transactor"
	output2 "homework/pkg/output"
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
	kafkaCFG, err := config.NewKafkaConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	pool, err := pool.PoolFromEnv(ctx, "DATABASE_URL")
	if err != nil {
		log.Fatalln(err)
	}

	commands := cli.NewCLI(cli.Deps{
		Service: getService(pool),
	})

	var (
		jobs        = getJobs(ctx, getLines())
		cliMessages = make(chan error, numJobs)
	)

	onCallProducer, err := getOnCallKafkaSender(ctx, kafkaCFG)
	if err != nil {
		log.Fatalln(err)
	}
	defer onCallProducer.Close()

	kafkaMessages, handler := oncall.NewTopicHandler()
	onCallConsumer, err := getOnCallKafkaReceiver(kafkaCFG, handler)
	if err != nil {
		log.Fatalln(err)
	}
	defer onCallConsumer.Close()

	app := app.NewApp(commands, jobs, onCallProducer)
	err = app.Start(ctx, numWorkers, cliMessages)
	if err != nil {
		log.Fatal(err)
	}

	output := output2.NewController[output2.Message[string]]()
	output.Add(output2.BuildMessageChan[string](output2.CLI, app.GetOutput()))
	output.Add(output2.BuildMessageChan[string](output2.CLI, commands.GetOutput()))
	output.Add(output2.BuildMessageChan[string](output2.Kafka, kafkaMessages))

	cfg, err := config.NewOutputConfig()
	if err != nil {
		log.Fatal(err)
	}

	go run(ctx, cancel, app, output2.FilterMessageChan[string](cfg.Filter, output.Subscribe()))

	app.Wait()
	output.Close()
	commands.Close()
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
