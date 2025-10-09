package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
	serviceserrors "github.com/DavydAbbasov/spy-cat/internal/servies_errors"
	"github.com/rs/zerolog/log"
)

type MissionService interface {
	CreateMission(ctx context.Context, p domain.CreateMissionParams) (domain.Mission, error)
	AssignCat(ctx context.Context, missionID int64, catID *int64) error
}
type MissionRepository interface {
	BeginTx(ctx context.Context) (Tx, error)
	InsertMission(ctx context.Context, tx Tx, m *domain.Mission) (int64, error)
	InsertGoals(ctx context.Context, tx Tx, missionID int64, goals []domain.MissionGoal) error
	AssignCat(ctx context.Context, tx Tx, missionID int64, catID *int64) error
}
type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
type missionService struct {
	repo MissionRepository
}

func NewMissionService(repo MissionRepository) MissionService {
	return &missionService{
		repo: repo,
	}
}
func (s *missionService) CreateMission(ctx context.Context, p domain.CreateMissionParams) (domain.Mission, error) {
	p.Title = strings.TrimSpace(p.Title)
	p.Description = strings.TrimSpace(p.Description)

	if p.Title == "" {
		return domain.Mission{}, serviceserrors.ErrInvalidCreateMission
	}

	goals := make([]domain.MissionGoal, 0, len(p.Goals))
	for _, g := range p.Goals {
		name := strings.TrimSpace(g.Name)
		country := strings.ToUpper(strings.TrimSpace(g.Country))
		notes := strings.TrimSpace(g.Notes)

		if name == "" {
			return domain.Mission{}, serviceserrors.ErrInvalidGoalName
		}
		if len(country) != 2 {
			return domain.Mission{}, serviceserrors.ErrInvalidCountry
		}

		goals = append(goals, domain.MissionGoal{
			Name:    name,
			Country: country,
			Notes:   notes,
			Status:  domain.GoalTodo,
		})

	}
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return domain.Mission{}, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Warn().Err(err).Msg("rollback failed")
		}
	}()

	m := domain.Mission{
		Title:       p.Title,
		Description: p.Description,
		Status:      domain.StatusPlanned,
		CatID:       nil,
	}

	id, err := s.repo.InsertMission(ctx, tx, &m)
	if err != nil {
		return domain.Mission{}, err
	}
	if len(goals) > 0 {
		if err := s.repo.InsertGoals(ctx, tx, id, goals); err != nil {
			return domain.Mission{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Mission{}, err
	}

	m.ID = id
	return m, nil
}
func (s *missionService) AssignCat(ctx context.Context, missionID int64, catID *int64) error {
	if missionID <= 0 {
		return serviceserrors.ErrInvalidCreateMission
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = s.repo.AssignCat(ctx, tx, missionID, catID)
	if err != nil {
		if errors.Is(err, serviceserrors.ErrMissionNotFound) {
			return serviceserrors.ErrMissionNotFound
		}
		if errors.Is(err, serviceserrors.ErrCatNotFound) {
			return serviceserrors.ErrCatNotFound
		}
		return err
	}

	return tx.Commit(ctx)
}
