package testscommon

import (
	"github.com/multiversx/mx-chain-core-go/data/outport"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

// OutportBlockConverterMock -
type OutportBlockConverterMock struct {
	HandleShardOutportBlockCalled   func(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.ShardOutportBlock, error)
	HandleShardOutportBlockV2Called func(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.ShardOutportBlockV2, error)
	HandleMetaOutportBlockCalled    func(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.MetaOutportBlock, error)
}

// HandleShardOutportBlock -
func (o *OutportBlockConverterMock) HandleShardOutportBlock(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.ShardOutportBlock, error) {
	if o.HandleShardOutportBlockCalled != nil {
		return o.HandleShardOutportBlockCalled(outportBlock)
	}

	return &hyperOutportBlocks.ShardOutportBlock{}, nil
}

// HandleShardOutportBlockV2 -
func (o *OutportBlockConverterMock) HandleShardOutportBlockV2(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.ShardOutportBlockV2, error) {
	if o.HandleShardOutportBlockCalled != nil {
		return o.HandleShardOutportBlockV2Called(outportBlock)
	}

	return &hyperOutportBlocks.ShardOutportBlockV2{}, nil
}

// HandleMetaOutportBlock -
func (o *OutportBlockConverterMock) HandleMetaOutportBlock(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.MetaOutportBlock, error) {
	if o.HandleMetaOutportBlockCalled != nil {
		return o.HandleMetaOutportBlockCalled(outportBlock)
	}

	return &hyperOutportBlocks.MetaOutportBlock{}, nil
}

// IsInterfaceNil -
func (o *OutportBlockConverterMock) IsInterfaceNil() bool {
	return o == nil
}
