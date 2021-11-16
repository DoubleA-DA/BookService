// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protof

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

// SRVClient is the client API for SRV service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SRVClient interface {
	GetBook(ctx context.Context, in *Limit, opts ...grpc.CallOption) (*Dispbook, error)
	InsBook(ctx context.Context, in *Book, opts ...grpc.CallOption) (*Ack, error)
	InsReview(ctx context.Context, in *Review, opts ...grpc.CallOption) (*Ack, error)
	GetReview(ctx context.Context, in *Review, opts ...grpc.CallOption) (*Dispreview, error)
}

type sRVClient struct {
	cc grpc.ClientConnInterface
}

func NewSRVClient(cc grpc.ClientConnInterface) SRVClient {
	return &sRVClient{cc}
}

func (c *sRVClient) GetBook(ctx context.Context, in *Limit, opts ...grpc.CallOption) (*Dispbook, error) {
	out := new(Dispbook)
	err := c.cc.Invoke(ctx, "/protof.SRV/get_book", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sRVClient) InsBook(ctx context.Context, in *Book, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/protof.SRV/ins_book", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sRVClient) InsReview(ctx context.Context, in *Review, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/protof.SRV/ins_review", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sRVClient) GetReview(ctx context.Context, in *Review, opts ...grpc.CallOption) (*Dispreview, error) {
	out := new(Dispreview)
	err := c.cc.Invoke(ctx, "/protof.SRV/get_review", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SRVServer is the server API for SRV service.
// All implementations must embed UnimplementedSRVServer
// for forward compatibility
type SRVServer interface {
	GetBook(context.Context, *Limit) (*Dispbook, error)
	InsBook(context.Context, *Book) (*Ack, error)
	InsReview(context.Context, *Review) (*Ack, error)
	GetReview(context.Context, *Review) (*Dispreview, error)
	mustEmbedUnimplementedSRVServer()
}

// UnimplementedSRVServer must be embedded to have forward compatible implementations.
type UnimplementedSRVServer struct {
}

func (UnimplementedSRVServer) GetBook(context.Context, *Limit) (*Dispbook, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBook not implemented")
}
func (UnimplementedSRVServer) InsBook(context.Context, *Book) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsBook not implemented")
}
func (UnimplementedSRVServer) InsReview(context.Context, *Review) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsReview not implemented")
}
func (UnimplementedSRVServer) GetReview(context.Context, *Review) (*Dispreview, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReview not implemented")
}
func (UnimplementedSRVServer) mustEmbedUnimplementedSRVServer() {}

// UnsafeSRVServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SRVServer will
// result in compilation errors.
type UnsafeSRVServer interface {
	mustEmbedUnimplementedSRVServer()
}

func RegisterSRVServer(s grpc.ServiceRegistrar, srv SRVServer) {
	s.RegisterService(&SRV_ServiceDesc, srv)
}

func _SRV_GetBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Limit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SRVServer).GetBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protof.SRV/get_book",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SRVServer).GetBook(ctx, req.(*Limit))
	}
	return interceptor(ctx, in, info, handler)
}

func _SRV_InsBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Book)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SRVServer).InsBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protof.SRV/ins_book",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SRVServer).InsBook(ctx, req.(*Book))
	}
	return interceptor(ctx, in, info, handler)
}

func _SRV_InsReview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Review)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SRVServer).InsReview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protof.SRV/ins_review",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SRVServer).InsReview(ctx, req.(*Review))
	}
	return interceptor(ctx, in, info, handler)
}

func _SRV_GetReview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Review)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SRVServer).GetReview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protof.SRV/get_review",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SRVServer).GetReview(ctx, req.(*Review))
	}
	return interceptor(ctx, in, info, handler)
}

// SRV_ServiceDesc is the grpc.ServiceDesc for SRV service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SRV_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protof.SRV",
	HandlerType: (*SRVServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "get_book",
			Handler:    _SRV_GetBook_Handler,
		},
		{
			MethodName: "ins_book",
			Handler:    _SRV_InsBook_Handler,
		},
		{
			MethodName: "ins_review",
			Handler:    _SRV_InsReview_Handler,
		},
		{
			MethodName: "get_review",
			Handler:    _SRV_GetReview_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "prt.proto",
}
