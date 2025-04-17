package userservice

import (
	"GameApp/entity"
	"GameApp/param"
	"GameApp/repository/mysql/mysqluser"
	"fmt"
)

func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {

	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    mysqluser.GetMD5Hash(req.Password),
		Role:        entity.UserRole,
	}
	// create new mysqluser in storage
	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return param.RegisterResponse{}, fmt.Errorf("failed to register mysqluser: %w", err)
	}

	//	 return created mysqluser
	return param.RegisterResponse{User: createdUser}, nil

}
