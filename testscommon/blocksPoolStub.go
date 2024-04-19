package testscommon

import "github.com/multiversx/mx-chain-core-go/data/outport"

// BlocksPoolStub -
type BlocksPoolStub struct {
	PutCalled             func(hash []byte, data []byte) error
	GetCalled             func(hash []byte) ([]byte, error)
	PutBlockCalled        func(hash []byte, outportBlock *outport.OutportBlock, round uint64) error
	GetBlockCalled        func(hash []byte) (*outport.OutportBlock, error)
	UpdateMetaStateCalled func(round uint64)
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
func (b *BlocksPoolStub) PutBlock(hash []byte, outportBlock *outport.OutportBlock, round uint64) error {
	if b.PutBlockCalled != nil {
		return b.PutBlockCalled(hash, outportBlock, round)
	}

	return nil
}

// GetBlock -
func (b *BlocksPoolStub) GetBlock(hash []byte) (*outport.OutportBlock, error) {
	if b.GetBlockCalled != nil {
		return b.GetBlockCalled(hash)
	}

	return &outport.OutportBlock{}, nil
}

// UpdateMetaState -
func (b *BlocksPoolStub) UpdateMetaState(round uint64) {
	if b.UpdateMetaStateCalled != nil {
		b.UpdateMetaStateCalled(round)
	}
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
