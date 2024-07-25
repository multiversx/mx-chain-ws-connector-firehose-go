package process_test

import (
	"encoding/json"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/stretchr/testify/require"

	data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
)

type fieldsGetter interface {
	GetShardID() uint32
	GetHeaderGasConsumption() *data.HeaderGasConsumption
	GetAlteredAccounts() map[string]*data.AlteredAccount
	GetNotarizedHeadersHashes() []string
	GetNumberOfShards() uint32
	GetSignersIndexes() []uint64
	GetHighestFinalBlockNonce() uint64
	GetHighestFinalBlockHash() []byte
}

type fieldsGetterV1 interface {
	fieldsGetter
	GetTransactionPool() *data.TransactionPool
}

type fieldsGetterV2 interface {
	fieldsGetter
	GetTransactionPool() *data.TransactionPoolV2
}

type blockDataGetter interface {
	GetShardID() uint32
	GetHeaderType() string
	GetHeaderHash() []byte
	GetBody() *data.Body
	GetIntraShardMiniBlocks() []*data.MiniBlock
	GetScheduledRootHash() []byte
	GetScheduledAccumulatedFees() []byte
	GetScheduledDeveloperFees() []byte
	GetScheduledGasProvided() uint64
	GetScheduledGasPenalized() uint64
	GetScheduledGasRefunded() uint64
}

const (
	outportBlockHeaderV1JSONPath  = "../testscommon/testdata/outportBlockHeaderV1.json"
	outportBlockHeaderV2JSONPath  = "../testscommon/testdata/outportBlockHeaderV2.json"
	outportBlockMetaBlockJSONPath = "../testscommon/testdata/outportBlockMetaBlock.json"
)

func TestOutportBlockConverter(t *testing.T) {
	jsonBytes, err := os.ReadFile("../testscommon/testdata/outportBlocks.json")
	require.NoError(t, err, "failed to read test data")

	ob := &outport.OutportBlock{}
	err = json.Unmarshal(jsonBytes, ob)
	require.NoError(t, err, "failed to unmarshal test block")

	converter, err := process.NewOutportBlockConverter(gogoProtoMarshaller, protoMarshaller)
	require.Nil(t, err)

	shardOutportBlock, err := converter.HandleShardOutportBlockV2(ob)
	if err != nil {
		panic(err)
	}

	header := &block.HeaderV2{}
	err = gogoProtoMarshaller.Unmarshal(header, ob.BlockData.HeaderBytes)
	require.NoError(t, err, "failed to unmarshall outport block header bytes")

	checkHeaderV2ShardV2(t, header, shardOutportBlock)
	checkFieldsV2(t, ob, shardOutportBlock)
	checkBlockData(t, ob.BlockData, shardOutportBlock.BlockData)

	j, err := json.Marshal(shardOutportBlock)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./shardoutport.json", j, 0655)
	if err != nil {
		panic(err)
	}

	jsonBytes, err = os.ReadFile(outportBlockHeaderV2JSONPath)
	require.NoError(t, err, "failed to read test data")

}

func TestOutportBlockConverter_HandleShardOutportBlockV2(t *testing.T) {
	t.Parallel()

	jsonBytes, err := os.ReadFile(outportBlockHeaderV1JSONPath)
	require.NoError(t, err, "failed to read test data")

	ob := &outport.OutportBlock{}
	err = json.Unmarshal(jsonBytes, ob)
	require.NoError(t, err, "failed to unmarshal test block")

	converter, err := process.NewOutportBlockConverter(gogoProtoMarshaller, protoMarshaller)
	require.Nil(t, err)

	shardOutportBlock, err := converter.HandleShardOutportBlockV2(ob)
	if err != nil {
		panic(err)
	}

	header := &block.Header{}
	err = gogoProtoMarshaller.Unmarshal(header, ob.BlockData.HeaderBytes)
	require.NoError(t, err, "failed to unmarshall outport block header bytes")

	checkHeaderV1ShardV2(t, header, shardOutportBlock)
	checkFieldsV2(t, ob, shardOutportBlock)
	checkBlockData(t, ob.BlockData, shardOutportBlock.BlockData)

	j, err := json.Marshal(shardOutportBlock)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./shardoutport.json", j, 0655)
	if err != nil {
		panic(err)
	}

	jsonBytes, err = os.ReadFile(outportBlockHeaderV2JSONPath)
	require.NoError(t, err, "failed to read test data")

	ob = &outport.OutportBlock{}
	err = json.Unmarshal(jsonBytes, ob)
	require.NoError(t, err, "failed to unmarshal test block")

	converter, err = process.NewOutportBlockConverter(gogoProtoMarshaller, protoMarshaller)
	require.Nil(t, err)

	shardOutportBlock, err = converter.HandleShardOutportBlockV2(ob)
	if err != nil {
		panic(err)
	}

	headerV2 := &block.HeaderV2{}
	err = gogoProtoMarshaller.Unmarshal(headerV2, ob.BlockData.HeaderBytes)
	require.NoError(t, err, "failed to unmarshall outport block header bytes")

	checkHeaderV2ShardV2(t, headerV2, shardOutportBlock)
	checkFieldsV2(t, ob, shardOutportBlock)
	checkBlockData(t, ob.BlockData, shardOutportBlock.BlockData)

	j, err = json.Marshal(shardOutportBlock)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./shardoutpor2.json", j, 0655)
	if err != nil {
		panic(err)
	}

}

func TestHeaderConverter(t *testing.T) {
	t.Parallel()

	jsonBytes, err := os.ReadFile(outportBlockHeaderV1JSONPath)
	require.NoError(t, err, "failed to read test data")

	ob := &outport.OutportBlock{}
	err = json.Unmarshal(jsonBytes, ob)
	require.NoError(t, err, "failed to unmarshal test block")

	converter, err := process.NewOutportBlockConverter(gogoProtoMarshaller, protoMarshaller)
	require.Nil(t, err)
	shardOutportBlock, err := converter.HandleShardOutportBlock(ob)
	require.NoError(t, err, "failed to marshal to standard outport")

	header := &block.Header{}
	err = gogoProtoMarshaller.Unmarshal(header, ob.BlockData.HeaderBytes)
	require.NoError(t, err, "failed to unmarshall outport block header bytes")

	checkHeaderV1(t, header, shardOutportBlock)
	checkFieldsV1(t, ob, shardOutportBlock)
	checkBlockData(t, ob.BlockData, shardOutportBlock.BlockData)
}

func TestHeaderV2Converter(t *testing.T) {
	t.Parallel()

	jsonBytes, err := os.ReadFile(outportBlockHeaderV2JSONPath)
	require.NoError(t, err, "failed to read test data")

	ob := outport.OutportBlock{}
	err = json.Unmarshal(jsonBytes, &ob)
	require.NoError(t, err, "failed to unmarshal test block")

	converter, err := process.NewOutportBlockConverter(gogoProtoMarshaller, protoMarshaller)
	require.Nil(t, err)
	shardOutportBlock, err := converter.HandleShardOutportBlock(&ob)
	require.Nil(t, err)
	require.NoError(t, err, "failed to marshal to standard outport")

	header := block.HeaderV2{}
	err = gogoProtoMarshaller.Unmarshal(&header, ob.BlockData.HeaderBytes)
	require.NoError(t, err, "failed to unmarshall outport block header bytes")

	checkHeaderV2(t, &header, shardOutportBlock)
	checkFieldsV1(t, &ob, shardOutportBlock)
	checkBlockData(t, ob.BlockData, shardOutportBlock.BlockData)
}

func TestMetaBlockConverter(t *testing.T) {
	t.Parallel()

	jsonBytes, err := os.ReadFile(outportBlockMetaBlockJSONPath)
	require.NoError(t, err, "failed to read test data")

	ob := outport.OutportBlock{}
	err = json.Unmarshal(jsonBytes, &ob)
	require.NoError(t, err, "failed to unmarshal test block")

	converter, _ := process.NewOutportBlockConverter(gogoProtoMarshaller, protoMarshaller)
	metaOutportBlock, err := converter.HandleMetaOutportBlock(&ob)
	require.NoError(t, err, "failed to marshal to standard outport")

	header := block.MetaBlock{}
	err = gogoProtoMarshaller.Unmarshal(&header, ob.BlockData.HeaderBytes)
	require.NoError(t, err, "failed to unmarshall outport block header bytes")

	checkHeaderMeta(t, &header, metaOutportBlock)
	checkFieldsV1(t, &ob, metaOutportBlock)
	checkBlockData(t, ob.BlockData, metaOutportBlock.BlockData)
}

func checkHeaderV1ShardV2(t *testing.T, header *block.Header, shardOutportBlock *data.ShardOutportBlockV2) {
	require.Equal(t, header.Nonce, shardOutportBlock.BlockData.Header.Nonce)
	require.Equal(t, header.PrevHash, shardOutportBlock.BlockData.Header.PrevHash)
	require.Equal(t, header.PrevRandSeed, shardOutportBlock.BlockData.Header.PrevRandSeed)
	require.Equal(t, header.RandSeed, shardOutportBlock.BlockData.Header.RandSeed)
	require.Equal(t, header.PubKeysBitmap, shardOutportBlock.BlockData.Header.PubKeysBitmap)
	require.Equal(t, header.ShardID, shardOutportBlock.BlockData.Header.ShardID)
	require.Equal(t, header.TimeStamp, shardOutportBlock.BlockData.Header.TimeStamp)
	require.Equal(t, header.Round, shardOutportBlock.BlockData.Header.Round)
	require.Equal(t, header.Epoch, shardOutportBlock.BlockData.Header.Epoch)
	require.Equal(t, header.BlockBodyType.String(), shardOutportBlock.BlockData.Header.BlockBodyType.String())
	require.Equal(t, header.Signature, shardOutportBlock.BlockData.Header.Signature)
	require.Equal(t, header.LeaderSignature, shardOutportBlock.BlockData.Header.LeaderSignature)
	require.Equal(t, header.RootHash, shardOutportBlock.BlockData.Header.RootHash)
	require.Equal(t, header.MetaBlockHashes, shardOutportBlock.BlockData.Header.MetaBlockHashes)
	require.Equal(t, header.TxCount, shardOutportBlock.BlockData.Header.TxCount)
	require.Equal(t, header.EpochStartMetaHash, shardOutportBlock.BlockData.Header.EpochStartMetaHash)
	require.Equal(t, header.ReceiptsHash, shardOutportBlock.BlockData.Header.ReceiptsHash)
	require.Equal(t, header.ChainID, shardOutportBlock.BlockData.Header.ChainID)
	require.Equal(t, header.SoftwareVersion, shardOutportBlock.BlockData.Header.SoftwareVersion)
	require.Equal(t, header.Reserved, shardOutportBlock.BlockData.Header.Reserved)
	require.Equal(t, mustCastBigInt(t, header.AccumulatedFees), shardOutportBlock.BlockData.Header.AccumulatedFees)
	require.Equal(t, mustCastBigInt(t, header.DeveloperFees), shardOutportBlock.BlockData.Header.DeveloperFees)

	// Block data - Header - Mini block headers.
	for i, miniBlockHeader := range header.MiniBlockHeaders {
		require.Equal(t, miniBlockHeader.Hash, shardOutportBlock.BlockData.Header.MiniBlockHeaders[i].Hash)
		require.Equal(t, miniBlockHeader.SenderShardID, shardOutportBlock.BlockData.Header.MiniBlockHeaders[i].SenderShardID)
		require.Equal(t, miniBlockHeader.ReceiverShardID, shardOutportBlock.BlockData.Header.MiniBlockHeaders[i].ReceiverShardID)
		require.Equal(t, miniBlockHeader.TxCount, shardOutportBlock.BlockData.Header.MiniBlockHeaders[i].TxCount)
		require.Equal(t, miniBlockHeader.Type, shardOutportBlock.BlockData.Header.MiniBlockHeaders[i].Type)
		require.Equal(t, miniBlockHeader.Reserved, shardOutportBlock.BlockData.Header.MiniBlockHeaders[i].Reserved)
	}

	// Block data - Header - Peer changes.
	for i, peerChange := range header.PeerChanges {
		require.Equal(t, peerChange.PubKey, shardOutportBlock.BlockData.Header.PeerChanges[i].PubKey)
		require.Equal(t, peerChange.ShardIdDest, shardOutportBlock.BlockData.Header.PeerChanges[i].ShardIdDest)
	}
}

func checkHeaderV2ShardV2(t *testing.T, headerV2 *block.HeaderV2, shardOutportBlock *data.ShardOutportBlockV2) {
	// Block data - Header.
	header := headerV2.Header
	require.Equal(t, header.Nonce, shardOutportBlock.BlockData.Header.Nonce)
	require.Equal(t, header.PrevHash, shardOutportBlock.BlockData.Header.PrevHash)
	require.Equal(t, header.PrevRandSeed, shardOutportBlock.BlockData.Header.PrevRandSeed)
	require.Equal(t, header.RandSeed, shardOutportBlock.BlockData.Header.RandSeed)
	require.Equal(t, header.PubKeysBitmap, shardOutportBlock.BlockData.Header.PubKeysBitmap)
	require.Equal(t, header.ShardID, shardOutportBlock.BlockData.Header.ShardID)
	require.Equal(t, header.TimeStamp, shardOutportBlock.BlockData.Header.TimeStamp)
	require.Equal(t, header.Round, shardOutportBlock.BlockData.Header.Round)
	require.Equal(t, header.Epoch, shardOutportBlock.BlockData.Header.Epoch)
	require.Equal(t, header.BlockBodyType.String(), shardOutportBlock.BlockData.Header.BlockBodyType.String())
	require.Equal(t, header.Signature, shardOutportBlock.BlockData.Header.Signature)
	require.Equal(t, header.LeaderSignature, shardOutportBlock.BlockData.Header.LeaderSignature)
	require.Equal(t, header.RootHash, shardOutportBlock.BlockData.Header.RootHash)
	require.Equal(t, header.MetaBlockHashes, shardOutportBlock.BlockData.Header.MetaBlockHashes)
	require.Equal(t, header.TxCount, shardOutportBlock.BlockData.Header.TxCount)
	require.Equal(t, header.EpochStartMetaHash, shardOutportBlock.BlockData.Header.EpochStartMetaHash)
	require.Equal(t, header.ReceiptsHash, shardOutportBlock.BlockData.Header.ReceiptsHash)
	require.Equal(t, header.ChainID, shardOutportBlock.BlockData.Header.ChainID)
	require.Equal(t, header.SoftwareVersion, shardOutportBlock.BlockData.Header.SoftwareVersion)
	require.Equal(t, header.Reserved, shardOutportBlock.BlockData.Header.Reserved)
	require.Equal(t, mustCastBigInt(t, header.AccumulatedFees), shardOutportBlock.BlockData.Header.AccumulatedFees)
	require.Equal(t, mustCastBigInt(t, header.DeveloperFees), shardOutportBlock.BlockData.Header.DeveloperFees)

	// Block data - Header - Mini block headers.
	for i, miniBlockHeader := range header.MiniBlockHeaders {
		require.Equal(t, miniBlockHeader.Hash, shardOutportBlock.BlockData.Header.MiniBlockHeaders[i].Hash)
		require.Equal(t, miniBlockHeader.SenderShardID, shardOutportBlock.BlockData.Header.MiniBlockHeaders[i].SenderShardID)
		require.Equal(t, miniBlockHeader.ReceiverShardID, shardOutportBlock.BlockData.Header.MiniBlockHeaders[i].ReceiverShardID)
		require.Equal(t, miniBlockHeader.TxCount, shardOutportBlock.BlockData.Header.MiniBlockHeaders[i].TxCount)
		require.Equal(t, miniBlockHeader.Type.String(), shardOutportBlock.BlockData.Header.MiniBlockHeaders[i].Type.String())
		require.Equal(t, miniBlockHeader.Reserved, shardOutportBlock.BlockData.Header.MiniBlockHeaders[i].Reserved)
	}

	// Block data - Header - Peer changes.
	for i, peerChange := range header.PeerChanges {
		require.Equal(t, peerChange.PubKey, shardOutportBlock.BlockData.Header.PeerChanges[i].PubKey)
		require.Equal(t, peerChange.ShardIdDest, shardOutportBlock.BlockData.Header.PeerChanges[i].ShardIdDest)
	}

	require.Equal(t, headerV2.ScheduledRootHash, shardOutportBlock.BlockData.ScheduledRootHash)
	require.Equal(t, mustCastBigInt(t, headerV2.ScheduledAccumulatedFees), shardOutportBlock.BlockData.ScheduledAccumulatedFees)
	require.Equal(t, mustCastBigInt(t, headerV2.ScheduledDeveloperFees), shardOutportBlock.BlockData.ScheduledDeveloperFees)
	require.Equal(t, headerV2.ScheduledGasProvided, shardOutportBlock.BlockData.ScheduledGasProvided)
	require.Equal(t, headerV2.ScheduledGasPenalized, shardOutportBlock.BlockData.ScheduledGasPenalized)
	require.Equal(t, headerV2.ScheduledGasRefunded, shardOutportBlock.BlockData.ScheduledGasRefunded)
}

func checkHeaderV1(t *testing.T, header *block.Header, fireOutportBlock *data.ShardOutportBlock) {
	// Block data - Header.
	require.Equal(t, header.Nonce, fireOutportBlock.BlockData.Header.Nonce)
	require.Equal(t, header.PrevHash, fireOutportBlock.BlockData.Header.PrevHash)
	require.Equal(t, header.PrevRandSeed, fireOutportBlock.BlockData.Header.PrevRandSeed)
	require.Equal(t, header.RandSeed, fireOutportBlock.BlockData.Header.RandSeed)
	require.Equal(t, header.PubKeysBitmap, fireOutportBlock.BlockData.Header.PubKeysBitmap)
	require.Equal(t, header.ShardID, fireOutportBlock.BlockData.Header.ShardID)
	require.Equal(t, header.TimeStamp, fireOutportBlock.BlockData.Header.TimeStamp)
	require.Equal(t, header.Round, fireOutportBlock.BlockData.Header.Round)
	require.Equal(t, header.Epoch, fireOutportBlock.BlockData.Header.Epoch)
	require.Equal(t, header.BlockBodyType.String(), fireOutportBlock.BlockData.Header.BlockBodyType.String())
	require.Equal(t, header.Signature, fireOutportBlock.BlockData.Header.Signature)
	require.Equal(t, header.LeaderSignature, fireOutportBlock.BlockData.Header.LeaderSignature)
	require.Equal(t, header.RootHash, fireOutportBlock.BlockData.Header.RootHash)
	require.Equal(t, header.MetaBlockHashes, fireOutportBlock.BlockData.Header.MetaBlockHashes)
	require.Equal(t, header.TxCount, fireOutportBlock.BlockData.Header.TxCount)
	require.Equal(t, header.EpochStartMetaHash, fireOutportBlock.BlockData.Header.EpochStartMetaHash)
	require.Equal(t, header.ReceiptsHash, fireOutportBlock.BlockData.Header.ReceiptsHash)
	require.Equal(t, header.ChainID, fireOutportBlock.BlockData.Header.ChainID)
	require.Equal(t, header.SoftwareVersion, fireOutportBlock.BlockData.Header.SoftwareVersion)
	require.Equal(t, header.Reserved, fireOutportBlock.BlockData.Header.Reserved)
	require.Equal(t, mustCastBigInt(t, header.AccumulatedFees), fireOutportBlock.BlockData.Header.AccumulatedFees)
	require.Equal(t, mustCastBigInt(t, header.DeveloperFees), fireOutportBlock.BlockData.Header.DeveloperFees)

	// Block data - Header - Mini block headers.
	for i, miniBlockHeader := range header.MiniBlockHeaders {
		require.Equal(t, miniBlockHeader.Hash, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].Hash)
		require.Equal(t, miniBlockHeader.SenderShardID, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].SenderShardID)
		require.Equal(t, miniBlockHeader.ReceiverShardID, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].ReceiverShardID)
		require.Equal(t, miniBlockHeader.TxCount, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].TxCount)
		require.Equal(t, miniBlockHeader.Type, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].Type)
		require.Equal(t, miniBlockHeader.Reserved, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].Reserved)
	}

	// Block data - Header - Peer changes.
	for i, peerChange := range header.PeerChanges {
		require.Equal(t, peerChange.PubKey, fireOutportBlock.BlockData.Header.PeerChanges[i].PubKey)
		require.Equal(t, peerChange.ShardIdDest, fireOutportBlock.BlockData.Header.PeerChanges[i].ShardIdDest)
	}
}

func checkHeaderV2(t *testing.T, headerV2 *block.HeaderV2, fireOutportBlock *data.ShardOutportBlock) {
	// Block data - Header.
	header := headerV2.Header
	require.Equal(t, header.Nonce, fireOutportBlock.BlockData.Header.Nonce)
	require.Equal(t, header.PrevHash, fireOutportBlock.BlockData.Header.PrevHash)
	require.Equal(t, header.PrevRandSeed, fireOutportBlock.BlockData.Header.PrevRandSeed)
	require.Equal(t, header.RandSeed, fireOutportBlock.BlockData.Header.RandSeed)
	require.Equal(t, header.PubKeysBitmap, fireOutportBlock.BlockData.Header.PubKeysBitmap)
	require.Equal(t, header.ShardID, fireOutportBlock.BlockData.Header.ShardID)
	require.Equal(t, header.TimeStamp, fireOutportBlock.BlockData.Header.TimeStamp)
	require.Equal(t, header.Round, fireOutportBlock.BlockData.Header.Round)
	require.Equal(t, header.Epoch, fireOutportBlock.BlockData.Header.Epoch)
	require.Equal(t, header.BlockBodyType.String(), fireOutportBlock.BlockData.Header.BlockBodyType.String())
	require.Equal(t, header.Signature, fireOutportBlock.BlockData.Header.Signature)
	require.Equal(t, header.LeaderSignature, fireOutportBlock.BlockData.Header.LeaderSignature)
	require.Equal(t, header.RootHash, fireOutportBlock.BlockData.Header.RootHash)
	require.Equal(t, header.MetaBlockHashes, fireOutportBlock.BlockData.Header.MetaBlockHashes)
	require.Equal(t, header.TxCount, fireOutportBlock.BlockData.Header.TxCount)
	require.Equal(t, header.EpochStartMetaHash, fireOutportBlock.BlockData.Header.EpochStartMetaHash)
	require.Equal(t, header.ReceiptsHash, fireOutportBlock.BlockData.Header.ReceiptsHash)
	require.Equal(t, header.ChainID, fireOutportBlock.BlockData.Header.ChainID)
	require.Equal(t, header.SoftwareVersion, fireOutportBlock.BlockData.Header.SoftwareVersion)
	require.Equal(t, header.Reserved, fireOutportBlock.BlockData.Header.Reserved)
	require.Equal(t, mustCastBigInt(t, header.AccumulatedFees), fireOutportBlock.BlockData.Header.AccumulatedFees)
	require.Equal(t, mustCastBigInt(t, header.DeveloperFees), fireOutportBlock.BlockData.Header.DeveloperFees)

	// Block data - Header - Mini block headers.
	for i, miniBlockHeader := range header.MiniBlockHeaders {
		require.Equal(t, miniBlockHeader.Hash, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].Hash)
		require.Equal(t, miniBlockHeader.SenderShardID, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].SenderShardID)
		require.Equal(t, miniBlockHeader.ReceiverShardID, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].ReceiverShardID)
		require.Equal(t, miniBlockHeader.TxCount, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].TxCount)
		require.Equal(t, miniBlockHeader.Type.String(), fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].Type.String())
		require.Equal(t, miniBlockHeader.Reserved, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].Reserved)
	}

	// Block data - Header - Peer changes.
	for i, peerChange := range header.PeerChanges {
		require.Equal(t, peerChange.PubKey, fireOutportBlock.BlockData.Header.PeerChanges[i].PubKey)
		require.Equal(t, peerChange.ShardIdDest, fireOutportBlock.BlockData.Header.PeerChanges[i].ShardIdDest)
	}

	require.Equal(t, headerV2.ScheduledRootHash, fireOutportBlock.BlockData.ScheduledRootHash)
	require.Equal(t, mustCastBigInt(t, headerV2.ScheduledAccumulatedFees), fireOutportBlock.BlockData.ScheduledAccumulatedFees)
	require.Equal(t, mustCastBigInt(t, headerV2.ScheduledDeveloperFees), fireOutportBlock.BlockData.ScheduledDeveloperFees)
	require.Equal(t, headerV2.ScheduledGasProvided, fireOutportBlock.BlockData.ScheduledGasProvided)
	require.Equal(t, headerV2.ScheduledGasPenalized, fireOutportBlock.BlockData.ScheduledGasPenalized)
	require.Equal(t, headerV2.ScheduledGasRefunded, fireOutportBlock.BlockData.ScheduledGasRefunded)
}

func checkHeaderMeta(t *testing.T, header *block.MetaBlock, fireOutportBlock *data.MetaOutportBlock) {
	require.Equal(t, header.Nonce, fireOutportBlock.BlockData.Header.Nonce)
	require.Equal(t, header.Epoch, fireOutportBlock.BlockData.Header.Epoch)
	require.Equal(t, header.Round, fireOutportBlock.BlockData.Header.Round)
	require.Equal(t, header.TimeStamp, fireOutportBlock.BlockData.Header.TimeStamp)

	for i, si := range header.ShardInfo {
		require.Equal(t, si.ShardID, fireOutportBlock.BlockData.Header.ShardInfo[i].ShardID)
		require.Equal(t, si.HeaderHash, fireOutportBlock.BlockData.Header.ShardInfo[i].HeaderHash)
		require.Equal(t, si.HeaderHash, fireOutportBlock.BlockData.Header.ShardInfo[i].HeaderHash)

		for j, mbh := range header.MiniBlockHeaders {
			require.Equal(t, mbh.Hash, fireOutportBlock.BlockData.Header.MiniBlockHeaders[j].Hash)
			require.Equal(t, mbh.SenderShardID, fireOutportBlock.BlockData.Header.MiniBlockHeaders[j].SenderShardID)
			require.Equal(t, mbh.ReceiverShardID, fireOutportBlock.BlockData.Header.MiniBlockHeaders[j].ReceiverShardID)
			require.Equal(t, mbh.TxCount, fireOutportBlock.BlockData.Header.MiniBlockHeaders[j].TxCount)
			require.Equal(t, mbh.Type, fireOutportBlock.BlockData.Header.MiniBlockHeaders[j].Type)
			require.Equal(t, mbh.Reserved, fireOutportBlock.BlockData.Header.MiniBlockHeaders[j].Reserved)
		}

		require.Equal(t, si.PrevRandSeed, fireOutportBlock.BlockData.Header.ShardInfo[i].PrevRandSeed)
		require.Equal(t, si.PubKeysBitmap, fireOutportBlock.BlockData.Header.ShardInfo[i].PubKeysBitmap)
		require.Equal(t, si.Signature, fireOutportBlock.BlockData.Header.ShardInfo[i].Signature)
		require.Equal(t, si.Round, fireOutportBlock.BlockData.Header.ShardInfo[i].Round)
		require.Equal(t, si.PrevHash, fireOutportBlock.BlockData.Header.ShardInfo[i].PrevHash)
		require.Equal(t, si.Nonce, fireOutportBlock.BlockData.Header.ShardInfo[i].Nonce)
		require.Equal(t, mustCastBigInt(t, si.AccumulatedFees), fireOutportBlock.BlockData.Header.ShardInfo[i].AccumulatedFees)
		require.Equal(t, mustCastBigInt(t, si.DeveloperFees), fireOutportBlock.BlockData.Header.ShardInfo[i].DeveloperFees)
		require.Equal(t, si.NumPendingMiniBlocks, fireOutportBlock.BlockData.Header.ShardInfo[i].NumPendingMiniBlocks)
		require.Equal(t, si.LastIncludedMetaNonce, fireOutportBlock.BlockData.Header.ShardInfo[i].LastIncludedMetaNonce)
		require.Equal(t, si.ShardID, fireOutportBlock.BlockData.Header.ShardInfo[i].ShardID)
		require.Equal(t, si.TxCount, fireOutportBlock.BlockData.Header.ShardInfo[i].TxCount)
	}

	for i, peerInfo := range header.PeerInfo {
		require.Equal(t, peerInfo.Address, fireOutportBlock.BlockData.Header.PeerInfo[i].Address)
		require.Equal(t, peerInfo.PublicKey, fireOutportBlock.BlockData.Header.PeerInfo[i].PublicKey)
		require.Equal(t, peerInfo.Action, fireOutportBlock.BlockData.Header.PeerInfo[i].Action)
		require.Equal(t, peerInfo.TimeStamp, fireOutportBlock.BlockData.Header.PeerInfo[i].TimeStamp)
		require.Equal(t, peerInfo.ValueChange, fireOutportBlock.BlockData.Header.PeerInfo[i].ValueChange)
	}

	require.Equal(t, header.Signature, fireOutportBlock.BlockData.Header.Signature)
	require.Equal(t, header.LeaderSignature, fireOutportBlock.BlockData.Header.LeaderSignature)
	require.Equal(t, header.PubKeysBitmap, fireOutportBlock.BlockData.Header.PubKeysBitmap)
	require.Equal(t, header.PrevHash, fireOutportBlock.BlockData.Header.PrevHash)
	require.Equal(t, header.PrevRandSeed, fireOutportBlock.BlockData.Header.PrevRandSeed)
	require.Equal(t, header.RandSeed, fireOutportBlock.BlockData.Header.RandSeed)
	require.Equal(t, header.RootHash, fireOutportBlock.BlockData.Header.RootHash)
	require.Equal(t, header.ValidatorStatsRootHash, fireOutportBlock.BlockData.Header.ValidatorStatsRootHash)

	for i, mb := range header.MiniBlockHeaders {
		require.Equal(t, mb.Hash, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].Hash)
		require.Equal(t, mb.SenderShardID, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].SenderShardID)
		require.Equal(t, mb.ReceiverShardID, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].ReceiverShardID)
		require.Equal(t, mb.TxCount, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].TxCount)
		require.Equal(t, mb.Type, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].Type)
		require.Equal(t, mb.Reserved, fireOutportBlock.BlockData.Header.MiniBlockHeaders[i].Reserved)
	}

	require.Equal(t, header.ReceiptsHash, fireOutportBlock.BlockData.Header.ReceiptsHash)

	for i, lfh := range header.EpochStart.LastFinalizedHeaders {
		require.Equal(t, lfh.ShardID, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].ShardID)
		require.Equal(t, lfh.Epoch, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].Epoch)
		require.Equal(t, lfh.Round, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].Round)
		require.Equal(t, lfh.Nonce, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].Nonce)
		require.Equal(t, lfh.HeaderHash, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].HeaderHash)
		require.Equal(t, lfh.RootHash, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].RootHash)
		require.Equal(t, lfh.ScheduledRootHash, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].ScheduledRootHash)
		require.Equal(t, lfh.FirstPendingMetaBlock, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].FirstPendingMetaBlock)
		require.Equal(t, lfh.LastFinishedMetaBlock, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].LastFinishedMetaBlock)

		for j, mbh := range lfh.PendingMiniBlockHeaders {
			require.Equal(t, mbh.Hash, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].PendingMiniBlockHeaders[j].Hash)
			require.Equal(t, mbh.SenderShardID, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].PendingMiniBlockHeaders[j].SenderShardID)
			require.Equal(t, mbh.ReceiverShardID, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].PendingMiniBlockHeaders[j].ReceiverShardID)
			require.Equal(t, mbh.TxCount, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].PendingMiniBlockHeaders[j].TxCount)
			require.Equal(t, mbh.Type, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].PendingMiniBlockHeaders[j].Type)
			require.Equal(t, mbh.Reserved, fireOutportBlock.BlockData.Header.EpochStart.LastFinalizedHeaders[i].PendingMiniBlockHeaders[j].Reserved)
		}
	}

	require.Equal(t, mustCastBigInt(t, header.EpochStart.Economics.TotalSupply), fireOutportBlock.BlockData.Header.EpochStart.Economics.TotalSupply)
	require.Equal(t, mustCastBigInt(t, header.EpochStart.Economics.TotalToDistribute), fireOutportBlock.BlockData.Header.EpochStart.Economics.TotalToDistribute)
	require.Equal(t, mustCastBigInt(t, header.EpochStart.Economics.TotalNewlyMinted), fireOutportBlock.BlockData.Header.EpochStart.Economics.TotalNewlyMinted)
	require.Equal(t, mustCastBigInt(t, header.EpochStart.Economics.RewardsPerBlock), fireOutportBlock.BlockData.Header.EpochStart.Economics.RewardsPerBlock)
	require.Equal(t, mustCastBigInt(t, header.EpochStart.Economics.RewardsForProtocolSustainability), fireOutportBlock.BlockData.Header.EpochStart.Economics.RewardsForProtocolSustainability)
	require.Equal(t, mustCastBigInt(t, header.EpochStart.Economics.NodePrice), fireOutportBlock.BlockData.Header.EpochStart.Economics.NodePrice)
	require.Equal(t, header.EpochStart.Economics.PrevEpochStartRound, fireOutportBlock.BlockData.Header.EpochStart.Economics.PrevEpochStartRound)
	require.Equal(t, header.EpochStart.Economics.PrevEpochStartHash, fireOutportBlock.BlockData.Header.EpochStart.Economics.PrevEpochStartHash)

	require.Equal(t, header.ChainID, fireOutportBlock.BlockData.Header.ChainID)
	require.Equal(t, header.SoftwareVersion, fireOutportBlock.BlockData.Header.SoftwareVersion)
	require.Equal(t, mustCastBigInt(t, header.AccumulatedFees), fireOutportBlock.BlockData.Header.AccumulatedFees)
	require.Equal(t, mustCastBigInt(t, header.AccumulatedFeesInEpoch), fireOutportBlock.BlockData.Header.AccumulatedFeesInEpoch)
	require.Equal(t, mustCastBigInt(t, header.DeveloperFees), fireOutportBlock.BlockData.Header.DeveloperFees)
	require.Equal(t, mustCastBigInt(t, header.DevFeesInEpoch), fireOutportBlock.BlockData.Header.DevFeesInEpoch)
	require.Equal(t, header.TxCount, fireOutportBlock.BlockData.Header.TxCount)
	require.Equal(t, header.Reserved, fireOutportBlock.BlockData.Header.Reserved)
}

func checkFieldsV1(t *testing.T, outportBlock *outport.OutportBlock, fireOutportBlock fieldsGetterV1) {
	// Asserting values.
	require.Equal(t, outportBlock.ShardID, fireOutportBlock.GetShardID())
	require.Equal(t, outportBlock.NotarizedHeadersHashes, fireOutportBlock.GetNotarizedHeadersHashes())
	require.Equal(t, outportBlock.NumberOfShards, fireOutportBlock.GetNumberOfShards())
	require.Equal(t, outportBlock.SignersIndexes, fireOutportBlock.GetSignersIndexes())
	require.Equal(t, outportBlock.HighestFinalBlockNonce, fireOutportBlock.GetHighestFinalBlockNonce())
	require.Equal(t, outportBlock.HighestFinalBlockHash, fireOutportBlock.GetHighestFinalBlockHash())

	// Transaction pool - Transactions.
	for k, v := range outportBlock.TransactionPool.Transactions {
		// Transaction pool - Transactions. - TxInfo.
		require.Equal(t, v.Transaction.Nonce, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Nonce)
		require.Equal(t, mustCastBigInt(t, v.Transaction.Value), fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Value)
		require.Equal(t, v.Transaction.RcvAddr, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.RcvAddr)
		require.Equal(t, v.Transaction.RcvUserName, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.RcvUserName)
		require.Equal(t, v.Transaction.SndAddr, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.SndAddr)
		require.Equal(t, v.Transaction.GasPrice, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.GasPrice)
		require.Equal(t, v.Transaction.GasLimit, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.GasLimit)
		require.Equal(t, v.Transaction.Data, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Data)
		require.Equal(t, v.Transaction.ChainID, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.ChainID)
		require.Equal(t, v.Transaction.Version, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Version)
		require.Equal(t, v.Transaction.Signature, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Signature)
		require.Equal(t, v.Transaction.Options, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Options)
		require.Equal(t, v.Transaction.GuardianAddr, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.GuardianAddr)
		require.Equal(t, v.Transaction.GuardianSignature, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.GuardianSignature)

		// Transaction pool - Transactions - Tx Info - Fee info.
		require.Equal(t, v.FeeInfo.GasUsed, fireOutportBlock.GetTransactionPool().Transactions[k].FeeInfo.GasUsed)
		require.Equal(t, mustCastBigInt(t, v.FeeInfo.Fee), fireOutportBlock.GetTransactionPool().Transactions[k].FeeInfo.Fee)
		require.Equal(t, mustCastBigInt(t, v.FeeInfo.InitialPaidFee), fireOutportBlock.GetTransactionPool().Transactions[k].FeeInfo.InitialPaidFee)

		require.Equal(t, v.ExecutionOrder, fireOutportBlock.GetTransactionPool().Transactions[k].ExecutionOrder)
	}

	// Transaction pool - Smart Contract results.
	for k, v := range outportBlock.TransactionPool.SmartContractResults {
		// Transaction pool - Smart Contract results - SmartContractResult.
		require.Equal(t, v.SmartContractResult.Nonce, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.Nonce)
		require.Equal(t, mustCastBigInt(t, v.SmartContractResult.Value), fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.Value)
		require.Equal(t, v.SmartContractResult.RcvAddr, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.RcvAddr)
		require.Equal(t, v.SmartContractResult.SndAddr, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.SndAddr)
		require.Equal(t, v.SmartContractResult.RelayerAddr, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.RelayerAddr)
		require.Equal(t, mustCastBigInt(t, v.SmartContractResult.RelayedValue), fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.RelayedValue)
		require.Equal(t, v.SmartContractResult.Code, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.Code)
		require.Equal(t, v.SmartContractResult.Data, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.Data)
		require.Equal(t, v.SmartContractResult.PrevTxHash, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.PrevTxHash)
		require.Equal(t, v.SmartContractResult.OriginalTxHash, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.OriginalTxHash)
		require.Equal(t, v.SmartContractResult.GasLimit, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.GasLimit)
		require.Equal(t, v.SmartContractResult.GasPrice, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.GasPrice)
		require.Equal(t, int64(v.SmartContractResult.CallType), fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.CallType)
		require.Equal(t, v.SmartContractResult.CodeMetadata, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.CodeMetadata)
		require.Equal(t, v.SmartContractResult.ReturnMessage, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.ReturnMessage)
		require.Equal(t, v.SmartContractResult.OriginalSender, fireOutportBlock.GetTransactionPool().SmartContractResults[k].SmartContractResult.OriginalSender)

		// Transaction pool - Smart Contract results - Fee info.
		require.Equal(t, v.FeeInfo.GasUsed, fireOutportBlock.GetTransactionPool().SmartContractResults[k].FeeInfo.GasUsed)
		require.Equal(t, mustCastBigInt(t, v.FeeInfo.Fee), fireOutportBlock.GetTransactionPool().SmartContractResults[k].FeeInfo.Fee)
		require.Equal(t, mustCastBigInt(t, v.FeeInfo.InitialPaidFee), fireOutportBlock.GetTransactionPool().SmartContractResults[k].FeeInfo.InitialPaidFee)

		// Transaction pool - Smart Contract results - Execution Order.
		require.Equal(t, v.ExecutionOrder, fireOutportBlock.GetTransactionPool().SmartContractResults[k].ExecutionOrder)
	}

	// Transaction Pool - Rewards
	for k, v := range outportBlock.TransactionPool.Rewards {
		// Transaction Pool - Rewards - Reward info
		require.Equal(t, v.Reward.Round, fireOutportBlock.GetTransactionPool().Rewards[k].Reward.Round)
		require.Equal(t, mustCastBigInt(t, v.Reward.Value), fireOutportBlock.GetTransactionPool().Rewards[k].Reward.Value)
		require.Equal(t, v.Reward.RcvAddr, fireOutportBlock.GetTransactionPool().Rewards[k].Reward.RcvAddr)
		require.Equal(t, v.Reward.Epoch, fireOutportBlock.GetTransactionPool().Rewards[k].Reward.Epoch)

		// Transaction Pool - Rewards - Execution Order
		require.Equal(t, v.ExecutionOrder, fireOutportBlock.GetTransactionPool().Rewards[k].ExecutionOrder)
	}

	// Transaction Pool - Receipts
	for k, v := range outportBlock.TransactionPool.Receipts {
		// Transaction Pool - Receipts - Receipt info
		require.Equal(t, mustCastBigInt(t, v.Value), fireOutportBlock.GetTransactionPool().Receipts[k].Value)
		require.Equal(t, v.SndAddr, fireOutportBlock.GetTransactionPool().Receipts[k].SndAddr)
		require.Equal(t, v.Data, fireOutportBlock.GetTransactionPool().Receipts[k].Data)
		require.Equal(t, v.TxHash, fireOutportBlock.GetTransactionPool().Receipts[k].TxHash)
	}

	// Transaction Pool - Invalid Txs
	for k, v := range outportBlock.TransactionPool.InvalidTxs {
		// Transaction Pool - Invalid Txs - Tx Info
		require.Equal(t, v.Transaction.Nonce, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.Nonce)
		require.Equal(t, mustCastBigInt(t, v.Transaction.Value), fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.Value)
		require.Equal(t, v.Transaction.RcvAddr, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.RcvAddr)
		require.Equal(t, v.Transaction.RcvUserName, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.RcvUserName)
		require.Equal(t, v.Transaction.SndAddr, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.SndAddr)
		require.Equal(t, v.Transaction.GasPrice, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.GasPrice)
		require.Equal(t, v.Transaction.GasLimit, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.GasLimit)
		require.Equal(t, v.Transaction.Data, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.Data)
		require.Equal(t, v.Transaction.ChainID, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.ChainID)
		require.Equal(t, v.Transaction.Version, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.Version)
		require.Equal(t, v.Transaction.Signature, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.Signature)
		require.Equal(t, v.Transaction.Options, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.Options)
		require.Equal(t, v.Transaction.GuardianAddr, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.GuardianAddr)
		require.Equal(t, v.Transaction.GuardianSignature, fireOutportBlock.GetTransactionPool().InvalidTxs[k].Transaction.GuardianSignature)
		require.Equal(t, v.ExecutionOrder, fireOutportBlock.GetTransactionPool().InvalidTxs[k].ExecutionOrder)

		// Transaction pool - Invalid Txs - Fee info.
		require.Equal(t, v.FeeInfo.GasUsed, fireOutportBlock.GetTransactionPool().InvalidTxs[k].FeeInfo.GasUsed)
		require.Equal(t, mustCastBigInt(t, v.FeeInfo.Fee), fireOutportBlock.GetTransactionPool().InvalidTxs[k].FeeInfo.Fee)
		require.Equal(t, mustCastBigInt(t, v.FeeInfo.InitialPaidFee), fireOutportBlock.GetTransactionPool().InvalidTxs[k].FeeInfo.InitialPaidFee)

		require.Equal(t, v.ExecutionOrder, fireOutportBlock.GetTransactionPool().InvalidTxs[k].ExecutionOrder)
	}

	// Transaction Pool - Logs
	for i, l := range outportBlock.TransactionPool.Logs {
		// Transaction Pool - Logs - Log Data
		require.Equal(t, l.TxHash, fireOutportBlock.GetTransactionPool().Logs[i].TxHash)

		// Transaction Pool - Logs - Log data - Log
		require.Equal(t, l.Log.Address, fireOutportBlock.GetTransactionPool().Logs[i].Log.Address)

		for k, e := range outportBlock.TransactionPool.Logs[i].Log.Events {
			require.Equal(t, e.Address, fireOutportBlock.GetTransactionPool().Logs[i].Log.Events[k].Address)
			require.Equal(t, e.Identifier, fireOutportBlock.GetTransactionPool().Logs[i].Log.Events[k].Identifier)
			require.Equal(t, e.Topics, fireOutportBlock.GetTransactionPool().Logs[i].Log.Events[k].Topics)
			require.Equal(t, e.Data, fireOutportBlock.GetTransactionPool().Logs[i].Log.Events[k].Data)
			require.Equal(t, e.AdditionalData, fireOutportBlock.GetTransactionPool().Logs[i].Log.Events[k].AdditionalData)
		}
	}

	// Transaction Pool - ScheduledExecutedSCRSHashesPrevBlock
	for i, s := range outportBlock.TransactionPool.ScheduledExecutedSCRSHashesPrevBlock {
		require.Equal(t, s, fireOutportBlock.GetTransactionPool().ScheduledExecutedSCRSHashesPrevBlock[i])
	}

	// Transaction Pool - ScheduledExecutedInvalidTxsHashesPrevBlock
	for i, s := range outportBlock.TransactionPool.ScheduledExecutedInvalidTxsHashesPrevBlock {
		require.Equal(t, s, fireOutportBlock.GetTransactionPool().ScheduledExecutedInvalidTxsHashesPrevBlock[i])
	}

	// Header gas consumption.
	require.Equal(t, outportBlock.HeaderGasConsumption.GasProvided, fireOutportBlock.GetHeaderGasConsumption().GasProvided)
	require.Equal(t, outportBlock.HeaderGasConsumption.GasRefunded, fireOutportBlock.GetHeaderGasConsumption().GasRefunded)
	require.Equal(t, outportBlock.HeaderGasConsumption.GasPenalized, fireOutportBlock.GetHeaderGasConsumption().GasPenalized)
	require.Equal(t, outportBlock.HeaderGasConsumption.MaxGasPerBlock, fireOutportBlock.GetHeaderGasConsumption().MaxGasPerBlock)

	// Altered accounts.
	for key, account := range outportBlock.AlteredAccounts {
		acc := fireOutportBlock.GetAlteredAccounts()[key]

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
			if token.MetaData != nil {
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
		}

		//  Altered accounts - Additional account data.
		if account.AdditionalData != nil {
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
}

func checkFieldsV2(t *testing.T, outportBlock *outport.OutportBlock, fireOutportBlock fieldsGetterV2) {
	// Asserting values.
	require.Equal(t, outportBlock.ShardID, fireOutportBlock.GetShardID())
	require.Equal(t, outportBlock.NotarizedHeadersHashes, fireOutportBlock.GetNotarizedHeadersHashes())
	require.Equal(t, outportBlock.NumberOfShards, fireOutportBlock.GetNumberOfShards())
	require.Equal(t, outportBlock.SignersIndexes, fireOutportBlock.GetSignersIndexes())
	require.Equal(t, outportBlock.HighestFinalBlockNonce, fireOutportBlock.GetHighestFinalBlockNonce())
	require.Equal(t, outportBlock.HighestFinalBlockHash, fireOutportBlock.GetHighestFinalBlockHash())

	// Transaction pool - Transactions.
	for k, v := range outportBlock.TransactionPool.Transactions {
		// Transaction pool - Transactions. - TxInfo.
		require.Equal(t, v.Transaction.Nonce, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Nonce)
		require.Equal(t, mustCastBigInt(t, v.Transaction.Value), fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Value)
		require.Equal(t, v.Transaction.RcvAddr, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.RcvAddr)
		require.Equal(t, v.Transaction.RcvUserName, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.RcvUserName)
		require.Equal(t, v.Transaction.SndAddr, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.SndAddr)
		require.Equal(t, v.Transaction.GasPrice, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.GasPrice)
		require.Equal(t, v.Transaction.GasLimit, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.GasLimit)
		require.Equal(t, v.Transaction.Data, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Data)
		require.Equal(t, v.Transaction.ChainID, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.ChainID)
		require.Equal(t, v.Transaction.Version, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Version)
		require.Equal(t, v.Transaction.Signature, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Signature)
		require.Equal(t, v.Transaction.Options, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Options)
		require.Equal(t, v.Transaction.GuardianAddr, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.GuardianAddr)
		require.Equal(t, v.Transaction.GuardianSignature, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.GuardianSignature)

		// Transaction pool - Transactions - Tx Info - Fee info.
		require.Equal(t, v.FeeInfo.GasUsed, fireOutportBlock.GetTransactionPool().Transactions[k].FeeInfo.GasUsed)
		require.Equal(t, mustCastBigInt(t, v.FeeInfo.Fee), fireOutportBlock.GetTransactionPool().Transactions[k].FeeInfo.Fee)
		require.Equal(t, mustCastBigInt(t, v.FeeInfo.InitialPaidFee), fireOutportBlock.GetTransactionPool().Transactions[k].FeeInfo.InitialPaidFee)
	}

	// Transaction pool - Smart Contract results.
	for k, v := range outportBlock.TransactionPool.SmartContractResults {
		tx := fireOutportBlock.GetTransactionPool().Transactions[k]
		// Transaction pool - Smart Contract results - SmartContractResult.
		require.Equal(t, v.SmartContractResult.Nonce, tx.Transaction.Nonce)
		require.Equal(t, mustCastBigInt(t, v.SmartContractResult.Value), tx.Transaction.Value)
		require.Equal(t, v.SmartContractResult.RcvAddr, tx.Transaction.RcvAddr)
		require.Equal(t, v.SmartContractResult.SndAddr, tx.Transaction.SndAddr)
		require.Equal(t, v.SmartContractResult.RelayerAddr, tx.Transaction.RelayerAddr)
		require.Equal(t, mustCastBigInt(t, v.SmartContractResult.RelayedValue), tx.Transaction.RelayedValue)
		require.Equal(t, v.SmartContractResult.Code, tx.Transaction.Code)
		require.Equal(t, v.SmartContractResult.Data, tx.Transaction.Data)
		require.Equal(t, v.SmartContractResult.PrevTxHash, tx.Transaction.PrevTxHash)
		require.Equal(t, v.SmartContractResult.OriginalTxHash, tx.Transaction.OriginalTxHash)
		require.Equal(t, v.SmartContractResult.GasLimit, tx.Transaction.GasLimit)
		require.Equal(t, v.SmartContractResult.GasPrice, tx.Transaction.GasPrice)
		require.Equal(t, int64(v.SmartContractResult.CallType), tx.Transaction.CallType)
		require.Equal(t, v.SmartContractResult.CodeMetadata, tx.Transaction.CodeMetadata)
		require.Equal(t, v.SmartContractResult.ReturnMessage, tx.Transaction.ReturnMessage)
		require.Equal(t, v.SmartContractResult.OriginalSender, tx.Transaction.OriginalSender)
		require.Equal(t, data.TxType_SCR, tx.Transaction.TxType)

		require.Equal(t, v.FeeInfo.GasUsed, tx.FeeInfo.GasUsed)
		require.Equal(t, mustCastBigInt(t, v.FeeInfo.Fee), tx.FeeInfo.Fee)
		require.Equal(t, mustCastBigInt(t, v.FeeInfo.InitialPaidFee), tx.FeeInfo.InitialPaidFee)

		// Transaction pool - Smart Contract results - Execution Order.
		require.Equal(t, v.ExecutionOrder, tx.Transaction.ExecutionOrder)
	}

	// Transaction Pool - Rewards
	for k, v := range outportBlock.TransactionPool.Rewards {
		// Transaction Pool - Rewards - Reward info
		require.Equal(t, v.Reward.Round, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Round)
		require.Equal(t, mustCastBigInt(t, v.Reward.Value), fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Value)
		require.Equal(t, v.Reward.RcvAddr, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.RcvAddr)
		require.Equal(t, v.Reward.Epoch, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Epoch)
		require.Equal(t, []byte("metachain"), fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.SndAddr)

		// Transaction Pool - Rewards - Execution Order
		require.Equal(t, v.ExecutionOrder, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.ExecutionOrder)

		require.Equal(t, data.TxType_Reward, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.TxType)
	}

	// Transaction Pool - Receipts
	for k, v := range outportBlock.TransactionPool.Receipts {
		// Transaction Pool - Receipts - Receipt info
		require.Equal(t, mustCastBigInt(t, v.Value), fireOutportBlock.GetTransactionPool().Transactions[k].Receipt.Value)
		require.Equal(t, v.SndAddr, fireOutportBlock.GetTransactionPool().Transactions[k].Receipt.SndAddr)
		require.Equal(t, v.Data, fireOutportBlock.GetTransactionPool().Transactions[k].Receipt.Data)
		require.Equal(t, v.TxHash, fireOutportBlock.GetTransactionPool().Transactions[k].Receipt.TxHash)
	}

	// Transaction Pool - Invalid Txs
	for k, v := range outportBlock.TransactionPool.InvalidTxs {
		// Transaction Pool - Invalid Txs - Tx Info
		require.Equal(t, v.Transaction.Nonce, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Nonce)
		require.Equal(t, mustCastBigInt(t, v.Transaction.Value), fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Value)
		require.Equal(t, v.Transaction.RcvAddr, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.RcvAddr)
		require.Equal(t, v.Transaction.RcvUserName, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.RcvUserName)
		require.Equal(t, v.Transaction.SndAddr, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.SndAddr)
		require.Equal(t, v.Transaction.GasPrice, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.GasPrice)
		require.Equal(t, v.Transaction.GasLimit, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.GasLimit)
		require.Equal(t, v.Transaction.Data, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Data)
		require.Equal(t, v.Transaction.ChainID, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.ChainID)
		require.Equal(t, v.Transaction.Version, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Version)
		require.Equal(t, v.Transaction.Signature, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Signature)
		require.Equal(t, v.Transaction.Options, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.Options)
		require.Equal(t, v.Transaction.GuardianAddr, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.GuardianAddr)
		require.Equal(t, v.Transaction.GuardianSignature, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.GuardianSignature)
		require.Equal(t, v.ExecutionOrder, fireOutportBlock.GetTransactionPool().Transactions[k].Transaction.ExecutionOrder)

		// Transaction pool - Invalid Txs - Fee info.
		require.Equal(t, v.FeeInfo.GasUsed, fireOutportBlock.GetTransactionPool().Transactions[k].FeeInfo.GasUsed)
		require.Equal(t, mustCastBigInt(t, v.FeeInfo.Fee), fireOutportBlock.GetTransactionPool().Transactions[k].FeeInfo.Fee)
		require.Equal(t, mustCastBigInt(t, v.FeeInfo.InitialPaidFee), fireOutportBlock.GetTransactionPool().Transactions[k].FeeInfo.InitialPaidFee)
	}

	// Transaction Pool - Logs
	for _, l := range outportBlock.TransactionPool.Logs {
		require.NotNil(t, fireOutportBlock.GetTransactionPool().Transactions[l.TxHash])

		// Transaction Pool - Logs - Log Data

		// Transaction Pool - Logs - Log data - Log
		var ok bool
		for i, ll := range fireOutportBlock.GetTransactionPool().Transactions[l.TxHash].Logs {
			if reflect.DeepEqual(l.Log.Address, ll.Address) {
				ok = true
			}

			for k, e := range l.Log.Events {
				require.Equal(t, e.Address, fireOutportBlock.GetTransactionPool().Transactions[l.TxHash].Logs[i].Events[k].Address)
				require.Equal(t, e.Identifier, fireOutportBlock.GetTransactionPool().Transactions[l.TxHash].Logs[i].Events[k].Identifier)
				require.Equal(t, e.Topics, fireOutportBlock.GetTransactionPool().Transactions[l.TxHash].Logs[i].Events[k].Topics)
				require.Equal(t, e.Data, fireOutportBlock.GetTransactionPool().Transactions[l.TxHash].Logs[i].Events[k].Data)
				require.Equal(t, e.AdditionalData, fireOutportBlock.GetTransactionPool().Transactions[l.TxHash].Logs[i].Events[k].AdditionalData)
			}
		}

		require.True(t, ok)
	}

	// Transaction Pool - ScheduledExecutedSCRSHashesPrevBlock
	for i, s := range outportBlock.TransactionPool.ScheduledExecutedSCRSHashesPrevBlock {
		require.Equal(t, s, fireOutportBlock.GetTransactionPool().ScheduledExecutedSCRSHashesPrevBlock[i])
	}

	// Transaction Pool - ScheduledExecutedInvalidTxsHashesPrevBlock
	for i, s := range outportBlock.TransactionPool.ScheduledExecutedInvalidTxsHashesPrevBlock {
		require.Equal(t, s, fireOutportBlock.GetTransactionPool().ScheduledExecutedInvalidTxsHashesPrevBlock[i])
	}

	// Header gas consumption.
	require.Equal(t, outportBlock.HeaderGasConsumption.GasProvided, fireOutportBlock.GetHeaderGasConsumption().GasProvided)
	require.Equal(t, outportBlock.HeaderGasConsumption.GasRefunded, fireOutportBlock.GetHeaderGasConsumption().GasRefunded)
	require.Equal(t, outportBlock.HeaderGasConsumption.GasPenalized, fireOutportBlock.GetHeaderGasConsumption().GasPenalized)
	require.Equal(t, outportBlock.HeaderGasConsumption.MaxGasPerBlock, fireOutportBlock.GetHeaderGasConsumption().MaxGasPerBlock)

	// Altered accounts.
	for key, account := range outportBlock.AlteredAccounts {
		acc := fireOutportBlock.GetAlteredAccounts()[key]

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
			if token.MetaData != nil {
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
		}

		//  Altered accounts - Additional account data.
		if account.AdditionalData != nil {
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
}

func checkBlockData(t *testing.T, blockData *outport.BlockData, getter blockDataGetter) {
	require.Equal(t, blockData.ShardID, getter.GetShardID())
	require.Equal(t, blockData.HeaderType, getter.GetHeaderType())
	require.Equal(t, blockData.HeaderHash, getter.GetHeaderHash())

	for i, miniBlock := range blockData.GetBody().MiniBlocks {
		require.Equal(t, miniBlock.TxHashes, getter.GetBody().MiniBlocks[i].TxHashes)
		require.Equal(t, miniBlock.ReceiverShardID, getter.GetBody().MiniBlocks[i].ReceiverShardID)
		require.Equal(t, miniBlock.SenderShardID, getter.GetBody().MiniBlocks[i].SenderShardID)
		require.Equal(t, data.Type(miniBlock.Type), getter.GetBody().MiniBlocks[i].Type)
		require.Equal(t, miniBlock.Reserved, getter.GetBody().MiniBlocks[i].Reserved)
	}

	for i, intraShardMiniBlock := range blockData.GetIntraShardMiniBlocks() {
		require.Equal(t, intraShardMiniBlock.TxHashes, getter.GetIntraShardMiniBlocks()[i].TxHashes)
		require.Equal(t, intraShardMiniBlock.ReceiverShardID, getter.GetIntraShardMiniBlocks()[i].ReceiverShardID)
		require.Equal(t, intraShardMiniBlock.SenderShardID, getter.GetIntraShardMiniBlocks()[i].SenderShardID)
		require.Equal(t, intraShardMiniBlock.Reserved, getter.GetIntraShardMiniBlocks()[i].Reserved)
	}
}

func mustCastBigInt(t *testing.T, i *big.Int) []byte {
	t.Helper()

	converter, _ := process.NewOutportBlockConverter(gogoProtoMarshaller, protoMarshaller)
	buf, err := converter.CastBigInt(i)
	require.NoError(t, err, "failed to cast from big.Int")

	return buf
}
