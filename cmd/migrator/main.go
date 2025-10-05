package main

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/DavydAbbasov/spy-cat/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	log "github.com/rs/zerolog/log"
)

func main() {
	_ = godotenv.Load(".env")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("load config")
	}

	user := url.QueryEscape(cfg.Postgres.User)
	pass := url.QueryEscape(cfg.Postgres.Password)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user,
		pass,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("migrate new")
	}
	defer m.Close()

	if err := m.Up(); err != nil {

		if errors.Is(err, migrate.ErrNoChange) {
			log.Info().Msg("no migrations to apply")

			return
		}
		log.Fatal().Err(err).Msg("migrate up")
	}
	log.Info().Msg("migrations applied")
}
