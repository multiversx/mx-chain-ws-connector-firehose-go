package factory

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/stretchr/testify/require"
)

func TestGetOutportCaster(t *testing.T) {
	t.Run("get heaverV1 caster", func(t *testing.T) {

		ob := &outport.OutportBlock{BlockData: &outport.BlockData{HeaderType: header}}

		_, err := GetOutportCaster(ob)
		require.NoError(t, err, "failed to retrieve caster")

	})

	t.Run("get heaverV2 caster", func(t *testing.T) {

		ob := &outport.OutportBlock{BlockData: &outport.BlockData{HeaderType: headerV2}}

		_, err := GetOutportCaster(ob)
		require.NoError(t, err, "failed to retrieve caster")

	})

	t.Run("get metaBlock caster", func(t *testing.T) {

		ob := &outport.OutportBlock{BlockData: &outport.BlockData{HeaderType: metaBlock}}

		_, err := GetOutportCaster(ob)
		require.NoError(t, err, "failed to retrieve caster")

	})

	t.Run("get unknown caster", func(t *testing.T) {

		ob := &outport.OutportBlock{BlockData: &outport.BlockData{HeaderType: "Unknown"}}

		_, err := GetOutportCaster(ob)
		require.Equal(t, errors.New("unknown header type: Unknown"), err)
	})

}
