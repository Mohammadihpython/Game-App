package httpmsg

import (
	"GameApp/pkg/richerror"
	"errors"
	"net/http"
)

func CodeAndMessage(err error) (message string, code int) {
	// use errors.As instead of check directly err.(richerror.RichError)
	//Since Go 1.13, errors can be wrapped using the fmt.Errorf function with the %w verb. Therefore,
	//type assertion or type switch on errors fails on wrapped errors.
	//The preferred way for checking for a specific error type is to use the errors.
	//As function from the standard library as this function traverses the chain of the wrapped errors while checking for a specific error type.
	var richError richerror.RichError
	switch {
	case errors.As(err, &richError):

		var re richerror.RichError
		msg := re.Message()
		// because we get recursive message and code if it didn't set in uper layer
		// when we get server error not secure to send message of them  this to client
		code := MapKindToHttpStatusCode(re.Kind())
		if code >= 500 {
			msg = "internal server error"
		}
		return msg, code
	default:
		return err.Error(), http.StatusBadRequest
	}
}

func MapKindToHttpStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError

	}
}
