package process

import (
	"math/big"

	"github.com/multiversx/mx-chain-storage-go/types"
)

const (
	// FirehosePrefix -
	FirehosePrefix = firehosePrefix

	// BlockPrefix -
	BlockPrefix = blockPrefix
)

// GetActivePersisters -
func (ps *pruningStorer) GetActivePersisters() []*persisterData {
	ps.persistersMut.RLock()
	defer ps.persistersMut.RUnlock()

	return ps.activePersisters
}

// GetActivePersister -
func (ps *pruningStorer) GetActivePersister(index int) types.Persister {
	ps.persistersMut.RLock()
	defer ps.persistersMut.RUnlock()

	return ps.activePersisters[index].persister
}

// GetPersisterPaths -
func (ps *pruningStorer) GetPersisterPaths() ([]string, error) {
	return ps.getPersisterPaths()
}

// CastBigInt -
func (o *outportBlockConverter) CastBigInt(i *big.Int) ([]byte, error) {
	return o.castBigInt(i)
}
