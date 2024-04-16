package process_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var isFullDBSync = true

func TestNewPruningStorer(t *testing.T) {
	t.Parallel()

	t.Run("invalid number of persisters to keep, should fail", func(t *testing.T) {
		t.Parallel()

		dbConfig := config.DBConfig{
			FilePath:          t.TempDir(),
			Type:              "LvlDBSerial",
			BatchDelaySeconds: 2,
			MaxBatchSize:      100,
			MaxOpenFiles:      10,
		}

		ps, err := process.NewPruningStorer(dbConfig, testscommon.NewCacherMock(), 0, isFullDBSync)
		require.Nil(t, ps)
		require.Equal(t, process.ErrInvalidNumberOfPersisters, err)
	})

	t.Run("invalid path, should fail", func(t *testing.T) {
		t.Parallel()

		dbConfig := config.DBConfig{
			FilePath:          "",
			Type:              "LvlDBSerial",
			BatchDelaySeconds: 2,
			MaxBatchSize:      100,
			MaxOpenFiles:      10,
		}

		ps, err := process.NewPruningStorer(dbConfig, testscommon.NewCacherMock(), 2, isFullDBSync)
		require.Nil(t, ps)
		require.Equal(t, process.ErrInvalidFilePath, err)
	})

	t.Run("not supported db type, should fail", func(t *testing.T) {
		t.Parallel()

		dbConfig := config.DBConfig{
			FilePath:          t.TempDir(),
			Type:              "not supported",
			BatchDelaySeconds: 2,
			MaxBatchSize:      100,
			MaxOpenFiles:      10,
		}

		ps, err := process.NewPruningStorer(dbConfig, testscommon.NewCacherMock(), 2, isFullDBSync)
		require.Nil(t, ps)
		require.Equal(t, process.ErrNotSupportedDBType, err)
	})

	t.Run("empty db path, should create default dir", func(t *testing.T) {
		t.Parallel()

		tmpPath := t.TempDir()

		dbConfig := config.DBConfig{
			FilePath:          tmpPath,
			Type:              "LvlDBSerial",
			BatchDelaySeconds: 2,
			MaxBatchSize:      100,
			MaxOpenFiles:      10,
		}

		ps, err := process.NewPruningStorer(dbConfig, testscommon.NewCacherMock(), 2, isFullDBSync)
		require.Nil(t, err)
		require.NotNil(t, ps)

		files, err := ioutil.ReadDir(tmpPath)
		require.Nil(t, err)
		require.False(t, ps.IsInterfaceNil())

		require.Equal(t, "0", files[0].Name())
	})

	t.Run("non empty db path, should create based on already existing dbs", func(t *testing.T) {
		t.Parallel()

		tmpPath := t.TempDir()

		_ = os.MkdirAll(filepath.Join(tmpPath, "3"), os.ModePerm)

		dbConfig := config.DBConfig{
			FilePath:          tmpPath,
			Type:              "LvlDBSerial",
			BatchDelaySeconds: 2,
			MaxBatchSize:      100,
			MaxOpenFiles:      10,
		}

		numPersistersToKeep := 2

		ps, err := process.NewPruningStorer(dbConfig, testscommon.NewCacherMock(), numPersistersToKeep, isFullDBSync)
		require.Nil(t, err)
		require.NotNil(t, ps)

		numAlreadyExistingDBs := 1
		require.Equal(t, numAlreadyExistingDBs, len(ps.GetActivePersisters()))
	})

	t.Run("non empty db path, should create based on numPersistersToKeep", func(t *testing.T) {
		t.Parallel()

		tmpPath := t.TempDir()

		_ = os.MkdirAll(filepath.Join(tmpPath, "1"), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(tmpPath, "2"), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(tmpPath, "3"), os.ModePerm)

		dbConfig := config.DBConfig{
			FilePath:          tmpPath,
			Type:              "LvlDB",
			BatchDelaySeconds: 2,
			MaxBatchSize:      100,
			MaxOpenFiles:      10,
		}

		numPersistersToKeep := 2

		ps, err := process.NewPruningStorer(dbConfig, testscommon.NewCacherMock(), numPersistersToKeep, isFullDBSync)
		require.Nil(t, err)
		require.NotNil(t, ps)

		require.Equal(t, numPersistersToKeep, len(ps.GetActivePersisters()))
	})
}

func TestPruningStorer_getPersisterPaths(t *testing.T) {
	t.Parallel()

	tmpPath := t.TempDir()

	_ = os.MkdirAll(filepath.Join(tmpPath, "1"), os.ModePerm)
	_ = os.MkdirAll(filepath.Join(tmpPath, "2"), os.ModePerm)
	_ = os.MkdirAll(filepath.Join(tmpPath, "4"), os.ModePerm)
	_ = os.MkdirAll(filepath.Join(tmpPath, "3"), os.ModePerm)

	dbConfig := config.DBConfig{
		FilePath:          tmpPath,
		Type:              "LvlDB",
		BatchDelaySeconds: 2,
		MaxBatchSize:      100,
		MaxOpenFiles:      10,
	}

	numPersistersToKeep := 2
	ps, err := process.NewPruningStorer(dbConfig, testscommon.NewCacherMock(), numPersistersToKeep, isFullDBSync)
	require.Nil(t, err)
	require.NotNil(t, ps)

	persistersPaths, err := ps.GetPersisterPaths()
	require.Nil(t, err)

	require.Equal(t, "4", filepath.Base(persistersPaths[0]))
	require.Equal(t, "3", filepath.Base(persistersPaths[1]))
	require.Equal(t, "2", filepath.Base(persistersPaths[2]))
}

func TestPruningStorer_Get(t *testing.T) {
	t.Parallel()

	t.Run("will get from cache", func(t *testing.T) {
		t.Parallel()

		tmpPath := t.TempDir()

		_ = os.MkdirAll(filepath.Join(tmpPath, "3"), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(tmpPath, "4"), os.ModePerm)

		dbConfig := config.DBConfig{
			FilePath:          tmpPath,
			Type:              "LvlDB",
			BatchDelaySeconds: 2,
			MaxBatchSize:      100,
			MaxOpenFiles:      10,
		}

		numPersistersToKeep := 2
		ps, err := process.NewPruningStorer(dbConfig, testscommon.NewCacherMock(), numPersistersToKeep, isFullDBSync)
		require.Nil(t, err)
		require.NotNil(t, ps)

		key := []byte("key4")
		value := []byte("value4")

		err = ps.Put(key, value)
		require.Nil(t, err)

		ret, err := ps.Get(key)
		require.Nil(t, err)

		require.Equal(t, value, ret)
	})

	t.Run("will get from persisters", func(t *testing.T) {
		t.Parallel()

		tmpPath := t.TempDir()

		_ = os.MkdirAll(filepath.Join(tmpPath, "3"), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(tmpPath, "4"), os.ModePerm)

		dbConfig := config.DBConfig{
			FilePath:          tmpPath,
			Type:              "LvlDB",
			BatchDelaySeconds: 2,
			MaxBatchSize:      100,
			MaxOpenFiles:      10,
		}

		cacher := testscommon.NewCacherMock()

		numPersistersToKeep := 2
		ps, err := process.NewPruningStorer(dbConfig, cacher, numPersistersToKeep, isFullDBSync)
		require.Nil(t, err)
		require.NotNil(t, ps)

		key := []byte("key4")
		value := []byte("value4")

		err = ps.Put(key, value)
		require.Nil(t, err)

		err = ps.Dump()
		require.Nil(t, err)

		ret, err := ps.Get(key)
		require.Nil(t, err)

		require.Equal(t, value, ret)
	})

	t.Run("should fail if not found in persisters", func(t *testing.T) {
		t.Parallel()

		tmpPath := t.TempDir()

		_ = os.MkdirAll(filepath.Join(tmpPath, "3"), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(tmpPath, "4"), os.ModePerm)

		dbConfig := config.DBConfig{
			FilePath:          tmpPath,
			Type:              "LvlDB",
			BatchDelaySeconds: 2,
			MaxBatchSize:      100,
			MaxOpenFiles:      10,
		}

		cacher := &testscommon.CacherStub{
			GetCalled: func(key []byte) (value interface{}, ok bool) {
				return nil, false
			},
			PutCalled: func(key []byte, value interface{}, sizeInBytes int) (evicted bool) {
				assert.Fail(t, "should have not been called")
				return true
			},
		}

		numPersistersToKeep := 2
		ps, err := process.NewPruningStorer(dbConfig, cacher, numPersistersToKeep, isFullDBSync)
		require.Nil(t, err)
		require.NotNil(t, ps)

		key := []byte("key4")

		ret, err := ps.Get(key)
		require.Nil(t, ret)
		require.Error(t, err)
	})
}

func TestPruningStorer_Put(t *testing.T) {
	t.Parallel()

	t.Run("should put to last active persister", func(t *testing.T) {
		t.Parallel()

		tmpPath := t.TempDir()

		_ = os.MkdirAll(filepath.Join(tmpPath, "3"), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(tmpPath, "4"), os.ModePerm)

		dbConfig := config.DBConfig{
			FilePath:          tmpPath,
			Type:              "LvlDB",
			BatchDelaySeconds: 2,
			MaxBatchSize:      100,
			MaxOpenFiles:      10,
		}

		numPersistersToKeep := 2
		ps, err := process.NewPruningStorer(dbConfig, testscommon.NewCacherMock(), numPersistersToKeep, isFullDBSync)
		require.Nil(t, err)
		require.NotNil(t, ps)

		persistersPaths, err := ps.GetPersisterPaths()
		require.Nil(t, err)

		require.Equal(t, "4", filepath.Base(persistersPaths[0]))
		require.Equal(t, "3", filepath.Base(persistersPaths[1]))

		key := []byte("key4")
		value := []byte("value4")

		persister4 := ps.GetActivePersister(0)
		_, err = persister4.Get(key)
		require.Error(t, err)

		persister3 := ps.GetActivePersister(1)
		_, err = persister3.Get(key)
		require.Error(t, err)

		err = ps.Put(key, value)
		require.Nil(t, err)

		err = ps.Dump()
		require.Nil(t, err)

		ret, err := persister4.Get(key)
		require.Nil(t, err)
		require.Equal(t, value, ret)
	})
}

func TestPruningStorer_Prune(t *testing.T) {
	t.Parallel()

	t.Run("with one already existing db", func(t *testing.T) {
		t.Parallel()

		tmpPath := t.TempDir()

		_ = os.MkdirAll(filepath.Join(tmpPath, "4"), os.ModePerm)

		dbConfig := config.DBConfig{
			FilePath:          tmpPath,
			Type:              "LvlDB",
			BatchDelaySeconds: 2,
			MaxBatchSize:      100,
			MaxOpenFiles:      10,
		}

		numPersistersToKeep := 2
		ps, err := process.NewPruningStorer(dbConfig, testscommon.NewCacherMock(), numPersistersToKeep, isFullDBSync)
		require.Nil(t, err)
		require.NotNil(t, ps)

		persistersPaths, err := ps.GetPersisterPaths()
		require.Nil(t, err)

		require.Equal(t, "4", filepath.Base(persistersPaths[0]))

		persister4 := ps.GetActivePersister(0)
		_ = persister4.Put([]byte("key4"), []byte("value4"))

		err = ps.Prune(5)
		require.Nil(t, err)

		// new persister will not have the data but the last one will have it
		persister5 := ps.GetActivePersister(0)
		key4Val, err := persister5.Get([]byte("key4"))
		require.Nil(t, key4Val)
		require.Error(t, err)

		persister4 = ps.GetActivePersister(1)
		key4Val, err = persister4.Get([]byte("key4"))
		require.Nil(t, err)
		require.Equal(t, []byte("value4"), key4Val)
	})

	t.Run("with multiple already existing dbs", func(t *testing.T) {
		t.Parallel()

		tmpPath := t.TempDir()

		_ = os.MkdirAll(filepath.Join(tmpPath, "1"), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(tmpPath, "2"), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(tmpPath, "4"), os.ModePerm)
		_ = os.MkdirAll(filepath.Join(tmpPath, "3"), os.ModePerm)

		dbConfig := config.DBConfig{
			FilePath:          tmpPath,
			Type:              "LvlDB",
			BatchDelaySeconds: 2,
			MaxBatchSize:      100,
			MaxOpenFiles:      10,
		}

		numPersistersToKeep := 2
		ps, err := process.NewPruningStorer(dbConfig, testscommon.NewCacherMock(), numPersistersToKeep, isFullDBSync)
		require.Nil(t, err)
		require.NotNil(t, ps)

		persistersPaths, err := ps.GetPersisterPaths()
		require.Nil(t, err)

		require.Equal(t, "4", filepath.Base(persistersPaths[0]))
		require.Equal(t, "3", filepath.Base(persistersPaths[1]))
		require.Equal(t, "2", filepath.Base(persistersPaths[2]))

		require.Equal(t, 4, len(persistersPaths))

		persister4 := ps.GetActivePersister(0)
		_ = persister4.Put([]byte("key4"), []byte("value4"))

		persister3 := ps.GetActivePersister(1)
		_ = persister3.Put([]byte("key3"), []byte("value3"))

		err = ps.Prune(5)
		require.Nil(t, err)

		persistersPaths, err = ps.GetPersisterPaths()
		require.Nil(t, err)

		require.Equal(t, numPersistersToKeep, len(persistersPaths))

		// new persister will not have the data but the last one will have it
		persister5 := ps.GetActivePersister(0)
		key4Val, err := persister5.Get([]byte("key4"))
		require.Nil(t, key4Val)
		require.Error(t, err)

		persister4 = ps.GetActivePersister(1)
		key4Val, err = persister4.Get([]byte("key4"))
		require.Nil(t, err)
		require.Equal(t, []byte("value4"), key4Val)
	})
}

func TestPruningStorer_Dump(t *testing.T) {
	t.Parallel()

	tmpPath := t.TempDir()

	_ = os.MkdirAll(filepath.Join(tmpPath, "4"), os.ModePerm)

	dbConfig := config.DBConfig{
		FilePath:          tmpPath,
		Type:              "LvlDB",
		BatchDelaySeconds: 2,
		MaxBatchSize:      100,
		MaxOpenFiles:      10,
	}

	numPersistersToKeep := 2
	ps, err := process.NewPruningStorer(dbConfig, testscommon.NewCacherMock(), numPersistersToKeep, !isFullDBSync)
	require.Nil(t, err)
	require.NotNil(t, ps)

	_ = ps.Put([]byte("key1"), []byte("value1"))
	_ = ps.Put([]byte("key2"), []byte("value2"))
	_ = ps.Put([]byte("key3"), []byte("value3"))
	_ = ps.Put([]byte("key4"), []byte("value4"))

	persister4 := ps.GetActivePersister(0)
	key4Val, err := persister4.Get([]byte("key4"))
	require.Nil(t, key4Val)
	require.Error(t, err)

	err = ps.Dump()
	require.Nil(t, err)

	persister4 = ps.GetActivePersister(0)
	key4Val, err = persister4.Get([]byte("key4"))
	require.NotNil(t, key4Val)
	require.Nil(t, err)
}
