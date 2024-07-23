package app

import (
	"context"
	"fmt"
	"homework/internal/dto"
	"homework/internal/model"
	"homework/pkg/output"
	"math"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type (
	cli interface {
		Run(ctx context.Context, args []string)
		GetChangeNumberWorkers() <-chan int
		GetOutput() <-chan string
		Exit() <-chan struct{}
	}

	onCallProducer interface {
		SendAsyncMessage(message dto.OnCallMessage) error
	}

	App struct {
		cli           cli
		jobs          <-chan []string
		numberWorkers int

		stopWorker  chan struct{}
		startWorker chan struct{}

		wg sync.WaitGroup

		isStarted atomic.Bool
		onCall    onCallProducer
		output    *output.Controller[string]
	}
)

func NewApp(commands cli, jobs <-chan []string, onCall onCallProducer) *App {
	output := output.NewController[string]()
	output.Add(commands.GetOutput())

	return &App{
		stopWorker:  make(chan struct{}),
		startWorker: make(chan struct{}),
		cli:         commands,
		jobs:        jobs,
		onCall:      onCall,
		output:      output,
	}
}

func (a *App) Start(ctx context.Context, n int) error {
	if a.isStarted.Swap(true) {
		return ErrHasAlreadyStarted
	}

	go a.changeNumberWorkers(a.cli.GetChangeNumberWorkers())
	for i := 0; i < n; i++ {
		go a.worker(ctx, i)
		a.wg.Add(1)
		a.numberWorkers++
	}
	return nil
}

func (a *App) worker(ctx context.Context, n int) {
	defer a.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-a.jobs:
			if !ok {
				return
			}

			a.output.Push(fmt.Sprintf("start: job=%s, n=%v, time=%s\n", job, n, time.Now().Format(model.TimeFormat)))
			a.cli.Run(ctx, job)
			a.output.Push(fmt.Sprintf("stop: job=%s, n=%v, time=%s\n", job, n, time.Now().Format(model.TimeFormat)))

			_ = a.onCall.SendAsyncMessage(dto.OnCallMessage{
				CalledAt: time.Now(),
				Method:   job[0],
				Args:     strings.Join(job[1:], " "),
			})

		case <-a.startWorker:
			go a.worker(ctx, rand.Intn(math.MaxInt))
			a.wg.Add(1)
		case <-a.stopWorker:
			return
		}
	}
}

func (a *App) Wait() {
	a.wg.Wait()
}

func (a *App) GetOutput() <-chan string {
	return a.output.Subscribe()
}

func (a *App) Exit() <-chan struct{} {
	return a.cli.Exit()
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
