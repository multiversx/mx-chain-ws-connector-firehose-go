package factory

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/data/outport"
	"google.golang.org/protobuf/proto"

	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
)

type headerType = string

const (
	header    headerType = "Header"
	headerV2  headerType = "HeaderV2"
	metaBlock headerType = "MetaBlock"
)

// CastOutportBlock will cast to a proto.Message depending on the headerType.
// The proto.Message can be an instance of *data.ShardOutportBlock or *data.MetaOutportBlock.
func CastOutportBlock(outportBlock *outport.OutportBlock) (proto.Message, error) {
	caster, err := GetOutportCaster(outportBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to get caster: %w", err)
	}

	return caster.Cast()
}

// GetOutportCaster will retrieve the implementation used to cast the outportBlock based on the headerType.
func GetOutportCaster(outportBlock *outport.OutportBlock) (process.HeaderCaster, error) {
	switch outportBlock.BlockData.HeaderType {
	case header:
		caster, err := process.NewHeaderV1Caster(outportBlock)
		if err != nil {
			return nil, fmt.Errorf("failed to create new header v1 caster: %w", err)
		}

		return caster, nil
	case headerV2:
		caster, err := process.NewHeaderV2Caster(outportBlock)
		if err != nil {
			return nil, fmt.Errorf("failed to create new header v1 caster: %w", err)
		}

		return caster, nil
	case metaBlock:
		caster, err := process.NewMetaBlockCaster(outportBlock)
		if err != nil {
			return nil, fmt.Errorf("failed to create new header v1 caster: %w", err)
		}

		return caster, nil

	default:
		return nil, fmt.Errorf("unknown header type: %s", outportBlock.BlockData.HeaderType)
	}
}

func GetOutportCasterFunc(outportBlock *outport.OutportBlock) func(block *outport.OutportBlock) (proto.Message, error) {
	switch outportBlock.BlockData.HeaderType {

	case header:
		return func(block *outport.OutportBlock) (proto.Message, error) {
			return process.CastHeaderV1(block)
		}
	}

	return nil
}
