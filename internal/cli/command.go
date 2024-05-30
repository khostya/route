package cli

import "fmt"

const (
	help = "help"

	deliveryOrder = "delivery"
	returnOrder   = "return"
	issueOrder    = "issue"
	listOrder     = "list"

	refundOrder   = "refund"
	refundedOrder = "refunded"

	exit = "exit"
)

var (
	ErrExit = fmt.Errorf(exit)
)

type command struct {
	name        string
	description string
	usage       string
}

func (c command) String() string {
	return fmt.Sprintf("%s\n  %s\n  %s", c.name, c.description, c.usage)
}
