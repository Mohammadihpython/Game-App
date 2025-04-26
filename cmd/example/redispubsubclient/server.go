package main

import (
	"GameApp/adaptor/adaptor/redis"
	"GameApp/conf"
	"GameApp/contract/golang/matching"
	"GameApp/entity"
	"GameApp/pkg/slice"
	"context"
	"encoding/base64"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
)

func main() {
	cfg := conf.Load()
	redisAdaptor := redis.New(cfg.Redis)
	topic := "matching.users_matched"

	subsciber := redisAdaptor.Client().Subscribe(context.Background(), topic)

	for {
		msg, err := subsciber.ReceiveMessage(context.Background())
		if err != nil {
			log.Println(err)
		}
		payload, err := base64.StdEncoding.DecodeString(msg.Payload)
		if err != nil {
			log.Println(err)

		}
		pbMu := matching.MatchedUsers{}
		if err := proto.Unmarshal(payload, &pbMu); err != nil {
			log.Println(err)
		}
		mu := entity.MatchedUsers{
			Category: entity.Category(pbMu.Category),
			UserIDs:  slice.MapFromUint64ToUint(pbMu.UserIds),
		}
		fmt.Println("received message from " + msg.Channel + "channel.")
		fmt.Println(mu)

	}
}
