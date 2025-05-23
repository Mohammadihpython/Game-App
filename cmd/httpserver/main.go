package main

import (
	presenceClient "GameApp/adaptor/presence"
	"GameApp/adaptor/redis"
	"GameApp/conf"
	"GameApp/delivery/httpserver"
	"GameApp/repository/migrator"
	"GameApp/repository/mysql"
	"GameApp/repository/mysql/mysqlaccesscontrol"
	"GameApp/repository/mysql/mysqluser"
	"GameApp/repository/redis/redismatching"
	"GameApp/service/authorizationservice"
	"GameApp/service/authservice"
	"GameApp/service/backofficeuserservice"
	"GameApp/service/matchingservice"
	"GameApp/service/userservice"
	"GameApp/validator/matchingsvalidator"
	"GameApp/validator/uservalidator"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println("start Echo server")
	cfg := conf.Load()
	fmt.Println(cfg)
	// TODO add command for migrations to dont run automatically
	mgr := migrator.New(cfg.Mysql, "../../repository/mysql/migrations")
	mgr.Up()
	userSvc, authSvc, userValidator, backofficeSVC, authorizationSVC, matchingSVC, matchingV, presenceSVC := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator, authorizationSVC, backofficeSVC, matchingSVC, matchingV, presenceSVC)

	go func() {
		server.Serve()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Application.GracefulShutdownTimeout)
	defer cancel()
	if err := server.Router.Shutdown(ctx); err != nil {
		fmt.Println("shutdown server error:", err)
	}
	fmt.Println("close server")
	<-ctx.Done()
	time.Sleep(5 * time.Second)

}

func setupServices(cfg conf.Config) (
	userservice.Service,
	authservice.Service,
	uservalidator.Validator,
	backofficeuserservice.Service,
	authorizationservice.Service,
	matchingservice.Service,
	matchingsvalidator.Validator,
	presenceClient.Client,
) {

	authSvc := authservice.New(cfg.Auth)
	mysqlRepo := mysql.New(cfg.Mysql)

	userMysql := mysqluser.New(mysqlRepo)
	userSvc := userservice.New(userMysql, authSvc)

	userValidator := uservalidator.New(userMysql)

	backofficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(mysqlRepo)

	authorizationSvc := authorizationservice.New(aclMysql)

	// we must create a redis client and pass it to matching service
	matchingv := matchingsvalidator.New()
	redisAdaptor := redis.New(cfg.Redis)
	matchingRepo := redismatching.New(cfg.RedisMatching, redisAdaptor)

	presenceAdaptor := presenceClient.New(":8080", nil)
	pub := redis.New(cfg.Redis)
	matchingSVC := matchingservice.New(cfg.MatchingService, matchingRepo, presenceAdaptor, pub)

	presenceC := presenceClient.New(":8086", grpc.DialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))

	return userSvc, authSvc, userValidator, backofficeUserSvc, authorizationSvc, matchingSVC, matchingv, presenceC
}
