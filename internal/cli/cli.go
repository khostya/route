package cli

import (
	"bufio"
	"flag"
	"fmt"
	"homework/internal/model"
	"homework/internal/service"
	"slices"
)

type (
	orderService interface {
		Deliver(order service.DeliverOrderParam) error
		ListUserOrders(userID string, count int) ([]model.Order, error)
		RefundedOrders(param service.RefundedOrdersParam) ([]model.Order, error)
		ReturnOrder(id string) error
		IssueOrders(ids []string) error
		RefundOrder(param service.RefundOrderParam) error
	}

	Deps struct {
		Service orderService
		Out     *bufio.Writer
	}

	CLI struct {
		service                 orderService
		out                     *bufio.Writer
		commandList             []command
		changeNumberWorkersChan chan int
	}
)

func NewCLI(d Deps) *CLI {
	changeNumberWorkers := make(chan int)
	return &CLI{
		service:                 d.Service,
		out:                     d.Out,
		commandList:             newCommandList(d.Service),
		changeNumberWorkersChan: changeNumberWorkers,
	}
}

func (c CLI) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("command isn't set")
	}
	defer c.out.Flush()

	commandName := args[0]
	switch commandName {
	case help:
		c.help()
		return nil
	case workers:
		c.changeNumberWorkers(args[1:])
		return nil
	case exit:
		return ErrExit
	default:
		handlerIndex := slices.IndexFunc(c.commandList, func(h command) bool {
			return h.name == commandName
		})
		if handlerIndex == -1 {
			break
		}
		out := c.commandList[handlerIndex].handler(args[1:])
		if out == "" {
			return nil
		}
		fmt.Fprintln(c.out, out)
		return nil
	}

	return fmt.Errorf("command isn't set")
}

func (c CLI) GetChangeNumberWorkers() <-chan int {
	return c.changeNumberWorkersChan
}

func (c CLI) changeNumberWorkers(args []string) string {
	var n int

	fs := flag.NewFlagSet(workers, flag.ContinueOnError)
	fs.IntVar(&n, "n", -1, workersUsage)
	if err := fs.Parse(args); err != nil {
		return err.Error()
	}
	if n <= 0 {
		return "N isn`t set"
	}

	c.changeNumberWorkersChan <- n
	return ""
}

func (c CLI) help() {
	fmt.Fprintln(c.out, "command list:")
	for _, cmd := range c.commandList {
		fmt.Fprintln(c.out, cmd)
	}
}

func (c CLI) Close() {
	close(c.changeNumberWorkersChan)
}
