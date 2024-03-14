package testdata

import (
	"errors"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/data/smartContractResult"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/marshal"
)

var errNilMarshaller = errors.New("nil marshaller")

type blockData struct {
	marshaller marshal.Marshalizer
}

// NewBlockData will create block data component for testing
func NewBlockData(marshaller marshal.Marshalizer) (*blockData, error) {
	if check.IfNil(marshaller) {
		return nil, errNilMarshaller
	}

	return &blockData{marshaller: marshaller}, nil
}

// OutportBlockV1 -
func (bd *blockData) OutportBlockV1() *outport.OutportBlock {
	header := &block.MetaBlock{
		TimeStamp: 1234,
	}
	headerBytes, _ := bd.marshaller.Marshal(header)

	return &outport.OutportBlock{
		ShardID: core.MetachainShardId,
		BlockData: &outport.BlockData{
			HeaderBytes: headerBytes,
			HeaderType:  "MetaBlock",
			HeaderHash:  []byte("headerHash1"),
			Body: &block.Body{
				MiniBlocks: []*block.MiniBlock{
					{
						TxHashes:        [][]byte{},
						ReceiverShardID: 1,
						SenderShardID:   1,
					},
				},
			},
		},
		HeaderGasConsumption: &outport.HeaderGasConsumption{},
		TransactionPool: &outport.TransactionPool{
			Transactions: map[string]*outport.TxInfo{
				"txHash1": {
					Transaction: &transaction.Transaction{
						Nonce:    1,
						GasPrice: 1,
						GasLimit: 1,
					},
					FeeInfo: &outport.FeeInfo{
						GasUsed: 1,
					},
					ExecutionOrder: 2,
				},
			},
			SmartContractResults: map[string]*outport.SCRInfo{
				"scrHash1": {
					SmartContractResult: &smartContractResult.SmartContractResult{
						Nonce:    2,
						GasLimit: 2,
						GasPrice: 2,
						CallType: 2,
					},
					FeeInfo: &outport.FeeInfo{
						GasUsed: 2,
					},
					ExecutionOrder: 0,
				},
			},
			Logs: []*outport.LogData{
				{
					Log: &transaction.Log{
						Address: []byte("logaddr1"),
						Events:  []*transaction.Event{},
					},
					TxHash: "logHash1",
				},
			},
		},
		NumberOfShards: 2,
	}
}

// RevertBlockV1 -
func (bd *blockData) RevertBlockV1() *outport.BlockData {
	header := &block.Header{
		ShardID:   1,
		TimeStamp: 1234,
	}
	headerBytes, _ := bd.marshaller.Marshal(header)

	return &outport.BlockData{
		ShardID:     1,
		HeaderBytes: headerBytes,
		HeaderType:  "Header",
		HeaderHash:  []byte("headerHash1"),
		Body: &block.Body{
			MiniBlocks: []*block.MiniBlock{
				{
					TxHashes:        [][]byte{},
					ReceiverShardID: 1,
					SenderShardID:   1,
				},
			},
		},
	}
}

// FinalizedBlockV1 -
func (bd *blockData) FinalizedBlockV1() *outport.FinalizedBlock {
	return &outport.FinalizedBlock{
		ShardID:    1,
		HeaderHash: []byte("headerHash1"),
	}
}
