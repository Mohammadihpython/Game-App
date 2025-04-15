package mymiddleware

import (
	"GameApp/constant"
	"GameApp/service/authservice"
	ejwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

//
//func Auth(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		// دریافت توکن از هدر
//		authHeader := c.Request().Header.Get("Authorization")
//		if authHeader == "" {
//			return echo.NewHTTPError(http.StatusUnauthorized, "توکن یافت نشد")
//		}
//
//		// جدا کردن Bearer از توکن
//		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
//
//		// کلید رمزنگاری
//		secretKey := []byte("Hmdsfksdf")
//
//		// اعتبارسنجی توکن
//		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, fmt.Errorf("روش امضای نامعتبر: %v", token.Header["alg"])
//			}
//			return secretKey, nil
//		})
//
//		if err != nil {
//			return echo.NewHTTPError(http.StatusUnauthorized, "توکن نامعتبر")
//		}
//
//		// بررسی اعتبار claims
//		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//			// ذخیره اطلاعات کاربر در context
//
//			c.Set("user_id", uint(claims["user_id"].(float64)))
//			c.Set("username", claims["username"])
//			return next(c)
//		}
//
//		return echo.NewHTTPError(http.StatusUnauthorized, "توکن منقضی شده")
//	}
//}

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
