package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/DavydAbbasov/spy-cat/internal/config"
	"github.com/DavydAbbasov/spy-cat/internal/lib/catapi"
	postgres "github.com/DavydAbbasov/spy-cat/internal/lib/postgresql"
	"github.com/joho/godotenv"

	catrepository "github.com/DavydAbbasov/spy-cat/internal/repository/cat_repo"
	missionrepository "github.com/DavydAbbasov/spy-cat/internal/repository/mission_repo"

	catservice "github.com/DavydAbbasov/spy-cat/internal/service/cat_service"
	missionservice "github.com/DavydAbbasov/spy-cat/internal/service/mission_service"

	log "github.com/rs/zerolog/log"
)

func Run() error {
	_ = godotenv.Load(".env")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	db, err := postgres.NewConn(cfg, cfg.Postgres.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	defer db.Close()

	//catApi
	breedClient := catapi.NewClient(
		cfg.CatAPI.BaseURL,
		cfg.CatAPI.APIKey,
		cfg.CatAPI.Timeout,
	)
	// repository
	catRepo := catrepository.NewCatRepository(db)
	missionRepo := missionrepository.NewMissionRepository(db)

	// services
	catSvc := catservice.NewCatService(catRepo, breedClient)
	missionSvc := missionservice.NewMissionService(missionRepo)

	httpServer := &http.Server{
		Addr:    cfg.HTTP.Addr,
		Handler: NewRouter(catSvc, missionSvc),

		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()

	log.Info().Msg(fmt.Sprintf("listening on %s", cfg.HTTP.Addr))

	<-ctx.Done()
	log.Info().Msg("shutting down server gracefully")

	httpServer.SetKeepAlivesEnabled(false)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown server")
	}

	log.Info().Msg("server gracefully shutdown")

	return nil
}
