package uservalidator

import (
	"GameApp/param"
	"GameApp/pkg/phonenumber"
	"GameApp/pkg/richerror"
	"github.com/go-ozzo/ozzo-validation/v4"
	_ "github.com/go-ozzo/ozzo-validation/v4/is"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo}
}

func (v Validator) ValidateRegisterRequest(req param.RegisterRequest) error {
	// TODO : We should verify phone number by verification code
	const op = "uservalidator.ValidateRegisterRequest"
	// Validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return richerror.New(op).WithMessage("phone number is not valid").WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"phone number": req.PhoneNumber})

	}

	if isUnique, err := v.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			// %w wrap the error and show us the last errors corrupted for this error
			//ارور های قبلی را که مربوط به این خطا هست زا نیز نشان می دهد
			return richerror.New(op).WithWrappedError(err)
		}
		if !isUnique {
			richerror.New(op).WithMessage("phone number is not unique").WithKind(richerror.KindInvalid)
		}
	}

	//validate name

	// TODO : Add support for Persion word or not ASID word
	if len(req.Name) < 3 {
		return richerror.New(op).
			WithMessage("name is too short").WithKind(richerror.KindInvalid)

	}

	// TODO validate password with regex
	//TODO add 8 to config
	if len(req.Password) < 8 {
		validation.Validate(req.Password){
			
		}
		return richerror.New(op).
			WithMessage("password is must grader than 8 character").WithKind(richerror.KindInvalid)
	}
	return nil
}
