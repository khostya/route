package cli

import (
	"context"
	"flag"
	"github.com/shopspring/decimal"
	"homework/internal/dto"
	"homework/internal/model"
	"homework/internal/model/wrapper"
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

func (e executor) parseRefundOrder(args []string) (dto.RefundOrderParam, error) {
	var param dto.RefundOrderParam

	fs := flag.NewFlagSet(refundOrder, flag.ContinueOnError)
	fs.StringVar(&param.RecipientID, userIdParam, "", userIdParamUsage)
	fs.StringVar(&param.ID, orderIdParam, "", orderIdParamUsage)
	if err := fs.Parse(args); err != nil {
		return dto.RefundOrderParam{}, err
	}

	if param.ID == "" {
		return dto.RefundOrderParam{}, ErrIdIsEmpty
	}
	if param.RecipientID == "" {
		return dto.RefundOrderParam{}, ErrUserIsEmpty
	}

	return param, nil
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

func (e executor) parseDeliverOrder(args []string) (dto.DeliverOrderParam, error) {
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
		return dto.DeliverOrderParam{}, err
	}

	if expString == "" {
		return dto.DeliverOrderParam{}, ErrExpIsEmpty
	}
	if ID == "" {
		return dto.DeliverOrderParam{}, ErrIdIsEmpty
	}
	if userID == "" {
		return dto.DeliverOrderParam{}, ErrUserIsEmpty
	}
	if weightInKg <= 0 {
		return dto.DeliverOrderParam{}, ErrWeightInKgInNotValid
	}
	if priceInRubFloat64 < 0 {
		return dto.DeliverOrderParam{}, ErrPriceInRubIsNotValid
	}

	priceInRub := wrapper.PriceInRub(decimal.NewFromFloat(priceInRubFloat64))
	wrapperIsEmpty := wrapperType == ""
	if !wrapperIsEmpty && !slices.Contains(wrapper.GetAllWrapperTypes(), wrapper.WrapperType(wrapperType)) {
		return dto.DeliverOrderParam{}, ErrWrapperIsNotValid
	}

	exp, err := time.Parse(model.TimeFormat, expString)
	if err != nil {
		return dto.DeliverOrderParam{}, err
	}

	wrapper, err := wrapper.NewDefaultWrapper(wrapper.WrapperType(wrapperType))
	if !wrapperIsEmpty && err != nil {
		return dto.DeliverOrderParam{}, err
	}

	return dto.DeliverOrderParam{
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

func (e executor) parseListOrders(args []string) (dto.ListUserOrdersParam, error) {
	var param dto.ListUserOrdersParam

	fs := flag.NewFlagSet(listOrders, flag.ContinueOnError)
	fs.StringVar(&param.UserId, userIdParam, "", userIdParamUsage)
	fs.UintVar(&param.Count, sizeParam, math.MaxUint, sizeParamUsage)

	if err := fs.Parse(args); err != nil {
		return param, err
	}

	if param.UserId == "" {
		return param, ErrUserIsEmpty
	}
	if param.Count <= 0 {
		return param, ErrSizeIsNotValid
	}

	return param, nil
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

func (e executor) parseListRefunded(args []string) (dto.PageParam, error) {
	var param dto.PageParam

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
