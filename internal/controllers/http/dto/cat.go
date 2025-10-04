package dto

type CreateCatRequest struct {
	Name            string  `json:"name"              validate:"required,min=2,max=64"`
	YearsExperience int64   `json:"years_experience"  validate:"required,min=0,max=60"`
	Breed           string  `json:"breed"             validate:"required"`
	Salary          float64 `json:"salary"            validate:"required,gte=0,lte=1000000"`
}

type CreateCatResponse struct {
	ID int64 `json:"id"`
}
