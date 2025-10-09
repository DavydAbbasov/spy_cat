package dto

import "github.com/DavydAbbasov/spy-cat/internal/domain"

type CatResponse struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	YearsExperience int64   `json:"years_experience"`
	Breed           string  `json:"breed"`
	Salary          float64 `json:"salary"`
}

type CreateCatRequest struct {
	Name            string  `json:"name"              validate:"required,min=2,max=64"`
	YearsExperience int64   `json:"years_experience"  validate:"gte=0,lte=60"`
	Breed           string  `json:"breed"             validate:"required,min=2,max=64"`
	Salary          float64 `json:"salary"            validate:"gte=0,lte=1000000"`
}

type CreateCatResponse struct {
	ID int64 `json:"id"`
}
type UpdateSalaryRequest struct {
	Salary float64 `json:"salary" validate:"required,gte=0,lte=1000000"`
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
	Items      []CatResponse `json:"items"`
	Limit      int           `json:"limit"`
	Offset     int           `json:"offset"`
	NextOffset int           `json:"next_offset"`
}

// type ErrorResponse struct {
// 	Code    string      `json:"code"    example:"INVALID_INPUT"`
// 	Message string      `json:"message" example:"validation error"`
// 	Details interface{} `json:"details,omitempty"`
// }
type DeleteCatResponse struct {
	Deleted bool  `json:"deleted"`
	ID      int64 `json:"id"`
}

// mapping
func ToNewCatDomain(req CreateCatRequest) domain.Cat {
	return domain.Cat{
		Name:            req.Name,
		YearsExperience: req.YearsExperience,
		Breed:           req.Breed,
		Salary:          req.Salary,
	}
}
func ToCatResponse(c domain.Cat) CatResponse {
	return CatResponse{
		ID:              c.ID,
		Name:            c.Name,
		YearsExperience: c.YearsExperience,
		Breed:           c.Breed,
		Salary:          c.Salary,
	}
}

func ToCatResponses(items []domain.Cat) []CatResponse {
	out := make([]CatResponse, 0, len(items))
	for _, it := range items {
		out = append(out, ToCatResponse(it))
	}
	return out
}
