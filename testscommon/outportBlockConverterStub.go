package testscommon

import (
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

// OutportBlockConverterStub -
type OutportBlockConverterStub struct {
	HandleShardOutportBlockCalled func(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.ShardOutportBlock, error)
	HandleMetaOutportBlockCalled  func(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.MetaOutportBlock, error)
}

// HandleShardOutportBlock -
func (o *OutportBlockConverterStub) HandleShardOutportBlock(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.ShardOutportBlock, error) {
	if o.HandleShardOutportBlockCalled != nil {
		return o.HandleShardOutportBlockCalled(outportBlock)
	}

	return &hyperOutportBlocks.ShardOutportBlock{}, nil
}

// HandleMetaOutportBlock -
func (o *OutportBlockConverterStub) HandleMetaOutportBlock(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.MetaOutportBlock, error) {
	if o.HandleMetaOutportBlockCalled != nil {
		return o.HandleMetaOutportBlockCalled(outportBlock)
	}

	return &hyperOutportBlocks.MetaOutportBlock{}, nil
}

// IsInterfaceNil -
func (o *OutportBlockConverterStub) IsInterfaceNil() bool {
	return o == nil
}
