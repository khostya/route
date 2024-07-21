package main

import (
	"context"
	"fmt"
	"homework/cmd"
	"homework/config"
	"homework/internal/infrastructure/app/oncall"
	"homework/internal/tracer"
	"homework/pkg/output"
	"os"
	"os/signal"
	"syscall"
)

const name = "orders-grpc"

func main() {
	controller := output.NewController[output.Message[string]]()
	outputCFG := config.MustNewOutputConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	tracerCloser := tracer.MustSetup(name)
	defer tracerCloser.Close()

	command, closeDB := cmd.GetOrderService(ctx)
	producer := cmd.GetOnCallKafkaSender(ctx)
	grpcWG := startGrpcServer(ctx, cancel, command, producer)

	if outputCFG.Filter == output.Kafka {
		kafkaMessages, handler := oncall.NewTopicHandler()
		onCallConsumer := cmd.GetOnCallKafkaReceiver(handler)
		controller.Add(output.BuildMessageChan[string](output.Kafka, kafkaMessages))
		defer onCallConsumer.Close()
	}

	filtered := output.FilterMessageChan(outputCFG.Filter, controller.Subscribe())
	go run(ctx, cancel, filtered)

	grpcWG.Wait()
	closeDB()
	_, _ = fmt.Fprintln(os.Stdout, "done")
}
