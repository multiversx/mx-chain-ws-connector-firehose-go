package testscommon

import (
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
)

// BlocksPoolMock -
type BlocksPoolMock struct {
	PutCalled               func(hash []byte, data []byte) error
	GetCalled               func(hash []byte) ([]byte, error)
	PutBlockCalled          func(hash []byte, value []byte, round uint64, shardID uint32) error
	UpdateMetaStateCalled   func(checkpoint *data.BlockCheckpoint) error
	GetLastCheckpointCalled func() (*data.BlockCheckpoint, error)
	CloseCalled             func() error
}

// Put -
func (b *BlocksPoolMock) Put(hash []byte, data []byte) error {
	if b.PutCalled != nil {
		return b.PutCalled(hash, data)
	}

	return nil
}

// Get -
func (b *BlocksPoolMock) Get(hash []byte) ([]byte, error) {
	if b.GetCalled != nil {
		return b.GetCalled(hash)
	}

	return []byte{}, nil
}

// PutBlock -
func (b *BlocksPoolMock) PutBlock(hash []byte, value []byte, round uint64, shardID uint32) error {
	if b.PutBlockCalled != nil {
		return b.PutBlockCalled(hash, value, round, shardID)
	}

	return nil
}

// UpdateMetaState -
func (b *BlocksPoolMock) UpdateMetaState(checkpoint *data.BlockCheckpoint) error {
	if b.UpdateMetaStateCalled != nil {
		return b.UpdateMetaStateCalled(checkpoint)
	}

	return nil
}

// GetLastCheckpoint -
func (b *BlocksPoolMock) GetLastCheckpoint() (*data.BlockCheckpoint, error) {
	if b.GetLastCheckpointCalled != nil {
		return b.GetLastCheckpointCalled()
	}

	return &data.BlockCheckpoint{
		LastNonces: make(map[uint32]uint64),
	}, nil
}

// Close -
func (b *BlocksPoolMock) Close() error {
	if b.CloseCalled != nil {
		return b.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (b *BlocksPoolMock) IsInterfaceNil() bool {
	return b == nil
}
