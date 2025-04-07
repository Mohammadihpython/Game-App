package main

import (
	"GameApp/conf"
	"GameApp/delicery/httpserver"
	"GameApp/repository/mysql"
	"GameApp/service/authservice"
	"GameApp/service/userservice"
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

	userSvc, authSvc := setupServices(cfg)
	server := httpserver.New(cfg, authSvc, userSvc)
	server.Serve()

	//httpserver.HandleFunc("/users/register",userRegisterHandler)
	//httpserver.ListenAndServe(:8080,nil)
	// another way is to use multiplexer
	//mux := http.NewServeMux()
	//mux.HandleFunc("/users/register", userRegisterHandler)
	//mux.HandleFunc("/health", HealthCheck)
	//mux.HandleFunc("/users/login", userLoginHandler)
	//mux.HandleFunc("/users/profile", userProfileHandler)
	//
	//server = httpserver.New(cfg, authSvc, userSvc)
	//log.Fatal(server.ListenAndServe())

}

func setupServices(cfg conf.Config) (userservice.Service, authservice.Service) {
	authSvc := authservice.New(cfg.Auth)
	mysqlRepo := mysql.New(cfg.Mysql)

	userSvc := userservice.New(mysqlRepo, authSvc)
	return userSvc, authSvc
}
