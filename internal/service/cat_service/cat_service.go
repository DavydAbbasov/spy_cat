package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	servieserrors "github.com/DavydAbbasov/spy-cat/internal/servies_errors"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
)

type CatService interface {
	CreateCat(ctx context.Context, cat *domain.Cat) (int64, error)
	ListCats(ctx context.Context, p domain.ListCatsParams) ([]domain.Cat, error)
	GetCat(ctx context.Context, id int64) (domain.Cat, error)
	DeleteCat(ctx context.Context, id int64) (int64, error)
}
type CatRepository interface {
	CreateCat(ctx context.Context, cat *domain.Cat) (int64, error)
	ListCats(ctx context.Context, p domain.ListCatsParams) ([]domain.Cat, error)
	GetCat(ctx context.Context, id int64) (domain.Cat, error)
	DeleteCat(ctx context.Context, id int64) (int64, error)
}
type catService struct {
	repo CatRepository
}

func NewCatService(repo CatRepository) CatService {
	return &catService{repo: repo}
}

func (s *catService) CreateCat(ctx context.Context, cat *domain.Cat) (int64, error) {
	cat.Name = strings.TrimSpace(cat.Name)
	if cat.Name == "" {
		return 0, errors.New("name is required")
	}

	return s.repo.CreateCat(ctx, cat)
}

func (s *catService) GetCat(ctx context.Context, id int64) (domain.Cat, error) {
	if id <= 0 {
		return domain.Cat{}, errors.New("invalid id")
	}

	cat, err := s.repo.GetCat(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Cat{}, servieserrors.ErrCatNotFound
		}
		return domain.Cat{}, err
	}

	return cat, nil
}
func (s *catService) ListCats(ctx context.Context, p domain.ListCatsParams) ([]domain.Cat, error) {

	if p.MinYears != nil && p.MaxYears != nil && *p.MinYears > *p.MaxYears {
		return nil, errors.New("min years cannot be greater than max years")
	}

	return s.repo.ListCats(ctx, p)
}
func (s *catService) DeleteCat(ctx context.Context, id int64) (int64, error) {
	if id <= 0 {
		return 0, errors.New("invalid id")
	}

	affected, err := s.repo.DeleteCat(ctx, id)
	if err != nil {
		return 0, err
	}

	if affected == 0 {
		return 0, servieserrors.ErrCatNotFound
	}
	return 0, nil
}
