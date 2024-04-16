package disabled

import (
	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

type disabledHyperOutportBlocksPool struct {
}

// NewDisabledHyperOutportBlocksPool will create a new disabled blocks pool component
func NewDisabledHyperOutportBlocksPool() *disabledHyperOutportBlocksPool {
	return &disabledHyperOutportBlocksPool{}
}

// UpdateMetaState does nothing
func (bp *disabledHyperOutportBlocksPool) UpdateMetaState(round uint64) {
}

// PutBlock returns nil
func (bp *disabledHyperOutportBlocksPool) PutBlock(hash []byte, outportBlock *data.HyperOutportBlock, currentRound uint64) error {
	return nil
}

// GetBlock returns nil
func (bp *disabledHyperOutportBlocksPool) GetBlock(hash []byte) (*data.HyperOutportBlock, error) {
	return nil, nil
}

func (bp *disabledHyperOutportBlocksPool) Close() error {
	return nil
}

// IsInterfaceNil returns nil if there is no value under the interface
func (bp *disabledHyperOutportBlocksPool) IsInterfaceNil() bool {
	return bp == nil
}
