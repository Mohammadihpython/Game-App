package main

import (
	"GameApp/adaptor/adaptor/redis"
	"GameApp/conf"
	"GameApp/delicery/httpserver"
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
	"fmt"
)

func main() {
	fmt.Println("start Echo server")
	cfg := conf.Load()
	fmt.Println(cfg)
	// TODO add command for migrations to dont run automatically
	mgr := migrator.New(cfg.Mysql)
	mgr.Up()
	userSvc, authSvc, userValidator, backofficeSVC, authorizationSVC, matchingSVC, matchingV := setupServices(cfg)
	server := httpserver.New(cfg, authSvc, userSvc, userValidator, authorizationSVC, backofficeSVC, matchingSVC, matchingV)
	server.Serve()

}

func setupServices(cfg conf.Config) (
	userservice.Service,
	authservice.Service,
	uservalidator.Validator,
	backofficeuserservice.Service,
	authorizationservice.Service,
	matchingservice.Service,
	matchingsvalidator.Validator,
) {

	authSvc := authservice.New(cfg.Auth)
	mysqlRepo := mysql.New(cfg.Mysql)

	userMysql := mysqluser.New(mysqlRepo)
	userSvc := userservice.New(userMysql, authSvc)

	userValidator := uservalidator.New(userMysql)

	backofficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(mysqlRepo)

	authorizationSvc := authorizationservice.New(aclMysql)

	// we must create an redis client and pass it to matching service
	matcingv := matchingsvalidator.New()
	redisAdaptor := redis.New(cfg.Redis)
	matcingRepo := redismatching.New(redisAdaptor)
	matchingSVC := matchingservice.New(cfg.MatchingService, matcingRepo)

	return userSvc, authSvc, userValidator, backofficeUserSvc, authorizationSvc, matchingSVC, matcingv
}
