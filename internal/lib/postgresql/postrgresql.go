package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/DavydAbbasov/spy-cat/internal/config"
	"github.com/rs/zerolog/log"
)

func NewConn(cfg *config.Config, dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error get database driver %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connection to database %w", err)
	}

	log.Info().Msg("database Repos connection is success")

	return db, nil
}
