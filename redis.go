package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var ctx = context.Background()
	rdb.ZAdd(ctx, "footbal", redis.Z{
		Score: float64(time.Now().UnixMicro()),
	})

}
