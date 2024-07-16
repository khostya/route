package main

import (
	"context"
	"fmt"
	"homework/cmd"
	"homework/config"
	"homework/internal/app"
	"homework/internal/cli"
	"homework/internal/infrastructure/app/oncall"
	"homework/pkg/output"
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

	outputCFG := config.MustNewOutputConfig()

	controller := output.NewController[output.Message[string]]()
	jobs := getJobs(ctx, getLines())
	orderService, closePG := cmd.GetOrderService(ctx)
	commands := cli.NewCLI(cli.Deps{Service: orderService})

	onCallProducer := cmd.GetOnCallKafkaSender(ctx)
	defer onCallProducer.Close()

	app := app.NewApp(commands, jobs, onCallProducer)
	err := app.Start(ctx, numWorkers)
	if err != nil {
		log.Fatalln(err)
	}
	controller.Add(output.BuildMessageChan[string](output.CLI, app.GetOutput()))

	kafkaMessages, handler := oncall.NewTopicHandler()
	onCallConsumer := cmd.GetOnCallKafkaReceiver(handler)
	controller.Add(output.BuildMessageChan[string](output.Kafka, kafkaMessages))
	defer onCallConsumer.Close()

	output := output.FilterMessageChan(outputCFG.Filter, controller.Subscribe())
	go run(ctx, cancel, app, output)

	app.Wait()
	controller.Close()
	commands.Close()
	closePG()
	_, _ = fmt.Fprintln(os.Stdout, "done")
}
