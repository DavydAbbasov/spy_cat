package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
	serviceerrors "github.com/DavydAbbasov/spy-cat/internal/servies_errors"
	"github.com/rs/zerolog/log"
)

type MissionService interface {
	CreateMission(ctx context.Context, p domain.CreateMissionParams) (domain.Mission, error)
	AssignCat(ctx context.Context, missionID int64, catID *int64) error
	GetMission(ctx context.Context, id int64) (domain.Mission, []domain.MissionGoal, error)
	List(ctx context.Context, f domain.MissionFilter) ([]domain.MissionListItem, int, error)
	UpdateStatus(ctx context.Context, p domain.UpdateMissionStatusParams) (domain.Mission, error)
	AddGoal(ctx context.Context, missionID int64, p domain.CreateGoalParams) (domain.MissionGoal, error)
}
type MissionRepository interface {
	BeginTx(ctx context.Context) (Tx, error)
	InsertMission(ctx context.Context, tx Tx, m *domain.Mission) (int64, error)
	InsertGoals(ctx context.Context, tx Tx, missionID int64, goals []domain.MissionGoal) error
	AssignCat(ctx context.Context, tx Tx, missionID int64, catID *int64) error
	GetMission(ctx context.Context, id int64) (domain.Mission, error)
	GetMissionGoals(ctx context.Context, missionID int64) ([]domain.MissionGoal, error)
	ListMissions(ctx context.Context, f domain.MissionFilter) ([]domain.MissionListItem, int, error)
	UpdateStatusIfCurrent(ctx context.Context, id int64, newStatus, expected domain.MissionStatus) (domain.Mission, bool, error)
	InsertGoal(ctx context.Context, missionID int64, p domain.CreateGoalParams) (domain.MissionGoal, error)
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
		return domain.Mission{}, serviceerrors.ErrInvalidCreateMission
	}

	goals := make([]domain.MissionGoal, 0, len(p.Goals))
	for _, g := range p.Goals {
		name := strings.TrimSpace(g.Name)
		country := strings.ToUpper(strings.TrimSpace(g.Country))
		notes := strings.TrimSpace(g.Notes)

		if name == "" {
			return domain.Mission{}, serviceerrors.ErrInvalidGoalName
		}
		if len(country) != 2 {
			return domain.Mission{}, serviceerrors.ErrInvalidCountry
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
		return serviceerrors.ErrInvalidCreateMission
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	err = s.repo.AssignCat(ctx, tx, missionID, catID)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrMissionNotFound) {
			return serviceerrors.ErrMissionNotFound
		}
		if errors.Is(err, serviceerrors.ErrCatNotFound) {
			return serviceerrors.ErrCatNotFound
		}
		return err
	}

	return tx.Commit(ctx)
}
func (s *missionService) GetMission(ctx context.Context, id int64) (domain.Mission, []domain.MissionGoal, error) {
	if id <= 0 {
		return domain.Mission{}, nil, serviceerrors.ErrMissionNotFound
	}

	mission, err := s.repo.GetMission(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Mission{}, nil, serviceerrors.ErrMissionNotFound
		}
		return domain.Mission{}, nil, err
	}

	goals, err := s.repo.GetMissionGoals(ctx, id)
	if err != nil {
		return domain.Mission{}, nil, err
	}

	return mission, goals, nil
}

func (s *missionService) List(ctx context.Context, f domain.MissionFilter) ([]domain.MissionListItem, int, error) {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	return s.repo.ListMissions(ctx, f)
}
func (s *missionService) UpdateStatus(ctx context.Context, p domain.UpdateMissionStatusParams) (domain.Mission, error) {
	if p.ID <= 0 {
		return domain.Mission{}, serviceerrors.ErrMissionNotFound
	}

	m, err := s.repo.GetMission(ctx, p.ID)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrMissionNotFound) {
			return domain.Mission{}, serviceerrors.ErrMissionNotFound
		}
		return domain.Mission{}, err
	}

	newStatus := domain.MissionStatus(strings.TrimSpace(string(p.Status)))
	if newStatus != domain.StatusPlanned && newStatus != domain.StatusActive && newStatus != domain.StatusCompleted {
		return domain.Mission{}, serviceerrors.ErrInvalidStatus
	}

	if m.Status == newStatus {
		return m, nil
	}

	if !domain.CanTransition(m.Status, newStatus) {
		return domain.Mission{}, serviceerrors.ErrInvalidTransition
	}

	updated, ok, err := s.repo.UpdateStatusIfCurrent(ctx,
		p.ID,
		newStatus,
		m.Status,
	)
	if err != nil {
		return domain.Mission{}, err
	}
	if !ok {
		return domain.Mission{}, serviceerrors.ErrConflict
	}

	return updated, nil
}
func (s *missionService) AddGoal(ctx context.Context, missionID int64, p domain.CreateGoalParams) (domain.MissionGoal, error) {
	if missionID <= 0 {
		return domain.MissionGoal{}, serviceerrors.ErrMissionNotFound
	}

	name := strings.TrimSpace(p.Name)
	if name == "" || len(name) < 2 || len(name) > 64 {
		return domain.MissionGoal{}, serviceerrors.ErrInvalidGoalName
	}

	country := strings.ToUpper(strings.TrimSpace(p.Country))
	if len(country) != 2 {
		return domain.MissionGoal{}, serviceerrors.ErrInvalidCountry
	}

	notes := strings.TrimSpace(p.Notes)

	m, err := s.repo.GetMission(ctx, missionID)
	if err != nil {
		if errors.Is(err, serviceerrors.ErrMissionNotFound) {
			return domain.MissionGoal{}, serviceerrors.ErrMissionNotFound
		}
		return domain.MissionGoal{}, err
	}
	if m.Status == domain.StatusCompleted {
		return domain.MissionGoal{}, serviceerrors.ErrMissionCompleted
	}

	goal, err := s.repo.InsertGoal(ctx, missionID, domain.CreateGoalParams{
		Name:    name,
		Country: country,
		Notes:   notes,
	})
	if err != nil {
		if errors.Is(err, serviceerrors.ErrMissionNotFound) {
			return domain.MissionGoal{}, serviceerrors.ErrMissionNotFound
		}
		return domain.MissionGoal{}, err
	}

	return goal, nil
}
