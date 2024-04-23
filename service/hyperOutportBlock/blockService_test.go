package hyperOutportBlock

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
)

func TestService_GetHyperOutportBlockByHash(t *testing.T) {
	t.Parallel()
	queue := process.NewHyperOutportBlocksQueue()
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
	bs, err := NewService(&handler, queue)
	require.NoError(t, err)
	hash := "437a88d24178dea0060afd74f1282c23b34947cf96adcf71cdfa0f3f7bdcdc73"
	expectedHash, _ := hex.DecodeString(hash)
	outportBlock, err := bs.GetHyperOutportBlockByHash(context.Background(), &data.BlockHashRequest{Hash: hash})
	require.NoError(t, err, "couldn't get block by hash")
	require.Equal(t, expectedHash, outportBlock.MetaOutportBlock.BlockData.HeaderHash)
}

func TestService_GetHyperOutportBlockByNonce(t *testing.T) {
	t.Parallel()
	queue := process.NewHyperOutportBlocksQueue()
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
	bs, err := NewService(&handler, queue)
	require.NoError(t, err)
	nonce := uint64(1)
	outportBlock, err := bs.GetHyperOutportBlockByNonce(context.Background(), &data.BlockNonceRequest{Nonce: nonce})
	require.NoError(t, err, "couldn't get block by nonce")
	require.Equal(t, nonce, outportBlock.MetaOutportBlock.BlockData.Header.Nonce)
}

func Test111(t *testing.T) {

	a := make([]int, 0)
	i := 0

	go func() {
		for {
			a = append(a, i)
			i++
		}
	}()

	time.Sleep(5 * time.Second)

	for _, j := range a {
		fmt.Println(j)
	}
}
