package playhandler

import (
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRouter(e *echo.Echo) {
	e.GET("/ws", h.GameHandler)

}
