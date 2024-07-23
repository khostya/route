package main

import (
	"bufio"
	"context"
	"os"
	"strings"
)

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
