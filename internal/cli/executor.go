package cli

import (
	"context"
	"flag"
	"github.com/shopspring/decimal"
	"homework/internal/model"
	"homework/internal/model/wrapper"
	"homework/internal/service"
	"math"
	"slices"
	"strings"
	"time"
)

type Executor struct {
	service orderService
}

func (e Executor) refundOrder(ctx context.Context, args []string) string {
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

	err := e.service.RefundOrder(ctx, service.RefundOrderParam{
		ID:          ID,
		RecipientID: userID,
	})
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e Executor) issueOrders(ctx context.Context, args []string) string {
	err := e.service.IssueOrders(ctx, args)
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e Executor) returnOrder(ctx context.Context, args []string) string {
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

	err := e.service.ReturnOrder(ctx, ID)
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e Executor) deliverOrder(ctx context.Context, args []string) string {
	var (
		ID, userID        string
		expString         string
		wrapperType       string
		weightInKg        float64
		priceInRubFloat64 float64
	)

	fs := flag.NewFlagSet(deliverOrder, flag.ContinueOnError)
	fs.StringVar(&expString, expParam, "", expString)
	fs.StringVar(&userID, userIdParam, "", userID)
	fs.StringVar(&ID, orderIdParam, "", orderIdParamUsage)
	fs.StringVar(&wrapperType, wrapperParam, "", wrapperParamUsage)
	fs.Float64Var(&weightInKg, weightInKgParam, 0, weightInKgUsage)
	fs.Float64Var(&priceInRubFloat64, priceInRubParam, 0, priceInRubParamUsage)
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
	if weightInKg <= 0 {
		return ErrWeightInKgInNotValid.Error()
	}
	if priceInRubFloat64 < 0 {
		return ErrPriceInRubIsNotValid.Error()
	}

	priceInRub := wrapper.PriceInRub(decimal.NewFromFloat(priceInRubFloat64))
	wrapperIsEmpty := wrapperType == ""
	if !wrapperIsEmpty && !slices.Contains(wrapper.GetAllWrapperTypes(), wrapper.WrapperType(wrapperType)) {
		return ErrWrapperIsNotValid.Error()
	}

	exp, err := time.Parse(model.TimeFormat, expString)
	if err != nil {
		return err.Error()
	}

	wrapper, err := wrapper.NewDefaultWrapper(wrapper.WrapperType(wrapperType))
	if !wrapperIsEmpty && err != nil {
		return err.Error()
	}

	err = e.service.Deliver(ctx, service.DeliverOrderParam{
		ID:             ID,
		RecipientID:    userID,
		ExpirationDate: exp,
		WeightInGram:   weightInKg * 1000,
		Wrapper:        wrapper,
		PriceInRub:     priceInRub,
	})
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e Executor) listOrders(ctx context.Context, args []string) string {
	var (
		userID string
		size   uint
	)

	fs := flag.NewFlagSet(listOrders, flag.ContinueOnError)
	fs.StringVar(&userID, userIdParam, "", userIdParamUsage)
	fs.UintVar(&size, sizeParam, math.MaxUint, sizeParamUsage)

	if err := fs.Parse(args); err != nil {
		return err.Error()
	}

	if userID == "" {
		return ErrUserIsEmpty.Error()
	}
	if size <= 0 {
		return ErrSizeIsNotValid.Error()
	}

	list, err := e.service.ListUserOrders(ctx, userID, size)
	if err != nil {
		return err.Error()
	}
	return e.stringOrders(list)
}

func (e Executor) listRefunded(ctx context.Context, args []string) string {
	var (
		size uint
		page uint
	)

	fs := flag.NewFlagSet(deliverOrder, flag.ContinueOnError)
	fs.UintVar(&size, sizeParam, math.MaxUint, sizeParamUsage)
	fs.UintVar(&page, pageParam, 1, pageParamUsage)
	if err := fs.Parse(args); err != nil {
		return err.Error()
	}

	if page <= 0 {
		return ErrPageIsNotValid.Error()
	}
	if size <= 0 {
		return ErrSizeIsNotValid.Error()
	}

	list, err := e.service.RefundedOrders(ctx, service.RefundedOrdersParam{Page: page - 1, Size: size})
	if err != nil {
		return err.Error()
	}
	return e.stringOrders(list)
}

func (e Executor) stringOrders(orders []model.Order) string {
	var builder strings.Builder

	for i := 0; i < len(orders)-1; i++ {
		builder.WriteString(orders[i].String() + "\n")
	}
	builder.WriteString(orders[len(orders)-1].String())

	return builder.String()
}
