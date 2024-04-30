package testscommon

// PublisherStub -
type PublisherStub struct {
	PublishBlockCalled func(headerHash []byte) error
	CloseCalled        func() error
}

// PublishBlock -
func (p *PublisherStub) PublishBlock(headerHash []byte) error {
	if p.PublishBlockCalled != nil {
		return p.PublishBlockCalled(headerHash)
	}

	return nil
}

// Close -
func (p *PublisherStub) Close() error {
	if p.CloseCalled != nil {
		return p.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (p *PublisherStub) IsInterfaceNil() bool {
	return p == nil
}
