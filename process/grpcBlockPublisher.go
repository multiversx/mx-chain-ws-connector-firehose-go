package process

import (
	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

type grpcBlockPublisher struct {
	server        GRPCServer
	blocksChannel *chan *data.HyperOutportBlock
}

// NewGRPCBlockPublisher is the publisher set up when serving hyperOutportBlocks via gRPC.
func NewGRPCBlockPublisher(server GRPCServer, blocksChannel *chan *data.HyperOutportBlock) (*grpcBlockPublisher, error) {
	server.Start()

	return &grpcBlockPublisher{
		server:        server,
		blocksChannel: blocksChannel,
	}, nil
}

// PublishHyperBlock will do nothing for now, as they are available via gRPC endpoints.
func (g *grpcBlockPublisher) PublishHyperBlock(block *data.HyperOutportBlock) error {
	*g.blocksChannel <- block

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
