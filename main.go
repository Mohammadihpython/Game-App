package main

import (
	"GameApp/conf"
	"GameApp/delicery/httpserver"
	"GameApp/repository/migrator"
	"GameApp/repository/mysql"
	"GameApp/service/authservice"
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
	userSvc, authSvc, userValidator := setupServices(cfg)
	server := httpserver.New(cfg, authSvc, userSvc, userValidator)
	server.Serve()

}

func setupServices(cfg conf.Config) (userservice.Service, authservice.Service, uservalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)
	mysqlRepo := mysql.New(cfg.Mysql)

	userSvc := userservice.New(mysqlRepo, authSvc)
	userValidator := uservalidator.New(mysqlRepo)
	return userSvc, authSvc, userValidator
}
