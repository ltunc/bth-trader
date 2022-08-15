// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: api/proto/trader.proto

package bth

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TraderClient is the client API for Trader service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TraderClient interface {
	// AddOrder submits a new order on the exchange
	AddOrder(ctx context.Context, in *AddOrderRequest, opts ...grpc.CallOption) (*AddOrderResponse, error)
	// CancelOrder cancels an open order
	CancelOrder(ctx context.Context, in *CancelOrderRequest, opts ...grpc.CallOption) (*CancelOrderResponse, error)
	// OrderStatus request status of particular order
	OrderStatus(ctx context.Context, in *OrderStatusRequest, opts ...grpc.CallOption) (*OrderStatusResponse, error)
	// StreamOrders opens stream to receive update on order statuses as they become available
	StreamOrders(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Trader_StreamOrdersClient, error)
}

type traderClient struct {
	cc grpc.ClientConnInterface
}

func NewTraderClient(cc grpc.ClientConnInterface) TraderClient {
	return &traderClient{cc}
}

func (c *traderClient) AddOrder(ctx context.Context, in *AddOrderRequest, opts ...grpc.CallOption) (*AddOrderResponse, error) {
	out := new(AddOrderResponse)
	err := c.cc.Invoke(ctx, "/bth.Trader/AddOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traderClient) CancelOrder(ctx context.Context, in *CancelOrderRequest, opts ...grpc.CallOption) (*CancelOrderResponse, error) {
	out := new(CancelOrderResponse)
	err := c.cc.Invoke(ctx, "/bth.Trader/CancelOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traderClient) OrderStatus(ctx context.Context, in *OrderStatusRequest, opts ...grpc.CallOption) (*OrderStatusResponse, error) {
	out := new(OrderStatusResponse)
	err := c.cc.Invoke(ctx, "/bth.Trader/OrderStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *traderClient) StreamOrders(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Trader_StreamOrdersClient, error) {
	stream, err := c.cc.NewStream(ctx, &Trader_ServiceDesc.Streams[0], "/bth.Trader/StreamOrders", opts...)
	if err != nil {
		return nil, err
	}
	x := &traderStreamOrdersClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Trader_StreamOrdersClient interface {
	Recv() (*OrderStatusResponse, error)
	grpc.ClientStream
}

type traderStreamOrdersClient struct {
	grpc.ClientStream
}

func (x *traderStreamOrdersClient) Recv() (*OrderStatusResponse, error) {
	m := new(OrderStatusResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TraderServer is the server API for Trader service.
// All implementations must embed UnimplementedTraderServer
// for forward compatibility
type TraderServer interface {
	// AddOrder submits a new order on the exchange
	AddOrder(context.Context, *AddOrderRequest) (*AddOrderResponse, error)
	// CancelOrder cancels an open order
	CancelOrder(context.Context, *CancelOrderRequest) (*CancelOrderResponse, error)
	// OrderStatus request status of particular order
	OrderStatus(context.Context, *OrderStatusRequest) (*OrderStatusResponse, error)
	// StreamOrders opens stream to receive update on order statuses as they become available
	StreamOrders(*Empty, Trader_StreamOrdersServer) error
	mustEmbedUnimplementedTraderServer()
}

// UnimplementedTraderServer must be embedded to have forward compatible implementations.
type UnimplementedTraderServer struct {
}

func (UnimplementedTraderServer) AddOrder(context.Context, *AddOrderRequest) (*AddOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddOrder not implemented")
}
func (UnimplementedTraderServer) CancelOrder(context.Context, *CancelOrderRequest) (*CancelOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelOrder not implemented")
}
func (UnimplementedTraderServer) OrderStatus(context.Context, *OrderStatusRequest) (*OrderStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderStatus not implemented")
}
func (UnimplementedTraderServer) StreamOrders(*Empty, Trader_StreamOrdersServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamOrders not implemented")
}
func (UnimplementedTraderServer) mustEmbedUnimplementedTraderServer() {}

// UnsafeTraderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TraderServer will
// result in compilation errors.
type UnsafeTraderServer interface {
	mustEmbedUnimplementedTraderServer()
}

func RegisterTraderServer(s grpc.ServiceRegistrar, srv TraderServer) {
	s.RegisterService(&Trader_ServiceDesc, srv)
}

func _Trader_AddOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraderServer).AddOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bth.Trader/AddOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraderServer).AddOrder(ctx, req.(*AddOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Trader_CancelOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraderServer).CancelOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bth.Trader/CancelOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraderServer).CancelOrder(ctx, req.(*CancelOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Trader_OrderStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TraderServer).OrderStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bth.Trader/OrderStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TraderServer).OrderStatus(ctx, req.(*OrderStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Trader_StreamOrders_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TraderServer).StreamOrders(m, &traderStreamOrdersServer{stream})
}

type Trader_StreamOrdersServer interface {
	Send(*OrderStatusResponse) error
	grpc.ServerStream
}

type traderStreamOrdersServer struct {
	grpc.ServerStream
}

func (x *traderStreamOrdersServer) Send(m *OrderStatusResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Trader_ServiceDesc is the grpc.ServiceDesc for Trader service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Trader_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bth.Trader",
	HandlerType: (*TraderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddOrder",
			Handler:    _Trader_AddOrder_Handler,
		},
		{
			MethodName: "CancelOrder",
			Handler:    _Trader_CancelOrder_Handler,
		},
		{
			MethodName: "OrderStatus",
			Handler:    _Trader_OrderStatus_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamOrders",
			Handler:       _Trader_StreamOrders_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/proto/trader.proto",
}
