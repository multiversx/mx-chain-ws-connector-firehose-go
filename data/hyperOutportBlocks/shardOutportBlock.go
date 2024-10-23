package hyperOutportBlocks

import "fmt"

// GetHeaderNonce will return header nonce
func (x *ShardOutportBlock) GetHeaderNonce() (uint64, error) {
	if x.BlockData == nil {
		return 0, fmt.Errorf("shard outport block: nil block data")
	}
	if x.BlockData.Header == nil {
		return 0, fmt.Errorf("shard outport block: nil header")
	}

	return x.BlockData.Header.Nonce, nil
}

// GetHeaderNonce will return header nonce
func (y *ShardOutportBlockV2) GetHeaderNonce() (uint64, error) {
	if y.BlockData == nil {
		return 0, fmt.Errorf("shard outport block: nil block data")
	}
	if y.BlockData.Header == nil {
		return 0, fmt.Errorf("shard outport block: nil header")
	}

	return y.BlockData.Header.Nonce, nil
}
