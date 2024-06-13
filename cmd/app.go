package main

import (
	"bufio"
	"fmt"
	"homework/internal/cli"
	"math"
	"math/rand"
	"sync"
	"time"
)

type app struct {
	cli           *cli.CLI
	jobs          <-chan []string
	numberWorkers int

	stopWorker  chan struct{}
	startWorker chan struct{}
}

func newApp(cli *cli.CLI, jobs <-chan []string) *app {
	return &app{
		stopWorker:  make(chan struct{}),
		startWorker: make(chan struct{}),
		cli:         cli,
		jobs:        jobs,
	}
}
func (a *app) RunWorkers(n int, result chan<- error, out *bufio.Writer) *sync.WaitGroup {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		go a.worker(i, &wg, result, out)
		wg.Add(1)
		a.numberWorkers++
	}
	return &wg
}

func (a *app) worker(n int, wg *sync.WaitGroup, result chan<- error, out *bufio.Writer) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-a.jobs:
			if !ok {
				return
			}
			_, _ = fmt.Fprintf(out, "start: job=%s, n=%v, time=%s\n", job, n, time.Now().Format(time.RFC3339))
			result <- a.cli.Run(job)
			_, _ = fmt.Fprintf(out, "stop: job=%s, n=%v, time=%s\n", job, n, time.Now().Format(time.RFC3339))

			_ = out.Flush()
		case <-a.startWorker:
			go a.worker(rand.Intn(math.MaxInt), wg, result, out)
			wg.Add(1)
		case <-a.stopWorker:
			return
		}
	}
}

func (a *app) changeNumberWorkers(workers <-chan int) {
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
