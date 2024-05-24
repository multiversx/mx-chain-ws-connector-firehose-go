package hyperOutportBlocks

import "fmt"

// GetHeaderNonce will return header nonce
func (x *MetaOutportBlock) GetHeaderNonce() (uint64, error) {
	if x.BlockData == nil {
		return 0, fmt.Errorf("meta outport block: nil block data")
	}
	if x.BlockData.Header == nil {
		return 0, fmt.Errorf("meta outport block: nil header")
	}

	return x.BlockData.Header.Nonce, nil
}
