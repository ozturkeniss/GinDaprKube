// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: api/proto/basket/basket.proto

package basket

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

const (
	BasketService_GetBasket_FullMethodName      = "/basket.BasketService/GetBasket"
	BasketService_AddItem_FullMethodName        = "/basket.BasketService/AddItem"
	BasketService_RemoveItem_FullMethodName     = "/basket.BasketService/RemoveItem"
	BasketService_UpdateQuantity_FullMethodName = "/basket.BasketService/UpdateQuantity"
	BasketService_ClearBasket_FullMethodName    = "/basket.BasketService/ClearBasket"
)

// BasketServiceClient is the client API for BasketService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BasketServiceClient interface {
	GetBasket(ctx context.Context, in *GetBasketRequest, opts ...grpc.CallOption) (*GetBasketResponse, error)
	AddItem(ctx context.Context, in *AddItemRequest, opts ...grpc.CallOption) (*AddItemResponse, error)
	RemoveItem(ctx context.Context, in *RemoveItemRequest, opts ...grpc.CallOption) (*RemoveItemResponse, error)
	UpdateQuantity(ctx context.Context, in *UpdateQuantityRequest, opts ...grpc.CallOption) (*UpdateQuantityResponse, error)
	ClearBasket(ctx context.Context, in *ClearBasketRequest, opts ...grpc.CallOption) (*ClearBasketResponse, error)
}

type basketServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBasketServiceClient(cc grpc.ClientConnInterface) BasketServiceClient {
	return &basketServiceClient{cc}
}

func (c *basketServiceClient) GetBasket(ctx context.Context, in *GetBasketRequest, opts ...grpc.CallOption) (*GetBasketResponse, error) {
	out := new(GetBasketResponse)
	err := c.cc.Invoke(ctx, BasketService_GetBasket_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basketServiceClient) AddItem(ctx context.Context, in *AddItemRequest, opts ...grpc.CallOption) (*AddItemResponse, error) {
	out := new(AddItemResponse)
	err := c.cc.Invoke(ctx, BasketService_AddItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basketServiceClient) RemoveItem(ctx context.Context, in *RemoveItemRequest, opts ...grpc.CallOption) (*RemoveItemResponse, error) {
	out := new(RemoveItemResponse)
	err := c.cc.Invoke(ctx, BasketService_RemoveItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basketServiceClient) UpdateQuantity(ctx context.Context, in *UpdateQuantityRequest, opts ...grpc.CallOption) (*UpdateQuantityResponse, error) {
	out := new(UpdateQuantityResponse)
	err := c.cc.Invoke(ctx, BasketService_UpdateQuantity_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *basketServiceClient) ClearBasket(ctx context.Context, in *ClearBasketRequest, opts ...grpc.CallOption) (*ClearBasketResponse, error) {
	out := new(ClearBasketResponse)
	err := c.cc.Invoke(ctx, BasketService_ClearBasket_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BasketServiceServer is the server API for BasketService service.
// All implementations must embed UnimplementedBasketServiceServer
// for forward compatibility
type BasketServiceServer interface {
	GetBasket(context.Context, *GetBasketRequest) (*GetBasketResponse, error)
	AddItem(context.Context, *AddItemRequest) (*AddItemResponse, error)
	RemoveItem(context.Context, *RemoveItemRequest) (*RemoveItemResponse, error)
	UpdateQuantity(context.Context, *UpdateQuantityRequest) (*UpdateQuantityResponse, error)
	ClearBasket(context.Context, *ClearBasketRequest) (*ClearBasketResponse, error)
	mustEmbedUnimplementedBasketServiceServer()
}

// UnimplementedBasketServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBasketServiceServer struct {
}

func (UnimplementedBasketServiceServer) GetBasket(context.Context, *GetBasketRequest) (*GetBasketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBasket not implemented")
}
func (UnimplementedBasketServiceServer) AddItem(context.Context, *AddItemRequest) (*AddItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddItem not implemented")
}
func (UnimplementedBasketServiceServer) RemoveItem(context.Context, *RemoveItemRequest) (*RemoveItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveItem not implemented")
}
func (UnimplementedBasketServiceServer) UpdateQuantity(context.Context, *UpdateQuantityRequest) (*UpdateQuantityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateQuantity not implemented")
}
func (UnimplementedBasketServiceServer) ClearBasket(context.Context, *ClearBasketRequest) (*ClearBasketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearBasket not implemented")
}
func (UnimplementedBasketServiceServer) mustEmbedUnimplementedBasketServiceServer() {}

// UnsafeBasketServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BasketServiceServer will
// result in compilation errors.
type UnsafeBasketServiceServer interface {
	mustEmbedUnimplementedBasketServiceServer()
}

func RegisterBasketServiceServer(s grpc.ServiceRegistrar, srv BasketServiceServer) {
	s.RegisterService(&BasketService_ServiceDesc, srv)
}

func _BasketService_GetBasket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBasketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasketServiceServer).GetBasket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BasketService_GetBasket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasketServiceServer).GetBasket(ctx, req.(*GetBasketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasketService_AddItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasketServiceServer).AddItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BasketService_AddItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasketServiceServer).AddItem(ctx, req.(*AddItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasketService_RemoveItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasketServiceServer).RemoveItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BasketService_RemoveItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasketServiceServer).RemoveItem(ctx, req.(*RemoveItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasketService_UpdateQuantity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateQuantityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasketServiceServer).UpdateQuantity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BasketService_UpdateQuantity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasketServiceServer).UpdateQuantity(ctx, req.(*UpdateQuantityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BasketService_ClearBasket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearBasketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BasketServiceServer).ClearBasket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BasketService_ClearBasket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BasketServiceServer).ClearBasket(ctx, req.(*ClearBasketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BasketService_ServiceDesc is the grpc.ServiceDesc for BasketService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BasketService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "basket.BasketService",
	HandlerType: (*BasketServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBasket",
			Handler:    _BasketService_GetBasket_Handler,
		},
		{
			MethodName: "AddItem",
			Handler:    _BasketService_AddItem_Handler,
		},
		{
			MethodName: "RemoveItem",
			Handler:    _BasketService_RemoveItem_Handler,
		},
		{
			MethodName: "UpdateQuantity",
			Handler:    _BasketService_UpdateQuantity_Handler,
		},
		{
			MethodName: "ClearBasket",
			Handler:    _BasketService_ClearBasket_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/basket/basket.proto",
}
