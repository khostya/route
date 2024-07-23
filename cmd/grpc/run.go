package main

import (
	"bufio"
	"context"
	"homework/pkg/output"
	"log"
	"os"
)

func run(ctx context.Context, cancelFunc context.CancelFunc, messages <-chan output.Message[string]) {
	out := bufio.NewWriter(os.Stdout)
	defer cancelFunc()

	for {
		select {
		case <-ctx.Done():
			return
		case message := <-messages:
			log.Println(message.GetMessage())
			_ = out.Flush()
		}
	}
}
