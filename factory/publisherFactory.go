package factory

import (
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

	server := NewServer(cfg.GRPC, blocksPool)

	go func() {
		if serverErr := server.Start(); serverErr != nil {
			log.Error("Failed to start server", "error", serverErr)
		}
	}()

	return process.NewGrpcPublisher(blocksPool, blockCreator, marshaller)
}
