package process

import (
	"fmt"
	"math/big"

	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"google.golang.org/protobuf/proto"

	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
)

var (
	protoMarshaller = &marshal.GogoProtoMarshalizer{}
	caster          coreData.BigIntCaster
)

type outportBlockCaster struct {
	outportBlock *outport.OutportBlock
}

// NewOutportBlockCaster returns a new instance of outportBlockCaster.
func NewOutportBlockCaster(outportBlock *outport.OutportBlock) *outportBlockCaster {
	return &outportBlockCaster{outportBlock: outportBlock}
}

// MarshalTo will marshal the underlying outport.OutportBlock to the given proto.Message
func (o *outportBlockCaster) MarshalTo(m proto.Message) error {
	//obBytes, err := protoMarshaller.Marshal(o.outportBlock)
	//if err != nil {
	//	return fmt.Errorf("failed to marshal gogo outport block: %w", err)
	//}

	standardOb := &data.OutportBlock{}

	// Asserting values.
	standardOb.ShardID = o.outportBlock.ShardID
	standardOb.NotarizedHeadersHashes = o.outportBlock.NotarizedHeadersHashes
	standardOb.NumberOfShards = o.outportBlock.NumberOfShards
	standardOb.SignersIndexes = o.outportBlock.SignersIndexes
	standardOb.HighestFinalBlockNonce = o.outportBlock.HighestFinalBlockNonce
	standardOb.HighestFinalBlockHash = o.outportBlock.HighestFinalBlockHash

	err := o.convertBlockData(standardOb.BlockData)
	if err != nil {
		return fmt.Errorf("failed to convert block data: %w", err)
	}
	//
	//// Transaction pool - Transactions.
	//for k, v := range ob.TransactionPool.Transactions {
	//	// Transaction pool - Transactions. - TxInfo.
	//	require.Equal(t, v.Transaction.Nonce, standardOb.TransactionPool.Transactions[k].Transaction.Nonce)
	//	require.Equal(t, castBigInt(t, v.Transaction.Value), standardOb.TransactionPool.Transactions[k].Transaction.Value)
	//	require.Equal(t, v.Transaction.RcvAddr, standardOb.TransactionPool.Transactions[k].Transaction.RcvAddr)
	//	require.Equal(t, v.Transaction.RcvUserName, standardOb.TransactionPool.Transactions[k].Transaction.RcvUserName)
	//	require.Equal(t, v.Transaction.SndAddr, standardOb.TransactionPool.Transactions[k].Transaction.SndAddr)
	//	require.Equal(t, v.Transaction.GasPrice, standardOb.TransactionPool.Transactions[k].Transaction.GasPrice)
	//	require.Equal(t, v.Transaction.GasLimit, standardOb.TransactionPool.Transactions[k].Transaction.GasLimit)
	//	require.Equal(t, v.Transaction.Data, standardOb.TransactionPool.Transactions[k].Transaction.Data)
	//	require.Equal(t, v.Transaction.ChainID, standardOb.TransactionPool.Transactions[k].Transaction.ChainID)
	//	require.Equal(t, v.Transaction.Version, standardOb.TransactionPool.Transactions[k].Transaction.Version)
	//	require.Equal(t, v.Transaction.Signature, standardOb.TransactionPool.Transactions[k].Transaction.Signature)
	//	require.Equal(t, v.Transaction.Options, standardOb.TransactionPool.Transactions[k].Transaction.Options)
	//	require.Equal(t, v.Transaction.GuardianAddr, standardOb.TransactionPool.Transactions[k].Transaction.GuardianAddr)
	//	require.Equal(t, v.Transaction.GuardianSignature, standardOb.TransactionPool.Transactions[k].Transaction.GuardianSignature)
	//
	//	// Transaction pool - Transactions - Tx Info - Fee info.
	//	require.Equal(t, v.FeeInfo.GasUsed, standardOb.TransactionPool.Transactions[k].FeeInfo.GasUsed)
	//	require.Equal(t, castBigInt(t, v.FeeInfo.Fee), standardOb.TransactionPool.Transactions[k].FeeInfo.Fee)
	//	require.Equal(t, castBigInt(t, v.FeeInfo.InitialPaidFee), standardOb.TransactionPool.Transactions[k].FeeInfo.InitialPaidFee)
	//
	//	require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.Transactions[k].ExecutionOrder)
	//}
	//
	//// Transaction pool - Smart Contract results.
	//for k, v := range ob.TransactionPool.SmartContractResults {
	//	// Transaction pool - Smart Contract results - SmartContractResult.
	//	require.Equal(t, v.SmartContractResult.Nonce, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.Nonce)
	//	require.Equal(t, castBigInt(t, v.SmartContractResult.Value), standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.Value)
	//	require.Equal(t, v.SmartContractResult.RcvAddr, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.RcvAddr)
	//	require.Equal(t, v.SmartContractResult.SndAddr, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.SndAddr)
	//	require.Equal(t, v.SmartContractResult.RelayerAddr, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.RelayerAddr)
	//	require.Equal(t, castBigInt(t, v.SmartContractResult.RelayedValue), standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.RelayedValue)
	//	require.Equal(t, v.SmartContractResult.Code, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.Code)
	//	require.Equal(t, v.SmartContractResult.Data, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.Data)
	//	require.Equal(t, v.SmartContractResult.PrevTxHash, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.PrevTxHash)
	//	require.Equal(t, v.SmartContractResult.OriginalTxHash, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.OriginalTxHash)
	//	require.Equal(t, v.SmartContractResult.GasLimit, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.GasLimit)
	//	require.Equal(t, v.SmartContractResult.GasPrice, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.GasPrice)
	//	require.Equal(t, int(v.SmartContractResult.CallType), standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.CallType)
	//	require.Equal(t, v.SmartContractResult.CodeMetadata, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.CodeMetadata)
	//	require.Equal(t, v.SmartContractResult.ReturnMessage, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.ReturnMessage)
	//	require.Equal(t, v.SmartContractResult.OriginalSender, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.OriginalSender)
	//
	//	// Transaction pool - Smart Contract results - Fee info.
	//	require.Equal(t, v.FeeInfo.GasUsed, standardOb.TransactionPool.SmartContractResults[k].FeeInfo.GasUsed)
	//	require.Equal(t, castBigInt(t, v.FeeInfo.Fee), standardOb.TransactionPool.SmartContractResults[k].FeeInfo.Fee)
	//	require.Equal(t, castBigInt(t, v.FeeInfo.InitialPaidFee), standardOb.TransactionPool.SmartContractResults[k].FeeInfo.InitialPaidFee)
	//
	//	// Transaction pool - Smart Contract results - Execution Order.
	//	require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.SmartContractResults[k].ExecutionOrder)
	//}
	//
	//// Transaction Pool - Rewards
	//for k, v := range ob.TransactionPool.Rewards {
	//	// Transaction Pool - Rewards - Reward info
	//	require.Equal(t, v.Reward.Round, standardOb.TransactionPool.Rewards[k].Reward.Round)
	//	require.Equal(t, castBigInt(t, v.Reward.Value), standardOb.TransactionPool.Rewards[k].Reward.Value)
	//	require.Equal(t, v.Reward.RcvAddr, standardOb.TransactionPool.Rewards[k].Reward.RcvAddr)
	//	require.Equal(t, v.Reward.Epoch, standardOb.TransactionPool.Rewards[k].Reward.Epoch)
	//
	//	// Transaction Pool - Rewards - Execution Order
	//	require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.Rewards[k].ExecutionOrder)
	//}
	//
	//// Transaction Pool - Receipts
	//for k, v := range ob.TransactionPool.Receipts {
	//	// Transaction Pool - Receipts - Receipt info
	//	require.Equal(t, castBigInt(t, v.Value), standardOb.TransactionPool.Receipts[k].Value)
	//	require.Equal(t, v.SndAddr, standardOb.TransactionPool.Receipts[k].SndAddr)
	//	require.Equal(t, v.Data, standardOb.TransactionPool.Receipts[k].Data)
	//	require.Equal(t, v.TxHash, standardOb.TransactionPool.Receipts[k].TxHash)
	//}
	//
	//// Transaction Pool - Invalid Txs
	//for k, v := range ob.TransactionPool.InvalidTxs {
	//	// Transaction Pool - Invalid Txs - Tx Info
	//	require.Equal(t, v.Transaction.Nonce, standardOb.TransactionPool.InvalidTxs[k].Transaction.Nonce)
	//	require.Equal(t, castBigInt(t, v.Transaction.Value), standardOb.TransactionPool.InvalidTxs[k].Transaction.Value)
	//	require.Equal(t, v.Transaction.RcvAddr, standardOb.TransactionPool.InvalidTxs[k].Transaction.RcvAddr)
	//	require.Equal(t, v.Transaction.RcvUserName, standardOb.TransactionPool.InvalidTxs[k].Transaction.RcvUserName)
	//	require.Equal(t, v.Transaction.SndAddr, standardOb.TransactionPool.InvalidTxs[k].Transaction.SndAddr)
	//	require.Equal(t, v.Transaction.GasPrice, standardOb.TransactionPool.InvalidTxs[k].Transaction.GasPrice)
	//	require.Equal(t, v.Transaction.GasLimit, standardOb.TransactionPool.InvalidTxs[k].Transaction.GasLimit)
	//	require.Equal(t, v.Transaction.Data, standardOb.TransactionPool.InvalidTxs[k].Transaction.Data)
	//	require.Equal(t, v.Transaction.ChainID, standardOb.TransactionPool.InvalidTxs[k].Transaction.ChainID)
	//	require.Equal(t, v.Transaction.Version, standardOb.TransactionPool.InvalidTxs[k].Transaction.Version)
	//	require.Equal(t, v.Transaction.Signature, standardOb.TransactionPool.InvalidTxs[k].Transaction.Signature)
	//	require.Equal(t, v.Transaction.Options, standardOb.TransactionPool.InvalidTxs[k].Transaction.Options)
	//	require.Equal(t, v.Transaction.GuardianAddr, standardOb.TransactionPool.InvalidTxs[k].Transaction.GuardianAddr)
	//	require.Equal(t, v.Transaction.GuardianSignature, standardOb.TransactionPool.InvalidTxs[k].Transaction.GuardianSignature)
	//	require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.InvalidTxs[k].ExecutionOrder)
	//
	//	// Transaction pool - Invalid Txs - Fee info.
	//	require.Equal(t, v.FeeInfo.GasUsed, standardOb.TransactionPool.InvalidTxs[k].FeeInfo.GasUsed)
	//	require.Equal(t, castBigInt(t, v.FeeInfo.Fee), standardOb.TransactionPool.InvalidTxs[k].FeeInfo.Fee)
	//	require.Equal(t, castBigInt(t, v.FeeInfo.InitialPaidFee), standardOb.TransactionPool.InvalidTxs[k].FeeInfo.InitialPaidFee)
	//
	//	require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.InvalidTxs[k].ExecutionOrder)
	//}
	//
	//// Transaction Pool - Logs
	//for i, l := range ob.TransactionPool.Logs {
	//	// Transaction Pool - Logs - Log Data
	//	require.Equal(t, l.TxHash, standardOb.TransactionPool.Logs[i].TxHash)
	//
	//	// Transaction Pool - Logs - Log data - Log
	//	require.Equal(t, l.Log, standardOb.TransactionPool.Logs[i].Log.Address)
	//
	//	for j, e := range ob.TransactionPool.Logs[i].Log.Events {
	//		require.Equal(t, e.Address, standardOb.TransactionPool.Logs[i].Log.Events[j].Address)
	//		require.Equal(t, e.Identifier, standardOb.TransactionPool.Logs[i].Log.Events[j].Identifier)
	//		require.Equal(t, e.Topics, standardOb.TransactionPool.Logs[i].Log.Events[j].Topics)
	//		require.Equal(t, e.Data, standardOb.TransactionPool.Logs[i].Log.Events[j].Data)
	//		require.Equal(t, e.AdditionalData, standardOb.TransactionPool.Logs[i].Log.Events[j].AdditionalData)
	//	}
	//}
	//
	//// Transaction Pool - ScheduledExecutedSCRSHashesPrevBlock
	//for i, s := range ob.TransactionPool.ScheduledExecutedSCRSHashesPrevBlock {
	//	require.Equal(t, s, standardOb.TransactionPool.ScheduledExecutedSCRSHashesPrevBlock[i])
	//}
	//
	//// Transaction Pool - ScheduledExecutedInvalidTxsHashesPrevBlock
	//for i, s := range ob.TransactionPool.ScheduledExecutedInvalidTxsHashesPrevBlock {
	//	require.Equal(t, s, standardOb.TransactionPool.ScheduledExecutedInvalidTxsHashesPrevBlock[i])
	//}
	//
	//// Header gas consumption.
	//require.Equal(t, ob.HeaderGasConsumption.GasProvided, standardOb.HeaderGasConsumption.GasProvided)
	//require.Equal(t, ob.HeaderGasConsumption.GasRefunded, standardOb.HeaderGasConsumption.GasRefunded)
	//require.Equal(t, ob.HeaderGasConsumption.GasPenalized, standardOb.HeaderGasConsumption.GasPenalized)
	//require.Equal(t, ob.HeaderGasConsumption.MaxGasPerBlock, standardOb.HeaderGasConsumption.MaxGasPerBlock)
	//
	//// Altered accounts.
	//for key, account := range ob.AlteredAccounts {
	//	acc := standardOb.AlteredAccounts[key]
	//
	//	require.Equal(t, account.Address, acc.Address)
	//	require.Equal(t, account.Nonce, acc.Nonce)
	//	require.Equal(t, account.Balance, acc.Balance)
	//	require.Equal(t, account.Balance, acc.Balance)
	//
	//	// Altered accounts - Account token data.
	//	for i, token := range account.Tokens {
	//		require.Equal(t, token.Nonce, acc.Tokens[i].Nonce)
	//		require.Equal(t, token.Identifier, acc.Tokens[i].Identifier)
	//		require.Equal(t, token.Balance, acc.Tokens[i].Balance)
	//		require.Equal(t, token.Properties, acc.Tokens[i].Properties)
	//
	//		// Altered accounts - Account token data - Metadata.
	//		require.Equal(t, token.MetaData.Nonce, acc.Tokens[i].MetaData.Nonce)
	//		require.Equal(t, token.MetaData.Name, acc.Tokens[i].MetaData.Name)
	//		require.Equal(t, token.MetaData.Creator, acc.Tokens[i].MetaData.Creator)
	//		require.Equal(t, token.MetaData.Royalties, acc.Tokens[i].MetaData.Royalties)
	//		require.Equal(t, token.MetaData.Hash, acc.Tokens[i].MetaData.Hash)
	//		require.Equal(t, token.MetaData.URIs, acc.Tokens[i].MetaData.URIs)
	//		require.Equal(t, token.MetaData.Attributes, acc.Tokens[i].MetaData.Attributes)
	//
	//		// Altered accounts - Account token data - Additional data.
	//		require.Equal(t, token.AdditionalData.IsNFTCreate, acc.Tokens[i].AdditionalData.IsNFTCreate)
	//	}
	//
	//	//  Altered accounts - Additional account data.
	//	require.Equal(t, account.AdditionalData.IsSender, acc.AdditionalData.IsSender)
	//	require.Equal(t, account.AdditionalData.BalanceChanged, acc.AdditionalData.BalanceChanged)
	//	require.Equal(t, account.AdditionalData.CurrentOwner, acc.AdditionalData.CurrentOwner)
	//	require.Equal(t, account.AdditionalData.UserName, acc.AdditionalData.UserName)
	//	require.Equal(t, account.AdditionalData.DeveloperRewards, acc.AdditionalData.DeveloperRewards)
	//	require.Equal(t, account.AdditionalData.CodeHash, acc.AdditionalData.CodeHash)
	//	require.Equal(t, account.AdditionalData.RootHash, acc.AdditionalData.RootHash)
	//	require.Equal(t, account.AdditionalData.CodeMetadata, acc.AdditionalData.CodeMetadata)
	//}
	//
	//err = proto.Unmarshal(obBytes, m)
	//if err != nil {
	//	return fmt.Errorf("failed to unmarshal to proto message [%s]: %w",
	//		m.ProtoReflect().Descriptor().Name(), err)
	//}

	return nil
}

func (o *outportBlockCaster) convertBlockData(bd *data.BlockData) error {
	newBd := &data.BlockData{}
	newBd.ShardID = o.outportBlock.BlockData.ShardID
	newBd.HeaderHash = o.outportBlock.BlockData.HeaderHash
	newBd.HeaderType = o.outportBlock.BlockData.HeaderType

	// convert Header
	header := block.Header{}
	err := protoMarshaller.Unmarshal(&header, o.outportBlock.BlockData.HeaderBytes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal block header: %w", err)
	}

	newBd.Header.Nonce = header.Nonce
	newBd.Header.PrevHash = header.PrevHash
	newBd.Header.PrevRandSeed = header.PrevRandSeed
	newBd.Header.RandSeed = header.RandSeed
	newBd.Header.PubKeysBitmap = header.PubKeysBitmap
	newBd.Header.ShardID = header.ShardID
	newBd.Header.TimeStamp = header.TimeStamp
	newBd.Header.TimeStamp = header.TimeStamp
	newBd.Header.Round = header.Round
	newBd.Header.Epoch = header.Epoch
	newBd.Header.BlockBodyType = data.Type(header.BlockBodyType)
	newBd.Header.Signature = header.Signature
	newBd.Header.LeaderSignature = header.LeaderSignature
	newBd.Header.RootHash = header.RootHash
	newBd.Header.MetaBlockHashes = header.MetaBlockHashes
	newBd.Header.TxCount = header.TxCount
	newBd.Header.EpochStartMetaHash = header.EpochStartMetaHash
	newBd.Header.ReceiptsHash = header.ReceiptsHash
	newBd.Header.ChainID = header.ChainID
	newBd.Header.SoftwareVersion = header.SoftwareVersion
	newBd.Header.Reserved = header.Reserved

	accumulatedFees, err := castBigInt(header.AccumulatedFees)
	if err != nil {
		return fmt.Errorf("failed to cast accumulated fees: %w", err)
	}
	newBd.Header.AccumulatedFees = accumulatedFees

	developerFees, err := castBigInt(header.DeveloperFees)
	if err != nil {
		return fmt.Errorf("failed to cast developerFees fees: %w", err)
	}
	newBd.Header.DeveloperFees = developerFees

	// Block data - Header - Mini block headers.
	for i, miniBlockHeader := range header.MiniBlockHeaders {
		newBd.Header.MiniBlockHeaders[i].Hash = miniBlockHeader.Hash
		newBd.Header.MiniBlockHeaders[i].SenderShardID = miniBlockHeader.SenderShardID
		newBd.Header.MiniBlockHeaders[i].ReceiverShardID = miniBlockHeader.ReceiverShardID
		newBd.Header.MiniBlockHeaders[i].TxCount = miniBlockHeader.TxCount
		newBd.Header.MiniBlockHeaders[i].Type = data.Type(miniBlockHeader.Type)
		newBd.Header.MiniBlockHeaders[i].Reserved = miniBlockHeader.Reserved
	}

	// Block data - Header - Peer changes.
	for i, peerChange := range header.PeerChanges {
		newBd.Header.PeerChanges[i].PubKey = peerChange.PubKey
		newBd.Header.PeerChanges[i].ShardIdDest = peerChange.ShardIdDest
	}

	for i, intraShardMiniBlock := range o.outportBlock.BlockData.IntraShardMiniBlocks {
		newBd.IntraShardMiniBlocks[i].TxHashes = intraShardMiniBlock.TxHashes
		newBd.IntraShardMiniBlocks[i].ReceiverShardID = intraShardMiniBlock.ReceiverShardID
		newBd.IntraShardMiniBlocks[i].SenderShardID = intraShardMiniBlock.SenderShardID
		newBd.IntraShardMiniBlocks[i].Type = data.Type(intraShardMiniBlock.SenderShardID)
		newBd.IntraShardMiniBlocks[i].Reserved = intraShardMiniBlock.Reserved
	}

	bd = newBd
	return nil
}

func (o *outportBlockCaster) convertTransactionPool(tp *data.TransactionPool) error {
	// Transaction pool - Transactions.
	for k, v := range o.outportBlock.TransactionPool.Transactions {
		// Transaction pool - Transactions. - TxInfo.
		tp.Transactions[k].Transaction.Nonce = v.Transaction.Nonce
		value, err := castBigInt(v.Transaction.Value)
		if err != nil {
			return fmt.Errorf("failed to cast transaction value: %w", err)
		}
		tp.Transactions[k].Transaction.Value = value
		tp.Transactions[k].Transaction.RcvAddr = v.Transaction.RcvAddr
		tp.Transactions[k].Transaction.RcvUserName = v.Transaction.RcvUserName
		tp.Transactions[k].Transaction.SndAddr = v.Transaction.SndAddr
		tp.Transactions[k].Transaction.GasPrice = v.Transaction.GasPrice
		tp.Transactions[k].Transaction.GasLimit = v.Transaction.GasLimit
		tp.Transactions[k].Transaction.Data = v.Transaction.Data
		tp.Transactions[k].Transaction.ChainID = v.Transaction.ChainID
		tp.Transactions[k].Transaction.Version = v.Transaction.Version
		tp.Transactions[k].Transaction.Signature = v.Transaction.Signature
		tp.Transactions[k].Transaction.Options = v.Transaction.Options
		tp.Transactions[k].Transaction.GuardianAddr = v.Transaction.GuardianAddr
		tp.Transactions[k].Transaction.GuardianSignature = v.Transaction.GuardianSignature

		tp.Transactions[k].FeeInfo.GasUsed = v.FeeInfo.GasUsed

		fee, err := castBigInt(v.FeeInfo.Fee)
		if err != nil {
			return fmt.Errorf("failed to cast fee info fee: %w", err)
		}
		tp.Transactions[k].FeeInfo.Fee = fee

		initialPaidFee, err := castBigInt(v.FeeInfo.InitialPaidFee)
		if err != nil {
			return fmt.Errorf("failed to cast fee info initial paid fee: %w", err)
		}
		tp.Transactions[k].FeeInfo.InitialPaidFee = initialPaidFee

		tp.Transactions[k].ExecutionOrder = v.ExecutionOrder

		//	require.Equal(t, v.Transaction.Nonce, standardOb.TransactionPool.Transactions[k].Transaction.Nonce)
		//	require.Equal(t, castBigInt(t, v.Transaction.Value), standardOb.TransactionPool.Transactions[k].Transaction.Value)
		//	require.Equal(t, v.Transaction.RcvAddr, standardOb.TransactionPool.Transactions[k].Transaction.RcvAddr)
		//	require.Equal(t, v.Transaction.RcvUserName, standardOb.TransactionPool.Transactions[k].Transaction.RcvUserName)
		//	require.Equal(t, v.Transaction.SndAddr, standardOb.TransactionPool.Transactions[k].Transaction.SndAddr)
		//	require.Equal(t, v.Transaction.GasPrice, standardOb.TransactionPool.Transactions[k].Transaction.GasPrice)
		//	require.Equal(t, v.Transaction.GasLimit, standardOb.TransactionPool.Transactions[k].Transaction.GasLimit)
		//	require.Equal(t, v.Transaction.Data, standardOb.TransactionPool.Transactions[k].Transaction.Data)
		//	require.Equal(t, v.Transaction.ChainID, standardOb.TransactionPool.Transactions[k].Transaction.ChainID)
		//	require.Equal(t, v.Transaction.Version, standardOb.TransactionPool.Transactions[k].Transaction.Version)
		//	require.Equal(t, v.Transaction.Signature, standardOb.TransactionPool.Transactions[k].Transaction.Signature)
		//	require.Equal(t, v.Transaction.Options, standardOb.TransactionPool.Transactions[k].Transaction.Options)
		//	require.Equal(t, v.Transaction.GuardianAddr, standardOb.TransactionPool.Transactions[k].Transaction.GuardianAddr)
		//	require.Equal(t, v.Transaction.GuardianSignature, standardOb.TransactionPool.Transactions[k].Transaction.GuardianSignature)
		//
		//	// Transaction pool - Transactions - Tx Info - Fee info.
		//	require.Equal(t, v.FeeInfo.GasUsed, standardOb.TransactionPool.Transactions[k].FeeInfo.GasUsed)
		//	require.Equal(t, castBigInt(t, v.FeeInfo.Fee), standardOb.TransactionPool.Transactions[k].FeeInfo.Fee)
		//	require.Equal(t, castBigInt(t, v.FeeInfo.InitialPaidFee), standardOb.TransactionPool.Transactions[k].FeeInfo.InitialPaidFee)
		//
		//	require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.Transactions[k].ExecutionOrder)
		//}
		//
		//// Transaction pool - Smart Contract results.
		//for k, v := range ob.TransactionPool.SmartContractResults {
		//	// Transaction pool - Smart Contract results - SmartContractResult.
		//	require.Equal(t, v.SmartContractResult.Nonce, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.Nonce)
		//	require.Equal(t, castBigInt(t, v.SmartContractResult.Value), standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.Value)
		//	require.Equal(t, v.SmartContractResult.RcvAddr, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.RcvAddr)
		//	require.Equal(t, v.SmartContractResult.SndAddr, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.SndAddr)
		//	require.Equal(t, v.SmartContractResult.RelayerAddr, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.RelayerAddr)
		//	require.Equal(t, castBigInt(t, v.SmartContractResult.RelayedValue), standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.RelayedValue)
		//	require.Equal(t, v.SmartContractResult.Code, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.Code)
		//	require.Equal(t, v.SmartContractResult.Data, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.Data)
		//	require.Equal(t, v.SmartContractResult.PrevTxHash, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.PrevTxHash)
		//	require.Equal(t, v.SmartContractResult.OriginalTxHash, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.OriginalTxHash)
		//	require.Equal(t, v.SmartContractResult.GasLimit, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.GasLimit)
		//	require.Equal(t, v.SmartContractResult.GasPrice, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.GasPrice)
		//	require.Equal(t, int(v.SmartContractResult.CallType), standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.CallType)
		//	require.Equal(t, v.SmartContractResult.CodeMetadata, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.CodeMetadata)
		//	require.Equal(t, v.SmartContractResult.ReturnMessage, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.ReturnMessage)
		//	require.Equal(t, v.SmartContractResult.OriginalSender, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.OriginalSender)
		//
		//	// Transaction pool - Smart Contract results - Fee info.
		//	require.Equal(t, v.FeeInfo.GasUsed, standardOb.TransactionPool.SmartContractResults[k].FeeInfo.GasUsed)
		//	require.Equal(t, castBigInt(t, v.FeeInfo.Fee), standardOb.TransactionPool.SmartContractResults[k].FeeInfo.Fee)
		//	require.Equal(t, castBigInt(t, v.FeeInfo.InitialPaidFee), standardOb.TransactionPool.SmartContractResults[k].FeeInfo.InitialPaidFee)
		//
		//	// Transaction pool - Smart Contract results - Execution Order.
		//	require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.SmartContractResults[k].ExecutionOrder)
		//}
		//
		//// Transaction Pool - Rewards
		//for k, v := range ob.TransactionPool.Rewards {
		//	// Transaction Pool - Rewards - Reward info
		//	require.Equal(t, v.Reward.Round, standardOb.TransactionPool.Rewards[k].Reward.Round)
		//	require.Equal(t, castBigInt(t, v.Reward.Value), standardOb.TransactionPool.Rewards[k].Reward.Value)
		//	require.Equal(t, v.Reward.RcvAddr, standardOb.TransactionPool.Rewards[k].Reward.RcvAddr)
		//	require.Equal(t, v.Reward.Epoch, standardOb.TransactionPool.Rewards[k].Reward.Epoch)
		//
		//	// Transaction Pool - Rewards - Execution Order
		//	require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.Rewards[k].ExecutionOrder)
		//}
		//
		//// Transaction Pool - Receipts
		//for k, v := range ob.TransactionPool.Receipts {
		//	// Transaction Pool - Receipts - Receipt info
		//	require.Equal(t, castBigInt(t, v.Value), standardOb.TransactionPool.Receipts[k].Value)
		//	require.Equal(t, v.SndAddr, standardOb.TransactionPool.Receipts[k].SndAddr)
		//	require.Equal(t, v.Data, standardOb.TransactionPool.Receipts[k].Data)
		//	require.Equal(t, v.TxHash, standardOb.TransactionPool.Receipts[k].TxHash)
		//}
		//
		//// Transaction Pool - Invalid Txs
		//for k, v := range ob.TransactionPool.InvalidTxs {
		//	// Transaction Pool - Invalid Txs - Tx Info
		//	require.Equal(t, v.Transaction.Nonce, standardOb.TransactionPool.InvalidTxs[k].Transaction.Nonce)
		//	require.Equal(t, castBigInt(t, v.Transaction.Value), standardOb.TransactionPool.InvalidTxs[k].Transaction.Value)
		//	require.Equal(t, v.Transaction.RcvAddr, standardOb.TransactionPool.InvalidTxs[k].Transaction.RcvAddr)
		//	require.Equal(t, v.Transaction.RcvUserName, standardOb.TransactionPool.InvalidTxs[k].Transaction.RcvUserName)
		//	require.Equal(t, v.Transaction.SndAddr, standardOb.TransactionPool.InvalidTxs[k].Transaction.SndAddr)
		//	require.Equal(t, v.Transaction.GasPrice, standardOb.TransactionPool.InvalidTxs[k].Transaction.GasPrice)
		//	require.Equal(t, v.Transaction.GasLimit, standardOb.TransactionPool.InvalidTxs[k].Transaction.GasLimit)
		//	require.Equal(t, v.Transaction.Data, standardOb.TransactionPool.InvalidTxs[k].Transaction.Data)
		//	require.Equal(t, v.Transaction.ChainID, standardOb.TransactionPool.InvalidTxs[k].Transaction.ChainID)
		//	require.Equal(t, v.Transaction.Version, standardOb.TransactionPool.InvalidTxs[k].Transaction.Version)
		//	require.Equal(t, v.Transaction.Signature, standardOb.TransactionPool.InvalidTxs[k].Transaction.Signature)
		//	require.Equal(t, v.Transaction.Options, standardOb.TransactionPool.InvalidTxs[k].Transaction.Options)
		//	require.Equal(t, v.Transaction.GuardianAddr, standardOb.TransactionPool.InvalidTxs[k].Transaction.GuardianAddr)
		//	require.Equal(t, v.Transaction.GuardianSignature, standardOb.TransactionPool.InvalidTxs[k].Transaction.GuardianSignature)
		//	require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.InvalidTxs[k].ExecutionOrder)
		//
		//	// Transaction pool - Invalid Txs - Fee info.
		//	require.Equal(t, v.FeeInfo.GasUsed, standardOb.TransactionPool.InvalidTxs[k].FeeInfo.GasUsed)
		//	require.Equal(t, castBigInt(t, v.FeeInfo.Fee), standardOb.TransactionPool.InvalidTxs[k].FeeInfo.Fee)
		//	require.Equal(t, castBigInt(t, v.FeeInfo.InitialPaidFee), standardOb.TransactionPool.InvalidTxs[k].FeeInfo.InitialPaidFee)
		//
		//	require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.InvalidTxs[k].ExecutionOrder)
		//}
		//
		//// Transaction Pool - Logs
		//for i, l := range ob.TransactionPool.Logs {
		//	// Transaction Pool - Logs - Log Data
		//	require.Equal(t, l.TxHash, standardOb.TransactionPool.Logs[i].TxHash)
		//
		//	// Transaction Pool - Logs - Log data - Log
		//	require.Equal(t, l.Log, standardOb.TransactionPool.Logs[i].Log.Address)
		//
		//	for j, e := range ob.TransactionPool.Logs[i].Log.Events {
		//		require.Equal(t, e.Address, standardOb.TransactionPool.Logs[i].Log.Events[j].Address)
		//		require.Equal(t, e.Identifier, standardOb.TransactionPool.Logs[i].Log.Events[j].Identifier)
		//		require.Equal(t, e.Topics, standardOb.TransactionPool.Logs[i].Log.Events[j].Topics)
		//		require.Equal(t, e.Data, standardOb.TransactionPool.Logs[i].Log.Events[j].Data)
		//		require.Equal(t, e.AdditionalData, standardOb.TransactionPool.Logs[i].Log.Events[j].AdditionalData)
		//	}
		//}
		//
		//// Transaction Pool - ScheduledExecutedSCRSHashesPrevBlock
		//for i, s := range ob.TransactionPool.ScheduledExecutedSCRSHashesPrevBlock {
		//	require.Equal(t, s, standardOb.TransactionPool.ScheduledExecutedSCRSHashesPrevBlock[i])
		//}
		//
		//// Transaction Pool - ScheduledExecutedInvalidTxsHashesPrevBlock
		//for i, s := range ob.TransactionPool.ScheduledExecutedInvalidTxsHashesPrevBlock {
		//	require.Equal(t, s, standardOb.TransactionPool.ScheduledExecutedInvalidTxsHashesPrevBlock[i])
		//}
	}

	return nil
}

func castBigInt(i *big.Int) ([]byte, error) {
	buf := make([]byte, caster.Size(i))

	_, err := caster.MarshalTo(i, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to cast big Int: %w", err)
	}

	return buf, nil
}
