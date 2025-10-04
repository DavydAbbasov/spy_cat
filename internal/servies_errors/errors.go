package servieserrors

import "errors"

type ServiceError struct {
	Msg string
}

func newServiceError(msg string) error {
	return &ServiceError{
		Msg: msg,
	}
}
func (b *ServiceError) Error() string {
	return b.Msg
}

var (
	ErrCatNotFound  = errors.New("cat not found")
	ErrBreedInvalid = errors.New("breed is invalid")
)
