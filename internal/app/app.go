package app

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/DavydAbbasov/spy-cat/config"
	"github.com/DavydAbbasov/spy-cat/internal/controllers/http/validator"
	postgres "github.com/DavydAbbasov/spy-cat/internal/repository"

	"github.com/DavydAbbasov/spy-cat/internal/service"
	"github.com/gin-gonic/gin"
	log "github.com/rs/zerolog/log"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	db, err := postgres.NewStorage(cfg, cfg.Postgres.DSN())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	defer db.Close()

	v := validator.NewValidator()
	catSvc := service.CatService()
	catHandler := handlers.NewCatHandler(catSvc, v, cfg.HTTP.HandlerTimeout)

	r := gin.Default()
	r.POST("/cats", catHandler.CreateCatHandler)
	// r.GET("/cats", catHandler.ListCatsHandler)

	httpServer := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      r,
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

	<-ctx.Done()
	log.Info().Msg("shutting down server gracefully")

	httpServer.SetKeepAlivesEnabled(false)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown server")
	}

	// if err = envBox.Close(); err != nil {
	// 	log.Error().Err(err).Msg("failed to close connections")
	// }

	log.Info().Msg("server gracefully shutdown")

	return nil
}
