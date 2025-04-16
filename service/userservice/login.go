package userservice

import (
	"GameApp/param"
	"GameApp/pkg/richerror"
	"GameApp/repository/mysql/user"
	"fmt"
)

func (s Service) Login(req param.LoginRequest) (param.LoginResponse, error) {
	const op = "userservice.Login"
	//  check the existence of phone number from repository
	//	get the user by phone number
	// TODO : Its better to  separate method for check existence check  of get user by phone number
	user, err := s.repo.GetUserByPhone(req.PhoneNumber)
	if err != nil {
		return param.LoginResponse{}, richerror.New(op).
			WithMessage("failed to get user by phone number").
			WithMeta(map[string]interface{}{"Req": req}).WithWrappedError(err)
	}
	if user.Password != user.GetMD5Hash(req.Password) {
		return param.LoginResponse{}, fmt.Errorf("user or password not valid ")

	}
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to create token: %w", err)
	}
	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("failed to refresh token: %w", err)
	}

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

	//	compare user.password with re.password

}
