package testscommon

import (
	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

// GRPCBlocksHandlerStub -
type GRPCBlocksHandlerStub struct {
	FetchHyperBlockByHashCalled  func(hash []byte) (*data.HyperOutportBlock, error)
	FetchHyperBlockByNonceCalled func(nonce uint64) (*data.HyperOutportBlock, error)
}

// FetchHyperBlockByHash -
func (g *GRPCBlocksHandlerStub) FetchHyperBlockByHash(hash []byte) (*data.HyperOutportBlock, error) {
	if g.FetchHyperBlockByHashCalled != nil {
		return g.FetchHyperBlockByHashCalled(hash)
	}
	return &data.HyperOutportBlock{}, nil
}

// FetchHyperBlockByNonce -
func (g *GRPCBlocksHandlerStub) FetchHyperBlockByNonce(nonce uint64) (*data.HyperOutportBlock, error) {
	if g.FetchHyperBlockByNonceCalled != nil {
		return g.FetchHyperBlockByNonceCalled(nonce)
	}
	return &data.HyperOutportBlock{}, nil
}

func (g *GRPCBlocksHandlerStub) IsInterfaceNil() bool {
	return g == nil
}
