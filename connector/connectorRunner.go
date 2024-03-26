package connector

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/factory"
)

var log = logger.GetOrCreate("connectorRunner")

// ErrNilConfig signals that a nil config structure was provided
var ErrNilConfig = errors.New("nil configs provided")

type connectorRunner struct {
	config       *config.Config
	importDBMode bool
}

// NewConnectorRunner will create a new connector runner instance
func NewConnectorRunner(cfg *config.Config, importDBMode bool) (*connectorRunner, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	return &connectorRunner{
		config:       cfg,
		importDBMode: importDBMode,
	}, nil
}

// Start will trigger connector service
func (cr *connectorRunner) Start() error {
	storer, err := factory.CreateStorer(*cr.config, cr.importDBMode)
	if err != nil {
		return err
	}

	dataProcessor, err := factory.CreateDataProcessor(*cr.config, storer)
	if err != nil {
		return fmt.Errorf("cannot create ws firehose data processor, error: %w", err)
	}

	wsClient, err := factory.CreateWSConnector(cr.config.WebSocket, dataProcessor)
	if err != nil {
		return fmt.Errorf("cannot create ws firehose connector, error: %w", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("starting ws client...")

	<-interrupt

	log.Info("application closing, calling Close on all subcomponents...")

	err = storer.Close()
	if err != nil {
		return err
	}

	err = wsClient.Close()
	if err != nil {
		return err
	}

	return err
}
