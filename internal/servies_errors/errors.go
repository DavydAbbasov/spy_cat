package servieserrors

import "errors"

var (
	ErrCatNotFound  = errors.New("cat not found")
	ErrBreedInvalid = errors.New("breed is invalid")
	ErrInvalidSalary= errors.New("salary is invalid")
)
