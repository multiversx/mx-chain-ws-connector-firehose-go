package hyperOutportBlock_test

import (
	"context"
	"encoding/hex"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"

	data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/service/hyperOutportBlock"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/testscommon"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	t.Run("nil blocks handler", func(t *testing.T) {
		t.Parallel()

		bs, err := hyperOutportBlock.NewService(context.TODO(), nil)
		require.Nil(t, bs)
		require.Equal(t, process.ErrNilGRPCBlocksHandler, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		bs, err := hyperOutportBlock.NewService(context.TODO(), &testscommon.GRPCBlocksHandlerStub{})
		require.Nil(t, err)
		require.NotNil(t, bs)
	})
}

func TestService_GetHyperOutportBlockByHash(t *testing.T) {
	t.Parallel()

	handler := testscommon.GRPCBlocksHandlerStub{
		FetchHyperBlockByHashCalled: func(hash []byte) (*data.HyperOutportBlock, error) {
			return &data.HyperOutportBlock{
				MetaOutportBlock: &data.MetaOutportBlock{
					BlockData: &data.MetaBlockData{
						HeaderHash: hash,
					},
				},
			}, nil
		},
	}
	bs, err := hyperOutportBlock.NewService(context.TODO(), &handler)
	require.NoError(t, err)

	hash := "437a88d24178dea0060afd74f1282c23b34947cf96adcf71cdfa0f3f7bdcdc73"
	expectedHash, _ := hex.DecodeString(hash)

	outportBlock, err := bs.GetHyperOutportBlockByHash(context.Background(), &data.BlockHashRequest{Hash: hash})
	require.NoError(t, err, "couldn't get block by hash")
	require.Equal(t, expectedHash, outportBlock.MetaOutportBlock.BlockData.HeaderHash)
}

func TestService_GetHyperOutportBlockByNonce(t *testing.T) {
	t.Parallel()

	handler := testscommon.GRPCBlocksHandlerStub{
		FetchHyperBlockByNonceCalled: func(nonce uint64) (*data.HyperOutportBlock, error) {
			return &data.HyperOutportBlock{
				MetaOutportBlock: &data.MetaOutportBlock{
					BlockData: &data.MetaBlockData{
						Header: &data.MetaHeader{Nonce: nonce},
					},
				},
			}, nil
		},
	}
	bs, err := hyperOutportBlock.NewService(context.TODO(), &handler)
	require.NoError(t, err)

	nonce := uint64(1)
	outportBlock, err := bs.GetHyperOutportBlockByNonce(context.Background(), &data.BlockNonceRequest{Nonce: nonce})
	require.NoError(t, err, "couldn't get block by nonce")
	require.Equal(t, nonce, outportBlock.MetaOutportBlock.BlockData.Header.Nonce)
}

func TestService_Poll(t *testing.T) {
	t.Parallel()

	t.Run("service context done, should stop polling", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())

		fetchBlockCalled := uint32(0)
		blocksHandler := &testscommon.GRPCBlocksHandlerStub{
			FetchHyperBlockByNonceCalled: func(nonce uint64) (*data.HyperOutportBlock, error) {
				atomic.AddUint32(&fetchBlockCalled, 1)
				return &data.HyperOutportBlock{}, nil
			},
		}

		bs, err := hyperOutportBlock.NewService(ctx, blocksHandler)
		require.Nil(t, err)

		nonce := uint64(10)
		stream := &testscommon.GRPCServerStreamStub{}
		pdDuration := durationpb.New(time.Duration(100) * time.Millisecond)

		wg := &sync.WaitGroup{}
		wg.Add(1)

		go func() {
			err = bs.Poll(nonce, stream, pdDuration)
			require.Nil(t, err)

			wg.Done()
		}()

		time.Sleep(1100 * time.Millisecond)

		cancel()

		wg.Wait()

		require.GreaterOrEqual(t, atomic.LoadUint32(&fetchBlockCalled), uint32(10))
	})

	t.Run("should fail if not able to send to the stream", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		blocksHandler := &testscommon.GRPCBlocksHandlerStub{
			FetchHyperBlockByNonceCalled: func(nonce uint64) (*data.HyperOutportBlock, error) {
				return &data.HyperOutportBlock{}, nil
			},
		}

		bs, err := hyperOutportBlock.NewService(ctx, blocksHandler)
		require.Nil(t, err)

		nonce := uint64(10)

		expectedErr := errors.New("expected error")

		sendBlockCalled := uint32(0)
		stream := &testscommon.GRPCServerStreamStub{
			SendCalled: func(block *data.HyperOutportBlock) error {
				if atomic.LoadUint32(&sendBlockCalled) <= 1 {
					atomic.AddUint32(&sendBlockCalled, 1)
					return nil
				}

				return expectedErr
			},
		}
		pdDuration := durationpb.New(time.Duration(100) * time.Millisecond)

		wg := &sync.WaitGroup{}
		wg.Add(1)

		go func() {
			err = bs.Poll(nonce, stream, pdDuration)
			require.True(t, errors.Is(err, expectedErr))

			wg.Done()
		}()

		time.Sleep(300 * time.Millisecond)

		wg.Wait()

		require.GreaterOrEqual(t, atomic.LoadUint32(&sendBlockCalled), uint32(1))
	})
}

func TestService_HyperOutportBlocksStreamByHash(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := testscommon.GRPCBlocksHandlerStub{
		FetchHyperBlockByHashCalled: func(hash []byte) (*data.HyperOutportBlock, error) {
			return &data.HyperOutportBlock{
				MetaOutportBlock: &data.MetaOutportBlock{
					BlockData: &data.MetaBlockData{
						HeaderHash: hash,
						Header:     &data.MetaHeader{},
					},
				},
			}, nil
		},
	}
	bs, err := hyperOutportBlock.NewService(ctx, &handler)
	require.NoError(t, err)

	streamRequest := &data.BlockHashStreamRequest{
		Hash: hex.EncodeToString([]byte("hash1")),
		PollingInterval: &durationpb.Duration{
			Seconds: 1,
		},
	}

	wasCalled := uint32(0)
	stream := &testscommon.GRPCServerStreamStub{
		SendCalled: func(block *data.HyperOutportBlock) error {
			atomic.AddUint32(&wasCalled, 1)
			return nil
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		err = bs.HyperOutportBlockStreamByHash(streamRequest, stream)
		require.Nil(t, err)

		wg.Done()
	}()

	time.Sleep(1100 * time.Millisecond)

	cancel()

	wg.Wait()

	require.GreaterOrEqual(t, atomic.LoadUint32(&wasCalled), uint32(2))
}

func TestService_HyperOutportBlocksStreamByNonce(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nonce := uint64(10)

	handler := testscommon.GRPCBlocksHandlerStub{
		FetchHyperBlockByNonceCalled: func(nonce uint64) (*data.HyperOutportBlock, error) {
			return &data.HyperOutportBlock{
				MetaOutportBlock: &data.MetaOutportBlock{
					BlockData: &data.MetaBlockData{
						Header: &data.MetaHeader{
							Nonce: nonce,
						},
					},
				},
			}, nil
		},
	}
	bs, err := hyperOutportBlock.NewService(ctx, &handler)
	require.NoError(t, err)

	streamRequest := &data.BlockNonceStreamRequest{
		Nonce: nonce,
		PollingInterval: &durationpb.Duration{
			Seconds: 1,
		},
	}

	wasCalled := uint32(0)
	stream := &testscommon.GRPCServerStreamStub{
		SendCalled: func(block *data.HyperOutportBlock) error {
			atomic.AddUint32(&wasCalled, 1)
			return nil
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		err = bs.HyperOutportBlockStreamByNonce(streamRequest, stream)
		require.Nil(t, err)

		wg.Done()
	}()

	time.Sleep(1100 * time.Millisecond)

	cancel()

	wg.Wait()

	require.GreaterOrEqual(t, atomic.LoadUint32(&wasCalled), uint32(2))
}
