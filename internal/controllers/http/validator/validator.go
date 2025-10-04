package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	goval "github.com/go-playground/validator/v10"
)

var ErrHandlerValidationFailed = errors.New("handler_validation_failed")

type Validator struct {
	validator *goval.Validate
}

func NewValidator() *Validator {
	v := goval.New()
	return &Validator{
		validator: v,
	}
}

func (v *Validator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		return fmt.Errorf("%w: %v", ErrHandlerValidationFailed, err)
	}
	return nil
}

func DecodeJSON[T any](v *Validator, r *http.Request) (*T, error) {
	var payload T

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	defer r.Body.Close()

	if err := dec.Decode(&payload); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrHandlerValidationFailed, err)
	}

	if err := ensureEOF(dec); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrHandlerValidationFailed, err)
	}

	if err := v.Validate(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
func ensureEOF(dec *json.Decoder) error {
	var extra any

	if err := dec.Decode(&extra); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return errors.New("unexpected trailing data in JSON body")
}
