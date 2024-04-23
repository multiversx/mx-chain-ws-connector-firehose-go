package main

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

func main() {
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(16045217)))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := data.NewHyperOutportBlockServiceClient(conn)
	stream, err := client.HyperOutportBlockStream(context.TODO(), &empty.Empty{})
	if err != nil {
		panic(err)
	}

	for {
		recv, err := stream.Recv()
		if err != nil {
			panic(err)
		}

		fmt.Println(recv)
	}
}
