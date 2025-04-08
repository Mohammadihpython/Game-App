package httpserver

import (
	"GameApp/pkg/httpmsg"
	"GameApp/pkg/richerror"
	"GameApp/service/userservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s Server) userRegister(c echo.Context) error {

	var req userservice.RegisterRequest
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, richerror.New("userservice.userRegister").
			WithWrappedError(err).
			WithMessage("wrong password or phone number").
			WithKind(richerror.KindInvalid),
		)
	}
	res, err := s.userSvc.Register(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, richerror.New("userservice.userRegister").
			WithWrappedError(err).
			WithMessage("wrong password or hhhh phone number").
			WithKind(richerror.KindInvalid),
		)
	}
	return c.JSON(http.StatusCreated, res)

}

func (s Server) userLogin(c echo.Context) error {
	var req userservice.LoginRequest
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	res, err := s.userSvc.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, richerror.New("userservice.userLogin").
			WithWrappedError(err).
			WithMessage("wrong password or phone number").
			WithKind(richerror.KindInvalid),
		)

	}
	return c.JSON(http.StatusOK, res)

}

func (s Server) userProfile(c echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")
	claims, err := s.authSvc.ParsToken(authToken)
	if err != nil {
		msg, code := httpmsg.CodeAndMessage(err)
		return echo.NewHTTPError(code, msg)
	}
	res, err := s.userSvc.Profile(userservice.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			richerror.New("httpserver.userProfile").
				WithWrappedError(err).
				WithMessage("internal service error").
				WithKind(richerror.KindUnexpected),
		)
	}

	return c.JSON(http.StatusOK, res)
}
