package disabled

import (
	"github.com/multiversx/mx-chain-core-go/data/outport"
)

type disabledOutportBlocksPool struct {
}

// NewDisabledOutportBlocksPool will create a new disabled outport block pool component
func NewDisabledOutportBlocksPool() *disabledOutportBlocksPool {
	return &disabledOutportBlocksPool{}
}

// UpdateMetaState does nothing
func (bp *disabledOutportBlocksPool) UpdateMetaState(round uint64) {
}

// PutBlock returns nil
func (bp *disabledOutportBlocksPool) PutBlock(hash []byte, outportBlock *outport.OutportBlock, currentRound uint64) error {
	return nil
}

// GetBlock returns nil
func (bp *disabledOutportBlocksPool) GetBlock(hash []byte) (*outport.OutportBlock, error) {
	return nil, nil
}

// Close will trigger close on blocks pool component
func (bp *disabledOutportBlocksPool) Close() error {
	return nil
}

// IsInterfaceNil returns nil if there is no value under the interface
func (bp *disabledOutportBlocksPool) IsInterfaceNil() bool {
	return bp == nil
}
