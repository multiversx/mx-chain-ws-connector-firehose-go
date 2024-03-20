package testscommon

import "github.com/multiversx/mx-chain-core-go/data/outport"

// BlocksPoolStub -
type BlocksPoolStub struct {
	PutBlockCalled func(hash []byte, outportBlock *outport.OutportBlock, round uint64) error
	GetBlockCalled func(hash []byte) (*outport.OutportBlock, error)
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

// UpdateMetaRound -
func (b *BlocksPoolStub) UpdateMetaRound(round uint64) {
}

// IsInterfaceNil -
func (b *BlocksPoolStub) IsInterfaceNil() bool {
	return b == nil
}
