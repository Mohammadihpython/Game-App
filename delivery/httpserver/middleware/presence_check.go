package middleware

import (
	"GameApp/adaptor/presence"
	"GameApp/param"
	"GameApp/pkg/cliam"
	"GameApp/pkg/timestamp"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UpsertPresence(service presence.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := cliam.GetClaims(c)
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
