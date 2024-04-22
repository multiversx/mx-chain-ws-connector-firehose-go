package testscommon

import (
	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

type GRPCBlocksHandlerStub struct {
	FetchHyperBlockByHashCalled  func(hash []byte) (*data.HyperOutportBlock, error)
	FetchHyperBlockByNonceCalled func(nonce uint64) (*data.HyperOutportBlock, error)
}

func (g *GRPCBlocksHandlerStub) FetchHyperBlockByHash(hash []byte) (*data.HyperOutportBlock, error) {
	if g.FetchHyperBlockByHashCalled != nil {
		return g.FetchHyperBlockByHashCalled(hash)
	}
	return &data.HyperOutportBlock{}, nil
}

func (g *GRPCBlocksHandlerStub) FetchHyperBlockByNonce(nonce uint64) (*data.HyperOutportBlock, error) {
	if g.FetchHyperBlockByNonceCalled != nil {
		return g.FetchHyperBlockByNonceCalled(nonce)
	}
	return &data.HyperOutportBlock{}, nil
}
