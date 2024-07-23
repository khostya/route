// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v4.25.1
// source: order/v1/order.proto

package order

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OrderStatus int32

const (
	OrderStatus_ORDER_STATUS_ANY       OrderStatus = 0
	OrderStatus_ORDER_STATUS_DELIVERED OrderStatus = 1
	OrderStatus_ORDER_STATUS_ISSUED    OrderStatus = 2
	OrderStatus_ORDER_STATUS_REFUNDED  OrderStatus = 3
)

// Enum value maps for OrderStatus.
var (
	OrderStatus_name = map[int32]string{
		0: "ORDER_STATUS_ANY",
		1: "ORDER_STATUS_DELIVERED",
		2: "ORDER_STATUS_ISSUED",
		3: "ORDER_STATUS_REFUNDED",
	}
	OrderStatus_value = map[string]int32{
		"ORDER_STATUS_ANY":       0,
		"ORDER_STATUS_DELIVERED": 1,
		"ORDER_STATUS_ISSUED":    2,
		"ORDER_STATUS_REFUNDED":  3,
	}
)

func (x OrderStatus) Enum() *OrderStatus {
	p := new(OrderStatus)
	*p = x
	return p
}

func (x OrderStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrderStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_order_v1_order_proto_enumTypes[0].Descriptor()
}

func (OrderStatus) Type() protoreflect.EnumType {
	return &file_order_v1_order_proto_enumTypes[0]
}

func (x OrderStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrderStatus.Descriptor instead.
func (OrderStatus) EnumDescriptor() ([]byte, []int) {
	return file_order_v1_order_proto_rawDescGZIP(), []int{0}
}

type WrapperType int32

const (
	WrapperType_WRAPPER_TYPE_NONE    WrapperType = 0
	WrapperType_WRAPPER_TYPE_BOX     WrapperType = 1
	WrapperType_WRAPPER_TYPE_PACKAGE WrapperType = 2
	WrapperType_WRAPPER_TYPE_STRETCH WrapperType = 3
)

// Enum value maps for WrapperType.
var (
	WrapperType_name = map[int32]string{
		0: "WRAPPER_TYPE_NONE",
		1: "WRAPPER_TYPE_BOX",
		2: "WRAPPER_TYPE_PACKAGE",
		3: "WRAPPER_TYPE_STRETCH",
	}
	WrapperType_value = map[string]int32{
		"WRAPPER_TYPE_NONE":    0,
		"WRAPPER_TYPE_BOX":     1,
		"WRAPPER_TYPE_PACKAGE": 2,
		"WRAPPER_TYPE_STRETCH": 3,
	}
)

func (x WrapperType) Enum() *WrapperType {
	p := new(WrapperType)
	*p = x
	return p
}

func (x WrapperType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WrapperType) Descriptor() protoreflect.EnumDescriptor {
	return file_order_v1_order_proto_enumTypes[1].Descriptor()
}

func (WrapperType) Type() protoreflect.EnumType {
	return &file_order_v1_order_proto_enumTypes[1]
}

func (x WrapperType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WrapperType.Descriptor instead.
func (WrapperType) EnumDescriptor() ([]byte, []int) {
	return file_order_v1_order_proto_rawDescGZIP(), []int{1}
}

type DeliverOrderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderID     string                 `protobuf:"bytes,1,opt,name=orderID,proto3" json:"orderID,omitempty"`
	UserID      string                 `protobuf:"bytes,2,opt,name=userID,proto3" json:"userID,omitempty"`
	Exp         *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=exp,proto3" json:"exp,omitempty"`
	WrapperType WrapperType            `protobuf:"varint,4,opt,name=wrapperType,proto3,enum=order.WrapperType" json:"wrapperType,omitempty"`
	WeightInKg  float32                `protobuf:"fixed32,5,opt,name=weightInKg,proto3" json:"weightInKg,omitempty"`
	PriceInRub  float32                `protobuf:"fixed32,6,opt,name=priceInRub,proto3" json:"priceInRub,omitempty"`
}

func (x *DeliverOrderRequest) Reset() {
	*x = DeliverOrderRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_v1_order_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliverOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliverOrderRequest) ProtoMessage() {}

func (x *DeliverOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_order_v1_order_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliverOrderRequest.ProtoReflect.Descriptor instead.
func (*DeliverOrderRequest) Descriptor() ([]byte, []int) {
	return file_order_v1_order_proto_rawDescGZIP(), []int{0}
}

func (x *DeliverOrderRequest) GetOrderID() string {
	if x != nil {
		return x.OrderID
	}
	return ""
}

func (x *DeliverOrderRequest) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *DeliverOrderRequest) GetExp() *timestamppb.Timestamp {
	if x != nil {
		return x.Exp
	}
	return nil
}

func (x *DeliverOrderRequest) GetWrapperType() WrapperType {
	if x != nil {
		return x.WrapperType
	}
	return WrapperType_WRAPPER_TYPE_NONE
}

func (x *DeliverOrderRequest) GetWeightInKg() float32 {
	if x != nil {
		return x.WeightInKg
	}
	return 0
}

func (x *DeliverOrderRequest) GetPriceInRub() float32 {
	if x != nil {
		return x.PriceInRub
	}
	return 0
}

type ReturnOrderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ReturnOrderRequest) Reset() {
	*x = ReturnOrderRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_v1_order_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReturnOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReturnOrderRequest) ProtoMessage() {}

func (x *ReturnOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_order_v1_order_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReturnOrderRequest.ProtoReflect.Descriptor instead.
func (*ReturnOrderRequest) Descriptor() ([]byte, []int) {
	return file_order_v1_order_proto_rawDescGZIP(), []int{1}
}

func (x *ReturnOrderRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type IssueOrdersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids []string `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`
}

func (x *IssueOrdersRequest) Reset() {
	*x = IssueOrdersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_v1_order_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IssueOrdersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IssueOrdersRequest) ProtoMessage() {}

func (x *IssueOrdersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_order_v1_order_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IssueOrdersRequest.ProtoReflect.Descriptor instead.
func (*IssueOrdersRequest) Descriptor() ([]byte, []int) {
	return file_order_v1_order_proto_rawDescGZIP(), []int{2}
}

func (x *IssueOrdersRequest) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

type RefundOrderRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID  string `protobuf:"bytes,1,opt,name=userID,proto3" json:"userID,omitempty"`
	OrderID string `protobuf:"bytes,2,opt,name=orderID,proto3" json:"orderID,omitempty"`
}

func (x *RefundOrderRequest) Reset() {
	*x = RefundOrderRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_v1_order_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefundOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefundOrderRequest) ProtoMessage() {}

func (x *RefundOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_order_v1_order_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefundOrderRequest.ProtoReflect.Descriptor instead.
func (*RefundOrderRequest) Descriptor() ([]byte, []int) {
	return file_order_v1_order_proto_rawDescGZIP(), []int{3}
}

func (x *RefundOrderRequest) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

func (x *RefundOrderRequest) GetOrderID() string {
	if x != nil {
		return x.OrderID
	}
	return ""
}

type ListOrdersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID *string      `protobuf:"bytes,1,opt,name=userID,proto3,oneof" json:"userID,omitempty"`
	Size   *uint32      `protobuf:"varint,2,opt,name=size,proto3,oneof" json:"size,omitempty"`
	Page   *uint32      `protobuf:"varint,3,opt,name=page,proto3,oneof" json:"page,omitempty"`
	Status *OrderStatus `protobuf:"varint,4,opt,name=status,proto3,enum=order.OrderStatus,oneof" json:"status,omitempty"`
}

func (x *ListOrdersRequest) Reset() {
	*x = ListOrdersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_v1_order_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListOrdersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListOrdersRequest) ProtoMessage() {}

func (x *ListOrdersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_order_v1_order_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListOrdersRequest.ProtoReflect.Descriptor instead.
func (*ListOrdersRequest) Descriptor() ([]byte, []int) {
	return file_order_v1_order_proto_rawDescGZIP(), []int{4}
}

func (x *ListOrdersRequest) GetUserID() string {
	if x != nil && x.UserID != nil {
		return *x.UserID
	}
	return ""
}

func (x *ListOrdersRequest) GetSize() uint32 {
	if x != nil && x.Size != nil {
		return *x.Size
	}
	return 0
}

func (x *ListOrdersRequest) GetPage() uint32 {
	if x != nil && x.Page != nil {
		return *x.Page
	}
	return 0
}

func (x *ListOrdersRequest) GetStatus() OrderStatus {
	if x != nil && x.Status != nil {
		return *x.Status
	}
	return OrderStatus_ORDER_STATUS_ANY
}

type ListOrdersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Orders []*ListOrdersResponse_Order `protobuf:"bytes,1,rep,name=orders,proto3" json:"orders,omitempty"`
}

func (x *ListOrdersResponse) Reset() {
	*x = ListOrdersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_v1_order_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListOrdersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListOrdersResponse) ProtoMessage() {}

func (x *ListOrdersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_order_v1_order_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListOrdersResponse.ProtoReflect.Descriptor instead.
func (*ListOrdersResponse) Descriptor() ([]byte, []int) {
	return file_order_v1_order_proto_rawDescGZIP(), []int{5}
}

func (x *ListOrdersResponse) GetOrders() []*ListOrdersResponse_Order {
	if x != nil {
		return x.Orders
	}
	return nil
}

type ListOrdersResponse_Order struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	RecipientID string      `protobuf:"bytes,2,opt,name=recipientID,proto3" json:"recipientID,omitempty"`
	Status      OrderStatus `protobuf:"varint,3,opt,name=status,proto3,enum=order.OrderStatus" json:"status,omitempty"`
}

func (x *ListOrdersResponse_Order) Reset() {
	*x = ListOrdersResponse_Order{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_v1_order_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListOrdersResponse_Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListOrdersResponse_Order) ProtoMessage() {}

func (x *ListOrdersResponse_Order) ProtoReflect() protoreflect.Message {
	mi := &file_order_v1_order_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListOrdersResponse_Order.ProtoReflect.Descriptor instead.
func (*ListOrdersResponse_Order) Descriptor() ([]byte, []int) {
	return file_order_v1_order_proto_rawDescGZIP(), []int{5, 0}
}

func (x *ListOrdersResponse_Order) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ListOrdersResponse_Order) GetRecipientID() string {
	if x != nil {
		return x.RecipientID
	}
	return ""
}

func (x *ListOrdersResponse_Order) GetStatus() OrderStatus {
	if x != nil {
		return x.Status
	}
	return OrderStatus_ORDER_STATUS_ANY
}

var File_order_v1_order_proto protoreflect.FileDescriptor

var file_order_v1_order_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65,
	0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69,
	0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xb5, 0x02, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x07, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe0, 0x41,
	0x02, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49,
	0x44, 0x12, 0x22, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x0a, 0xe0, 0x41, 0x02, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x3b, 0x0a, 0x03, 0x65, 0x78, 0x70, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x0d,
	0xe0, 0x41, 0x02, 0xfa, 0x42, 0x07, 0xb2, 0x01, 0x04, 0x08, 0x01, 0x40, 0x01, 0x52, 0x03, 0x65,
	0x78, 0x70, 0x12, 0x39, 0x0a, 0x0b, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e,
	0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x42, 0x03, 0xe0, 0x41, 0x02,
	0x52, 0x0b, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2d, 0x0a,
	0x0a, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x4b, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x02, 0x42, 0x0d, 0xe0, 0x41, 0x02, 0xfa, 0x42, 0x07, 0x0a, 0x05, 0x25, 0x00, 0x00, 0x00, 0x00,
	0x52, 0x0a, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x4b, 0x67, 0x12, 0x2d, 0x0a, 0x0a,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x52, 0x75, 0x62, 0x18, 0x06, 0x20, 0x01, 0x28, 0x02,
	0x42, 0x0d, 0xe0, 0x41, 0x02, 0xfa, 0x42, 0x07, 0x0a, 0x05, 0x2d, 0x00, 0x00, 0x00, 0x00, 0x52,
	0x0a, 0x70, 0x72, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x52, 0x75, 0x62, 0x22, 0x30, 0x0a, 0x12, 0x52,
	0x65, 0x74, 0x75, 0x72, 0x6e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1a, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xe0,
	0x41, 0x02, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x02, 0x69, 0x64, 0x22, 0x39, 0x0a,
	0x12, 0x49, 0x73, 0x73, 0x75, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x42, 0x11, 0xe0, 0x41, 0x02, 0xfa, 0x42, 0x0b, 0x92, 0x01, 0x08, 0x08, 0x01, 0x22, 0x04, 0x72,
	0x02, 0x10, 0x01, 0x52, 0x03, 0x69, 0x64, 0x73, 0x22, 0x5e, 0x0a, 0x12, 0x52, 0x65, 0x66, 0x75,
	0x6e, 0x64, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22,
	0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a,
	0xe0, 0x41, 0x02, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x44, 0x12, 0x24, 0x0a, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x0a, 0xe0, 0x41, 0x02, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52,
	0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x44, 0x22, 0xdf, 0x01, 0x0a, 0x11, 0x4c, 0x69, 0x73,
	0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27,
	0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a,
	0xe0, 0x41, 0x02, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x48, 0x00, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x44, 0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x0a, 0xe0, 0x41, 0x02, 0xfa, 0x42, 0x04, 0x2a, 0x02, 0x20,
	0x00, 0x48, 0x01, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x0a, 0xe0, 0x41, 0x02, 0xfa,
	0x42, 0x04, 0x2a, 0x02, 0x20, 0x00, 0x48, 0x02, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x2f, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x12, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x48, 0x03, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x88,
	0x01, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x42, 0x07, 0x0a,
	0x05, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x42,
	0x09, 0x0a, 0x07, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xb4, 0x01, 0x0a, 0x12, 0x4c,
	0x69, 0x73, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x37, 0x0a, 0x06, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1f, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x52, 0x06, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x1a, 0x65, 0x0a, 0x05, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74,
	0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x2a, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x2a, 0x73, 0x0a, 0x0b, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x14, 0x0a, 0x10, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53,
	0x5f, 0x41, 0x4e, 0x59, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x45, 0x4c, 0x49, 0x56, 0x45, 0x52, 0x45, 0x44,
	0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54,
	0x55, 0x53, 0x5f, 0x49, 0x53, 0x53, 0x55, 0x45, 0x44, 0x10, 0x02, 0x12, 0x19, 0x0a, 0x15, 0x4f,
	0x52, 0x44, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x52, 0x45, 0x46, 0x55,
	0x4e, 0x44, 0x45, 0x44, 0x10, 0x03, 0x2a, 0x6e, 0x0a, 0x0b, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65,
	0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x11, 0x57, 0x52, 0x41, 0x50, 0x50, 0x45, 0x52,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10,
	0x57, 0x52, 0x41, 0x50, 0x50, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x42, 0x4f, 0x58,
	0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x57, 0x52, 0x41, 0x50, 0x50, 0x45, 0x52, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x50, 0x41, 0x43, 0x4b, 0x41, 0x47, 0x45, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14,
	0x57, 0x52, 0x41, 0x50, 0x50, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x54, 0x52,
	0x45, 0x54, 0x43, 0x48, 0x10, 0x03, 0x32, 0x99, 0x05, 0x0a, 0x05, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x12, 0x69, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x12, 0x1a, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x25, 0x92, 0x41, 0x07, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x3a, 0x01, 0x2a, 0x22, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x2f, 0x64, 0x65, 0x6c, 0x76, 0x65, 0x72, 0x12, 0x67, 0x0a, 0x0b, 0x52,
	0x65, 0x74, 0x75, 0x72, 0x6e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x2e, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x25, 0x92,
	0x41, 0x07, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x3a,
	0x01, 0x2a, 0x32, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2f, 0x72, 0x65,
	0x74, 0x75, 0x72, 0x6e, 0x12, 0x66, 0x0a, 0x0b, 0x49, 0x73, 0x73, 0x75, 0x65, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x73, 0x12, 0x19, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x49, 0x73, 0x73, 0x75,
	0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x24, 0x92, 0x41, 0x07, 0x0a, 0x05, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x3a, 0x01, 0x2a, 0x32, 0x0f, 0x2f, 0x76, 0x31,
	0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2f, 0x69, 0x73, 0x73, 0x75, 0x65, 0x12, 0x67, 0x0a, 0x0b,
	0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x2e, 0x52, 0x65, 0x66, 0x75, 0x6e, 0x64, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x25,
	0x92, 0x41, 0x07, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15,
	0x3a, 0x01, 0x2a, 0x32, 0x10, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2f, 0x72,
	0x65, 0x66, 0x75, 0x6e, 0x64, 0x12, 0x5f, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x73, 0x12, 0x18, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1c, 0x92, 0x41, 0x07, 0x0a, 0x05, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x76, 0x31, 0x2f,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x73, 0x1a, 0x89, 0x01, 0x92, 0x41, 0x85, 0x01, 0x0a, 0x05, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x12, 0x15, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x20, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x20, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x1a, 0x65, 0x0a, 0x20, 0x46,
	0x69, 0x6e, 0x64, 0x20, 0x6f, 0x75, 0x74, 0x20, 0x6d, 0x6f, 0x72, 0x65, 0x20, 0x61, 0x62, 0x6f,
	0x75, 0x74, 0x20, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x12,
	0x41, 0x68, 0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x65, 0x63, 0x6f, 0x73, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x3f,
	0x74, 0x61, 0x62, 0x3d, 0x72, 0x65, 0x61, 0x64, 0x6d, 0x65, 0x2d, 0x6f, 0x76, 0x2d, 0x66, 0x69,
	0x6c, 0x65, 0x42, 0x39, 0x92, 0x41, 0x17, 0x12, 0x15, 0x0a, 0x0e, 0x6f, 0x7a, 0x6f, 0x6e, 0x20,
	0x72, 0x6f, 0x75, 0x74, 0x65, 0x20, 0x32, 0x35, 0x36, 0x32, 0x03, 0x31, 0x2e, 0x30, 0x5a, 0x1d,
	0x68, 0x6f, 0x6d, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x3b, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_order_v1_order_proto_rawDescOnce sync.Once
	file_order_v1_order_proto_rawDescData = file_order_v1_order_proto_rawDesc
)

func file_order_v1_order_proto_rawDescGZIP() []byte {
	file_order_v1_order_proto_rawDescOnce.Do(func() {
		file_order_v1_order_proto_rawDescData = protoimpl.X.CompressGZIP(file_order_v1_order_proto_rawDescData)
	})
	return file_order_v1_order_proto_rawDescData
}

var file_order_v1_order_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_order_v1_order_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_order_v1_order_proto_goTypes = []any{
	(OrderStatus)(0),                 // 0: order.OrderStatus
	(WrapperType)(0),                 // 1: order.WrapperType
	(*DeliverOrderRequest)(nil),      // 2: order.DeliverOrderRequest
	(*ReturnOrderRequest)(nil),       // 3: order.ReturnOrderRequest
	(*IssueOrdersRequest)(nil),       // 4: order.IssueOrdersRequest
	(*RefundOrderRequest)(nil),       // 5: order.RefundOrderRequest
	(*ListOrdersRequest)(nil),        // 6: order.ListOrdersRequest
	(*ListOrdersResponse)(nil),       // 7: order.ListOrdersResponse
	(*ListOrdersResponse_Order)(nil), // 8: order.ListOrdersResponse.Order
	(*timestamppb.Timestamp)(nil),    // 9: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),            // 10: google.protobuf.Empty
}
var file_order_v1_order_proto_depIdxs = []int32{
	9,  // 0: order.DeliverOrderRequest.exp:type_name -> google.protobuf.Timestamp
	1,  // 1: order.DeliverOrderRequest.wrapperType:type_name -> order.WrapperType
	0,  // 2: order.ListOrdersRequest.status:type_name -> order.OrderStatus
	8,  // 3: order.ListOrdersResponse.orders:type_name -> order.ListOrdersResponse.Order
	0,  // 4: order.ListOrdersResponse.Order.status:type_name -> order.OrderStatus
	2,  // 5: order.Order.DeliverOrder:input_type -> order.DeliverOrderRequest
	3,  // 6: order.Order.ReturnOrder:input_type -> order.ReturnOrderRequest
	4,  // 7: order.Order.IssueOrders:input_type -> order.IssueOrdersRequest
	5,  // 8: order.Order.RefundOrder:input_type -> order.RefundOrderRequest
	6,  // 9: order.Order.ListOrders:input_type -> order.ListOrdersRequest
	10, // 10: order.Order.DeliverOrder:output_type -> google.protobuf.Empty
	10, // 11: order.Order.ReturnOrder:output_type -> google.protobuf.Empty
	10, // 12: order.Order.IssueOrders:output_type -> google.protobuf.Empty
	10, // 13: order.Order.RefundOrder:output_type -> google.protobuf.Empty
	7,  // 14: order.Order.ListOrders:output_type -> order.ListOrdersResponse
	10, // [10:15] is the sub-list for method output_type
	5,  // [5:10] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_order_v1_order_proto_init() }
func file_order_v1_order_proto_init() {
	if File_order_v1_order_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_order_v1_order_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*DeliverOrderRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_order_v1_order_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ReturnOrderRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_order_v1_order_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*IssueOrdersRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_order_v1_order_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*RefundOrderRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_order_v1_order_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*ListOrdersRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_order_v1_order_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*ListOrdersResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_order_v1_order_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*ListOrdersResponse_Order); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_order_v1_order_proto_msgTypes[4].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_order_v1_order_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_order_v1_order_proto_goTypes,
		DependencyIndexes: file_order_v1_order_proto_depIdxs,
		EnumInfos:         file_order_v1_order_proto_enumTypes,
		MessageInfos:      file_order_v1_order_proto_msgTypes,
	}.Build()
	File_order_v1_order_proto = out.File
	file_order_v1_order_proto_rawDesc = nil
	file_order_v1_order_proto_goTypes = nil
	file_order_v1_order_proto_depIdxs = nil
}
