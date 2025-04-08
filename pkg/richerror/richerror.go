package richerror

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
)

type Op string

type RichError struct {
	operation    Op
	wrappedError error
	message      string
	kind         Kind
	meta         map[string]interface{}
}

// First Approach to Handle Parameter of RichError
// First Approach Example
// richError.New(richerror.Op("userService.Profile),"unexpected error",.....)
//func New(args ...interface{}) RichError {
//	r := RichError{}
//	if len(args) == 0 {
//		return r
//	}
//	for _, arg := range args {
//		switch arg.(type) {
//		case Kind:
//			r.Kind = arg.(Kind)
//		case string:
//			r.message = arg.(string)
//		case error:
//			r.wrappedError = arg.(error)
//		case map[string]interface{}:
//			r.meta = arg.(map[string]interface{})
//		case Op:
//			r.operation = arg.(Op)
//
//		}
//	}
//	return r
//}

// Second Approach To Handel Params use composite method
// Example: richerror.New().WithMessage("not found error").WithMeta()

func New(operation Op) RichError {
	return RichError{operation: operation}

}
func (r RichError) WithMessage(message string) RichError {
	r.message = message
	return r
}
func (r RichError) WithKind(kind Kind) RichError {
	r.kind = kind
	return r
}
func (r RichError) WithMeta(meta map[string]interface{}) RichError {
	r.meta = meta
	return r

}
func (r RichError) WithWrappedError(err error) RichError {
	r.wrappedError = err
	return r
}

func (r RichError) Error() string {
	return r.message

}
func (r RichError) Kind() Kind {
	// doing recursive to wrapped error to get kind of that error
	if r.kind != 0 {
		return r.kind
	}
	re, ok := r.wrappedError.(RichError)
	if !ok {
		return 0

	}
	return re.Kind()
}
func (r RichError) Message() string {
	if r.message != "" {
		return r.message
	}
	re, ok := r.wrappedError.(RichError)
	if !ok {
		return r.wrappedError.Error()

	}
	return re.Message()
}
