package dummy

import (
	"context"

	"github.com/google/uuid"
	"github.com/multiversx/mx-chain-storage-go/types"
	"google.golang.org/grpc"

	"github.com/multiversx/mx-chain-ws-connector-template-go/api/dummy"
)

// BlockService returns blocks based on nonce or hash from cache.
type BlockService struct {
	cacher types.Cacher
	dummy.UnimplementedBlockServiceServer
}

// GetBlockByHash retrieves a block from the hash.
func (bs *BlockService) GetBlockByHash(ctx context.Context, req *dummy.BlockHashRequest) (*dummy.Block, error) {
	return &dummy.Block{Hash: uuid.Must(uuid.NewUUID()).String(), Nonce: 0}, nil
}

// GetBlockByNonce retrieve a block from the nonce.
func (bs *BlockService) GetBlockByNonce(ctx context.Context, req *dummy.BlockNonceRequest) (*dummy.Block, error) {
	return &dummy.Block{Hash: uuid.Must(uuid.NewUUID()).String(), Nonce: 0}, nil
}

// Register service in the server.
func (bs *BlockService) Register(server *grpc.Server) error {
	dummy.RegisterBlockServiceServer(server, bs)
	return nil
}
