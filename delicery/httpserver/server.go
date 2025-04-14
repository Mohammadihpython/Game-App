package httpserver

import (
	"GameApp/conf"
	"GameApp/delicery/httpserver/userHandler"
	"GameApp/service/authservice"
	"GameApp/service/userservice"
	"GameApp/validator/uservalidator"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config      conf.Config
	userHandler userHandler.Handler
}

func New(cfg conf.Config, authSvc authservice.Service, userSvc userservice.Service, validator uservalidator.Validator) Server {
	return Server{
		config:      cfg,
		userHandler: userHandler.New(authSvc, userSvc, validator),
	}
}

func (s Server) Serve() {

	e := echo.New()

	//	 Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//	Routes
	e.GET("/", s.healthCheck)
	s.userHandler.SetUserRouter(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))

}
