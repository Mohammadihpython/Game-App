package redismatching

import (
	"GameApp/entity"
	"GameApp/pkg/richerror"
	"GameApp/pkg/timestamp"
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func (d DB) AddToWaitingList(userID uint, category entity.Category) error {
	const OP = richerror.Op("redismatching.AddToWaitingList")
	fmt.Println("AddToWaitingList", userID, category)
	_, err := d.adaptor.Client().ZAdd(
		context.Background(),
		fmt.Sprintf("%s:%s", d.config.WaitingListPrefix, category),
		redis.Z{
			Score:  float64(timestamp.Now()),
			Member: fmt.Sprintf("%d", userID),
		}).Result()
	if err != nil {
		return richerror.New(OP).WithWrappedError(err).WithKind(richerror.KindUnexpected)

	}
	return nil

}

func (d DB) GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error) {
	const OP = richerror.Op("redismatching.GetWaitingListByCategory")
	minOrder := fmt.Sprintf("%d", timestamp.Add(-1*time.Hour))
	maxOrder := strconv.Itoa(int(timestamp.Now()))
	list, err := d.adaptor.Client().ZRangeByScoreWithScores(ctx, d.getCategory(category), &redis.ZRangeBy{
		Min:    minOrder,
		Max:    maxOrder,
		Offset: 0,
		Count:  0,
	}).Result()
	if err != nil {
		return nil, richerror.New(OP).WithWrappedError(err)
	}
	var result = make([]entity.WaitingMember, 0)
	for _, l := range list {
		fmt.Println("userID from redis ", l.Member)
		userID, _ := strconv.Atoi(l.Member.(string))
		result = append(result, entity.WaitingMember{
			UserID:    uint(userID),
			Timestamp: int64(l.Score),
			Category:  category,
		})

	}

	fmt.Println("its result of users", result)
	return result, nil
}

func (d DB) getCategory(category entity.Category) string {
	return fmt.Sprintf("%s:%s", d.config.WaitingListPrefix, category)
}

func (d DB) RemoveUsersFromWaitingList(userID []uint, category entity.Category) {
	ctx, cansel := context.WithTimeout(context.Background(), time.Second*5)
	defer cansel()
	members := make([]any, 0)
	for _, u := range userID {
		members = append(members, fmt.Sprintf("%d", u))
	}
	numberOfRemoved, err := d.adaptor.Client().ZRem(ctx, d.getCategory(category), members...).Result()
	if err != nil {
		log.Errorf("remove from waiting list error: %s", err.Error())
		// TODO update metrics
	}
	log.Printf("remove from waiting list: %d", numberOfRemoved)
	//	 TODO update metric
	return
}
