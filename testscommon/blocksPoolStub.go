package testscommon

import "github.com/multiversx/mx-chain-core-go/data/outport"

// BlocksPoolStub -
type BlocksPoolStub struct {
	PutBlockCalled func(hash []byte, outportBlock *outport.OutportBlock) error
	GetBlockCalled func(hash []byte) (*outport.OutportBlock, error)
}

// PutBlock -
func (b *BlocksPoolStub) PutBlock(hash []byte, outportBlock *outport.OutportBlock) error {
	if b.PutBlockCalled != nil {
		return b.PutBlockCalled(hash, outportBlock)
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

// IsInterfaceNil -
func (b *BlocksPoolStub) IsInterfaceNil() bool {
	return b == nil
}
