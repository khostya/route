package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	app2 "homework/internal/app"
	"homework/internal/cli"
	"homework/internal/service"
	"homework/internal/storage"
	"homework/internal/storage/transactor"
	pool "homework/pkg/postgres"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	numJobs    = 2
	numWorkers = 2
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	pool, err := getPool(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	commands, err := getCommands(out, pool)
	if err != nil {
		log.Fatalln(err)
	}

	var (
		jobs   = getJobs(ctx, getLines())
		result = make(chan error, numJobs)
	)

	app := app2.NewApp(ctx, commands, jobs, numWorkers, result, out)

	go func() {
		defer cancel()

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
	pool.Close()
	_, _ = fmt.Fprintln(out, "done")
}

func getLines() chan string {
	lines := make(chan string)

	go func(lines chan<- string) {
		defer close(lines)
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if !scanner.Scan() {
				return
			}
			line := scanner.Text()
			lines <- line
		}
	}(lines)

	return lines
}

func getJobs(ctx context.Context, lines chan string) <-chan []string {
	jobs := make(chan []string, numJobs)

	go func(jobs chan<- []string, lines <-chan string) {
		defer close(jobs)
		for {
			select {
			case _ = <-ctx.Done():
				return
			case line, ok := <-lines:
				if !ok {
					return
				}
				args := strings.Split(line, " ")
				jobs <- args
			}
		}
	}(jobs, lines)

	return jobs
}

func getCommands(out *bufio.Writer, pool *pgxpool.Pool) (*cli.CLI, error) {
	transactionManager := transactor.NewTransactionManager(pool)

	storage := storage.NewStorage(&transactionManager)

	var orderService = service.NewOrder(service.Deps{
		Storage:            storage,
		TransactionManager: &transactionManager,
	})
	return cli.NewCLI(cli.Deps{
		Service: &orderService,
		Out:     out,
	}), nil
}

func getPool(ctx context.Context) (*pgxpool.Pool, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, errors.New("Unable to parse DATABASE_URL")
	}

	pool, err := pool.Pool(ctx, url)
	if err != nil {
		return nil, err
	}
	return pool, err
}
