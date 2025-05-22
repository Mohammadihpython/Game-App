package httpserver

import (
	"GameApp/adaptor/presence"
	"GameApp/conf"
	"GameApp/delivery/httpserver/backofficeuserhandler"
	"GameApp/delivery/httpserver/matchinghandler"
	"GameApp/delivery/httpserver/playhandler"
	"GameApp/delivery/httpserver/userHandler"
	"GameApp/pkg/logger"
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
	"go.uber.org/zap"
)

type Server struct {
	config                conf.Config
	userHandler           userHandler.Handler
	backofficeUserHandler backofficeuserhandler.Handler
	matchingHandler       matchinghandler.Handler

	playHandler playhandler.Handler
	Router      *echo.Echo
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
	return Server{
		config:                cfg,
		userHandler:           userHandler.New(cfg.Auth, authSvc, userSvc, validator, cfg.Auth.SignKey, presenceSVC),
		backofficeUserHandler: backofficeuserhandler.New(cfg.Auth, authSvc, authorizationSvc, backOfficeUseSVC),
		matchingHandler:       matchinghandler.New(cfg.Auth, authSvc, matchingSVC, matchingValidator, presenceSVC),
		playHandler:           playhandler.Handler{},
		Router:                echo.New(),
	}
}

func (s Server) Serve() {

	//	 Middleware
	//s.Router.Use(middleware.Logger())

	// use Zap logger for echo instead of its logger
	s.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogContentLength: true,
		LogRemoteIP:      true,
		LogLatency:       true,
		LogUserAgent:     true,
		LogURIPath:       true,
		LogRequestID:     true,
		LogMethod:        true,

		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Logger.Named("http-server").Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
				zap.String("method", v.Method),
				zap.Duration("latency", v.Latency),
				zap.Int64("response_size", v.ResponseSize),
				zap.String("request_id", v.RequestID),
				zap.String("user-agent", v.UserAgent),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
				zap.String("remote_IP", v.RemoteIP),
			)

			return nil
		},
	}))
	s.Router.Use(middleware.RequestID())
	s.Router.Use(middleware.Recover())

	//	Routes
	s.Router.GET("/", s.healthCheck)
	s.userHandler.SetRouter(s.Router)
	s.backofficeUserHandler.SetBackOfficeUserRouter(s.Router)
	s.matchingHandler.SetRouter(s.Router)
	s.playHandler.SetRouter(s.Router)
	s.Router.Logger.Fatal(s.Router.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))

}
