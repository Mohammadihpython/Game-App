package matchinghandler

import (
	"GameApp/delivery/httpserver/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRouter(e *echo.Echo) {

	e.POST("/users/add-to-game", h.AddToWaitingList, middleware.Auth(h.authSVC, h.authConfig), middleware.UpsertPresence(h.presenceClient))

}
