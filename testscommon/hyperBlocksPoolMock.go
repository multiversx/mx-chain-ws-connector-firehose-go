package testscommon

import (
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
)

// HyperBlocksPoolMock -
type HyperBlocksPoolMock struct {
	PutCalled               func(hash []byte, data []byte) error
	GetCalled               func(hash []byte) ([]byte, error)
	PutBlockCalled          func(hash []byte, outportBlock process.OutportBlockHandler) error
	GetMetaBlockCalled      func(hash []byte) (*hyperOutportBlocks.MetaOutportBlock, error)
	GetShardBlockCalled     func(hash []byte) (*hyperOutportBlocks.ShardOutportBlock, error)
	UpdateMetaStateCalled   func(checkpoint *data.BlockCheckpoint) error
	GetLastCheckpointCalled func() (*data.BlockCheckpoint, error)
	CloseCalled             func() error
}

// Get -
func (b *HyperBlocksPoolMock) Get(hash []byte) ([]byte, error) {
	if b.GetCalled != nil {
		return b.GetCalled(hash)
	}

	return []byte{}, nil
}

// Put -
func (b *HyperBlocksPoolMock) Put(hash []byte, data []byte) error {
	if b.PutCalled != nil {
		return b.PutCalled(hash, data)
	}

	return nil
}

// PutBlock -
func (b *HyperBlocksPoolMock) PutBlock(hash []byte, outportBlock process.OutportBlockHandler) error {
	if b.PutBlockCalled != nil {
		return b.PutBlockCalled(hash, outportBlock)
	}

	return nil
}

// GetMetaBlock -
func (b *HyperBlocksPoolMock) GetMetaBlock(hash []byte) (*hyperOutportBlocks.MetaOutportBlock, error) {
	if b.GetMetaBlockCalled != nil {
		return b.GetMetaBlockCalled(hash)
	}

	return &hyperOutportBlocks.MetaOutportBlock{}, nil
}

// GetShardBlock -
func (b *HyperBlocksPoolMock) GetShardBlock(hash []byte) (*hyperOutportBlocks.ShardOutportBlock, error) {
	if b.GetShardBlockCalled != nil {
		return b.GetShardBlockCalled(hash)
	}

	return &hyperOutportBlocks.ShardOutportBlock{}, nil
}

// UpdateMetaState -
func (b *HyperBlocksPoolMock) UpdateMetaState(checkpoint *data.BlockCheckpoint) error {
	if b.UpdateMetaStateCalled != nil {
		return b.UpdateMetaStateCalled(checkpoint)
	}

	return nil
}

// GetLastCheckpoint -
func (b *HyperBlocksPoolMock) GetLastCheckpoint() (*data.BlockCheckpoint, error) {
	if b.GetLastCheckpointCalled != nil {
		return b.GetLastCheckpointCalled()
	}

	return &data.BlockCheckpoint{
		LastNonces: make(map[uint32]uint64),
	}, nil
}

// Close -
func (b *HyperBlocksPoolMock) Close() error {
	if b.CloseCalled != nil {
		return b.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (b *HyperBlocksPoolMock) IsInterfaceNil() bool {
	return b == nil
}
