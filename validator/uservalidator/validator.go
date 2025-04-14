package uservalidator

import "GameApp/entity"

const (
	PhoneNumberRegex = `^09[0-9]{9}$`
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	GetUserByPhone(phoneNumber string) (entity.User, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo}
}
