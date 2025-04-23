package main

import (
	"GameApp/contract/golang/presence"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//var grpcClient = &grpc.ClientConn{}
	conn, err := grpc.Dial(":8070", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := presence.NewPresenceServiceClient(conn)
	res, err := client.GetPresence(context.Background(), &presence.GetPresenceRequest{UserIds: []uint64{1, 2, 3}})
	if err != nil {
		panic(err)
	}
	for _, item := range res.Items {
		fmt.Println(item)
	}
}
