package httpserver

import (
	"GameApp/conf"
	"GameApp/service/authservice"
	"GameApp/service/userservice"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config  conf.Config
	authSvc authservice.Service
	userSvc userservice.Service
}

func New(cfg conf.Config, authSvc authservice.Service, userSvc userservice.Service) Server {
	return Server{
		config:  cfg,
		authSvc: authSvc,
		userSvc: userSvc,
	}
}

func (s Server) Serve() {

	e := echo.New()

	//	 Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//	Routes
	e.GET("/", s.healthCheck)
	e.POST("/users/register", s.userRegister)
	e.POST("/users/login", s.userLogin)
	e.GET("/users/profile", s.userProfile)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))

}
