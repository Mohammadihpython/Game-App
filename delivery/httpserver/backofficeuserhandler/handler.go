package backofficeuserhandler

import (
	"GameApp/pkg/richerror"
	"GameApp/service/authorizationservice"
	"GameApp/service/authservice"
	"GameApp/service/backofficeuserservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	authConfig        authservice.Config
	authSVC           authservice.Service
	authorizationSVC  authorizationservice.Service
	backofficeUserSVC backofficeuserservice.Service
}

func New(authConfig authservice.Config,
	authSVC authservice.Service,
	authorizationSVC authorizationservice.Service,
	backofficeUserSVC backofficeuserservice.Service) Handler {
	return Handler{
		authConfig:        authConfig,
		authSVC:           authSVC,
		authorizationSVC:  authorizationSVC,
		backofficeUserSVC: backofficeUserSVC,
	}

}

func (h Handler) UserList(c echo.Context) error {
	const op = "backofficeuserHadler.UserList"
	users, err := h.backofficeUserSVC.ListAllUsers()
	if err != nil {
		return richerror.New(op).WithKind(richerror.KindUnexpected)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"users": users,
	})

}
