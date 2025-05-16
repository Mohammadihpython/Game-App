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
	redispresence "GameApp/repository/redis/presence"
	"GameApp/repository/redis/redismatching"
	"GameApp/scheduler"
	"GameApp/service/authorizationservice"
	"GameApp/service/authservice"
	"GameApp/service/backofficeuserservice"
	"GameApp/service/matchingservice"
	"GameApp/service/presenceservice"
	"GameApp/service/userservice"
	"GameApp/validator/matchingsvalidator"
	"GameApp/validator/uservalidator"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"sync"
)

func main() {
	fmt.Println("start Echo server")
	cfg := conf.Load()
	fmt.Println(cfg)
	// TODO add command for migrations to dont run automatically
	mgr := migrator.New(cfg.Mysql)
	mgr.Up()

	PresenceGrpcConn, err := grpc.NewClient(":8070", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}
	defer PresenceGrpcConn.Close()
	userSvc, authSvc, userValidator, backofficeSVC, authorizationSVC, matchingSVC, matchingV, presenceSVC := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator, authorizationSVC, backofficeSVC, matchingSVC, matchingV, presenceSVC)
	go func() {
		server.Serve()
	}()
	done := make(chan bool)
	var wg sync.WaitGroup
	go func() {
		sch := scheduler.New(cfg.Scheduler, matchingSVC)
		wg.Add(1)
		sch.Start(done, &wg)
	}()

	// تابع signal وفتی دستور بستن برنامه میاد اون رو داخل چنل میریزه
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// اینحا ما یک کانتکس میساریم و وقتی برتامه میخواد بسته بشه به بقیه سرویس ها اطلاع میدهد
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Application.GracefulShutdownTimeout)
	defer cancel()
	if err := server.Router.Shutdown(ctx); err != nil {
		fmt.Println("shutdown server error:", err)
	}

	fmt.Println("close server")
	// اینجا با پر کردن جنل done اسکژولز های خودمون میکیم که کار بسته تمومش کن
	done <- true
	// اینجا یا کاری نداره و بسته میشه یا تایم ان درخواست تمام میشه
	<-ctx.Done()

	wg.Wait()

}

func setupServices(cfg conf.Config) (
	userservice.Service,
	authservice.Service,
	uservalidator.Validator,
	backofficeuserservice.Service,
	authorizationservice.Service,
	matchingservice.Service,
	matchingsvalidator.Validator,
	presenceservice.Service,
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

	redisAdaptor := redis.New(cfg.Redis)
	matcingv := matchingsvalidator.New()

	matcingRepo := redismatching.New(cfg.RedisMatching, redisAdaptor)

	presenceSClient := presenceClient.New(":8086")

	matchingSVC := matchingservice.New(cfg.MatchingService, matcingRepo, presenceSClient, redisAdaptor)

	presenceRepo := redispresence.New(redisAdaptor)

	presenceSVC := presenceservice.New(cfg.Presence, presenceRepo)

	return userSvc, authSvc, userValidator, backofficeUserSvc, authorizationSvc, matchingSVC, matcingv, presenceSVC
}
