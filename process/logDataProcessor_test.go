package process

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	outportcore "github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
	"github.com/stretchr/testify/require"
)

var protoMarshaller = &marshal.GogoProtoMarshalizer{}

func createOutportBlock() *outportcore.OutportBlock {
	header := &block.Header{
		Nonce:     1,
		PrevHash:  []byte("prev hash"),
		TimeStamp: 100,
	}
	headerBytes, _ := protoMarshaller.Marshal(header)

	return &outportcore.OutportBlock{
		BlockData: &outportcore.BlockData{
			HeaderHash:  []byte("hash"),
			HeaderBytes: headerBytes,
			HeaderType:  string(core.ShardHeaderV1),
		},
	}
}

func TestNewLogDataProcessor(t *testing.T) {
	t.Parallel()

	t.Run("nil io writer, should return error", func(t *testing.T) {
		t.Parallel()

		fi, err := NewLogDataProcessor(nil, block.NewEmptyBlockCreatorsContainer(), protoMarshaller)
		require.Nil(t, fi)
		require.Equal(t, errNilWriter, err)
	})

	t.Run("nil block creator, should return error", func(t *testing.T) {
		t.Parallel()

		fi, err := NewLogDataProcessor(&testscommon.IoWriterStub{}, nil, protoMarshaller)
		require.Nil(t, fi)
		require.Equal(t, errNilBlockCreator, err)
	})

	t.Run("nil marshaller, should return error", func(t *testing.T) {
		t.Parallel()

		fi, err := NewLogDataProcessor(&testscommon.IoWriterStub{}, block.NewEmptyBlockCreatorsContainer(), nil)
		require.Nil(t, fi)
		require.Equal(t, errNilMarshaller, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		fi, err := NewLogDataProcessor(&testscommon.IoWriterStub{}, block.NewEmptyBlockCreatorsContainer(), protoMarshaller)
		require.Nil(t, err)
		require.False(t, check.IfNil(fi))
	})
}

func TestFirehoseIndexer_SaveBlock(t *testing.T) {
	t.Parallel()

	t.Run("nil outport block, should return error", func(t *testing.T) {
		t.Parallel()

		fi, _ := NewLogDataProcessor(&testscommon.IoWriterStub{}, block.NewEmptyBlockCreatorsContainer(), protoMarshaller)

		err := fi.ProcessPayload(nil, outportcore.TopicSaveBlock)
		require.Equal(t, errNilOutportBlockData, err)

		outportBlock := createOutportBlock()
		outportBlock.BlockData = nil
		outportBlockBytes, _ := protoMarshaller.Marshal(outportBlock)
		err = fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock)
		require.Equal(t, errNilOutportBlockData, err)
	})

	t.Run("invalid payload, cannot unmarshall, should return error", func(t *testing.T) {
		t.Parallel()

		fi, _ := NewLogDataProcessor(&testscommon.IoWriterStub{}, block.NewEmptyBlockCreatorsContainer(), protoMarshaller)

		err := fi.ProcessPayload([]byte("invalid payload"), outportcore.TopicSaveBlock)
		require.NotNil(t, err)
	})

	t.Run("unknown block creator, should return error", func(t *testing.T) {
		t.Parallel()

		outportBlock := createOutportBlock()
		outportBlock.BlockData.HeaderType = "unknown"

		ioWriterCalledCt := 0
		ioWriter := &testscommon.IoWriterStub{
			WriteCalled: func(p []byte) (n int, err error) {
				ioWriterCalledCt++
				return 0, nil
			},
		}

		container := block.NewEmptyBlockCreatorsContainer()
		_ = container.Add(core.ShardHeaderV1, block.NewEmptyHeaderCreator())

		fi, _ := NewLogDataProcessor(ioWriter, container, protoMarshaller)

		outportBlockBytes, _ := protoMarshaller.Marshal(outportBlock)
		err := fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock)
		require.NotNil(t, err)
		require.Equal(t, 0, ioWriterCalledCt)
	})

	t.Run("cannot unmarshall to get header from bytes, should return error", func(t *testing.T) {
		t.Parallel()

		ioWriterCalledCt := 0
		ioWriter := &testscommon.IoWriterStub{
			WriteCalled: func(p []byte) (n int, err error) {
				ioWriterCalledCt++
				return 0, nil
			},
		}

		outportBlock := createOutportBlock()
		outportBlockBytes, _ := protoMarshaller.Marshal(outportBlock)

		unmarshalCalled := false
		errUnmarshal := errors.New("err unmarshal")
		marshaller := &testscommon.MarshallerStub{
			UnmarshalCalled: func(obj interface{}, buff []byte) error {
				defer func() {
					unmarshalCalled = true
				}()

				if !unmarshalCalled {
					require.Equal(t, outportBlockBytes, buff)

					err := protoMarshaller.Unmarshal(obj, buff)
					require.Nil(t, err)

					return nil

				}
				return errUnmarshal
			},
		}

		container := block.NewEmptyBlockCreatorsContainer()
		_ = container.Add(core.ShardHeaderV1, block.NewEmptyHeaderCreator())

		fi, _ := NewLogDataProcessor(ioWriter, container, marshaller)
		err := fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock)
		require.Equal(t, errUnmarshal, err)
	})

	t.Run("cannot write in console, should return error", func(t *testing.T) {
		t.Parallel()

		ioWriterCalledCt := 0
		err1 := errors.New("err1")
		err2 := errors.New("err2")
		ioWriter := &testscommon.IoWriterStub{
			WriteCalled: func(p []byte) (n int, err error) {
				defer func() {
					ioWriterCalledCt++
				}()

				switch ioWriterCalledCt {
				case 0:
					return 0, err1
				case 1:
					return 0, nil
				case 2:
					return 0, err2
				}

				return 0, nil
			},
		}

		container := block.NewEmptyBlockCreatorsContainer()
		_ = container.Add(core.ShardHeaderV1, block.NewEmptyHeaderCreator())

		fi, _ := NewLogDataProcessor(ioWriter, container, protoMarshaller)

		outportBlock := createOutportBlock()
		outportBlockBytes, _ := protoMarshaller.Marshal(outportBlock)

		err := fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock)
		require.True(t, strings.Contains(err.Error(), err1.Error()))

		err = fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock)
		require.True(t, strings.Contains(err.Error(), err2.Error()))

		err = fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock)
		require.Nil(t, err)

		require.Equal(t, 5, ioWriterCalledCt)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		header := &block.Header{
			Nonce:     1,
			PrevHash:  []byte("prev hash"),
			TimeStamp: 100,
		}
		headerBytes, err := protoMarshaller.Marshal(header)
		require.Nil(t, err)

		outportBlock := &outportcore.OutportBlock{
			BlockData: &outportcore.BlockData{
				HeaderHash:  []byte("hash"),
				HeaderBytes: headerBytes,
				HeaderType:  string(core.ShardHeaderV1),
			},
		}
		outportBlockBytes, err := protoMarshaller.Marshal(outportBlock)
		require.Nil(t, err)

		ioWriterCalledCt := 0
		ioWriter := &testscommon.IoWriterStub{
			WriteCalled: func(p []byte) (n int, err error) {
				defer func() {
					ioWriterCalledCt++
				}()

				switch ioWriterCalledCt {
				case 0:
					require.Equal(t, []byte("FIRE BLOCK_BEGIN 1\n"), p)
				case 1:
					require.Equal(t, []byte(fmt.Sprintf("FIRE BLOCK_END 1 %s 100 %x\n",
						hex.EncodeToString(header.PrevHash),
						outportBlockBytes)), p)
				default:
					require.Fail(t, "should not write again")
				}
				return 0, nil
			},
		}

		container := block.NewEmptyBlockCreatorsContainer()
		_ = container.Add(core.ShardHeaderV1, block.NewEmptyHeaderCreator())

		fi, _ := NewLogDataProcessor(ioWriter, container, protoMarshaller)

		err = fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock)
		require.Nil(t, err)
		require.Equal(t, 2, ioWriterCalledCt)
	})

}

func TestFirehoseIndexer_NoOperationFunctions(t *testing.T) {
	t.Parallel()

	fi, _ := NewLogDataProcessor(&testscommon.IoWriterStub{}, block.NewEmptyBlockCreatorsContainer(), protoMarshaller)

	err := fi.ProcessPayload([]byte("payload"), "invalid topic")
	require.True(t, strings.Contains(err.Error(), errInvalidOperationType.Error()))
	require.True(t, strings.Contains(err.Error(), "payload"))
	require.True(t, strings.Contains(err.Error(), "invalid topic"))

	require.Nil(t, fi.ProcessPayload([]byte("payload"), outportcore.TopicRevertIndexedBlock))
	require.Nil(t, fi.ProcessPayload([]byte("payload"), outportcore.TopicSaveRoundsInfo))
	require.Nil(t, fi.ProcessPayload([]byte("payload"), outportcore.TopicSaveValidatorsRating))
	require.Nil(t, fi.ProcessPayload([]byte("payload"), outportcore.TopicSaveValidatorsPubKeys))
	require.Nil(t, fi.ProcessPayload([]byte("payload"), outportcore.TopicSaveAccounts))
	require.Nil(t, fi.ProcessPayload([]byte("payload"), outportcore.TopicFinalizedBlock))
	require.Nil(t, fi.Close())
}
