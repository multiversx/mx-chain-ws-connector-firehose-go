package hyperOutportBlock

import (
	"context"

	"google.golang.org/grpc"

	api "github.com/multiversx/mx-chain-ws-connector-template-go/api/hyperOutportBlocks"
)

// Service returns blocks based on nonce or hash from cache.
type Service struct {
	//TODO: unused for now, placeholder for upcoming PRs
	//cacher types.Cacher
	api.UnimplementedHyperOutportBlockServiceServer
}

// GetBlockByHash retrieves a block from the hash.
func (bs *Service) GetBlockByHash(ctx context.Context, req *api.BlockHashRequest) (*api.HyperOutportBlock, error) {
	return &api.HyperOutportBlock{}, nil
}

// GetBlockByNonce retrieve a block from the nonce.
func (bs *Service) GetBlockByNonce(ctx context.Context, req *api.BlockNonceRequest) (*api.HyperOutportBlock, error) {
	return &api.HyperOutportBlock{}, nil
}

// Register service in the service.
func (bs *Service) Register(server *grpc.Server) {
	api.RegisterHyperOutportBlockServiceServer(server, bs)
}
