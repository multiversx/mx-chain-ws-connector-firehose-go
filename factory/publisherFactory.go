package factory

import (
	"os"

	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
)

// CreatePublisher will create publisher component
func CreatePublisher(
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

	return process.NewGrpcPublisher(blocksPool, blockCreator, marshaller)
}
