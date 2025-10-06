package domain

type Cat struct {
	ID              int64
	Name            string
	YearsExperience int64
	Breed           string
	Salary          float64
}
type ListCatsParams struct {
	Name     *string
	Breed    *string
	MinYears *int
	MaxYears *int
	Limit    int
	Offset   int
}
type UpdateSalaryParams struct {
    ID     int64
    Salary float64
}