package testscommon

// BlocksPoolStub -
type BlocksPoolStub struct {
	PutBlockCalled        func(hash []byte, data []byte, round uint64, shardID uint32) error
	GetBlockCalled        func(hash []byte) ([]byte, error)
	UpdateMetaStateCalled func(round uint64)
	CloseCalled           func() error
}

// PutBlock -
func (b *BlocksPoolStub) PutBlock(hash []byte, data []byte, round uint64, shardID uint32) error {
	if b.PutBlockCalled != nil {
		return b.PutBlockCalled(hash, data, round, shardID)
	}

	return nil
}

// GetBlock -
func (b *BlocksPoolStub) GetBlock(hash []byte) ([]byte, error) {
	if b.GetBlockCalled != nil {
		return b.GetBlockCalled(hash)
	}

	return []byte{}, nil
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
