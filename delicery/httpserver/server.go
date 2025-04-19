package httpserver

import (
	"GameApp/conf"
	"GameApp/delicery/httpserver/backofficeuserhandler"
	"GameApp/delicery/httpserver/matchinghandler"
	"GameApp/delicery/httpserver/userHandler"
	"GameApp/service/authorizationservice"
	"GameApp/service/authservice"
	"GameApp/service/backofficeuserservice"
	"GameApp/service/matchingservice"
	"GameApp/service/userservice"
	"GameApp/validator/matchingsvalidator"
	"GameApp/validator/uservalidator"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config                conf.Config
	userHandler           userHandler.Handler
	backofficeUserHandler backofficeuserhandler.Handler
	matchingHandler       matchinghandler.Handler
}

func New(cfg conf.Config,
	authSvc authservice.Service,
	userSvc userservice.Service,
	validator uservalidator.Validator,
	authorizationSvc authorizationservice.Service,
	backOfficeUseSVC backofficeuserservice.Service,
	matchingSVC matchingservice.Service,
	matchingValidator matchingsvalidator.Validator,

) Server {
	fmt.Println(cfg)
	return Server{
		config:                cfg,
		userHandler:           userHandler.New(cfg.Auth, authSvc, userSvc, validator, cfg.Auth.SignKey),
		backofficeUserHandler: backofficeuserhandler.New(cfg.Auth, authSvc, authorizationSvc, backOfficeUseSVC),
		matchingHandler:       matchinghandler.New(cfg.Auth, authSvc, matchingSVC, matchingValidator),
	}
}

func (s Server) Serve() {

	e := echo.New()

	//	 Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//	Routes
	e.GET("/", s.healthCheck)
	s.userHandler.SetRouter(e)
	s.backofficeUserHandler.SetBackOfficeUserRouter(e)
	s.matchingHandler.SetRouter(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))

}
