package testscommon

import data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"

// PublisherMock -
type PublisherMock struct {
	PublishHyperBlockCalled func(hyperOutportBlock *data.HyperOutportBlock) error
	CloseCalled             func() error
}

// PublishHyperBlock -
func (p *PublisherMock) PublishHyperBlock(hyperOutportBlock *data.HyperOutportBlock) error {
	if p.PublishHyperBlockCalled != nil {
		return p.PublishHyperBlockCalled(hyperOutportBlock)
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
