package process_test

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	outportcore "github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/stretchr/testify/require"

	data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/testscommon"
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
			&testscommon.HyperBlocksPoolMock{},
			&testscommon.OutportBlockConverterMock{},
			defaultFirstCommitableBlocks,
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilPublisher, err)
	})

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherMock{},
			nil,
			&testscommon.HyperBlocksPoolMock{},
			&testscommon.OutportBlockConverterMock{},
			defaultFirstCommitableBlocks,
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilMarshaller, err)
	})

	t.Run("nil blocks pool", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherMock{},
			&testscommon.MarshallerMock{},
			nil,
			&testscommon.OutportBlockConverterMock{},
			defaultFirstCommitableBlocks,
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilBlocksPool, err)
	})

	t.Run("nil block converter", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherMock{},
			&testscommon.MarshallerStub{},
			&testscommon.HyperBlocksPoolMock{},
			nil,
			defaultFirstCommitableBlocks,
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilOutportBlocksConverter, err)
	})

	t.Run("nil first commitable blocks", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherMock{},
			&testscommon.MarshallerStub{},
			&testscommon.HyperBlocksPoolMock{},
			&testscommon.OutportBlockConverterMock{},
			nil,
		)
		require.Nil(t, dp)
		require.Equal(t, process.ErrNilFirstCommitableBlocks, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		dp, err := process.NewDataProcessor(
			&testscommon.PublisherMock{},
			&testscommon.MarshallerStub{},
			&testscommon.HyperBlocksPoolMock{},
			&testscommon.OutportBlockConverterMock{},
			defaultFirstCommitableBlocks,
		)
		require.Nil(t, err)
		require.False(t, dp.IsInterfaceNil())
	})
}

func TestDataProcessor_ProcessPayload_NotImplementedTopics(t *testing.T) {
	t.Parallel()

	dp, _ := process.NewDataProcessor(
		&testscommon.PublisherMock{},
		&testscommon.MarshallerStub{},
		&testscommon.HyperBlocksPoolMock{},
		&testscommon.OutportBlockConverterMock{},
		defaultFirstCommitableBlocks,
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

	outportBlockConverter, err := process.NewOutportBlockConverter(gogoProtoMarshaller, protoMarshaller)
	require.Nil(t, err)

	t.Run("nil outport block data, should return error", func(t *testing.T) {
		t.Parallel()

		dp, _ := process.NewDataProcessor(
			&testscommon.PublisherMock{},
			gogoProtoMarshaller,
			&testscommon.HyperBlocksPoolMock{},
			&testscommon.OutportBlockConverterMock{},
			defaultFirstCommitableBlocks,
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
			&testscommon.PublisherMock{},
			protoMarshaller,
			&testscommon.HyperBlocksPoolMock{},
			&testscommon.OutportBlockConverterMock{},
			defaultFirstCommitableBlocks,
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
			&testscommon.PublisherMock{},
			gogoProtoMarshaller,
			&testscommon.HyperBlocksPoolMock{
				PutBlockCalled: func(hash []byte, outportBlock process.OutportBlockHandler) error {
					putBlockWasCalled = true
					return nil
				},
			},
			outportBlockConverter,
			defaultFirstCommitableBlocks,
		)

		err := dp.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
		require.Nil(t, err)

		require.True(t, putBlockWasCalled)
	})

	t.Run("should not publish block until first commitable block", func(t *testing.T) {
		t.Parallel()

		round := uint64(10)

		firstCommitableBlocks := map[uint32]uint64{
			core.MetachainShardId: round + 1,
		}

		outportBlock := createChainMetaOutportBlock()
		outportBlockBytes, err := gogoProtoMarshaller.Marshal(outportBlock)
		require.Nil(t, err)

		metaOutportBlock, err := outportBlockConverter.HandleMetaOutportBlock(outportBlock)
		require.Nil(t, err)

		publishWasCalled := false
		dp, _ := process.NewDataProcessor(
			&testscommon.PublisherMock{
				PublishBlockCalled: func(headerHash []byte) error {
					publishWasCalled = true
					return nil
				},
			},
			gogoProtoMarshaller,
			&testscommon.HyperBlocksPoolMock{
				PutBlockCalled: func(hash []byte, outportBlock process.OutportBlockHandler) error {
					require.Equal(t, metaOutportBlock, outportBlock)

					return nil
				},
			},
			outportBlockConverter,
			firstCommitableBlocks,
		)

		err = dp.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
		require.Nil(t, err)

		require.False(t, publishWasCalled)
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
			&testscommon.PublisherMock{
				PublishBlockCalled: func(headerHash []byte) error {
					publishWasCalled = true
					return nil
				},
			},
			gogoProtoMarshaller,
			&testscommon.HyperBlocksPoolMock{
				PutBlockCalled: func(hash []byte, outportBlock process.OutportBlockHandler) error {
					require.Equal(t, metaOutportBlock, outportBlock)

					return nil
				},
			},
			outportBlockConverter,
			defaultFirstCommitableBlocks,
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
		&testscommon.PublisherMock{
			CloseCalled: func() error {
				return expectedErr
			},
		},
		&testscommon.MarshallerStub{},
		&testscommon.HyperBlocksPoolMock{},
		&testscommon.OutportBlockConverterMock{},
		defaultFirstCommitableBlocks,
	)
	require.Nil(t, err)

	err = dp.Close()
	require.Equal(t, expectedErr, err)
}
