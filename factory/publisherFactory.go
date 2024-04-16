package factory

import (
	"fmt"
	"os"

	"github.com/multiversx/mx-chain-core-go/marshal"

	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
)

// CreatePublisher will create publisher component
func CreatePublisher(cfg config.Config,
	isGrpcServer bool,
	blockCreator process.BlockContainerHandler,
	marshaller marshal.Marshalizer,
	blocksPool process.HyperOutportBlocksPool,
) (process.Publisher, error) {
	if !isGrpcServer {
		return process.NewFirehosePublisher(
			os.Stdout,
			blockCreator,
			marshaller,
		)
	}

	hyperOutportBlocksPool, err := dataPool.NewHyperOutportBlocksPool(blocksPool, marshaller)
	if err != nil {
		return nil, fmt.Errorf("failed to create hyper blocks pool: %w", err)
	}

	server := NewServer(cfg.GRPC, hyperOutportBlocksPool)

	go func() {
		if serverErr := server.Start(); serverErr != nil {
			log.Error("Failed to start server", "error", serverErr)
		}
	}()

	return process.NewGrpcPublisher(hyperOutportBlocksPool, blockCreator, marshaller)
}
