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
	resetCheckpoints bool
}

// NewConnectorRunner will create a new connector runner instance
func NewConnectorRunner(cfg *config.Config, dbMode string, enableGrpcServer bool, resetCheckpoints bool) (*connectorRunner, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}

	return &connectorRunner{
		config:           cfg,
		dbMode:           common.DBMode(dbMode),
		enableGrpcServer: enableGrpcServer,
		resetCheckpoints: resetCheckpoints,
	}, nil
}

// Run will trigger connector service
func (cr *connectorRunner) Run() error {
	gogoProtoMarshaller := &marshal.GogoProtoMarshalizer{}
	protoMarshaller := &process.ProtoMarshaller{}

	firstCommitableBlocks, err := common.ConvertFirstCommitableBlocks(cr.config.DataPool.FirstCommitableBlocks)
	if err != nil {
		return err
	}

	outportBlockConverter, err := process.NewOutportBlockConverter(gogoProtoMarshaller, protoMarshaller)
	if err != nil {
		return err
	}

	blockContainer, err := factory.CreateBlockContainer()
	if err != nil {
		return err
	}

	blocksStorer, err := factory.CreateStorer(*cr.config, cr.dbMode)
	if err != nil {
		return err
	}

	argsBlocksPool := process.DataPoolArgs{
		Storer:                blocksStorer,
		Marshaller:            protoMarshaller,
		MaxDelta:              cr.config.DataPool.MaxDelta,
		CleanupInterval:       cr.config.DataPool.PruningWindow,
		FirstCommitableBlocks: common.DeepCopyNoncesMap(firstCommitableBlocks),
		ResetCheckpoints:      cr.resetCheckpoints,
	}
	dataPool, err := process.NewDataPool(argsBlocksPool)
	if err != nil {
		return err
	}

	blocksPool, err := process.NewBlocksPool(
		dataPool,
		protoMarshaller,
	)
	if err != nil {
		return err
	}

	dataAggregator, err := process.NewDataAggregator(blocksPool)
	if err != nil {
		return err
	}

	hyperBlockPublisher, err := factory.CreateHyperBlockPublisher(cr.config, cr.enableGrpcServer, blockContainer, blocksPool, dataAggregator)
	if err != nil {
		return fmt.Errorf("cannot create hyper block publisher: %w", err)
	}

	publisherHandlerArgs := process.PublisherHandlerArgs{
		Handler:                     hyperBlockPublisher,
		OutportBlocksPool:           blocksPool,
		DataAggregator:              dataAggregator,
		RetryDurationInMilliseconds: cr.config.Publisher.RetryDurationInMiliseconds,
		Marshalizer:                 protoMarshaller,
		FirstCommitableBlocks:       common.DeepCopyNoncesMap(firstCommitableBlocks),
	}

	publisherHandler, err := process.NewPublisherHandler(publisherHandlerArgs)
	if err != nil {
		return fmt.Errorf("cannot create common publisher: %w", err)
	}

	dataProcessor, err := process.NewDataProcessor(
		publisherHandler,
		gogoProtoMarshaller,
		blocksPool,
		outportBlockConverter,
		common.DeepCopyNoncesMap(firstCommitableBlocks),
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

	err = publisherHandler.Close()
	if err != nil {
		log.Error(err.Error())
	}

	err = wsClient.Close()
	if err != nil {
		log.Error(err.Error())
	}

	err = blocksPool.Close()
	if err != nil {
		log.Error(err.Error())
	}

	return err
}
