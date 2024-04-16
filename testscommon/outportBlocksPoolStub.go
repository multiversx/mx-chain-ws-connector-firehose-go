package testscommon

import "github.com/multiversx/mx-chain-core-go/data/outport"

// OutportBlocksPoolStub -
type OutportBlocksPoolStub struct {
	PutBlockCalled        func(hash []byte, outportBlock *outport.OutportBlock, round uint64) error
	GetBlockCalled        func(hash []byte) (*outport.OutportBlock, error)
	UpdateMetaStateCalled func(round uint64)
	CloseCalled           func() error
}

// PutBlock -
func (b *OutportBlocksPoolStub) PutBlock(hash []byte, outportBlock *outport.OutportBlock, round uint64) error {
	if b.PutBlockCalled != nil {
		return b.PutBlockCalled(hash, outportBlock, round)
	}

	return nil
}

// GetBlock -
func (b *OutportBlocksPoolStub) GetBlock(hash []byte) (*outport.OutportBlock, error) {
	if b.GetBlockCalled != nil {
		return b.GetBlockCalled(hash)
	}

	return &outport.OutportBlock{}, nil
}

// UpdateMetaState -
func (b *OutportBlocksPoolStub) UpdateMetaState(round uint64) {
	if b.UpdateMetaStateCalled != nil {
		b.UpdateMetaStateCalled(round)
	}
}

// Close -
func (b *OutportBlocksPoolStub) Close() error {
	if b.CloseCalled != nil {
		return b.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (b *OutportBlocksPoolStub) IsInterfaceNil() bool {
	return b == nil
}
