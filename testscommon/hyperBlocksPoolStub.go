package testscommon

import (
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

// HyperBlocksPoolStub -
type HyperBlocksPoolStub struct {
	GetCalled             func(hash []byte) ([]byte, error)
	PutMetaBlockCalled    func(hash []byte, outportBlock *hyperOutportBlocks.MetaOutportBlock) error
	PutShardBlockCalled   func(hash []byte, outportBlock *hyperOutportBlocks.ShardOutportBlock) error
	GetMetaBlockCalled    func(hash []byte) (*hyperOutportBlocks.MetaOutportBlock, error)
	GetShardBlockCalled   func(hash []byte) (*hyperOutportBlocks.ShardOutportBlock, error)
	UpdateMetaStateCalled func(checkpoint *data.BlockCheckpoint) error
	CloseCalled           func() error
}

// Get -
func (b *HyperBlocksPoolStub) Get(hash []byte) ([]byte, error) {
	if b.GetCalled != nil {
		return b.GetCalled(hash)
	}

	return []byte{}, nil
}

// PutMetaBlock -
func (b *HyperBlocksPoolStub) PutMetaBlock(hash []byte, outportBlock *hyperOutportBlocks.MetaOutportBlock) error {
	if b.PutMetaBlockCalled != nil {
		return b.PutMetaBlockCalled(hash, outportBlock)
	}

	return nil
}

// PutShardBlock -
func (b *HyperBlocksPoolStub) PutShardBlock(hash []byte, outportBlock *hyperOutportBlocks.ShardOutportBlock) error {
	if b.PutShardBlockCalled != nil {
		return b.PutShardBlockCalled(hash, outportBlock)
	}

	return nil
}

// GetMetaBlock -
func (b *HyperBlocksPoolStub) GetMetaBlock(hash []byte) (*hyperOutportBlocks.MetaOutportBlock, error) {
	if b.GetMetaBlockCalled != nil {
		return b.GetMetaBlockCalled(hash)
	}

	return &hyperOutportBlocks.MetaOutportBlock{}, nil
}

// GetShardBlock -
func (b *HyperBlocksPoolStub) GetShardBlock(hash []byte) (*hyperOutportBlocks.ShardOutportBlock, error) {
	if b.GetShardBlockCalled != nil {
		return b.GetShardBlockCalled(hash)
	}

	return &hyperOutportBlocks.ShardOutportBlock{}, nil
}

// UpdateMetaState -
func (b *HyperBlocksPoolStub) UpdateMetaState(checkpoint *data.BlockCheckpoint) error {
	if b.UpdateMetaStateCalled != nil {
		return b.UpdateMetaStateCalled(checkpoint)
	}

	return nil
}

// Close -
func (b *HyperBlocksPoolStub) Close() error {
	if b.CloseCalled != nil {
		return b.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (b *HyperBlocksPoolStub) IsInterfaceNil() bool {
	return b == nil
}
