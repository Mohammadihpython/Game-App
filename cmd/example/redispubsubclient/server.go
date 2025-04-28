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

		payload := protobufEncoder.DecoderEvent(entity.MatchingUsersMatchedEvent, msg.Payload)
		p, ok := payload.(entity.MatchedUsers)
		if !ok {
			// log the error
			return
		}
		fmt.Println(p.Category, p.UserIDs)
		fmt.Println("received message from " + msg.Channel + "channel.")

	}
}
