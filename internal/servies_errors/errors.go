package servieserrors

import "errors"

var (
	ErrCatNotFound     = errors.New("cat not found")
	ErrBreedInvalid    = errors.New("breed invalid")
	ErrInvalidSalary   = errors.New("salary invalid")
	ErrExternalService = errors.New("external service")
)

