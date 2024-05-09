package common

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/config"
)

// DBMode defines db mode type
type DBMode string

const (
	// FullPersisterDBMode defines the db mode where each event will be saved
	// to cache and persister
	FullPersisterDBMode DBMode = "full-persister"

	// ImportDBMode defines the db mode where the events will be saved only to cache
	ImportDBMode DBMode = "import-db"

	// OptimizedPersisterDBMode defines the db mode where the events will be saved to
	// cache, and they will be dumped to persister when necessary
	OptimizedPersisterDBMode DBMode = "optimized-persister"
)

// ConvertFirstCommitableBlocks will convert first commitable blocks map
func ConvertFirstCommitableBlocks(blocks []config.FirstCommitableBlock) (map[uint32]uint64, error) {
	newBlocks := make(map[uint32]uint64)

	for _, firstCommitableBlock := range blocks {
		shardID, err := core.ConvertShardIDToUint32(firstCommitableBlock.ShardID)
		if err != nil {
			return nil, err
		}

		newBlocks[shardID] = firstCommitableBlock.Nonce
	}

	return newBlocks, nil
}
