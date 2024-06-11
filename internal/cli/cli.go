package cli

import (
	"bufio"
	"fmt"
	"homework/internal/model"
	"homework/internal/service"
	"slices"
)

type (
	orderService interface {
		Deliver(order service.DeliverOrderParam) error
		Orders(userID string, count int) ([]model.Order, error)
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
		service     orderService
		out         *bufio.Writer
		commandList []command
	}
)

func NewCLI(d Deps) CLI {
	return CLI{
		service:     d.Service,
		out:         d.Out,
		commandList: newCommandList(d.Service),
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

func (c CLI) help() {
	fmt.Fprintln(c.out, "command list:")
	for _, cmd := range c.commandList {
		fmt.Fprintln(c.out, cmd)
	}
}
