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
	config *config.Config
	dbMode string
}

// NewConnectorRunner will create a new connector runner instance
func NewConnectorRunner(cfg *config.Config, dbMode string) (*connectorRunner, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	return &connectorRunner{
		config: cfg,
		dbMode: dbMode,
	}, nil
}

// Run will trigger connector service
func (cr *connectorRunner) Run() error {
	// TODO: move variable to config
	isGrpcServerActivated := false

	protoMarshaller := &marshal.GogoProtoMarshalizer{}

	blockContainer, err := factory.CreateBlockContainer()
	if err != nil {
		return err
	}

	outportBlockDataPool, err := factory.CreateBlocksPool(*cr.config, cr.dbMode, protoMarshaller)
	if err != nil {
		return err
	}

	// TODO: add separate config section for hyper blocks grpc data pool
	hyperOutportBlockPool, err := factory.CreateHyperBlocksPool(isGrpcServerActivated, *cr.config, cr.dbMode, protoMarshaller)
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

	publisher, err := factory.CreatePublisher(isGrpcServerActivated, blockContainer, protoMarshaller, hyperOutportBlockPool)
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

	err = dataProcessor.Close()
	if err != nil {
		return err
	}

	err = wsClient.Close()
	if err != nil {
		return err
	}

	err = outportBlocksPool.Close()
	if err != nil {
		return err
	}

	err = hyperOutportBlockPool.Close()
	if err != nil {
		return err
	}

	return err
}
