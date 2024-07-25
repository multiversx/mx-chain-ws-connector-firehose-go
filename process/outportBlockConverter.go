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
	"github.com/multiversx/mx-chain-core-go/data/receipt"
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

	miniBlocks := copyMiniBlocks(blockData.Body.MiniBlocks)
	intraShardMiniBlocks := copyMiniBlocks(blockData.IntraShardMiniBlocks)

	shardOutportBlock.BlockData = &hyperOutportBlocks.BlockData{
		ShardID:    blockData.ShardID,
		HeaderType: blockData.HeaderType,
		HeaderHash: blockData.HeaderHash,
		Body: &hyperOutportBlocks.Body{
			MiniBlocks: miniBlocks,
		},
		IntraShardMiniBlocks: intraShardMiniBlocks,
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

	shardBlockData.Header = &hyperOutportBlocks.Header{}
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
	shardOutportBlockV2.TransactionPool.ScheduledExecutedSCRSHashesPrevBlock = outportTxPool.ScheduledExecutedSCRSHashesPrevBlock
	shardOutportBlockV2.TransactionPool.ScheduledExecutedInvalidTxsHashesPrevBlock = outportTxPool.ScheduledExecutedInvalidTxsHashesPrevBlock

	if len(outportTxPool.Transactions) == 0 {
		return nil
	}

	shardOutportBlockV2.TransactionPool.Transactions = make(map[string]*hyperOutportBlocks.TxInfoV2)
	err := o.copyTransactions(outportTxPool.Transactions, shardOutportBlockV2.TransactionPool)
	if err != nil {
		return fmt.Errorf("failed to copy transactions: %w", err)
	}

	err = o.copySmartContractResults(outportTxPool.SmartContractResults, shardOutportBlockV2.TransactionPool)
	if err != nil {
		return fmt.Errorf("failed to copy smart contract results: %w", err)
	}

	err = o.copyRewards(outportTxPool.Rewards, shardOutportBlockV2.TransactionPool)
	if err != nil {
		return fmt.Errorf("failed to copy rewards: %w", err)
	}

	err = o.copyReceipts(outportTxPool.Receipts, shardOutportBlockV2.TransactionPool)
	if err != nil {
		return fmt.Errorf("failed to copy receipts: %w", err)
	}

	err = o.copyInvalidTxs(outportTxPool.InvalidTxs, shardOutportBlockV2.TransactionPool)
	if err != nil {
		return fmt.Errorf("failed to copy invalid txs: %w", err)
	}

	err = o.copyLogs(outportTxPool.Logs, shardOutportBlockV2.TransactionPool)
	if err != nil {
		return fmt.Errorf("failed to copy logs: %w", err)
	}

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
			var md *hyperOutportBlocks.TokenMetaData
			if tokenData.MetaData != nil {
				md = &hyperOutportBlocks.TokenMetaData{
					Nonce:      tokenData.MetaData.Nonce,
					Name:       tokenData.MetaData.Name,
					Creator:    tokenData.MetaData.Creator,
					Royalties:  tokenData.MetaData.Royalties,
					Hash:       tokenData.MetaData.Hash,
					URIs:       tokenData.MetaData.URIs,
					Attributes: tokenData.MetaData.Attributes,
				}
			}

			tokens[i] = &hyperOutportBlocks.AccountTokenData{
				Nonce:      tokenData.Nonce,
				Identifier: tokenData.Identifier,
				Balance:    tokenData.Balance,
				Properties: tokenData.Properties,
				MetaData:   md,
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

func (o *outportBlockConverter) copyTransactions(sourceTxs map[string]*outport.TxInfo, transactionPool *hyperOutportBlocks.TransactionPoolV2) error {
	for txHash, txInfo := range sourceTxs {
		destTxInfo := &hyperOutportBlocks.TxInfoV2{}

		// TxInfo - Transaction
		if txInfo.Transaction != nil {
			value, err := o.castBigInt(txInfo.Transaction.Value)
			if err != nil {
				return fmt.Errorf("failed to cast transaction [%s] value: %w", txHash, err)
			}

			destTxInfo.Transaction = &hyperOutportBlocks.WrappedTx{
				Nonce:             txInfo.Transaction.Nonce,
				Value:             value,
				RcvAddr:           txInfo.Transaction.RcvAddr,
				RcvUserName:       txInfo.Transaction.RcvUserName,
				SndAddr:           txInfo.Transaction.SndAddr,
				SndUserName:       txInfo.Transaction.SndUserName,
				GasPrice:          txInfo.Transaction.GasPrice,
				GasLimit:          txInfo.Transaction.GasLimit,
				Data:              txInfo.Transaction.Data,
				ChainID:           txInfo.Transaction.ChainID,
				Version:           txInfo.Transaction.Version,
				Signature:         txInfo.Transaction.Signature,
				Options:           txInfo.Transaction.Options,
				GuardianAddr:      txInfo.Transaction.GuardianAddr,
				GuardianSignature: txInfo.Transaction.GuardianSignature,
				ExecutionOrder:    txInfo.ExecutionOrder,
				TxType:            hyperOutportBlocks.TxType_UserTx,
			}

			// TxInfo - FeeInfo
			if txInfo.FeeInfo != nil {
				fee, err := o.castBigInt(txInfo.FeeInfo.Fee)
				if err != nil {
					return fmt.Errorf("failed to cast transaction [%s] fee: %w", txHash, err)
				}

				initialPaidFee, err := o.castBigInt(txInfo.FeeInfo.InitialPaidFee)
				if err != nil {
					return fmt.Errorf("failed to cast transaction [%s] initial paid fee: %w", txHash, err)
				}

				destTxInfo.FeeInfo = &hyperOutportBlocks.FeeInfo{
					GasUsed:        txInfo.FeeInfo.GasUsed,
					Fee:            fee,
					InitialPaidFee: initialPaidFee,
				}
			}

			transactionPool.Transactions[txHash] = destTxInfo
		}
	}

	return nil
}

func (o *outportBlockConverter) copySmartContractResults(sourceSCRs map[string]*outport.SCRInfo, transactionPool *hyperOutportBlocks.TransactionPoolV2) error {
	for scrHash, scrInfo := range sourceSCRs {
		value, err := o.castBigInt(scrInfo.SmartContractResult.Value)
		if err != nil {
			return fmt.Errorf("failed to cast smart contract [%s] value: %w", scrHash, err)
		}

		relayedValue, err := o.castBigInt(scrInfo.SmartContractResult.RelayedValue)
		if err != nil {
			return fmt.Errorf("failed to cast relayed transaction [%s] value: %w", scrHash, err)
		}

		wrappedTx := &hyperOutportBlocks.WrappedTx{
			Nonce:          scrInfo.SmartContractResult.Nonce,
			Value:          value,
			RcvAddr:        scrInfo.SmartContractResult.RcvAddr,
			SndAddr:        scrInfo.SmartContractResult.SndAddr,
			GasPrice:       scrInfo.SmartContractResult.GasPrice,
			GasLimit:       scrInfo.SmartContractResult.GasLimit,
			Data:           scrInfo.SmartContractResult.Data,
			RelayerAddr:    scrInfo.SmartContractResult.RelayerAddr,
			RelayedValue:   relayedValue,
			Code:           scrInfo.SmartContractResult.Code,
			PrevTxHash:     scrInfo.SmartContractResult.PrevTxHash,
			OriginalTxHash: scrInfo.SmartContractResult.OriginalTxHash,
			CallType:       int64(scrInfo.SmartContractResult.CallType),
			CodeMetadata:   scrInfo.SmartContractResult.CodeMetadata,
			ReturnMessage:  scrInfo.SmartContractResult.ReturnMessage,
			OriginalSender: scrInfo.SmartContractResult.OriginalSender,
			ExecutionOrder: scrInfo.ExecutionOrder,
			TxType:         hyperOutportBlocks.TxType_SCR,
		}

		fee, err := o.castBigInt(scrInfo.FeeInfo.Fee)
		if err != nil {
			return fmt.Errorf("failed to cast transaction [%s] fee: %w", scrHash, err)
		}

		initialPaidFee, err := o.castBigInt(scrInfo.FeeInfo.InitialPaidFee)
		if err != nil {
			return fmt.Errorf("failed to cast transaction [%s] initial paid fee: %w", scrHash, err)
		}

		txInfo := &hyperOutportBlocks.TxInfoV2{
			Transaction: wrappedTx,
			FeeInfo: &hyperOutportBlocks.FeeInfo{
				GasUsed:        scrInfo.FeeInfo.GasUsed,
				Fee:            fee,
				InitialPaidFee: initialPaidFee,
			},
		}

		if _, ok := transactionPool.Transactions[scrHash]; !ok {
			transactionPool.Transactions[scrHash] = txInfo
		} else {
			if transactionPool.Transactions[scrHash].ResultTxs == nil {
				transactionPool.Transactions[scrHash].ResultTxs = make([]*hyperOutportBlocks.WrappedTx, 0)
			}
			transactionPool.Transactions[scrHash].ResultTxs = append(transactionPool.Transactions[scrHash].ResultTxs, wrappedTx)
		}
	}

	return nil
}

func (o *outportBlockConverter) copyRewards(sourceRewards map[string]*outport.RewardInfo, transactionPool *hyperOutportBlocks.TransactionPoolV2) error {
	for hash, reward := range sourceRewards {
		value, err := o.castBigInt(reward.Reward.Value)
		if err != nil {
			return fmt.Errorf("failed to cast reward tx value: %w", err)
		}

		transactionPool.Transactions[hash] = &hyperOutportBlocks.TxInfoV2{
			Transaction: &hyperOutportBlocks.WrappedTx{
				Value:          value,
				SndAddr:        []byte("metachain"),
				RcvAddr:        reward.Reward.RcvAddr,
				Round:          reward.Reward.Round,
				Epoch:          reward.Reward.Epoch,
				ExecutionOrder: reward.ExecutionOrder,
				TxType:         hyperOutportBlocks.TxType_Reward,
			},
		}
	}

	return nil
}

func (o *outportBlockConverter) copyReceipts(sourceReceipts map[string]*receipt.Receipt, transactionPool *hyperOutportBlocks.TransactionPoolV2) error {
	for hash, r := range sourceReceipts {
		value, err := o.castBigInt(r.Value)
		if err != nil {
			return fmt.Errorf("failed to cast receipt tx value: %w", err)
		}

		transactionPool.Transactions[hash] = &hyperOutportBlocks.TxInfoV2{
			Receipt: &hyperOutportBlocks.Receipt{
				Value:   value,
				SndAddr: r.SndAddr,
				Data:    r.Data,
				TxHash:  r.TxHash,
			},
		}
	}

	return nil
}

func (o *outportBlockConverter) copyInvalidTxs(sourceInvalidTxs map[string]*outport.TxInfo, transactionPool *hyperOutportBlocks.TransactionPoolV2) error {
	for txHash, invalidTx := range sourceInvalidTxs {
		value, err := o.castBigInt(invalidTx.Transaction.Value)
		if err != nil {
			return fmt.Errorf("failed to cast receipt tx value: %w", err)
		}

		transactionPool.Transactions[txHash] = &hyperOutportBlocks.TxInfoV2{
			Transaction: &hyperOutportBlocks.WrappedTx{
				Nonce:             invalidTx.Transaction.Nonce,
				Value:             value,
				RcvAddr:           invalidTx.Transaction.RcvAddr,
				RcvUserName:       invalidTx.Transaction.RcvUserName,
				SndAddr:           invalidTx.Transaction.SndAddr,
				SndUserName:       invalidTx.Transaction.SndUserName,
				GasPrice:          invalidTx.Transaction.GasPrice,
				GasLimit:          invalidTx.Transaction.GasLimit,
				Data:              invalidTx.Transaction.Data,
				ChainID:           invalidTx.Transaction.ChainID,
				Version:           invalidTx.Transaction.Version,
				Signature:         invalidTx.Transaction.Signature,
				Options:           invalidTx.Transaction.Options,
				GuardianAddr:      invalidTx.Transaction.GuardianAddr,
				GuardianSignature: invalidTx.Transaction.GuardianSignature,
				ExecutionOrder:    invalidTx.ExecutionOrder,
				TxType:            hyperOutportBlocks.TxType_Invalid,
			},
		}
	}

	return nil
}

func (o *outportBlockConverter) copyLogs(sourceLogs []*outport.LogData, transactionPool *hyperOutportBlocks.TransactionPoolV2) error {
	for _, logData := range sourceLogs {
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

		if txInfo, ok := transactionPool.Transactions[logData.TxHash]; !ok {
			if transactionPool.Transactions[logData.TxHash].Logs == nil {
				transactionPool.Transactions[logData.TxHash].Logs = make([]*hyperOutportBlocks.Log, 0)
			}
		} else {
			txInfo.Logs = append(transactionPool.Transactions[logData.TxHash].Logs, ll)
		}

	}

	return nil
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

func copyMiniBlocks(sourceMiniBlocks []*block.MiniBlock) []*hyperOutportBlocks.MiniBlock {
	var destMiniBlocks []*hyperOutportBlocks.MiniBlock
	if sourceMiniBlocks != nil {
		destMiniBlocks = []*hyperOutportBlocks.MiniBlock{}
		for _, mb := range sourceMiniBlocks {
			destMiniBlocks = append(destMiniBlocks, &hyperOutportBlocks.MiniBlock{
				TxHashes:        mb.TxHashes,
				ReceiverShardID: mb.ReceiverShardID,
				SenderShardID:   mb.SenderShardID,
				Type:            hyperOutportBlocks.Type(mb.Type),
				Reserved:        mb.Reserved,
			})
		}
	}

	return destMiniBlocks
}

// IsInterfaceNil returns nil if there is no value under the interface
func (o *outportBlockConverter) IsInterfaceNil() bool {
	return o == nil
}
