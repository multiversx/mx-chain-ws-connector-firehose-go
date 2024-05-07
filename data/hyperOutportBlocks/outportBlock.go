package hyperOutportBlocks

import "fmt"

// GetHeaderRound will return header round
func (x *ShardOutportBlock) GetHeaderRound() (uint64, error) {
	if x.BlockData == nil {
		return 0, fmt.Errorf("shard outport block: nil block data")
	}
	if x.BlockData.Header == nil {
		return 0, fmt.Errorf("shard outport block: nil header")
	}

	return x.BlockData.Header.Round, nil
}

// GetHeaderRound will return header round
func (x *MetaOutportBlock) GetHeaderRound() (uint64, error) {
	if x.BlockData == nil {
		return 0, fmt.Errorf("meta outport block: nil block data")
	}
	if x.BlockData.Header == nil {
		return 0, fmt.Errorf("meta outport block: nil header")
	}

	return x.BlockData.Header.Round, nil
}
