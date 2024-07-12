package process

import (
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

type outportBlockConverter struct {
	gogoProtoMarshaller marshal.Marshalizer
	protoMarshaller     marshal.Marshalizer
	bigIntCaster        coreData.BigIntCaster
}

// NewOutportBlockConverter will return a component than can convert outport.OutportBlock to either data.ShardOutportBlock.
// or data.MetaOutportBlock
func NewOutportBlockConverter(
	gogoMarshaller marshal.Marshalizer,
	protoMarshaller marshal.Marshalizer,
) (*outportBlockConverter, error) {
	if check.IfNil(gogoMarshaller) {
		return nil, fmt.Errorf("%w: for gogo proto marshaller", ErrNilMarshaller)
	}
	if check.IfNil(protoMarshaller) {
		return nil, fmt.Errorf("%w: for proto marshaller", ErrNilMarshaller)
	}

	return &outportBlockConverter{
		gogoProtoMarshaller: gogoMarshaller,
		protoMarshaller:     protoMarshaller,
		bigIntCaster:        coreData.BigIntCaster{},
	}, nil
}

func (o *outportBlockConverter) HandleShardOutportBlockV2(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.ShardOutportBlockV2, error) {
	headerType := outportBlock.BlockData.HeaderType

	// check if the header type is supported by this function.
	if headerType != string(core.ShardHeaderV1) && headerType != string(core.ShardHeaderV2) {
		return nil, fmt.Errorf("cannot convert to shard outport block. header type: %s not supported", headerType)
	}

	shardOutportBlock := &hyperOutportBlocks.ShardOutportBlockV2{
		BlockData: &hyperOutportBlocks.BlockData{},
	}
	shardOutportBlock.ShardID = outportBlock.ShardID
	blockData, err := o.handleBlockData(outportBlock.BlockData, shardOutportBlock.BlockData)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate block data: %w", err)
	}
	err = o.handleTransactionPool(outportBlock.TransactionPool, shardOutportBlock.TransactionPool)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate transacion pool: %w", err)
	}
	//TODO: add these
	//handleHeaderGasConsumption()
	//handleAlteredAccounts()
	shardOutportBlock.BlockData = blockData

	shardOutportBlock.NotarizedHeadersHashes = outportBlock.NotarizedHeadersHashes
	shardOutportBlock.NumberOfShards = outportBlock.NumberOfShards
	shardOutportBlock.SignersIndexes = outportBlock.SignersIndexes
	shardOutportBlock.HighestFinalBlockNonce = outportBlock.HighestFinalBlockNonce
	shardOutportBlock.HighestFinalBlockHash = outportBlock.HighestFinalBlockHash

	return shardOutportBlock, nil
}

func (o *outportBlockConverter) handleBlockData(blockData *outport.BlockData, shardBlockData *hyperOutportBlocks.BlockData) (*hyperOutportBlocks.BlockData, error) {
	shardBlockData.ShardID = blockData.ShardID
	shardBlockData.HeaderType = blockData.HeaderType
	shardBlockData.HeaderHash = blockData.HeaderHash

	var err error
	switch blockData.HeaderType {
	case string(core.ShardHeaderV1):
		err = o.handleHeaderV1(blockData.HeaderBytes, shardBlockData)
		if err != nil {
			return nil, err
		}

	case string(core.ShardHeaderV2):
		err = o.handleHeaderV2(blockData.HeaderBytes, shardBlockData)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown header type [%s]", blockData.HeaderType)
	}

	return shardBlockData, nil
}

func (o *outportBlockConverter) handleHeaderV1(headerBytes []byte, shardBlockData *hyperOutportBlocks.BlockData) error {
	blockHeader := block.Header{}
	err := o.gogoProtoMarshaller.Unmarshal(&blockHeader, headerBytes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	header := shardBlockData.Header
	header.Nonce = blockHeader.Nonce
	header.PrevHash = blockHeader.PrevHash
	header.PrevRandSeed = blockHeader.PrevRandSeed
	header.RandSeed = blockHeader.RandSeed
	header.PubKeysBitmap = blockHeader.PubKeysBitmap
	header.ShardID = blockHeader.ShardID
	header.TimeStamp = blockHeader.TimeStamp
	header.Round = blockHeader.Round
	header.Epoch = blockHeader.Epoch
	header.BlockBodyType = hyperOutportBlocks.Type(blockHeader.BlockBodyType)
	header.Signature = blockHeader.Signature
	header.LeaderSignature = blockHeader.LeaderSignature

	miniBlockHeaders := make([]*hyperOutportBlocks.MiniBlockHeader, 0)
	for _, miniBlockHeader := range blockHeader.MiniBlockHeaders {
		mb := &hyperOutportBlocks.MiniBlockHeader{
			Hash:            miniBlockHeader.Hash,
			SenderShardID:   miniBlockHeader.SenderShardID,
			ReceiverShardID: miniBlockHeader.ReceiverShardID,
			TxCount:         miniBlockHeader.TxCount,
			Type:            hyperOutportBlocks.Type(miniBlockHeader.Type),
			Reserved:        miniBlockHeader.Reserved,
		}
		miniBlockHeaders = append(miniBlockHeaders, mb)
	}
	header.MiniBlockHeaders = miniBlockHeaders

	peerChanges := make([]*hyperOutportBlocks.PeerChange, 0)
	for _, peerChange := range blockHeader.PeerChanges {
		pc := &hyperOutportBlocks.PeerChange{
			PubKey:      peerChange.PubKey,
			ShardIdDest: peerChange.ShardIdDest,
		}

		peerChanges = append(peerChanges, pc)
	}
	header.PeerChanges = peerChanges

	header.RootHash = blockHeader.RootHash
	header.MetaBlockHashes = blockHeader.MetaBlockHashes
	header.TxCount = blockHeader.TxCount
	header.EpochStartMetaHash = blockHeader.EpochStartMetaHash
	header.ReceiptsHash = blockHeader.ReceiptsHash
	header.ChainID = blockHeader.ChainID
	header.SoftwareVersion = blockHeader.SoftwareVersion
	header.Reserved = blockHeader.Reserved

	accumulatedFees, err := o.castBigInt(blockHeader.AccumulatedFees)
	if err != nil {
		return fmt.Errorf("failed to cast accumulated fees: %w", err)
	}
	developerFees, err := o.castBigInt(blockHeader.DeveloperFees)
	if err != nil {
		return fmt.Errorf("failed to cast developer fees: %w", err)
	}

	header.AccumulatedFees = accumulatedFees
	header.DeveloperFees = developerFees

	return nil
}

func (o *outportBlockConverter) handleHeaderV2(headerBytes []byte, shardBlockData *hyperOutportBlocks.BlockData) error {
	blockHeader := block.HeaderV2{}
	err := o.gogoProtoMarshaller.Unmarshal(&blockHeader, headerBytes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	header := shardBlockData.Header
	header.Nonce = blockHeader.Header.Nonce
	header.PrevHash = blockHeader.Header.PrevHash
	header.PrevRandSeed = blockHeader.Header.PrevRandSeed
	header.RandSeed = blockHeader.Header.RandSeed
	header.PubKeysBitmap = blockHeader.Header.PubKeysBitmap
	header.ShardID = blockHeader.Header.ShardID
	header.TimeStamp = blockHeader.Header.TimeStamp
	header.Round = blockHeader.Header.Round
	header.Epoch = blockHeader.Header.Epoch
	header.BlockBodyType = hyperOutportBlocks.Type(blockHeader.Header.BlockBodyType)
	header.Signature = blockHeader.Header.Signature
	header.LeaderSignature = blockHeader.Header.LeaderSignature

	miniBlockHeaders := make([]*hyperOutportBlocks.MiniBlockHeader, 0)
	for _, miniBlockHeader := range blockHeader.Header.MiniBlockHeaders {
		mb := &hyperOutportBlocks.MiniBlockHeader{
			Hash:            miniBlockHeader.Hash,
			SenderShardID:   miniBlockHeader.SenderShardID,
			ReceiverShardID: miniBlockHeader.ReceiverShardID,
			TxCount:         miniBlockHeader.TxCount,
			Type:            hyperOutportBlocks.Type(miniBlockHeader.Type),
			Reserved:        miniBlockHeader.Reserved,
		}
		miniBlockHeaders = append(miniBlockHeaders, mb)
	}
	header.MiniBlockHeaders = miniBlockHeaders

	peerChanges := make([]*hyperOutportBlocks.PeerChange, 0)
	for _, peerChange := range blockHeader.Header.PeerChanges {
		pc := &hyperOutportBlocks.PeerChange{
			PubKey:      peerChange.PubKey,
			ShardIdDest: peerChange.ShardIdDest,
		}

		peerChanges = append(peerChanges, pc)
	}
	header.PeerChanges = peerChanges

	header.RootHash = blockHeader.Header.RootHash
	header.MetaBlockHashes = blockHeader.Header.MetaBlockHashes
	header.TxCount = blockHeader.Header.TxCount
	header.EpochStartMetaHash = blockHeader.Header.EpochStartMetaHash
	header.ReceiptsHash = blockHeader.Header.ReceiptsHash
	header.ChainID = blockHeader.Header.ChainID
	header.SoftwareVersion = blockHeader.Header.SoftwareVersion
	header.Reserved = blockHeader.Header.Reserved

	accumulatedFees, err := o.castBigInt(blockHeader.Header.AccumulatedFees)
	if err != nil {
		return fmt.Errorf("failed to cast accumulated fees: %w", err)
	}
	developerFees, err := o.castBigInt(blockHeader.Header.DeveloperFees)
	if err != nil {
		return fmt.Errorf("failed to cast developer fees: %w", err)
	}

	header.AccumulatedFees = accumulatedFees
	header.DeveloperFees = developerFees

	shardBlockData.ScheduledRootHash = blockHeader.ScheduledRootHash
	shardBlockData.ScheduledAccumulatedFees, err = o.castBigInt(blockHeader.ScheduledAccumulatedFees)
	if err != nil {
		return fmt.Errorf("failed to cast scheduled fees: %w", err)
	}
	shardBlockData.ScheduledDeveloperFees, err = o.castBigInt(blockHeader.ScheduledDeveloperFees)
	if err != nil {
		return fmt.Errorf("failed to cast scheduled developer fees: %w", err)
	}
	shardBlockData.ScheduledGasProvided = blockHeader.ScheduledGasProvided
	shardBlockData.ScheduledGasPenalized = blockHeader.ScheduledGasPenalized
	shardBlockData.ScheduledGasRefunded = blockHeader.ScheduledGasRefunded

	return nil
}

func (o *outportBlockConverter) handleTransactionPool(outportTxPool *outport.TransactionPool, shardTxPool *hyperOutportBlocks.TransactionPoolV2) error {
	transactions := make(map[string]*hyperOutportBlocks.TxInfoV2, 0)
	var err error
	for txHash, tx := range outportTxPool.Transactions {
		txInfo := &hyperOutportBlocks.TxInfoV2{}

		// TxInfo - Transaction
		txInfo.Transaction.Nonce = tx.Transaction.Nonce
		txInfo.Transaction.Value, err = o.castBigInt(tx.Transaction.Value)
		if err != nil {
			return fmt.Errorf("failed to cast transaction [%s] value: %w", txHash, err)
		}
		txInfo.Transaction.RcvAddr = tx.Transaction.RcvAddr
		txInfo.Transaction.RcvUserName = tx.Transaction.RcvUserName
		txInfo.Transaction.SndAddr = tx.Transaction.SndAddr
		txInfo.Transaction.SndUserName = tx.Transaction.SndUserName
		txInfo.Transaction.GasPrice = tx.Transaction.GasPrice
		txInfo.Transaction.GasLimit = tx.Transaction.GasLimit
		txInfo.Transaction.Data = tx.Transaction.Data
		txInfo.Transaction.ChainID = tx.Transaction.ChainID
		txInfo.Transaction.Version = tx.Transaction.Version
		txInfo.Transaction.Signature = tx.Transaction.Signature
		txInfo.Transaction.Options = tx.Transaction.Options
		txInfo.Transaction.GuardianAddr = tx.Transaction.GuardianAddr
		txInfo.Transaction.GuardianSignature = tx.Transaction.GuardianSignature

		// TxInfo - FeeInfo
		txInfo.FeeInfo.GasUsed = tx.FeeInfo.GasUsed
		txInfo.FeeInfo.Fee, err = o.castBigInt(tx.FeeInfo.Fee)
		if err != nil {
			return fmt.Errorf("failed to cast transaction [%s] fee: %w", txHash, err)
		}

		txInfo.FeeInfo.InitialPaidFee, err = o.castBigInt(tx.FeeInfo.InitialPaidFee)
		if err != nil {
			return fmt.Errorf("failed to cast transaction [%s] initial paid fee: %w", txHash, err)
		}

		txInfo.ExecutionOrder = tx.ExecutionOrder

		if sc, ok := outportTxPool.SmartContractResults[txHash]; ok {
			txInfo.SmartContractResults.SmartContractResult.Nonce = sc.SmartContractResult.Nonce
			txInfo.SmartContractResults.SmartContractResult.Value, err = o.castBigInt(sc.SmartContractResult.Value)
			if err != nil {
				return fmt.Errorf("failed to cast transaction [%s] smart contract result value: %s", txHash, sc.SmartContractResult.Value)
			}
			txInfo.SmartContractResults.SmartContractResult.RcvAddr = sc.SmartContractResult.RcvAddr
			txInfo.SmartContractResults.SmartContractResult.SndAddr = sc.SmartContractResult.SndAddr
			txInfo.SmartContractResults.SmartContractResult.RelayerAddr = sc.SmartContractResult.RelayerAddr
			txInfo.SmartContractResults.SmartContractResult.RelayedValue, err = o.castBigInt(sc.SmartContractResult.RelayedValue)
			if err != nil {
				return fmt.Errorf("failed to cast transaction [%s] smart contract result relayed value: %s", txHash, sc.SmartContractResult.Value)
			}
			txInfo.SmartContractResults.SmartContractResult.Code = sc.SmartContractResult.Code
			txInfo.SmartContractResults.SmartContractResult.Data = sc.SmartContractResult.Data
			txInfo.SmartContractResults.SmartContractResult.PrevTxHash = sc.SmartContractResult.PrevTxHash
			txInfo.SmartContractResults.SmartContractResult.OriginalTxHash = sc.SmartContractResult.OriginalTxHash
			txInfo.SmartContractResults.SmartContractResult.GasLimit = sc.SmartContractResult.GasLimit
			txInfo.SmartContractResults.SmartContractResult.GasPrice = sc.SmartContractResult.GasPrice
			txInfo.SmartContractResults.SmartContractResult.CallType = int64(sc.SmartContractResult.CallType)
			txInfo.SmartContractResults.SmartContractResult.CodeMetadata = sc.SmartContractResult.CodeMetadata
			txInfo.SmartContractResults.SmartContractResult.ReturnMessage = sc.SmartContractResult.ReturnMessage
			txInfo.SmartContractResults.SmartContractResult.OriginalSender = sc.SmartContractResult.OriginalSender
		}

		// Logs
		for _, logData := range outportTxPool.Logs {
			events := make([]*hyperOutportBlocks.Event, len(logData.Log.Events))
			for i, event := range logData.Log.Events {
				e := &hyperOutportBlocks.Event{}

				e.Address = event.Address
				e.Identifier = event.Identifier
				e.Topics = event.Topics
				e.Data = event.Data
				e.AdditionalData = event.AdditionalData

				events[i] = e
			}

			ll := &hyperOutportBlocks.Log{
				Address: logData.Log.Address,
				Events:  events,
			}

			if transactions[logData.TxHash].Logs == nil {
				transactions[logData.TxHash].Logs = make([]*hyperOutportBlocks.Log, 0)
			}

			transactions[logData.TxHash].Logs = append(transactions[logData.TxHash].Logs, ll)
		}

		transactions[txHash] = txInfo

		for _, l := range outportTxPool.Logs {
			if _, ok := transactions[l.TxHash]; ok {
				events := make([]*hyperOutportBlocks.Event, len(l.Log.Events))
				for i, event := range l.Log.Events {
					e := &hyperOutportBlocks.Event{}

					e.Address = event.Address
					e.Identifier = event.Identifier
					e.Topics = event.Topics
					e.Data = event.Data
					e.AdditionalData = event.AdditionalData

					events[i] = e
				}

				ll := &hyperOutportBlocks.Log{
					Address: l.Log.Address,
					Events:  events,
				}

				transactions[l.TxHash].Logs = append(transactions[l.TxHash].Logs, ll)
			}
		}
	}

	shardTxPool.Transactions = transactions
	return nil
}

// HandleShardOutportBlock will convert an outport.OutportBlock to data.ShardOutportBlock.
func (o *outportBlockConverter) HandleShardOutportBlock(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.ShardOutportBlock, error) {
	headerType := outportBlock.BlockData.HeaderType

	// check if the header type is supported by this function.
	if headerType != string(core.ShardHeaderV1) && headerType != string(core.ShardHeaderV2) {
		return nil, fmt.Errorf("cannot convert to shard outport block. header type: %s not supported", headerType)
	}

	// marshal with gogo, since the outportBlock is gogo protobuf (coming from the node).
	bytes, err := o.gogoProtoMarshaller.Marshal(outportBlock)
	if err != nil {
		return nil, fmt.Errorf("marshal shard outport block error: %w", err)
	}

	shardOutportBlock := &hyperOutportBlocks.ShardOutportBlock{}
	// unmarshall into google protobuf. This is the proto that will be later consumed.
	err = o.protoMarshaller.Unmarshal(shardOutportBlock, bytes)
	if err != nil {
		return nil, fmt.Errorf("unmarshal shard outport block error: %w", err)
	}

	// ShardHeaderV1 marshals 1 to 1 into *data.ShardOutportBlock.
	if headerType == string(core.ShardHeaderV1) {
		return shardOutportBlock, nil
	}

	// ShardHeaderV2 does not marshal 1 to 1. A few fields need to be injected.
	header := block.HeaderV2{}
	err = o.gogoProtoMarshaller.Unmarshal(&header, outportBlock.BlockData.HeaderBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	miniBlockHeaders := make([]*hyperOutportBlocks.MiniBlockHeader, 0)
	for _, miniBlockHeader := range header.Header.MiniBlockHeaders {
		mb := &hyperOutportBlocks.MiniBlockHeader{
			Hash:            miniBlockHeader.Hash,
			SenderShardID:   miniBlockHeader.SenderShardID,
			ReceiverShardID: miniBlockHeader.ReceiverShardID,
			TxCount:         miniBlockHeader.TxCount,
			Type:            hyperOutportBlocks.Type(miniBlockHeader.Type),
			Reserved:        miniBlockHeader.Reserved,
		}
		miniBlockHeaders = append(miniBlockHeaders, mb)
	}

	peerChanges := make([]*hyperOutportBlocks.PeerChange, 0)
	for _, peerChange := range header.Header.PeerChanges {
		pc := &hyperOutportBlocks.PeerChange{
			PubKey:      peerChange.PubKey,
			ShardIdDest: peerChange.ShardIdDest,
		}

		peerChanges = append(peerChanges, pc)
	}

	accumulatedFees, err := o.castBigInt(header.Header.AccumulatedFees)
	if err != nil {
		return nil, fmt.Errorf("failed to cast accumulated fees: %w", err)
	}
	developerFees, err := o.castBigInt(header.Header.DeveloperFees)
	if err != nil {
		return nil, fmt.Errorf("failed to cast developer fees: %w", err)
	}

	shardOutportBlock.BlockData.Header = &hyperOutportBlocks.Header{
		Nonce:              header.Header.Nonce,
		PrevHash:           header.Header.PrevHash,
		PrevRandSeed:       header.Header.PrevRandSeed,
		RandSeed:           header.Header.RandSeed,
		PubKeysBitmap:      header.Header.PubKeysBitmap,
		ShardID:            header.Header.ShardID,
		TimeStamp:          header.Header.TimeStamp,
		Round:              header.Header.Round,
		Epoch:              header.Header.Epoch,
		BlockBodyType:      hyperOutportBlocks.Type(header.Header.BlockBodyType),
		Signature:          header.Header.Signature,
		LeaderSignature:    header.Header.LeaderSignature,
		MiniBlockHeaders:   miniBlockHeaders,
		PeerChanges:        peerChanges,
		RootHash:           header.Header.RootHash,
		MetaBlockHashes:    header.Header.MetaBlockHashes,
		TxCount:            header.Header.TxCount,
		EpochStartMetaHash: header.Header.EpochStartMetaHash,
		ReceiptsHash:       header.Header.ReceiptsHash,
		ChainID:            header.Header.ChainID,
		SoftwareVersion:    header.Header.SoftwareVersion,
		AccumulatedFees:    accumulatedFees,
		DeveloperFees:      developerFees,
		Reserved:           header.Header.Reserved,
	}

	scheduledAccumulatedFees, err := o.castBigInt(header.ScheduledAccumulatedFees)
	if err != nil {
		return nil, fmt.Errorf("failed to cast scheduled accumulated fees: %w", err)
	}
	scheduledDeveloperFees, err := o.castBigInt(header.ScheduledDeveloperFees)
	if err != nil {
		return nil, fmt.Errorf("failed to cast scheduled developer fees: %w", err)
	}

	shardOutportBlock.BlockData.ScheduledRootHash = header.ScheduledRootHash
	shardOutportBlock.BlockData.ScheduledAccumulatedFees = scheduledAccumulatedFees
	shardOutportBlock.BlockData.ScheduledDeveloperFees = scheduledDeveloperFees
	shardOutportBlock.BlockData.ScheduledGasProvided = header.ScheduledGasProvided
	shardOutportBlock.BlockData.ScheduledGasPenalized = header.ScheduledGasPenalized
	shardOutportBlock.BlockData.ScheduledGasRefunded = header.ScheduledGasRefunded

	return shardOutportBlock, nil
}

// HandleMetaOutportBlock will convert an outport.OutportBlock to data.MetaOutportBlock.
func (o *outportBlockConverter) HandleMetaOutportBlock(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.MetaOutportBlock, error) {
	headerType := outportBlock.BlockData.HeaderType

	// check if the header type is supported by this function.
	if headerType != string(core.MetaHeader) {
		return nil, fmt.Errorf("cannot convert to meta outport block. header type: %s not supported", outportBlock.BlockData.HeaderType)
	}

	// marshal with gogo, since the outportBlock is gogo protobuf (coming from the node).
	bytes, err := o.gogoProtoMarshaller.Marshal(outportBlock)
	if err != nil {
		return nil, fmt.Errorf("marshal metaBlockCaster error: %w", err)
	}

	// unmarshall into google protobuf. This is the proto that will be used in firehose.
	metaOutportBlock := &hyperOutportBlocks.MetaOutportBlock{}
	err = o.protoMarshaller.Unmarshal(metaOutportBlock, bytes)
	if err != nil {
		return nil, fmt.Errorf("unmarshal metaBlockCaster error: %w", err)
	}

	return metaOutportBlock, nil
}

func (o *outportBlockConverter) castBigInt(i *big.Int) ([]byte, error) {
	buf := make([]byte, o.bigIntCaster.Size(i))
	_, err := o.bigIntCaster.MarshalTo(i, buf)

	return buf, err
}

// IsInterfaceNil returns nil if there is no value under the interface
func (o *outportBlockConverter) IsInterfaceNil() bool {
	return o == nil
}
