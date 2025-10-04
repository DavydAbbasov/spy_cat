package postgres

import (
	"database/sql"
	"fmt"

	"github.com/DavydAbbasov/spy-cat/internal/config"
	_ "github.com/lib/pq"
	log "github.com/rs/zerolog/log"
	// _ "github.com/tommy-muehle/go-mnd/v2/config"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(cfg *config.Config, dsn string) (*Storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error get database driver %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connection to database %w", err)
	}

	log.Info().Msg("database Repos connection is success")

	r := &Storage{
		db: db,
	}

	return r, nil
}

func (r *Storage) Close() error {
	return r.db.Close()
}
