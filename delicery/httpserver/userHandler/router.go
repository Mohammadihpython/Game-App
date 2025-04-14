package userHandler

import (
	"GameApp/delicery/mymiddleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetUserRouter(e *echo.Echo) {
	e.POST("/users/register", h.userRegister)
	e.POST("/users/login", h.userLogin)
	e.GET("/users/profile", h.userProfile, mymiddleware.Auth)

}
