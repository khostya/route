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

	orderService, closePG := cmd.GetOrderService(ctx)
	commands := cli.NewCLI(cli.Deps{
		Service: orderService,
	})

	onCallProducer := cmd.GetOnCallKafkaSender(ctx)
	defer onCallProducer.Close()

	jobs := getJobs(ctx, getLines())
	app := app.NewApp(commands, jobs, onCallProducer)
	err := app.Start(ctx, numWorkers)
	if err != nil {
		log.Fatalln(err)
	}

	if outputCFG.Filter == output.Kafka {
		kafkaMessages, handler := oncall.NewTopicHandler()
		onCallConsumer := cmd.GetOnCallKafkaReceiver(handler)
		controller.Add(output.BuildMessageChan[string](output.Kafka, kafkaMessages))
		defer onCallConsumer.Close()
	}

	controller.Add(output.BuildMessageChan[string](output.CLI, app.GetOutput()))
	output := output.FilterMessageChan[string](outputCFG.Filter, controller.Subscribe())
	go run(ctx, cancel, app, output)

	app.Wait()
	controller.Close()
	commands.Close()
	closePG()
	_, _ = fmt.Fprintln(os.Stdout, "done")
}
