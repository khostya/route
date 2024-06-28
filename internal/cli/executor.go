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

type executor struct {
	service orderService
}

func newExecutor(service orderService) executor {
	return executor{service: service}
}

func (e executor) refundOrder(ctx context.Context, args []string) string {
	param, err := e.parseRefundOrder(args)
	if err != nil {
		return err.Error()
	}

	err = e.service.RefundOrder(ctx, param)
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e executor) parseRefundOrder(args []string) (service.RefundOrderParam, error) {
	var (
		ID, userID string
	)

	fs := flag.NewFlagSet(refundOrder, flag.ContinueOnError)
	fs.StringVar(&userID, userIdParam, "", userIdParamUsage)
	fs.StringVar(&ID, orderIdParam, "", orderIdParamUsage)
	if err := fs.Parse(args); err != nil {
		return service.RefundOrderParam{}, err
	}

	if ID == "" {
		return service.RefundOrderParam{}, ErrIdIsEmpty
	}
	if userID == "" {
		return service.RefundOrderParam{}, ErrUserIsEmpty
	}

	return service.RefundOrderParam{
		ID:          ID,
		RecipientID: userID,
	}, nil
}

func (e executor) issueOrders(ctx context.Context, args []string) string {
	err := e.service.IssueOrders(ctx, args)
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e executor) returnOrder(ctx context.Context, args []string) string {
	id, err := e.parseReturnOrder(args)
	if err != nil {
		return err.Error()
	}

	err = e.service.ReturnOrder(ctx, id)
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e executor) parseReturnOrder(args []string) (string, error) {
	var (
		ID string
	)

	fs := flag.NewFlagSet(deliverOrder, flag.ContinueOnError)
	fs.StringVar(&ID, orderIdParam, "", orderIdParamUsage)
	err := fs.Parse(args)
	if ID == "" {
		return "", ErrIdIsEmpty
	}
	return ID, err
}

func (e executor) deliverOrder(ctx context.Context, args []string) string {
	param, err := e.parseDeliverOrder(args)
	if err != nil {
		return err.Error()
	}

	err = e.service.Deliver(ctx, param)
	if err == nil {
		return ""
	}
	return err.Error()
}

func (e executor) parseDeliverOrder(args []string) (service.DeliverOrderParam, error) {
	var (
		ID, userID        string
		expString         string
		wrapperType       string
		weightInKg        float64
		priceInRubFloat64 float64
	)

	fs := flag.NewFlagSet(deliverOrder, flag.ContinueOnError)
	fs.StringVar(&expString, expParam, "", expParamUsage)
	fs.StringVar(&userID, userIdParam, "", userIdParamUsage)
	fs.StringVar(&ID, orderIdParam, "", orderIdParamUsage)
	fs.StringVar(&wrapperType, wrapperParam, "", wrapperParamUsage)
	fs.Float64Var(&weightInKg, weightInKgParam, 0, weightInKgUsage)
	fs.Float64Var(&priceInRubFloat64, priceInRubParam, 0, priceInRubParamUsage)
	if err := fs.Parse(args); err != nil {
		return service.DeliverOrderParam{}, err
	}

	if expString == "" {
		return service.DeliverOrderParam{}, ErrExpIsEmpty
	}
	if ID == "" {
		return service.DeliverOrderParam{}, ErrIdIsEmpty
	}
	if userID == "" {
		return service.DeliverOrderParam{}, ErrUserIsEmpty
	}
	if weightInKg <= 0 {
		return service.DeliverOrderParam{}, ErrWeightInKgInNotValid
	}
	if priceInRubFloat64 < 0 {
		return service.DeliverOrderParam{}, ErrPriceInRubIsNotValid
	}

	priceInRub := wrapper.PriceInRub(decimal.NewFromFloat(priceInRubFloat64))
	wrapperIsEmpty := wrapperType == ""
	if !wrapperIsEmpty && !slices.Contains(wrapper.GetAllWrapperTypes(), wrapper.WrapperType(wrapperType)) {
		return service.DeliverOrderParam{}, ErrWrapperIsNotValid
	}

	exp, err := time.Parse(model.TimeFormat, expString)
	if err != nil {
		return service.DeliverOrderParam{}, err
	}

	wrapper, err := wrapper.NewDefaultWrapper(wrapper.WrapperType(wrapperType))
	if !wrapperIsEmpty && err != nil {
		return service.DeliverOrderParam{}, err
	}

	return service.DeliverOrderParam{
		ID:             ID,
		RecipientID:    userID,
		ExpirationDate: exp,
		WeightInGram:   weightInKg * 1000,
		Wrapper:        wrapper,
		PriceInRub:     priceInRub,
	}, nil
}

func (e executor) listOrders(ctx context.Context, args []string) string {
	param, err := e.parseListOrders(args)
	if err != nil {
		return err.Error()
	}

	list, err := e.service.ListUserOrders(ctx, param)
	if err != nil {
		return err.Error()
	}

	return e.stringOrders(list)
}

func (e executor) parseListOrders(args []string) (service.ListUserOrdersParam, error) {
	var (
		userID string
		size   uint
	)

	fs := flag.NewFlagSet(listOrders, flag.ContinueOnError)
	fs.StringVar(&userID, userIdParam, "", userIdParamUsage)
	fs.UintVar(&size, sizeParam, math.MaxUint, sizeParamUsage)

	if err := fs.Parse(args); err != nil {
		return service.ListUserOrdersParam{}, err
	}

	if userID == "" {
		return service.ListUserOrdersParam{}, ErrUserIsEmpty
	}
	if size <= 0 {
		return service.ListUserOrdersParam{}, ErrSizeIsNotValid
	}

	return service.ListUserOrdersParam{UserId: userID, Count: size}, nil
}

func (e executor) listRefunded(ctx context.Context, args []string) string {
	param, err := e.parseListRefunded(args)
	if err != nil {
		return err.Error()
	}

	list, err := e.service.RefundedOrders(ctx, param)
	if err != nil {
		return err.Error()
	}
	return e.stringOrders(list)
}

func (e executor) parseListRefunded(args []string) (service.RefundedOrdersParam, error) {
	var param service.RefundedOrdersParam

	fs := flag.NewFlagSet(deliverOrder, flag.ContinueOnError)
	fs.UintVar(&param.Size, sizeParam, math.MaxUint, sizeParamUsage)
	fs.UintVar(&param.Page, pageParam, 1, pageParamUsage)
	if err := fs.Parse(args); err != nil {
		return param, err
	}

	if param.Page <= 0 {
		return param, ErrPageIsNotValid
	}
	if param.Size <= 0 {
		return param, ErrSizeIsNotValid
	}

	return param, nil
}

func (e executor) stringOrders(orders []model.Order) string {
	var builder strings.Builder

	for i := 0; i < len(orders)-1; i++ {
		builder.WriteString(orders[i].String() + "\n")
	}
	builder.WriteString(orders[len(orders)-1].String())

	return builder.String()
}
