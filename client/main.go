package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "github.com/multiversx/mx-chain-ws-connector-template-go/api/hyperOutportBlocks"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("localhost:8000", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := api.NewHyperOutportBlockServiceClient(conn)

	h, err := client.GetHyperOutportBlockByHash(context.TODO(), &api.BlockHashRequest{Hash: "290527fcdd2d3f565f609a82f870be14d40f5aff3de73e87d72be2520c667c98"})
	if err != nil {
		panic(err)
	}

	fmt.Println(h)
}
