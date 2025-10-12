package main

import (
	log "github.com/rs/zerolog/log"

	"github.com/DavydAbbasov/spy-cat/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Panic().Err(err).Msg("Application execution error")
	}
}
