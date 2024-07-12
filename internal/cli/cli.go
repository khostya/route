package cli

import (
	"context"
	"flag"
	"homework/internal/dto"
	"homework/internal/model"
	"homework/pkg/output"
	"slices"
)

type (
	orderService interface {
		Deliver(ctx context.Context, order dto.DeliverOrderParam) error
		ListUserOrders(ctx context.Context, param dto.ListUserOrdersParam) ([]model.Order, error)
		RefundedOrders(ctx context.Context, param dto.PageParam) ([]model.Order, error)
		ReturnOrder(ctx context.Context, id string) error
		IssueOrders(ctx context.Context, ids []string) error
		RefundOrder(ctx context.Context, param dto.RefundOrderParam) error
	}

	Deps struct {
		Service orderService
	}

	CLI struct {
		service                 orderService
		out                     *output.Controller[string]
		commandList             []command
		changeNumberWorkersChan *output.Controller[int]
		exit                    chan struct{}
	}
)

func NewCLI(d Deps) *CLI {
	return &CLI{
		service:                 d.Service,
		commandList:             newCommandList(d.Service),
		changeNumberWorkersChan: output.NewController[int](),
		out:                     output.NewController[string](),
		exit:                    make(chan struct{}, 1),
	}
}

func (c CLI) Run(ctx context.Context, args []string) {
	if len(args) == 0 {
		c.out.Push("command isn't set")
		return
	}

	commandName := args[0]
	switch commandName {
	case help:
		c.help()
		return
	case workers:
		c.changeNumberWorkers(args[1:])
		return
	case exit:
		close(c.exit)
		return
	default:
		handlerIndex := slices.IndexFunc(c.commandList, func(h command) bool {
			return h.name == commandName
		})
		if handlerIndex == -1 {
			break
		}
		out := c.commandList[handlerIndex].handler(ctx, args[1:])
		if out != "" {
			c.out.Push(out)
		}
		return
	}

	c.out.Push("command isn't set")
	return
}

func (c CLI) GetChangeNumberWorkers() <-chan int {
	return c.changeNumberWorkersChan.Subscribe()
}

func (c CLI) GetOutput() chan string {
	return c.out.Subscribe()
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

	c.changeNumberWorkersChan.Push(n)
	return ""
}

func (c CLI) help() {
	c.out.Push("command list:")
	for _, cmd := range c.commandList {
		c.out.Push(cmd.String())
	}
}

func (c CLI) Close() {
	c.out.Close()
	c.changeNumberWorkersChan.Close()
}

func (c CLI) Exit() <-chan struct{} {
	return c.exit
}
