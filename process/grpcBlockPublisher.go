package process

import (
	data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

type grpcBlockPublisher struct {
	server GRPCServer
}

// NewGRPCBlockPublisher is the publisher set up when serving hyperOutportBlocks via gRPC.
func NewGRPCBlockPublisher(server GRPCServer) (*grpcBlockPublisher, error) {
	server.Start()

	return &grpcBlockPublisher{
		server: server,
	}, nil
}

// PublishHyperBlock will do nothing for now, as they are available via gRPC endpoints.
func (g *grpcBlockPublisher) PublishHyperBlock(_ *data.HyperOutportBlock) error {
	return nil
}

// Close will terminate the underlying gRPC server.
func (g *grpcBlockPublisher) Close() error {
	g.server.Close()
	return nil
}

// IsInterfaceNil checks whether it is nil.
func (g *grpcBlockPublisher) IsInterfaceNil() bool {
	return g == nil
}
