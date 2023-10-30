// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: proto/template.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	TimeAsk_AskForTime_FullMethodName = "/proto.TimeAsk/AskForTime"
)

// TimeAskClient is the client API for TimeAsk service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TimeAskClient interface {
	AskForTime(ctx context.Context, in *AskForTimeMessage, opts ...grpc.CallOption) (*TimeMessage, error)
}

type timeAskClient struct {
	cc grpc.ClientConnInterface
}

func NewTimeAskClient(cc grpc.ClientConnInterface) TimeAskClient {
	return &timeAskClient{cc}
}

func (c *timeAskClient) AskForTime(ctx context.Context, in *AskForTimeMessage, opts ...grpc.CallOption) (*TimeMessage, error) {
	out := new(TimeMessage)
	err := c.cc.Invoke(ctx, TimeAsk_AskForTime_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TimeAskServer is the server API for TimeAsk service.
// All implementations must embed UnimplementedTimeAskServer
// for forward compatibility
type TimeAskServer interface {
	AskForTime(context.Context, *AskForTimeMessage) (*TimeMessage, error)
	mustEmbedUnimplementedTimeAskServer()
}

// UnimplementedTimeAskServer must be embedded to have forward compatible implementations.
type UnimplementedTimeAskServer struct {
}

func (UnimplementedTimeAskServer) AskForTime(context.Context, *AskForTimeMessage) (*TimeMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AskForTime not implemented")
}
func (UnimplementedTimeAskServer) mustEmbedUnimplementedTimeAskServer() {}

// UnsafeTimeAskServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TimeAskServer will
// result in compilation errors.
type UnsafeTimeAskServer interface {
	mustEmbedUnimplementedTimeAskServer()
}

func RegisterTimeAskServer(s grpc.ServiceRegistrar, srv TimeAskServer) {
	s.RegisterService(&TimeAsk_ServiceDesc, srv)
}

func _TimeAsk_AskForTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AskForTimeMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TimeAskServer).AskForTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TimeAsk_AskForTime_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TimeAskServer).AskForTime(ctx, req.(*AskForTimeMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// TimeAsk_ServiceDesc is the grpc.ServiceDesc for TimeAsk service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TimeAsk_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TimeAsk",
	HandlerType: (*TimeAskServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AskForTime",
			Handler:    _TimeAsk_AskForTime_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/template.proto",
}

const (
	Broadcast_PublishReceive_FullMethodName = "/proto.Broadcast/PublishReceive"
	Broadcast_Bob_FullMethodName            = "/proto.Broadcast/bob"
)

// BroadcastClient is the client API for Broadcast service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BroadcastClient interface {
	PublishReceive(ctx context.Context, opts ...grpc.CallOption) (Broadcast_PublishReceiveClient, error)
	Bob(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type broadcastClient struct {
	cc grpc.ClientConnInterface
}

func NewBroadcastClient(cc grpc.ClientConnInterface) BroadcastClient {
	return &broadcastClient{cc}
}

func (c *broadcastClient) PublishReceive(ctx context.Context, opts ...grpc.CallOption) (Broadcast_PublishReceiveClient, error) {
	stream, err := c.cc.NewStream(ctx, &Broadcast_ServiceDesc.Streams[0], Broadcast_PublishReceive_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &broadcastPublishReceiveClient{stream}
	return x, nil
}

type Broadcast_PublishReceiveClient interface {
	Send(*Publish) error
	Recv() (*Publish, error)
	grpc.ClientStream
}

type broadcastPublishReceiveClient struct {
	grpc.ClientStream
}

func (x *broadcastPublishReceiveClient) Send(m *Publish) error {
	return x.ClientStream.SendMsg(m)
}

func (x *broadcastPublishReceiveClient) Recv() (*Publish, error) {
	m := new(Publish)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *broadcastClient) Bob(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Broadcast_Bob_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BroadcastServer is the server API for Broadcast service.
// All implementations must embed UnimplementedBroadcastServer
// for forward compatibility
type BroadcastServer interface {
	PublishReceive(Broadcast_PublishReceiveServer) error
	Bob(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedBroadcastServer()
}

// UnimplementedBroadcastServer must be embedded to have forward compatible implementations.
type UnimplementedBroadcastServer struct {
}

func (UnimplementedBroadcastServer) PublishReceive(Broadcast_PublishReceiveServer) error {
	return status.Errorf(codes.Unimplemented, "method PublishReceive not implemented")
}
func (UnimplementedBroadcastServer) Bob(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bob not implemented")
}
func (UnimplementedBroadcastServer) mustEmbedUnimplementedBroadcastServer() {}

// UnsafeBroadcastServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BroadcastServer will
// result in compilation errors.
type UnsafeBroadcastServer interface {
	mustEmbedUnimplementedBroadcastServer()
}

func RegisterBroadcastServer(s grpc.ServiceRegistrar, srv BroadcastServer) {
	s.RegisterService(&Broadcast_ServiceDesc, srv)
}

func _Broadcast_PublishReceive_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BroadcastServer).PublishReceive(&broadcastPublishReceiveServer{stream})
}

type Broadcast_PublishReceiveServer interface {
	Send(*Publish) error
	Recv() (*Publish, error)
	grpc.ServerStream
}

type broadcastPublishReceiveServer struct {
	grpc.ServerStream
}

func (x *broadcastPublishReceiveServer) Send(m *Publish) error {
	return x.ServerStream.SendMsg(m)
}

func (x *broadcastPublishReceiveServer) Recv() (*Publish, error) {
	m := new(Publish)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Broadcast_Bob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BroadcastServer).Bob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broadcast_Bob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BroadcastServer).Bob(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Broadcast_ServiceDesc is the grpc.ServiceDesc for Broadcast service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Broadcast_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Broadcast",
	HandlerType: (*BroadcastServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "bob",
			Handler:    _Broadcast_Bob_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PublishReceive",
			Handler:       _Broadcast_PublishReceive_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/template.proto",
}