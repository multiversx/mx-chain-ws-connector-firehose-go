package testscommon

import (
	data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

// GRPCBlocksHandlerMock -
type GRPCBlocksHandlerMock struct {
	FetchHyperBlockByHashCalled  func(hash []byte) (*data.HyperOutportBlock, error)
	FetchHyperBlockByNonceCalled func(nonce uint64) (*data.HyperOutportBlock, error)
}

// FetchHyperBlockByHash -
func (g *GRPCBlocksHandlerMock) FetchHyperBlockByHash(hash []byte) (*data.HyperOutportBlock, error) {
	if g.FetchHyperBlockByHashCalled != nil {
		return g.FetchHyperBlockByHashCalled(hash)
	}
	return &data.HyperOutportBlock{}, nil
}

// FetchHyperBlockByNonce -
func (g *GRPCBlocksHandlerMock) FetchHyperBlockByNonce(nonce uint64) (*data.HyperOutportBlock, error) {
	if g.FetchHyperBlockByNonceCalled != nil {
		return g.FetchHyperBlockByNonceCalled(nonce)
	}
	return &data.HyperOutportBlock{}, nil
}

// IsInterfaceNil -
func (g *GRPCBlocksHandlerMock) IsInterfaceNil() bool {
	return g == nil
}
