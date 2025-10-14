package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

type whereParts struct {
	sql  string
	args []any
}

func buildWhere(f domain.MissionFilter) whereParts {
	var conds []string
	var args []any
	i := 1

	// status = $1
	if f.Status != nil {
		conds = append(conds, fmt.Sprintf("status = $%d", i))
		args = append(args, *f.Status)
		i++
	}

	// cat_id = $N
	if f.CatID != nil {
		conds = append(conds, fmt.Sprintf("cat_id = $%d", i))
		args = append(args, *f.CatID)
		i++
	}

	// title ILIKE $N
	if f.Q != nil {
		q := strings.TrimSpace(*f.Q)
		if q != "" {
			conds = append(conds, fmt.Sprintf("title ILIKE $%d", i))
			args = append(args, "%"+q+"%")
			i++
		}
	}

	if len(conds) == 0 {
		return whereParts{}
	}

	return whereParts{
		sql:  " WHERE " + strings.Join(conds, " AND "),
		args: args,
	}
}
func (r *MissionRepo) queryItems(ctx context.Context, w whereParts, limit, offset int) ([]domain.MissionListItem, error) {
	sel := `
	SELECT id, title, status, cat_id, created_at
	FROM missions
	`
	order := `
	ORDER BY created_at
	DESC, id DESC
	`

	limitPos := len(w.args) + 1
	offsetPos := limitPos + 1

	q := sel + w.sql + order + fmt.Sprintf(" LIMIT $%d OFFSET $%d", limitPos, offsetPos)

	args := make([]any, 0, len(w.args)+2)
	args = append(args, w.args...)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.MissionListItem, 0, limit)
	for rows.Next() {
		var it domain.MissionListItem
		var cat sql.NullInt64

		if err := rows.Scan(
			&it.ID,
			&it.Title,
			&it.Status,
			&cat,
			&it.CreatedAt); err != nil {
			return nil, err
		}

		if cat.Valid {
			v := cat.Int64
			it.CatID = &v
		}

		items = append(items, it)
	}

	return items, rows.Err()
}

func (r *MissionRepo) queryTotal(ctx context.Context, w whereParts) (int, error) {
	q := `
	SELECT count(*)
	FROM missions
	`

	var total int
	if err := r.db.QueryRowContext(ctx, q+w.sql, w.args...).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *MissionRepo) ListMissions(ctx context.Context, f domain.MissionFilter) ([]domain.MissionListItem, int, error) {
	w := buildWhere(f)

	items, err := r.queryItems(ctx, w, f.Limit, f.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("items: %w", err)
	}

	total, err := r.queryTotal(ctx, w)
	if err != nil {
		return nil, 0, fmt.Errorf("count: %w", err)
	}

	return items, total, nil
}
func (r *MissionRepo) UpdateStatusIfCurrent(ctx context.Context, id int64, newStatus, expected domain.MissionStatus) (domain.Mission, bool, error) {
	q := `
	UPDATE missions
	SET status = $2, updated_at = now()
	WHERE id = $1 AND status = $3
	RETURNING id, title, description, status, cat_id, created_at, updated_at;
	`
	var m domain.Mission
	err := r.db.QueryRowContext(ctx, q, id, newStatus, expected).
		Scan(
			&m.ID,
			&m.Title,
			&m.Description,
			&m.Status,
			&m.CatID,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Mission{}, false, nil
	}
	if err != nil {
		return domain.Mission{}, false, err
	}
	
	return m, true, nil
}
