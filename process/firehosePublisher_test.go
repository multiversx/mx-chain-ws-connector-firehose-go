package process_test

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/stretchr/testify/require"

	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
)

var (
	protoMarshaller     = &process.ProtoMarshaller{}
	gogoProtoMarshaller = &marshal.GogoProtoMarshalizer{}
)

func createHyperOutportBlock() *data.HyperOutportBlock {
	hyperOutportBlock := &data.HyperOutportBlock{
		MetaOutportBlock: &data.MetaOutportBlock{
			ShardID: 1,
			BlockData: &data.MetaBlockData{
				Header: &data.MetaHeader{
					Nonce:     1,
					PrevHash:  []byte("prev hash"),
					TimeStamp: 100,
				},
				HeaderType: string(core.ShardHeaderV1),
				HeaderHash: []byte("hash"),
			},
			NotarizedHeadersHashes: []string{},
			NumberOfShards:         0,
			SignersIndexes:         []uint64{},
			HighestFinalBlockNonce: 0,
			HighestFinalBlockHash:  []byte{},
		},
		NotarizedHeadersOutportData: []*data.NotarizedHeaderOutportData{},
	}

	return hyperOutportBlock
}

func createContainer() process.BlockContainerHandler {
	container := block.NewEmptyBlockCreatorsContainer()

	_ = container.Add(core.MetaHeader, block.NewEmptyMetaBlockCreator())
	_ = container.Add(core.ShardHeaderV1, block.NewEmptyHeaderCreator())
	_ = container.Add(core.ShardHeaderV2, block.NewEmptyHeaderV2Creator())

	return container
}

func TestNewFirehosePublisher(t *testing.T) {
	t.Parallel()

	t.Run("nil io writer, should return error", func(t *testing.T) {
		t.Parallel()

		fp, err := process.NewFirehosePublisher(nil, createContainer(), protoMarshaller)
		require.Nil(t, fp)
		require.Equal(t, process.ErrNilWriter, err)
	})

	t.Run("nil block creator, should return error", func(t *testing.T) {
		t.Parallel()

		fp, err := process.NewFirehosePublisher(&testscommon.IoWriterStub{}, nil, protoMarshaller)
		require.Nil(t, fp)
		require.Equal(t, process.ErrNilBlockCreator, err)
	})

	t.Run("nil marshaller, should return error", func(t *testing.T) {
		t.Parallel()

		fp, err := process.NewFirehosePublisher(&testscommon.IoWriterStub{}, createContainer(), nil)
		require.Nil(t, fp)
		require.Equal(t, process.ErrNilMarshaller, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		fp, err := process.NewFirehosePublisher(&testscommon.IoWriterStub{}, createContainer(), protoMarshaller)
		require.Nil(t, err)
		require.False(t, check.IfNil(fp))
	})
}

func TestFirehosePublisher_PublishHyperBlock(t *testing.T) {
	t.Parallel()

	t.Run("unknown block creator for header type, should return error", func(t *testing.T) {
		t.Parallel()

		outportBlock := createHyperOutportBlock()
		outportBlock.MetaOutportBlock.BlockData.HeaderType = "unknown"

		ioWriterCalledCt := 0
		ioWriter := &testscommon.IoWriterStub{
			WriteCalled: func(p []byte) (n int, err error) {
				ioWriterCalledCt++
				return 0, nil
			},
		}

		fi, _ := process.NewFirehosePublisher(ioWriter, createContainer(), gogoProtoMarshaller)

		err := fi.PublishHyperBlock(outportBlock)
		require.NotNil(t, err)
		require.Equal(t, 1, ioWriterCalledCt) // 1 write comes from constructor
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

		expectedErr := errors.New("expected err")
		marshaller := &testscommon.MarshallerStub{
			MarshalCalled: func(obj interface{}) ([]byte, error) {
				return nil, expectedErr
			},
		}

		fp, _ := process.NewFirehosePublisher(ioWriter, createContainer(), marshaller)

		outportBlock := createHyperOutportBlock()
		err := fp.PublishHyperBlock(outportBlock)
		require.Equal(t, expectedErr, err)
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
					return 0, nil
				case 1:
					return 0, err1
				case 2:
					return 0, nil
				case 3:
					return 0, err2
				}

				return 0, nil
			},
		}

		fp, _ := process.NewFirehosePublisher(ioWriter, createContainer(), protoMarshaller)

		outportBlock := createHyperOutportBlock()

		err := fp.PublishHyperBlock(outportBlock)
		require.True(t, strings.Contains(err.Error(), err1.Error()))

		err = fp.PublishHyperBlock(outportBlock)
		require.Nil(t, err)

		err = fp.PublishHyperBlock(outportBlock)
		require.True(t, errors.Is(err, err2))

		require.Equal(t, 4, ioWriterCalledCt)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		outportBlock := &data.HyperOutportBlock{
			MetaOutportBlock: &data.MetaOutportBlock{
				ShardID: 1,
				BlockData: &data.MetaBlockData{
					Header: &data.MetaHeader{
						Nonce:     1,
						PrevHash:  []byte("prev hash"),
						TimeStamp: 100,
					},
					HeaderType: string(core.ShardHeaderV1),
					HeaderHash: []byte("hash"),
				},
				NotarizedHeadersHashes: []string{},
				NumberOfShards:         0,
				SignersIndexes:         []uint64{},
				HighestFinalBlockNonce: 0,
				HighestFinalBlockHash:  []byte{},
			},
			NotarizedHeadersOutportData: []*data.NotarizedHeaderOutportData{},
		}
		outportBlockBytes, err := protoMarshaller.Marshal(outportBlock)
		require.Nil(t, err)

		ioWriterCalledCt := 0
		header := outportBlock.MetaOutportBlock.BlockData.Header
		ioWriter := &testscommon.IoWriterStub{
			WriteCalled: func(p []byte) (n int, err error) {
				defer func() {
					ioWriterCalledCt++
				}()

				switch ioWriterCalledCt {
				case 0:
				case 1:
					num := header.GetNonce()
					parentNum := num - 1
					libNum := parentNum
					encodedMvxBlock := base64.StdEncoding.EncodeToString(outportBlockBytes)

					require.Equal(t, []byte(fmt.Sprintf("%s %s %d %s %d %s %d %d %s\n",
						process.FirehosePrefix,
						process.BlockPrefix,
						num,
						hex.EncodeToString(outportBlock.MetaOutportBlock.BlockData.HeaderHash),
						parentNum,
						hex.EncodeToString(header.PrevHash),
						libNum,
						header.TimeStamp,
						encodedMvxBlock)), p)
				default:
					require.Fail(t, "should not write again")
				}
				return 0, nil
			},
		}

		fp, _ := process.NewFirehosePublisher(ioWriter, createContainer(), protoMarshaller)

		err = fp.PublishHyperBlock(outportBlock)
		require.Nil(t, err)
		require.Equal(t, 2, ioWriterCalledCt)
	})
}

func TestFirehosePublisher_Close(t *testing.T) {
	t.Parallel()

	closeError := errors.New("error closing")
	writer := &testscommon.IoWriterStub{
		CloseCalled: func() error {
			return closeError
		},
	}

	fi, _ := process.NewFirehosePublisher(writer, createContainer(), protoMarshaller)
	err := fi.Close()
	require.Equal(t, closeError, err)
}
