package main

import (
	"homework/internal/cli"
	"sync"
)

type app struct {
	cli  cli.CLI
	jobs <-chan []string
}

func (a app) worker(wg *sync.WaitGroup, result chan<- error) {
	defer wg.Done()

	for args := range a.jobs {
		result <- a.cli.Run(args)
	}
}
