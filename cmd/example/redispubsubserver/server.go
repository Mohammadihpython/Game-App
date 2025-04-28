package main

import (
	"GameApp/adaptor/redis"
	"GameApp/conf"
	"GameApp/contract/goproto/matching"
	"GameApp/entity"
	"GameApp/pkg/slice"
	"context"
	"encoding/base64"
	"google.golang.org/protobuf/proto"
)

func main() {
	cfg := conf.Load()
	redisAdaptor := redis.New(cfg.Redis)
	topic := "matching.users_matched"

	matchedUsers := entity.MatchedUsers{
		Category: entity.FootballCategory,
		UserIDs:  []uint{1, 4},
	}
	protobufmu := matching.MatchedUsers{
		Category: string(matchedUsers.Category),
		UserIds:  slice.MapFromUintToUint64(matchedUsers.UserIDs),
	}
	payload, err := proto.Marshal(&protobufmu)
	if err != nil {
		panic(err)
	}
	payloadStr := base64.StdEncoding.EncodeToString(payload)
	if err := redisAdaptor.Client().Publish(context.Background(), topic, payloadStr).Err(); err != nil {
		panic(err)
	}

}
