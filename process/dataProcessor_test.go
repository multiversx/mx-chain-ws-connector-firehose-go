package process_test

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	outportcore "github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/stretchr/testify/require"

	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
)

func createOutportBlock() *outportcore.OutportBlock {
	header := &block.Header{
		Nonce:     1,
		PrevHash:  []byte("prev hash"),
		TimeStamp: 100,
	}
	headerBytes, _ := gogoProtoMarshaller.Marshal(header)

	return &outportcore.OutportBlock{
		ShardID: 1,
		BlockData: &outportcore.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  string(core.ShardHeaderV1),
			HeaderHash:  []byte("hash"),
		},
	}
}

func createChainMetaOutportBlock() *outportcore.OutportBlock {
	header := &block.MetaBlock{
		Nonce:     1,
		PrevHash:  []byte("prev hash"),
		TimeStamp: 100,
	}
	headerBytes, _ := gogoProtoMarshaller.Marshal(header)

	return &outportcore.OutportBlock{
		ShardID: core.MetachainShardId,
		BlockData: &outportcore.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  string(core.MetaHeader),
			HeaderHash:  []byte("hash"),
		},
	}
}

func createShardOutportBlock() *data.ShardOutportBlock {
	return &data.ShardOutportBlock{
		ShardID: 1,
		BlockData: &data.BlockData{
			ShardID: 1,
			Header: &data.Header{
				Nonce: 10,
				Round: 10,
			},
			HeaderHash: []byte("hash_shard"),
		},
	}
}

func createMetaOutportBlock() *data.MetaOutportBlock {
	return &data.MetaOutportBlock{
		ShardID: core.MetachainShardId,
		BlockData: &data.MetaBlockData{
			ShardID: 1,
			Header: &data.MetaHeader{
				Nonce: 10,
				Round: 10,
			},
			HeaderHash: []byte("hash_meta"),
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
			&testscommon.HyperBlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			&testscommon.OutportBlockConverterStub{},
			0,
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilPublisher, err)
	})

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			nil,
			&testscommon.HyperBlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			&testscommon.OutportBlockConverterStub{},
			0,
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
			&testscommon.OutportBlockConverterStub{},
			0,
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilHyperBlocksPool, err)
	})

	t.Run("nil data aggregator", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			&testscommon.MarshallerStub{},
			&testscommon.HyperBlocksPoolStub{},
			nil,
			&testscommon.OutportBlockConverterStub{},
			0,
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilDataAggregator, err)
	})

	t.Run("nil block converter", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			&testscommon.MarshallerStub{},
			&testscommon.HyperBlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			nil,
			0,
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilOutportBlocksConverter, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			&testscommon.MarshallerStub{},
			&testscommon.HyperBlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			&testscommon.OutportBlockConverterStub{},
			0,
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
		&testscommon.HyperBlocksPoolStub{},
		&testscommon.DataAggregatorStub{},
		&testscommon.OutportBlockConverterStub{},
		0,
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

	outportBlockConverter := process.NewOutportBlockConverter(gogoProtoMarshaller, protoMarshaller)

	t.Run("nil outport block data, should return error", func(t *testing.T) {
		t.Parallel()

		dp, _ := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			gogoProtoMarshaller,
			&testscommon.HyperBlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			&testscommon.OutportBlockConverterStub{},
			0,
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
			&testscommon.HyperBlocksPoolStub{},
			&testscommon.DataAggregatorStub{},
			&testscommon.OutportBlockConverterStub{},
			0,
		)

		err := dp.ProcessPayload([]byte("invalid payload"), outportcore.TopicSaveBlock, 1)
		require.NotNil(t, err)
	})

	t.Run("shard outport block, should work", func(t *testing.T) {
		t.Parallel()

		outportBlock := createOutportBlock()
		outportBlockBytes, _ := gogoProtoMarshaller.Marshal(outportBlock)

		putBlockWasCalled := false
		dp, _ := process.NewDataProcessor(
			&testscommon.PublisherStub{},
			gogoProtoMarshaller,
			&testscommon.HyperBlocksPoolStub{
				PutShardBlockCalled: func(hash []byte, outportBlock *data.ShardOutportBlock) error {
					putBlockWasCalled = true
					return nil
				},
			},
			&testscommon.DataAggregatorStub{},
			outportBlockConverter,
			0,
		)

		err := dp.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
		require.Nil(t, err)

		require.True(t, putBlockWasCalled)
	})

	t.Run("meta outport block, should work", func(t *testing.T) {
		t.Parallel()

		outportBlock := createChainMetaOutportBlock()
		outportBlockBytes, err := gogoProtoMarshaller.Marshal(outportBlock)
		require.Nil(t, err)

		metaOutportBlock, err := outportBlockConverter.HandleMetaOutportBlock(outportBlock)
		require.Nil(t, err)

		publishWasCalled := false
		dp, _ := process.NewDataProcessor(
			&testscommon.PublisherStub{
				PublishHyperBlockCalled: func(hyperOutportBlock *data.HyperOutportBlock) error {
					publishWasCalled = true
					return nil
				},
			},
			gogoProtoMarshaller,
			&testscommon.HyperBlocksPoolStub{},
			&testscommon.DataAggregatorStub{
				ProcessHyperBlockCalled: func(outportBlock *data.MetaOutportBlock) (*data.HyperOutportBlock, error) {
					return &data.HyperOutportBlock{
						MetaOutportBlock: metaOutportBlock,
					}, nil
				},
			},
			outportBlockConverter,
			0,
		)

		err = dp.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
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
		&testscommon.HyperBlocksPoolStub{},
		&testscommon.DataAggregatorStub{},
		&testscommon.OutportBlockConverterStub{},
		0,
	)
	require.Nil(t, err)

	err = dp.Close()
	require.Equal(t, expectedErr, err)
}
