// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: grpcBlockService.proto

package hyperOutportBlocks

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

// BlockStreamClient is the client API for BlockStream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlockStreamClient interface {
	BlocksByHash(ctx context.Context, in *BlockHashStreamRequest, opts ...grpc.CallOption) (BlockStream_BlocksByHashClient, error)
	BlocksByNonce(ctx context.Context, in *BlockNonceStreamRequest, opts ...grpc.CallOption) (BlockStream_BlocksByNonceClient, error)
}

type blockStreamClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockStreamClient(cc grpc.ClientConnInterface) BlockStreamClient {
	return &blockStreamClient{cc}
}

func (c *blockStreamClient) BlocksByHash(ctx context.Context, in *BlockHashStreamRequest, opts ...grpc.CallOption) (BlockStream_BlocksByHashClient, error) {
	stream, err := c.cc.NewStream(ctx, &BlockStream_ServiceDesc.Streams[0], "/proto.BlockStream/BlocksByHash", opts...)
	if err != nil {
		return nil, err
	}
	x := &blockStreamBlocksByHashClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BlockStream_BlocksByHashClient interface {
	Recv() (*HyperOutportBlock, error)
	grpc.ClientStream
}

type blockStreamBlocksByHashClient struct {
	grpc.ClientStream
}

func (x *blockStreamBlocksByHashClient) Recv() (*HyperOutportBlock, error) {
	m := new(HyperOutportBlock)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *blockStreamClient) BlocksByNonce(ctx context.Context, in *BlockNonceStreamRequest, opts ...grpc.CallOption) (BlockStream_BlocksByNonceClient, error) {
	stream, err := c.cc.NewStream(ctx, &BlockStream_ServiceDesc.Streams[1], "/proto.BlockStream/BlocksByNonce", opts...)
	if err != nil {
		return nil, err
	}
	x := &blockStreamBlocksByNonceClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BlockStream_BlocksByNonceClient interface {
	Recv() (*HyperOutportBlock, error)
	grpc.ClientStream
}

type blockStreamBlocksByNonceClient struct {
	grpc.ClientStream
}

func (x *blockStreamBlocksByNonceClient) Recv() (*HyperOutportBlock, error) {
	m := new(HyperOutportBlock)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BlockStreamServer is the server API for BlockStream service.
// All implementations must embed UnimplementedBlockStreamServer
// for forward compatibility
type BlockStreamServer interface {
	BlocksByHash(*BlockHashStreamRequest, BlockStream_BlocksByHashServer) error
	BlocksByNonce(*BlockNonceStreamRequest, BlockStream_BlocksByNonceServer) error
	mustEmbedUnimplementedBlockStreamServer()
}

// UnimplementedBlockStreamServer must be embedded to have forward compatible implementations.
type UnimplementedBlockStreamServer struct {
}

func (UnimplementedBlockStreamServer) BlocksByHash(*BlockHashStreamRequest, BlockStream_BlocksByHashServer) error {
	return status.Errorf(codes.Unimplemented, "method BlocksByHash not implemented")
}
func (UnimplementedBlockStreamServer) BlocksByNonce(*BlockNonceStreamRequest, BlockStream_BlocksByNonceServer) error {
	return status.Errorf(codes.Unimplemented, "method BlocksByNonce not implemented")
}
func (UnimplementedBlockStreamServer) mustEmbedUnimplementedBlockStreamServer() {}

// UnsafeBlockStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockStreamServer will
// result in compilation errors.
type UnsafeBlockStreamServer interface {
	mustEmbedUnimplementedBlockStreamServer()
}

func RegisterBlockStreamServer(s grpc.ServiceRegistrar, srv BlockStreamServer) {
	s.RegisterService(&BlockStream_ServiceDesc, srv)
}

func _BlockStream_BlocksByHash_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BlockHashStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BlockStreamServer).BlocksByHash(m, &blockStreamBlocksByHashServer{stream})
}

type BlockStream_BlocksByHashServer interface {
	Send(*HyperOutportBlock) error
	grpc.ServerStream
}

type blockStreamBlocksByHashServer struct {
	grpc.ServerStream
}

func (x *blockStreamBlocksByHashServer) Send(m *HyperOutportBlock) error {
	return x.ServerStream.SendMsg(m)
}

func _BlockStream_BlocksByNonce_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BlockNonceStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BlockStreamServer).BlocksByNonce(m, &blockStreamBlocksByNonceServer{stream})
}

type BlockStream_BlocksByNonceServer interface {
	Send(*HyperOutportBlock) error
	grpc.ServerStream
}

type blockStreamBlocksByNonceServer struct {
	grpc.ServerStream
}

func (x *blockStreamBlocksByNonceServer) Send(m *HyperOutportBlock) error {
	return x.ServerStream.SendMsg(m)
}

// BlockStream_ServiceDesc is the grpc.ServiceDesc for BlockStream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlockStream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.BlockStream",
	HandlerType: (*BlockStreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BlocksByHash",
			Handler:       _BlockStream_BlocksByHash_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "BlocksByNonce",
			Handler:       _BlockStream_BlocksByNonce_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "grpcBlockService.proto",
}

// BlockFetchClient is the client API for BlockFetch service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlockFetchClient interface {
	GetBlockByHash(ctx context.Context, in *BlockHashRequest, opts ...grpc.CallOption) (*HyperOutportBlock, error)
	GetBlockByNonce(ctx context.Context, in *BlockNonceRequest, opts ...grpc.CallOption) (*HyperOutportBlock, error)
}

type blockFetchClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockFetchClient(cc grpc.ClientConnInterface) BlockFetchClient {
	return &blockFetchClient{cc}
}

func (c *blockFetchClient) GetBlockByHash(ctx context.Context, in *BlockHashRequest, opts ...grpc.CallOption) (*HyperOutportBlock, error) {
	out := new(HyperOutportBlock)
	err := c.cc.Invoke(ctx, "/proto.BlockFetch/GetBlockByHash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockFetchClient) GetBlockByNonce(ctx context.Context, in *BlockNonceRequest, opts ...grpc.CallOption) (*HyperOutportBlock, error) {
	out := new(HyperOutportBlock)
	err := c.cc.Invoke(ctx, "/proto.BlockFetch/GetBlockByNonce", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlockFetchServer is the server API for BlockFetch service.
// All implementations must embed UnimplementedBlockFetchServer
// for forward compatibility
type BlockFetchServer interface {
	GetBlockByHash(context.Context, *BlockHashRequest) (*HyperOutportBlock, error)
	GetBlockByNonce(context.Context, *BlockNonceRequest) (*HyperOutportBlock, error)
	mustEmbedUnimplementedBlockFetchServer()
}

// UnimplementedBlockFetchServer must be embedded to have forward compatible implementations.
type UnimplementedBlockFetchServer struct {
}

func (UnimplementedBlockFetchServer) GetBlockByHash(context.Context, *BlockHashRequest) (*HyperOutportBlock, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockByHash not implemented")
}
func (UnimplementedBlockFetchServer) GetBlockByNonce(context.Context, *BlockNonceRequest) (*HyperOutportBlock, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockByNonce not implemented")
}
func (UnimplementedBlockFetchServer) mustEmbedUnimplementedBlockFetchServer() {}

// UnsafeBlockFetchServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockFetchServer will
// result in compilation errors.
type UnsafeBlockFetchServer interface {
	mustEmbedUnimplementedBlockFetchServer()
}

func RegisterBlockFetchServer(s grpc.ServiceRegistrar, srv BlockFetchServer) {
	s.RegisterService(&BlockFetch_ServiceDesc, srv)
}

func _BlockFetch_GetBlockByHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockHashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockFetchServer).GetBlockByHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.BlockFetch/GetBlockByHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockFetchServer).GetBlockByHash(ctx, req.(*BlockHashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockFetch_GetBlockByNonce_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockNonceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockFetchServer).GetBlockByNonce(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.BlockFetch/GetBlockByNonce",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockFetchServer).GetBlockByNonce(ctx, req.(*BlockNonceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BlockFetch_ServiceDesc is the grpc.ServiceDesc for BlockFetch service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlockFetch_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.BlockFetch",
	HandlerType: (*BlockFetchServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBlockByHash",
			Handler:    _BlockFetch_GetBlockByHash_Handler,
		},
		{
			MethodName: "GetBlockByNonce",
			Handler:    _BlockFetch_GetBlockByNonce_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpcBlockService.proto",
}
