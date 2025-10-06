package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type CatRepository struct {
	db *sql.DB
}

func NewCatRepository(db *sql.DB) *CatRepository {
	return &CatRepository{
		db: db,
	}
}
func (r *CatRepository) CreateCat(ctx context.Context, c *domain.Cat) (int64, error) {
	var id int64

	q := `
		INSERT INTO cats(name, years_experience, breed, salary)
		VALUES ($1,$2,$3,$4)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, q, c.Name, c.YearsExperience, c.Breed, c.Salary).Scan(&id)
	return id, err
}
func (r *CatRepository) GetCat(ctx context.Context, id int64) (domain.Cat, error) {
	var c domain.Cat

	q := `SELECT id, name, years_experience, breed, salary
	      FROM cats
		  WHERE id = $1`

	err := r.db.QueryRowContext(ctx, q, id).Scan(&c.ID, &c.Name, &c.YearsExperience, &c.Breed, &c.Salary)
	return c, err
}
func (r *CatRepository) ListCats(ctx context.Context, p domain.ListCatsParams) ([]domain.Cat, error) {
	if p.Limit <= 0 {
		p.Limit = 50
	}
	if p.Offset < 0 {
		p.Offset = 0
	}

	q := `
		SELECT id, name, years_experience, breed, salary
		FROM cats
		ORDER BY id DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.QueryContext(ctx, q, p.Limit, p.Offset)
	if err != nil {
		return nil, fmt.Errorf("list cats: %w", err)
	}
	defer rows.Close()

	out := make([]domain.Cat, 0, p.Limit)

	for rows.Next() {
		var c domain.Cat
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.YearsExperience,
			&c.Breed,
			&c.Salary,
		); err != nil {
			return nil, fmt.Errorf("scan cat: %w", err)
		}
		out = append(out, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iter: %w", err)
	}

	return out, nil
}
func (r *CatRepository) DeleteCat(ctx context.Context, id int64) (int64, error) {
	q := `
	DELETE
	FROM cats
	WHERE id = $1;`

	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return 0, fmt.Errorf("delete cat: %w", err)
	}
	n, _ := res.RowsAffected()
	return n, err
}
