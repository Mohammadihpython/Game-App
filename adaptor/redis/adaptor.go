package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string `koanf:"host"`
	Password string `koanf:"password"`
	Port     string `koanf:"port"`
	DB       int    `koanf:"db"`
}

type Adaptor struct {
	client *redis.Client
}

func New(config Config) Adaptor {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
	return Adaptor{client: rdb}

}

func (a *Adaptor) Client() *redis.Client {
	return a.client
}
