package factory

import (
	"testing"

	"github.com/multiversx/mx-chain-communication-go/websocket"
	"github.com/multiversx/mx-chain-communication-go/websocket/data"
	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/stretchr/testify/require"
)

func createConfig() config.WebSocketConfig {
	return config.WebSocketConfig{
		URL:                "localhost",
		MarshallerType:     "json",
		RetryDurationInSec: 1,
		WithAcknowledge:    false,
		BlockingAckOnError: false,
		Mode:               data.ModeClient,
	}
}

func TestCreateWSConnector(t *testing.T) {
	t.Parallel()

	t.Run("invalid marshaller, should return error", func(t *testing.T) {
		t.Parallel()

		cfg := createConfig()
		cfg.MarshallerType = "invalid"
		ws, err := CreateWSConnector(cfg, websocket.NewNilPayloadHandler())
		require.NotNil(t, err)
		require.Nil(t, ws)
	})

	t.Run("cannot create ws host, should return error", func(t *testing.T) {
		t.Parallel()

		cfg := createConfig()
		cfg.RetryDurationInSec = 0
		ws, err := CreateWSConnector(cfg, websocket.NewNilPayloadHandler())
		require.NotNil(t, err)
		require.Nil(t, ws)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		cfg := createConfig()
		ws, err := CreateWSConnector(cfg, websocket.NewNilPayloadHandler())
		require.Nil(t, err)
		require.NotNil(t, ws)

		_ = ws.Close()
	})
}
