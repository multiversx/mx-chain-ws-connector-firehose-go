package testscommon

import (
	data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

// DataAggregatorMock -
type DataAggregatorMock struct {
	ProcessHyperBlockCalled func(outportBlock *data.MetaOutportBlock) (*data.HyperOutportBlock, error)
}

// ProcessHyperBlock -
func (d *DataAggregatorMock) ProcessHyperBlock(outportBlock *data.MetaOutportBlock) (*data.HyperOutportBlock, error) {
	if d.ProcessHyperBlockCalled != nil {
		return d.ProcessHyperBlockCalled(outportBlock)
	}

	return &data.HyperOutportBlock{}, nil
}

// IsInterfaceNil -
func (d *DataAggregatorMock) IsInterfaceNil() bool {
	return d == nil
}
