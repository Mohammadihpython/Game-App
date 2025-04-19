package userHandler

import (
	"GameApp/delicery/httpserver/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRouter(e *echo.Echo) {
	e.POST("/users/register", h.userRegister)
	e.POST("/users/login", h.userLogin)
	e.GET("/users/profile", h.userProfile, middleware.Auth(h.authSvc, h.authConfig))

}
