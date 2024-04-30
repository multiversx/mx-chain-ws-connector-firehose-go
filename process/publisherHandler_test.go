package process_test

import (
	"sync/atomic"
	"testing"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/testscommon"
	"github.com/stretchr/testify/require"
)

const (
	defaultRetryDuration        = 100
	defaultFirstCommitableBlock = 0
)

func TestNewPublisherHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil publisher", func(t *testing.T) {
		t.Parallel()

		ph, err := process.NewPublisherHandler(
			nil,
			&testscommon.HyperBlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			defaultRetryDuration,
			defaultFirstCommitableBlock,
		)
		require.Nil(t, ph)
		require.Equal(t, process.ErrNilPublisher, err)
	})

	t.Run("nil outport blocks pool", func(t *testing.T) {
		t.Parallel()

		ph, err := process.NewPublisherHandler(
			&testscommon.HyperBlockPublisherStub{},
			nil,
			&testscommon.DataAggregatorStub{},
			defaultRetryDuration,
			defaultFirstCommitableBlock,
		)
		require.Nil(t, ph)
		require.Equal(t, process.ErrNilHyperBlocksPool, err)
	})

	t.Run("nil data aggregator", func(t *testing.T) {
		t.Parallel()

		ph, err := process.NewPublisherHandler(
			&testscommon.HyperBlockPublisherStub{},
			&testscommon.HyperBlocksPoolStub{},
			nil,
			defaultRetryDuration,
			defaultFirstCommitableBlock,
		)
		require.Nil(t, ph)
		require.Equal(t, process.ErrNilDataAggregator, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		ph, err := process.NewPublisherHandler(
			&testscommon.HyperBlockPublisherStub{},
			&testscommon.HyperBlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			defaultRetryDuration,
			defaultFirstCommitableBlock,
		)
		require.Nil(t, err)
		require.False(t, ph.IsInterfaceNil())
	})
}

func TestPublisherHandler_PublishBlock(t *testing.T) {
	t.Parallel()

	t.Run("should not process hyper block until first commitable block", func(t *testing.T) {
		t.Parallel()

		round := uint64(10)

		metaOutportBlock := createMetaOutportBlock()
		metaOutportBlock.BlockData.Header.Round = round
		hyperOutportBlock := createHyperOutportBlock()
		hyperOutportBlock.MetaOutportBlock = metaOutportBlock

		updateMetaStateCalled := uint32(0)
		ph, err := process.NewPublisherHandler(
			&testscommon.HyperBlockPublisherStub{},
			&testscommon.HyperBlocksPoolStub{
				GetMetaBlockCalled: func(hash []byte) (*hyperOutportBlocks.MetaOutportBlock, error) {
					return metaOutportBlock, nil
				},
				UpdateMetaStateCalled: func(checkpoint *data.BlockCheckpoint) error {
					atomic.AddUint32(&updateMetaStateCalled, 1)
					return nil
				},
			},
			&testscommon.DataAggregatorStub{
				ProcessHyperBlockCalled: func(outportBlock *hyperOutportBlocks.MetaOutportBlock) (*hyperOutportBlocks.HyperOutportBlock, error) {
					require.Fail(t, "should not have been called")
					return nil, nil
				},
			},
			defaultRetryDuration,
			round+1,
		)
		require.Nil(t, err)

		err = ph.PublishBlock([]byte("headerHash"))
		require.Nil(t, err)

		require.GreaterOrEqual(t, atomic.LoadUint32(&updateMetaStateCalled), uint32(1))
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		round := uint64(10)

		metaOutportBlock := createMetaOutportBlock()
		metaOutportBlock.BlockData.Header.Round = round
		hyperOutportBlock := createHyperOutportBlock()
		hyperOutportBlock.MetaOutportBlock = metaOutportBlock

		updateMetaStateCalled := uint32(0)
		ph, err := process.NewPublisherHandler(
			&testscommon.HyperBlockPublisherStub{},
			&testscommon.HyperBlocksPoolStub{
				GetMetaBlockCalled: func(hash []byte) (*hyperOutportBlocks.MetaOutportBlock, error) {
					return metaOutportBlock, nil
				},
				UpdateMetaStateCalled: func(checkpoint *data.BlockCheckpoint) error {
					atomic.AddUint32(&updateMetaStateCalled, 1)
					return nil
				},
			},
			&testscommon.DataAggregatorStub{
				ProcessHyperBlockCalled: func(outportBlock *hyperOutportBlocks.MetaOutportBlock) (*hyperOutportBlocks.HyperOutportBlock, error) {
					return hyperOutportBlock, nil
				},
			},
			defaultRetryDuration,
			round-1,
		)
		require.Nil(t, err)

		err = ph.PublishBlock([]byte("headerHash"))
		require.Nil(t, err)

		require.GreaterOrEqual(t, atomic.LoadUint32(&updateMetaStateCalled), uint32(1))
	})
}

func TestPublisherHandler_Close(t *testing.T) {
	t.Parallel()

	publisherCloseCalled := false
	ph, err := process.NewPublisherHandler(
		&testscommon.HyperBlockPublisherStub{
			CloseCalled: func() error {
				publisherCloseCalled = true
				return nil
			},
		},
		&testscommon.HyperBlocksPoolStub{},
		&testscommon.DataAggregatorStub{},
		defaultRetryDuration,
		defaultFirstCommitableBlock,
	)
	require.Nil(t, err)

	err = ph.Close()
	require.Nil(t, err)

	require.True(t, publisherCloseCalled)
}
