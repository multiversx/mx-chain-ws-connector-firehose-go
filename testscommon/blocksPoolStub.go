package testscommon

import (
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
)

// BlocksPoolStub -
type BlocksPoolStub struct {
	PutCalled             func(hash []byte, data []byte) error
	GetCalled             func(hash []byte) ([]byte, error)
	PutBlockCalled        func(hash []byte, value []byte, round uint64, shardID uint32) error
	UpdateMetaStateCalled func(checkpoint *data.BlockCheckpoint) error
	CloseCalled           func() error
}

// Put -
func (b *BlocksPoolStub) Put(hash []byte, data []byte) error {
	if b.PutCalled != nil {
		return b.PutCalled(hash, data)
	}

	return nil
}

// Get -
func (b *BlocksPoolStub) Get(hash []byte) ([]byte, error) {
	if b.GetCalled != nil {
		return b.GetCalled(hash)
	}

	return []byte{}, nil
}

// PutBlock -
func (b *BlocksPoolStub) PutBlock(hash []byte, value []byte, round uint64, shardID uint32) error {
	if b.PutBlockCalled != nil {
		return b.PutBlockCalled(hash, value, round, shardID)
	}

	return nil
}

// UpdateMetaState -
func (b *BlocksPoolStub) UpdateMetaState(checkpoint *data.BlockCheckpoint) error {
	if b.UpdateMetaStateCalled != nil {
		return b.UpdateMetaStateCalled(checkpoint)
	}

	return nil
}

// Close -
func (b *BlocksPoolStub) Close() error {
	if b.CloseCalled != nil {
		return b.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (b *BlocksPoolStub) IsInterfaceNil() bool {
	return b == nil
}
