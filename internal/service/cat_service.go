package service

import (
	"context"
	"errors"
	"strings"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
)

type CatService interface {
	CreateCat(ctx context.Context, cat *domain.Cat) (int64, error)
}
type CatRepository interface {
	CreateCat(ctx context.Context, cat *domain.Cat) (int64, error)
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
