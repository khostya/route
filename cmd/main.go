package main

import (
	"bufio"
	"errors"
	"fmt"
	"homework-1/internal/cli"
	"homework-1/internal/module"
	"homework-1/internal/storage"
	"os"
	"strings"
)

const (
	fileName = "orders.json"
)

func main() {
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	storageJSON, err := storage.NewStorage(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	orderService := module.NewModule(module.Deps{
		Storage: storageJSON,
	})
	commands := cli.NewCLI(cli.Deps{
		Module: orderService,
		Out:    out,
	})

	scanner := bufio.NewScanner(os.Stdin)
	app := app{scanner: scanner, cli: commands}

	for {
		err := app.run()
		if errors.Is(err, cli.ErrExit) {
			break
		}
		if err != nil {
			fmt.Fprintln(out, err)
			out.Flush()
		}
	}

	fmt.Fprintln(out, "done")
}

type app struct {
	scanner *bufio.Scanner
	cli     cli.CLI
}

func (a app) run() error {
	a.scanner.Scan()

	line := a.scanner.Text()
	args := strings.Split(line, " ")
	return a.cli.Run(args)
}
