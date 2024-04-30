package process_test

import (
	"testing"

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

	metaOutportBlock := createMetaOutportBlock()
	hyperOutportBlock := createHyperOutportBlock()
	hyperOutportBlock.MetaOutportBlock = metaOutportBlock

	ph, err := process.NewPublisherHandler(
		&testscommon.HyperBlockPublisherStub{},
		&testscommon.HyperBlocksPoolStub{
			GetMetaBlockCalled: func(hash []byte) (*hyperOutportBlocks.MetaOutportBlock, error) {
				return metaOutportBlock, nil
			},
		},
		&testscommon.DataAggregatorStub{
			ProcessHyperBlockCalled: func(outportBlock *hyperOutportBlocks.MetaOutportBlock) (*hyperOutportBlocks.HyperOutportBlock, error) {
				return hyperOutportBlock, nil
			},
		},
		defaultRetryDuration,
		defaultFirstCommitableBlock,
	)
	require.Nil(t, err)

	err = ph.PublishBlock([]byte("headerHash"))
	require.Nil(t, err)
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
