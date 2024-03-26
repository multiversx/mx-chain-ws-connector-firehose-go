package process

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core/check"
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

// NewPruningStorer will create a new instance of pruning storer
func NewPruningStorer(cfg config.DBConfig, cacher types.Cacher, numPersistersToKeep int) (*pruningStorer, error) {
	if check.IfNil(cacher) {
		return nil, ErrNilCacher
	}
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

	err := createDirIfNotExisting(basePath)
	if err != nil {
		return nil, err
	}

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

func createDirIfNotExisting(path string) error {
	err := os.Mkdir(path, os.ModePerm)
	if err == nil {
		return nil
	}

	if !os.IsExist(err) {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("path exists but is not a directory")
	}

	return nil
}

func reverseSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

// Get will get data from cacher or storer
func (ps *pruningStorer) Get(key []byte) ([]byte, error) {
	v, ok := ps.cacher.Get(key)
	if ok {
		return v.([]byte), nil
	}

	return ps.getFromPersister(key)
}

func (ps *pruningStorer) getFromPersister(key []byte) ([]byte, error) {
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

// Put will try to put data to cacher and storer
func (ps *pruningStorer) Put(key, data []byte) error {
	ps.cacher.Put(key, data, len(data))

	return nil
}

func (ps *pruningStorer) putInPersister(key, data []byte) error {
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

// Dump will dump cached data to persister
func (ps *pruningStorer) Dump() error {
	return ps.dumpDataToPersister()
}

func (ps *pruningStorer) dumpDataToPersister() error {
	cacherKeys := ps.cacher.Keys()

	for _, key := range cacherKeys {
		v, ok := ps.cacher.Get(key)
		if !ok {
			log.Warn("failed to get key from cache", "key", hex.EncodeToString(key))
		}

		data, ok := v.([]byte)
		if !ok {
			log.Warn("failed to convert to byte array")
			continue
		}

		err := ps.putInPersister(key, data)
		if err != nil {
			log.Error("failed to dump data to persister", "key", hex.EncodeToString(key))
		}
	}

	return nil
}

// Prune will update active peristers list and will delete old persisters
func (ps *pruningStorer) Prune(index uint64) error {
	err := ps.updateActivePersisters(index)
	if err != nil {
		return err
	}

	return ps.cleanupOldPersisters()
}

func (ps *pruningStorer) cleanupOldPersisters() error {
	ps.persistersMut.Lock()
	defer ps.persistersMut.Unlock()

	numPersistersToKeep := ps.numPersistersToKeep
	if len(ps.activePersisters) > ps.numPersistersToKeep {
		numPersistersToKeep = len(ps.activePersisters)
	}

	persistersPaths, err := ps.getPersisterPaths()
	if err != nil {
		return err
	}

	if len(persistersPaths) <= numPersistersToKeep {
		// should not remove old persisters
		return nil
	}

	persistersPathsToRemove := persistersPaths[numPersistersToKeep:]

	for _, path := range persistersPathsToRemove {
		err := os.RemoveAll(path)
		if err != nil {
			log.Warn("failed to remove db dir", "path", path)
			continue
		}
	}

	return nil
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

	// numTmpPersistersToKeep := ps.getTmpNumPersistersToKeep()
	numTmpPersistersToKeep := ps.getTmpNumPersistersToKeep()

	newActivePersisters := make([]types.Persister, 0)
	for i := 0; i < numTmpPersistersToKeep; i++ {
		newActivePersisters = append(newActivePersisters, ps.activePersisters[i])
	}

	newActivePersisters = append([]types.Persister{newPersister}, newActivePersisters...)

	if len(ps.activePersisters) >= ps.numPersistersToKeep {
		inactivePersister := ps.activePersisters[len(ps.activePersisters)-1]
		err = inactivePersister.Close()
		if err != nil {
			return err
		}
	}

	ps.activePersisters = newActivePersisters

	return nil
}

func (ps *pruningStorer) getTmpNumPersistersToKeep() int {
	numPersistersToKeep := ps.numPersistersToKeep
	if len(ps.activePersisters) < numPersistersToKeep {
		numPersistersToKeep = len(ps.activePersisters)
	}
	numPersistersToKeep--

	if numPersistersToKeep < 1 {
		numPersistersToKeep = 1
	}

	return numPersistersToKeep
}

// IsInterfaceNil returns true if there is no value under the interface
func (ps *pruningStorer) IsInterfaceNil() bool {
	return ps == nil
}
