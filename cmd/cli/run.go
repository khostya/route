package main

import (
	"bufio"
	"context"
	"fmt"
	"homework/internal/app"
	"homework/pkg/output"
	"os"
)

func run(ctx context.Context, cancelFunc context.CancelFunc, app *app.App, messages <-chan output.Message[string]) {
	out := bufio.NewWriter(os.Stdout)
	defer cancelFunc()

	for {
		select {
		case <-ctx.Done():
			return
		case <-app.Exit():
			return
		case message := <-messages:
			fmt.Fprintln(out, message.GetMessage())
			_ = out.Flush()
		}
	}
}
