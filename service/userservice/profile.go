package userservice

import (
	"GameApp/param"
	"GameApp/pkg/richerror"
	"context"
	"fmt"
)

func (s Service) Profile(ctx context.Context, req param.ProfileRequest) (param.ProfileResponse, error) {
	// get mysqluser by ID
	user, err := s.repo.GetUserByID(ctx, req.UserID)
	if err != nil {
		//I have not expected the repository call return
		//"record not found" error,
		// because I assume the input is sanitized
		// TODO we can use Rich Error
		return param.ProfileResponse{}, richerror.New("userService.Profile").WithWrappedError(err)
	}
	fmt.Println(param.ProfileResponse{Name: user.Name})
	return param.ProfileResponse{Name: user.Name}, nil

}
