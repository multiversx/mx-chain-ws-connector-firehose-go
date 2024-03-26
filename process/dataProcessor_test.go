package process_test

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	outportcore "github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
	"github.com/stretchr/testify/require"
)

func createOutportBlock() *outportcore.OutportBlock {
	header := &block.Header{
		Nonce:     1,
		PrevHash:  []byte("prev hash"),
		TimeStamp: 100,
	}
	headerBytes, _ := protoMarshaller.Marshal(header)

	return &outportcore.OutportBlock{
		ShardID: 1,
		BlockData: &outportcore.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  string(core.ShardHeaderV1),
			HeaderHash:  []byte("hash"),
		},
	}
}

func createMetaOutportBlock() *outportcore.OutportBlock {
	header := &block.MetaBlock{
		Nonce:     1,
		PrevHash:  []byte("prev hash"),
		TimeStamp: 100,
	}
	headerBytes, _ := protoMarshaller.Marshal(header)

	return &outportcore.OutportBlock{
		ShardID: core.MetachainShardId,
		BlockData: &outportcore.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  string(core.MetaHeader),
			HeaderHash:  []byte("hash"),
		},
	}
}

func TestNewDataProcessor(t *testing.T) {
	t.Parallel()

	t.Run("nil publisher", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			nil,
			&testscommon.MarshallerStub{},
			&testscommon.BlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			createContainer(),
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilPublisher, err)
	})

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			nil,
			&testscommon.BlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			createContainer(),
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilMarshaller, err)
	})

	t.Run("nil blocks pool", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			&testscommon.MarshallerStub{},
			nil,
			&testscommon.DataAggregatorStub{},
			createContainer(),
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilBlocksPool, err)
	})

	t.Run("nil data aggregator", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			&testscommon.MarshallerStub{},
			&testscommon.BlocksPoolStub{},
			nil,
			createContainer(),
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilDataAggregator, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			&testscommon.MarshallerStub{},
			&testscommon.BlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			createContainer(),
		)
		require.Nil(t, err)
		require.False(t, dp.IsInterfaceNil())
	})
}

func TestDataProcessor_ProcessPayload_NotImplementedTopics(t *testing.T) {
	t.Parallel()

	dp, _ := process.NewDataProcessor(
		&testscommon.PublisherStub{},
		&testscommon.MarshallerStub{},
		&testscommon.BlocksPoolStub{},
		&testscommon.DataAggregatorStub{},
		createContainer(),
	)

	require.Nil(t, dp.ProcessPayload([]byte("payload"), "random topic", 1))
	require.Nil(t, dp.ProcessPayload([]byte("payload"), outportcore.TopicSaveRoundsInfo, 1))
	require.Nil(t, dp.ProcessPayload([]byte("payload"), outportcore.TopicSaveValidatorsRating, 1))
	require.Nil(t, dp.ProcessPayload([]byte("payload"), outportcore.TopicSaveValidatorsPubKeys, 1))
	require.Nil(t, dp.ProcessPayload([]byte("payload"), outportcore.TopicSaveAccounts, 1))
	require.Nil(t, dp.ProcessPayload([]byte("payload"), outportcore.TopicFinalizedBlock, 1))
}

func TestDataProcessor_ProcessPayload(t *testing.T) {
	t.Parallel()

	t.Run("nil outport block, should return error", func(t *testing.T) {
		t.Parallel()

		dp, _ := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			protoMarshaller,
			&testscommon.BlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			createContainer(),
		)

		err := dp.ProcessPayload(nil, outportcore.TopicSaveBlock, 1)
		require.Equal(t, process.ErrNilOutportBlockData, err)

		outportBlock := createOutportBlock()
		outportBlock.BlockData = nil
		outportBlockBytes, _ := protoMarshaller.Marshal(outportBlock)

		err = dp.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
		require.Equal(t, process.ErrNilOutportBlockData, err)
	})

	t.Run("invalid payload, cannot unmarshall, should return error", func(t *testing.T) {
		t.Parallel()

		dp, _ := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			protoMarshaller,
			&testscommon.BlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			createContainer(),
		)

		err := dp.ProcessPayload([]byte("invalid payload"), outportcore.TopicSaveBlock, 1)
		require.NotNil(t, err)
	})

	t.Run("shard outport block, should work", func(t *testing.T) {
		t.Parallel()

		outportBlock := createOutportBlock()
		outportBlockBytes, _ := protoMarshaller.Marshal(outportBlock)

		putBlockWasCalled := false
		dp, _ := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			protoMarshaller,
			&testscommon.BlocksPoolStub{
				PutBlockCalled: func(hash []byte, outportBlock *outportcore.OutportBlock, round uint64) error {
					putBlockWasCalled = true
					return nil
				},
			},
			&testscommon.DataAggregatorStub{},
			createContainer(),
		)

		err := dp.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
		require.Nil(t, err)

		require.True(t, putBlockWasCalled)
	})

	t.Run("meta outport block, should work", func(t *testing.T) {
		t.Parallel()

		outportBlock := createMetaOutportBlock()
		outportBlockBytes, _ := protoMarshaller.Marshal(outportBlock)

		publishWasCalled := false
		dp, _ := process.NewDataProcessor(
			&testscommon.PublisherStub{
				PublishHyperBlockCalled: func(hyperOutportBlock *data.HyperOutportBlock) error {
					publishWasCalled = true
					return nil
				},
			},
			protoMarshaller,
			&testscommon.BlocksPoolStub{},
			&testscommon.DataAggregatorStub{
				ProcessHyperBlockCalled: func(outportBlock *outportcore.OutportBlock) (*data.HyperOutportBlock, error) {
					return &data.HyperOutportBlock{
						MetaOutportBlock: outportBlock,
					}, nil
				},
			},
			createContainer(),
		)

		err := dp.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
		require.Nil(t, err)

		require.True(t, publishWasCalled)
	})
}

func TestDataProcessor_Close(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")

	dp, err := process.NewDataProcessor(
		&testscommon.PublisherStub{
			CloseCalled: func() error {
				return expectedErr
			},
		},
		&testscommon.MarshallerStub{},
		&testscommon.BlocksPoolStub{},
		&testscommon.DataAggregatorStub{},
		createContainer(),
	)
	require.Nil(t, err)

	err = dp.Close()
	require.Equal(t, expectedErr, err)
}
