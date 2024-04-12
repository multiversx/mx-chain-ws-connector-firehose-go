package process

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
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

	round, err := gp.getHeaderRound(hyperOutportBlock.MetaOutportBlock)
	if err != nil {
		return err
	}

	return gp.hyperBlocksPool.PutBlock(blockHash, hyperOutportBlock, round)
}

func (gp *grpcPublisher) getHeaderRound(outportBlock *outport.OutportBlock) (uint64, error) {
	blockCreator, err := gp.blockCreator.Get(core.HeaderType(outportBlock.BlockData.HeaderType))
	if err != nil {
		return 0, err
	}

	header, err := block.GetHeaderFromBytes(gp.marshaller, blockCreator, outportBlock.BlockData.HeaderBytes)
	if err != nil {
		return 0, err
	}

	return header.GetRound(), nil
}

func (g *grpcPublisher) Close() error {
	return nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (gp *grpcPublisher) IsInterfaceNil() bool {
	return gp == nil
}
