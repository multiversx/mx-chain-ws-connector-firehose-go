package process

import (
	"encoding/hex"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-storage-go/types"
)

type importDBStorer struct {
	cacher types.Cacher
}

// NewImportDBStorer creates a new import db storer instancer
func NewImportDBStorer(cacher types.Cacher) (*importDBStorer, error) {
	if check.IfNil(cacher) {
		return nil, ErrNilCacher
	}

	return &importDBStorer{
		cacher: cacher,
	}, nil
}

// Get will get value from cache
func (is *importDBStorer) Get(key []byte) ([]byte, error) {
	v, ok := is.cacher.Get(key)
	if ok {
		return v.([]byte), nil
	}

	return nil, fmt.Errorf("key %s not found", hex.EncodeToString(key))
}

// Put will put data into cacher
func (is *importDBStorer) Put(key []byte, data []byte) error {
	is.cacher.Put(key, data, len(data))
	return nil
}

// Prune returns nil
func (is *importDBStorer) Prune(_ uint64) error {
	return nil
}

// Dump returns nil
func (is *importDBStorer) Dump() error {
	return nil
}

// IsInterfaceNil returns nil if there is no value under the interface
func (is *importDBStorer) IsInterfaceNil() bool {
	return is == nil
}
