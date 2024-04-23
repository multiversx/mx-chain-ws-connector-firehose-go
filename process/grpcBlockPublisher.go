package process

import (
	"fmt"

	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

type grpcBlockPublisher struct {
	server GRPCServer
	queue  HyperOutportBlocksQueue
}

// NewGRPCBlockPublisher is the publisher set up when serving hyperOutportBlocks via gRPC.
func NewGRPCBlockPublisher(server GRPCServer, queue HyperOutportBlocksQueue) (*grpcBlockPublisher, error) {
	server.Start()

	return &grpcBlockPublisher{
		server: server,
		queue:  queue,
	}, nil
}

// PublishHyperBlock will do nothing for now, as they are available via gRPC endpoints.
func (g *grpcBlockPublisher) PublishHyperBlock(block *data.HyperOutportBlock) error {
	err := g.queue.Enqueue(block)
	if err != nil {
		return fmt.Errorf("failed to enqueue block: %w", err)
	}

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
