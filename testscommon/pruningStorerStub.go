package testscommon

// PruningStorerStub -
type PruningStorerStub struct {
	GetCalled   func(key []byte) ([]byte, error)
	PutCalled   func(key []byte, data []byte) error
	PruneCalled func(index uint64) error
}

// Get -
func (p *PruningStorerStub) Get(key []byte) ([]byte, error) {
	if p.GetCalled != nil {
		return p.GetCalled(key)
	}

	return nil, nil
}

// Put -
func (p *PruningStorerStub) Put(key []byte, data []byte) error {
	if p.PutCalled != nil {
		return p.PutCalled(key, data)
	}

	return nil
}

// Prune -
func (p *PruningStorerStub) Prune(index uint64) error {
	if p.PruneCalled != nil {
		return p.PruneCalled(index)
	}

	return nil
}

// IsInterfaceNil -
func (p *PruningStorerStub) IsInterfaceNil() bool {
	return p == nil
}
