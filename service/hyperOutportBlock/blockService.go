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
	//TODO: unused for now, placeholder for upcoming PRs
	BlocksPool process.HyperOutportBlocksPool
	Converter  process.OutportBlockConverter
	data.UnimplementedHyperOutportBlockServiceServer
}

// GetHyperOutportBlockByHash retrieves the hyperBlock stored in block pool and converts it to standard proto.
func (bs *Service) GetHyperOutportBlockByHash(ctx context.Context, req *data.BlockHashRequest) (*data.HyperOutportBlock, error) {
	decodeString, err := hex.DecodeString(req.Hash)
	if err != nil {
		return nil, err
	}
	hyperOutportBlock, err := bs.BlocksPool.GetBlock(decodeString)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hyperOutportBlock: %v", err)
	}

	return hyperOutportBlock, nil

}

// GetBlockByNonce retrieve a block from the nonce.
func (bs *Service) GetHyperOutportBlockByNonce(ctx context.Context, req *data.BlockNonceRequest) (*data.HyperOutportBlock, error) {
	return &data.HyperOutportBlock{}, nil
}

//// Register service in the service.
//func (bs *Service) Register(server *grpc.Server) {
//	api.RegisterHyperOutportBlockServiceServer(server, bs)
//}
