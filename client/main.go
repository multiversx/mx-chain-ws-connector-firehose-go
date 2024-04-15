package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("localhost:8000", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := data.NewHyperOutportBlockServiceClient(conn)

	h, err := client.GetHyperOutportBlockByHash(context.TODO(), &data.BlockHashRequest{Hash: "94a3a0e610350452da40bed3e1dc83b82b135e59835b02f20c46a4ba131ffb50"})
	if err != nil {
		panic(err)
	}

	fmt.Println(h)
}
