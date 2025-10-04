package service

import (
	"context"
	"errors"
	"strings"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
)

type CatService interface {
	CreateCat(ctx context.Context, cat *domain.Cat) (int64, error)
	// ListCats(ctx context.Context, p domain.ListCatsParams) ([]domain.Cat, int64, error)
	GetCat(ctx context.Context, id int64) (domain.Cat, error)
}
type CatRepository interface {
	CreateCat(ctx context.Context, cat *domain.Cat) (int64, error)
	GetCat(ctx context.Context, id int64) (domain.Cat, error)
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
	return s.repo.GetCat(ctx, id)
}
