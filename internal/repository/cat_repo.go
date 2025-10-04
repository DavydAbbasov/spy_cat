package postgres

import (
	"context"
	"database/sql"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
)

type CatRepo struct {
	db *sql.DB
}

func NewActivityRepo(db *sql.DB) (*CatRepo, error) {
	return &CatRepo{
		db: db,
	}, nil
}
func (r *CatRepo) CreateCat(ctx context.Context, c *domain.Cat) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO cats(name, years_experience, breed, salary)
		 VALUES ($1,$2,$3,$4) RETURNING id`,
		c.Name, c.YearsExperience, c.Breed, c.Salary,
	).Scan(&id)
	return id, err
}
