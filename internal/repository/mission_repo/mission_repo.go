package postgres

import (
	"context"
	"database/sql"
	"errors"

	serviceerrors "github.com/DavydAbbasov/spy-cat/internal/servies_errors"
	"github.com/rs/zerolog/log"

	"github.com/DavydAbbasov/spy-cat/internal/domain"
	service "github.com/DavydAbbasov/spy-cat/internal/service/mission_service"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type pgTx struct{ tx *sql.Tx }

func (t *pgTx) Commit(ctx context.Context) error   { return t.tx.Commit() }
func (t *pgTx) Rollback(ctx context.Context) error { return t.tx.Rollback() }

type MissionRepo struct {
	db *sql.DB
}

func NewMissionRepository(db *sql.DB) *MissionRepo {
	return &MissionRepo{db: db}
}

func (r *MissionRepo) BeginTx(ctx context.Context) (service.Tx, error) {
	raw, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	return &pgTx{tx: raw}, nil
}

func (r *MissionRepo) InsertMission(ctx context.Context, tx service.Tx, m *domain.Mission) (int64, error) {
	pgtx := tx.(*pgTx)

	const q = `
		INSERT INTO missions (title, description, status, cat_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	err := pgtx.tx.QueryRowContext(ctx, q, m.Title, m.Description, m.Status, m.CatID).Scan(&m.ID)
	return m.ID, err
}

func (r *MissionRepo) InsertGoals(ctx context.Context, tx service.Tx, missionID int64, goals []domain.MissionGoal) error {
	pgtx := tx.(*pgTx)

	const q = `
		INSERT INTO mission_goals (mission_id, name, country, notes)
		VALUES ($1, $2, $3, $4);
	`

	for _, g := range goals {
		if _, err := pgtx.tx.ExecContext(ctx, q, missionID, g.Name, g.Country, g.Notes); err != nil {
			return err
		}
	}
	return nil
}
func (r *MissionRepo) AssignCat(ctx context.Context, tx service.Tx, missionID int64, catID *int64) error {
	pgtx := tx.(*pgTx)

	q := `
		UPDATE missions
		SET cat_id = $2, updated_at = now()
		WHERE id = $1
		RETURNING id;
	`

	var id int64
	err := pgtx.tx.QueryRowContext(ctx, q, missionID, catID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Int64("mission_id", missionID).Msg("mission not found")
			return serviceerrors.ErrMissionNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return serviceerrors.ErrCatNotFound
		}
		return err
	}

	return nil
}
func (r *MissionRepo) GetMission(ctx context.Context, id int64) (domain.Mission, error) {
	var m domain.Mission

	q := `
		SELECT id, title, description,status, cat_id, created_at, updated_at
		FROM missions
		WHERE id = $1;`

	err := r.db.QueryRowContext(ctx, q, id).
		Scan(
			&m.ID,
			&m.Title,
			&m.Description,
			&m.Status,
			&m.CatID,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Mission{}, serviceerrors.ErrMissionNotFound
		}
		return domain.Mission{}, err
	}
	return m, nil
}
func (r *MissionRepo) GetMissionGoals(ctx context.Context, missionID int64) ([]domain.MissionGoal, error) {
	q := `
		SELECT id, mission_id, name, country, notes, status, created_at, updated_at
		FROM mission_goals
		WHERE mission_id = $1
		ORDER BY id;
	`
	rows, err := r.db.QueryContext(ctx, q, missionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.MissionGoal
	for rows.Next() {
		var g domain.MissionGoal
		if err := rows.Scan(&g.ID, &g.MissionID, &g.Name, &g.Country, &g.Notes, &g.Status, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
