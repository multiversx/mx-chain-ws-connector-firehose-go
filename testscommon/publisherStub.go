package testscommon

import data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"

// PublisherStub -
type PublisherStub struct {
	PublishHyperBlockCalled func(hyperOutportBlock *data.HyperOutportBlock) error
	CloseCalled             func() error
}

// PublishHyperBlock -
func (p *PublisherStub) PublishHyperBlock(hyperOutportBlock *data.HyperOutportBlock) error {
	if p.PublishHyperBlockCalled != nil {
		return p.PublishHyperBlockCalled(hyperOutportBlock)
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
