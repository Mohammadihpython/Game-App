package main

import (
	"GameApp/adaptor/redis"
	"GameApp/conf"
	"GameApp/entity"
	"GameApp/pkg/protobufEncoder"
	"context"
	"fmt"
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

		matchedUsers := protobufEncoder.DecoderEvent(entity.MatchingUsersMatchedEvent, msg.Payload)

		fmt.Println(matchedUsers.Category, matchedUsers.UserIDs)
		fmt.Println("received message from " + msg.Channel + "channel.")

	}
}
