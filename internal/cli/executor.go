package cli

import (
	"flag"
	"homework/internal/model"
	"homework/internal/service"
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
	fs.StringVar(&userID, userIdParam, "", userIdParamUsage)
	fs.StringVar(&ID, orderIdParam, "", orderIdParamUsage)
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
	fs.StringVar(&ID, orderIdParam, "", orderIdParamUsage)
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
	fs.StringVar(&expString, expParam, "", expString)
	fs.StringVar(&userID, userIdParam, "", userID)
	fs.StringVar(&ID, orderIdParam, "", orderIdParamUsage)
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

func (e Executor) listOrders(args []string) string {
	var (
		userID string
		size   int
	)

	fs := flag.NewFlagSet(listOrders, flag.ContinueOnError)
	fs.StringVar(&userID, userIdParam, "", userIdParamUsage)
	fs.IntVar(&size, sizeParam, math.MaxInt, sizeParamUsage)

	if err := fs.Parse(args); err != nil {
		return err.Error()
	}

	if userID == "" {
		return ErrUserIsEmpty.Error()
	}
	if size <= 0 {
		return ErrSizeIsNotValid.Error()
	}

	list, err := e.service.ListUserOrders(userID, size)
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
	fs.IntVar(&size, sizeParam, math.MaxInt, sizeParamUsage)
	fs.IntVar(&page, pageParam, 1, pageParamUsage)
	if err := fs.Parse(args); err != nil {
		return err.Error()
	}

	if page <= 0 {
		return ErrPageIsNotValid.Error()
	}
	if size <= 0 {
		return ErrSizeIsNotValid.Error()
	}

	list, err := e.service.RefundedOrders(service.RefundedOrdersParam{Page: page - 1, Size: size})
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
