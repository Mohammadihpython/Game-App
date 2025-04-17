package middleware

import (
	"GameApp/entity"
	"GameApp/pkg/cliam"
	"GameApp/service/authorizationservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AccessCheck(service authorizationservice.Service, permissions ...entity.PermissionTitle) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := cliam.GetClaims(c)
			isAllowed, err := service.CheckAccess(claims.UserID, claims.Role, permissions...)
			if err != nil || !isAllowed {
				// TODO - Log unexpected error
				return c.JSON(http.StatusForbidden, echo.Map{
					"message": "user not allowed",
				})
			}
			return next(c)
		}
	}
}
