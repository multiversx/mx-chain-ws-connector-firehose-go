package testscommon

// PruningStorerStub -
type PruningStorerStub struct {
	GetCalled   func(key []byte) ([]byte, error)
	PutCalled   func(key []byte, data []byte) error
	PruneCalled func(index uint64) error
	DumpCalled  func() error
	CloseCalled func() error
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

// Dump -
func (p *PruningStorerStub) Dump() error {
	if p.DumpCalled != nil {
		return p.DumpCalled()
	}

	return nil
}

// Close -
func (p *PruningStorerStub) Close() error {
	if p.CloseCalled != nil {
		return p.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (p *PruningStorerStub) IsInterfaceNil() bool {
	return p == nil
}
