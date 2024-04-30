package hyperOutportBlock_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

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

	t.Run("nil context", func(t *testing.T) {
		t.Parallel()

		bs, err := hyperOutportBlock.NewService(nil, &testscommon.GRPCBlocksHandlerStub{})
		require.Nil(t, bs)
		require.Equal(t, process.ErrNilBlockServiceContext, err)
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

// func TestService_Poll(t *testing.T) {
// 	t.Parallel()

// 	bs, err := hyperOutportBlock.NewService(context.TODO(), &testscommon.GRPCBlocksHandlerStub{})
// 	require.Nil(t, err)
// }
