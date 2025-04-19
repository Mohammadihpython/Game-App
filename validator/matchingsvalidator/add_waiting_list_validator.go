package matchingsvalidator

import (
	"GameApp/entity"
	"GameApp/param"
	"GameApp/pkg/richerror"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) AddToWaitingListValidator(req param.AddToWaitingListRequest) (error, map[string]string) {
	const op = "matchingsvalidator.AddToWaitingListValidator"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Category, validation.Required, validation.By(v.IsCategoryValid)),
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

func (v Validator) IsCategoryValid(value interface{}) error {
	category := value.(entity.Category)
	if !category.IsValid() {
		return fmt.Errorf("category is invalid")

	}
	return nil

}
