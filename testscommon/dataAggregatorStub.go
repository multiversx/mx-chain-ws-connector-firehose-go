package testscommon

import (
	data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

// DataAggregatorStub -
type DataAggregatorStub struct {
	ProcessHyperBlockCalled func(outportBlock *data.MetaOutportBlock) (*data.HyperOutportBlock, error)
}

// ProcessHyperBlock -
func (d *DataAggregatorStub) ProcessHyperBlock(outportBlock *data.MetaOutportBlock) (*data.HyperOutportBlock, error) {
	if d.ProcessHyperBlockCalled != nil {
		return d.ProcessHyperBlockCalled(outportBlock)
	}

	return &data.HyperOutportBlock{}, nil
}

// IsInterfaceNil -
func (d *DataAggregatorStub) IsInterfaceNil() bool {
	return d == nil
}
