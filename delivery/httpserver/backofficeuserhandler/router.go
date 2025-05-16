package backofficeuserhandler

import (
	"GameApp/delivery/httpserver/middleware"
	"GameApp/entity"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetBackOfficeUserRouter(e *echo.Echo) {
	eg := e.Group("back-office/users")

	eg.GET("/", h.UserList, middleware.Auth(h.authSVC, h.authConfig),
		middleware.AccessCheck(h.authorizationSVC, entity.UserListPermission))

}
