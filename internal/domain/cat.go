package domain

type Cat struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	YearsExperience int64   `json:"years_experience"`
	Breed           string  `json:"breed"`
	Salary          float64 `json:"salary"`
}
type ListCatsParams struct {
	Name     *string
	Breed    *string
	MinYears *int
	MaxYears *int
	Limit    int
	Offset   int
}
