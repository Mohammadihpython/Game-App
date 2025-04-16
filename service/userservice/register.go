package userservice

import (
	"GameApp/entity"
	"GameApp/param"
	"GameApp/repository/mysql/user"
	"fmt"
)

func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {

	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    user.GetMD5Hash(req.Password),
	}
	// create new user in storage
	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return param.RegisterResponse{}, fmt.Errorf("failed to register user: %w", err)
	}

	//	 return created user
	return param.RegisterResponse{User: createdUser}, nil

}
