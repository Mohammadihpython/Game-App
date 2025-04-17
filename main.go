package main

import (
	"GameApp/conf"
	"GameApp/delicery/httpserver"
	"GameApp/repository/migrator"
	"GameApp/repository/mysql"
	"GameApp/repository/mysql/mysqlaccesscontrol"
	"GameApp/repository/mysql/mysqluser"
	"GameApp/service/authorizationservice"
	"GameApp/service/authservice"
	"GameApp/service/backofficeuserservice"
	"GameApp/service/userservice"
	"GameApp/validator/uservalidator"
	"fmt"
	"time"
)

const (
	SECRET                = "Hmdsfksdf"
	AccessExpirationTime  = time.Hour * 24
	RefreshExpirationTime = time.Hour * 24 * 7
	AccessSubject         = "at"
	RefreshSubject        = "rt"
)

func main() {
	fmt.Println("start Echo server")
	//cfg := conf.Load()
	cfg := conf.Config{
		HTTPServer: conf.HTTPServer{Port: 8080},
		Auth: authservice.Config{
			SignKey:               SECRET,
			AccessExpirationTime:  AccessExpirationTime,
			RefreshExpirationTime: RefreshExpirationTime,
			AccessSubject:         AccessSubject,
			RefreshSubject:        RefreshSubject,
		},
		Mysql: mysql.Config{
			Host:     "localhost",
			Port:     3308,
			Username: "Hamed",
			Password: "hmah8013",
			DBName:   "gameappDB",
		},
	}
	// TODO add command for migrations to dont run automatically
	mgr := migrator.New(cfg.Mysql)
	mgr.Up()
	userSvc, authSvc, userValidator, backofficeSVC, authorizationSVC := setupServices(cfg)
	server := httpserver.New(cfg, authSvc, userSvc, userValidator, authorizationSVC, backofficeSVC)
	server.Serve()

}

func setupServices(cfg conf.Config) (
	userservice.Service,
	authservice.Service,
	uservalidator.Validator,
	backofficeuserservice.Service,
	authorizationservice.Service) {

	authSvc := authservice.New(cfg.Auth)
	mysqlRepo := mysql.New(cfg.Mysql)

	userMysql := mysqluser.New(mysqlRepo)
	userSvc := userservice.New(userMysql, authSvc)

	userValidator := uservalidator.New(userMysql)

	backofficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(mysqlRepo)

	authorizationSvc := authorizationservice.New(aclMysql)

	return userSvc, authSvc, userValidator, backofficeUserSvc, authorizationSvc
}
