package process

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/stretchr/testify/require"

	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
)

const outportBlockJSONPath = "../testscommon/testdata/outportBlock.json"

var caster coreData.BigIntCaster

func TestOutportBlockAdapter_MarshalTo(t *testing.T) {
	jsonBytes, err := os.ReadFile(outportBlockJSONPath)
	require.NoError(t, err, "failed to read test data")

	ob := outport.OutportBlock{}
	err = json.Unmarshal(jsonBytes, &ob)
	require.NoError(t, err, "failed to unmarshal test block")

	adapter := NewOutportBlockCaster(&ob)
	standardOb := data.OutportBlock{}
	err = adapter.MarshalTo(&standardOb)
	require.NoError(t, err, "failed to marshal to standard outport")

	header := block.Header{}
	err = protoMarshaller.Unmarshal(&header, ob.BlockData.HeaderBytes)
	require.NoError(t, err, "failed to unmarshall outport block header bytes")

	// Asserting values.
	require.Equal(t, ob.ShardID, standardOb.ShardID)
	require.Equal(t, ob.NotarizedHeadersHashes, standardOb.NotarizedHeadersHashes)
	require.Equal(t, ob.NumberOfShards, standardOb.NumberOfShards)
	require.Equal(t, ob.SignersIndexes, standardOb.SignersIndexes)
	require.Equal(t, ob.HighestFinalBlockNonce, standardOb.HighestFinalBlockNonce)
	require.Equal(t, ob.HighestFinalBlockHash, standardOb.HighestFinalBlockHash)

	// Block data.
	require.Equal(t, ob.BlockData.ShardID, standardOb.BlockData.ShardID)
	require.Equal(t, ob.BlockData.HeaderHash, standardOb.BlockData.HeaderHash)
	require.Equal(t, ob.BlockData.HeaderType, standardOb.BlockData.HeaderType)

	// Block data - Header.
	require.Equal(t, header.Nonce, standardOb.BlockData.Header.Nonce)
	require.Equal(t, header.PrevHash, standardOb.BlockData.Header.PrevHash)
	require.Equal(t, header.PrevRandSeed, standardOb.BlockData.Header.PrevRandSeed)
	require.Equal(t, header.RandSeed, standardOb.BlockData.Header.RandSeed)
	require.Equal(t, header.PubKeysBitmap, standardOb.BlockData.Header.PubKeysBitmap)
	require.Equal(t, header.ShardID, standardOb.BlockData.Header.ShardID)
	require.Equal(t, header.TimeStamp, standardOb.BlockData.Header.TimeStamp)
	require.Equal(t, header.Round, standardOb.BlockData.Header.Round)
	require.Equal(t, header.Epoch, standardOb.BlockData.Header.Epoch)
	require.Equal(t, header.BlockBodyType.String(), standardOb.BlockData.Header.BlockBodyType.String())
	require.Equal(t, header.Signature, standardOb.BlockData.Header.Signature)
	require.Equal(t, header.LeaderSignature, standardOb.BlockData.Header.LeaderSignature)
	require.Equal(t, header.RootHash, standardOb.BlockData.Header.RootHash)
	require.Equal(t, header.MetaBlockHashes, standardOb.BlockData.Header.MetaBlockHashes)
	require.Equal(t, header.TxCount, standardOb.BlockData.Header.TxCount)
	require.Equal(t, header.EpochStartMetaHash, standardOb.BlockData.Header.EpochStartMetaHash)
	require.Equal(t, header.ReceiptsHash, standardOb.BlockData.Header.ReceiptsHash)
	require.Equal(t, header.ChainID, standardOb.BlockData.Header.ChainID)
	require.Equal(t, header.SoftwareVersion, standardOb.BlockData.Header.SoftwareVersion)
	require.Equal(t, header.Reserved, standardOb.BlockData.Header.Reserved)
	require.Equal(t, castBigInt(t, header.AccumulatedFees), standardOb.BlockData.Header.AccumulatedFees)
	require.Equal(t, castBigInt(t, header.DeveloperFees), standardOb.BlockData.Header.DeveloperFees)

	// Block data - Header - Mini block headers.
	for i, miniBlockHeader := range header.MiniBlockHeaders {
		require.Equal(t, miniBlockHeader.Hash, standardOb.BlockData.Header.MiniBlockHeaders[i].Hash)
		require.Equal(t, miniBlockHeader.SenderShardID, standardOb.BlockData.Header.MiniBlockHeaders[i].SenderShardID)
		require.Equal(t, miniBlockHeader.ReceiverShardID, standardOb.BlockData.Header.MiniBlockHeaders[i].ReceiverShardID)
		require.Equal(t, miniBlockHeader.TxCount, standardOb.BlockData.Header.MiniBlockHeaders[i].TxCount)
		require.Equal(t, miniBlockHeader.Type, standardOb.BlockData.Header.MiniBlockHeaders[i].Type)
		require.Equal(t, miniBlockHeader.Reserved, standardOb.BlockData.Header.MiniBlockHeaders[i].Reserved)
	}

	// Block data - Header - Peer changes.
	for i, peerChange := range header.PeerChanges {
		require.Equal(t, peerChange.PubKey, standardOb.BlockData.Header.PeerChanges[i].PubKey)
		require.Equal(t, peerChange.ShardIdDest, standardOb.BlockData.Header.PeerChanges[i].ShardIdDest)

	}

	for i, intraShardMiniBlock := range ob.BlockData.IntraShardMiniBlocks {
		require.Equal(t, intraShardMiniBlock.TxHashes, standardOb.BlockData.IntraShardMiniBlocks[i].TxHashes)
		require.Equal(t, intraShardMiniBlock.ReceiverShardID, standardOb.BlockData.IntraShardMiniBlocks[i].TxHashes)
		require.Equal(t, intraShardMiniBlock.SenderShardID, standardOb.BlockData.IntraShardMiniBlocks[i].SenderShardID)
		require.Equal(t, intraShardMiniBlock.Type.String(), standardOb.BlockData.IntraShardMiniBlocks[i].Type.String())
		require.Equal(t, intraShardMiniBlock.Reserved, standardOb.BlockData.IntraShardMiniBlocks[i].Reserved)
	}

	// Transaction pool - Transactions.
	for k, v := range ob.TransactionPool.Transactions {
		// Transaction pool - Transactions. - TxInfo.
		require.Equal(t, v.Transaction.Nonce, standardOb.TransactionPool.Transactions[k].Transaction.Nonce)
		require.Equal(t, castBigInt(t, v.Transaction.Value), standardOb.TransactionPool.Transactions[k].Transaction.Value)
		require.Equal(t, v.Transaction.RcvAddr, standardOb.TransactionPool.Transactions[k].Transaction.RcvAddr)
		require.Equal(t, v.Transaction.RcvUserName, standardOb.TransactionPool.Transactions[k].Transaction.RcvUserName)
		require.Equal(t, v.Transaction.SndAddr, standardOb.TransactionPool.Transactions[k].Transaction.SndAddr)
		require.Equal(t, v.Transaction.GasPrice, standardOb.TransactionPool.Transactions[k].Transaction.GasPrice)
		require.Equal(t, v.Transaction.GasLimit, standardOb.TransactionPool.Transactions[k].Transaction.GasLimit)
		require.Equal(t, v.Transaction.Data, standardOb.TransactionPool.Transactions[k].Transaction.Data)
		require.Equal(t, v.Transaction.ChainID, standardOb.TransactionPool.Transactions[k].Transaction.ChainID)
		require.Equal(t, v.Transaction.Version, standardOb.TransactionPool.Transactions[k].Transaction.Version)
		require.Equal(t, v.Transaction.Signature, standardOb.TransactionPool.Transactions[k].Transaction.Signature)
		require.Equal(t, v.Transaction.Options, standardOb.TransactionPool.Transactions[k].Transaction.Options)
		require.Equal(t, v.Transaction.GuardianAddr, standardOb.TransactionPool.Transactions[k].Transaction.GuardianAddr)
		require.Equal(t, v.Transaction.GuardianSignature, standardOb.TransactionPool.Transactions[k].Transaction.GuardianSignature)

		// Transaction pool - Transactions - Tx Info - Fee info.
		require.Equal(t, v.FeeInfo.GasUsed, standardOb.TransactionPool.Transactions[k].FeeInfo.GasUsed)
		require.Equal(t, castBigInt(t, v.FeeInfo.Fee), standardOb.TransactionPool.Transactions[k].FeeInfo.Fee)
		require.Equal(t, castBigInt(t, v.FeeInfo.InitialPaidFee), standardOb.TransactionPool.Transactions[k].FeeInfo.InitialPaidFee)

		require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.Transactions[k].ExecutionOrder)
	}

	// Transaction pool - Smart Contract results.
	for k, v := range ob.TransactionPool.SmartContractResults {
		// Transaction pool - Smart Contract results - SmartContractResult.
		require.Equal(t, v.SmartContractResult.Nonce, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.Nonce)
		require.Equal(t, castBigInt(t, v.SmartContractResult.Value), standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.Value)
		require.Equal(t, v.SmartContractResult.RcvAddr, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.RcvAddr)
		require.Equal(t, v.SmartContractResult.SndAddr, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.SndAddr)
		require.Equal(t, v.SmartContractResult.RelayerAddr, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.RelayerAddr)
		require.Equal(t, castBigInt(t, v.SmartContractResult.RelayedValue), standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.RelayedValue)
		require.Equal(t, v.SmartContractResult.Code, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.Code)
		require.Equal(t, v.SmartContractResult.Data, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.Data)
		require.Equal(t, v.SmartContractResult.PrevTxHash, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.PrevTxHash)
		require.Equal(t, v.SmartContractResult.OriginalTxHash, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.OriginalTxHash)
		require.Equal(t, v.SmartContractResult.GasLimit, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.GasLimit)
		require.Equal(t, v.SmartContractResult.GasPrice, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.GasPrice)
		require.Equal(t, int(v.SmartContractResult.CallType), standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.CallType)
		require.Equal(t, v.SmartContractResult.CodeMetadata, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.CodeMetadata)
		require.Equal(t, v.SmartContractResult.ReturnMessage, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.ReturnMessage)
		require.Equal(t, v.SmartContractResult.OriginalSender, standardOb.TransactionPool.SmartContractResults[k].SmartContractResult.OriginalSender)

		// Transaction pool - Smart Contract results - Fee info.
		require.Equal(t, v.FeeInfo.GasUsed, standardOb.TransactionPool.SmartContractResults[k].FeeInfo.GasUsed)
		require.Equal(t, castBigInt(t, v.FeeInfo.Fee), standardOb.TransactionPool.SmartContractResults[k].FeeInfo.Fee)
		require.Equal(t, castBigInt(t, v.FeeInfo.InitialPaidFee), standardOb.TransactionPool.SmartContractResults[k].FeeInfo.InitialPaidFee)

		// Transaction pool - Smart Contract results - Execution Order.
		require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.SmartContractResults[k].ExecutionOrder)
	}

	// Transaction Pool - Rewards
	for k, v := range ob.TransactionPool.Rewards {
		// Transaction Pool - Rewards - Reward info
		require.Equal(t, v.Reward.Round, standardOb.TransactionPool.Rewards[k].Reward.Round)
		require.Equal(t, castBigInt(t, v.Reward.Value), standardOb.TransactionPool.Rewards[k].Reward.Value)
		require.Equal(t, v.Reward.RcvAddr, standardOb.TransactionPool.Rewards[k].Reward.RcvAddr)
		require.Equal(t, v.Reward.Epoch, standardOb.TransactionPool.Rewards[k].Reward.Epoch)

		// Transaction Pool - Rewards - Execution Order
		require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.Rewards[k].ExecutionOrder)
	}

	// Transaction Pool - Receipts
	for k, v := range ob.TransactionPool.Receipts {
		// Transaction Pool - Receipts - Receipt info
		require.Equal(t, castBigInt(t, v.Value), standardOb.TransactionPool.Receipts[k].Value)
		require.Equal(t, v.SndAddr, standardOb.TransactionPool.Receipts[k].SndAddr)
		require.Equal(t, v.Data, standardOb.TransactionPool.Receipts[k].Data)
		require.Equal(t, v.TxHash, standardOb.TransactionPool.Receipts[k].TxHash)
	}

	// Transaction Pool - Invalid Txs
	for k, v := range ob.TransactionPool.InvalidTxs {
		// Transaction Pool - Invalid Txs - Tx Info
		require.Equal(t, v.Transaction.Nonce, standardOb.TransactionPool.InvalidTxs[k].Transaction.Nonce)
		require.Equal(t, castBigInt(t, v.Transaction.Value), standardOb.TransactionPool.InvalidTxs[k].Transaction.Value)
		require.Equal(t, v.Transaction.RcvAddr, standardOb.TransactionPool.InvalidTxs[k].Transaction.RcvAddr)
		require.Equal(t, v.Transaction.RcvUserName, standardOb.TransactionPool.InvalidTxs[k].Transaction.RcvUserName)
		require.Equal(t, v.Transaction.SndAddr, standardOb.TransactionPool.InvalidTxs[k].Transaction.SndAddr)
		require.Equal(t, v.Transaction.GasPrice, standardOb.TransactionPool.InvalidTxs[k].Transaction.GasPrice)
		require.Equal(t, v.Transaction.GasLimit, standardOb.TransactionPool.InvalidTxs[k].Transaction.GasLimit)
		require.Equal(t, v.Transaction.Data, standardOb.TransactionPool.InvalidTxs[k].Transaction.Data)
		require.Equal(t, v.Transaction.ChainID, standardOb.TransactionPool.InvalidTxs[k].Transaction.ChainID)
		require.Equal(t, v.Transaction.Version, standardOb.TransactionPool.InvalidTxs[k].Transaction.Version)
		require.Equal(t, v.Transaction.Signature, standardOb.TransactionPool.InvalidTxs[k].Transaction.Signature)
		require.Equal(t, v.Transaction.Options, standardOb.TransactionPool.InvalidTxs[k].Transaction.Options)
		require.Equal(t, v.Transaction.GuardianAddr, standardOb.TransactionPool.InvalidTxs[k].Transaction.GuardianAddr)
		require.Equal(t, v.Transaction.GuardianSignature, standardOb.TransactionPool.InvalidTxs[k].Transaction.GuardianSignature)
		require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.InvalidTxs[k].ExecutionOrder)

		// Transaction pool - Invalid Txs - Fee info.
		require.Equal(t, v.FeeInfo.GasUsed, standardOb.TransactionPool.InvalidTxs[k].FeeInfo.GasUsed)
		require.Equal(t, castBigInt(t, v.FeeInfo.Fee), standardOb.TransactionPool.InvalidTxs[k].FeeInfo.Fee)
		require.Equal(t, castBigInt(t, v.FeeInfo.InitialPaidFee), standardOb.TransactionPool.InvalidTxs[k].FeeInfo.InitialPaidFee)

		require.Equal(t, v.ExecutionOrder, standardOb.TransactionPool.InvalidTxs[k].ExecutionOrder)
	}

	// Transaction Pool - Logs
	for i, l := range ob.TransactionPool.Logs {
		// Transaction Pool - Logs - Log Data
		require.Equal(t, l.TxHash, standardOb.TransactionPool.Logs[i].TxHash)

		// Transaction Pool - Logs - Log data - Log
		require.Equal(t, l.Log, standardOb.TransactionPool.Logs[i].Log.Address)

		for j, e := range ob.TransactionPool.Logs[i].Log.Events {
			require.Equal(t, e.Address, standardOb.TransactionPool.Logs[i].Log.Events[j].Address)
			require.Equal(t, e.Identifier, standardOb.TransactionPool.Logs[i].Log.Events[j].Identifier)
			require.Equal(t, e.Topics, standardOb.TransactionPool.Logs[i].Log.Events[j].Topics)
			require.Equal(t, e.Data, standardOb.TransactionPool.Logs[i].Log.Events[j].Data)
			require.Equal(t, e.AdditionalData, standardOb.TransactionPool.Logs[i].Log.Events[j].AdditionalData)
		}
	}

	// Transaction Pool - ScheduledExecutedSCRSHashesPrevBlock
	for i, s := range ob.TransactionPool.ScheduledExecutedSCRSHashesPrevBlock {
		require.Equal(t, s, standardOb.TransactionPool.ScheduledExecutedSCRSHashesPrevBlock[i])
	}

	// Transaction Pool - ScheduledExecutedInvalidTxsHashesPrevBlock
	for i, s := range ob.TransactionPool.ScheduledExecutedInvalidTxsHashesPrevBlock {
		require.Equal(t, s, standardOb.TransactionPool.ScheduledExecutedInvalidTxsHashesPrevBlock[i])
	}

	// Header gas consumption.
	require.Equal(t, ob.HeaderGasConsumption.GasProvided, standardOb.HeaderGasConsumption.GasProvided)
	require.Equal(t, ob.HeaderGasConsumption.GasRefunded, standardOb.HeaderGasConsumption.GasRefunded)
	require.Equal(t, ob.HeaderGasConsumption.GasPenalized, standardOb.HeaderGasConsumption.GasPenalized)
	require.Equal(t, ob.HeaderGasConsumption.MaxGasPerBlock, standardOb.HeaderGasConsumption.MaxGasPerBlock)

	// Altered accounts.
	for key, account := range ob.AlteredAccounts {
		acc := standardOb.AlteredAccounts[key]

		require.Equal(t, account.Address, acc.Address)
		require.Equal(t, account.Nonce, acc.Nonce)
		require.Equal(t, account.Balance, acc.Balance)
		require.Equal(t, account.Balance, acc.Balance)

		// Altered accounts - Account token data.
		for i, token := range account.Tokens {
			require.Equal(t, token.Nonce, acc.Tokens[i].Nonce)
			require.Equal(t, token.Identifier, acc.Tokens[i].Identifier)
			require.Equal(t, token.Balance, acc.Tokens[i].Balance)
			require.Equal(t, token.Properties, acc.Tokens[i].Properties)

			// Altered accounts - Account token data - Metadata.
			require.Equal(t, token.MetaData.Nonce, acc.Tokens[i].MetaData.Nonce)
			require.Equal(t, token.MetaData.Name, acc.Tokens[i].MetaData.Name)
			require.Equal(t, token.MetaData.Creator, acc.Tokens[i].MetaData.Creator)
			require.Equal(t, token.MetaData.Royalties, acc.Tokens[i].MetaData.Royalties)
			require.Equal(t, token.MetaData.Hash, acc.Tokens[i].MetaData.Hash)
			require.Equal(t, token.MetaData.URIs, acc.Tokens[i].MetaData.URIs)
			require.Equal(t, token.MetaData.Attributes, acc.Tokens[i].MetaData.Attributes)

			// Altered accounts - Account token data - Additional data.
			require.Equal(t, token.AdditionalData.IsNFTCreate, acc.Tokens[i].AdditionalData.IsNFTCreate)
		}

		//  Altered accounts - Additional account data.
		require.Equal(t, account.AdditionalData.IsSender, acc.AdditionalData.IsSender)
		require.Equal(t, account.AdditionalData.BalanceChanged, acc.AdditionalData.BalanceChanged)
		require.Equal(t, account.AdditionalData.CurrentOwner, acc.AdditionalData.CurrentOwner)
		require.Equal(t, account.AdditionalData.UserName, acc.AdditionalData.UserName)
		require.Equal(t, account.AdditionalData.DeveloperRewards, acc.AdditionalData.DeveloperRewards)
		require.Equal(t, account.AdditionalData.CodeHash, acc.AdditionalData.CodeHash)
		require.Equal(t, account.AdditionalData.RootHash, acc.AdditionalData.RootHash)
		require.Equal(t, account.AdditionalData.CodeMetadata, acc.AdditionalData.CodeMetadata)
	}
}

func castBigInt(t *testing.T, i *big.Int) []byte {
	t.Helper()

	buf := make([]byte, caster.Size(i))

	_, err := caster.MarshalTo(i, buf)
	require.NoError(t, err, "failed to cast from big.Int")

	return buf
}
