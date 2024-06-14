package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	app2 "homework/internal/app"
	"homework/internal/cli"
	"homework/internal/service"
	"homework/internal/storage"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	fileName   = "orders.json"
	numJobs    = 2
	numWorkers = 2
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	commands := getCommands(out)

	var (
		jobs   = getJobs(ctx)
		result = make(chan error, numJobs)
	)

	app := app2.NewApp(commands, jobs, numWorkers, result, out)

	go func() {
		defer cancel()
		defer app.Stop()

		for {
			select {
			case _ = <-ctx.Done():
				return
			case err := <-result:
				if errors.Is(err, cli.ErrExit) {
					return
				}
				if err != nil {
					_, _ = fmt.Fprintln(out, err)
					_ = out.Flush()
				}
			}
		}
	}()

	app.Wait()
	commands.Close()
	_, _ = fmt.Fprintln(out, "done")
}

func getJobs(ctx context.Context) <-chan []string {
	var (
		lines = make(chan string)
		jobs  = make(chan []string, numJobs)
	)

	go func(lines chan<- string) {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if !scanner.Scan() {
				close(lines)
				return
			}
			line := scanner.Text()
			lines <- line
		}
	}(lines)

	go func(jobs chan<- []string, lines <-chan string) {
		for {
			select {
			case _ = <-ctx.Done():
				close(jobs)
				return
			case line, ok := <-lines:
				if !ok {
					close(jobs)
					return
				}
				args := strings.Split(line, " ")
				jobs <- args
			}
		}
	}(jobs, lines)

	return jobs
}

func getCommands(out *bufio.Writer) *cli.CLI {
	storageJSON, err := storage.NewStorage(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var orderService = service.NewOrder(service.Deps{
		Storage: storageJSON,
	})
	return cli.NewCLI(cli.Deps{
		Service: &orderService,
		Out:     out,
	})
}
