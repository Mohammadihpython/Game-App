package redismatching

import "GameApp/adaptor/adaptor/redis"

type Config struct {
	WaitingListPrefix string `koanf:"waiting_list_prefix"`
}

type DB struct {
	config  Config
	adaptor redis.Adaptor
}

func New(config Config, adaptor redis.Adaptor) DB {
	return DB{config, adaptor}

}
