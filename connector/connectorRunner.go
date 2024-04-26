package connector

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/common"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/config"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/factory"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
)

var log = logger.GetOrCreate("connectorRunner")

// ErrNilConfig signals that a nil config structure was provided
var ErrNilConfig = errors.New("nil configs provided")

type connectorRunner struct {
	config           *config.Config
	dbMode           common.DBMode
	enableGrpcServer bool
}

// NewConnectorRunner will create a new connector runner instance
func NewConnectorRunner(cfg *config.Config, dbMode string, enableGrpcServer bool) (*connectorRunner, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	return &connectorRunner{
		config:           cfg,
		dbMode:           common.DBMode(dbMode),
		enableGrpcServer: enableGrpcServer,
	}, nil
}

// Run will trigger connector service
func (cr *connectorRunner) Run() error {
	gogoProtoMarshaller := &marshal.GogoProtoMarshalizer{}
	protoMarshaller := &process.ProtoMarshaller{}

	outportBlockConverter := process.NewOutportBlockConverter(gogoProtoMarshaller, protoMarshaller)

	blockContainer, err := factory.CreateBlockContainer()
	if err != nil {
		return err
	}

	blocksStorer, err := factory.CreateStorer(*cr.config, cr.dbMode)
	if err != nil {
		return err
	}

	baseBlocksPool, err := process.NewBlocksPool(
		blocksStorer,
		protoMarshaller,
		cr.config.DataPool.MaxDelta,
		cr.config.DataPool.PruningWindow,
		cr.config.DataPool.FirstCommitableBlock,
	)
	if err != nil {
		return err
	}

	outportBlocksPool, err := process.NewHyperOutportBlocksPool(
		baseBlocksPool,
		protoMarshaller,
	)
	if err != nil {
		return err
	}

	dataAggregator, err := process.NewDataAggregator(outportBlocksPool)
	if err != nil {
		return err
	}

	server, err := factory.CreateGRPCServer(cr.enableGrpcServer, cr.config.GRPC, outportBlocksPool, dataAggregator)
	if err != nil {
		return err
	}

	publisher, err := factory.CreatePublisher(cr.enableGrpcServer, blockContainer)
	if err != nil {
		return fmt.Errorf("cannot create publisher: %w", err)
	}

	dataProcessor, err := process.NewDataProcessor(
		publisher,
		gogoProtoMarshaller,
		outportBlocksPool,
		dataAggregator,
		outportBlockConverter,
		cr.config.DataPool.FirstCommitableBlock,
	)
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

	err = outportBlocksPool.Close()
	if err != nil {
		log.Error(err.Error())
	}

	err = wsClient.Close()
	if err != nil {
		log.Error(err.Error())
	}

	if server != nil {
		server.Close()
	}

	return err
}
