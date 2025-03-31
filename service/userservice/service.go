package userservice

import (
	"GameApp/entity"
	"GameApp/pkg/phonenumber"
	"errors"
	"fmt"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	RegisterUser(u entity.User) (entity.User, error)
}

type Service struct {
	repo Repository
}
type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone"`
}
type RegisterResponse struct {
	User entity.User `json:"user"`
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO : We should verify phone number by verification code
	// Validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, errors.New("invalid phone")
	}
	// check uniqueness of phone number
	// we used short hand if here
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			// %w wrap the error and show us the last errors corrupted for this error
			//ارور های قبلی را که مربوط به این خطا هست زا نیز نشان می دهد
			return RegisterResponse{}, fmt.Errorf("failed to check if phone is unique: %w", err)
		}
		if !isUnique {
			return RegisterResponse{}, errors.New("phone number is already used")
		}
	}

	//validate name
	// TODO : Add support for Persion word or not ASID word
	if len(req.Name) < 3 {
		return RegisterResponse{}, errors.New("name is too short")
	}
	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}
	// create new user in storage
	createdUser, err := s.repo.RegisterUser(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("failed to register user: %w", err)
	}

	//	 return created user
	return RegisterResponse{User: createdUser}, nil

}

func New(repo Repository) Service {
	return Service{repo: repo}

}
