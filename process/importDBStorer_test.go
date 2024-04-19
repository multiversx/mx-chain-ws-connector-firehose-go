package process_test

import (
	"testing"

	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
	"github.com/stretchr/testify/require"
)

func TestNewImportDBStorer(t *testing.T) {
	t.Parallel()

	t.Run("nil cacher", func(t *testing.T) {
		t.Parallel()

		is, err := process.NewImportDBStorer(nil)
		require.Nil(t, is)
		require.Equal(t, process.ErrNilCacher, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		is, err := process.NewImportDBStorer(testscommon.NewCacherMock())
		require.Nil(t, err)
		require.False(t, is.IsInterfaceNil())

		key := []byte("key1")
		val := []byte("val1")

		err = is.Put(key, val)
		require.Nil(t, err)

		ret, err := is.Get(key)
		require.Nil(t, err)
		require.Equal(t, val, ret)

		ret, err = is.Get([]byte("not existing"))
		require.Nil(t, ret)
		require.Error(t, err)

		require.Nil(t, is.Prune(2))
	})
}
