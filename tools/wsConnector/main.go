package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	wsData "github.com/multiversx/mx-chain-communication-go/websocket/data"
	wsFactory "github.com/multiversx/mx-chain-communication-go/websocket/factory"
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testdata"
)

// tool to simulate observer outport flow
// 		component (separate ws client) that triggers a goroutine, that push events at an interval
// 		testing different scenarios: shards + meta

func main() {
	marshaller := &marshal.GogoProtoMarshalizer{}

	publisher0, err := NewWSPublisher(marshaller, 100)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = publisher0.PushEvents(0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("pushed shard 0 event")

	time.Sleep(time.Second * 4)

	publisherM, err := NewWSPublisher(marshaller, 1000)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = publisherM.PushEvents(core.MetachainShardId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("pushed meta event")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interrupt
	fmt.Println("closing app at user's signal")

	err = publisher0.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = publisherM.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

type blockDataHandler interface {
	OutportMetaBlockV1() *outport.OutportBlock
	OutportShardBlockV1() *outport.OutportBlock
}

type wsPublisher struct {
	wsConnector  *wsObsClient
	blockData    blockDataHandler
	duration     time.Duration
	hashIncr     uint32
	hashIncrMeta uint32
	marshaller   marshal.Marshalizer

	cancelFunc func()
	closeChan  chan struct{}
	mutState   sync.RWMutex
}

func NewWSPublisher(marshaller marshal.Marshalizer, durationInMs int) (*wsPublisher, error) {
	if check.IfNil(marshaller) {
		return nil, fmt.Errorf("nil marshaller")
	}
	if durationInMs <= 0 {
		return nil, fmt.Errorf("invalid duration in milliseconds")
	}

	wsClient, err := newWSObsClient(marshaller)
	if err != nil {
		return nil, err
	}

	blockData, err := testdata.NewBlockData(marshaller)
	if err != nil {
		return nil, err
	}

	duration := time.Millisecond * time.Duration(durationInMs)

	return &wsPublisher{
		wsConnector:  wsClient,
		blockData:    blockData,
		duration:     duration,
		closeChan:    make(chan struct{}),
		hashIncrMeta: 2,
		marshaller:   &marshal.GogoProtoMarshalizer{},
	}, nil
}

func (wp *wsPublisher) PushEvents(shardID uint32) error {
	wp.mutState.Lock()
	defer wp.mutState.Unlock()

	var ctx context.Context
	ctx, wp.cancelFunc = context.WithCancel(context.Background())

	go wp.run(ctx, shardID)

	return nil
}

func (wp *wsPublisher) run(ctx context.Context, shardID uint32) {
	ticker := time.NewTicker(wp.duration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			err := wp.wsConnector.Close()
			if err != nil {
				fmt.Println("failed to close connector", "error", err.Error())
			}

			return
		case <-ticker.C:
			outportData := wp.getOutportData(shardID)
			wp.wsConnector.PushEventsRequest(outportData)
		}
	}
}

func (ws *wsPublisher) getOutportData(shardID uint32) *outport.OutportBlock {
	baseHeaderHash := []byte("headerHash")
	baseMetaHeaderHash := []byte("metaheaderHash")

	if shardID == core.MetachainShardId {
		ws.hashIncrMeta++
		outportBlock := ws.blockData.OutportMetaBlockV1()
		metaHeaderHash := append(baseMetaHeaderHash, []byte(fmt.Sprintf("%d", ws.hashIncrMeta))...)
		outportBlock.BlockData.HeaderHash = metaHeaderHash
		header := &block.MetaBlock{
			Round: uint64(ws.hashIncrMeta),
		}
		headerBytes, _ := ws.marshaller.Marshal(header)
		outportBlock.BlockData.HeaderBytes = headerBytes

		shardHeaderHash := append(baseHeaderHash, []byte(fmt.Sprintf("%d", ws.hashIncr+1))...)

		outportBlock.NotarizedHeadersHashes = []string{hex.EncodeToString(shardHeaderHash)}

		return outportBlock
	}

	ws.hashIncr++

	outportBlock := ws.blockData.OutportShardBlockV1()
	headerHash := append(baseHeaderHash, []byte(fmt.Sprintf("%d", ws.hashIncr))...)
	outportBlock.BlockData.HeaderHash = headerHash

	header := &block.Header{
		Round: uint64(ws.hashIncr),
	}
	headerBytes, _ := ws.marshaller.Marshal(header)
	outportBlock.BlockData.HeaderBytes = headerBytes

	return outportBlock
}

func (wp *wsPublisher) Close() error {
	wp.mutState.RLock()
	defer wp.mutState.RUnlock()

	if wp.cancelFunc != nil {
		wp.cancelFunc()
	}

	return nil
}

// senderHost defines the actions that a host sender should do
type senderHost interface {
	Send(payload []byte, topic string) error
	Close() error
	IsInterfaceNil() bool
}

type wsObsClient struct {
	marshaller marshal.Marshalizer
	senderHost senderHost
}

// newWSObsClient will create a new instance of observer websocket client
func newWSObsClient(marshaller marshal.Marshalizer) (*wsObsClient, error) {
	var log = logger.GetOrCreate("hostdriver")

	port := 22111
	wsHost, err := wsFactory.CreateWebSocketHost(wsFactory.ArgsWebSocketHost{
		WebSocketConfig: wsData.WebSocketConfig{
			URL:                     "localhost:" + fmt.Sprintf("%d", port),
			WithAcknowledge:         true,
			Mode:                    "client",
			RetryDurationInSec:      5,
			BlockingAckOnError:      false,
			AcknowledgeTimeoutInSec: 60,
			Version:                 1,
		},
		Marshaller: marshaller,
		Log:        log,
	})
	if err != nil {
		return nil, err
	}

	return &wsObsClient{
		marshaller: marshaller,
		senderHost: wsHost,
	}, nil
}

// SaveBlock will handle the saving of block
func (o *wsObsClient) PushEventsRequest(outportBlock *outport.OutportBlock) error {
	return o.handleAction(outportBlock, outport.TopicSaveBlock)
}

// RevertIndexedBlock will handle the action of reverting the indexed block
func (o *wsObsClient) RevertEventsRequest(blockData *outport.BlockData) error {
	return o.handleAction(blockData, outport.TopicRevertIndexedBlock)
}

// FinalizedBlock will handle the finalized block
func (o *wsObsClient) FinalizedEventsRequest(finalizedBlock *outport.FinalizedBlock) error {
	return o.handleAction(finalizedBlock, outport.TopicFinalizedBlock)
}

func (o *wsObsClient) handleAction(args interface{}, topic string) error {
	marshalledPayload, err := o.marshaller.Marshal(args)
	if err != nil {
		return fmt.Errorf("%w while marshaling block for topic %s", err, topic)
	}

	err = o.senderHost.Send(marshalledPayload, topic)
	if err != nil {
		return fmt.Errorf("%w while sending data on route for topic %s", err, topic)
	}

	return nil
}

// Close will close sender host connection
func (o *wsObsClient) Close() error {
	return o.senderHost.Close()
}
