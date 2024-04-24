package process

import (
	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

type grpcBlockPublisher struct{}

// NewGRPCBlockPublisher is the publisher set up when serving hyperOutportBlocks via gRPC.
func NewGRPCBlockPublisher() *grpcBlockPublisher {
	return &grpcBlockPublisher{}
}

// PublishHyperBlock will do nothing for now, as they are available via gRPC endpoints.
func (g *grpcBlockPublisher) PublishHyperBlock(_ *data.HyperOutportBlock) error {
	return nil
}

// Close -
func (g *grpcBlockPublisher) Close() error {
	return nil
}

// IsInterfaceNil checks whether it is nil.
func (g *grpcBlockPublisher) IsInterfaceNil() bool {
	return g == nil
}
