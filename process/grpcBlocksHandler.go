package process

import (
	"encoding/binary"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
)

type grpcBlocksHandler struct {
	outportBlocksPool DataPool
	dataAggregator    DataAggregator
}

// NewGrpcBlocksHandler will create a new grpc blocks handler component able to fetch hyper outport blocks data to blocks pool
// which will then be consumed by the grpc server
func NewGrpcBlocksHandler(
	outportBlocksPool DataPool,
	blockCreator BlockContainerHandler,
	marshaller marshal.Marshalizer,
	dataAggregator DataAggregator,
) (*grpcBlocksHandler, error) {
	if check.IfNil(outportBlocksPool) {
		return nil, ErrNilBlocksPool
	}
	if check.IfNil(dataAggregator) {
		return nil, ErrNilDataAggregator
	}

	return &grpcBlocksHandler{
		outportBlocksPool: outportBlocksPool,
		dataAggregator:    dataAggregator,
	}, nil
}

// FetchHyperBlockByHash will fetch hyper block from pool by hash
func (gb *grpcBlocksHandler) FetchHyperBlockByHash(hash []byte) (*data.HyperOutportBlock, error) {
	metaOutportBlock, err := gb.outportBlocksPool.GetBlock(hash)
	if err != nil {
		return nil, err
	}

	hyperOutportBlock, err := gb.dataAggregator.ProcessHyperBlock(metaOutportBlock)
	if err != nil {
		return nil, err
	}

	return hyperOutportBlock, nil
}

// FetchHyperBlockByNonce will fetch hyper block from pool by nonce
func (gb *grpcBlocksHandler) FetchHyperBlockByNonce(nonce uint64) (*data.HyperOutportBlock, error) {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, nonce)

	hash, err := gb.outportBlocksPool.Get(nonceBytes)
	if err != nil {
		return nil, err
	}

	return gb.FetchHyperBlockByHash(hash)
}

// IsInterfaceNil checks if the underlying pointer is nil
func (gb *grpcBlocksHandler) IsInterfaceNil() bool {
	return gb == nil
}
