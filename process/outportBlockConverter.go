package process

import (
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/alteredAccount"
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

	shardOutportBlock := &hyperOutportBlocks.ShardOutportBlockV2{}
	shardOutportBlock.ShardID = outportBlock.ShardID
	err := o.handleBlockData(outportBlock.BlockData, shardOutportBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate block data: %w", err)
	}
	err = o.handleTransactionPool(outportBlock.TransactionPool, shardOutportBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate transacion pool: %w", err)
	}
	handleHeaderGasConsumption(outportBlock.HeaderGasConsumption, shardOutportBlock)
	handleAlteredAccounts(outportBlock.AlteredAccounts, shardOutportBlock)

	shardOutportBlock.NotarizedHeadersHashes = outportBlock.NotarizedHeadersHashes
	shardOutportBlock.NumberOfShards = outportBlock.NumberOfShards
	shardOutportBlock.SignersIndexes = outportBlock.SignersIndexes
	shardOutportBlock.HighestFinalBlockNonce = outportBlock.HighestFinalBlockNonce
	shardOutportBlock.HighestFinalBlockHash = outportBlock.HighestFinalBlockHash

	return shardOutportBlock, nil
}

func (o *outportBlockConverter) handleBlockData(blockData *outport.BlockData, shardOutportBlock *hyperOutportBlocks.ShardOutportBlockV2) error {
	if blockData == nil {
		return nil
	}

	shardOutportBlock.BlockData = &hyperOutportBlocks.BlockData{
		ShardID:              blockData.ShardID,
		Header:               nil,
		HeaderType:           blockData.HeaderType,
		HeaderHash:           blockData.HeaderHash,
		Body:                 nil,
		IntraShardMiniBlocks: nil,
	}

	var err error
	switch blockData.HeaderType {
	case string(core.ShardHeaderV1):
		err = o.handleHeaderV1(blockData.HeaderBytes, shardOutportBlock.BlockData)
		if err != nil {
			return fmt.Errorf("failed to handle header v1: %w", err)
		}

	case string(core.ShardHeaderV2):
		err = o.handleHeaderV2(blockData.HeaderBytes, shardOutportBlock.BlockData)
		if err != nil {
			return fmt.Errorf("failed to handle header v2: %w", err)
		}

	default:
		return fmt.Errorf("unknown header type [%s]", blockData.HeaderType)
	}

	return nil
}

func (o *outportBlockConverter) handleHeaderV1(headerBytes []byte, shardBlockData *hyperOutportBlocks.BlockData) error {
	blockHeader := block.Header{}
	err := o.gogoProtoMarshaller.Unmarshal(&blockHeader, headerBytes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	shardBlockData.Header = &hyperOutportBlocks.Header{}
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

func (o *outportBlockConverter) handleTransactionPool(outportTxPool *outport.TransactionPool, shardOutportBlockV2 *hyperOutportBlocks.ShardOutportBlockV2) error {
	if outportTxPool == nil {
		return nil
	}

	shardOutportBlockV2.TransactionPool = &hyperOutportBlocks.TransactionPoolV2{}

	if outportTxPool.Transactions != nil {
		shardOutportBlockV2.TransactionPool.Transactions = make(map[string]*hyperOutportBlocks.TxInfoV2, len(outportTxPool.Transactions))
	}

	if outportTxPool.Rewards != nil {
		shardOutportBlockV2.TransactionPool.Rewards = make(map[string]*hyperOutportBlocks.RewardInfo, len(outportTxPool.Rewards))
	}

	if outportTxPool.Receipts != nil {
		shardOutportBlockV2.TransactionPool.Receipts = make(map[string]*hyperOutportBlocks.Receipt, len(outportTxPool.Receipts))
	}

	if outportTxPool.InvalidTxs != nil {
		shardOutportBlockV2.TransactionPool.InvalidTxs = make(map[string]*hyperOutportBlocks.TxInfo, len(outportTxPool.InvalidTxs))
	}

	shardOutportBlockV2.TransactionPool.ScheduledExecutedSCRSHashesPrevBlock = outportTxPool.ScheduledExecutedSCRSHashesPrevBlock
	shardOutportBlockV2.TransactionPool.ScheduledExecutedInvalidTxsHashesPrevBlock = outportTxPool.ScheduledExecutedInvalidTxsHashesPrevBlock

	for txHash, tx := range outportTxPool.Transactions {
		txInfo := &hyperOutportBlocks.TxInfoV2{}

		// TxInfo - Transaction
		if tx.Transaction != nil {
			value, err := o.castBigInt(tx.Transaction.Value)
			if err != nil {
				return fmt.Errorf("failed to cast transaction [%s] value: %w", txHash, err)
			}

			txInfo.Transaction = &hyperOutportBlocks.Transaction{
				Nonce:             tx.Transaction.Nonce,
				Value:             value,
				RcvAddr:           tx.Transaction.RcvAddr,
				RcvUserName:       tx.Transaction.RcvUserName,
				SndAddr:           tx.Transaction.SndAddr,
				SndUserName:       tx.Transaction.SndUserName,
				GasPrice:          tx.Transaction.GasPrice,
				GasLimit:          tx.Transaction.GasLimit,
				Data:              tx.Transaction.Data,
				ChainID:           tx.Transaction.ChainID,
				Version:           tx.Transaction.Version,
				Signature:         tx.Transaction.Signature,
				Options:           tx.Transaction.Options,
				GuardianAddr:      tx.Transaction.GuardianAddr,
				GuardianSignature: tx.Transaction.GuardianSignature,
			}

			// TxInfo - FeeInfo
			if tx.FeeInfo != nil {
				fee, err := o.castBigInt(tx.FeeInfo.Fee)
				if err != nil {
					return fmt.Errorf("failed to cast transaction [%s] fee: %w", txHash, err)
				}

				initialPaidFee, err := o.castBigInt(tx.FeeInfo.InitialPaidFee)
				if err != nil {
					return fmt.Errorf("failed to cast transaction [%s] initial paid fee: %w", txHash, err)
				}

				txInfo.FeeInfo = &hyperOutportBlocks.FeeInfo{
					GasUsed:        tx.FeeInfo.GasUsed,
					Fee:            fee,
					InitialPaidFee: initialPaidFee,
				}
			}

			txInfo.ExecutionOrder = tx.ExecutionOrder

			if sc, ok := outportTxPool.SmartContractResults[txHash]; ok {
				txInfo.SmartContractResults = &hyperOutportBlocks.SCRInfo{
					SmartContractResult: &hyperOutportBlocks.SmartContractResult{},
				}
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

			if outportTxPool.Logs != nil {
				txInfo.Logs = make([]*hyperOutportBlocks.Log, len(outportTxPool.Logs))
			}

			// Logs
			for i, logData := range outportTxPool.Logs {
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

				txInfo.Logs[i] = ll
			}

			shardOutportBlockV2.TransactionPool.Transactions[txHash] = txInfo
		}
	}

	for txHash, reward := range outportTxPool.Rewards {
		value, err := o.castBigInt(reward.Reward.Value)
		if err != nil {
			return fmt.Errorf("failed to cast reward tx value: %w", err)
		}

		shardOutportBlockV2.TransactionPool.Rewards[txHash] = &hyperOutportBlocks.RewardInfo{
			Reward: &hyperOutportBlocks.RewardTx{
				Round:   reward.Reward.Round,
				Value:   value,
				RcvAddr: reward.Reward.RcvAddr,
				Epoch:   reward.Reward.Epoch,
			},
			ExecutionOrder: reward.ExecutionOrder,
		}
	}

	for txHash, receipt := range outportTxPool.Receipts {
		value, err := o.castBigInt(receipt.Value)
		if err != nil {
			return fmt.Errorf("failed to cast receipt tx value: %w", err)
		}
		shardOutportBlockV2.TransactionPool.Receipts[txHash] = &hyperOutportBlocks.Receipt{
			Value:   value,
			SndAddr: receipt.SndAddr,
			Data:    receipt.Data,
			TxHash:  receipt.TxHash,
		}
	}

	//TODO: add invalid txs
	//shardOutportBlockV2.TransactionPool.InvalidTxs = outportTxPool.InvalidTxs

	return nil
}

func handleAlteredAccounts(alteredAccounts map[string]*alteredAccount.AlteredAccount, shardOutportBlock *hyperOutportBlocks.ShardOutportBlockV2) {
	if alteredAccounts == nil {
		return
	}

	shardAlteredAccounts := make(map[string]*hyperOutportBlocks.AlteredAccount, len(alteredAccounts))
	for key, alteredAcc := range alteredAccounts {
		tokens := make([]*hyperOutportBlocks.AccountTokenData, len(alteredAcc.Tokens))
		for i, tokenData := range alteredAcc.Tokens {
			tokens[i] = &hyperOutportBlocks.AccountTokenData{
				Nonce:      tokenData.Nonce,
				Identifier: tokenData.Identifier,
				Balance:    tokenData.Balance,
				Properties: tokenData.Properties,
				MetaData: &hyperOutportBlocks.TokenMetaData{
					Nonce:      tokenData.MetaData.Nonce,
					Name:       tokenData.MetaData.Name,
					Creator:    tokenData.MetaData.Creator,
					Royalties:  tokenData.MetaData.Royalties,
					Hash:       tokenData.MetaData.Hash,
					URIs:       tokenData.MetaData.URIs,
					Attributes: tokenData.MetaData.Attributes,
				},
				AdditionalData: &hyperOutportBlocks.AdditionalAccountTokenData{
					IsNFTCreate: tokenData.AdditionalData.IsNFTCreate,
				},
			}
		}

		shardAlteredAccounts[key] = &hyperOutportBlocks.AlteredAccount{
			Address: alteredAcc.Address,
			Nonce:   alteredAcc.Nonce,
			Balance: alteredAcc.Balance,
			Tokens:  tokens,
			AdditionalData: &hyperOutportBlocks.AdditionalAccountData{
				IsSender:         alteredAcc.AdditionalData.IsSender,
				BalanceChanged:   alteredAcc.AdditionalData.BalanceChanged,
				CurrentOwner:     alteredAcc.AdditionalData.CurrentOwner,
				UserName:         alteredAcc.AdditionalData.UserName,
				DeveloperRewards: alteredAcc.AdditionalData.DeveloperRewards,
				CodeHash:         alteredAcc.AdditionalData.CodeHash,
				RootHash:         alteredAcc.AdditionalData.RootHash,
				CodeMetadata:     alteredAcc.AdditionalData.CodeMetadata,
			},
		}
	}

	shardOutportBlock.AlteredAccounts = shardAlteredAccounts
}

func handleHeaderGasConsumption(consumption *outport.HeaderGasConsumption, shardOutportBlock *hyperOutportBlocks.ShardOutportBlockV2) {
	if consumption == nil {
		return
	}

	shardOutportBlock.HeaderGasConsumption = &hyperOutportBlocks.HeaderGasConsumption{
		GasProvided:    consumption.GasProvided,
		GasRefunded:    consumption.GasRefunded,
		GasPenalized:   consumption.GasPenalized,
		MaxGasPerBlock: consumption.MaxGasPerBlock,
	}
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
