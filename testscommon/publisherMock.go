package testscommon

// PublisherMock -
type PublisherMock struct {
	PublishBlockCalled func(headerHash []byte) error
	CloseCalled        func() error
}

// PublishBlock -
func (p *PublisherMock) PublishBlock(headerHash []byte) error {
	if p.PublishBlockCalled != nil {
		return p.PublishBlockCalled(headerHash)
	}

	return nil
}

// Close -
func (p *PublisherMock) Close() error {
	if p.CloseCalled != nil {
		return p.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (p *PublisherMock) IsInterfaceNil() bool {
	return p == nil
}
