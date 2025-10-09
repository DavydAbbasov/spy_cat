package dto

type ErrorResponse struct {
	Code    string      `json:"code"    example:"INVALID_INPUT"`
	Message string      `json:"message" example:"validation error"`
	Details interface{} `json:"details,omitempty"`
}
