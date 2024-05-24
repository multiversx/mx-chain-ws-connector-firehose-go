package testscommon

import "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"

// HyperBlockPublisherStub -
type HyperBlockPublisherStub struct {
	PublishHyperBlockCalled func(hyperOutportBlock *hyperOutportBlocks.HyperOutportBlock) error
	CloseCalled             func() error
}

// PublishHyperBlock -
func (h *HyperBlockPublisherStub) PublishHyperBlock(hyperOutportBlock *hyperOutportBlocks.HyperOutportBlock) error {
	if h.PublishHyperBlockCalled != nil {
		return h.PublishHyperBlockCalled(hyperOutportBlock)
	}

	return nil
}

// Close -
func (h *HyperBlockPublisherStub) Close() error {
	if h.CloseCalled != nil {
		return h.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (h *HyperBlockPublisherStub) IsInterfaceNil() bool {
	return h == nil
}
