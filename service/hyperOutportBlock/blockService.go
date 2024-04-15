package hyperOutportBlock

import (
	"context"
	"encoding/hex"
	"fmt"

	api "github.com/multiversx/mx-chain-ws-connector-template-go/api/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
)

// Service returns blocks based on nonce or hash from cache.
type Service struct {
	//TODO: unused for now, placeholder for upcoming PRs
	BlocksPool process.HyperOutportBlocksPool
	Converter  process.OutportBlockConverter
	api.UnimplementedHyperOutportBlockServiceServer
}

// GetHyperOutportBlockByHash retrieves the hyperBlock stored in block pool and converts it to standard proto.
func (bs *Service) GetHyperOutportBlockByHash(ctx context.Context, req *api.BlockHashRequest) (*api.HyperOutportBlock, error) {
	decodeString, err := hex.DecodeString(req.Hash)
	if err != nil {
		return nil, err
	}
	hyperOutportBlock, err := bs.BlocksPool.GetBlock(decodeString)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hyperOutportBlock: %v", err)
	}

	metaOutportBlock, err := bs.Converter.HandleMetaOutportBlock(hyperOutportBlock.MetaOutportBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to metaOutportBlock: %v", err)
	}

	notarizedHeadersOutportData := make([]*api.NotarizedHeaderOutportData, len(hyperOutportBlock.NotarizedHeadersOutportData))
	for _, h := range hyperOutportBlock.NotarizedHeadersOutportData {
		shardOutportBlock, err := bs.Converter.HandleShardOutportBlock(h.OutportBlock)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to shardOutportBlock: %v", err)
		}

		n := &api.NotarizedHeaderOutportData{
			NotarizedHeaderHash: h.NotarizedHeaderHash,
			OutportBlock:        shardOutportBlock,
		}

		notarizedHeadersOutportData = append(notarizedHeadersOutportData, n)
	}

	return &api.HyperOutportBlock{
		OutportBlock:                metaOutportBlock,
		NotarizedHeadersOutportData: notarizedHeadersOutportData,
	}, nil

}

// GetBlockByNonce retrieve a block from the nonce.
func (bs *Service) GetHyperOutportBlockByNonce(ctx context.Context, req *api.BlockNonceRequest) (*api.HyperOutportBlock, error) {
	return &api.HyperOutportBlock{}, nil
}

//// Register service in the service.
//func (bs *Service) Register(server *grpc.Server) {
//	api.RegisterHyperOutportBlockServiceServer(server, bs)
//}
