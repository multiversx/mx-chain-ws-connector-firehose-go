package connector

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/factory"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process/dataPool"
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
	protoMarshaller := &marshal.GogoProtoMarshalizer{}

	outportBlocksStorer, err := factory.CreateStorer(*cr.config, cr.importDBMode)
	if err != nil {
		return err
	}

	// TODO: separate config for hyper outport blocks storer
	hyperOutportBlocksStorer, err := factory.CreateStorer(*cr.config, cr.importDBMode)
	if err != nil {
		return err
	}

	blockContainer, err := factory.CreateBlockContainer()
	if err != nil {
		return err
	}

	outportBlockDataPool, err := dataPool.NewBlocksPool(outportBlocksStorer, protoMarshaller, cr.config.DataPool.NumberOfShards, cr.config.DataPool.MaxDelta, cr.config.DataPool.PruningWindow)
	if err != nil {
		return err
	}

	// TODO: add separate config section for hyper blocks grpc data pool
	hyperOutportBlockDataPool, err := dataPool.NewBlocksPool(hyperOutportBlocksStorer, protoMarshaller, cr.config.DataPool.NumberOfShards, cr.config.DataPool.MaxDelta, cr.config.DataPool.PruningWindow)
	if err != nil {
		return err
	}

	outportBlocksPool, err := dataPool.NewOutportBlocksPool(outportBlockDataPool, protoMarshaller)
	if err != nil {
		return err
	}

	dataAggregator, err := process.NewDataAggregator(outportBlocksPool)
	if err != nil {
		return err
	}

	// TODO: move variable to config
	isGrpcServerActivated := false

	publisher, err := factory.CreatePublisher(isGrpcServerActivated, blockContainer, protoMarshaller, hyperOutportBlockDataPool)
	if err != nil {
		return err
	}

	dataProcessor, err := process.NewDataProcessor(publisher, protoMarshaller, outportBlocksPool, dataAggregator, blockContainer)
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

	err = outportBlockDataPool.Close()
	if err != nil {
		return err
	}

	err = hyperOutportBlockDataPool.Close()
	if err != nil {
		return err
	}

	err = wsClient.Close()
	if err != nil {
		return err
	}

	return err
}
