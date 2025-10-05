package dto

import "github.com/DavydAbbasov/spy-cat/internal/domain"

type CreateCatRequest struct {
	Name            string  `json:"name"              validate:"required,min=2,max=64"`
	YearsExperience int64   `json:"years_experience"  validate:"required,min=0,max=60"`
	Breed           string  `json:"breed"             validate:"required"`
	Salary          float64 `json:"salary"            validate:"required,gte=0,lte=1000000"`
}

type CreateCatResponse struct {
	ID int64 `json:"id"`
}
type GetCatsQuery struct {
	Name     *string `form:"name"       binding:"omitempty,min=1"`
	Breed    *string `form:"breed"      binding:"omitempty,min=1"`
	MinYears *int    `form:"min_years"  binding:"omitempty,min=0"`
	MaxYears *int    `form:"max_years"  binding:"omitempty,min=0,gtefield=MinYears"`
	Limit    int     `form:"limit,default=10"  binding:"omitempty,min=1,max=200"`
	Offset   int     `form:"offset,default=0"  binding:"omitempty,min=0"`
}
type GetCatsResponse struct {
	Items      []domain.Cat `json:"items"`
	Limit      int          `json:"limit"`
	Offset     int          `json:"offset"`
	NextOffset int          `json:"next_offset"`
}

type ErrorResponse struct {
	Code    string      `json:"code"    example:"INVALID_INPUT"`
	Message string      `json:"message" example:"validation error"`
	Details interface{} `json:"details,omitempty"`
}
