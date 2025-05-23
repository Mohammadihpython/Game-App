package conf

import (
	"GameApp/adaptor/redis"
	"GameApp/repository/mysql"
	"GameApp/repository/redis/redismatching"
	"GameApp/scheduler"
	"GameApp/service/authservice"
	"GameApp/service/matchingservice"
	"GameApp/service/presenceservice"
	"time"
)

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
}

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	Application     Application            `koanf:"application"`
	HTTPServer      HTTPServer             `koanf:"http_server"`
	Auth            authservice.Config     `koanf:"auth"`
	Mysql           mysql.Config           `koanf:"mysql"`
	MatchingService matchingservice.Config `koanf:"matching_service"`
	Redis           redis.Config           `koanf:"redis"`
	Presence        presenceservice.Config `koanf:"presence_service"`
	Scheduler       scheduler.Config       `koanf:"scheduler"`
	RedisMatching   redismatching.Config   `koanf:"redis_matching"`
}
