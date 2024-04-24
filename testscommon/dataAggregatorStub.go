package testscommon

import (
	"github.com/multiversx/mx-chain-core-go/data/outport"

	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

// DataAggregatorStub -
type DataAggregatorStub struct {
	ProcessHyperBlockCalled func(outportBlock *outport.OutportBlock) (*data.HyperOutportBlock, error)
}

// ProcessHyperBlock -
func (d *DataAggregatorStub) ProcessHyperBlock(outportBlock *outport.OutportBlock) (*data.HyperOutportBlock, error) {
	if d.ProcessHyperBlockCalled != nil {
		return d.ProcessHyperBlockCalled(outportBlock)
	}

	return &data.HyperOutportBlock{}, nil
}

// IsInterfaceNil -
func (d *DataAggregatorStub) IsInterfaceNil() bool {
	return d == nil
}
