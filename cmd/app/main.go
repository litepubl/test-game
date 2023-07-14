package main

import (
	"os"

	"github.com/litepubl/test-game/pkg/app"
	"github.com/rs/zerolog/log"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Error().Err(err).Msg("Application not created")
		os.Exit(1)
	}

	defer app.Close()
	err = app.Run()
	if err != nil {
		log.Error().Err(err).Msg("Error shutdown server")
	}
}
