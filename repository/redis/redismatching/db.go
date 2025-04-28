package redismatching

import (
	"GameApp/adaptor/redis"
)

type Config struct {
	WaitingListPrefix string `koanf:"waiting_list"`
}

type DB struct {
	config  Config
	adaptor redis.Adaptor
}

func New(config Config, adaptor redis.Adaptor) DB {
	return DB{config, adaptor}

}
