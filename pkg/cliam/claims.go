package cliam

import (
	"GameApp/constant"
	"GameApp/service/authservice"
	"github.com/labstack/echo/v4"
)

func GetClaims(c echo.Context) *authservice.Claims {
	claims := c.Get(constant.AuthMiddlewareContextKey)

	//	convert claims object to authService claims object
	cl, ok := claims.(*authservice.Claims)
	if !ok {
		panic("not found claims")

	}
	return cl

}
