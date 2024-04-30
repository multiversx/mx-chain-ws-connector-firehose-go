package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"

	data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(16045218)),
	}

	// set up the client connection. Must have the connector running on the specified port.
	cc, err := grpc.Dial(":8000", opts...)
	if err != nil {
		panic(err)
	}
	defer cc.Close()

	// create the client
	client := data.NewHyperOutportBlockServiceClient(cc)

	// initiate streaming.
	stream, err := client.HyperOutportBlockStreamByNonce(context.Background(), &data.BlockNonceStreamRequest{
		Nonce:           0,
		PollingInterval: durationpb.New(time.Second),
	})
	if err != nil {
		panic(err)
	}

	// listen on the interrupt and terminate signals.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// try receiving blocks on the provided stream.
	go func() {
		for {
			recv, err := stream.Recv()
			if err != nil {
				panic(err)
			}

			fmt.Println(recv.MetaOutportBlock.BlockData.Header.Nonce)
		}
	}()

	<-interrupt
	stream.CloseSend()
}
