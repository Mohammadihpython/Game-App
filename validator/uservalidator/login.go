package uservalidator

import (
	"GameApp/param"
	"GameApp/pkg/richerror"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

func (v Validator) ValidateLoginRequest(req param.LoginRequest) (error, map[string]string) {
	const op = "uservalidator.ValidateLoginRequest"
	// TODO : Add support for Persion word or not ASID word

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.PhoneNumber, validation.Required, validation.Match(regexp.MustCompile(PhoneNumberRegex)),
			validation.By(v.DosePhoneNumberExist)),
		validation.Field(&req.Password, validation.Required),
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

func (v Validator) DosePhoneNumberExist(value interface{}) error {

	phoneNumber := value.(string)
	if _, err := v.repo.GetUserByPhone(phoneNumber); err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil

}
