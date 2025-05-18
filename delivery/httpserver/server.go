package httpserver

import (
	"GameApp/adaptor/presence"
	"GameApp/conf"
	"GameApp/delivery/httpserver/backofficeuserhandler"
	"GameApp/delivery/httpserver/matchinghandler"
	"GameApp/delivery/httpserver/userHandler"
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
	Router                *echo.Echo
}

func New(cfg conf.Config,
	authSvc authservice.Service,
	userSvc userservice.Service,
	validator uservalidator.Validator,
	authorizationSvc authorizationservice.Service,
	backOfficeUseSVC backofficeuserservice.Service,
	matchingSVC matchingservice.Service,
	matchingValidator matchingsvalidator.Validator,
	presenceSVC presence.Client,

) Server {
	fmt.Println(cfg)
	return Server{
		config:                cfg,
		userHandler:           userHandler.New(cfg.Auth, authSvc, userSvc, validator, cfg.Auth.SignKey, presenceSVC),
		backofficeUserHandler: backofficeuserhandler.New(cfg.Auth, authSvc, authorizationSvc, backOfficeUseSVC),
		matchingHandler:       matchinghandler.New(cfg.Auth, authSvc, matchingSVC, matchingValidator, presenceSVC),
		Router:                echo.New(),
	}
}

func (s Server) Serve() {

	//	 Middleware
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())

	//	Routes
	s.Router.GET("/", s.healthCheck)
	s.userHandler.SetRouter(s.Router)
	s.backofficeUserHandler.SetBackOfficeUserRouter(s.Router)
	s.matchingHandler.SetRouter(s.Router)
	fmt.Println(s.config.HTTPServer.Port)
	s.Router.Logger.Fatal(s.Router.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))

}
