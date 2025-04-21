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
	if r.message != "" {
		return r.wrappedError.Error()
	}
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
