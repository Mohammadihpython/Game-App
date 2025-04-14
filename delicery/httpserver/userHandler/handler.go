package userHandler

import (
	"GameApp/param"
	"GameApp/pkg/httpmsg"
	"GameApp/pkg/richerror"
	"GameApp/service/authservice"
	"GameApp/service/userservice"
	"GameApp/validator/uservalidator"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
}

func New(authSVC authservice.Service, userSvc userservice.Service, userValidator uservalidator.Validator) Handler {
	return Handler{authSvc: authSVC, userSvc: userSvc, userValidator: userValidator}
}

func (h Handler) userRegister(c echo.Context) error {
	var req param.RegisterRequest
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, richerror.New("userservice.userRegister").
			WithWrappedError(err).
			WithMessage("wrong password or phone number").
			WithKind(richerror.KindInvalid),
		)
	}
	if err, fieldError := h.userValidator.ValidateRegisterRequest(req); err != nil {
		fmt.Println(err.Error())
		fmt.Println(fieldError)

		msg, code := httpmsg.CodeAndMessage(err)
		fmt.Println(msg, code)
		return c.JSON(
			code, echo.Map{
				"message": msg,
				"errors":  fieldError,
			})
	}

	res, err := h.userSvc.Register(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, richerror.New("userservice.userRegister").
			WithWrappedError(err).
			WithMessage("wrong password or hhhh phone number").
			WithKind(richerror.KindInvalid),
		)
	}
	return c.JSON(http.StatusCreated, res)

}

func (h Handler) userLogin(c echo.Context) error {
	var req param.LoginRequest
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	res, err := h.userSvc.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, richerror.New("userservice.userLogin").
			WithWrappedError(err).
			WithMessage("wrong password or phone number").
			WithKind(richerror.KindInvalid),
		)

	}
	return c.JSON(http.StatusOK, res)

}

func (h Handler) userProfile(c echo.Context) error {
	//authToken := c.Request().Header.Get("Authorization")
	//claims, err := h.authSvc.ParsToken(authToken)
	//if err != nil {
	//	msg, code := httpmsg.CodeAndMessage(err)
	//	return echo.NewHTTPError(code, msg)
	//}
	userID := c.Get("user_id")
	fmt.Println(userID)
	res, err := h.userSvc.Profile(param.ProfileRequest{UserID: userID.(uint)})
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
