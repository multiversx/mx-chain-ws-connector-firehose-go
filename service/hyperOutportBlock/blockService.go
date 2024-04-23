package hyperOutportBlock

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/multiversx/mx-chain-core-go/core/check"

	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
)

// Service returns blocks based on nonce or hash from cache.
type Service struct {
	blocksHandler process.GRPCBlocksHandler
	queue         process.HyperOutportBlocksQueue
	data.UnimplementedHyperOutportBlockServiceServer
}

// NewService returns a new instance of the hyperOutportBlock service.
func NewService(blocksHandler process.GRPCBlocksHandler, queue process.HyperOutportBlocksQueue) (*Service, error) {
	if check.IfNil(blocksHandler) {
		return nil, process.ErrNilOutportBlockData
	}

	return &Service{blocksHandler: blocksHandler, queue: queue}, nil
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

func (bs *Service) HyperOutportBlockStream(_ *empty.Empty, stream data.HyperOutportBlockService_HyperOutportBlockStreamServer) error {
	for {
		select {
		// Exit on stream context done
		case <-stream.Context().Done():
			return nil
		default:
			block, err := bs.queue.Dequeue()
			if err != nil {
				return fmt.Errorf("failed to retrieve hyperOutportBlock: %w", err)
			}

			err = stream.Send(block)
			if err != nil {
				return fmt.Errorf("failed to send hyperOutportBlock: %w", err)
			}
		}
	}
}
