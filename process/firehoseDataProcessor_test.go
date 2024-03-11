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
	"github.com/stretchr/testify/require"

	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
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

func createContainer() BlockContainerHandler {
	container := block.NewEmptyBlockCreatorsContainer()
	_ = container.Add(core.ShardHeaderV1, block.NewEmptyHeaderCreator())

	return container
}

func TestNewFirehoseDataProcessor(t *testing.T) {
	t.Parallel()

	t.Run("nil io writer, should return error", func(t *testing.T) {
		t.Parallel()

		fi, err := NewFirehoseDataProcessor(nil, createContainer(), protoMarshaller)
		require.Nil(t, fi)
		require.Equal(t, errNilWriter, err)
	})

	t.Run("nil block creator, should return error", func(t *testing.T) {
		t.Parallel()

		fi, err := NewFirehoseDataProcessor(&testscommon.IoWriterStub{}, nil, protoMarshaller)
		require.Nil(t, fi)
		require.Equal(t, errNilBlockCreator, err)
	})

	t.Run("nil marshaller, should return error", func(t *testing.T) {
		t.Parallel()

		fi, err := NewFirehoseDataProcessor(&testscommon.IoWriterStub{}, createContainer(), nil)
		require.Nil(t, fi)
		require.Equal(t, errNilMarshaller, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		fi, err := NewFirehoseDataProcessor(&testscommon.IoWriterStub{}, createContainer(), protoMarshaller)
		require.Nil(t, err)
		require.False(t, check.IfNil(fi))
	})
}

func TestFirehoseIndexer_SaveBlock(t *testing.T) {
	t.Parallel()

	t.Run("nil outport block, should return error", func(t *testing.T) {
		t.Parallel()

		fi, _ := NewFirehoseDataProcessor(&testscommon.IoWriterStub{}, createContainer(), protoMarshaller)

		err := fi.ProcessPayload(nil, outportcore.TopicSaveBlock, 1)
		require.Equal(t, errNilOutportBlockData, err)

		outportBlock := createOutportBlock()
		outportBlock.BlockData = nil
		outportBlockBytes, _ := protoMarshaller.Marshal(outportBlock)

		err = fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
		require.Equal(t, errNilOutportBlockData, err)
	})

	t.Run("invalid payload, cannot unmarshall, should return error", func(t *testing.T) {
		t.Parallel()

		fi, _ := NewFirehoseDataProcessor(&testscommon.IoWriterStub{}, createContainer(), protoMarshaller)

		err := fi.ProcessPayload([]byte("invalid payload"), outportcore.TopicSaveBlock, 1)
		require.NotNil(t, err)
	})

	t.Run("unknown block creator for header type, should return error", func(t *testing.T) {
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

		fi, _ := NewFirehoseDataProcessor(ioWriter, createContainer(), protoMarshaller)

		outportBlockBytes, _ := protoMarshaller.Marshal(outportBlock)
		err := fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
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

		unmarshalCalledBefore := false
		errUnmarshal := errors.New("err unmarshal")
		marshaller := &testscommon.MarshallerStub{
			UnmarshalCalled: func(obj interface{}, buff []byte) error {
				defer func() {
					unmarshalCalledBefore = true
				}()

				if unmarshalCalledBefore {
					return errUnmarshal
				}

				require.Equal(t, outportBlockBytes, buff)
				err := protoMarshaller.Unmarshal(obj, buff)
				require.Nil(t, err)

				return nil
			},
		}

		fi, _ := NewFirehoseDataProcessor(ioWriter, createContainer(), marshaller)
		err := fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
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

		fi, _ := NewFirehoseDataProcessor(ioWriter, createContainer(), protoMarshaller)

		outportBlock := createOutportBlock()
		outportBlockBytes, _ := protoMarshaller.Marshal(outportBlock)

		err := fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
		require.True(t, strings.Contains(err.Error(), err1.Error()))

		err = fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
		require.True(t, strings.Contains(err.Error(), err2.Error()))

		err = fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
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

		fi, _ := NewFirehoseDataProcessor(ioWriter, createContainer(), protoMarshaller)

		err = fi.ProcessPayload(outportBlockBytes, outportcore.TopicSaveBlock, 1)
		require.Nil(t, err)
		require.Equal(t, 2, ioWriterCalledCt)
	})

}

func TestFirehoseIndexer_NoOperationFunctions(t *testing.T) {
	t.Parallel()

	fi, _ := NewFirehoseDataProcessor(&testscommon.IoWriterStub{}, createContainer(), protoMarshaller)

	require.Nil(t, fi.ProcessPayload([]byte("payload"), "random topic", 1))
	require.Nil(t, fi.ProcessPayload([]byte("payload"), outportcore.TopicSaveRoundsInfo, 1))
	require.Nil(t, fi.ProcessPayload([]byte("payload"), outportcore.TopicSaveValidatorsRating, 1))
	require.Nil(t, fi.ProcessPayload([]byte("payload"), outportcore.TopicSaveValidatorsPubKeys, 1))
	require.Nil(t, fi.ProcessPayload([]byte("payload"), outportcore.TopicSaveAccounts, 1))
	require.Nil(t, fi.ProcessPayload([]byte("payload"), outportcore.TopicFinalizedBlock, 1))
}

func TestFirehoseIndexer_Close(t *testing.T) {
	t.Parallel()

	closeError := errors.New("error closing")
	writer := &testscommon.IoWriterStub{
		CloseCalled: func() error {
			return closeError
		},
	}

	fi, _ := NewFirehoseDataProcessor(writer, createContainer(), protoMarshaller)
	err := fi.Close()
	require.Equal(t, closeError, err)
}
