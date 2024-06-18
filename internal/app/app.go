package app

import (
	"bufio"
	"context"
	"fmt"
	"homework/internal/cli"
	"homework/internal/model"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type App struct {
	cli           *cli.CLI
	jobs          <-chan []string
	numberWorkers int

	stopWorker  chan struct{}
	startWorker chan struct{}

	wg sync.WaitGroup

	isStarted atomic.Bool
}

func NewApp(commands *cli.CLI, jobs <-chan []string) *App {
	return &App{
		stopWorker:  make(chan struct{}),
		startWorker: make(chan struct{}),
		cli:         commands,
		jobs:        jobs,
	}
}

func (a *App) Start(ctx context.Context, n int, result chan<- error, out *bufio.Writer) error {
	if a.isStarted.Swap(true) {
		return ErrHasAlreadyStarted
	}

	go a.changeNumberWorkers(a.cli.GetChangeNumberWorkers())
	for i := 0; i < n; i++ {
		go a.worker(ctx, i, result, out)
		a.wg.Add(1)
		a.numberWorkers++
	}
	return nil
}

func (a *App) worker(ctx context.Context, n int, result chan<- error, out *bufio.Writer) {
	defer a.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-a.jobs:
			if !ok {
				return
			}
			_, _ = fmt.Fprintf(out, "start: job=%s, n=%v, time=%s\n", job, n, time.Now().Format(model.TimeFormat))
			result <- a.cli.Run(ctx, job)
			_, _ = fmt.Fprintf(out, "stop: job=%s, n=%v, time=%s\n", job, n, time.Now().Format(model.TimeFormat))

			_ = out.Flush()
		case <-a.startWorker:
			go a.worker(ctx, rand.Intn(math.MaxInt), result, out)
			a.wg.Add(1)
		case <-a.stopWorker:
			return
		}
	}
}

func (a *App) Wait() {
	a.wg.Wait()
}

func (a *App) changeNumberWorkers(workers <-chan int) {
	for number := range workers {
		for a.numberWorkers < number {
			a.startWorker <- struct{}{}
			a.numberWorkers++
		}
		for a.numberWorkers > number {
			a.stopWorker <- struct{}{}
			a.numberWorkers--
		}
	}
}
