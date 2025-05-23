package redispresence

import (
	"GameApp/adaptor/redis"
)

type DB struct {
	adaptor redis.Adaptor
}

func New(adaptor redis.Adaptor) DB {
	return DB{adaptor}

}
