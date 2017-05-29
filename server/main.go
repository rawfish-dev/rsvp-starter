package main

import (
	"github.com/rawfish-dev/rsvp-starter/server/api"
	"github.com/rawfish-dev/rsvp-starter/server/config"
)

func main() {
	loadedConfig := config.LoadConfig()

	reactReduxBasicsAPI := api.NewAPI(loadedConfig)

	reactReduxBasicsAPI.Run()
}
