package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"homework-1/internal/model"
	"homework-1/internal/module"
	"math"
	"time"
)

type Module interface {
	Delivery(order module.Order) error
	ListOrder(userID string, count int) ([]model.Order, error)
	RefundedOrder(param module.RefundedOrderParam) ([]model.Order, error)
	ReturnOrder(id string) error
	IssueOrders(ids []string) error
	RefundOrder(param module.RefundOrderParam) error
}

type Deps struct {
	Module Module
	Out    *bufio.Writer
}

type CLI struct {
	Deps
	commandList []command
}

func NewCLI(d Deps) CLI {
	return CLI{
		Deps: d,
		commandList: []command{
			{
				name:        help,
				usage:       help,
				description: "справка",
			},
			{
				name:        refundOrder,
				usage:       "refund --id=<id заказа> --user=<id>",
				description: "На вход принимается ID пользователя и ID заказа. Заказ может быть возвращен в течение двух дней с момента выдачи. Также необходимо проверить, что заказ выдавался с нашего ПВЗ.",
			},
			{
				name:        refundedOrder,
				usage:       "refunded --count=20 --offset=10",
				description: "получить список всех заказов, которые вернули клиенты: Метод должен выдавать список пагинированно.",
			},
			{
				name:        listOrder,
				usage:       "list --user=<id>",
				description: "На вход принимается ID пользователя как обязательный параметр и опциональные параметры. Параметры позволяют получать только последние N заказов или заказы клиента, находящиеся в нашем ПВЗ.",
			},
			{
				name:        issueOrder,
				usage:       "issue <id заказа 1> ... <id заказа N>",
				description: "Можно выдавать только те заказы, которые были приняты от курьера и чей срок хранения меньше текущей даты. Все ID заказов должны принадлежать только одному клиенту.",
			},
			{
				name:        returnOrder,
				usage:       "return --id=<id заказа>",
				description: "На вход принимается ID заказа. Метод должен удалять заказ из вашего файла. Можно вернуть только те заказы, у которых вышел срок хранения и если заказы находятся в пвз.",
			},
			{
				name:        deliveryOrder,
				usage:       "delivery --id=<id заказа> --user=<id> --exp=2h",
				description: "На вход принимается ID заказа, ID получателя и срок хранения. Заказ нельзя принять дважды. Если срок хранения в прошлом, приложение должно выдать ошибку. Список принятых заказов необходимо сохранять в файл. Формат файла остается на выбор автора.",
			},
			{
				name:        exit,
				usage:       exit,
				description: "завершить выполение",
			},
		},
	}
}

func (c CLI) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("command isn't set")
	}
	defer c.Out.Flush()

	commandName := args[0]
	switch commandName {
	case help:
		c.help()
		return nil
	case refundedOrder:
		return c.refundedOrder(args[1:])
	case refundOrder:
		return c.refundOrder(args[1:])
	case deliveryOrder:
		return c.deliveryOrder(args[1:])
	case issueOrder:
		return c.issueOrders(args[1:])
	case returnOrder:
		return c.returnOrder(args[1:])
	case listOrder:
		return c.listOrder(args[1:])
	case exit:
		return ErrExit
	}
	return fmt.Errorf("command isn't set")
}

func (c CLI) refundOrder(args []string) error {
	var (
		ID, userID string
	)

	fs := flag.NewFlagSet(deliveryOrder, flag.ContinueOnError)
	fs.StringVar(&userID, "user", "", "use --user=<id>")
	fs.StringVar(&ID, "id", "", "use --user=<id>")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if ID == "" {
		return errors.New("id is empty")
	}
	if userID == "" {
		return errors.New("user is empty")
	}

	return c.Module.RefundOrder(module.RefundOrderParam{
		ID:          ID,
		RecipientID: userID,
	})
}

func (c CLI) issueOrders(args []string) error {
	return c.Module.IssueOrders(args)
}

func (c CLI) returnOrder(args []string) error {
	var (
		ID string
	)

	fs := flag.NewFlagSet(deliveryOrder, flag.ContinueOnError)
	fs.StringVar(&ID, "id", "", "use --user=<id>")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if ID == "" {
		return errors.New("id is empty")
	}
	return c.Module.ReturnOrder(ID)
}

func (c CLI) deliveryOrder(args []string) error {
	var (
		ID, userID string
		exp        time.Duration
	)

	fs := flag.NewFlagSet(deliveryOrder, flag.ContinueOnError)
	fs.DurationVar(&exp, "exp", time.Duration(0), "use --exp=2h")
	fs.StringVar(&userID, "user", "", "use --user=<id>")
	fs.StringVar(&ID, "id", "", "use --user=<id>")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if exp == time.Duration(0) {
		return errors.New("exp is empty")
	}
	if ID == "" {
		return errors.New("id is empty")
	}
	if userID == "" {
		return errors.New("user is empty")
	}

	return c.Module.Delivery(module.Order{
		ID:             ID,
		RecipientID:    userID,
		ExpirationDate: time.Now().Add(exp),
	})
}

func (c CLI) listOrder(args []string) error {
	var (
		userID string
		count  int
	)

	fs := flag.NewFlagSet(deliveryOrder, flag.ContinueOnError)
	fs.StringVar(&userID, "user", "", "use --user=<id>")
	fs.IntVar(&count, "count", math.MaxInt, "use --count=<int>")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if userID == "" {
		return errors.New("user is empty")
	}

	list, err := c.Module.ListOrder(userID, count)
	if err != nil {
		return err
	}
	for _, order := range list {
		fmt.Fprintln(c.Out, order)
	}
	return nil
}

func (c CLI) refundedOrder(args []string) error {
	var (
		count  int
		offset int
	)

	fs := flag.NewFlagSet(deliveryOrder, flag.ContinueOnError)
	fs.IntVar(&count, "count", math.MaxInt, "use --count=<int>")
	fs.IntVar(&offset, "offset", 0, "use --offset=<int>")
	if err := fs.Parse(args); err != nil {
		return err
	}

	list, err := c.Module.RefundedOrder(module.RefundedOrderParam{Offset: offset, Count: count})
	if err != nil {
		return err
	}
	for _, order := range list {
		fmt.Fprintln(c.Out, order)
	}
	return nil
}

func (c CLI) help() {
	fmt.Fprintln(c.Out, "command list:")
	for _, cmd := range c.commandList {
		fmt.Fprintln(c.Out, cmd)
	}
}
