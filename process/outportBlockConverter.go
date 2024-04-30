package process

import (
	"fmt"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
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
) *outportBlockConverter {
	return &outportBlockConverter{
		gogoProtoMarshaller: gogoMarshaller,
		protoMarshaller:     protoMarshaller,
		bigIntCaster:        coreData.BigIntCaster{},
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
