package factory

import (
	"os"

	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process/dataPool"
)

// CreatePublisher will create publisher component
func CreatePublisher(
	isGrpcServer bool,
	blockCreator process.BlockContainerHandler,
	marshaller marshal.Marshalizer,
	blocksPool process.DataPool,
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
		return nil, err
	}

	return process.NewGrpcPublisher(hyperOutportBlocksPool, blockCreator, marshaller)
}
