package process

import (
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"

	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
)

type outportBlockConverter struct {
	gogoProtoMarshalizer marshal.Marshalizer
	protoMarshalizer     marshal.Marshalizer
	bigIntCaster         coreData.BigIntCaster
}

// NewOutportBlockConverter will return a component than can convert outport.OutportBlock to either data.ShardOutportBlock.
// or data.MetaOutportBlock
func NewOutportBlockConverter() *outportBlockConverter {
	return &outportBlockConverter{
		gogoProtoMarshalizer: &marshal.GogoProtoMarshalizer{},
		protoMarshalizer:     &ProtoMarshalizer{},
		bigIntCaster:         coreData.BigIntCaster{},
	}
}

// HandleShardOutportBlock will convert an outport.OutportBlock to data.ShardOutportBlock.
func (o *outportBlockConverter) HandleShardOutportBlock(outportBlock *outport.OutportBlock) (*data.ShardOutportBlock, error) {
	headerType := outportBlock.BlockData.HeaderType

	// check if the header type is supported by this function.
	if headerType != string(core.ShardHeaderV1) && headerType != string(core.ShardHeaderV2) {
		return nil, fmt.Errorf("cannot convert to shard outport block. header type: %s not supported", headerType)
	}

	// marshal with gogo, since the outportBlock is gogo protobuf (coming from the node).
	bytes, err := o.gogoProtoMarshalizer.Marshal(outportBlock)
	if err != nil {
		return nil, fmt.Errorf("marshal shard outport block error: %s", err)
	}

	shardOutportBlock := &data.ShardOutportBlock{}
	// unmarshall into google protobuf. This is the proto that will be used in firehose.
	err = o.protoMarshalizer.Unmarshal(shardOutportBlock, bytes)
	if err != nil {
		return nil, fmt.Errorf("unmarshal shard outport block error: %s", err)
	}

	// ShardHeaderV1 marshals 1 to 1 into *data.ShardOutportBlock.
	if headerType == string(core.ShardHeaderV1) {
		return shardOutportBlock, nil
	}

	// ShardHeaderV2 does not marshal 1 to 1. A few fields need to be injected.
	header := block.HeaderV2{}
	err = o.gogoProtoMarshalizer.Unmarshal(&header, outportBlock.BlockData.HeaderBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	miniBlockHeaders := make([]*data.MiniBlockHeader, 0)
	for _, miniBlockHeader := range header.Header.MiniBlockHeaders {
		mb := &data.MiniBlockHeader{
			Hash:            miniBlockHeader.Hash,
			SenderShardID:   miniBlockHeader.SenderShardID,
			ReceiverShardID: miniBlockHeader.ReceiverShardID,
			TxCount:         miniBlockHeader.TxCount,
			Type:            data.Type(miniBlockHeader.Type),
			Reserved:        miniBlockHeader.Reserved,
		}
		miniBlockHeaders = append(miniBlockHeaders, mb)
	}

	peerChanges := make([]*data.PeerChange, 0)
	for _, peerChange := range header.Header.PeerChanges {
		pc := &data.PeerChange{
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

	shardOutportBlock.BlockData.Header = &data.Header{
		Nonce:              header.Header.Nonce,
		PrevHash:           header.Header.PrevHash,
		PrevRandSeed:       header.Header.PrevRandSeed,
		RandSeed:           header.Header.RandSeed,
		PubKeysBitmap:      header.Header.PubKeysBitmap,
		ShardID:            header.Header.ShardID,
		TimeStamp:          header.Header.TimeStamp,
		Round:              header.Header.Round,
		Epoch:              header.Header.Epoch,
		BlockBodyType:      data.Type(header.Header.BlockBodyType),
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
func (o *outportBlockConverter) HandleMetaOutportBlock(outportBlock *outport.OutportBlock) (*data.MetaOutportBlock, error) {
	headerType := outportBlock.BlockData.HeaderType

	// check if the header type is supported by this function.
	if headerType != string(core.MetaHeader) {
		return nil, fmt.Errorf("cannot convert to meta outport block. header type: %s not supported", outportBlock.BlockData.HeaderType)
	}

	// marshal with gogo, since the outportBlock is gogo protobuf (coming from the node).
	bytes, err := o.gogoProtoMarshalizer.Marshal(outportBlock)
	if err != nil {
		return nil, fmt.Errorf("marshal metaBlockCaster error: %w", err)
	}

	// unmarshall into google protobuf. This is the proto that will be used in firehose.
	metaOutportBlock := &data.MetaOutportBlock{}
	err = o.protoMarshalizer.Unmarshal(metaOutportBlock, bytes)
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
