package main

import (
	"GameApp/adaptor/redis"
	"GameApp/conf"
	"GameApp/delivery/grpcserver/presenceserver"
	redispresence "GameApp/repository/redis/presence"
	"GameApp/service/presenceservice"
)

func main() {
	cfg := conf.Load()
	redisAdaptor := redis.New(cfg.Redis)
	presenceRepo := redispresence.New(redisAdaptor)
	presenceSVC := presenceservice.New(cfg.Presence, presenceRepo)

	server := presenceserver.New(presenceSVC)
	server.Start()
}
