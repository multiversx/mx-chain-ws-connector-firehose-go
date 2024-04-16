package process

import (
	"github.com/multiversx/mx-chain-core-go/marshal"

	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

type grpcPublisher struct {
	hyperBlocksPool HyperOutportBlocksPool
	blockCreator    BlockContainerHandler
	marshaller      marshal.Marshalizer
}

// NewGrpcPublisher will create a new grpc publisher component able to publish hyper outport blocks data to blocks pool
// which will then be consumed by the grpc server
func NewGrpcPublisher(
	hyperBlocksPool HyperOutportBlocksPool,
	blockCreator BlockContainerHandler,
	marshaller marshal.Marshalizer,
) (*grpcPublisher, error) {
	return &grpcPublisher{
		hyperBlocksPool: hyperBlocksPool,
		blockCreator:    blockCreator,
		marshaller:      marshaller,
	}, nil
}

func (gp *grpcPublisher) PublishHyperBlock(hyperOutportBlock *data.HyperOutportBlock) error {
	blockHash := hyperOutportBlock.MetaOutportBlock.BlockData.HeaderHash

	round := hyperOutportBlock.MetaOutportBlock.BlockData.Header.Round

	return gp.hyperBlocksPool.PutBlock(blockHash, hyperOutportBlock, round)
}

func (g *grpcPublisher) Close() error {
	return nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (gp *grpcPublisher) IsInterfaceNil() bool {
	return gp == nil
}
