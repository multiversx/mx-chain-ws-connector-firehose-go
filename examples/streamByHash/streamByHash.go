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

	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(16045218)),
	}

	cc, err := grpc.Dial(":8000", opts...)
	if err != nil {
		panic(err)
	}
	client := data.NewHyperOutportBlockServiceClient(cc)

	stream, err := client.HyperOutportBlockStreamByHash(context.Background(), &data.BlockHashStreamRequest{
		Hash:            "97227073a4f6702eefe6d4925a49bc6786947d084fb60250b39edabcba9ef6de",
		PollingInterval: durationpb.New(time.Second),
	})
	if err != nil {
		panic(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

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
