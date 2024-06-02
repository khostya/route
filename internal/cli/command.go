package cli

import "fmt"

const (
	help = "help"

	deliverOrder = "deliver"
	returnOrder  = "return"
	issueOrder   = "issue"
	listOrder    = "list"
	refundOrder  = "refund"
	listRefunded = "refunded"

	exit = "exit"
)

type (
	commandHandler func([]string) string

	command struct {
		name        string
		description string
		usage       string
	}
)

func (c command) String() string {
	return fmt.Sprintf("%s\n   %s\n   %s", c.name, c.description, c.usage)
}

func newCommandList() []command {
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
		},
		{
			name:        listRefunded,
			usage:       listRefundedUsage,
			description: listRefundedDescription,
		},
		{
			name:        listOrder,
			usage:       listOrderUsage,
			description: listOrderDescription,
		},
		{
			name:        issueOrder,
			usage:       issueOrderUsage,
			description: issueOrderDescription,
		},
		{
			name:        returnOrder,
			usage:       returnOrderUsage,
			description: returnOrderDescription,
		},
		{
			name:        deliverOrder,
			usage:       deliverOrderUsage,
			description: deliverOrderDescription,
		},
		{
			name:        exit,
			usage:       exit,
			description: exitDescription,
		},
	}
}
