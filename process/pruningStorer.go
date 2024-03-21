package process

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/multiversx/mx-chain-storage-go/leveldb"
	"github.com/multiversx/mx-chain-storage-go/storageUnit"
	"github.com/multiversx/mx-chain-storage-go/types"
	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
)

type pruningStorer struct {
	activePersisters []types.Persister
	persistersMut    sync.RWMutex
	cacher           types.Cacher

	dbConf              config.DBConfig
	numPersistersToKeep int
}

func NewPruningStorer(cfg config.DBConfig, cacher types.Cacher, numPersistersToKeep int) (*pruningStorer, error) {
	if cfg.FilePath == "" {
		return nil, ErrInvalidFilePath
	}
	if numPersistersToKeep < 1 {
		return nil, ErrInvalidNumberOfPersisters
	}

	ps := &pruningStorer{
		dbConf:              cfg,
		cacher:              cacher,
		numPersistersToKeep: numPersistersToKeep,
	}

	err := ps.initPersisters()
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func (ps *pruningStorer) initPersisters() error {
	persisterPaths, err := ps.getPersisterPaths()
	if err != nil {
		return err
	}

	persistersSlice := make([]types.Persister, 0)

	numPersistersToInit := ps.numPersistersToKeep
	if len(persisterPaths) < numPersistersToInit {
		numPersistersToInit = len(persisterPaths)
	}

	for i := 0; i < numPersistersToInit; i++ {
		persister, err := ps.createPersister(persisterPaths[i])
		if err != nil {
			return err
		}

		persistersSlice = append(persistersSlice, persister)
	}

	ps.activePersisters = persistersSlice

	return nil
}

func (ps *pruningStorer) createPersister(path string) (types.Persister, error) {
	var dbType = storageUnit.DBType(ps.dbConf.Type)
	switch dbType {
	case storageUnit.LvlDB:
		return leveldb.NewDB(path, ps.dbConf.BatchDelaySeconds, ps.dbConf.MaxBatchSize, ps.dbConf.MaxOpenFiles)
	case storageUnit.LvlDBSerial:
		return leveldb.NewSerialDB(path, ps.dbConf.BatchDelaySeconds, ps.dbConf.MaxBatchSize, ps.dbConf.MaxOpenFiles)
	default:
		return nil, ErrNotSupportedDBType
	}
}

func (ps *pruningStorer) getPersisterPaths() ([]string, error) {
	basePath := ps.dbConf.FilePath

	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	persistersPaths := make([]string, 0)

	for _, file := range files {
		if file.IsDir() {
			path := filepath.Join(basePath, file.Name())
			persistersPaths = append(persistersPaths, path)
		}
	}

	// set latest dir to first position
	persistersPaths = reverseSlice(persistersPaths)

	if len(persistersPaths) == 0 {
		// add "0" db dir if not already existing db
		path := filepath.Join(basePath, "0")
		persistersPaths = append(persistersPaths, path)
	}

	return persistersPaths, nil
}

func reverseSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func (ps *pruningStorer) Get(key []byte) ([]byte, error) {
	v, ok := ps.cacher.Get(key)
	if ok {
		return v.([]byte), nil
	}

	ps.persistersMut.RLock()
	defer ps.persistersMut.RUnlock()

	for idx := 0; idx < len(ps.activePersisters); idx++ {
		val, err := ps.activePersisters[idx].Get(key)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		// if found in persistence unit, add it to cache and return
		_ = ps.cacher.Put(key, val, len(val))

		return val, nil
	}

	return nil, fmt.Errorf("key %s not found", hex.EncodeToString(key))
}

func (ps *pruningStorer) Put(key, data []byte) error {
	ps.cacher.Put(key, data, len(data))

	ps.persistersMut.RLock()
	defer ps.persistersMut.RUnlock()

	// always put to last active persister
	persisterToUse := ps.activePersisters[0]

	err := persisterToUse.Put(key, data)
	if err != nil {
		ps.cacher.Remove(key)
		return err
	}

	return nil
}

func (ps *pruningStorer) DBChange(index uint64) error {
	return ps.updateActivePersisters(index)
}

func (ps *pruningStorer) updateActivePersisters(index uint64) error {
	ps.persistersMut.Lock()
	defer ps.persistersMut.Unlock()

	basePath := ps.dbConf.FilePath
	newPath := filepath.Join(basePath, fmt.Sprintf("%d", index))

	newPersister, err := ps.createPersister(newPath)
	if err != nil {
		return err
	}

	numPersistersToKeep := ps.numPersistersToKeep
	if len(ps.activePersisters) < numPersistersToKeep {
		numPersistersToKeep = len(ps.activePersisters)
	}
	numPersistersToKeep--

	if numPersistersToKeep < 1 {
		numPersistersToKeep = 1
	}

	newActivePersisters := make([]types.Persister, 0)
	for i := 0; i < numPersistersToKeep; i++ {
		newActivePersisters = append(newActivePersisters, ps.activePersisters[i])
	}

	newActivePersisters = append([]types.Persister{newPersister}, newActivePersisters...)

	ps.activePersisters = newActivePersisters

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (ps *pruningStorer) IsInterfaceNil() bool {
	return ps == nil
}
