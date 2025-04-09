package uservalidator

import (
	"GameApp/param"
	"GameApp/pkg/richerror"
	"errors"
	"fmt"
	"github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
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

func (v Validator) ValidateRegisterRequest(req param.RegisterRequest) (error, map[string]string) {
	const op = "uservalidator.ValidateRegisterRequest"
	// TODO : Add support for Persion word or not ASID word

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 60)),
		validation.Field(&req.Password, validation.Required, validation.Match(regexp.MustCompile(`^[a-zA-Z0-9!&^%$#]{8,}$`))),
		validation.Field(&req.PhoneNumber, validation.Required, validation.Match(regexp.MustCompile(`^09[0-9]{9}$`)),
			validation.By(v.CheckPhoneNumberUnique)),
	); err != nil {

		fieldError := make(map[string]string)
		fmt.Printf("%T", err)
		var errv validation.Errors
		ok := errors.As(err, &errv)
		if ok {
			for key, value := range errv {
				if value != nil {
					fieldError[key] = value.Error()
				}
			}
		}
		return richerror.New(op).WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"req": req}).WithWrappedError(err), fieldError

	}
	// TODO : We should verify phone number by verification code

	return nil, nil
}

func (v Validator) CheckPhoneNumberUnique(value interface{}) error {
	phoneNumber := value.(string)

	if isUnique, err := v.repo.IsPhoneNumberUnique(phoneNumber); err != nil || !isUnique {
		if err != nil {
			// %w wrap the error and show us the last errors corrupted for this error
			//ارور های قبلی را که مربوط به این خطا هست زا نیز نشان می دهد
			fmt.Println(err.Error())
			return err
		}
		if !isUnique {
			return fmt.Errorf("phone number is not uniqe")
		}
	}
	return nil
}
