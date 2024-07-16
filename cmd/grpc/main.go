package main

import (
	"context"
	"fmt"
	"homework/cmd"
	"homework/config"
	"homework/internal/infrastructure/app/oncall"
	"homework/pkg/output"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	controller := output.NewController[output.Message[string]]()
	outputCFG := config.MustNewOutputConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	command, closeDB := cmd.GetOrderService(ctx)
	producer := cmd.GetOnCallKafkaSender(ctx)

	grpcWG := startGrpcServer(ctx, cancel, command, producer)

	out, handler := oncall.NewTopicHandler()
	controller.Add(output.BuildMessageChan[string](output.Kafka, out))

	consumer := cmd.GetOnCallKafkaReceiver(handler)
	defer consumer.Close()

	filtered := output.FilterMessageChan(outputCFG.Filter, controller.Subscribe())
	go run(ctx, cancel, filtered)

	grpcWG.Wait()
	closeDB()
	_, _ = fmt.Fprintln(os.Stdout, "done")
}
