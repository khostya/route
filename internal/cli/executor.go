package cli

import (
	"flag"
	"homework-1/internal/model"
	"homework-1/internal/service"
	"math"
	"strings"
	"time"
)

type Executor struct {
	service orderService
}

func (e Executor) refundOrder(args []string) string {
	var (
		ID, userID string
	)

	fs := flag.NewFlagSet(refundOrder, flag.ContinueOnError)
	fs.StringVar(&userID, "user", "", refundOrderUsage)
	fs.StringVar(&ID, "id", "", refundOrderUsage)
	if err := fs.Parse(args); err != nil {
		return err.Error()
	}

	if ID == "" {
		return ErrIdIsEmpty.Error()
	}
	if userID == "" {
		return ErrUserIsEmpty.Error()
	}

	err := e.service.RefundOrder(service.RefundOrderParam{
		ID:          ID,
		RecipientID: userID,
	})
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e Executor) issueOrders(args []string) string {
	err := e.service.IssueOrders(args)
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e Executor) returnOrder(args []string) string {
	var (
		ID string
	)

	fs := flag.NewFlagSet(deliverOrder, flag.ContinueOnError)
	fs.StringVar(&ID, "id", "", returnOrderUsage)
	if err := fs.Parse(args); err != nil {
		return err.Error()
	}

	if ID == "" {
		return ErrIdIsEmpty.Error()
	}

	err := e.service.ReturnOrder(ID)
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e Executor) deliverOrder(args []string) string {
	var (
		ID, userID string
		expString  string
	)

	fs := flag.NewFlagSet(deliverOrder, flag.ContinueOnError)
	fs.StringVar(&expString, "exp", "", deliverOrderUsage)
	fs.StringVar(&userID, "user", "", deliverOrderUsage)
	fs.StringVar(&ID, "id", "", deliverOrderUsage)
	if err := fs.Parse(args); err != nil {
		return err.Error()
	}

	if expString == "" {
		return ErrExpIsEmpty.Error()
	}
	if ID == "" {
		return ErrIdIsEmpty.Error()
	}
	if userID == "" {
		return ErrUserIsEmpty.Error()
	}

	exp, err := time.Parse(time.RFC3339, expString)
	if err != nil {
		return err.Error()
	}
	if exp.Before(time.Now()) {
		return ErrExpIsNotValid.Error()
	}

	err = e.service.Deliver(service.DeliverOrderParam{
		ID:             ID,
		RecipientID:    userID,
		ExpirationDate: exp,
	})
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e Executor) listOrder(args []string) string {
	var (
		userID string
		size   int
	)

	fs := flag.NewFlagSet(listOrder, flag.ContinueOnError)
	fs.StringVar(&userID, "user", "", listOrderUsage)
	fs.IntVar(&size, "size", math.MaxInt, listOrderUsage)

	if err := fs.Parse(args); err != nil {
		return err.Error()
	}

	if userID == "" {
		return ErrUserIsEmpty.Error()
	}

	list, err := e.service.ListOrder(userID, size)
	if err != nil {
		return err.Error()
	}
	return e.stringOrders(list)
}

func (e Executor) listRefunded(args []string) string {
	var (
		size int
		page int
	)

	fs := flag.NewFlagSet(deliverOrder, flag.ContinueOnError)
	fs.IntVar(&size, "size", math.MaxInt, deliverOrderUsage)
	fs.IntVar(&page, "page", 1, deliverOrderUsage)
	if err := fs.Parse(args); err != nil {
		return err.Error()
	}

	if page <= 0 {
		return ErrPageIsNotValid.Error()
	}

	list, err := e.service.ListRefunded(service.RefundedOrderParam{Page: page - 1, Size: size})
	if err != nil {
		return err.Error()
	}
	return e.stringOrders(list)
}

func (e Executor) stringOrders(orders []model.Order) string {
	var builder strings.Builder

	for _, order := range orders {
		builder.WriteString(order.String())
	}

	return builder.String()
}
