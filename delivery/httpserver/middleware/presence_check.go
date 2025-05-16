package middleware

import (
	"GameApp/param"
	"GameApp/pkg/cliam"
	"GameApp/pkg/timestamp"
	"GameApp/service/presenceservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpsertPresence(service presenceservice.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := cliam.GetClaims(c)
			service.Upsert(c.Request().Context(), param.UpsertPresenceRequest{
				UserID:    claims.UserID,
				Timestamp: timestamp.Now(),
			})
			_, err := service.Upsert(c.Request().Context(), param.UpsertPresenceRequest{
				UserID:    claims.UserID,
				Timestamp: timestamp.Now(),
			})
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": err.Error(),
				})
			}
			return next(c)

		}

	}

}
