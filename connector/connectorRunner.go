package connector

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"

	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/factory"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
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
	protoMarshaller := &marshal.GogoProtoMarshalizer{}

	blockContainer, err := factory.CreateBlockContainer()
	if err != nil {
		return err
	}

	blocksStorer, err := factory.CreateStorer(*cr.config, cr.dbMode)
	if err != nil {
		return err
	}

	outportBlocksPool, err := process.NewBlocksPool(blocksStorer, protoMarshaller, cr.config.DataPool.MaxDelta, cr.config.DataPool.PruningWindow)
	if err != nil {
		return err
	}

	dataAggregator, err := process.NewDataAggregator(outportBlocksPool)
	if err != nil {
		return err
	}

	var (
		server *factory.GRPCServer
		writer process.Writer
	)
	if cr.config.GRPC.Enable {
		writer = &fakeWriter{}
		handler, err := process.NewGRPCBlocksHandler(outportBlocksPool, dataAggregator)
		if err != nil {
			return fmt.Errorf("couldn't create grpc blocks handler, error: %w", err)
		}
		server = factory.NewServer(cr.config.GRPC, handler)

		go func() {
			if err := server.Start(); err != nil {
				log.Error("couldn't start grpc server", "error", err)
			}
		}()
	} else {
		writer = os.Stdout
	}

	publisher, err := process.NewFirehosePublisher(
		writer,
		blockContainer,
		&process.ProtoMarshaller{},
	)
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

	err = outportBlocksPool.Close()
	if err != nil {
		log.Error(err.Error())
	}

	err = wsClient.Close()
	if err != nil {
		log.Error(err.Error())
	}

	if server != nil {
		server.Stop()
	}

	return err
}

type fakeWriter struct {
	err      error
	duration time.Duration
}

func (f *fakeWriter) Write(p []byte) (int, error) {
	time.Sleep(f.duration)
	if f.err != nil {
		return 0, f.err
	}

	return len(p), nil
}

func (f *fakeWriter) Close() error {
	return nil
}
