package redismatching

import (
	"GameApp/entity"
	"GameApp/pkg/richerror"
	"GameApp/pkg/timestamp"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const WaithingListPrefix = "waitingList"

func (d DB) AddToWaitingList(userID uint, category entity.Category) error {
	const OP = richerror.Op("redismatching.AddToWaitingList")

	_, err := d.adaptor.Client().ZAdd(
		context.Background(),
		fmt.Sprintf("&s:%s", WaithingListPrefix, category),
		redis.Z{
			Score:  float64(timestamp.Now()),
			Member: fmt.Sprintf("%d", userID),
		}).Result()
	if err != nil {
		return richerror.New(OP).WithWrappedError(err).WithKind(richerror.KindUnexpected)

	}
	return nil

}
