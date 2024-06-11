package cli

import (
	"fmt"
)

const (
	help = "help"

	deliverOrder = "deliver"
	returnOrder  = "return"
	issueOrders  = "issue"
	listOrders   = "list"
	refundOrder  = "refund"
	listRefunded = "refunded"
	procs        = "procs"

	exit = "exit"
)

type (
	commandHandler func([]string) string

	command struct {
		name        string
		description string
		usage       string
		handler     commandHandler
	}
)

func (c command) String() string {
	return fmt.Sprintf("%s\n   %s\n   %s", c.name, c.description, c.usage)
}

func newCommandList(service orderService) []command {
	handlers := newHandlers(service)

	return []command{
		{
			name:        help,
			usage:       help,
			description: helpDescription,
		},
		{
			name:        refundOrder,
			usage:       refundOrderUsage,
			description: refundOrderDescription,
			handler:     handlers.mustFind(refundOrder).handle,
		},
		{
			name:        listRefunded,
			usage:       listRefundedUsage,
			description: listRefundedDescription,
			handler:     handlers.mustFind(listRefunded).handle,
		},
		{
			name:        listOrders,
			usage:       listOrdersUsage,
			description: listOrdersDescription,
			handler:     handlers.mustFind(listOrders).handle,
		},
		{
			name:        issueOrders,
			usage:       issueOrdersUsage,
			description: issueOrdersDescription,
			handler:     handlers.mustFind(issueOrders).handle,
		},
		{
			name:        returnOrder,
			usage:       returnOrderUsage,
			description: returnOrderDescription,
			handler:     handlers.mustFind(returnOrder).handle,
		},
		{
			name:        deliverOrder,
			usage:       deliverOrderUsage,
			description: deliverOrderDescription,
			handler:     handlers.mustFind(deliverOrder).handle,
		},
		{
			name:        procs,
			usage:       procsUsage,
			description: procsDescription,
			handler:     handlers.mustFind(procs).handle,
		},
		{
			name:        exit,
			usage:       exit,
			description: exitDescription,
		},
	}
}
