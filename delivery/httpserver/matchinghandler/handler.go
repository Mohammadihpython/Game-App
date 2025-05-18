package matchinghandler

import (
	"GameApp/adaptor/presence"
	"GameApp/param"
	"GameApp/pkg/httpmsg"
	"GameApp/pkg/richerror"
	"GameApp/service/authservice"
	"GameApp/service/matchingservice"
	"GameApp/validator/matchingsvalidator"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	authConfig     authservice.Config
	authSVC        authservice.Service
	matchingSVC    matchingservice.Service
	validator      matchingsvalidator.Validator
	presenceClient presence.Client
}

func New(authConfig authservice.Config,
	authSVC authservice.Service,
	matchingSVC matchingservice.Service,
	validator matchingsvalidator.Validator,
	presenceClient presence.Client,

) Handler {
	return Handler{authConfig,
		authSVC,
		matchingSVC,
		validator,
		presenceClient,
	}
}

func (h Handler) AddToWaitingList(c echo.Context) error {
	const OP = "matchinghandler.AddToWaitingList"
	var req param.AddToWaitingListRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err, fieldError := h.validator.AddToWaitingListValidator(req)
	if err != nil && fieldError != nil {
		msg, code := httpmsg.CodeAndMessage(err)
		fmt.Println(msg, code)
		return c.JSON(
			code, echo.Map{
				"message": msg,
				"errors":  fieldError,
			})
	}
	res, err := h.matchingSVC.AddToWaitingList(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, richerror.New(OP).
			WithWrappedError(err).
			WithMessage("wrong password or hhhh phone number").
			WithKind(richerror.KindInvalid),
		)
	}
	return c.JSON(http.StatusOK, res)

}
