package main

import (
	"bufio"
	"fmt"
	"homework/internal/cli"
	"sync"
	"time"
)

type app struct {
	cli  cli.CLI
	jobs <-chan []string
}

func (a app) worker(n int, wg *sync.WaitGroup, result chan<- error, out *bufio.Writer) {
	defer wg.Done()

	for args := range a.jobs {
		_, _ = fmt.Fprintf(out, "start: job=%s, n=%v, time=%s\n", args, n, time.Now().Format(time.RFC3339))
		result <- a.cli.Run(args)
		_, _ = fmt.Fprintf(out, "stop: job=%s, n=%v, time=%s\n", args, n, time.Now().Format(time.RFC3339))

		_ = out.Flush()
	}
}
