package dummy

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/multiversx/mx-chain-ws-connector-template-go/api/dummy"
)

type BlockService struct {
	dummy.UnimplementedBlockServiceServer
}

func (bs *BlockService) GetBlockByHash(ctx context.Context, req *dummy.BlockHashRequest) (*dummy.Block, error) {
	return &dummy.Block{Hash: uuid.Must(uuid.NewUUID()).String(), Nonce: 0}, nil
}

func (bs *BlockService) GetBlockByNonce(ctx context.Context, req *dummy.BlockNonceRequest) (*dummy.Block, error) {
	return &dummy.Block{Hash: uuid.Must(uuid.NewUUID()).String(), Nonce: 0}, nil
}

func (bs *BlockService) Register(server *grpc.Server) error {
	dummy.RegisterBlockServiceServer(server, bs)
	return nil
}
