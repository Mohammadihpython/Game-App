package main

import (
	presenceClient "GameApp/adaptor/presence"
	"GameApp/param"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//var grpcClient = &grpc.ClientConn{}
	conn, err := grpc.NewClient(":8070", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := presenceClient.New(conn)
	res, err := client.GetPresence(context.Background(), param.GetPresenceRequest{UserIDs: []uint{1, 2, 3}})
	if err != nil {
		panic(err)
	}
	for _, item := range res.Items {
		fmt.Println(item)
	}
}
