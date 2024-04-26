package hyperOutportBlock

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"google.golang.org/grpc"

	data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
)

var (
	log = logger.GetOrCreate("service")
)

type serverStream interface {
	grpc.ServerStream
	Send(*data.HyperOutportBlock) error
}

// Service returns blocks based on nonce or hash from cache.
type Service struct {
	ctx           context.Context
	blocksHandler process.GRPCBlocksHandler

	data.UnimplementedHyperOutportBlockServiceServer
}

// NewService returns a new instance of the hyperOutportBlock service.
func NewService(ctx context.Context, blocksHandler process.GRPCBlocksHandler) (*Service, error) {
	if check.IfNil(blocksHandler) {
		return nil, process.ErrNilOutportBlockData
	}

	if ctx == nil {
		return nil, process.ErrNilBlockServiceContext
	}

	return &Service{ctx: ctx, blocksHandler: blocksHandler}, nil
}

// GetHyperOutportBlockByHash retrieves the hyperBlock stored in block pool and converts it to standard proto.
func (bs *Service) GetHyperOutportBlockByHash(ctx context.Context, req *data.BlockHashRequest) (*data.HyperOutportBlock, error) {
	return bs.fetchBlockByHash(req.Hash)
}

// GetHyperOutportBlockByNonce retrieve a block from the nonce.
func (bs *Service) GetHyperOutportBlockByNonce(ctx context.Context, req *data.BlockNonceRequest) (*data.HyperOutportBlock, error) {
	return bs.fetchBlockByNonce(req.Nonce)
}

// HyperOutportBlockStreamByHash will return a stream on which the incoming hyperBlocks are being sent.
func (bs *Service) HyperOutportBlockStreamByHash(req *data.BlockHashStreamRequest, stream data.HyperOutportBlockService_HyperOutportBlockStreamByHashServer) error {
	hyperOutportBlock, err := bs.fetchBlockByHash(req.Hash)
	if err != nil {
		return err
	}

	// send the initial hyper outport block
	err = stream.Send(hyperOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to send stream to hyperOutportBlock: %w", err)
	}

	// start polling and retrieve the starting nonce
	nonce := hyperOutportBlock.MetaOutportBlock.BlockData.Header.Nonce + 1
	err = bs.poll(nonce, stream, req.PollingInterval)
	if err != nil {
		return fmt.Errorf("failure encountered while polling for hyper blocks: %w", err)
	}

	return nil
}

// HyperOutportBlockStreamByNonce will return a stream on which the incoming hyperBlocks are being sent.
func (bs *Service) HyperOutportBlockStreamByNonce(req *data.BlockNonceStreamRequest, stream data.HyperOutportBlockService_HyperOutportBlockStreamByNonceServer) error {
	hyperOutportBlock, err := bs.fetchBlockByNonce(req.Nonce)
	if err != nil {
		return err
	}

	// send the initial hyper outport block
	err = stream.Send(hyperOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to send stream to hyperOutportBlock: %w", err)
	}

	// start polling and retrieve the starting nonce
	nonce := hyperOutportBlock.MetaOutportBlock.BlockData.Header.Nonce + 1
	err = bs.poll(nonce, stream, req.PollingInterval)
	if err != nil {
		return fmt.Errorf("failure encountered while polling for hyper blocks: %w", err)
	}

	return nil
}

func (bs *Service) fetchBlockByNonce(nonce uint64) (*data.HyperOutportBlock, error) {
	hyperOutportBlock, err := bs.blocksHandler.FetchHyperBlockByNonce(nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hyperOutportBlock with nonce '%d': %w", nonce, err)
	}

	return hyperOutportBlock, nil
}

func (bs *Service) fetchBlockByHash(hash string) (*data.HyperOutportBlock, error) {
	decodeString, err := hex.DecodeString(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}
	hyperOutportBlock, err := bs.blocksHandler.FetchHyperBlockByHash(decodeString)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hyperOutportBlock with hash %q: %w", hash, err)
	}

	return hyperOutportBlock, nil
}

func (bs *Service) poll(nonce uint64, stream serverStream, pollingInterval *duration.Duration) error {
	// parse the provided pollingInterval
	timeDuration, err := protoToTimeDuration(pollingInterval)
	if err != nil {
		return fmt.Errorf("invalid polling interval: %w", err)
	}
	ticker := time.NewTicker(timeDuration)

	for {
		select {
		case <-bs.ctx.Done():
			ticker.Stop()
			return bs.ctx.Err()

		case <-stream.Context().Done():
			ticker.Stop()
			return stream.Context().Err()

		case <-ticker.C:
			// fetch the next hyperBlock.
			hb, fetchErr := bs.blocksHandler.FetchHyperBlockByNonce(nonce)
			if fetchErr != nil {
				// if the hyperBlock was not found. try again in the next iteration.
				log.Error(fmt.Errorf("failed to retrieve hyper block with nonce '%d': %w", nonce, fetchErr).Error())
				continue
			}

			// if found, send it on the stream.
			if sendErr := stream.Send(hb); sendErr != nil {
				return fmt.Errorf("failed to send hyperOutportBlock: %w", sendErr)
			}

			nonce++
		}
	}
}

func protoToTimeDuration(protoDur *duration.Duration) (time.Duration, error) {
	seconds := strconv.FormatInt(protoDur.GetSeconds(), 10)    // Convert int64 to string
	nanos := strconv.FormatInt(int64(protoDur.GetNanos()), 10) // Convert int32 to string
	return time.ParseDuration(seconds + "s" + nanos + "ns")
}
