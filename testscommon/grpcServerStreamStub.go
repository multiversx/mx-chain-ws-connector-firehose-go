package testscommon

import (
	"context"

	data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
	"google.golang.org/grpc/metadata"
)

// GRPCServerStreamStub -
type GRPCServerStreamStub struct {
	SetHeaderCalled  func(_ metadata.MD) error
	SendHeaderCalled func(_ metadata.MD) error
	SetTrailerCalled func(_ metadata.MD)
	ContextCalled    func() context.Context
	SendMsgCalled    func(m any) error
	RecvMsgCalled    func(m any) error
	SendCalled       func(block *data.HyperOutportBlock) error
}

// Send -
func (g *GRPCServerStreamStub) Send(block *data.HyperOutportBlock) error {
	if g.SendCalled != nil {
		return g.SendCalled(block)
	}

	return nil
}

// SetHeader -
func (g *GRPCServerStreamStub) SetHeader(md metadata.MD) error {
	if g.SetHeaderCalled != nil {
		return g.SetHeaderCalled(md)
	}

	return nil
}

// SendHeader -
func (g *GRPCServerStreamStub) SendHeader(md metadata.MD) error {
	if g.SendHeaderCalled != nil {
		return g.SendHeaderCalled(md)
	}

	return nil
}

// SetTrailer -
func (g *GRPCServerStreamStub) SetTrailer(md metadata.MD) {
	if g.SetTrailerCalled != nil {
		g.SetTrailerCalled(md)
	}
}

// Context -
func (g *GRPCServerStreamStub) Context() context.Context {
	if g.ContextCalled != nil {
		return g.ContextCalled()
	}

	return context.TODO()
}

// SendMsg -
func (g *GRPCServerStreamStub) SendMsg(m any) error {
	if g.SendMsgCalled != nil {
		return g.SendMsgCalled(m)
	}

	return nil
}

// RecvMsg -
func (g *GRPCServerStreamStub) RecvMsg(m any) error {
	if g.RecvMsgCalled != nil {
		return g.RecvMsgCalled(m)
	}

	return nil
}
