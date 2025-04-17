package middleware

import (
	"GameApp/constant"
	"GameApp/service/authservice"
	ejwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Auth(service authservice.Service, config authservice.Config) echo.MiddlewareFunc {
	return ejwt.WithConfig(ejwt.Config{

		ContextKey:    constant.AuthMiddlewareContextKey,
		SigningKey:    []byte(config.SignKey),
		SigningMethod: "Hs256",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims, err := service.ParsToken(auth)
			if err != nil {
				return nil, err
			}
			return claims, nil

		},
	})

}
