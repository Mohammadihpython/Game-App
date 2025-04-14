package userservice

import (
	"GameApp/entity"
)

type Repository interface {
	RegisterUser(u entity.User) (entity.User, error)
	GetUserByPhone(phoneNumber string) (entity.User, error)
	GetUserByID(userID uint) (entity.User, error)
}

type AuthService interface {
	CreateRefreshToken(user entity.User) (string, error)
	CreateAccessToken(user entity.User) (string, error)
}

type Service struct {
	repo Repository
	auth AuthService
}

func New(repo Repository, auth AuthService) Service {
	return Service{repo: repo, auth: auth}

}
