package factory

import (
	"fmt"
	"os"

	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/server"
)

// CreatePublisher will return the required Publisher implementation based on whether the hyperOutportBlock are
// served via gRPC or stdout.
func CreatePublisher(
	cfg *config.Config,
	enableGRPCServer bool,
	blockContainer process.BlockContainerHandler,
	outportBlocksPool process.DataPool,
	dataAggregator process.DataAggregator) (process.Publisher, error) {
	if enableGRPCServer {
		handler, err := process.NewGRPCBlocksHandler(outportBlocksPool, dataAggregator)
		if err != nil {
			return nil, fmt.Errorf("failed to create grpc blocks handler: %w", err)
		}

		blocksChannel := make(chan *data.HyperOutportBlock, cfg.GRPC.BufferSize)
		s, err := server.New(cfg.GRPC, handler, &blocksChannel)
		if err != nil {
			return nil, fmt.Errorf("failed to create grpc server: %w", err)
		}

		return process.NewGRPCBlockPublisher(s, &blocksChannel)
	}

	publisher, err := process.NewFirehosePublisher(
		os.Stdout,
		blockContainer,
		&process.ProtoMarshaller{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create firehose publisher: %w", err)
	}

	return publisher, nil
}
