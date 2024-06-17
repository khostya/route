package app

import (
	"bufio"
	"fmt"
	"homework/internal/cli"
	"math"
	"math/rand"
	"sync"
	"time"
)

type App struct {
	cli           *cli.CLI
	jobs          <-chan []string
	numberWorkers int

	stopWorker  chan struct{}
	startWorker chan struct{}

	wg   sync.WaitGroup
	stop chan struct{}
}

func NewApp(commands *cli.CLI, jobs <-chan []string) *App {
	app := &App{
		stopWorker:  make(chan struct{}),
		startWorker: make(chan struct{}),
		cli:         commands,
		jobs:        jobs,
		stop:        make(chan struct{}),
	}
	go app.changeNumberWorkers(commands.GetChangeNumberWorkers())
	return app
}

func (a *App) Start(n int, result chan<- error, out *bufio.Writer) {
	for i := 0; i < n; i++ {
		go a.worker(i, result, out)
		a.wg.Add(1)
		a.numberWorkers++
	}
}

func (a *App) worker(n int, result chan<- error, out *bufio.Writer) {
	defer a.wg.Done()

	for {
		select {
		case <-a.stop:
			return
		case job, ok := <-a.jobs:
			if !ok {
				return
			}
			_, _ = fmt.Fprintf(out, "start: job=%s, n=%v, time=%s\n", job, n, time.Now().Format(time.RFC3339))
			result <- a.cli.Run(job)
			_, _ = fmt.Fprintf(out, "stop: job=%s, n=%v, time=%s\n", job, n, time.Now().Format(time.RFC3339))

			_ = out.Flush()
		case <-a.startWorker:
			go a.worker(rand.Intn(math.MaxInt), result, out)
			a.wg.Add(1)
		case <-a.stopWorker:
			return
		}
	}
}

func (a *App) Wait() {
	a.wg.Wait()
}

func (a *App) Stop() {
	close(a.stop)
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
