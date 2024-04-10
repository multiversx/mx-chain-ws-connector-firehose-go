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
	bigIntCaster    coreData.BigIntCaster
)

type (
	headerV1Caster struct {
		outportBlock *outport.OutportBlock
	}
	headerV2Caster struct {
		outportBlock *outport.OutportBlock
	}
	metaBlockCaster struct {
		outportBlock *outport.OutportBlock
	}

	// HeaderCaster is the interface used for the caster to convert from outport.OutportBlock to
	// data.ShardOutportBlock or data.MetaOutportBlock
	HeaderCaster interface {
		Cast() (proto.Message, error)
	}
)

// NewHeaderV1Caster will return an instance of a caster used to cast outport.OutportBlocks with a headerV1.
func NewHeaderV1Caster(outportBlock *outport.OutportBlock) (*headerV1Caster, error) {
	return &headerV1Caster{
		outportBlock,
	}, nil
}

// NewHeaderV2Caster will return an instance of a caster used to cast outport.OutportBlocks with a headerV2.
func NewHeaderV2Caster(outportBlock *outport.OutportBlock) (*headerV2Caster, error) {
	return &headerV2Caster{
		outportBlock,
	}, nil
}

// NewMetaBlockCaster will return an instance of a caster used to cast outport.OutportBlocks with a metaHeader.
func NewMetaBlockCaster(outportBlock *outport.OutportBlock) (*metaBlockCaster, error) {
	return &metaBlockCaster{
		outportBlock,
	}, nil
}

// Cast will convert the underlying outport.OutportBlock to data.ShardOutportBlock.
func (h *headerV1Caster) Cast() (proto.Message, error) {
	return castShardOutportBlock(h.outportBlock)
}

// Cast will convert the underlying outport.OutportBlock to data.ShardOutportBlock.
func (h *headerV2Caster) Cast() (proto.Message, error) {
	fireOutportBlock, err := castShardOutportBlock(h.outportBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to cast: %w", err)
	}

	header := block.HeaderV2{}
	err = protoMarshaller.Unmarshal(&header, h.outportBlock.BlockData.HeaderBytes)
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

	accumulatedFees, err := castBigInt(header.Header.AccumulatedFees)
	if err != nil {
		return nil, fmt.Errorf("failed to cast accumulated fees: %w", err)
	}
	developerFees, err := castBigInt(header.Header.DeveloperFees)
	if err != nil {
		return nil, fmt.Errorf("failed to cast developer fees: %w", err)
	}

	fireOutportBlock.BlockData.Header = &data.Header{
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

	scheduledAccumulatedFees, err := castBigInt(header.ScheduledAccumulatedFees)
	if err != nil {
		return nil, fmt.Errorf("failed to cast scheduled accumulated fees: %w", err)
	}
	scheduledDeveloperFees, err := castBigInt(header.ScheduledDeveloperFees)
	if err != nil {
		return nil, fmt.Errorf("failed to cast scheduled developer fees: %w", err)
	}

	fireOutportBlock.BlockData.ScheduledRootHash = header.ScheduledRootHash
	fireOutportBlock.BlockData.ScheduledAccumulatedFees = scheduledAccumulatedFees
	fireOutportBlock.BlockData.ScheduledDeveloperFees = scheduledDeveloperFees
	fireOutportBlock.BlockData.ScheduledGasProvided = header.ScheduledGasProvided
	fireOutportBlock.BlockData.ScheduledGasPenalized = header.ScheduledGasPenalized
	fireOutportBlock.BlockData.ScheduledGasRefunded = header.ScheduledGasRefunded

	return fireOutportBlock, nil
}

// Cast will convert the underlying outport.OutportBlock to data.MetaOutportBlock.
func (h *metaBlockCaster) Cast() (proto.Message, error) {
	bytes, err := protoMarshaller.Marshal(h.outportBlock)
	if err != nil {
		return nil, fmt.Errorf("marshal metaBlockCaster error: %w", err)
	}

	metaOutportBlock := &data.MetaOutportBlock{}
	err = proto.Unmarshal(bytes, metaOutportBlock)
	if err != nil {
		return nil, fmt.Errorf("unmarshal metaBlockCaster error: %w", err)
	}

	return metaOutportBlock, nil
}

func castShardOutportBlock(outportBlock *outport.OutportBlock) (*data.ShardOutportBlock, error) {
	bytes, err := protoMarshaller.Marshal(outportBlock)
	if err != nil {
		return nil, fmt.Errorf("marshal shard outport block error: %s", err)
	}

	shardOutportBlock := &data.ShardOutportBlock{}
	err = proto.Unmarshal(bytes, shardOutportBlock)
	if err != nil {
		return nil, fmt.Errorf("unmarshal shard outport block error: %s", err)
	}

	return shardOutportBlock, nil
}

func castBigInt(i *big.Int) ([]byte, error) {
	buf := make([]byte, bigIntCaster.Size(i))
	_, err := bigIntCaster.MarshalTo(i, buf)

	return buf, err
}
