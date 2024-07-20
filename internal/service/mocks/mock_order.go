// Code generated by MockGen. DO NOT EDIT.
// Source: ./mocks/order.go
//
// Generated by this command:
//
//	mockgen -source ./mocks/order.go -destination=./mocks/mock_order.go -package=mock_service
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	dto "homework/internal/dto"
	model "homework/internal/model"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockorderService is a mock of orderService interface.
type MockorderService struct {
	ctrl     *gomock.Controller
	recorder *MockorderServiceMockRecorder
}

// MockorderServiceMockRecorder is the mock recorder for MockorderService.
type MockorderServiceMockRecorder struct {
	mock *MockorderService
}

// NewMockorderService creates a new mock instance.
func NewMockorderService(ctrl *gomock.Controller) *MockorderService {
	mock := &MockorderService{ctrl: ctrl}
	mock.recorder = &MockorderServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockorderService) EXPECT() *MockorderServiceMockRecorder {
	return m.recorder
}

// Deliver mocks base method.
func (m *MockorderService) Deliver(ctx context.Context, order dto.DeliverOrderParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deliver", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deliver indicates an expected call of Deliver.
func (mr *MockorderServiceMockRecorder) Deliver(ctx, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deliver", reflect.TypeOf((*MockorderService)(nil).Deliver), ctx, order)
}

// IssueOrders mocks base method.
func (m *MockorderService) IssueOrders(ctx context.Context, ids []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IssueOrders", ctx, ids)
	ret0, _ := ret[0].(error)
	return ret0
}

// IssueOrders indicates an expected call of IssueOrders.
func (mr *MockorderServiceMockRecorder) IssueOrders(ctx, ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IssueOrders", reflect.TypeOf((*MockorderService)(nil).IssueOrders), ctx, ids)
}

// ListUserOrders mocks base method.
func (m *MockorderService) ListUserOrders(ctx context.Context, param dto.ListUserOrdersParam) ([]model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserOrders", ctx, param)
	ret0, _ := ret[0].([]model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUserOrders indicates an expected call of ListUserOrders.
func (mr *MockorderServiceMockRecorder) ListUserOrders(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserOrders", reflect.TypeOf((*MockorderService)(nil).ListUserOrders), ctx, param)
}

// RefundOrder mocks base method.
func (m *MockorderService) RefundOrder(ctx context.Context, param dto.RefundOrderParam) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefundOrder", ctx, param)
	ret0, _ := ret[0].(error)
	return ret0
}

// RefundOrder indicates an expected call of RefundOrder.
func (mr *MockorderServiceMockRecorder) RefundOrder(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefundOrder", reflect.TypeOf((*MockorderService)(nil).RefundOrder), ctx, param)
}

// RefundedOrders mocks base method.
func (m *MockorderService) RefundedOrders(ctx context.Context, param dto.PageParam) ([]model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefundedOrders", ctx, param)
	ret0, _ := ret[0].([]model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefundedOrders indicates an expected call of RefundedOrders.
func (mr *MockorderServiceMockRecorder) RefundedOrders(ctx, param any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefundedOrders", reflect.TypeOf((*MockorderService)(nil).RefundedOrders), ctx, param)
}

// ReturnOrder mocks base method.
func (m *MockorderService) ReturnOrder(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReturnOrder", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReturnOrder indicates an expected call of ReturnOrder.
func (mr *MockorderServiceMockRecorder) ReturnOrder(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReturnOrder", reflect.TypeOf((*MockorderService)(nil).ReturnOrder), ctx, id)
}
