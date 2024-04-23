package hyperOutportBlock

import (
	"context"
	"encoding/hex"
	"fmt"

	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
)

// Service returns blocks based on nonce or hash from cache.
type Service struct {
	blocksHandler process.GRPCBlocksHandler
	data.UnimplementedHyperOutportBlockServiceServer
}

func NewService(blocksHandler process.GRPCBlocksHandler) (*Service, error) {
	if check.IfNil(blocksHandler) {
		return nil, process.ErrNilOutportBlockData
	}

	return &Service{blocksHandler: blocksHandler}, nil
}

// GetHyperOutportBlockByHash retrieves the hyperBlock stored in block pool and converts it to standard proto.
func (bs *Service) GetHyperOutportBlockByHash(ctx context.Context, req *data.BlockHashRequest) (*data.HyperOutportBlock, error) {
	decodeString, err := hex.DecodeString(req.Hash)
	if err != nil {
		return nil, err
	}
	hyperOutportBlock, err := bs.blocksHandler.FetchHyperBlockByHash(decodeString)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hyperOutportBlock: %v", err)
	}

	return hyperOutportBlock, nil
}

// GetHyperOutportBlockByNonce retrieve a block from the nonce.
func (bs *Service) GetHyperOutportBlockByNonce(ctx context.Context, req *data.BlockNonceRequest) (*data.HyperOutportBlock, error) {
	hyperOutportBlock, err := bs.blocksHandler.FetchHyperBlockByNonce(req.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hyperOutportBlock: %v", err)
	}

	return hyperOutportBlock, nil
}
