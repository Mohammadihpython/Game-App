package userservice

import (
	"GameApp/metrics"
	"GameApp/param"
	"GameApp/pkg/richerror"
	"GameApp/repository/mysql/mysqluser"
	"fmt"
	"time"
)

func (s Service) Login(req param.LoginRequest) (param.LoginResponse, error) {
	const op = "userservice.Login"
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RequestDuration().WithLabelValues("userService", "/login").Observe(duration)
	}()

	//  check the existence of phone number from repository
	//	get the mysqluser by phone number
	// TODO : Its better to  separate method for check existence check  of get mysqluser by phone number
	user, err := s.repo.GetUserByPhone(req.PhoneNumber)
	if err != nil {
		return param.LoginResponse{}, richerror.New(op).
			WithMessage("failed to get mysqluser by phone number").
			WithMeta(map[string]interface{}{"Req": req}).WithWrappedError(err)
	}
	if user.Password != mysqluser.GetMD5Hash(req.Password) {
		return param.LoginResponse{}, fmt.Errorf("mysqluser or password not valid ")

	}
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to create token: %w", err)
	}
	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to refresh token: %w", err)
	}
	metrics.RequestCounter().WithLabelValues("userService", "/login").Inc()
	return param.LoginResponse{
		User: param.UserInfo{
			ID:          user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
		Tokens: param.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil

	//	compare mysqluser.password with re.password

}
